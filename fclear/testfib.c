#include "ewasm.h"

static u32 fib(u32 n) {
	if (n < 2) return n;
	return fib(n-1)+fib(n-2);
}

static	byte	ret[32]={0,0,0,0, 0,0,0,10};
static	bytes32	key0;
#pragma clang diagnostic ignored "-Wmain-return-type"
void main() // __attribute__((export_name("main")))
{
	i32	in_len;
	if ((in_len=eth_getCallDataSize()) == 0) eth_finish(ret, 8);
	u32 	met;
	u32 n = 10;
	eth_callDataCopy(&met, 0, 4);
	switch (__builtin_bswap32(met)) {
	case 0x73181a7b:
		if (in_len >= 36) {
			// should be call FibValue(uint32) with Sig 0x73181a7b
			eth_callDataCopy(&met, 32, 4);
			n = __builtin_bswap32(met);
		} else {
			eth_revert(ret, 8);
		}
		break;
	case 0x8da5cb5b:
	{
		// should be call owner() with Sig 0x8da5cb5b
		//eth_storageLoad(key0, ret);
		ret[7] = 0;
		ret[31] = 0xfe;
		eth_finish(ret,32);
		return;
	}
		break;
	case 0x861731d5:
		// Constructor
		eth_getCaller(ret+12);
		eth_storageStore(key0, ret);
		eth_finish(ret, 0);
		return;
	}
	u32 res = __builtin_bswap32(fib(n));
	*(u32 *)(ret+4) = 0;
	*(u32 *)(ret+28) = res;
	eth_finish(ret,32);
}
