package sc

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"regexp"
	"strings"
	"syscall"
	"text/template"

	. "github.com/stevegt/goadapt"
)

type State struct {
	Name        string
	Label       string
	Transitions Transitions
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
	Cmdline string
	Package string
	Machine string
	Txt     string
	States  States
	// maintain ordered lists so we can provide reproducible output
	StateNames  []string
	EventNames  []string
	MethodNames []string
}

func Load(fh io.Reader, cmdline string) (m *Machine, err error) {
	defer Return(&err)
	m = &Machine{}
	m.Cmdline = cmdline
	m.States = make(States)
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		txt := scanner.Text()
		m.AddRule(txt)
		m.Txt += Spf("%s\n", txt)
	}
	err = scanner.Err()
	Ck(err)

	err = m.Verify()
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
	case "package":
		m.Package = parts[1]
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
		m.StateNames = appendUniq(m.StateNames, s.Name)
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
			m.EventNames = appendUniq(m.EventNames, event)
			m.MethodNames = appendUniq(m.MethodNames, method)
			Debug("added %s %s %v", state.Name, event, t)
			// Pprint(m)
		}
		Assert(found, "unknown source state %s: %s", srcpat, txt)
	default:
		Assert(false, "unrecognized entry: %s", txt)
	}
}

/*
type SCError struct {
	msg string
}

func (e SCError) Error() string {
	return e.msg
}

func NewSCError(args ...interface{}) SCError {
	return SCError{msg: FormatArgs(args)}
}
*/

func (m *Machine) Verify() (err error) {
	for stateName, state := range m.States {
		for _, eventName := range m.EventNames {
			_, ok := state.Transitions[eventName]
			if !ok {
				err = fmt.Errorf("%w: unhandled event: machine %v, state %v, event %v", syscall.ENOSYS, m.Package, stateName, eventName)
				return
			}
		}
	}
	return
}

func appendUniq(in []string, add string) (out []string) {
	out = in[:]
	if add == "" {
		return
	}
	found := false
	for _, s := range out {
		if s == add {
			found = true
			break
		}
	}
	if !found {
		out = append(out, add)
	}
	return
}

func (m *Machine) LsStates() (out []*State) {
	for _, name := range m.StateNames {
		out = append(out, m.States[name])
	}
	return
}

func (m *Machine) LsTransitions() (out []*Transition) {
	for _, srcName := range m.StateNames {
		src := m.States[srcName]
		// we iterate over m.EventNames instead of src.Transitions
		// here because we want to preserve ordering for reproducible
		// output
		for _, eventName := range m.EventNames {
			t, ok := src.Transitions[eventName]
			if ok {
				out = append(out, t)
			}
		}
	}
	return
}

/*
func (m *Machine) LsEvents() (out []string) {
	for _, srcName := range m.stateNames {
		src := m.States[srcName]
		for _, event := range src.events {
			found := false
			for _, e := range out {
				if event == e {
					found = true
					break
				}
			}
			if !found {
				out = append(out, event)
			}
		}
	}
	return
}

func (m *Machine) LsMethods() (out []string) {
	for _, srcName := range m.stateNames {
		src := m.States[srcName]
		for _, event := range src.events {
			found := false
			for _, e := range out {
				if event == e {
					found = true
					break
				}
			}
			if !found {
				out = append(out, event)
			}
		}
	}
	return
}
*/

//go:embed template/*
var fs embed.FS

func (m *Machine) ToDot() (out []byte) {
	t := template.Must(template.ParseFS(fs, "template/dot.ttmpl"))
	var buf bytes.Buffer
	err := t.Execute(&buf, m)
	Ck(err)
	out = buf.Bytes()
	return
}

func (m *Machine) ToGo() (out []byte) {
	t := template.Must(template.ParseFS(fs, "template/go.ttmpl"))
	var buf bytes.Buffer
	err := t.Execute(&buf, m)
	Ck(err)
	out = buf.Bytes()
	return
}
