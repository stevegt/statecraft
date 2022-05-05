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

	// start a goroutine that creates a stoplight state machine and
	// emits light change events
	ssm, events := light()

	// create a Car system that contains the callbacks for the car
	// state machine to call
	csys := &Car{Events: events}

	// create a car state machine
	csm := c.New(csys, c.Stopped)

	// run forever, processing events from stoplight in callbacks
	for event := range events {
		Pf("light is %-8v ", ssm.State)
		Pf("event %-8v ", event)
		state, err := csm.Tick(event)
		Ck(err)
		Pf("%-15s ", csys.action)
		Pl("car is", state)
	}
}

func light() (ssm *stoplight.Machine, events chan c.Event) {
	events = make(chan c.Event, 10)
	// create a stoplight state machine
	ssm = s.New(nil, stoplight.Red)
	go func() {
		for {
			lightState, err := ssm.Tick(s.Timer)
			Ck(err)
			events <- c.Event(lightState)
			switch lightState {
			case s.Green:
				time.Sleep(5 * time.Second)
			case s.Yellow:
				time.Sleep(2 * time.Second)
			case s.Red:
				time.Sleep(7 * time.Second)
			}
		}
	}()
	return
}

type Car struct {
	Events chan c.Event
	action string
}

func (csys *Car) Brake() {
	csys.action = "applying brake"
}

func (csys *Car) Gas() {
	csys.action = "applying gas"
}

func (csys *Car) Decide() {
	if rand.Float64() < rand.Float64() {
		csys.Events <- c.Event("Go")
	} else {
		csys.Events <- c.Event("Stop")
	}
}
