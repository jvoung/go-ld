// Copyright (c) 2014, Jan Voung
// All rights reserved.

// Test ELF file utilities.

package main

import (
	"debug/elf"
	"path"
	"testing"
)

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
	ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
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
	ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
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
}

func TestRelocatableELFARM(t *testing.T) {
	fname := path.Join(TestARMBaseDir(), "crtbegin.o")
	elf_file := ReadElfFileFname(fname)
	ExpectEq(t, elf.ELFCLASS32, elf_file.Header.Class)
	ExpectEq(t, elf.ELFDATA2LSB, elf_file.Header.Data)
	ExpectEq(t, elf.EV_CURRENT, elf_file.Header.EI_Version)
	// Currently still built with the NaCl OSABI...
	// Will eventually switch to NONE.
	ExpectEq(t, elf.OSABI(123), elf_file.Header.OSABI)
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
}

// Check a particular executable ELF file for a particular architecture
// and ELF class.
func CheckExecutableELF(t *testing.T) {

}

func TestExecutableELF(t *testing.T) {

}
