package engine

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestEngineRoot(t *testing.T) {
	tmp, err := ioutil.TempDir("", "docker-test-TestEngineCreateDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	dir := path.Join(tmp, "dir")
	eng, err := New(dir)
	if err != nil {
		t.Fatal(err)
	}
	if st, err := os.Stat(dir); err != nil {
		t.Fatal(err)
	} else if !st.IsDir() {
		t.Fatalf("engine.New() created something other than a directory at %s", dir)
	}
	if r := eng.Root(); r != dir {
		t.Fatalf("Expected: %v\nReceived: %v", dir, r)
	}
}

func TestEngineString(t *testing.T) {
	eng1 := newTestEngine(t)
	defer os.RemoveAll(eng1.Root())
	eng2 := newTestEngine(t)
	defer os.RemoveAll(eng2.Root())
	s1 := eng1.String()
	s2 := eng2.String()
	if eng1 == eng2 {
		t.Fatalf("Different engines should have different names (%v == %v)", s1, s2)
	}
}

func TestEngineLogf(t *testing.T) {
	eng := newTestEngine(t)
	defer os.RemoveAll(eng.Root())
	input := "Test log line"
	if n, err := eng.Logf("%s\n", input); err != nil {
		t.Fatal(err)
	} else if n < len(input) {
		t.Fatalf("Test: Logf() should print at least as much as the input\ninput=%d\nprinted=%d", len(input), n)
	}
}
