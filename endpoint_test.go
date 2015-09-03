package hnfire

import (
	"fmt"
	"testing"
)

func isStringer(fmt.Stringer) {}

func TestEndpoint(t *testing.T) {
	e := Endpoint("foobar")

	//Fatal if not a stringer
	isStringer(e)

	if e.Child("baz").String() != "foobar/baz" {
		t.Error("Wrong endpoint child")
	}
}
