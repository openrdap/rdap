// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package bootstrap

import (
	"testing"

	"github.com/skip2/rdap/test"
)

func TestHello(t *testing.T) {
	test.Start(test.Bootstrap)
	defer test.Finish()

	c := NewClient()
	err := c.DownloadAll()

	if err != nil {
		t.Fatalf("DownloadAll() error: %s", err)
	}
}
