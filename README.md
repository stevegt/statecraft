# Statecraft

Statecraft is a state machine compiler that generates graphviz dot
files and Go or Python code from a simple DSL.

# Installation

Install Go 1.17 or later -- to manage upgrades and multiple Go
versions on the same machine, I use https://github.com/syndbg/goenv/.

Install the `statecraft` binary -- rather than `go get`, recent Go
versions use `go install` to install binaries:


```
go install github.com/stevegt/statecraft@latest
```

# Quick Start 

1. Install; see above.
1. Write a `foo.statecraft` file that describes the state machine
   you want to generate.  See below and
   `./example/stoplight/car/car.statecraft` for the DSL syntax.
2. Run `statecraft foo.statecraft foo.dot` to generate a
   graphviz dot file.
3. Fix any errors thrown by the `statecraft` run.  As of this writing,
   the most likely errors will be cases where you need to add DSL `t`
   (transition) statements to handle events in states where you didn't
   expect them.  The easy way to handle these is with a wildcard
   source state -- see the bottom of
   `./example/stoplight/car.statecraft`.
4. Use `xdot` or your favorite graphviz viewer to visually inspect the
   dot file for the state machine you've created.  Fix bugs.

# Quick Start for Generating Go Code

1. Do the Quick Start steps above.
2. In the base directory of your project, run `go mod init
   your/repo/uri` if you haven't already.  
3. Run `mkdir foo`, where `foo` is the name of the state machine
   you're creating.
2. Run `statecraft foo.statecraft foo/foo.go` to generate the Go
   code for your state machine.
1. In your calling code, create a struct or other custom type with
   methods that satisfy the foo.Handlers interface you'll see in your
   generated foo/foo.go.  These are the handlers for the events you
   specified in `foo.statecraft`.  See the `Car` struct in
   `./example/stoplight/go/main.go`.
1. In your calling code, `import your/repo/uri/foo`, write an event
   loop of some sort, and in the loop, call foo.Tick() for each event.
   See the comments in `./example/stoplight/go/main.go`.
2. Optionally add `//go:generate statecraft foo.statecraft
   foo/foo.dot` and `//go:generate statecraft foo.statecraft
   foo/foo.go` statements to your calling code so `go generate` will
   run `statecraft` for you in the future.

# Quick Start for Generating Python Code

1. Do the Quick Start steps above.
2. Run `statecraft foo.statecraft foo.py` to generate the
   Python code for your state machine.
1. In your calling code, `import foo`, write an event loop of some
   sort, and in the loop, call foo.tick() for each event.  See the
   comments in `./example/stoplight/python/main.py`.
4. Remember to re-run `statecraft foo.statecraft foo.py` if
   you change your statecraft file.  A Makefile will help with this.

# Example

See ./example/stoplight for a demo of two interacting state machines.
The statecraft DSL for the `car` state machine from that example
looks like this:

```
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

// Declare package name -- this is used verbatim as the `package` name
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
```
