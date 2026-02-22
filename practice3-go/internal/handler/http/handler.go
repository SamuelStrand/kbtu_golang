package http

import (
	"net/http"

	usersUC "practice3/internal/usecase/users"
)

type Handler struct {
	Users *usersUC.Usecase
}

func New(users *usersUC.Usecase) *Handler {
	return &Handler{Users: users}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "ok"})
}
