// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

type ClientErrorType uint

const (
	BootstrapNotSupported ClientErrorType = iota
)

type ClientError struct {
	Type ClientErrorType
	Text string
}

func (c ClientError) Error() string {
	return c.Text
}
