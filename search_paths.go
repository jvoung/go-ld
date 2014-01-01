// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Determine full file paths based on search paths.

package main

import (
	"os"
	"path"
)

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func DetermineFilepaths(input_paths []string, search_paths []string) []string {
	out := make([]string, 0, len(input_paths))
	for _, p := range input_paths {
		if fileExists(p) {
			out = append(out, p)
			continue
		}
		result := ""
		for _, sp := range search_paths {
			joined := path.Join(sp, p)
			if fileExists(joined) {
				result = joined
				break
			}
		}
		if result != "" {
			out = append(out, result)
		} else {
			panic("Cannot find input file " + p)
		}
	}
	return out
}
