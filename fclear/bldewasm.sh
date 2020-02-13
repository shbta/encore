#! /bin/bash
#
# Requires llvm 8.0.0 or higher
# llvm 7 does not support import_module/import_name attributes
# llvm 8 may export __heap_base and __data_end
#

FILE=$1

clang -c -O3 -Wall -Wno-main-return-type --target=wasm32 ${FILE}.c
wasm-ld --no-entry --allow-undefined-file=ewasm.syms --export=main --strip-all ${FILE}.o -o ${FILE}.wasm
rm -f ${FILE}.o
#wasm-dis /tmp/${FILE}.wasm | sed -s 's/Main/main/' > /tmp/${FILE}.wat
#wasm-as /tmp/${FILE}.wat -o ${FILE}.wasm
