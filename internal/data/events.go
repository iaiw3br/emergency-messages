package data

type EventType string

const (
	EventTypeSend = "send"
)

func getEvents() []string {
	return []string{
		EventTypeSend,
	}
}
