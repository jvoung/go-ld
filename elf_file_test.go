// Copyright (c) 2014, Jan Voung
// All rights reserved.

// Test ELF file utilities.

package main

import (
	"debug/elf"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

func findSectionIndex(name string, f *ElfFile) int {
	for i := range f.Shdrs {
		if f.Shdrs[i].Sh_name == name {
			return i
		}
	}
	panic("Cannot find section w/ name: " + name)
}

func checkSymtabCrtbegin(t *testing.T, f *ElfFile,
	st SymbolTable, start_off int, start_size int) {
	// crtbegin.o should have:
	// *) U _pnacl_wrapper_start which does a few things and then calls _start
	// *) U __pnacl_irt_init which sets up the __nacl_read_tp function using
	//      the startup_info auxv.
	// *) T __pnacl_start, the entry point for PNaCl programs.
	sym, ok := st["_pnacl_wrapper_start"]
	ExpectEq(t, true, ok)
	ExpectEq(t, "_pnacl_wrapper_start", sym.St_name)
	ExpectEq(t, uint8(0), sym.St_other)
	ExpectEq(t, elf.SHN_UNDEF, sym.St_shndx)
	ExpectEq(t, uint64(0), sym.St_value)

	sym, ok = st["__pnacl_init_irt"]
	ExpectEq(t, true, ok)
	ExpectEq(t, "__pnacl_init_irt", sym.St_name)
	ExpectEq(t, uint8(0), sym.St_other)
	ExpectEq(t, elf.SHN_UNDEF, sym.St_shndx)
	ExpectEq(t, uint64(0), sym.St_value)

	text_index := findSectionIndex(".text", f)
	sym, ok = st["__pnacl_start"]
	ExpectEq(t, true, ok)
	ExpectEq(t, "__pnacl_start", sym.St_name)
	ExpectEq(t, uint8(0), sym.St_other)
	ExpectEq(t, elf.SectionIndex(text_index), sym.St_shndx)
	// Offset relative to the beginning of the file.
	ExpectEq(t, uint64(start_off), sym.St_value)
	ExpectEq(t, uint64(start_size), sym.St_size)
}

func TestRelocatableELFX8632(t *testing.T) {
	// Just using crtbegin.o for now.
	// Want to also test a .o coming from a .pexe.
	fname := path.Join(TestX8632BaseDir(), "crtbegin.o")
	elf_file := ReadElfFileFname(fname)
	ExpectEq(t, elf.ELFCLASS32, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	// ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
	ExpectEq(t, uint8(0), elf_file.Header.ABIVersion)
	ExpectEq(t, elf.ET_REL, elf_file.Header.Type)
	ExpectEq(t, elf.EM_386, elf_file.Header.Machine)
	ExpectEq(t, uint32(1), elf_file.Header.E_Version)
	ExpectEq(t, uint64(0), elf_file.Header.Entry)
	ExpectEq(t, uint64(0), elf_file.Header.Phoff)
	ExpectEq(t, uint64(452), elf_file.Header.Shoff)
	ExpectEq(t, uint32(0), elf_file.Header.Flags)
	ExpectEq(t, uint16(52), elf_file.Header.FileHeaderSize)
	ExpectEq(t, uint16(0), elf_file.Header.Phentsize)
	ExpectEq(t, uint16(0), elf_file.Header.Phnum)
	ExpectEq(t, uint16(40), elf_file.Header.Shentsize)
	ExpectEq(t, uint16(12), elf_file.Header.Shnum)
	ExpectEq(t, uint16(9), elf_file.Header.Shstrndx)
	ExpectEq(t, 0, len(elf_file.Phdrs))
	ExpectEq(t, 12, len(elf_file.Shdrs))
	ExpectEq(t,
		SectionHeader{Sh_name_index: 0, Sh_name: "", Sh_type: elf.SHT_NULL,
			Sh_flags: 0, Sh_addr: 0, Sh_offset: 0, Sh_size: 0,
			Sh_link: 0, Sh_info: 0,	Sh_addralign: 0, Sh_entsize: 0},
		elf_file.Shdrs[0])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 16, Sh_name: ".group",
			Sh_type: elf.SHT_GROUP,
			Sh_flags: 0, Sh_addr: 0, Sh_offset: 0x34, Sh_size: 8,
			Sh_link: 10, Sh_info: 2, Sh_addralign: 4, Sh_entsize: 4},
		elf_file.Shdrs[1])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 5, Sh_name: ".text",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_ALLOC | elf.SHF_EXECINSTR,
			Sh_addr: 0, Sh_offset: 0x40, Sh_size: 0x100,
			Sh_link: 0, Sh_info: 0,	Sh_addralign: 32, Sh_entsize: 0},
		elf_file.Shdrs[2])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 1, Sh_name: ".rel.text",
			Sh_type: elf.SHT_REL,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x5fc, Sh_size: 0x20,
			Sh_link: 10, Sh_info: 2, Sh_addralign: 4, Sh_entsize: 8},
		elf_file.Shdrs[3])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 75, Sh_name: ".data",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x140, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[4])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 11, Sh_name: ".bss",
			Sh_type: elf.SHT_NOBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x140, Sh_size: 0x40,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 8, Sh_entsize: 0},
		elf_file.Shdrs[5])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 81, Sh_name: ".note.NaCl.ABI.x86-32",
			Sh_type: elf.SHT_NOTE,
			Sh_flags: elf.SHF_ALLOC | elf.SHF_GROUP,
			Sh_addr: 0, Sh_offset: 0x140, Sh_size: 0x1c,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[6])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 39, Sh_name: ".eh_frame",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x15c, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[7])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 23, Sh_name: ".note.GNU-stack",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x15c, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[8])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 49, Sh_name: ".shstrtab",
			Sh_type: elf.SHT_STRTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x15c, Sh_size: 0x67,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[9])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 67, Sh_name: ".symtab",
			Sh_type: elf.SHT_SYMTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x3a4, Sh_size: 0x120,
			Sh_link: 11, Sh_info: 13, Sh_addralign: 4, Sh_entsize: 0x10},
		elf_file.Shdrs[10])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 59, Sh_name: ".strtab",
			Sh_type: elf.SHT_STRTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x4c4, Sh_size: 0x136,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[11])

	// Try reading the symbol table too.
	st := elf_file.ReadSymbols()
	ExpectEq(t, 10, len(st))

	// Check it more deeply.
	checkSymtabCrtbegin(t, &elf_file, st, 0x40, 128)
}

func TestRelocatableELFX8664(t *testing.T) {
	fname := path.Join(TestX8664BaseDir(), "crtbegin.o")
	elf_file := ReadElfFileFname(fname)
	// Will eventually be ELFCLASS32 also, for NaCl x86-64.
	ExpectEq(t, elf.ELFCLASS64, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	// ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
	ExpectEq(t, uint8(0), elf_file.Header.ABIVersion)
	ExpectEq(t, elf.ET_REL, elf_file.Header.Type)
	ExpectEq(t, elf.EM_X86_64, elf_file.Header.Machine)
	ExpectEq(t, uint32(1), elf_file.Header.E_Version)
	ExpectEq(t, uint64(0), elf_file.Header.Entry)
	ExpectEq(t, uint64(0), elf_file.Header.Phoff)
	ExpectEq(t, uint64(584), elf_file.Header.Shoff)
	ExpectEq(t, uint32(0), elf_file.Header.Flags)
	ExpectEq(t, uint16(64), elf_file.Header.FileHeaderSize)
	ExpectEq(t, uint16(0), elf_file.Header.Phentsize)
	ExpectEq(t, uint16(0), elf_file.Header.Phnum)
	ExpectEq(t, uint16(64), elf_file.Header.Shentsize)
	ExpectEq(t, uint16(12), elf_file.Header.Shnum)
	ExpectEq(t, uint16(9), elf_file.Header.Shstrndx)
	ExpectEq(t, 0, len(elf_file.Phdrs))
	ExpectEq(t, 12, len(elf_file.Shdrs))
	ExpectEq(t,
		SectionHeader{Sh_name_index: 0, Sh_name: "", Sh_type: elf.SHT_NULL,
			Sh_flags: 0, Sh_addr: 0, Sh_offset: 0, Sh_size: 0,
			Sh_link: 0, Sh_info: 0,	Sh_addralign: 0, Sh_entsize: 0},
		elf_file.Shdrs[0])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 17, Sh_name: ".group",
			Sh_type: elf.SHT_GROUP,
			Sh_flags: 0, Sh_addr: 0, Sh_offset: 0x40, Sh_size: 8,
			Sh_link: 10, Sh_info: 2, Sh_addralign: 4, Sh_entsize: 4},
		elf_file.Shdrs[1])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 6, Sh_name: ".text",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_ALLOC | elf.SHF_EXECINSTR,
			Sh_addr: 0, Sh_offset: 0x60, Sh_size: 0x160,
			Sh_link: 0, Sh_info: 0,	Sh_addralign: 32, Sh_entsize: 0},
		elf_file.Shdrs[2])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 1, Sh_name: ".rela.text",
			Sh_type: elf.SHT_RELA,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x830, Sh_size: 0xc0,
			Sh_link: 10, Sh_info: 2, Sh_addralign: 8, Sh_entsize: 0x18},
		elf_file.Shdrs[3])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 76, Sh_name: ".data",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x1c0, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[4])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 12, Sh_name: ".bss",
			Sh_type: elf.SHT_NOBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x1c0, Sh_size: 0x40,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 8, Sh_entsize: 0},
		elf_file.Shdrs[5])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 82, Sh_name: ".note.NaCl.ABI.x86-64",
			Sh_type: elf.SHT_NOTE,
			Sh_flags: elf.SHF_ALLOC | elf.SHF_GROUP,
			Sh_addr: 0, Sh_offset: 0x1c0, Sh_size: 0x1c,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[6])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 40, Sh_name: ".eh_frame",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x1dc, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[7])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 24, Sh_name: ".note.GNU-stack",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x1dc, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[8])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 50, Sh_name: ".shstrtab",
			Sh_type: elf.SHT_STRTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x1dc, Sh_size: 0x68,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[9])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 68, Sh_name: ".symtab",
			Sh_type: elf.SHT_SYMTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x548, Sh_size: 0x1b0,
			Sh_link: 11, Sh_info: 13, Sh_addralign: 8, Sh_entsize: 0x18},
		elf_file.Shdrs[10])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 60, Sh_name: ".strtab",
			Sh_type: elf.SHT_STRTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x6f8, Sh_size: 0x136,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[11])
	st := elf_file.ReadSymbols()
	ExpectEq(t, 10, len(st))

	// Check it more deeply.
	checkSymtabCrtbegin(t, &elf_file, st, 0x80, 160)
}

func TestRelocatableELFARM(t *testing.T) {
	fname := path.Join(TestARMBaseDir(), "crtbegin.o")
	elf_file := ReadElfFileFname(fname)
	ExpectEq(t, elf.ELFCLASS32, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	// ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
	ExpectEq(t, uint8(0), elf_file.Header.ABIVersion)
	ExpectEq(t, elf.ET_REL, elf_file.Header.Type)
	ExpectEq(t, elf.EM_ARM, elf_file.Header.Machine)
	ExpectEq(t, uint32(1), elf_file.Header.E_Version)
	ExpectEq(t, uint64(0), elf_file.Header.Entry)
	ExpectEq(t, uint64(0), elf_file.Header.Phoff)
	ExpectEq(t, uint64(460), elf_file.Header.Shoff)
	ExpectEq(t, uint32(0x5000000), elf_file.Header.Flags)
	ExpectEq(t, uint16(52), elf_file.Header.FileHeaderSize)
	ExpectEq(t, uint16(0), elf_file.Header.Phentsize)
	ExpectEq(t, uint16(0), elf_file.Header.Phnum)
	ExpectEq(t, uint16(40), elf_file.Header.Shentsize)
	ExpectEq(t, uint16(12), elf_file.Header.Shnum)
	ExpectEq(t, uint16(9), elf_file.Header.Shstrndx)
	ExpectEq(t, 0, len(elf_file.Phdrs))
	ExpectEq(t, 12, len(elf_file.Shdrs))
	ExpectEq(t,
		SectionHeader{Sh_name_index: 0, Sh_name: "", Sh_type: elf.SHT_NULL,
			Sh_flags: 0, Sh_addr: 0, Sh_offset: 0, Sh_size: 0,
			Sh_link: 0, Sh_info: 0,	Sh_addralign: 0, Sh_entsize: 0},
		elf_file.Shdrs[0])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 32, Sh_name: ".group",
			Sh_type: elf.SHT_GROUP,
			Sh_flags: 0, Sh_addr: 0, Sh_offset: 0x34, Sh_size: 8,
			Sh_link: 10, Sh_info: 5, Sh_addralign: 4, Sh_entsize: 4},
		elf_file.Shdrs[1])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 5, Sh_name: ".text",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_ALLOC | elf.SHF_EXECINSTR,
			Sh_addr: 0, Sh_offset: 0x40, Sh_size: 0xe8,
			Sh_link: 0, Sh_info: 0,	Sh_addralign: 16, Sh_entsize: 0},
		elf_file.Shdrs[2])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 1, Sh_name: ".rel.text",
			Sh_type: elf.SHT_REL,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x63c, Sh_size: 0x40,
			Sh_link: 10, Sh_info: 2, Sh_addralign: 4, Sh_entsize: 8},
		elf_file.Shdrs[3])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 94, Sh_name: ".data",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x128, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[4])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 11, Sh_name: ".bss",
			Sh_type: elf.SHT_NOBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x128, Sh_size: 0x40,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 8, Sh_entsize: 0},
		elf_file.Shdrs[5])
	SHT_ARM_ATTRIBUTES := elf.SectionType(0x70000003)
	ExpectEq(t,
		SectionHeader{Sh_name_index: 16, Sh_name: ".ARM.attributes",
			Sh_type: SHT_ARM_ATTRIBUTES,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x128, Sh_size: 0x26,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[6])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 39, Sh_name: ".note.NaCl.ABI.arm",
			Sh_type: elf.SHT_NOTE,
			Sh_flags: elf.SHF_ALLOC | elf.SHF_GROUP,
			Sh_addr: 0, Sh_offset: 0x150, Sh_size: 0x18,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[7])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 58, Sh_name: ".eh_frame",
			Sh_type: elf.SHT_PROGBITS,
			Sh_flags: elf.SHF_WRITE | elf.SHF_ALLOC,
			Sh_addr: 0, Sh_offset: 0x168, Sh_size: 0,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 4, Sh_entsize: 0},
		elf_file.Shdrs[8])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 68, Sh_name: ".shstrtab",
			Sh_type: elf.SHT_STRTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x168, Sh_size: 0x64,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[9])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 86, Sh_name: ".symtab",
			Sh_type: elf.SHT_SYMTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x3ac, Sh_size: 0x150,
			Sh_link: 11, Sh_info: 16, Sh_addralign: 4, Sh_entsize: 0x10},
		elf_file.Shdrs[10])
	ExpectEq(t,
		SectionHeader{Sh_name_index: 78, Sh_name: ".strtab",
			Sh_type: elf.SHT_STRTAB,
			Sh_flags: 0,
			Sh_addr: 0, Sh_offset: 0x4fc, Sh_size: 0x13f,
			Sh_link: 0, Sh_info: 0, Sh_addralign: 1, Sh_entsize: 0},
		elf_file.Shdrs[11])

	st := elf_file.ReadSymbols()
	ExpectEq(t, 13, len(st))
	// Check it more deeply.
	checkSymtabCrtbegin(t, &elf_file, st, 0x40, 88)
}

func alignTo(addr uint64, alignment uint64) uint64 {
	diff := (alignment - (addr % alignment)) % alignment
	return addr + diff
}

func checkExecutableX8632NaCl(t *testing.T, fname string) {
	elf_file := ReadElfFileFname(fname)
	ExpectEq(t, elf.ELFCLASS32, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	// ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
	ExpectEq(t, uint8(0), elf_file.Header.ABIVersion)
	ExpectEq(t, elf.ET_EXEC, elf_file.Header.Type)
	ExpectEq(t, elf.EM_386, elf_file.Header.Machine)
	ExpectEq(t, 6, len(elf_file.Phdrs))

	// Check Phdrs
	ExpectEq(t, elf_file.Phdrs[0].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[0].P_flags, elf.PF_R | elf.PF_X)
	ExpectEq(t, elf_file.Phdrs[0].P_offset, uint64(0x10000))
	ExpectEq(t, elf_file.Phdrs[0].P_vaddr, uint64(0x20000))
	ExpectEq(t, elf_file.Phdrs[0].P_paddr, uint64(0x20000))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[0].P_align, uint64(0x10000))

	ExpectEq(t, elf_file.Phdrs[1].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[1].P_flags, elf.PF_R)
	ExpectEq(t, elf_file.Phdrs[1].P_offset, uint64(0))
	ExpectEq(t, elf_file.Phdrs[1].P_vaddr, uint64(0x10020000))
	ExpectEq(t, elf_file.Phdrs[1].P_paddr, uint64(0x10020000))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[1].P_align, uint64(0x10000))

	ExpectEq(t, elf_file.Phdrs[2].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[2].P_flags, elf.PF_R | elf.PF_W)
	// relative to the size of the previous segment.
	ExpectEq(t, elf_file.Phdrs[2].P_offset,
		alignTo(elf_file.Phdrs[1].P_filesz, 32))
	ExpectEq(t, elf_file.Phdrs[2].P_vaddr,
		alignTo(uint64(0x10030000 + elf_file.Phdrs[1].P_filesz), 32))
	ExpectEq(t, elf_file.Phdrs[2].P_paddr,
		alignTo(uint64(0x10030000 + elf_file.Phdrs[1].P_filesz), 32))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[0].P_align, uint64(0x10000))
}

func checkExecutableX8664NaCl(t *testing.T, fname string) {
	elf_file := ReadElfFileFname(fname)
	ExpectEq(t, elf.ELFCLASS64, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	// ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
	ExpectEq(t, uint8(0), elf_file.Header.ABIVersion)
	ExpectEq(t, elf.ET_EXEC, elf_file.Header.Type)
	ExpectEq(t, elf.EM_X86_64, elf_file.Header.Machine)
	ExpectEq(t, 6, len(elf_file.Phdrs))

	// Check Phdrs
	ExpectEq(t, elf_file.Phdrs[0].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[0].P_flags, elf.PF_R | elf.PF_X)
	ExpectEq(t, elf_file.Phdrs[0].P_offset, uint64(0x10000))
	ExpectEq(t, elf_file.Phdrs[0].P_vaddr, uint64(0x20000))
	ExpectEq(t, elf_file.Phdrs[0].P_paddr, uint64(0x20000))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[0].P_align, uint64(0x10000))

	ExpectEq(t, elf_file.Phdrs[1].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[1].P_flags, elf.PF_R)
	ExpectEq(t, elf_file.Phdrs[1].P_offset, uint64(0))
	ExpectEq(t, elf_file.Phdrs[1].P_vaddr, uint64(0x10020000))
	ExpectEq(t, elf_file.Phdrs[1].P_paddr, uint64(0x10020000))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[1].P_align, uint64(0x10000))

	ExpectEq(t, elf_file.Phdrs[2].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[2].P_flags, elf.PF_R | elf.PF_W)
	// relative to the size of the previous segment.
	ExpectEq(t, elf_file.Phdrs[2].P_offset,
		alignTo(elf_file.Phdrs[1].P_filesz, 32))
	ExpectEq(t, elf_file.Phdrs[2].P_vaddr,
		alignTo(uint64(0x10030000 + elf_file.Phdrs[1].P_filesz), 32))
	ExpectEq(t, elf_file.Phdrs[2].P_paddr,
		alignTo(uint64(0x10030000 + elf_file.Phdrs[1].P_filesz), 32))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[0].P_align, uint64(0x10000))
}

func checkExecutableARMNaCl(t *testing.T, fname string) {
	elf_file := ReadElfFileFname(fname)
	ExpectEq(t, elf.ELFCLASS32, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	// ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
	ExpectEq(t, uint8(0), elf_file.Header.ABIVersion)
	ExpectEq(t, elf.ET_EXEC, elf_file.Header.Type)
	ExpectEq(t, elf.EM_ARM, elf_file.Header.Machine)
	// The ARM one doesn't have a GNU_STACK segment so it's only 5 segments.
	ExpectEq(t, 5, len(elf_file.Phdrs))

	// Check Phdrs
	ExpectEq(t, elf_file.Phdrs[0].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[0].P_flags, elf.PF_R | elf.PF_X)
	ExpectEq(t, elf_file.Phdrs[0].P_offset, uint64(0x10000))
	ExpectEq(t, elf_file.Phdrs[0].P_vaddr, uint64(0x20000))
	ExpectEq(t, elf_file.Phdrs[0].P_paddr, uint64(0x20000))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[0].P_align, uint64(0x10000))

	ExpectEq(t, elf_file.Phdrs[1].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[1].P_flags, elf.PF_R)
	ExpectEq(t, elf_file.Phdrs[1].P_offset, uint64(0))
	ExpectEq(t, elf_file.Phdrs[1].P_vaddr, uint64(0x10020000))
	ExpectEq(t, elf_file.Phdrs[1].P_paddr, uint64(0x10020000))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[1].P_align, uint64(0x10000))

	ExpectEq(t, elf_file.Phdrs[2].P_type, elf.PT_LOAD)
	ExpectEq(t, elf_file.Phdrs[2].P_flags, elf.PF_R | elf.PF_W)
	// relative to the size of the previous segment.
	// TODO(jvoung): what is the alignment requirement???
	ExpectEq(t, elf_file.Phdrs[2].P_offset,
		alignTo(elf_file.Phdrs[1].P_filesz, 8))
	ExpectEq(t, elf_file.Phdrs[2].P_vaddr,
		alignTo(uint64(0x10030000 + elf_file.Phdrs[1].P_filesz), 8))
	ExpectEq(t, elf_file.Phdrs[2].P_paddr,
		alignTo(uint64(0x10030000 + elf_file.Phdrs[1].P_filesz), 8))
	// Skip the sizes, because they can change.
	ExpectEq(t, elf_file.Phdrs[0].P_align, uint64(0x10000))
}

func lastModTime(files []string, can_skip bool) time.Time {
	t := time.Time{}
	for _, fname := range files {
		stat, err := os.Stat(fname)
		if err != nil {
      if can_skip {
        continue
      } else {
        panic("Failed to stat: " + fname)
      }
    }
		if stat.ModTime().After(t) {
			t = stat.ModTime()
		}
	}
	return t
}

func naclTestDataOld(infiles, outfiles []string) bool {
	max_in_mod := lastModTime(infiles, false)
	max_out_mod := lastModTime(outfiles, true)
	return max_in_mod.After(max_out_mod)
}

func TestNaClExecutable(t *testing.T) {
	// Use the test_binary shell script to generate a NaCl .nexe
	// then read it.
	infiles := []string{path.Join(TestBaseDir, "test_relocs.sh"),
		path.Join(TestBaseDir, "test_relocs.c")}
	outdirs := []string{TestX8632BaseDir(), TestX8664BaseDir(),
		TestARMBaseDir()}
	outfiles := []string{"test_relocs.o",
		"test_relocs.nexe",
		"test_relocs.nexe---test_relocs.final.pexe---.o"}
	joined_of := []string{}
	for _, od := range outdirs {
		for _, of := range outfiles {
			joined_of = append(joined_of, path.Join(od, of))
		}
	}
	if naclTestDataOld(infiles, joined_of) {
		fmt.Println("Need to regenerate relocs test binaries: test_relocs.sh")
		cmd := exec.Command(path.Join(TestBaseDir, "test_relocs.sh"))
		err := cmd.Run()
		if err != nil {
			t.Fatal(err)
		}
	}
	checkExecutableX8632NaCl(
		t, path.Join(TestX8632BaseDir(), "test_relocs.nexe"))
	checkExecutableX8664NaCl(
		t, path.Join(TestX8664BaseDir(), "test_relocs.nexe"))
	checkExecutableARMNaCl(
		t, path.Join(TestARMBaseDir(), "test_relocs.nexe"))
}
