// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Representation of ELF file (separate from go's debug/elf,
// modulo some constants and printing functions).

package main

import (
	"bytes"
	"encoding/binary"
	"debug/elf"
	"fmt"
	"io"
)

// Rounded-up ELF header (elf class 32 and 64 layout is the same order)
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
	Entry uint64 // Offset 0x18 (could have been uint32)
	Phoff uint64 // (could have been uint32)
	Shoff uint64 // (could have been uint32)
	Flags uint32
	FileHeaderSize uint16
	Phentsize uint16
	Phnum uint16
	Shentsize uint16
	Shnum uint16
	Shstrndx uint16
}

func (h *ElfFileHeader) String() string {
	return fmt.Sprintf("ELF Header:\n" +
		"  Class: %s\n" +
		"  Data: %s\n" +
		"  Version: %s\n" +
		"  OS/ABI: %s\n" +
		"  ABI Version: %d\n" +
		"  Type: %s\n" +
		"  Machine: %s\n" +
		"  Version: %d\n" +
		"  Entry point address: 0x%x\n" +
		"  Start of program headers: %d (bytes into file)\n" +
		"  Start of section headers: %d (bytes into file)\n" +
		"  Flags: 0x%x\n" +
		"  Size of this header: %d\n" +
		"  Size of program headers: %d\n" +
		"  Number of program headers: %d\n" +
		"  Size of section headers: %d\n" +
		"  Number of section headers: %d\n" +
		"  Section header string table index: %d\n",
		h.Class, h.Data, h.EI_Version, h.OSABI, h.ABIVersion, h.Type,
		h.Machine, h.E_Version, h.Entry, h.Phoff, h.Shoff, h.Flags,
        h.FileHeaderSize, h.Phentsize, h.Phnum,
        h.Shentsize, h.Shnum, h.Shstrndx)
}

// Read portion of the ELF file header that depends on the ELF class,
// returning the rounded-up fields.
func ReadElfHeaderWithClass(
	byte_reader io.Reader, class elf.Class, byte_order binary.ByteOrder) (
	entry uint64, phoff uint64, shoff uint64) {
	switch class {
	default: panic("Unknown ELF class " + string(class))
	case elf.ELFCLASS32:
		var e32, ph32, sh32 uint32
		err1 := binary.Read(byte_reader, byte_order, &e32)
		err2 := binary.Read(byte_reader, byte_order, &ph32)
		err3 := binary.Read(byte_reader, byte_order, &sh32)
		if err1 != nil || err2 != nil || err3 != nil {
			panic("Failed to read ELF header")
		}
		return uint64(e32), uint64(ph32), uint64(sh32)
	case elf.ELFCLASS64:
		err1 := binary.Read(byte_reader, byte_order, &entry)
		err2 := binary.Read(byte_reader, byte_order, &phoff)
		err3 := binary.Read(byte_reader, byte_order, &shoff)
		if err1 != nil || err2 != nil || err3 != nil {
			panic("Failed to read ELF header")
		}
		return entry, phoff, shoff
	}
	return
}

func ReadElfHeader(f_buf []byte) ElfFileHeader {
	class := elf.Class(f_buf[4])
	data := elf.Data(f_buf[5])
	ei_ver := elf.Version(f_buf[6])
	osabi := elf.OSABI(f_buf[7])
	abi_ver := uint8(f_buf[8])
	var byte_order binary.ByteOrder
	switch data {
	case elf.ELFDATA2LSB:
		byte_order = binary.LittleEndian
	case elf.ELFDATA2MSB:
		byte_order = binary.BigEndian
	default:
		panic("Unknown byte order")
	}
	byte_reader := bytes.NewReader(f_buf[16:])
	var typ elf.Type
	var machine elf.Machine
	var e_ver uint32
	err1 := binary.Read(byte_reader, byte_order, &typ)
	err2 := binary.Read(byte_reader, byte_order, &machine)
	err3 := binary.Read(byte_reader, byte_order, &e_ver)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read ELF machine")
	}
	entry, phoff, shoff := ReadElfHeaderWithClass(
		byte_reader, class, byte_order)
	header := ElfFileHeader{
		Class: class,
		Data: data,
		EI_Version: ei_ver,
		OSABI: osabi,
		ABIVersion: abi_ver,
		Type: typ,
		Machine: machine,
		E_Version: e_ver,
		Entry: entry,
		Phoff: phoff,
		Shoff: shoff }
	err1 = binary.Read(byte_reader, byte_order, &header.Flags)
	err2 = binary.Read(byte_reader, byte_order, &header.FileHeaderSize)
	err3 = binary.Read(byte_reader, byte_order, &header.Phentsize)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read ELF machine")
	}
    err1 = binary.Read(byte_reader, byte_order, &header.Phnum)
	err2 = binary.Read(byte_reader, byte_order, &header.Shentsize)
	err3 = binary.Read(byte_reader, byte_order, &header.Shnum)
	err4 := binary.Read(byte_reader, byte_order, &header.Shstrndx)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read ELF machine")
	}
	return header
}

// PHDRs (elf class 32 and 64 have slightly different layout...)
type ProgramHeader struct {
	P_type elf.ProgType
	P_flags elf.ProgFlag // flags after memsz for elf-class 32
	P_offset uint64
	P_vaddr uint64
	P_paddr uint64
	P_filesz uint64
	P_memsz uint64
	P_align uint64
}

// For a generic program header, use debug/elf's ProgHeader
func ReadProgramHeaders(
	f_buf []byte, fhdr *ElfFileHeader) []ProgramHeader {
	if fhdr.Class == elf.ELFCLASS32 {
		return []ProgramHeader{}
	}
	return []ProgramHeader{}
}

// Rounded-up Section Headers (elf class 32 and 64 layout is the same order)
type SectionHeader struct {
	Sh_name_index uint32
	Sh_name string
	Sh_type elf.SectionType
	Sh_flags elf.SectionFlag
	Sh_addr uint64 // or 32
	Sh_offset uint64 // or 32
	Sh_size uint64 // or 32
	Sh_link uint32
	Sh_info uint32
	Sh_addralign uint64 // or 32
	Sh_entsize uint64 // or 32
}

func ReadSectionHeaders(f_buf []byte, fhdr *ElfFileHeader) []SectionHeader {
	return []SectionHeader{}
}

type StringTable struct {

}

type SymbolTable struct {
	
}

type ElfFile struct {
	Header ElfFileHeader
	Phdrs []ProgramHeader
	Shdrs []SectionHeader
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

type Elf64Rel struct {
	R_off uint64
	R_info uint32
}

type Elf64Rela struct {
	R_off uint64
	R_info uint32
	R_addend int64
}
