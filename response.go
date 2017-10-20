// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"net/http"
	"time"

	"github.com/openrdap/rdap/bootstrap"
)

type Response struct {
	Object          RDAPObject
	BootstrapAnswer *bootstrap.Answer
	HTTP            []*HTTPResponse
}

type RDAPObject interface{}

type HTTPResponse struct {
	URL      string
	Response *http.Response
	Body     []byte
	Error    error
	Duration time.Duration
}
