package main

import (
	"io/ioutil"
	"os"

	. "github.com/stevegt/goadapt"
	"github.com/stevegt/statecraft/sc"
)

func main() {
	Assert(len(os.Args) == 3, "usage: %s in.statecraft out.dot")
	infn := os.Args[1]
	infh, err := os.Open(infn)
	Ck(err)

	outfn := os.Args[2]

	m, err := sc.Load(infh)
	Ck(err)
	// Pprint(m)
	buf := m.ToDot()

	err = ioutil.WriteFile(outfn, buf, 0644)
	Ck(err)
}
