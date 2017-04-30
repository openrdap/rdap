// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	DefaultCacheDir = ".openrdap"
)

type DiskCache struct {
	Timeout time.Duration
	Dir string
	cache   map[string][]byte
	mtime   map[string]time.Time
}

func NewDiskCache() *DiskCache {
	d := &DiskCache{
		cache: make(map[string][]byte),
		mtime: make(map[string]time.Time),
		Timeout: time.Hour * 24,
	}

	dir, err := homedir.Dir()

	if err != nil {
		panic("Can't determine your home directory")
	}

	d.Dir = filepath.Join(dir, DefaultCacheDir)

	return d
}

func (d *DiskCache) InitDir() error {
	fileinfo, err := os.Stat(d.Dir)
	if err == nil {
		if fileinfo.IsDir() {
			return nil
		} else {
			return errors.New("Cache dir is not a dir")
		}
	}

	if os.IsNotExist(err) {
		return os.Mkdir(d.Dir, 0775)
	} else {
		return err
	}
}

func (d *DiskCache) SetTimeout(timeout time.Duration) {
	d.Timeout = timeout
}

func (d *DiskCache) Save(filename string, data []byte) error {
	err := d.InitDir()
	if err != nil {
		return err
	}

	// Save and copy into place

	d.cache[filename] = make([]byte, len(data))
	copy(d.cache[filename], data)

	d.mtime[filename] = time.Now()

	return nil
}

func (d *DiskCache) Load(filename string) ([]byte, bool, error) {
	// Stat file
	// if in cache and mtime the same, return that

	// Otherwise try and reload the file, put into cache, return that

	data, ok := d.cache[filename]

	if !ok {
		return nil, false, nil
	}

	result := make([]byte, len(data))
	copy(result, data)

	return result, false, nil
}

func (d *DiskCache) IsStale(filename string) bool {
	
	mtime, ok := d.mtime[filename]

	if !ok {
		return true
	}

	expiry := mtime.Add(d.Timeout)

	if expiry.Before(time.Now()) {
		return true
	}

	return false
}
