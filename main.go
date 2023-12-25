package main

import (
	"io/ioutil"
	"os"
	"strings"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/statecraft/sc"
)

const version = "v0.7.0"

const usage string = `statecraft %s
usage: %s {infn} {outfn}`

// convert panic into clean exit
func exit() {
	rc := 0
	r := recover()
	if r != nil {
		switch concrete := r.(type) {
		case sc.SCErr:
			rc = concrete.Rc
			Fpf(os.Stderr, "%s\n", concrete.Error())
		default:
			// not ours -- re-raise
			panic(r)
		}
	}
	os.Exit(rc)
}

func main() {
	defer exit()
	var err error

	if len(os.Args) != 3 {
		Fpf(os.Stderr, Spf("%s\n", usage), version, os.Args[0])
		os.Exit(1)
	}
	infn := os.Args[1]
	infh, err := os.Open(infn)
	Ck(err)

	m, err := sc.Load(infh, strings.Join(os.Args, " "))
	// Pf("%T\n", err)
	_, ok := err.(sc.SCErr)
	if ok {
		// Pl("salkfdja")
		panic(err)
	}
	Ck(err)
	// Pprint(m)

	outfn := os.Args[2]

	var buf []byte

	if strings.HasSuffix(outfn, ".dot") {
		// XXX replace with ToLang
		buf = m.ToDot()
	} else if strings.HasSuffix(outfn, ".go") {
		// XXX replace with ToLang
		buf = m.ToGo()
	} else if strings.HasSuffix(outfn, ".py") {
		buf = m.ToLang("python")
	}

	err = ioutil.WriteFile(outfn, buf, 0644)
	Ck(err)
}
