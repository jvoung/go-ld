// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Representation of AR files.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type ARFileHeader struct {
	Filename string // Offset 0, up to 16 chars for short names.
	Timestamp string // Offset 16
	OwnerID string // Offset 28
	GroupID string // Offset 34
	FileMode string // Offset 40
	FileSize uint32 // Offset 48
}

func (h *ARFileHeader) String() string {
	return fmt.Sprintf("{Filename: %s\n Timestamp: %s\n OwnerID: %s\n " +
		"GroupID: %s\n FileMode: %s\n FileSize: %d}\n",
		h.Filename, h.Timestamp, h.OwnerID, h.GroupID,
		h.FileMode, h.FileSize)
}

type ARFileHeaderContents struct {
	Header ARFileHeader
	Contents []byte
}

type ARFile map[string] ARFileHeaderContents

// Specialized AR file holding only ELF files.
type ARElfFile struct {
	Header ARFileHeader
	File ElfFile
}

func translateFilename(fname string, lf_file []byte) string {
	if fname[0] == '/' {
		// It's a long filename, which is /[0-9]+, or one of the special files.
		fname = strings.TrimSpace(fname)
		if fname == "/" || fname == "//" {
			return fname
		}
		offset, err := strconv.Atoi(fname[1:])
		if err != nil {
			panic("Failed to parse long filename offset " + err.Error())
		}
		end := bytes.IndexByte(lf_file[offset:], '/')
		return string(lf_file[offset : offset + end])
	} else {
		end := strings.IndexByte(fname, '/')
		return fname[:end]
	}
}

func ReadPlainARFile(f *os.File) ARFile {
	ar_file := make(map[string] ARFileHeaderContents)
	per_file_header_size := 60
	hbuf := make([]byte, per_file_header_size)
	// Assume magic number header is already read.
	offset := int64(len(AR_MAGIC))
	special_long_filename_file := make([]byte, 0)
	// Go through the AR, reading more and more file-headers + file-bodies.
	for {
		n, err := f.ReadAt(hbuf, offset)
		if err == io.EOF && n == 0 {
			break
		}
		if err != nil {
			panic("Failed to read AR sub-file header: " + err.Error() + 
				" reading " + string(n))
		}
		// Okay, hbuf now has the header contents.
		offset += int64(n)
		fsize, err := strconv.Atoi(strings.TrimSpace(string(hbuf[48:58])))
		if err != nil {
			panic("Failed to parse AR file size: " + err.Error())
		}
		filename := translateFilename(string(hbuf[0:16]),
			special_long_filename_file)
		new_header := ARFileHeader {
			Filename: filename,
			Timestamp: strings.TrimSpace(string(hbuf[16:28])),
			OwnerID: strings.TrimSpace(string(hbuf[28:34])),
			GroupID: strings.TrimSpace(string(hbuf[34:40])),
			FileMode: strings.TrimSpace(string(hbuf[40:48])),
			FileSize: uint32(fsize) }
		body_buf := make([]byte, fsize)
		_, err2 := f.ReadAt(body_buf, offset)
		if err2 != nil {
			panic("Failed to read AR sub-file contents: " + err2.Error())
		}
		if filename == "/" {
			// Skipping the special GNU symbol-table file.
			// (not adding it to the ar_file map)
		} else if filename == "//" {
			// This is the long-filename file.
			// (not adding it to the ar_file map)
			special_long_filename_file = append(special_long_filename_file,
				body_buf...)
		} else {
			// Normal file, index it!
			ar_file[filename] = ARFileHeaderContents{new_header, body_buf}
		}
		offset += int64(fsize)
		// Data section should be aligned to 2 bytes.
		if offset % 2 != 0 {
			offset += 1
		}
	}
	return ar_file
}

func ReadARFile(f *os.File, typ FileType) ARFile {
	switch typ {
	case AR_FILE:
		return ReadPlainARFile(f)
	case THIN_AR_FILE:
		panic("TODO(jvoung): Handle thin archives")
	default:
		panic("Unknown AR file type: " + string(typ))
	}
}

func (f *ARFile) WrapARElf() map[string] ARElfFile {
	result := make(map[string] ARElfFile, len(*f))
	for fname, arsubfile := range(*f) {
		result[fname] = ARElfFile{
			Header: arsubfile.Header,
			File: ReadElfFile(arsubfile.Contents)}
	}
	return result
}
