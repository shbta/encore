#include <ewasm/ewasm.hpp>

static u64 fib(u32 n) {
	if (n < 2) return n;
	u64	result=2;
	u64	pre = 1;
	u64 next = 1;
	for (u32 i = 2; i < n; ++i) {
		result = pre + next;
		pre = next;
		next = result;
	}
	return result;
}

static ewasm_argument	arg1[]={{UINT64}};
static ewasm_argument	result1[]={{UINT64}};
static ewasm_argument	retAddr[]={{UINT160}};
static ewasm::method	_methods[]={
	{"constructor", 0},
	{"fib", 0x73181a7b, arg1, result1},
	{"owner", 0x8da5cb5b, 0, retAddr},
};

namespace ewasm {
static ABI myABI={_methods};
}

using namespace	ewasm;

extern "C" {
ewasm_ABI __Contract_ABI{myABI.nMethods, myABI.methods};
// assign C struct from C++ struct does not work
//ewasm_ABI __Contract_ABI=myABI;
}

static	byte	ret[32]={0,0,0,0, 0,0,0,10};
static	bytes32	key0(1), val32;
void ewasm_main(const u32 Id, const ewasm_method *mtdPtr)
{
	static_assert(sizeof(method) == sizeof(ewasm_method), "size of ewasm_method and method MUST equal");
	static_assert(sizeof(ABI) == sizeof(ewasm_ABI), "size of ewasm_ABI and ABI MUST equal");
	u32 n = 10;
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
	u64 res = __builtin_bswap64(fib(n));
	*(u64 *)(ret+4) = 0;
	*(u64 *)(ret+24) = res;
	eth_finish(ret,32);
#else
	result1[0]._nValue = fib(n);
#endif
}
