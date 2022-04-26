package sc

import (
	"bufio"
	"bytes"
	"embed"
	"io"
	"regexp"
	"strings"
	"text/template"

	. "github.com/stevegt/goadapt"
)

type State struct {
	Name        string
	Label       string
	Transitions Transitions
	events      []string
}

// map[eventName]*Transition
type Transitions map[string]*Transition

type Transition struct {
	Src    string
	Event  string
	Method string
	Dst    string
}

// map[stateName]*State
type States map[string]*State

type Machine struct {
	States     States
	stateNames []string
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
			Name:        parts[1],
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

		_, ok := m.States[dstname]
		Assert(ok, "unknown destination state %s: %s", dstname, txt)

		re := regexp.MustCompile(Spf("^%s$", srcpat))
		found := false
		Debug(txt)
		for srcName, state := range m.States {
			if !re.MatchString(string(srcName)) {
				continue
			}
			Debug("matched %s", srcName)
			found = true
			if state.Transitions[event] != nil {
				// first rule wins
				Debug("skipping")
				continue
			}
			t := &Transition{
				Src:    srcName,
				Event:  event,
				Method: method,
				Dst:    dstname,
			}
			state.Transitions[event] = t
			// maintain ordered list so we can provide reproducible output
			state.events = append(state.events, event)
			Debug("added %s %s %v", state.Name, event, t)
			// Pprint(m)
		}
		Assert(found, "unknown source state %s: %s", srcpat, txt)
	default:
		Assert(false, "unrecognized entry: %s")
	}
}

func (m *Machine) LsStates() (out []*State) {
	for _, name := range m.stateNames {
		out = append(out, m.States[name])
	}
	return
}

func (m *Machine) LsTransitions() (out []*Transition) {
	for _, srcName := range m.stateNames {
		src := m.States[srcName]
		for _, event := range src.events {
			t := src.Transitions[event]
			out = append(out, t)
		}
	}
	return
}

//go:embed template/*
var fs embed.FS

func (m *Machine) ToDot() (out []byte) {
	t := template.Must(template.ParseFS(fs, "template/tmpl.dot"))
	var buf bytes.Buffer
	err := t.Execute(&buf, m)
	Ck(err)
	out = buf.Bytes()
	return
}
