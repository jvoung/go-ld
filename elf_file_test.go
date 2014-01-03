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
}

// Check a particular executable ELF file for a particular architecture
// and ELF class.
func CheckExecutableELF(t *testing.T) {

}

func TestExecutableELF(t *testing.T) {

}
