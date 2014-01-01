// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Simple file-type detection utilities for ELF linker.

package main

import (
	"os"
	"strings"
)

const (
	AR_MAGIC = "!<arch>\n"
	THIN_AR_MAGIC = "!<thin>\n"
	ELF_MAGIC = "\x7fELF"
)

type FileType int
const (
	AR_FILE = iota
	THIN_AR_FILE
	ELF_FILE
)

func (t FileType) String() string {
	switch t {
	case AR_FILE: return "archive"
	case THIN_AR_FILE: return "thin archive"
	case ELF_FILE: return "ELF file"
	default: return "unknown file type"
	}
}

func ValidateFiles(files map[string] *os.File) map[string]FileType {
	types := make(map[string]FileType, len(files))
	for fname, f := range files {
		sniff_amt := 8
		buf := make([]byte, sniff_amt)
		n, err := f.ReadAt(buf, 0)
		if err != nil || n != sniff_amt {
			panic("Unable to read file: " + fname + " w/ err: " + err.Error())
		}
		magic := string(buf)
		if strings.HasPrefix(magic, ELF_MAGIC) {
			types[fname] = ELF_FILE
		} else if strings.HasPrefix(magic, AR_MAGIC) {
			types[fname] = AR_FILE
		} else if strings.HasPrefix(magic, THIN_AR_MAGIC) {
			types[fname] = THIN_AR_FILE
		} else {
			panic("File: " + fname + " is not ELF or .a file")
		}
	}
	return types
}
