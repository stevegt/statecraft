package main

import (
	"math/rand"
	"time"

	. "github.com/stevegt/goadapt"
	c "github.com/stevegt/statecraft/example/stoplight/car"
	"github.com/stevegt/statecraft/example/stoplight/stoplight"
	s "github.com/stevegt/statecraft/example/stoplight/stoplight"
)

//go:generate ../../statecraft stoplight/stoplight.statecraft stoplight/stoplight.go
//go:generate ../../statecraft stoplight/stoplight.statecraft stoplight/stoplight.dot
//go:generate ../../statecraft car/car.statecraft car/car.go
//go:generate ../../statecraft car/car.statecraft car/car.dot

func main() {

	// Start a goroutine that creates a stoplight state machine that
	// emits light change events on an events channel.
	//
	// It might help to ignore the fact that we have two state
	// machines in this code -- the car state machine is the important
	// one.  The stoplight state machine, the events channel, and the
	// light() function are just here to provide a source of external
	// events for car to react to -- see the event loop below.
	//
	// In your own application, you might instead get events from a
	// sensor or io.Reader, or you might write an io.Writer method
	// somewhere that directly calls the Tick() method in your state
	// machine's generated code.
	ssm, events := light()

	// Create an instance of a Handlers-compliant struct.  See the
	// Handlers interface in car/car.go and the comments below.
	handlers := &Car{events: events}

	// Create a car state machine, passing it the handlers and
	// initial state.
	car := c.New(handlers, c.Stopped)

	// Main event loop:  Run forever, passing events from stoplight to
	// car.Tick().
	for event := range events {
		Pf("light is %-8v ", ssm.State)
		Pf("event %-8v ", event)
		// Send event to car, get back new state.  However you get
		// events, you simply need to call your generated state
		// machine's Tick() method once for each event.
		state, err := car.Tick(event)
		Ck(err)
		Pf("%-15s ", handlers.action)
		Pl("car is", state)
	}
}

func light() (ssm *stoplight.Machine, events chan c.Event) {
	events = make(chan c.Event, 10)
	// create a stoplight state machine
	ssm = s.New(nil, stoplight.Red)
	go func() {
		for {
			// send timer event to the stoplight, getting back new
			// light state
			lightState, err := ssm.Tick(s.Timer)
			Ck(err)
			// send the new stoplight state to the car as
			// a car state machine input event, then
			// sleep until next timer event
			switch lightState {
			case s.Green:
				events <- c.Green
				time.Sleep(5 * time.Second)
			case s.Yellow:
				events <- c.Yellow
				time.Sleep(2 * time.Second)
			case s.Red:
				events <- c.Red
				time.Sleep(7 * time.Second)
			}
		}
	}()
	return
}

// Car is a set of event handlers that the car state machine will call
// -- see the car.Handlers interface in car.go.
//
// In this example, the events field is a channel we will use to send
// stoplight change events to the car state machine from light(), to
// be processed by the event loop in main().  Whether you need some
// sort of events channel like this depends completely on your
// application and event loop -- see the comments in main().
//
// The action field is just a place for handlers to store a text
// description of a transition, and isn't strictly needed for
// functionality.
//
// In other words, your application may not need either of these
// fields, or may need others -- it's up to you.  An empty struct with
// just the Handler interface methods may be enough.
type Car struct {
	events chan c.Event
	action string
}

func (handlers *Car) Brake() {
	handlers.action = "applying brake"
}

func (handlers *Car) Gas() {
	handlers.action = "applying gas"
}

func (handlers *Car) Decide() {
	if rand.Float64() < rand.Float64() {
		handlers.events <- c.Go
	} else {
		handlers.events <- c.Stop
	}
}
