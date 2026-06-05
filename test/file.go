// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package test

import (
	"log"
	"os"
	"path"
	"runtime"
)

var testDataPath string

func LoadFile(filename string) []byte {
	var body []byte

	if len(testDataPath) == 0 {
		testDataPath = findTestDataPath()
	}

	body, err := os.ReadFile(path.Join(testDataPath, filename))

	if err != nil {
		log.Panic(err)
	}

	return body
}

func findTestDataPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller() failed")
	}

	dir, _ := path.Split(filename)

	return path.Join(dir, "testdata")
}
