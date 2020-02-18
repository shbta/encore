#! /bin/bash
#
# Requires llvm 8.0.0 or higher
# llvm 7 does not support import_module/import_name attributes
# llvm 8 may export __heap_base and __data_end
#

FILE=$1
cat <<EOF > /tmp/${FILE}.c
#include <ewasm/ewasm.h>

static u32 fib(u32 n) {
	if (n < 2) return n;
	return fib(n-1)+fib(n-2);
}

static	byte	ret[8]={0,0,0,0, 0,0,0,10};
void main() // __attribute__((export_name("main")))
{
	i32	in_len;
	if ((in_len=eth_getCallDataSize()) < 4) eth_finish(ret, 8);
	u32 	met;
	u32 n = 10;
	if (in_len >= 36) {
		eth_callDataCopy(&met, 32, 4);
		n = __builtin_bswap32(met);
	}
	u32 res = __builtin_bswap32(fib(n));
	*(u32 *)(ret+4) = res;
	eth_finish(ret,8);
}
EOF

clang -c -O3 -Wall -I${HOME}/opt/ewasm/include -Wno-main-return-type --target=wasm32 /tmp/${FILE}.c
wasm-ld --no-entry --allow-undefined-file=ewasm.syms --export=main --strip-all ${FILE}.o -o ${FILE}.wasm
rm -f ${FILE}.o
#wasm-dis /tmp/${FILE}.wasm | sed -s 's/Main/main/' > /tmp/${FILE}.wat
#wasm-as /tmp/${FILE}.wat -o ${FILE}.wasm
