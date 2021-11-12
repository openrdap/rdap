// OpenRDAP
// Copyright 2017 Tom Harwood
// MIT License, see the LICENSE file.

package rdap

import (
	"fmt"
	"testing"

	"github.com/openrdap/rdap/test"
)

func verboseFunc() func(text string) {
	if testing.Verbose() {
		return func(text string) {
			fmt.Printf("# %s\n", text)
		}
	}

	return func(text string) {
	}
}

func TestClientQueryDomain(t *testing.T) {
	test.Start(test.Bootstrap)
	test.Start(test.Responses)
	defer test.Finish()

	client := &Client{
		Verbose: verboseFunc(),
	}

	domain, err := client.QueryDomain("example.cz")

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	} else if domain == nil {
		t.Errorf("Unexpected nil Domain")
	} else if domain.LDHName != "example.cz" {
		t.Errorf("Unexpected LDHName %s", domain.LDHName)
	}
}

func TestClientQueryDomain404(t *testing.T) {
	test.Start(test.Bootstrap)
	test.Start(test.Responses)
	defer test.Finish()

	client := &Client{
		Verbose: verboseFunc(),
	}

	_, err := client.QueryDomain("non-existent.cz")

	if err == nil {
		t.Errorf("Unexpected success")
	} else if !isClientError(ObjectDoesNotExist, err) {
		t.Errorf("Unexpected err %s", err)
	}
}

func TestClientQueryDomainWrongType(t *testing.T) {
	test.Start(test.Bootstrap)
	test.Start(test.Responses)
	defer test.Finish()

	client := &Client{
		Verbose: verboseFunc(),
	}

	_, err := client.QueryDomain("wrong-response-type.cz")

	if err == nil {
		t.Errorf("Unexpected success")
	} else if !isClientError(WrongResponseType, err) {
		t.Errorf("Unexpected err %s", err)
	}
}

func TestClientQueryDomainMalformed(t *testing.T) {
	test.Start(test.Bootstrap)
	test.Start(test.Responses)
	defer test.Finish()

	client := &Client{
		Verbose: verboseFunc(),
	}

	_, err := client.QueryDomain("malformed.cz")

	if err == nil {
		t.Errorf("Unexpected success")
	} else if !isClientError(NoWorkingServers, err) {
		t.Errorf("Unexpected err %s", err)
	}
}

// test Do()
// 1) success, 1 of each query
// 2) bootstrap not supported
// 3) bootstrap no match
// test Help...
