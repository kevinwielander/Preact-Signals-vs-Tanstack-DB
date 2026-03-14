package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"backend/domain"
	"backend/events"
)

type ResourceHandler struct {
	store      *events.Store
	mu         sync.RWMutex
	resources  map[string]*domain.Resource
	listAlarms func() []*domain.Alarm
	meID       string
}

func NewResourceHandler(store *events.Store, listAlarms func() []*domain.Alarm) *ResourceHandler {
	return &ResourceHandler{
		store:      store,
		resources:  make(map[string]*domain.Resource),
		listAlarms: listAlarms,
	}
}

func (h *ResourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DisplayName      string `json:"displayName"`
		Email            string `json:"email"`
		IsUserAssociated bool   `json:"isUserAssociated"`
		Thumbnail        string `json:"thumbnail"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := events.GenerateID()
	now := time.Now().UTC()

	evt, err := h.store.Append(id, domain.AggregateResource, domain.EventResourceCreated, domain.ResourceCreatedData{
		DisplayName:      req.DisplayName,
		Email:            req.Email,
		IsUserAssociated: req.IsUserAssociated,
		Thumbnail:        req.Thumbnail,
	})
	if err != nil {
		http.Error(w, "failed to create resource", http.StatusInternalServerError)
		return
	}

	resource := &domain.Resource{
		ID:               id,
		DisplayName:      req.DisplayName,
		Email:            req.Email,
		IsUserAssociated: req.IsUserAssociated,
		Thumbnail:        req.Thumbnail,
		EventHash:        evt.Hash,
		EventNumber:      evt.Version,
		CreatedOnOffset:  evt.Offset,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	h.mu.Lock()
	if h.meID == "" {
		h.meID = id
	}
	h.resources[id] = resource
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resource)
}

func (h *ResourceHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	h.mu.RLock()
	resource, ok := h.resources[id]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (h *ResourceHandler) List(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	resources := make([]*domain.Resource, 0, len(h.resources))
	for _, res := range h.resources {
		resources = append(resources, res)
	}
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

func (h *ResourceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var patch map[string]any
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	resource, ok := h.resources[id]
	if !ok {
		h.mu.Unlock()
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	for field, newValue := range patch {
		var oldValue any
		switch field {
		case "displayName":
			oldValue = resource.DisplayName
		case "email":
			oldValue = resource.Email
		case "isUserAssociated":
			oldValue = resource.IsUserAssociated
		case "thumbnail":
			oldValue = resource.Thumbnail
		default:
			continue
		}

		evt, _ := h.store.Append(id, domain.AggregateResource, domain.EventResourceFieldUpdated, domain.FieldUpdatedData{
			Field:    field,
			OldValue: oldValue,
			NewValue: newValue,
		})

		switch field {
		case "displayName":
			resource.DisplayName, _ = newValue.(string)
		case "email":
			resource.Email, _ = newValue.(string)
		case "isUserAssociated":
			resource.IsUserAssociated, _ = newValue.(bool)
		case "thumbnail":
			resource.Thumbnail, _ = newValue.(string)
		}

		resource.EventHash = evt.Hash
		resource.EventNumber = evt.Version
		resource.UpdatedAt = evt.Timestamp
	}
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resource)
}

func (h *ResourceHandler) Me(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	res, ok := h.resources[h.meID]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "no resources exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *ResourceHandler) GetAlarms(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	h.mu.RLock()
	_, ok := h.resources[id]
	h.mu.RUnlock()

	if !ok {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	all := h.listAlarms()
	alarms := make([]*domain.Alarm, 0)
	for _, alarm := range all {
		if containsStr(alarm.AssignedResources, id) {
			alarms = append(alarms, alarm)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alarms)
}

func (h *ResourceHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	evts := h.store.GetEvents(id)
	if len(evts) == 0 {
		http.Error(w, "resource not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(evts)
}
