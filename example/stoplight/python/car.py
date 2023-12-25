# AUTOMATICALLY GENERATED by ../..//statecraft car.statecraft python/car.py
# DO NOT EDIT
# Original .statecraft file contents included at bottom.

# States
class State:
    def __init__(self, name, label):
        self.name = name
        self.label = label
        
    def __str__(self):
        return self.name

    def label(self):
        return self.label
Stopped = State("Stopped", "Stopped at red light")
Deciding = State("Deciding", "Deciding whether to stop")
Going = State("Going", "Going through light")
Beyond = State("Beyond", "Beyond light already")

# Events
class Event:
    def __init__(self, name):
        self.name = name

    def __str__(self):
        return self.name
Green = Event("Green")
Stop = Event("Stop")
Go = Event("Go")
Red = Event("Red")
Yellow = Event("Yellow")

# Transition
class Transition:
    def __init__(self, src, event, method, dst):
        self.src = src
        self.event = event
        self.method = method
        self.dst = dst

# Machine
class Machine:

    def __init__(self, handlers, state):
        self.state = state
        self.g = {}
        self.g["Stopped"] = {}
        self.g["Stopped"]["Go"] = Transition(
            Stopped,
            Go,
            handlers.Gas,
            Going
        )
        self.g["Stopped"]["Green"] = Transition(
            Stopped,
            Green,
            handlers.Gas,
            Going
        )
        self.g["Stopped"]["Red"] = Transition(
            Stopped,
            Red,
            handlers.Brake,
            Stopped
        )
        self.g["Stopped"]["Stop"] = Transition(
            Stopped,
            Stop,
            handlers.Brake,
            Stopped
        )
        self.g["Stopped"]["Yellow"] = Transition(
            Stopped,
            Yellow,
            handlers.Decide,
            Deciding
        )
        self.g["Deciding"] = {}
        self.g["Deciding"]["Go"] = Transition(
            Deciding,
            Go,
            handlers.Gas,
            Going
        )
        self.g["Deciding"]["Green"] = Transition(
            Deciding,
            Green,
            handlers.Gas,
            Going
        )
        self.g["Deciding"]["Red"] = Transition(
            Deciding,
            Red,
            handlers.Brake,
            Stopped
        )
        self.g["Deciding"]["Stop"] = Transition(
            Deciding,
            Stop,
            handlers.Brake,
            Stopped
        )
        self.g["Deciding"]["Yellow"] = Transition(
            Deciding,
            Yellow,
            handlers.Decide,
            Deciding
        )
        self.g["Going"] = {}
        self.g["Going"]["Go"] = Transition(
            Going,
            Go,
            handlers.Gas,
            Going
        )
        self.g["Going"]["Green"] = Transition(
            Going,
            Green,
            handlers.Gas,
            Going
        )
        self.g["Going"]["Red"] = Transition(
            Going,
            Red,
            handlers.Gas,
            Beyond
        )
        self.g["Going"]["Stop"] = Transition(
            Going,
            Stop,
            handlers.Brake,
            Stopped
        )
        self.g["Going"]["Yellow"] = Transition(
            Going,
            Yellow,
            handlers.Decide,
            Deciding
        )
        self.g["Beyond"] = {}
        self.g["Beyond"]["Go"] = Transition(
            Beyond,
            Go,
            handlers.Gas,
            Going
        )
        self.g["Beyond"]["Green"] = Transition(
            Beyond,
            Green,
            handlers.Gas,
            Going
        )
        self.g["Beyond"]["Red"] = Transition(
            Beyond,
            Red,
            handlers.Brake,
            Stopped
        )
        self.g["Beyond"]["Stop"] = Transition(
            Beyond,
            Stop,
            handlers.Brake,
            Stopped
        )
        self.g["Beyond"]["Yellow"] = Transition(
            Beyond,
            Yellow,
            handlers.Decide,
            Deciding
        )

    def tick(self, event):
        if event.name not in self.g[self.state.name]:
            raise Exception(f'unhandled: state {self.state} event {event}')
        t = self.g[self.state.name][event.name]
        self.state = t.dst
        if t.method is not None:
            t.method()
        return self.state

'''
// Comments look like this.  We ignore blank lines.

// Case matters in this file -- the names you provide here are passed
// straight through to the Go code generator as variable and struct
// names. This means that if you generate the .go file in a
// subdirectory as a separate package from your calling code, you'll
// need to uppercase everything here so it will be exported.
// 
// The README and example code assumes that you will be generating the
// .dot and .go files in a subdirectory as a separate package from
// your calling code, so we uppercase everything in this example.
//
// You can instead choose to generate your .go in the same directory
// as the calling code, in which case everything in your .statecraft
// file can be lowercased.
//
// Another consideration is reserved words -- several are used in the
// generated go code, including 'State', 'Event', 'New', 'Machine',
// and 'Tick'.  You will see compiler errors if you use any of these
// words for your own state or event names.

// Declare package name -- this is used verbatim as the 'package' name
// at the top of the generated .go:
package car

// Declare states with an 's' followed by the state node description.
// The first word of the description is used as the state node name.
// Each state name must be unique.

s Stopped at red light
s Deciding whether to stop
s Going through light 
s Beyond light already

// Declare transitions with a 't' followed by the source state, event
// name, and destination state.  Declare an optional transition method
// name as part of the event name, after a slash.
// 
// Regular expressions can be used as wildcards in the source state
// field.  The first matching rule will be used.

t Going Green/Gas Going

t Deciding Stop/Brake Stopped 
t Deciding Go/Gas Going 
t Going Red/Gas Beyond

t .* Red/Brake Stopped 
t .* Yellow/Decide Deciding 
t .* Green/Gas Going
t .* Stop/Brake Stopped
t .* Go/Gas Going

'''

