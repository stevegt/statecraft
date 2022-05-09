package main

import (
	"io/ioutil"
	"os"
	"strings"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/statecraft/sc"
)

const usage string = `usage: %s {infn} {outfn}`

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
		Fpf(os.Stderr, Spf("%s\n", usage), os.Args[0])
		os.Exit(1)
	}
	infn := os.Args[1]
	infh, err := os.Open(infn)
	Ck(err)

	m, err := sc.Load(infh, strings.Join(os.Args, " "))
	_, ok := err.(sc.SCErr)
	if ok {
		panic(err)
	}
	Ck(err)
	// Pprint(m)

	outfn := os.Args[2]

	var buf []byte

	if strings.HasSuffix(outfn, ".dot") {
		buf = m.ToDot()
	} else if strings.HasSuffix(outfn, ".go") {
		buf = m.ToGo()
	}

	err = ioutil.WriteFile(outfn, buf, 0644)
	Ck(err)
}
