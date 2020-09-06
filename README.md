# Event Dispatcher Library
The Event Dispatcher library provides the necessary functionality that allows 
your application components to communicate with each other in a decoupled way,
by dispatching events and listening to them.

## Creating an Event Listener
The way to listen to an event is to register an event listener.

```
dispatcher := event.NewDispatcher()
dispatcher.On("event", func(ev interface{}) bool {
    if e, ok := ev.(*MyAppEvent); ok {
        // handle event
    }

    // return false to stop event propagation
    return true
})
```

## Dispatching an Event
In general, a single dispatcher, which holds all listeners, is needed.
When an event is dispatched, the dispatcher notifies all listeners registered with that event:

```
dispatcher := event.NewDispatcher()
dispatcher.On("event", func(ev interface{}) bool {
    if e, ok := ev.(*MyAppEvent); ok {
        // handle event
    }

    // return false to stop event propagation
    return true
})

ev := &MyAppEvent{}
dispatcher.Dispatch("event", ev)
```

## Example

```
package main

import (
	"fmt"

	"github.com/gvre/event"
)

type UserCreatedEvent struct {
	ID   int
	Name string
}

func main() {
	dispatcher := event.NewDispatcher()
	dispatcher.On("user.created", func(ev interface{}) bool {
		if e, ok := ev.(*UserCreatedEvent); ok {
			fmt.Println(e)
		}

		return true
	})

	ev := &UserCreatedEvent{
		ID:   1,
		Name: "user",
	}
	dispatcher.Dispatch("user.created", ev)
}
```

The above code in playground: https://play.golang.org/p/4ag-48koyop

## Naming Conventions
The event name can be any string, but optionally follows the following conventions:

- Use only lowercase letters, numbers, dots (.) and underscores (_);
- Prefix names with a namespace followed by a dot (e.g. `user.*`);
- End names with a verb that indicates what action has been taken (e.g. `user.created`).

## License
- MIT
