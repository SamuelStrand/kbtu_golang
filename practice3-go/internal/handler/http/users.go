package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	usersUC "practice3/internal/usecase/users"
	"practice3/pkg/modules"
)

type userCreateReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type userPatchReq struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Age   *int    `json:"age"`
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	users, err := h.Users.GetUsers(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": users})
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		return
	}

	user, err := h.Users.GetUserByID(r.Context(), id)
	if err != nil {
		switch err {
		case usersUC.ErrBadRequest:
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		case usersUC.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": "user not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": user})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req userCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid json"})
		return
	}

	newID, err := h.Users.CreateUser(r.Context(), modules.User{Name: req.Name, Email: req.Email, Age: req.Age})
	if err != nil {
		switch err {
		case usersUC.ErrBadRequest, usersUC.ErrEmailFormat:
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"success": true, "data": map[string]any{"id": newID}})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		return
	}

	var req userCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid json"})
		return
	}

	err = h.Users.UpdateUser(r.Context(), id, modules.User{Name: req.Name, Email: req.Email, Age: req.Age})
	if err != nil {
		switch err {
		case usersUC.ErrBadRequest, usersUC.ErrEmailFormat:
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": err.Error()})
		case usersUC.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": "user not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "updated"})
}

func (h *Handler) PatchUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		return
	}

	var req userPatchReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid json"})
		return
	}

	existing, err := h.Users.GetUserByID(r.Context(), id)
	if err != nil {
		switch err {
		case usersUC.ErrBadRequest:
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		case usersUC.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": "user not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		}
		return
	}

	upd := modules.User{Name: existing.Name, Email: existing.Email, Age: existing.Age}
	if req.Name != nil {
		upd.Name = *req.Name
	}
	if req.Email != nil {
		upd.Email = *req.Email
	}
	if req.Age != nil {
		upd.Age = *req.Age
	}

	if err := h.Users.UpdateUser(r.Context(), id, upd); err != nil {
		switch err {
		case usersUC.ErrBadRequest, usersUC.ErrEmailFormat:
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": err.Error()})
		case usersUC.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": "user not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "message": "updated"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		return
	}

	affected, err := h.Users.DeleteUser(r.Context(), id)
	if err != nil {
		switch err {
		case usersUC.ErrBadRequest:
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": "invalid id"})
		case usersUC.ErrNotFound:
			writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "message": "user not found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"success": false, "message": err.Error()})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": map[string]any{"rows_affected": affected}})
}
