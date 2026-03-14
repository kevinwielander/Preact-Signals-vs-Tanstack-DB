package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"backend/domain"
	"backend/events"
)

type AlarmHandler struct {
	store  *events.Store
	mu     sync.RWMutex
	alarms map[string]*domain.Alarm
}

func NewAlarmHandler(store *events.Store) *AlarmHandler {
	return &AlarmHandler{
		store:  store,
		alarms: make(map[string]*domain.Alarm),
	}
}

func (h *AlarmHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title             string   `json:"title"`
		Description       string   `json:"description"`
		Severity          string   `json:"severity"`
		AssignedResources []string `json:"assignedResources"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := events.GenerateID()
	now := time.Now().UTC()

	assigned := req.AssignedResources
	if assigned == nil {
		assigned = []string{}
	}

	evt, err := h.store.Append(id, domain.AggregateAlarm, domain.EventAlarmCreated, domain.AlarmCreatedData{
		Title:             req.Title,
		Description:       req.Description,
		Severity:          req.Severity,
		Status:            "active",
		AssignedResources: assigned,
	})
	if err != nil {
		http.Error(w, "failed to create alarm", http.StatusInternalServerError)
		return
	}

	alarm := &domain.Alarm{
		ID:                id,
		Title:             req.Title,
		Description:       req.Description,
		Severity:          req.Severity,
		Status:            "active",
		AssignedResources: assigned,
		EventHash:         evt.Hash,
		EventNumber:       evt.Version,
		CreatedOnOffset:   evt.Offset,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	h.mu.Lock()
	h.alarms[id] = alarm
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(alarm)
}

func (h *AlarmHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	h.mu.RLock()
	alarm, ok := h.alarms[id]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "alarm not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alarm)
}

func (h *AlarmHandler) List(w http.ResponseWriter, r *http.Request) {
	resourceID := r.URL.Query().Get("resourceId")

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	h.mu.RLock()
	alarms := make([]*domain.Alarm, 0, len(h.alarms))
	for _, alarm := range h.alarms {
		if resourceID != "" && !containsStr(alarm.AssignedResources, resourceID) {
			continue
		}
		alarms = append(alarms, alarm)
	}
	h.mu.RUnlock()

	slices.SortFunc(alarms, func(a, b *domain.Alarm) int {
		return b.CreatedOnOffset - a.CreatedOnOffset
	})

	total := len(alarms)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"items":    alarms[start:end],
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AlarmHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	alarm, ok := h.alarms[id]
	if !ok {
		h.mu.Unlock()
		http.Error(w, "alarm not found", http.StatusNotFound)
		return
	}

	for field, newValue := range patch {
		var oldValue any
		switch field {
		case "title":
			oldValue = alarm.Title
		case "description":
			oldValue = alarm.Description
		case "severity":
			oldValue = alarm.Severity
		case "status":
			oldValue = alarm.Status
		case "assignedResources":
			oldValue = alarm.AssignedResources
		default:
			continue
		}

		evt, _ := h.store.Append(id, domain.AggregateAlarm, domain.EventAlarmFieldUpdated, domain.FieldUpdatedData{
			Field:    field,
			OldValue: oldValue,
			NewValue: newValue,
		})

		switch field {
		case "title":
			alarm.Title, _ = newValue.(string)
		case "description":
			alarm.Description, _ = newValue.(string)
		case "severity":
			alarm.Severity, _ = newValue.(string)
		case "status":
			alarm.Status, _ = newValue.(string)
		case "assignedResources":
			if v, ok := newValue.([]any); ok {
				alarm.AssignedResources = make([]string, len(v))
				for i, item := range v {
					alarm.AssignedResources[i], _ = item.(string)
				}
			}
		}

		alarm.EventHash = evt.Hash
		alarm.EventNumber = evt.Version
		alarm.UpdatedAt = evt.Timestamp
	}
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alarm)
}

func (h *AlarmHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	evts := h.store.GetEvents(id)
	if len(evts) == 0 {
		http.Error(w, "alarm not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evts)
}

// ListAll returns all alarms — used by ResourceHandler for cross-queries.
func (h *AlarmHandler) ListAll() []*domain.Alarm {
	h.mu.RLock()
	defer h.mu.RUnlock()

	alarms := make([]*domain.Alarm, 0, len(h.alarms))
	for _, a := range h.alarms {
		alarms = append(alarms, a)
	}
	return alarms
}

func containsStr(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
