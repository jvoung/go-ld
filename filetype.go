// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Simple file-type detection utilities for ELF linker.

package main

const (
	AR_MAGIC = "!<arch>\n"
	THIN_AR_MAGIC = "!<thin>\n"
	ELF_MAGIC = ".ELF"
)

type FileType int
const (
	AR_FILE = iota
	THIN_AR_FILE
	ELF_FILE
	UNKNOWN_FILE
)

func DetermineFileType() FileType {
	return UNKNOWN_FILE
}

func ValidateFiles(files []string) (types map[string]FileType, err error) {
	types = make(map[string]FileType, len(files))
	err = nil
	return
}
