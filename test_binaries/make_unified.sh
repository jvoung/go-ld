#!/bin/bash

# Create a mondo native CRT file w/ everything.
# Then use that to link the test.o file into a nexe.
# Requires having run test_relocs.sh first.
set -e
set -u
set -x

readonly arches="i686 x86_64 arm"
readonly LD=${NACL_SDK_ROOT}/toolchain/linux_pnacl/host_x86_64/bin/le32-nacl-ld.gold

for arch in ${arches}; do
  OUT=test_binaries/${arch}
  IN=test_binaries/${arch}
  ${LD} -r --whole-archive -o ${OUT}/crtall.o \
      ${IN}/crtbegin.o ${IN}/libcrt_platform.a \
      ${IN}/libgcc.a ${IN}/libpnacl_irt_shim.a ${IN}/crtend.o

  # Now try to link the test_relocs.o w/ this unified crtall.o to make a nexe.
  case ${arch} in
      arm)
          EMUL=armelf_nacl ;;
      i686)
          EMUL=elf_nacl ;;
      x86_64)
          EMUL=elf64_nacl ;;
  esac
  ${LD} -nostdlib -m ${EMUL} --eh-frame-hdr --static --build-id \
      ${OUT}/crtall.o \
      ${IN}/test_relocs.nexe---test_relocs.final.pexe---.o \
      -o ${OUT}/test_relocs.unified.nexe
done
