package events

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Event struct {
	ID            string          `json:"id"`
	AggregateID   string          `json:"aggregateId"`
	AggregateType string          `json:"aggregateType"`
	EventType     string          `json:"eventType"`
	Data          json.RawMessage `json:"data"`
	Timestamp     time.Time       `json:"timestamp"`
	Version       int             `json:"version"`
	Offset        int             `json:"offset"`
	Hash          string          `json:"hash"`
}

type Store struct {
	mu      sync.RWMutex
	events  map[string][]Event
	offset  atomic.Int64
	onEvent func(Event)
}

func NewStore(onEvent func(Event)) *Store {
	return &Store{
		events:  make(map[string][]Event),
		onEvent: onEvent,
	}
}

func (s *Store) Append(aggregateID, aggregateType, eventType string, data any) (Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rawData, err := json.Marshal(data)
	if err != nil {
		return Event{}, err
	}

	version := len(s.events[aggregateID]) + 1
	offset := int(s.offset.Add(1))

	evt := Event{
		ID:            GenerateID(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		Data:          rawData,
		Timestamp:     time.Now().UTC(),
		Version:       version,
		Offset:        offset,
	}

	evt.Hash = computeHash(evt)
	s.events[aggregateID] = append(s.events[aggregateID], evt)

	if s.onEvent != nil {
		s.onEvent(evt)
	}

	return evt, nil
}

func (s *Store) GetEvents(aggregateID string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := s.events[aggregateID]
	result := make([]Event, len(events))
	copy(result, events)
	return result
}

func computeHash(evt Event) string {
	h := sha256.New()
	h.Write([]byte(evt.ID))
	h.Write([]byte(evt.AggregateID))
	h.Write(evt.Data)
	h.Write([]byte(evt.Timestamp.String()))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
