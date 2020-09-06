// Package event provides the necessary functionality for subscribing to events
// and get notified when they happen.
package event

import (
	"strings"
	"sync"
)

const namespaceSeparator string = "."

type callback func(interface{}) (propagate bool)

// Dispatcher holds all events.
type Dispatcher struct {
	mu     sync.RWMutex
	events map[string][]callback
}

// NewDispatcher returns a new event dispatcher.
func NewDispatcher() *Dispatcher {
	d := &Dispatcher{}
	d.events = make(map[string][]callback)

	return d
}

// On is used when an callback wants to subscribe to an event.
// Wildcard events are supported when a wildcard
// is the suffix of the name. For example:
// "event.*" is supported
// "event.*.name is not supported
func (d *Dispatcher) On(name string, fn callback) {
	name = strings.TrimSuffix(name, "*")

	d.mu.Lock()
	d.events[name] = append(d.events[name], fn)
	d.mu.Unlock()
}

// Dispatch passes the Event data to all listeners that have subscribed
// to the event the name parameter equals to. Wildcard listeners of the parent namespace
// will also receive the event. For example, a dispatched event with name "event.test"
// will be received by the listeners of "event.test" and "event.*".
// The "." character is removed from the end of the event name,
// so a name like "event." becomes "event".
func (d *Dispatcher) Dispatch(name string, ev interface{}) {
	name = strings.TrimSuffix(name, namespaceSeparator)

	d.mu.RLock()
	defer d.mu.RUnlock()

	// Check wildcard events.
	var events []callback
	if pos := strings.LastIndex(name, namespaceSeparator); pos > 0 {
		prefix := name[:pos+1]
		if _events, ok := d.events[prefix]; ok {
			events = append(events, _events...)
		}
	}

	// Check events whose name matches with the dispatched event name.
	if _events, ok := d.events[name]; ok {
		events = append(events, _events...)
	}

	if len(events) == 0 {
		return
	}

	for _, fn := range events {
		if !fn(ev) {
			// stop propagation
			break
		}
	}
}
