#!/bin/bash

# Set up test binaries from test_relocs.c
# TODO(jvoung): Have a CHECK syntax like the LLVM-lit tests?

set -e
set -u
set -x

readonly SRC=test_binaries/test_relocs.c
readonly arches="i686 x86_64 arm"

pnacl-clang $SRC -O1 -o /tmp/test_relocs.pexe
pnacl-finalize /tmp/test_relocs.pexe -o /tmp/test_relocs.final.pexe
for arch in ${arches}; do
  OUT=test_binaries/${arch}
  pnacl-clang -c -arch ${arch} ${SRC} -O1 -o ${OUT}/test_relocs.o --pnacl-allow-translate
  pnacl-translate -arch ${arch} /tmp/test_relocs.final.pexe -o ${OUT}/test_relocs.nexe -save-temps --pnacl-driver-verbose
done
