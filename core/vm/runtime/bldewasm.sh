#! /bin/bash

FILE=$1
cat <<EOF > /tmp/${FILE}.cpp
#include <stdint.h>

extern "C" void finish(char* _off, int32_t _len) __attribute__((import_module("ethereum"), import_name("finish")));

/*
extern "C"
int64_t max(int64_t a, int64_t b) {
   return a > b ? a: b;
}
*/

static	char	ret[8]={0,0,0,0, 0,0,0,10};
extern "C"
void Main() //__attribute__((export_name("main")))
{
	finish(ret, 8);
}
EOF

clang++ -c -O3 --target=wasm32 /tmp/${FILE}.cpp
wasm-ld --no-entry --allow-undefined-file=ewasm.syms --export=Main --strip-all ${FILE}.o -o /tmp/${FILE}.wasm
rm -f ${FILE}.o
wasm-dis /tmp/${FILE}.wasm | sed -s 's/Main/main/' > /tmp/${FILE}.wat
wasm-as /tmp/${FILE}.wat -o ${FILE}.wasm
