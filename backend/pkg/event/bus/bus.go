package bus

import (
	"context"
	"sync"
)

// Event represents a generic event.
type Event interface {
	Name() string
}

// Listener handles an event.
type Listener interface {
	Handle(ctx context.Context, event Event) error
}

// ListenerFunc is a function adapter for Listener.
type ListenerFunc func(ctx context.Context, event Event) error

func (f ListenerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// Bus is a simple in-memory event bus.
type Bus struct {
	listeners map[string][]Listener
	mu        sync.RWMutex
}

// New creates a new Event Bus.
func New() *Bus {
	return &Bus{
		listeners: make(map[string][]Listener),
	}
}

// Subscribe subscribes a listener to an event name.
func (b *Bus) Subscribe(eventName string, listener Listener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners[eventName] = append(b.listeners[eventName], listener)
}

// Publish publishes an event to all subscribers.
// It executes listeners synchronously for simplicity in this iteration,
// but can be extended to be async.
func (b *Bus) Publish(ctx context.Context, event Event) error {
	b.mu.RLock()
	listeners := b.listeners[event.Name()]
	b.mu.RUnlock()

	for _, listener := range listeners {
		if err := listener.Handle(ctx, event); err != nil {
			// For now, we just log the error or return it.
			// In a real system, we might want to continue executing other listeners
			// and aggregate errors.
			return err
		}
	}
	return nil
}
