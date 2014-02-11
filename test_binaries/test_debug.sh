#!/bin/bash

# Set up test binaries from test_relocs.c
# TODO(jvoung): Have a CHECK syntax like the LLVM-lit tests?

set -e
set -u
set -x

readonly arches="i686 x86_64 arm mips"
readonly TC_ROOT=${NACL_SDK_ROOT}/toolchain/linux_pnacl/bin

pushd test_binaries
echo "pwd: $(pwd)"

${TC_ROOT}/pnacl-clang++ -mllvm -inline-threshold=0 -O1 -g \
  test_debug1.cc test_debug2.cc test_debug3.cc \
  -o /tmp/test_debug.pexe
for arch in ${arches}; do
  OUT=${arch}
  ${TC_ROOT}/pnacl-translate -arch ${arch} /tmp/test_debug.pexe \
    -o ${OUT}/test_debug.nexe -save-temps --pnacl-driver-verbose \
    --allow-llvm-bitcode-input
done
popd
