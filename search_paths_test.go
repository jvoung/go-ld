// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Test for the simple search_paths utils.

package main

import (
	"path"
	"testing"
)

func TestNoPathsNoDirs(t *testing.T) {
	DetermineFilepaths([]string{}, []string{})
}

func TestOneSearchPath(t *testing.T) {
	sp := []string{TestX8632BaseDir()}
	files := []string{"libcrt_platform.a", "libgcc.a",
		path.Join(TestX8632BaseDir(), "libpnacl_irt_shim.a")}
	results := DetermineFilepaths(files, sp)
	if results[0] != path.Join(TestX8632BaseDir(), files[0]) {
		t.Errorf("Result %s is wrong", results[0])
	}
	if results[1] != path.Join(TestX8632BaseDir(), files[1]) {
		t.Errorf("Result %s is wrong", results[1])
	}
	if results[2] != files[2] {
		t.Errorf("Result %s is wrong", results[2])
	}
}

func CheckMultiSearchPaths(t *testing.T, sp []string) {
	files := []string{"libcrt_platform.a", "libgcc.a",
		path.Join(sp[1], "libpnacl_irt_shim.a")}
	results := DetermineFilepaths(files, sp)
	if results[0] != path.Join(sp[0], files[0]) {
		t.Errorf("Result %s is wrong", results[0])
	}
	if results[1] != path.Join(sp[0], files[1]) {
		t.Errorf("Result %s is wrong", results[1])
	}
	if results[2] != files[2] {
		t.Errorf("Result %s is wrong", results[2])
	}	
}

func TestTwoSearchPathsA(t *testing.T) {
	sp := []string{TestX8632BaseDir(), TestX8664BaseDir()}
	CheckMultiSearchPaths(t, sp)
}

func TestTwoSearchPathsB(t *testing.T) {
	sp := []string{TestX8664BaseDir(), TestX8632BaseDir()}
	CheckMultiSearchPaths(t, sp)
}

func TestThreeSearchPaths(t *testing.T) {
	sp := []string{TestARMBaseDir(), TestX8632BaseDir(), TestX8664BaseDir()}
	CheckMultiSearchPaths(t, sp)
}

func TestNoShadowPaths(t *testing.T) {
	sp := []string{TestARMBaseDir(), TestLibDir()}
	files := []string{"libcrt_platform.a", "lib_x8632_foo.a",
		path.Join(sp[0], "libpnacl_irt_shim.a")}
	results := DetermineFilepaths(files, sp)
	if results[0] != path.Join(sp[0], files[0]) {
		t.Errorf("Result %s is wrong", results[0])
	}
	if results[1] != path.Join(sp[1], files[1]) {
		t.Errorf("Result %s is wrong", results[1])
	}
	if results[2] != files[2] {
		t.Errorf("Result %s is wrong", results[2])
	}
}
