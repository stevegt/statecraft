# statecraft
State machine compiler that generates Go code and graphviz dot files from
a simple DSL.

See ./example/stoplight for a demo of two interacting state machines.
The `car` DSL from that example looks like this:

```
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

```
