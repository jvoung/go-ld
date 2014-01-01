// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Test for file-type detection utilities for ELF linker.

package main

import (
	"os"
	"path"
	"testing"
)

func CheckARFiles(t *testing.T, base_dir string) {
	fnames := [...]string{
		path.Join(base_dir, "libcrt_platform.a"),
		path.Join(base_dir, "libgcc.a"),
		path.Join(base_dir, "libpnacl_irt_shim.a") }

	fhandles := make(map[string] *os.File, len(fnames)) 
	for _, fname := range fnames {
		f, err := os.Open(fname)
		if err != nil {
			t.Fatal("Failed to open file:", fname, "error:", err)
		}
		defer f.Close()
		fhandles[fname] = f
	}

	file_map := ValidateFiles(fhandles)
	for fname, typ := range file_map {
		if typ != AR_FILE {
			t.Errorf("%s should be AR_FILE, but is %v", fname, typ)
		}
	}
}

func CheckELFFiles(t *testing.T, base_dir string) {
	fnames := [...]string{
		path.Join(base_dir, "crtbegin.o"),
		path.Join(base_dir, "crtend.o")}

	fhandles := make(map[string] *os.File, len(fnames)) 
	for _, fname := range fnames {
		f, err := os.Open(fname)
		if err != nil {
			t.Fatal("Failed to open file:", fname, "error:", err)
		}
		defer f.Close()
		fhandles[fname] = f
	}

	file_map := ValidateFiles(fhandles)
	for fname, typ := range file_map {
		if typ != ELF_FILE {
			t.Errorf("%s should be ELF_FILE, but is %v", fname, typ)
		}
	}
}

func CheckBaseDirs(t *testing.T, check_func func(*testing.T, string)) {
	base_dirs := [...]string{
		TestX8632BaseDir(), TestX8664BaseDir(), TestARMBaseDir()}
	for _, b := range base_dirs {
		check_func(t, b)
	}
}

func TestARFile(t *testing.T) {
	CheckBaseDirs(t, CheckARFiles)
}

func TestELFFile(t *testing.T) {
	CheckBaseDirs(t, CheckELFFiles)
}
