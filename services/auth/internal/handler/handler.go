package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"streaming-platform/services/auth/internal/utils"
	db "streaming-platform/shared/db/generated"

	"streaming-platform/services/auth/internal/types"
)

type Handler struct {
	q db.Querier
}

func NewHandler(q db.Querier) *Handler {
	return &Handler{q: q}
}

func (h *Handler) handlerLogin(w http.ResponseWriter, r *http.Request) {

	var payload types.UserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.q.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return
	}

	if !utils.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	//jwt token generation
	token, err := utils.CreateJWT(u.UserID)
	if err != nil {
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) handlerRegister(w http.ResponseWriter, r *http.Request) {

	var payload types.UserPayload
	err := utils.ParseJSON(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//check in db if user exists
	user_exists, err := h.q.CheckUserExists(context.Background(), payload.Email)
	if user_exists {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	//hash password
	hashed_password, err := utils.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//Create new user
	_, err = h.q.CreateUser(context.Background(), db.CreateUserParams{
		Email:    payload.Email,
		Password: hashed_password,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handlerSecret(w http.ResponseWriter, r *http.Request) {

	utils.WriteJSON(w, http.StatusOK, "Hey")
}

func (h *Handler) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := fields[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.ID)
		next(w, r.WithContext(ctx))
	}
}

func (h *Handler) RegisterHandler(r *http.ServeMux) {

	r.HandleFunc("/login", h.handlerLogin)
	r.HandleFunc("/register", h.handlerRegister)

	//Secret test endpoint with auth middleware
	r.HandleFunc("/secret", h.withAuth(h.handlerSecret))

}
