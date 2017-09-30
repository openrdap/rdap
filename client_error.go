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
	ObjectDoesNotExist
)

type ClientError struct {
	Type ClientErrorType
	Text string
}

func (c ClientError) Error() string {
	return c.Text
}

func isClientError(t ClientErrorType, err error) bool {
	if ce, ok := err.(*ClientError); ok {
		if ce.Type == t {
			return true
		}
	}

	return false
}
