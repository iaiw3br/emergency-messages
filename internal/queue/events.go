package queue

type EventType string

const (
	EventTypeSend      = "send"
	EventTypeDelivered = "delivered"
	EventTypeFailed    = "failed"
)

func getEvents() []string {
	return []string{
		EventTypeSend,
		EventTypeDelivered,
		EventTypeFailed,
	}
}
