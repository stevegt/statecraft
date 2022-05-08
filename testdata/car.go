package car

import "fmt"

// AUTOMATICALLY GENERATED by test go
// DO NOT EDIT
// Original .statecraft file contents included at bottom.

type Handlers interface {
    Gas()
    Brake()
    Decide()
}

// states
type State string
const (
    Stopped State = "Stopped" // Stopped at red light
    Deciding State = "Deciding" // Deciding whether to stop
    Going State = "Going" // Going through light
    Beyond State = "Beyond" // Beyond light already
)

// events
type Event string
const (
    
    Green Event = "Green"
    Stop Event = "Stop"
    Go Event = "Go"
    Red Event = "Red"
    Yellow Event = "Yellow"
)

type Transition struct {
    Src    State
    Event  Event
	Method func()
	Dst    State
}

type Transitions map[Event]Transition

type Graph map[State]Transitions

type Machine struct {
    g Graph
    State State
}

func New(handlers Handlers, initState State) (m *Machine) {
    m = &Machine{
        State: initState,
    }
    m.g = Graph{
        State("Stopped"):  Transitions{
            Event("Go"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Green"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Red"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Stop"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Yellow"): Transition{
                    Dst: State("Deciding"), 
                    Method: handlers.Decide, 
            },
        },
        State("Deciding"):  Transitions{
            Event("Go"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Green"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Red"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Stop"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Yellow"): Transition{
                    Dst: State("Deciding"), 
                    Method: handlers.Decide, 
            },
        },
        State("Going"):  Transitions{
            Event("Go"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Green"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Red"): Transition{
                    Dst: State("Beyond"), 
                    Method: handlers.Gas, 
            },
            Event("Stop"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Yellow"): Transition{
                    Dst: State("Deciding"), 
                    Method: handlers.Decide, 
            },
        },
        State("Beyond"):  Transitions{
            Event("Go"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Green"): Transition{
                    Dst: State("Going"), 
                    Method: handlers.Gas, 
            },
            Event("Red"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Stop"): Transition{
                    Dst: State("Stopped"), 
                    Method: handlers.Brake, 
            },
            Event("Yellow"): Transition{
                    Dst: State("Deciding"), 
                    Method: handlers.Decide, 
            },
        },
    }
    return
}

func (m *Machine) Tick(event Event) (newState State, err error) {
    src := m.g[m.State]
    t, ok := src[event]
    if !ok {
        err = fmt.Errorf("unhandled: state %s event %s", string(m.State), string(event))
        return
    }
    m.State = t.Dst
    if t.Method != nil {
        t.Method()
    }
    return m.State, nil
}

var txt string = `
// Comments look like this.  We ignore blank lines.

// Declare Go package and state machine name.

package car
machine Car

// Declare states with an 's' followed by the state node description.
// - the first word of the description is used as the state node name 
// - the state name must be unique

s Stopped at red light
s Deciding whether to stop
s Going through light 
s Beyond light already

// Declare transitions with a 't' followed by the source state, event
// name, and destination state.  Declare an optional transition method
// name as part of the event name, after a slash.
// Regular expressions can be used as wildcards in the source name.
// The first matching rule will be used.

t Going Green/Gas Going

t Deciding Stop/Brake Stopped 
t Deciding Go/Gas Going 
t Going Red/Gas Beyond

t .* Red/Brake Stopped 
t .* Yellow/Decide Deciding 
t .* Green/Gas Going
t .* Stop/Brake Stopped
t .* Go/Gas Going


` 
