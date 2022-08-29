package model

type EventType string

const (
	CREATE EventType = "CREATE"
	UPDATE EventType = "UPDATE"
	DELETE EventType = "DELETE"
)

type PublishEvent struct {
	EventType   `json:"event_type"`
	Transaction `json:"message" binding:"required"`
}
