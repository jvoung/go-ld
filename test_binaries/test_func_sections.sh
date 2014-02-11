#!/bin/bash

# Set up test binaries from test_relocs.c
# TODO(jvoung): Have a CHECK syntax like the LLVM-lit tests?

set -e
set -u
set -x

readonly arches="i686 x86_64 arm mips"
readonly TC_ROOT=${NACL_SDK_ROOT}/toolchain/linux_pnacl/bin

pushd test_binaries
for arch in ${arches}; do
  OUT=${arch}
  ${TC_ROOT}/pnacl-clang -c -arch ${arch} test_func_secs1.c -O1 \
    -o ${OUT}/test_func_secs1.o --pnacl-allow-translate -Wt,-ffunction-sections
  ${TC_ROOT}/pnacl-clang -c -arch ${arch} test_func_secs2.c -O1 \
    -o ${OUT}/test_func_secs2.o --pnacl-allow-translate -Wt,-ffunction-sections
  # TODO(jvoung): generate a --section-ordering-file?
  ${TC_ROOT}/pnacl-clang -arch ${arch} \
    ${OUT}/test_func_secs1.o \
    ${OUT}/test_func_secs2.o \
    -o ${OUT}/test_func_secs.nexe \
    -Wn,--section-ordering-file=test_func_secs_order.txt \
    --pnacl-allow-native \
    -ffunction-sections -save-temps --pnacl-driver-verbose
done
popd
