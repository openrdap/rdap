// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"fmt"
	"time"
)

type MemoryCache struct {
	Timeout time.Duration
	cache   map[string][]byte
	mtime   map[string]time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string][]byte),
		mtime: make(map[string]time.Time),
		Timeout: time.Hour * 24,
	}
}

func (m *MemoryCache) SetTimeout(timeout time.Duration) {
	m.Timeout = timeout
}

func (m *MemoryCache) Save(filename string, data []byte) error {
	m.cache[filename] = make([]byte, len(data))
	copy(m.cache[filename], data)

	m.mtime[filename] = time.Now()

	return nil
}

func (m *MemoryCache) Load(filename string) ([]byte, error) {
	data, ok := m.cache[filename]

	if !ok {
		return nil, fmt.Errorf("File %s not in cache", filename)
	}

	result := make([]byte, len(data))
	copy(result, data)

	return result, nil
}

func (m *MemoryCache) State(filename string) FileState {
	mtime, ok := m.mtime[filename]

	if !ok {
		return Absent
	}

	expiry := mtime.Add(m.Timeout)

	if expiry.Before(time.Now()) {
		return Expired
	}

	return Good

}
