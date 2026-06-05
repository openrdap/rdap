// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
)

const (
	defaultCacheDirName = "openrdap"
)

// A DiskCache caches Service Registry files on disk.
//
// By default they're saved as $XDG_CACHE_HOME/openrdap/{asn,dns,ipv4,ipv6}.json.
// File mtimes are used to calculate cache expiry.
//
// The cache directory is created automatically as needed.
type DiskCache struct {
	// Duration files are stored before they're considered expired.
	//
	// The default is 24 hours.
	Timeout time.Duration

	// Directory to store cached files in.
	//
	// The default is $XDG_CACHE_HOME/openrdap (falling back to
	// $HOME/.cache/openrdap).
	Dir string

	lastLoadedModTime map[string]time.Time
}

// NewDiskCache creates a new DiskCache.
func NewDiskCache() *DiskCache {
	d := &DiskCache{
		Timeout:           time.Hour * 24,
		lastLoadedModTime: make(map[string]time.Time),
	}

	// Honor $XDG_CACHE_HOME, falling back to $HOME/.cache. Relative values are
	// ignored, per the XDG Base Directory spec.
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if !filepath.IsAbs(cacheDir) {
		home, err := homedir.Dir()
		if err != nil {
			panic("Cannot determine home directory: HOME environment variable not set or inaccessible")
		}

		cacheDir = filepath.Join(home, ".cache")
	}

	d.Dir = filepath.Join(cacheDir, defaultCacheDirName)

	return d
}

// InitDir creates the cache directory if it does not already exist.
//
// Returns true if the directory was created, or false if it already exists/or
// on error.
func (d *DiskCache) InitDir() (bool, error) {
	fileInfo, err := os.Stat(d.Dir)
	if err == nil {
		if fileInfo.IsDir() {
			return false, nil
		}

		return false, errors.New("Cache dir is not a dir")
	}

	if os.IsNotExist(err) {
		if err = os.MkdirAll(d.Dir, 0775); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, err
}

// SetTimeout sets the duration each Service Registry file can be stored before
// its State() is Expired.
func (d *DiskCache) SetTimeout(timeout time.Duration) {
	d.Timeout = timeout
}

// Save saves the file |filename| with |data| to disk.
//
// The cache directory is created if necessary.
func (d *DiskCache) Save(filename string, data []byte) error {
	if _, err := d.InitDir(); err != nil {
		return err
	}

	if err := os.WriteFile(d.cacheDirPath(filename), data, 0664); err != nil {
		return err
	}

	fileModTime, err := d.modTime(filename)
	if err != nil {
		return fmt.Errorf("File %s failed to save correctly: %s", filename, err)
	}

	d.lastLoadedModTime[filename] = fileModTime

	return nil
}

// Load loads the file |filename| from disk.
//
// Since Service Registry files do not change much, the file is returned even
// if its State() is Expired.
//
// An error is returned if the file is not on disk.
func (d *DiskCache) Load(filename string) ([]byte, error) {
	fileModTime, err := d.modTime(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to load %s: %s", filename, err)
	}

	var bytes []byte

	bytes, err = os.ReadFile(d.cacheDirPath(filename))
	if err != nil {
		return nil, err
	}

	d.lastLoadedModTime[filename] = fileModTime

	return bytes, nil
}

// State returns the cache state of the file |filename|.
//
// The returned state is one of: Absent, Good, ShouldReload, Expired.
func (d *DiskCache) State(filename string) FileState {
	var expiry = time.Now().Add(-d.Timeout)
	var state = Absent

	fileModTime, err := d.modTime(filename)
	if err == nil {
		if fileModTime.After(expiry) {
			state = ShouldReload

			lastLoadedModTime, haveLoaded := d.lastLoadedModTime[filename]
			if haveLoaded && !fileModTime.After(lastLoadedModTime) {
				state = Good
			}
		} else {
			state = Expired
		}
	}

	return state
}

func (d *DiskCache) modTime(filename string) (time.Time, error) {
	var fileInfo os.FileInfo
	fileInfo, err := os.Stat(d.cacheDirPath(filename))

	if err != nil {
		return time.Time{}, err
	}

	return fileInfo.ModTime(), nil
}

func (d *DiskCache) cacheDirPath(filename string) string {
	return filepath.Join(d.Dir, filename)
}
