// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Driver for go-ld.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	flag.Parse()
	fmt.Printf("Writing to: %s\n", Outfile)
	fmt.Printf("With entry point func: %s\n", EntryPointFunc)
	fmt.Printf("Search Paths to: %s\n", SearchPaths)
	inputs := flag.Args()

	// Go through search-paths to figure out the actual filenames of libs too.
	// Other non-library inputs aren't found in the library paths. 
	lib_full_paths := DetermineFilepaths(LibraryFiles, SearchPaths)
	full_paths := append(inputs, lib_full_paths...)
	fmt.Printf("Full paths of inputs and libs: %v\n", full_paths)

	// Open the files.
	fhandles := make(map[string] *os.File, len(full_paths)) 
	for _, fname := range full_paths {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Print("Failed to open file:", fname, "error:", err)
			return
		}
		defer f.Close()
		fhandles[fname] = f
	}

	// Validate that the inputs are really ELF or .a files full of ELF.
	file_map := ValidateFiles(fhandles)
	fmt.Println("File types: ", file_map)

	// Map the files -> symbol tables.
	for fname, ftyp := range file_map {
		fhandle := fhandles[fname]
		switch ftyp {
		case ELF_FILE:
			body, err := ioutil.ReadAll(fhandle)
			if err != nil {
				panic("Failed to read in file")
			}
			ElfFileHeader := ReadElfHeader(body)
			fmt.Println("Read an ELF file", ElfFileHeader.String())
		default:
			continue
		}
	}

	// Resolve symbols to determine which files to pull in.

	// Figure out which relocations we need.

	// Pull in the files, and lay them out.

	// Fix-up the relocations?
}
