#include <ewasm/ewasm.hpp>

static u32 fib(u32 n) {
	if (n < 2) return n;
	return fib(n-1)+fib(n-2);
}

static ewasm_argument	arg1[]={{UINT64}};
static ewasm_argument	result1[]={{UINT64}};
static ewasm_argument	retAddr[]={{UINT160}};
ewasm::method	_methods[]={
	{"constructor", 0},
	{"fib", 0x73181a7b, arg1, result1},
	{"owner", 0x8da5cb5b, 0, retAddr},
};

namespace ewasm {
static ABI myABI={_methods};
}

using namespace	ewasm;

extern "C" {
ewasm_ABI __Contract_ABI=myABI;
}

static	byte	ret[32]={0,0,0,0, 0,0,0,10};
static	bytes32	key0(1), val32;
void ewasm_main(const u32 Id, const ewasm_method *mtdPtr)
{
	//static_assert(sizeof(nullArg) == 0, "size of empty arguments MUST be 0");
	u32 n = 10;
	debug_printMemHex((void *)&Id, sizeof(Id));
	switch (Id) {
	case 0x73181a7b:
		n = arg1[0]._nValue;
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
	result1[0]._nValue = fib(n);
#endif
}
