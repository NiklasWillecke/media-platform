package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"streaming-platform/services/auth/pkg/utils"
	"streaming-platform/services/upload/pkg/queue"
	dataStore "streaming-platform/services/upload/pkg/s3"
)

type Handler struct {
	*dataStore.S3Service
	*queue.Queue
}

func NewHandler(q *queue.Queue, d *dataStore.S3Service) *Handler {
	return &Handler{
		S3Service: d,
		Queue:     q,
	}
}

type presigend_request struct {
	name string
	size string
}

func (h *Handler) handlerGetPresigendUrl(w http.ResponseWriter, r *http.Request) {

	//Request with size and name
	var p presigend_request

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "userID not found", http.StatusUnauthorized)
		return
	}

	url := h.S3Service.CreatePresignedUrl(userID, p.name)

	utils.WriteJSON(w, http.StatusOK, map[string]string{"Url": url})

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

func (h *Handler) RegisterHandler(r *http.ServeMux) {
	r.HandleFunc("/register", WithAuth(h.handlerGetPresigendUrl))
}
