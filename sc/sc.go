package sc

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	. "github.com/stevegt/goadapt"
)

type Event string
type stateName string

type State struct {
	Name        stateName
	Label       string
	Transitions Transitions
	events      []Event
}

type Transitions map[Event]*Transition

type Transition struct {
	Method string
	Dst    stateName
}

type States map[stateName]*State

type Machine struct {
	States     States
	stateNames []stateName
}

func Load(fh io.Reader) (m *Machine, err error) {
	defer Return(&err)
	m = &Machine{}
	m.States = make(States)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		m.AddRule(scanner.Text())
	}
	err = scanner.Err()
	Ck(err)
	return
}

func (m *Machine) AddRule(txt string) {
	parts := strings.Fields(txt)
	if len(parts) == 0 {
		return
	}
	typ := parts[0]
	switch typ {
	case "//":
		return
	case "s":
		Assert(len(parts) >= 2, "missing state name: %s", txt)
		s := &State{
			Name:        stateName(parts[1]),
			Label:       strings.Join(parts[1:], " "),
			Transitions: make(Transitions),
		}
		_, ok := m.States[s.Name]
		Assert(!ok, "duplicate state name: %s", txt)
		m.States[s.Name] = s
		// maintain ordered list so we can provide reproducible output
		m.stateNames = append(m.stateNames, s.Name)
	case "t":
		Assert(len(parts) > 1, "missing transition src: %s", txt)
		Assert(len(parts) > 2, "missing transition event: %s", txt)
		Assert(len(parts) > 3, "missing transition dst: %s", txt)
		Assert(len(parts) < 5, "too many args: %#v", parts)

		split := func(tok string) (name, method string) {
			toks := strings.Split(tok, "/")
			Assert(len(toks) > 0, tok)
			Assert(len(toks) <= 2, tok)
			name = toks[0]
			if len(toks) == 2 {
				method = toks[1]
			}
			return
		}

		srcpat := parts[1]
		event, method := split(parts[2])
		dstname := parts[3]

		_, ok := m.States[stateName(dstname)]
		Assert(ok, "unknown destination state %s: %s", dstname, txt)

		re := regexp.MustCompile(Spf("^%s$", srcpat))
		found := false
		Debug(txt)
		for name, state := range m.States {
			if !re.MatchString(string(name)) {
				continue
			}
			Debug("matched %s", name)
			found = true
			if state.Transitions[Event(event)] != nil {
				// first rule wins
				Debug("skipping")
				continue
			}
			t := &Transition{
				Method: method,
				Dst:    stateName(dstname),
			}
			state.Transitions[Event(event)] = t
			// maintain ordered list so we can provide reproducible output
			state.events = append(state.events, Event(event))
			Debug("added %s %s %v", state.Name, event, t)
			// Pprint(m)
		}
		Assert(found, "unknown source state %s: %s", srcpat, txt)
	default:
		Assert(false, "unrecognized entry: %s")
	}
}

type Node struct {
	Name  string
	Label string
}

type Edge struct {
	Tail  string
	Head  string
	Label string
}

type Nodes []Node
type Edges []Edge

type Dot struct {
	Nodes Nodes
	Edges Edges
}

func NewDot() *Dot {
	dot := &Dot{}
	return dot
}

func (d *Dot) AddNode(name, label string) {
	for _, n := range d.Nodes {
		if n.Name == name {
			return
		}
	}
	d.Nodes = append(d.Nodes, Node{name, label})
}

func (d *Dot) AddEdge(tail, head, label string) {
	edge := Edge{tail, head, label}
	d.Edges = append(d.Edges, edge)
}

func (d *Dot) String() (out string) {
	out = "digraph \"\" {\n"
	for _, node := range d.Nodes {
		name := node.Name
		label := node.Label
		if false {
			out += Spf("    %s [label=\"%s\"];\n", name, label)
		} else {
			out += Spf("    %s;\n", name)
		}
	}
	for _, e := range d.Edges {
		out += Spf("    %s -> %s [label=\"%s\"];\n", e.Tail, e.Head, e.Label)
	}
	out += "}\n"
	return
}

func (m *Machine) ToDot() (txt string) {

	dot := NewDot()

	for _, name := range m.stateNames {
		state := m.States[name]
		dot.AddNode(string(state.Name), state.Label)
	}

	for _, srcName := range m.stateNames {
		src := m.States[srcName]
		for _, event := range src.events {
			t := src.Transitions[event]
			dstName := t.Dst
			dst := m.States[dstName]
			dot.AddEdge(string(src.Name), string(dst.Name), Spf("%s/%s", event, t.Method))
		}
	}

	txt = dot.String()
	return
}
