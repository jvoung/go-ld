// Copyright (c) 2014, Jan Voung
// All rights reserved.

// Lay out the linked files into segments, and adjust the symbol table
// with the new addresses.

package main

import "fmt"

func DoLayout(f_syms []SymbolTable, files []ElfFile) ElfFile {
	result := ElfFile{Body: make([]byte, 0, 0),
		Header: ElfFileHeader{},
		Phdrs:  make([]ProgramHeader, 0, 3),
		Shdrs:  make([]SectionHeader, 0, 0)}

	// Default layout order for PHDRs.
	// The segment to sections map from readelf also shows that
	// although the .note.* are part of the .rodata, they have
	// their own segment (of type NOTE) instead.
	// Same with .eh_frame_hdr, which is R only, but is its own
	// segment of type GNU_EH_FRAME.
	phdr_order := [][]string{{".text"}, // R+E
		{".note", ".rodata", ".reginfo", ".eh_frame_hdr"}, // R
		{".data", ".eh_frame", ".got", ".bss"}} // R + W
	// File offset order...
	fmt.Print(phdr_order)

	// Go through files in order, and copy the sections to
	// the result body, concatenating each section.
	return result
}
