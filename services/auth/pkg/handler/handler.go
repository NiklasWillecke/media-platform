package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"streaming-platform/services/auth/pkg/utils"
	db "streaming-platform/shared/db/generated"

	"streaming-platform/services/auth/pkg/types"
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

	cookie := http.Cookie{
		Name:     "sp_auth",
		Value:    token,
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   int((time.Hour * 24).Seconds()), // z.B. 1 Tag
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,                // in Dev oft false, in Prod true (HTTPS)
		SameSite: http.SameSiteLaxMode, // damit Browser Cookie cross-site mitsendet
		// Domain: optional in Prod, z.B. ".example.com"; für localhost weglassen
	}
	http.SetCookie(w, &cookie)

	//jwt token generation
	token, err = utils.CreateJWT(u.UserID)
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
	u, err := h.q.CreateUser(context.Background(), db.CreateUserParams{
		Email:    payload.Email,
		Password: hashed_password,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//jwt token generation
	token, err := utils.CreateJWT(u.UserID)
	if err != nil {
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) handlerSecret(w http.ResponseWriter, r *http.Request) {

	utils.WriteJSON(w, http.StatusOK, "Hey")
}

func WithAuth1(next http.HandlerFunc) http.HandlerFunc {
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

func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sp_auth")

		if err != nil || cookie == nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims, err := utils.ValidateJWT(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.ID)
		next(w, r.WithContext(ctx))
	}
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // explizit, kein *
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true") // Pflicht für Cookies

		// OPTIONS-Preflight direkt beantworten
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RegisterHandler(r *http.ServeMux) {

	r.HandleFunc("/login", h.handlerLogin)
	r.HandleFunc("/register", h.handlerRegister)

	//Secret test endpoint with auth middleware
	r.HandleFunc("/secret", WithAuth(h.handlerSecret))
}
