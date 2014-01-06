// Copyright (c) 2014, Jan Voung
// All rights reserved.

// More methods for testing.

package main

import (
	"runtime"
	"testing"
)

func ExpectEq(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			t.Errorf("(%s:%d) Expected %v == %v", file, line, v1, v2)
		} else {
			t.Errorf("Expected %v == %v", v1, v2)
		}
	}
}

func ExpectEqM(t *testing.T, v1, v2 interface{}, s string) {
	if v1 != v2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			t.Errorf("%s (%s:%d) Expected %v == %v", s, file, line, v1, v2)
		} else {
			t.Errorf("%s Expected %v == %v", s, v1, v2)
		}
	}
}

func AssertEq(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			t.Fatalf("(%s:%d) Expected %v == %v", file, line, v1, v2)
		} else {
			t.Fatalf("Expected %v == %v", v1, v2)
		}
	}
}

func AssertEqM(t *testing.T, v1, v2 interface{}, s string) {
	if v1 != v2 {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			t.Fatalf("%s (%s:%d) Expected %v == %v", s, file, line, v1, v2)
		} else {
			t.Fatalf("%s Expected %v == %v", s, v1, v2)
		}
	}
}
