package sc

import (
	"bytes"
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

	m, err := Load(infh, "test dot")
	Tassert(t, err == nil, err)
	got := m.ToDot()

	ref, err := ioutil.ReadFile(reffn)
	Tassert(t, err == nil, err)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(ref), string(got), false)
	Tassert(t, bytes.Equal(ref, got), dmp.DiffPrettyText(diffs))

}

func TestGo(t *testing.T) {
	infn := "../example/stoplight.statecraft"
	infh, err := os.Open(infn)
	Tassert(t, err == nil, err)

	m, err := Load(infh, "test go")
	Tassert(t, err == nil, err)
	got := m.ToGo()

	reffn := "testdata/stoplight.go"
	if false {
		err = ioutil.WriteFile(reffn, got, 0644)
		Ck(err)
	}

	ref, err := ioutil.ReadFile(reffn)
	Tassert(t, err == nil, err)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(ref), string(got), false)
	Tassert(t, bytes.Equal(ref, got), dmp.DiffPrettyText(diffs))

}
