package sc

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	. "github.com/stevegt/goadapt"
)

func TestDot(t *testing.T) {
	infn := "../example/stoplight.statecraft"
	infh, err := os.Open(infn)
	Tassert(t, err == nil, err)

	reffn := "testdata/stoplight.dot"

	m, err := Load(infh)
	Tassert(t, err == nil, err)
	txt := m.ToDot()

	ref, err := ioutil.ReadFile(reffn)
	Tassert(t, err == nil, err)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(ref), txt, false)
	Tassert(t, txt == string(ref), dmp.DiffPrettyText(diffs))

}