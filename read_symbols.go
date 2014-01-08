// Copyright (c) 2014, Jan Voung
// All rights reserved.

// Functions to read symbol table entries.
// Depends on elf_file.go.

package main

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"io"
)

type SymbolTableEntry struct {
	St_name_index uint32
	St_name string
	St_info uint8
	St_other uint8
	// read as a uint16
	// TODO(jvoung): handle files w/ many sections
	St_shndx elf.SectionIndex
	St_value uint64 // or uint32
	St_size uint64 // or uint32
}

type SymbolTable map[string] SymbolTableEntry

func readSymbolEntryPrefix(r io.Reader, bo binary.ByteOrder) (
	SymbolTableEntry, error) {
	st_entry := SymbolTableEntry{}
	err1 := binary.Read(r, bo, &st_entry.St_name_index)
	if err1 == io.EOF {
		return st_entry, err1
	} else if err1 != nil {
		panic("Failed to read st_name")
	}
	return st_entry, nil
}

func readSymbolEntry32(r io.Reader, bo binary.ByteOrder, strtab []byte) (
	SymbolTableEntry, error) {
	st_entry, err1 := readSymbolEntryPrefix(r, bo)
	if err1 == io.EOF {
		return st_entry, err1
	} else if err1 != nil {
		panic("Failed to read symbol table prefix")
	}
	var value, size uint32
	err1 = binary.Read(r, bo, &value)
	err2 := binary.Read(r, bo, &size)
	if err1 != nil || err2 != nil {
		panic("Failed to read st_value, size")
	}
	st_entry.St_value = uint64(value)
	st_entry.St_size = uint64(size)
	var shndx uint16
	err1 = binary.Read(r, bo, &st_entry.St_info)
	err2 = binary.Read(r, bo, &st_entry.St_other)
	err3 := binary.Read(r, bo, &shndx)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read st_info, other, or shndx")
	}
	st_entry.St_shndx = elf.SectionIndex(shndx)
	st_entry.St_name = StringFromStrtab(strtab, st_entry.St_name_index)
	return st_entry, nil
}

func readSymbolEntry64(r io.Reader, bo binary.ByteOrder, strtab []byte) (
	SymbolTableEntry, error) {
	st_entry, err1 := readSymbolEntryPrefix(r, bo)
    if err1 == io.EOF {
		return st_entry, err1
	} else if err1 != nil {
		panic("Failed to read symbol table prefix")
	}
	var shndx uint16
	err1 = binary.Read(r, bo, &st_entry.St_info)
	err2 := binary.Read(r, bo, &st_entry.St_other)
	err3 := binary.Read(r, bo, &shndx)
	if err1 != nil || err2 != nil || err3 != nil {
		panic("Failed to read st_info, other, or shndx")
	}
	st_entry.St_shndx = elf.SectionIndex(shndx)
	err1 = binary.Read(r, bo, &st_entry.St_value)
	err2 = binary.Read(r, bo, &st_entry.St_size)
	if err1 != nil || err2 != nil {
		panic("Failed to read st_value, size")
	}
	st_entry.St_name = StringFromStrtab(strtab, st_entry.St_name_index)
	return st_entry, nil
}

// Reads all the symbol-table entries from the ElfFile,
// and figures out all the actual symbol names from the string table.
func (f ElfFile) ReadSymbols() SymbolTable {
	result := make(map[string] SymbolTableEntry)
	st_index := -1
	for i := range f.Shdrs {
		if f.Shdrs[i].Sh_name == ".symtab" &&
			f.Shdrs[i].Sh_type == elf.SHT_SYMTAB {
			st_index = i
			break
		}
	}
	if st_index == -1 {
		panic("No symbol table!")
	}
	symtab_sec_hdr := f.Shdrs[st_index]
	symtab_slice := f.Body[symtab_sec_hdr.Sh_offset:
		symtab_sec_hdr.Sh_offset + symtab_sec_hdr.Sh_size]
	strtab_sec_hdr := f.Shdrs[symtab_sec_hdr.Sh_link]
	strtab_slice := f.Body[strtab_sec_hdr.Sh_offset:
		strtab_sec_hdr.Sh_offset + strtab_sec_hdr.Sh_size]
	byte_reader := bytes.NewReader(symtab_slice)
	byte_order := ToByteOrder(f.Header.Data)
	var reader_func func(io.Reader, binary.ByteOrder, []byte) (
		SymbolTableEntry, error)
	if f.Header.Class == elf.ELFCLASS32 {
		reader_func = readSymbolEntry32
	} else if f.Header.Class == elf.ELFCLASS64 {
		reader_func = readSymbolEntry64
	} else {
		panic("Unknown ELF class")
	}
	for {
		new_st_entry, err := reader_func(byte_reader, byte_order, strtab_slice)
		if err == io.EOF {
			break
		}
		if new_st_entry.St_name == "" {
			// Can we omit the nameless symbols for now?
			continue
		}
		result[new_st_entry.St_name] = new_st_entry
	}
	return result
}
