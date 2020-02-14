#include "ewasm.h"

static u32 fib(u32 n) {
	if (n < 2) return n;
	return fib(n-1)+fib(n-2);
}

static	byte	ret[32]={0,0,0,0, 0,0,0,10};
#pragma clang diagnostic ignored "-Wmain-return-type"
void main() // __attribute__((export_name("main")))
{
	i32	in_len;
	if ((in_len=eth_getCallDataSize()) < 4) eth_finish(ret, 8);
	u32 	met;
	u32 n = 10;
	if (in_len >= 36) {
		// should be call FibValue(uint32) with Sig 0x73181a7b
		eth_callDataCopy(&met, 32, 4);
		n = __builtin_bswap32(met);
	} else if (in_len == 4) {
		// should be call owner() with Sig 0x8da5cb5b
		ret[7] = 0;
		ret[31] = 0xfe;
		eth_finish(ret,32);
		return;
	}
	u32 res = __builtin_bswap32(fib(n));
	*(u32 *)(ret+4) = 0;
	*(u32 *)(ret+28) = res;
	eth_finish(ret,32);
}
