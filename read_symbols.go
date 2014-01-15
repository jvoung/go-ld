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

type SymbolTable []SymbolTableEntry

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
	sizeof_struct := 0
	if f.Header.Class == elf.ELFCLASS32 {
		reader_func = readSymbolEntry32
		sizeof_struct = 16
	} else if f.Header.Class == elf.ELFCLASS64 {
		reader_func = readSymbolEntry64
		sizeof_struct = 24
	} else {
		panic("Unknown ELF class")
	}
	result := make([]SymbolTableEntry, 0,
		symtab_sec_hdr.Sh_size / uint64(sizeof_struct))
	for {
		new_st_entry, err := reader_func(byte_reader, byte_order, strtab_slice)
		if err == io.EOF {
			break
		}
		if new_st_entry.St_name == "" {
			// Can we omit the nameless symbols for now?
			continue
		}
		// Uh... we need to be careful about local symbols w/ the same name
		// (and local may have the same name as a global)!
		result = append(result, new_st_entry)
	}
	return result
}

type SymLinkInfo struct {
	// Index into symbol table for the symbol.
	UndefinedSyms []int	
	ExportedSyms []int
}

func GetSymBind(i uint8) elf.SymBind {
	return elf.SymBind(i >> 4)
}

func GetSymLinkInfo(st SymbolTable) SymLinkInfo {
	info := SymLinkInfo{}
	for i, sym := range st {
		if sym.St_shndx == elf.SHN_UNDEF {
			info.UndefinedSyms = append(info.UndefinedSyms, i)
		} else if GetSymBind(sym.St_info) == elf.STB_GLOBAL {
			info.ExportedSyms = append(info.ExportedSyms, i)
		}
	}
	return info
}

func SymLinkInfoToHash(link_info SymLinkInfo,
	st SymbolTable) map[string] *SymbolTableEntry {
	result := make(map[string] *SymbolTableEntry)
	// UndefinedSyms and ExportedSyms should be unique and not have
	// local symbols, so we can use a map[string] at this point.
	for _, index := range link_info.UndefinedSyms {
		result[st[index].St_name] = &st[index]
	}
	for _, index := range link_info.ExportedSyms {
		result[st[index].St_name] = &st[index]
	}
	return result
}
