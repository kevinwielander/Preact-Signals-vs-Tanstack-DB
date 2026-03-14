package domain

import "time"

const (
	EventAlarmCreated      = "AlarmCreated"
	EventAlarmFieldUpdated = "AlarmFieldUpdated"
	AggregateAlarm         = "alarm"
)

type Alarm struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Severity          string    `json:"severity"`
	Status            string    `json:"status"`
	AssignedResources []string  `json:"assignedResources"`
	EventHash         string    `json:"eventHash"`
	EventNumber       int       `json:"eventNumber"`
	ArchivedOnOffset  *int      `json:"archivedOnOffset"`
	CreatedOnOffset   int       `json:"createdOnOffset"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type AlarmCreatedData struct {
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Severity          string   `json:"severity"`
	Status            string   `json:"status"`
	AssignedResources []string `json:"assignedResources"`
}

type FieldUpdatedData struct {
	Field     string `json:"field"`
	OldValue  any    `json:"oldValue"`
	NewValue  any    `json:"newValue"`
	ChangedBy string `json:"changedBy"`
}
