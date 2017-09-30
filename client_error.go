// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

type ClientErrorType uint

const (
	_ ClientErrorType = iota

	InputError
	BootstrapNotSupported
	BootstrapNoMatch
	WrongResponseType
	NoWorkingServers
)

type ClientError struct {
	Type ClientErrorType
	Text string
}

func (c ClientError) Error() string {
	return c.Text
}
