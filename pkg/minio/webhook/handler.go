package webhook

import "context"

// Handler interface for webhook events
type Handler interface {
	ProcessEvent(ctx context.Context, e *Event) error
}
