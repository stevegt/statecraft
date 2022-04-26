package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/goccy/go-graphviz"
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
	infn := os.Args[1]
	infh, err := os.Open(infn)
	Ck(err)

	// outfh := os.Stdout

	m, err := Load(infh)
	Ck(err)
	// Pprint(m)
	buf := m.ToDot()
	Pl(string(buf))
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

func (m *Machine) ToDot() (buf []byte) {

	gv := graphviz.New()
	g, err := gv.Graph()
	Ck(err)
	defer func() {
		err := g.Close()
		Ck(err)
		gv.Close()
	}()

	for name, state := range m.States {
		src, err := g.CreateNode(string(name))
		Ck(err)
		for event, t := range state.Transitions {
			dst, err := g.CreateNode(string(t.Dst)) // returns same node if already exists
			Ck(err)
			edge, err := g.CreateEdge(string(event), src, dst)
			Ck(err)
			edge.SetLabel(Spf("%s/%s", event, t.Method))
		}
	}

	var dotbuf bytes.Buffer
	err = gv.Render(g, "dot", &dotbuf)
	Ck(err)
	buf = dotbuf.Bytes()
	return
}
