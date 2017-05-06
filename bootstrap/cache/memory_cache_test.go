// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"bytes"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	m := NewMemoryCache()
	if m.State("not-in-cache.json") != Absent {
		t.Fatal("m.State() returned non-Absent for absent file")
	}

	var data []byte
	var err error

	data, err = m.Load("not-in-cache.json")

	if err == nil {
		t.Fatal("Load of not-in-cache.json unexpected result")
	}

	var testData []byte = []byte("test")

	err = m.Save("file.json", testData)

	if err != nil {
		t.Fatal("Save failed")
	}

	data, err = m.Load("file.json")

	if len(data) == 0 || err != nil || bytes.Compare(data, testData) != 0 {
		t.Fatal("Load of file.json unexpected result")
	}

	testData[0] = 'x'
	if data[0] != 't' {
		t.Fatalf("Cache doesn't contain a copy, contains %s", data)
	}

	if m.State("file.json") != Good {
		t.Fatal("m.State() returned non-Good for cached file")
	}

	m.Timeout = 0

	time.Sleep(time.Millisecond)

	if m.State("file.json") != Expired {
		t.Fatal("m.State() returned non-Expired for expired file")
	}

}
