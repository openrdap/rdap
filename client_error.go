// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

type ClientError struct {
	Text string
}

func (c ClientError) Error() string {
	return c.Text
}
