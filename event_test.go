package event_test

import (
	"testing"

	"github.com/gvre/event"
)

type UserCreatedEvent struct {
	ID   int
	Name string
}

func TestUserCreatedEvent(t *testing.T) {
	id := 1
	name := "user"

	dispatcher := event.NewDispatcher()
	dispatcher.On("user.created", func(ev interface{}) bool {
		if e, ok := ev.(*UserCreatedEvent); ok {
			if e.ID != id {
				t.Errorf("Invalid user id, got: %d, want: %d", e.ID, id)
			}

			if e.Name != name {
				t.Errorf("Invalid user name, got: %q, want: %q", e.Name, name)
			}
		}

		return true
	})

	ev := &UserCreatedEvent{
		ID:   id,
		Name: name,
	}
	dispatcher.Dispatch("user.created", ev)
}

func TestStopPropagation(t *testing.T) {
	var totalExecutions int

	dispatcher := event.NewDispatcher()
	dispatcher.On("user.created", func(ev interface{}) bool {
		totalExecutions++
		return false
	})
	dispatcher.On("user.created", func(ev interface{}) bool {
		totalExecutions++
		return true
	})

	ev := &UserCreatedEvent{}
	dispatcher.Dispatch("user.created", ev)

	expectedExecutions := 1
	if totalExecutions != expectedExecutions {
		t.Errorf("Invalid number of executions, got: %d, want: %d", totalExecutions, expectedExecutions)
	}
}

func TestWildcardListener(t *testing.T) {
	var totalExecutions int

	dispatcher := event.NewDispatcher()
	dispatcher.On("user.*", func(ev interface{}) bool {
		totalExecutions++
		return true
	})

	ev := &UserCreatedEvent{}
	dispatcher.Dispatch("user.created", ev)

	expectedExecutions := 1
	if totalExecutions != expectedExecutions {
		t.Errorf("Invalid number of executions, got: %d, want: %d", totalExecutions, expectedExecutions)
	}
}
