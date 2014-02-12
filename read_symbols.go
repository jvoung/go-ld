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
		result = append(result, new_st_entry)
	}
	return result
}

// Set of symbol table indices.
type IndexSet map[int] bool

// An undefined symbol is from fileA and resolves to a defined
// symbol in fileB. Represent fileB with an int and the other symbol
// with another int.
type Resolver struct {
  DefFileIndex int
  DefSymIndex int
}
type UndefResolveMap map[int] Resolver

type SymLinkInfo struct {
	// Index into symbol table for the symbol.
    // These are sets (but use a map to represent that).
	UndefinedSyms UndefResolveMap
	ExportedSyms IndexSet
    ExportedSymHash map[string] int
}

func GetSymBind(i uint8) elf.SymBind {
	return elf.SymBind(i >> 4)
}

func GetSymLinkInfo(st SymbolTable) SymLinkInfo {
	info := SymLinkInfo{ make(map[int] Resolver, 0),
                         make(map[int] bool, 0),
                         make(map[string] int, 0) }
	for i, sym := range st {
        // Symbol at index 0 is always UNDEF and w/out a name.
        if i == 0 {
            continue
        }
		if sym.St_shndx == elf.SHN_UNDEF {
			info.UndefinedSyms[i] = Resolver{}
		} else if GetSymBind(sym.St_info) == elf.STB_GLOBAL {
			info.ExportedSyms[i] = true
            info.ExportedSymHash[sym.St_name] = i
        }
	}
	return info
}

func SymLinkInfoToHash(link_info SymLinkInfo,
	st SymbolTable) map[string] *SymbolTableEntry {
	result := make(map[string] *SymbolTableEntry)
	// UndefinedSyms and ExportedSyms should be unique and not have
	// local symbols, so we can use a map[string] at this point.
	for index := range link_info.UndefinedSyms {
		result[st[index].St_name] = &st[index]
	}
	for index := range link_info.ExportedSyms {
		result[st[index].St_name] = &st[index]
	}
	return result
}
