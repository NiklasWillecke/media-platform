package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NiklasWillecke/media-platform/services/auth/internal/types"
	"github.com/NiklasWillecke/media-platform/services/auth/internal/utils"
	db "github.com/niklaswillecke/streaming-platform/shared/db/generated"
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

func (h *Handler) RegisterHandler(r *http.ServeMux) {
	fmt.Println("Hello")

	r.HandleFunc("/login", h.handlerLogin)
	r.HandleFunc("/register", h.handlerRegister)

}
