package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"practice2/internal/store"
	"practice2/internal/utils"
)

type TaskHandler struct {
	store *store.TaskStore
}

func NewTaskHandler(s *store.TaskStore) http.Handler {
	return &TaskHandler{store: s}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	idStr := strings.TrimSpace(q.Get("id"))
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			utils.WriteError(w, http.StatusBadRequest, "invalid id")
			return
		}

		task, ok := h.store.Get(id)
		if !ok {
			utils.WriteError(w, http.StatusNotFound, "task not found")
			return
		}

		utils.WriteJSON(w, http.StatusOK, task)
		return
	}

	var doneFilter *bool
	doneStr := strings.TrimSpace(q.Get("done"))
	if doneStr != "" {
		v, err := parseBoolStrict(doneStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid done")
			return
		}
		doneFilter = &v
	}

	tasks := h.store.List(doneFilter)
	utils.WriteJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req struct {
		Title string `json:"title"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid title")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid title")
		return
	}

	task := h.store.Create(title)
	utils.WriteJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(r.URL.Query().Get("id"))
	if idStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req struct {
		Done *bool `json:"done"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil || req.Done == nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid done")
		return
	}

	updatedTask, ok := h.store.UpdateDone(id, *req.Done)
	if !ok {
		utils.WriteError(w, http.StatusNotFound, "task not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedTask)
}

func parseBoolStrict(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, strconv.ErrSyntax
	}
}
