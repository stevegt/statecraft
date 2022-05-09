package main

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/statecraft/sc"
)

const usage string = `usage: %s {infn} {outfn}`

func main() {
	if len(os.Args) != 3 {
		Fpf(os.Stderr, Spf("%s\n", usage), os.Args[0])
		os.Exit(1)
	}
	infn := os.Args[1]
	infh, err := os.Open(infn)
	Ck(err)

	m, err := sc.Load(infh, strings.Join(os.Args, " "))
	if errors.Is(err, syscall.ENOSYS) {
		Pl(err)
		os.Exit(2)
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
