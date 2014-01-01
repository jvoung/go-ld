// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Representation of ELF file (separate from go's debug/elf,
// modulo some constants and printing functions).

package main

import (
	"debug/elf"
)

// Non-class dependent parts of the header.
type ElfFileHeader struct {
	// Offset 0-3 is the ELF magic number.
	Class elf.Class // Offset 4
	Data elf.Data // Offset 5 (endianness, etc.)
	EI_Version elf.Version // Offset 6
	OSABI elf.OSABI // Offset 7
	ABIVersion uint8 // Offset 8
	// Padding Offset 9-15
	Type elf.Type // Offset 0x10
	Machine elf.Machine // Offset 0x12
	E_Version uint32 // Offset 0x14
}

// Class-dependent parts of the header.
type ElfHeader32 struct {
	Entry uint32 // Offset 0x18
	Phoff uint32 // Offset 0x1c
	Shoff uint32 // Offset 0x20
	Flags uint32 // Offset 0x24
	FileHeaderSize uint16 // Offset 0x28
	Phentsize uint16
	Phnum uint16
	Shentsize uint16
	Shnum uint16
	Shstrndx uint16
}

type ElfHeader64 struct {
	Entry uint64 // Offset 0x18
	Phoff uint64 // Offset 0x20
	Shoff uint64 // Offset 0x28
	Flags uint32 // Offset 0x30
	FileHeaderSize uint16 // Offset 0x34
	Phentsize uint16
	Phnum uint16
	Shentsize uint16
	Shnum uint16
	Shstrndx uint16
}

// Rounded-up headers.
type ElfHeaderGen ElfHeader64

// Class-dependent phdrs.
type ProgramHdr32 struct {

}

type ProgramHdr64 struct {

}

// Class-dependent section headers.
type SectionHeader32 struct {
	Sh_name uint32
	Sh_name_string string
	Sh_type elf.SectionType
	Sh_flags elf.SectionFlag
	Sh_addr uint32
	Sh_offset uint32
	Sh_size uint32
	Sh_link uint32
	Sh_info uint32
	Sh_addralign uint32
	Sh_entsize uint32
}

type SectionHeader64 struct {
	Sh_name uint32
	Sh_name_string string
	Sh_type elf.SectionType
	Sh_flags elf.SectionFlag
	Sh_addr uint64
	Sh_offset uint64
	Sh_size uint64
	Sh_link uint32
	Sh_info uint32
	Sh_addralign uint64
	Sh_entsize uint64
}

type SectionHeader SectionHeader64

type ElfFile struct {
	
}

type StringTable struct {

}

type SymbolTable struct {
	
}

type Elf32Rel struct {
	R_off uint32
	R_info uint32
}

type Elf32Rela struct {
	R_off uint32
	R_info uint32
	R_addend int32
}
