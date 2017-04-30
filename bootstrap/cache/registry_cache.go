// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import "time"

type RegistryCache interface {
	SetTimeout(timeout time.Duration)
	Save(filename string, data []byte) error
	Load(filename string) ([]byte, bool, error)
	IsStale(filename string) bool
}
