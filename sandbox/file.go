// OpenRDAP
// Copyright 2018 Tom Harwood
// MIT License, see the LICENSE file.

package sandbox

import (
	"errors"
	"log"
	"os"
	"path"
	"runtime"
)

var sandboxPath string

// LoadFile reads the named sandboxed file and returns its contents. It returns
// an error if the file is not one of the recognised sandbox files.
func LoadFile(filename string) ([]byte, error) {
	if !IsFileInSandbox(filename) {
		return nil, errors.New("file not found in sandbox")
	}

	var body []byte

	if len(sandboxPath) == 0 {
		sandboxPath = findPackagePath()
	}

	body, err := os.ReadFile(path.Join(sandboxPath, filename))
	if err != nil {
		log.Panic(err)
	}

	return body, nil
}

func findPackagePath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("runtime.Caller() failed")
	}

	dir, _ := path.Split(filename)

	return dir
}

// IsFileInSandbox reports whether filename is one of the recognised sandbox
// files.
func IsFileInSandbox(filename string) bool {
	switch filename {
	case "DigiCert_RDAP_Pilot_Client_Certificate.p12",
		"DigiCert_RDAP_Pilot_Client_Certificate_Expired.p12",
		"DigiCert_RDAP_Pilot_Client_Certificate_Revoked.p12":
		return true
	default:
		return false
	}
}
