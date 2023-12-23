import re
import itertools

class SCErr(Exception):
    def __init__(self, msg, rc):
        self.msg = msg
        self.rc = rc

    def __str__(self):
        return self.msg

def scerr_if(cond, rc, args):
    if cond:
        msg = " ".join(str(a) for a in args)
        raise SCErr(msg, rc)

class Transition:
    def __init__(self, src, event, method, dst):
        self.src = src
        self.event = event
        self.method = method
        self.dst = dst

class State:
    def __init__(self, name, label):
        self.name = name
        self.label = label
        self.transitions = {}

class Machine:
    def __init__(self, cmdline):
        self.cmdline = cmdline
        self.package = None
        self.txt = ""
        self.states = {}
        self.state_names = []
        self.event_names = []
        self.method_names = []

    def add_rule(self, txt):
        parts = txt.split()

        if len(parts) == 0:
            return

        typ = parts[0]

        if typ == "//":
            return
        elif typ == "package":
            self.package = parts[1]
        elif typ == "s":
            scerr_if(len(parts) < 2, 2, "missing state name:", txt)

            s = State(parts[1], " ".join(parts[1:]))

            scerr_if(s.name in self.states, 3, "duplicate state name:", txt)

            self.states[s.name] = s
            self.state_names.append(s.name)

        elif typ == "t":
            scerr_if(len(parts) < 2, 4, "missing transition src:", txt)
            scerr_if(len(parts) < 3, 5, "missing transition event:", txt)
            scerr_if(len(parts) < 4, 6, "missing transition dst:", txt)
            scerr_if(len(parts) > 4, 7, "too many args:", parts)

            event, method = parts[2].split("/") if "/" in parts[2] else (parts[2], "")

            scerr_if(event == "", 8, "missing event:", parts[1])
            scerr_if("/" in event, 9, "too many slashes:", parts[1])

            srcpat = parts[1]
            dstname = parts[3]

            scerr_if(dstname not in self.states, 10, "unknown destination state", dstname, txt)

            re_src = re.compile("^%s$" % srcpat)

            for src_name, state in self.states.items():
                if re_src.match(src_name):
                    if state.transitions.get(event):
                        continue

                    t = Transition(src_name, event, method, dstname)

                    state.transitions[event] = t
                    self.event_names.append(event)
                    self.method_names.append(method)

            scerr_if(srcpat not in self.states, 11, "unknown source state", srcpat, txt)

        else:
            scerr_if(True, 12, "unrecognized entry:", txt)

    def load(self, fh):
        for line in fh:
            txt = line.strip()
            self.add_rule(txt)
            self.txt += f"{txt}\n"

        self.verify()

    def verify(self):
        for state in self.states.values():
            for event_name in self.event_names:
                scerr_if(event_name not in state.transitions, 13,
                         "unhandled event: machine", self.package, "state", state.name, "event", event_name)

    def ls_states(self):
        return [self.states[name] for name in self.state_names]

    def ls_transitions(self):
        transitions = []
        for event_name in self.event_names:
            for state in self.states.values():
                transition = state.transitions.get(event_name)
                if transition:
                    transitions.append(transition)
        return transitions
