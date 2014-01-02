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

func CheckMultiSearchPaths(t *testing.T, sp []string) {
	files := []string{"libcrt_platform.a", "libgcc.a",
		// Also add a fully-qualified library path.
		path.Join(sp[len(sp) - 1], "libpnacl_irt_shim.a")}
	results := DetermineFilepaths(files, sp)
	ExpectEq(t, results[0], path.Join(sp[0], files[0]))
	ExpectEq(t, results[1], path.Join(sp[0], files[1]))
	ExpectEq(t, results[2], files[2])
}

func TestOneSearchPath(t *testing.T) {
	sp := []string{TestX8632BaseDir()}
	CheckMultiSearchPaths(t, sp)
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
	sp := []string{TestARMBaseDir(), TestX8632BaseDir(),
		TestX8664BaseDir()}
	CheckMultiSearchPaths(t, sp)
}

func TestNoShadowPaths(t *testing.T) {
	sp := []string{TestARMBaseDir(), TestLibDir()}
	files := []string{"libcrt_platform.a", "libfoo_in_libdir.a",
		path.Join(sp[0], "libpnacl_irt_shim.a")}
	results := DetermineFilepaths(files, sp)
	ExpectEq(t, results[0], path.Join(sp[0], files[0]))
	ExpectEq(t, results[1], path.Join(sp[1], files[1]))
	ExpectEq(t, results[2], files[2])

	sp = []string{TestLibDir(), TestARMBaseDir()}
	results = DetermineFilepaths(files, sp)
	ExpectEq(t, results[0], path.Join(sp[1], files[0]))
	ExpectEq(t, results[1], path.Join(sp[0], files[1]))
	ExpectEq(t, results[2], files[2])
}
