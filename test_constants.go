// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Filename / directory constants to testing

package main

import "path"

const TestBaseDir = "test_binaries"
func TestX8632BaseDir() string {
	return path.Join(TestBaseDir, "i686")
}
func TestX8664BaseDir() string {
	return path.Join(TestBaseDir, "x86_64")
}
func TestARMBaseDir() string {
	return path.Join(TestBaseDir, "arm")
}
