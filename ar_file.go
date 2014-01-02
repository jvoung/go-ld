// Copyright (c) 2013, Jan Voung
// All rights reserved.

// Representation of AR files.

package main

import (
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
	// TODO(jvoung): Assuming the files are ELF, wrap the contents as ELF?
}

type ARFile map[string] ARFileHeaderContents

func ReadPlainARFile(f *os.File) ARFile {
	ar_file := make(map[string] ARFileHeaderContents)
	per_file_header_size := 60
	hbuf := make([]byte, per_file_header_size)
	// Assume magic number header is already read.
	offset := int64(len(AR_MAGIC))
	var special_long_filename_file []byte
	// Go through the AR, reading more and more file-headers + file-bodies.
	for {
		n, err := f.ReadAt(hbuf, offset)
		if err == io.EOF {
			fmt.Printf("Hit eof, reading %d chars, which is %s\n", n,
				hbuf[0])
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
		// TODO(jvoung): Handle the / which is the terminator, instead of using
		// TrimSpace. Also use the / to handle long filenames.
		filename := strings.TrimSpace(string(hbuf[0:16]))
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
		fmt.Println("Read stuff for", new_header.String())
		if filename == "/" {
			fmt.Println("Copied", copy(special_long_filename_file, body_buf))
		} else {
			ar_file[filename] = ARFileHeaderContents{new_header, body_buf}
		}
		offset += int64(fsize)
	}
	fmt.Println("Special long filename section:",
		string(special_long_filename_file))
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
