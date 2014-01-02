// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Representation of AR files.

package main

import (
	"os"
	"path"
	"testing"
)

// Test an archive with elf files.
// This has a GNU symbol table.
func TestELFARFileStructure(t *testing.T) {
	test_name := path.Join(TestX8632BaseDir(), "libcrt_platform.a")
	expected_subfiles := map[string]bool{
		"pnacl_irt.o": true,
		"setjmp.o": true,
		"string.o": true}
	f, err := os.Open(test_name)
	if err != nil {
		t.Fatal("Failed to open test AR file")
	}
	defer f.Close()
	ar_file := ReadPlainARFile(f)
	ExpectEq(t, len(expected_subfiles), len(ar_file))
	// Check that the contents are really ELF.
	for fname, hdr_contents := range ar_file {
		ExpectEq(t, string(hdr_contents.Contents[0:4]), ELF_MAGIC)
		ExpectEq(t, expected_subfiles[fname], true)
	}
}

// Test an archive with text files.
// This has long filenames and files with a space in the name.
func TestLongFilenames(t *testing.T) {
	test_name := path.Join(TestLibDir(), "liblong_filename.a")
	expected_subfiles := []string{
		"file_11.txt",
		"file_24.txt",
		"file_nil.txt",
		"file_quick_brown_fox_jumped.txt",
		"file with space in it.txt"}
	expected_contents := map[string]string{
		"file_11.txt": "0123456789\n",
		"file_24.txt": "55555\n55555\n55555\n55555\n",
		"file_nil.txt": "",
		"file_quick_brown_fox_jumped.txt": "the quick brown fox " +
			"jumps over the lazy dog\n",
		"file with space in it.txt": "This file has a space in its name.\n"}
	f, err := os.Open(test_name)
	if err != nil {
		t.Fatal("Failed to open test AR file")
	}
	defer f.Close()
	ar_file := ReadPlainARFile(f)
	ExpectEq(t, len(expected_subfiles), len(ar_file))
	for fname, hdr_contents := range ar_file {
		ec := expected_contents[fname]
		ExpectEq(t, uint32(len(ec)), hdr_contents.Header.FileSize)
		ExpectEq(t, ec, string(hdr_contents.Contents))
	}
}
