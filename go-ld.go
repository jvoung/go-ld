// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Driver for go-ld.

package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	fmt.Printf("Writing to %s\n", Outfile)
	fmt.Printf("Search Paths to %s\n", SearchPaths)
	inputs := flag.Args()

	// Go through search-paths to figure out the actual filenames of libs too.
	// Other non-library inputs aren't found in the library paths. 
	lib_full_paths := DetermineFilepaths(LibraryFiles, SearchPaths)
	full_paths := append(inputs, lib_full_paths...)
	fmt.Printf("Full paths of inputs and libs: %v\n", full_paths)

	// Validate that the inputs are really ELF or .a files full of ELF.
	file_map, err := ValidateFiles(full_paths)
	if err != nil {
		fmt.Println("Error: some files aren't ELF...")
		return
	}
	fmt.Println("File types: ", file_map)

	// Map the files -> symbol tables.

	// Resolve symbols to determine which files to pull in.

	// Figure out which relocations we need.

	// Pull in the files, and lay them out.

	// Fix-up the relocations?
}
