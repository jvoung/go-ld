// Copyright (c) 2014, Jan Voung
// All rights reserved.

// More methods for testing.

package main

import "testing"

func ExpectEq(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		t.Errorf("Expected %v == %v", v1, v2)
	}
}

func ExpectEqM(t *testing.T, v1, v2 interface{}, s string) {
	if v1 != v2 {
		t.Errorf("%s Expected %v == %v", s, v1, v2)
	}
}

func AssertEq(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		t.Fatalf("Expected %v == %v", v1, v2)
	}
}

func AssertEqM(t *testing.T, v1, v2 interface{}, s string) {
	if v1 != v2 {
		t.Fatalf("%s Expected %v == %v", s, v1, v2)
	}
}
