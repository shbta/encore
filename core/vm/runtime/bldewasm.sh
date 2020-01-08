#! /bin/bash
#
# Requires llvm 8.0.0 or higher
# llvm 7 does not support import_module/import_name attributes
# llvm 8 may export __heap_base and __data_end
#

FILE=$1
cat <<EOF > /tmp/${FILE}.c

typedef	int	i32;
typedef long long i64;
typedef __int128_t i128;
typedef	unsigned int	u32;
typedef unsigned long long u64;
typedef __uint128_t u128;

void eth_finish(char* _off, i32 _len) __attribute__((import_module("ethereum"), import_name("finish")));
i32  eth_getCallDataSize() __attribute__((import_module("ethereum"),import_name("getCallDataSize")));
void eth_callDataCopy(void *res, i32 dOff, i32 dLen) __attribute__((import_module("ethereum"),import_name("callDataCopy")));

// cpurchase   0xd6960697
// creceived	0x73fac6f0
// refund		0x590e1ae3
static	char	ret[8]={0,0,0,0, 0,0,0,10};
static	char	ret1[8]={0,0,0,20};
void main() // __attribute__((export_name("main")))
{
	if (eth_getCallDataSize() < 4) eth_finish(ret, 8);
	u32 	met;
	eth_callDataCopy(&met, 0, 4);
	switch (__builtin_bswap32(met)) {
	case 0xd6960697:
	case 0x73fac6f0:
	case 0x590e1ae3:
		eth_finish(ret1,4);
	}
}
EOF

clang -c -O3 -Wall -Wno-main-return-type --target=wasm32 /tmp/${FILE}.c
wasm-ld --no-entry --allow-undefined-file=ewasm.syms --export=main --strip-all ${FILE}.o -o ${FILE}.wasm
rm -f ${FILE}.o
#wasm-dis /tmp/${FILE}.wasm | sed -s 's/Main/main/' > /tmp/${FILE}.wat
#wasm-as /tmp/${FILE}.wat -o ${FILE}.wasm
