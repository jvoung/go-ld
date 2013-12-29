go-ld
=====

Linker Exercise: A very simple ELF linker that links object files and archives
into a static executable. This is just an exercise to dive into ELF and
the various sections / segments. May just start with a "readelf" type
of utility to begin with, instead of an actual linker.

Will use the constants from "debug/elf", but not much else so that I get
a better sense of the file formats.

- Start with just linking X86_32 ELF files.
- Expand to other architectures.


It really is just the basics. It does not handle:

- TLS
- .init, .fini, .init_array, .fini_array
- .eh_*
- gc-sections
- identical code folding
- etc.