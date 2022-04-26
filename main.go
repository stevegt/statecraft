package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
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
}

type Transitions map[Event]*Transition

type Transition struct {
	Method string
	Dst    stateName
}

type States map[stateName]*State

type Machine struct {
	States States
}

func main() {
	Assert(len(os.Args) == 3, "usage: %s in.statecraft out.dot")
	infn := os.Args[1]
	infh, err := os.Open(infn)
	Ck(err)

	outfn := os.Args[2]

	m, err := Load(infh)
	Ck(err)
	// Pprint(m)
	txt := m.ToDot()

	err = ioutil.WriteFile(outfn, []byte(txt), 0644)
	Ck(err)
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
			Debug("added %s %s %v", state.Name, event, t)
			// Pprint(m)
		}
		Assert(found, "unknown source state %s: %s", srcpat, txt)
	default:
		Assert(false, "unrecognized entry: %s")
	}
}

type Edge struct {
	tail  string
	head  string
	label string
}

type Nodes map[string]string
type Edges []Edge

type Dot struct {
	nodes Nodes
	edges Edges
}

func NewDot() *Dot {
	dot := &Dot{}
	dot.nodes = make(Nodes)
	return dot
}

func (d *Dot) AddNode(name, label string) {
	_, ok := d.nodes[name]
	if ok {
		return
	}
	d.nodes[name] = label
}

func (d *Dot) AddEdge(tail, head, label string) {
	edge := Edge{tail, head, label}
	d.edges = append(d.edges, edge)
}

func (d *Dot) String() (out string) {
	out = "digraph \"\" {\n"
	for name, label := range d.nodes {
		if false {
			out += Spf("    %s [label=\"%s\"];\n", name, label)
		} else {
			out += Spf("    %s;\n", name)
		}
	}
	for _, e := range d.edges {
		out += Spf("    %s -> %s [label=\"%s\"];\n", e.tail, e.head, e.label)
	}
	out += "}\n"
	return
}

func (m *Machine) ToDot() (txt string) {

	dot := NewDot()

	for _, src := range m.States {
		dot.AddNode(string(src.Name), src.Label)
		for event, t := range src.Transitions {
			dstName := t.Dst
			dst := m.States[dstName]
			dot.AddNode(string(dst.Name), dst.Label)
			dot.AddEdge(string(src.Name), string(dst.Name), Spf("%s/%s", event, t.Method))
		}
	}

	txt = dot.String()
	return
}
