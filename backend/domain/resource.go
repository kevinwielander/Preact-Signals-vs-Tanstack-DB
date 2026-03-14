package domain

import "time"

const (
	EventResourceCreated      = "ResourceCreated"
	EventResourceFieldUpdated = "ResourceFieldUpdated"
	AggregateResource         = "resource"
)

type Resource struct {
	ID               string    `json:"id"`
	DisplayName      string    `json:"displayName"`
	Email            string    `json:"email"`
	IsUserAssociated bool      `json:"isUserAssociated"`
	Thumbnail        string    `json:"thumbnail"`
	EventHash        string    `json:"eventHash"`
	EventNumber      int       `json:"eventNumber"`
	ArchivedOnOffset *int      `json:"archivedOnOffset"`
	CreatedOnOffset  int       `json:"createdOnOffset"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type ResourceCreatedData struct {
	DisplayName      string `json:"displayName"`
	Email            string `json:"email"`
	IsUserAssociated bool   `json:"isUserAssociated"`
	Thumbnail        string `json:"thumbnail"`
}
