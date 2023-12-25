package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
	. "github.com/stevegt/goadapt"
	"github.com/stevegt/statecraft/sc"
)

// regenerate testdata
const regen bool = false

func TestDot(t *testing.T) {
	infn := "example/stoplight/car.statecraft"
	infh, err := os.Open(infn)
	Tassert(t, err == nil, err)

	m, err := sc.Load(infh, "test dot")
	Tassert(t, err == nil, err)
	got := m.ToDot()

	reffn := "testdata/car.dot"
	if regen {
		err = ioutil.WriteFile(reffn, got, 0644)
		Ck(err)
	}

	ref, err := ioutil.ReadFile(reffn)
	Tassert(t, err == nil, err)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(ref), string(got), false)
	Tassert(t, bytes.Equal(ref, got), dmp.DiffPrettyText(diffs))

}

// testGeneric checks or regenerates the generated code for a given
// language.
func testGeneric(t *testing.T, lang, extension string) {
	infn := "example/stoplight/car.statecraft"
	infh, err := os.Open(infn)
	Tassert(t, err == nil, err)

	m, err := sc.Load(infh, "test "+lang)
	Tassert(t, err == nil, err)
	got := m.ToLang(lang)

	reffn := "testdata/car." + extension
	if regen {
		err = ioutil.WriteFile(reffn, got, 0644)
		Ck(err)
	}

	ref, err := ioutil.ReadFile(reffn)
	Tassert(t, err == nil, err)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(ref), string(got), false)
	Tassert(t, bytes.Equal(ref, got), dmp.DiffPrettyText(diffs))
}

func TestGo(t *testing.T) {
	testGeneric(t, "go", "go")
}

func TestPython(t *testing.T) {
	testGeneric(t, "python", "py")
}
