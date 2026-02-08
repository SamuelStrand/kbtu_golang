package store

import (
	"practice2/internal/models"
	"sync"
)

type TaskStore struct {
	mu     sync.RWMutex
	nextID int
	tasks  map[int]models.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		nextID: 1,
		tasks:  make(map[int]models.Task),
	}
}

func (s *TaskStore) List(doneFilter *bool) []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if doneFilter != nil && t.Done != *doneFilter {
			continue
		}
		out = append(out, t)
	}
	return out
}

func (s *TaskStore) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	return t, ok
}

func (s *TaskStore) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[t.ID] = t
	s.nextID++
	return t
}

func (s *TaskStore) UpdateDone(id int, done bool) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return models.Task{}, false
	}
	t.Done = done
	s.tasks[id] = t
	return t, true
}
