package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"practice5/internal/usecase"
)

type Handler struct {
	uc *usecase.UsersUsecase
}

func NewHandler(uc *usecase.UsersUsecase) *Handler {
	return &Handler{uc: uc}
}

// @Summary Health
// @Tags    health
// @Router  /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

// @Summary GetUsers
// @Tags    users
// @Param   limit      query int    false "limit"
// @Param   offset     query int    false "offset"
// @Param   order_by   query string false "order_by"
// @Param   order_dir  query string false "order_dir"
// @Param   id         query string false "id"
// @Param   name       query string false "name"
// @Param   email      query string false "email"
// @Param   gender     query string false "gender"
// @Param   birth_date query string false "birth_date"
// @Success 200 {object} object
// @Router  /users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit := parseInt(q.Get("limit"), 10)
	offset := parseInt(q.Get("offset"), 0)
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filters := usecase.UserFilters{
		ID:        q.Get("id"),
		Name:      q.Get("name"),
		Email:     q.Get("email"),
		Gender:    q.Get("gender"),
		BirthDate: q.Get("birth_date"),
	}

	orderBy := q.Get("order_by")
	orderDir := q.Get("order_dir")

	res, err := h.uc.GetPaginatedUsers(r.Context(), limit, offset, filters, orderBy, orderDir)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

// @Summary GetCommonFriends
// @Tags    users
// @Param   user1 query string true "user1"
// @Param   user2 query string true "user2"
// @Success 200 {object} object
// @Router  /users/common-friends [get]
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	u1 := strings.TrimSpace(q.Get("user1"))
	u2 := strings.TrimSpace(q.Get("user2"))
	if u1 == "" || u2 == "" {
		writeError(w, http.StatusBadRequest, "user1 and user2 are required")
		return
	}
	if !looksLikeUUID(u1) || !looksLikeUUID(u2) {
		writeError(w, http.StatusBadRequest, "user1 and user2 must look like UUID")
		return
	}

	friends, err := h.uc.GetCommonFriends(r.Context(), u1, u2)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": friends})
}

func looksLikeUUID(s string) bool {
	// Very small validation: 36 chars and 4 hyphens.
	if len(s) != 36 {
		return false
	}
	return s[8] == '-' && s[13] == '-' && s[18] == '-' && s[23] == '-'
}

func parseInt(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"success": false, "message": message})
}
