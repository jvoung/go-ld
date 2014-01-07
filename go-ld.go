// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Driver for go-ld.

package main

import (
	"flag"
	"fmt"
	"os"
)

type read_symbols_result struct {
	fname string
	st    interface{}
}

func read_symbols_task(fname string, ftyp FileType,
	fhandles map[string]*os.File,
	done_ch chan read_symbols_result) {
	fhandle := fhandles[fname]
	switch ftyp {
	case ELF_FILE:
		elf_file := ReadElfFileFD(fhandle)
		st := ReadSymbols(&elf_file)
		done_ch <- read_symbols_result{fname, st}
	case AR_FILE, THIN_AR_FILE:
		ar_file := ReadARFile(fhandle, ftyp)
		ar_elf := WrapARElf(&ar_file)
		ar_syms := make(map[string]SymbolTable, len(ar_elf))
		done_ch <- read_symbols_result{fname, ar_syms}
	default:
		panic("Unknown file type")
	}
}

func main() {
	flag.Parse()
	fmt.Printf("Writing to: %s\n", Outfile)
	fmt.Printf("With entry point func: %s\n", EntryPointFunc)
	fmt.Printf("Search Paths to: %s\n", SearchPaths)
	inputs := flag.Args()
	fmt.Println("Other inputs:", inputs)

	// Go through search-paths to figure out the actual filenames of libs too.
	// Other non-library inputs aren't found in the library paths.
	lib_full_paths := DetermineFilepaths(LibraryFiles, SearchPaths)
	full_paths := append(inputs, lib_full_paths...)
	fmt.Printf("Full paths of inputs and libs: %v\n", full_paths)

	// Open the files.
	fhandles := make(map[string]*os.File, len(full_paths))
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
	// TODO(jvoung): Replace the interface{} with some real interface.
	fname_symbols := make(map[string]interface{}, len(fhandles))
	// Channel for reading them in parallel.
	read_symbols := make(chan read_symbols_result, len(fhandles))
	for fname, ftyp := range file_map {
		go read_symbols_task(fname, ftyp, fhandles, read_symbols)
	}
	for i := 0; i < len(file_map); i++ {
		result := <-read_symbols
		fname_symbols[result.fname] = result.st
	}
	fmt.Println("Fname_symbols: ", fname_symbols)

	// Resolve symbols to determine which files to pull in.

	// Figure out which relocations we need.

	// Pull in the files, and lay them out.

	// Fix-up the relocations?
}
