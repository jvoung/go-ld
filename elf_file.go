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
	"io/ioutil"
	"os"
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

func toByteOrder(d elf.Data) binary.ByteOrder {
    var b binary.ByteOrder
    if d == elf.ELFDATA2LSB {
        b = binary.LittleEndian
    } else if d == elf.ELFDATA2MSB {
        b = binary.BigEndian
    } else {
        panic("Unknown byte order")
    }
    return b
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

func ReadElfHeader(buf []byte) ElfFileHeader {
	class := elf.Class(buf[4])
	data := elf.Data(buf[5])
	ei_ver := elf.Version(buf[6])
	osabi := elf.OSABI(buf[7])
	abi_ver := uint8(buf[8])
	byte_order := toByteOrder(data)
	// Initialize part of the struct for now (the non-byte-order dependent bits)
	header := ElfFileHeader{
		Class: class,
		Data: data,
		EI_Version: ei_ver,
		OSABI: osabi,
		ABIVersion: abi_ver}
	byte_reader := bytes.NewReader(buf[16:])
	err1 := binary.Read(byte_reader, byte_order, &header.Type)
	err2 := binary.Read(byte_reader, byte_order, &header.Machine)
	err3 := binary.Read(byte_reader, byte_order, &header.E_Version)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read ELF machine")
	}
	header.Entry, header.Phoff, header.Shoff = ReadElfHeaderWithClass(
		byte_reader, class, byte_order)
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

func readPhdr32(buf []byte, byte_order binary.ByteOrder) ProgramHeader {
	byte_reader := bytes.NewReader(buf)
    phdr := ProgramHeader{}
	err1 := binary.Read(byte_reader, byte_order, &phdr.P_type)
	if err1 != nil {
		panic("Failed to read phdr type")
	}
	var offset, vaddr, paddr, filesz, memsz, flags, align uint32
	err1 = binary.Read(byte_reader, byte_order, &offset)
	err2 := binary.Read(byte_reader, byte_order, &vaddr)
	err3 := binary.Read(byte_reader, byte_order, &paddr)
	err4 := binary.Read(byte_reader, byte_order, &filesz)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read phdr offset, vaddr, paddr, or filesz")
	}
	err1 = binary.Read(byte_reader, byte_order, &memsz)
	err2 = binary.Read(byte_reader, byte_order, &flags)
	err3 = binary.Read(byte_reader, byte_order, &align)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read phdr memsz, flags, or align")
	}
	phdr.P_flags = elf.ProgFlag(flags)
	phdr.P_offset = uint64(offset)
	phdr.P_vaddr = uint64(vaddr)
	phdr.P_paddr = uint64(paddr)
	phdr.P_vaddr = uint64(vaddr)
	phdr.P_filesz = uint64(filesz)
	phdr.P_memsz = uint64(memsz)
	phdr.P_align = uint64(align)
	return phdr
}

func readPhdr64(buf []byte, byte_order binary.ByteOrder) ProgramHeader {
		byte_reader := bytes.NewReader(buf)
    phdr := ProgramHeader{}
	err1 := binary.Read(byte_reader, byte_order, &phdr.P_type)
	err2 := binary.Read(byte_reader, byte_order, &phdr.P_flags)
	err3 := binary.Read(byte_reader, byte_order, &phdr.P_offset)
	err4 := binary.Read(byte_reader, byte_order, &phdr.P_vaddr)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read phdr type, flags, offset, or vaddr")
	}
	err1 = binary.Read(byte_reader, byte_order, &phdr.P_paddr)
	err2 = binary.Read(byte_reader, byte_order, &phdr.P_filesz)
	err3 = binary.Read(byte_reader, byte_order, &phdr.P_memsz)
	err4 = binary.Read(byte_reader, byte_order, &phdr.P_align)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read paddr, filesz, memsz, or align")
	}
	return phdr
}

// Read in the program headers of the program.
func ReadProgramHeaders(
	buf []byte, fhdr *ElfFileHeader) []ProgramHeader {
    phdrs := make([]ProgramHeader, fhdr.Phnum)
    byte_order := toByteOrder(fhdr.Data)
	var reader_func func([]byte, binary.ByteOrder) ProgramHeader
	if fhdr.Class == elf.ELFCLASS32 {
		reader_func = readPhdr32
	} else if fhdr.Class == elf.ELFCLASS64 {
		reader_func = readPhdr64
	} else {
		panic("Unknown ELF class")
	}
	offset := fhdr.Phoff
	if offset == 0 {
		return phdrs
	}
    for i := 0; i < int(fhdr.Phnum); i++ {
		new_phdr := reader_func(
			buf[offset : offset + uint64(fhdr.Phentsize)], byte_order)
		phdrs = append(phdrs, new_phdr)
	}
	return phdrs
}

// Rounded-up Section Headers (elf class 32 and 64 layout is the same order)
type SectionHeader struct {
	Sh_name_index uint32
	Sh_name string
	Sh_type elf.SectionType
	Sh_flags elf.SectionFlag // 32 or 64
	Sh_addr uint64 // or 32
	Sh_offset uint64 // or 32
	Sh_size uint64 // or 32
	Sh_link uint32
	Sh_info uint32
	Sh_addralign uint64 // or 32
	Sh_entsize uint64 // or 32
}

// A string table is just an sequence null-terminated strings.
// The first element of the string table is null (used for non-existent names).
type StringTable []byte

func readShdr32(buf []byte, byte_order binary.ByteOrder) SectionHeader {
	byte_reader := bytes.NewReader(buf)
    shdr := SectionHeader{}
	err1 := binary.Read(byte_reader, byte_order, &shdr.Sh_name_index)
	err2 := binary.Read(byte_reader, byte_order, &shdr.Sh_type)
	if err1 != nil || err2 != nil {
		panic("Failed to read shdr name-index, or type")
	}
	var flags, addr, offset, size uint32
	err1 = binary.Read(byte_reader, byte_order, &flags)
	err2 = binary.Read(byte_reader, byte_order, &addr)
	err3 := binary.Read(byte_reader, byte_order, &offset)
	err4 := binary.Read(byte_reader, byte_order, &size)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read shdr flags, addr, offset, or size")
	}
	shdr.Sh_flags = elf.SectionFlag(flags)
	shdr.Sh_addr = uint64(addr)
	shdr.Sh_offset = uint64(offset)
	shdr.Sh_size = uint64(size)
	var addralign, entsize uint32
	err1 = binary.Read(byte_reader, byte_order, &shdr.Sh_link)
	err2 = binary.Read(byte_reader, byte_order, &shdr.Sh_info)
	err3 = binary.Read(byte_reader, byte_order, &addralign)
	err4 = binary.Read(byte_reader, byte_order, &entsize)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read shdr link, info, addralign, or entsize")
	}
	shdr.Sh_addralign = uint64(addralign)
	shdr.Sh_entsize = uint64(entsize)
	return shdr
}

func readShdr64(buf []byte, byte_order binary.ByteOrder) SectionHeader {
	byte_reader := bytes.NewReader(buf)
	shdr := SectionHeader{}
	err1 := binary.Read(byte_reader, byte_order, &shdr.Sh_name_index)
	err2 := binary.Read(byte_reader, byte_order, &shdr.Sh_type)
	if err1 != nil || err2 != nil {
		panic("Failed to read shdr name-index, or type")
	}
	err1 = binary.Read(byte_reader, byte_order, &shdr.Sh_flags)
	err2 = binary.Read(byte_reader, byte_order, &shdr.Sh_addr)
	err3 := binary.Read(byte_reader, byte_order, &shdr.Sh_offset)
	err4 := binary.Read(byte_reader, byte_order, &shdr.Sh_size)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read shdr flags, addr, offset, or size")
	}
	err1 = binary.Read(byte_reader, byte_order, &shdr.Sh_link)
	err2 = binary.Read(byte_reader, byte_order, &shdr.Sh_info)
	err3 = binary.Read(byte_reader, byte_order, &shdr.Sh_addralign)
	err4 = binary.Read(byte_reader, byte_order, &shdr.Sh_entsize)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		panic("Failed to read shdr link, info, addralign, or entsize")
	}
	return shdr
}

func ReadSectionHeaders(buf []byte, fhdr *ElfFileHeader) []SectionHeader {
    shdrs := make([]SectionHeader, fhdr.Shnum)
    byte_order := toByteOrder(fhdr.Data)
	var reader_func func([]byte, binary.ByteOrder) SectionHeader
	if fhdr.Class == elf.ELFCLASS32 {
		reader_func = readShdr32
	} else if fhdr.Class == elf.ELFCLASS64 {
		reader_func = readShdr64
	} else {
		panic("Unknown ELF class")
	}
	offset := fhdr.Shoff
	if offset == 0 {
		return shdrs
	}
    for i := 0; i < int(fhdr.Shnum); i++ {
		new_shdr := reader_func(
			buf[offset : offset + uint64(fhdr.Shentsize)], byte_order)
		shdrs = append(shdrs, new_shdr)
	}
	// Also read the section header string table and fill out
	// the section names.
	sh_strtab_hdr := shdrs[fhdr.Shstrndx]
	sh_strtab := buf[
		sh_strtab_hdr.Sh_offset:sh_strtab_hdr.Sh_offset+sh_strtab_hdr.Sh_size]
	for _, shdr := range shdrs {
		name_start := shdr.Sh_name_index
		if name_start == 0 {
			continue
		}
		name_end := bytes.IndexByte(sh_strtab[name_start:], 0)
		shdr.Sh_name = string(sh_strtab[name_start:name_end])
	}
	return shdrs
}

type SymbolTable struct {
	// TODO(jvoung): fill this...
}

type ElfFile struct {
	Header ElfFileHeader
	Phdrs []ProgramHeader
	Shdrs []SectionHeader
}

func ReadElfFile(buf []byte) ElfFile {
	result := ElfFile{}
	result.Header = ReadElfHeader(buf)
	result.Phdrs = ReadProgramHeaders(buf, &result.Header)
	result.Shdrs = ReadSectionHeaders(buf, &result.Header)
	return result
}

// For testing.
func ReadElfFileFname(fname string) ElfFile {
	f, err := os.Open(fname)
	if err != nil {
		panic("Failed to open file: " + string(fname) +
			" error: " + err.Error())
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Failed to open file: " + string(fname) +
			" error: " + err.Error())
	}
	return ReadElfFile(body)
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
