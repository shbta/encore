#include <ewasm/ewasm.hpp>

static u32 fib(u32 n) {
	if (n < 2) return n;
	return fib(n-1)+fib(n-2);
}

static ewasm_argument	arg1{UINT64};
static ewasm_argument	result1{UINT64};
ewasm_method	_methods[]={
	{(char *)"constructor", 0, 0, 0,},
	{(char *)"fib", 0x73181a7b, 1, 1, &arg1, &result1},
	{(char *)"owner", 0x8da5cb5b, 0, 0,},
};

extern "C" ewasm_ABI __Contract_ABI={3, _methods};

using namespace	ewasm;

static	byte	ret[32]={0,0,0,0, 0,0,0,10};
static	bytes32	key0(1), val32;
extern "C" void ewasm_main(const u32 Id, const ewasm_method *mtdPtr)
{
	u32 n = 10;
	switch (Id) {
	case 0x73181a7b:
		n = arg1._nValue;
		break;
	case 0x8da5cb5b:
	{
		// should be call owner() with Sig 0x8da5cb5b
		debug_printStorageHex(&key0);
		eth_storageLoad(&key0, &val32);
		eth_finish(&val32,32);
		return;
	}
		break;
	case 0:
		// Constructor
		address	sender;
		eth_getCaller(&sender);
		bytes32	val32(sender);
		eth_storageStore(&key0, &val32);
		eth_finish(ret, 0);
		return;
	}
#ifdef	ommit
	u32 res = __builtin_bswap32(fib(n));
	*(u32 *)(ret+4) = 0;
	*(u32 *)(ret+28) = res;
	eth_finish(ret,32);
#else
	result1._nValue = fib(n);
#endif
}
