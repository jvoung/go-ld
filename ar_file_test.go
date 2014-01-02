// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Representation of AR files.

package main

import (
	"os"
	"path"
	"testing"
)

func TestARFileStructure(t *testing.T) {
	test_name := path.Join(TestX8632BaseDir(), "libcrt_platform.a")
	expected_subfiles := []string{"pnacl_irt.o", "setjmp.o", "string.o"}
	f, err := os.Open(test_name)
	if err != nil {
		t.Fatal("Failed to open test AR file")
	}
	defer f.Close()
	ar_file := ReadPlainARFile(f)
	if len(ar_file) != 3 {
		t.Errorf("Expected 3 files (%s) in test AR file %s, got %d",
			test_name, expected_subfiles, len(ar_file))
	}
	// TODO(jvoung): Check more...
}

func TestLongFilenames(t *testing.T) {
	// TODO(jvoung): test...
}
