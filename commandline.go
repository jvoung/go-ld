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

type search_paths []string
func (p *search_paths) String() string {
	return fmt.Sprint(*p)
}
func (p *search_paths) Set(value string) error {
	*p = append(*p, value)
	return nil
}

// Search paths, specified by -L <path1> -L <path2>.
// Unfortunately, the go flag library doesn't seem to support the space-less
// variant of flags "-L<path1>".
var SearchPaths search_paths
func init() {
	flag.Var(&SearchPaths, "L", "Add a library search path")
}
