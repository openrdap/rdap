// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
)

const (
	DefaultCacheDirName = ".openrdap"
)

type DiskCache struct {
	Timeout time.Duration
	Dir string

	lastLoadedModTime map[string]time.Time
}

func NewDiskCache() *DiskCache {
	d := &DiskCache{
		lastLoadedModTime: make(map[string]time.Time),
		Timeout: time.Hour * 24,
	}

	dir, err := homedir.Dir()

	if err != nil {
		panic("Can't determine your home directory")
	}

	d.Dir = filepath.Join(dir, DefaultCacheDirName)

	return d
}

func (d *DiskCache) InitDir() error {
	fileInfo, err := os.Stat(d.Dir)
	if err == nil {
		if fileInfo.IsDir() {
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

	err = ioutil.WriteFile(d.cacheDirPath(filename), data, 0664)
	if err != nil {
		return err
	}

	fileModTime, err := d.modTime(filename)
	if err == nil {
		d.lastLoadedModTime[filename] = fileModTime
	} else {
		return fmt.Errorf("File %s failed to save correctly: %s", filename, err)
	}

	return nil
}

func (d *DiskCache) Load(filename string) ([]byte, error) {
	err := d.InitDir()
	if err != nil {
		return nil, err
	}

	fileModTime, err := d.modTime(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to load %s: %s", filename, err)
	}

	var bytes []byte
	bytes, err = ioutil.ReadFile(d.cacheDirPath(filename))

	if err != nil {
		return nil, err
	}

	d.lastLoadedModTime[filename] = fileModTime

	return bytes, nil
}

func (d *DiskCache) State(filename string) FileState {
	err := d.InitDir()
	if err != nil {
		return Absent
	}

	var expiry time.Time = time.Now().Add(-d.Timeout)
	var state FileState = Absent

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
