// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
)

func TestNewDiskCacheDir(t *testing.T) {
	homedir.DisableCache = true
	defer func() { homedir.DisableCache = false }()

	home, err := os.MkdirTemp("", "home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(home)

	t.Setenv("HOME", home)

	// XDG_CACHE_HOME set (absolute).
	xdg := filepath.Join(home, "custom-cache")
	t.Setenv("XDG_CACHE_HOME", xdg)
	if got, want := NewDiskCache().Dir, filepath.Join(xdg, "openrdap"); got != want {
		t.Errorf("with XDG_CACHE_HOME set: got %q, want %q", got, want)
	}

	// XDG_CACHE_HOME unset -> $HOME/.cache/openrdap.
	t.Setenv("XDG_CACHE_HOME", "")
	if got, want := NewDiskCache().Dir, filepath.Join(home, ".cache", "openrdap"); got != want {
		t.Errorf("with XDG_CACHE_HOME unset: got %q, want %q", got, want)
	}

	// Relative XDG_CACHE_HOME is ignored per the XDG spec.
	t.Setenv("XDG_CACHE_HOME", "relative/path")
	if got, want := NewDiskCache().Dir, filepath.Join(home, ".cache", "openrdap"); got != want {
		t.Errorf("with relative XDG_CACHE_HOME: got %q, want %q", got, want)
	}
}

func TestDiskCache(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	rdapDir := filepath.Join(dir, ".openrdap")

	m1 := NewDiskCache()
	m1.Dir = rdapDir

	m2 := NewDiskCache()
	m2.Dir = rdapDir

	asn1 := []byte(string("file 1"))
	asn2 := []byte(string("file 2"))

	if m1.State("asn.json") != Absent {
		t.Fatalf("asn.json expected absent in m1")
	} else if m2.State("asn.json") != Absent {
		t.Fatalf("asn.json expected absent in m2")
	}

	if err := m1.Save("asn.json", asn1); err != nil {
		t.Fatalf("Save failed: %s", err)
	}

	if m1.State("asn.json") != Good {
		t.Fatalf("asn.json expected good in m1")
	} else if m2.State("asn.json") != ShouldReload {
		t.Fatalf("asn.json expected shouldreload in m2")
	}

	loaded1, err := m1.Load("asn.json")
	loaded2, err := m2.Load("asn.json")

	if m1.State("asn.json") != Good {
		t.Fatalf("asn.json expected good in m1")
	} else if m2.State("asn.json") != Good {
		t.Fatalf("asn.json expected good in m2")
	}

	if bytes.Compare(loaded1, asn1) != 0 {
		t.Fatalf("loaded1(%v) != asn1(%v)", loaded1, asn1)
	} else if bytes.Compare(loaded2, asn1) != 0 {
		t.Fatalf("loaded2(%v) != asn1(%v)", loaded2, asn1)
	}

	time.Sleep(time.Second)

	if err := m2.Save("asn.json", asn2); err != nil {
		t.Fatalf("Save failed: %s", err)
	}

	if m1.State("asn.json") != ShouldReload {
		t.Fatalf("asn.json expected shouldreload in m1")
	} else if m2.State("asn.json") != Good {
		t.Fatalf("asn.json expected good in m2")
	}

	m1.Timeout = 0
	m2.Timeout = 0

	if m1.State("asn.json") != Expired {
		t.Fatal("m1 timeout broken")
	} else if m2.State("asn.json") != Expired {
		t.Fatal("m2 timeout broken")
	}

	m1.Timeout = time.Hour
	m2.Timeout = time.Hour

	loaded1, err = m1.Load("asn.json")
	loaded2, err = m2.Load("asn.json")

	if m1.State("asn.json") != Good {
		t.Fatalf("asn.json expected good in m1")
	} else if m2.State("asn.json") != Good {
		t.Fatalf("asn.json expected good in m2")
	}

	if bytes.Compare(loaded1, asn2) != 0 {
		t.Fatalf("loaded1(%v) != asn2(%v)", loaded1, asn2)
	} else if bytes.Compare(loaded2, asn2) != 0 {
		t.Fatalf("loaded2(%v) != asn2(%v)", loaded2, asn2)
	}
}
