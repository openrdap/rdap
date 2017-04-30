// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDiskCache(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	rdapDir := filepath.Join(dir, ".openrdap")
	
	m := NewDiskCache()
	t.Logf("Default cache dir is %s, test cache dir is %s\n", m.Dir, rdapDir)
	m.Dir = rdapDir

	err = m.InitDir()
	if err != nil {
		t.Fatalf("InitDir failed: %s\n", err)
	}

	if !m.IsStale("not-in-cache.json") {
		t.Fatal("m.IsStale() returned false for non-existent file")
	}

	var data []byte
	var isNew bool

	data, isNew, err = m.Load("not-in-cache.json")

	if len(data) != 0 || isNew || err != nil {
		t.Fatal("Load of not-in-cache.json unexpected result")
	}

	var testData []byte = []byte("test")

	err = m.Save("file.json", testData)

	if err != nil {
		t.Fatal("Save failed")
	}

	data, isNew, err = m.Load("file.json")

	if len(data) == 0 || isNew || err != nil || bytes.Compare(data, testData) != 0 {
		t.Fatal("Load of not-in-cache.json unexpected result")
	}

	testData[0] = 'x'
	if data[0] != 't' {
		t.Fatalf("Cache doesn't contain a copy, contains %s", data)
	}

	if m.IsStale("file.json") {
		t.Fatal("m.IsStale returned true for hot cache")
	}

	m.Timeout = 0

	time.Sleep(time.Millisecond)

	if !m.IsStale("file.json") {
		t.Fatal("m.IsStale returned false for stale cache")
	}

}
