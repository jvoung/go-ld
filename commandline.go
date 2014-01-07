// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Commandline flags for go-ld driver.
// This is a very simple linker and does not support many options.

package main

import (
	"flag"
	"fmt"
)

// Filename for the output.
var Outfile string
func init() {
	const (
		defaultOutfile = "a.out"
		usage = "The output filename"
	)
	flag.StringVar(&Outfile, "output", defaultOutfile, usage)
	flag.StringVar(&Outfile, "o", defaultOutfile, usage + " (shorthand)")
}

// Search paths for "-l" libraries, specified by -L <path1> -L <path2>.
// Unfortunately, the go flag library doesn't seem to support the space-less
// variant of flags "-L<path1>".
type search_paths []string
func (p *search_paths) String() string {
	return fmt.Sprint(*p)
}
func (p *search_paths) Set(value string) error {
	*p = append(*p, value)
	return nil
}

// The parsed search paths.
var SearchPaths search_paths
func init() {
	flag.Var(&SearchPaths, "L", "Add a library (-l) search path")
}


// Libraries, specified by "-l=libfoo.a" or "-l libfoo.a"... I haven't
// looked into how to coax the go flag package into accepting "-lfoo"
// and "-l:libfoo.a".
//
// This also doesn't track where static libraries w.r.t. relocatable
// files. Just assume all libraries show up after relocatable files.
type lib_names []string
func (l *lib_names) String() string {
	return fmt.Sprint(*l)
}
func (l *lib_names) Set(value string) error {
	*l = append(*l, value)
	return nil
}

// The parsed library files paths.
var LibraryFiles lib_names
func init() {
	flag.Var(&LibraryFiles, "l", "Add a library as input")
}

// The entry point function.
var EntryPointFunc string
func init() {
	const (
		defaultEntry = "_start"
		usage = "Set the entry point function name (default _start)"
	)
	flag.StringVar(&EntryPointFunc, "entry", defaultEntry, usage)
	flag.StringVar(&EntryPointFunc, "e", defaultEntry, usage + " (shorthand)")
}
