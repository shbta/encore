#include <ewasm/ewasm.hpp>
#include <assert.h>

using namespace	ewasm;
struct	memPosition {
	uint32_t	nLong;
	uint32_t	nShort;
	uint64_t	fee;
	uint128_t	p_l;
};

static ewasm_argument	arg2[]={{UINT32}, {STRING}};
static ewasm_argument	arg3[]={{UINT32}, {UINT16}, {UINT16}};
static ewasm_argument	arg7[]={{UINT32}, {UINT32}, {UINT64},
				{UINT16}, {UINT16}, {BOOL}, {BOOL},
			};
static ewasm_argument	ret2[]={{UINT32},{UINT32}};
static ewasm_argument	retStr[]={{STRING}};
static ewasm::method	_methods[]={
	{"constructor", 0, arg2},
	{"owner", 0x8da5cb5b},
	{"name", 0x06fdde03, 0, retStr},
	{"dealClearing", 0xbe704381, arg7},
	{"getClientPosition", 0xf42a90d6, arg3, ret2},
};

namespace ewasm {
static ABI myABI={_methods};
}

extern "C" {
ewasm_ABI __Contract_ABI{myABI.nMethods, myABI.methods};
// assign C struct from C++ struct does not work
//ewasm_ABI __Contract_ABI=myABI;
}

static	uint64_t memSymbolIdx(uint16_t symb, uint16_t memb, uint32_t clt) {
	uint64_t	ret = (uint64_t)memb << 48;
	ret |= (uint64_t)symb << 32;
	ret |= (uint64_t)clt;
	return ret;
}

static void	doClear(u32 clt, u32 qty, u64 price, uint16_t sym, uint16_t memb,
			bool isOff, bool isBuy)
{
	auto idx = memSymbolIdx(sym, memb, clt);
	bytes32		key(idx), val32;
	eth_storageLoad(&key, &val32);
	memPosition *mp = (memPosition *)((void *)&val32);
	if (isOff) {
		if (isBuy) {
			assert(mp->nShort >= qty);
			mp->nShort -= qty;
		} else {
			assert(mp->nShort >= qty);
			mp->nLong -= qty;
		}
	} else {
		if (isBuy) mp->nLong += qty; else mp->nShort += qty;
	}
	//int128_t pl = price * qty * multi;
	eth_storageStore(&key, &val32);
}

static	bytes	cName("SHFE Clear");

void ewasm_main(const u32 Id, const ewasm_method *mtdPtr)
{
	static_assert(sizeof(memPosition) == 32, "memPosition size MUST == 32");
	bytes32		val32;
	switch (Id) {
	case 0x8da5cb5b:
	{
		// should be call owner() with Sig 0x8da5cb5b
		bytes32		key0(1);
		debug_printStorageHex(&key0);
		eth_storageLoad(&key0, &val32);
		eth_finish(&val32,32);
	}
		break;
	case 0x06fdde03:
		// name()
		debug_print(cName.data(), cName.size());
		retStr[0].pValue = cName;
		break;
	case 0:
		// Constructor
	{
		address	sender;
		bytes32		key0(1);
		eth_getCaller(&sender);
		bytes32	val32(sender);
		eth_storageStore(&key0, &val32);
		debug_print(arg2[1].pValue._data, arg2[1].pValue._size);
		eth_finish(nullptr, 0);
	}
		return;
	case 0xbe704381:
		// dealClearing, arg7
		doClear(arg7[0]._nValue, arg7[1]._nValue, arg7[2]._nValue,
				arg7[3]._nValue, arg7[4]._nValue, arg7[5]._nValue,
				arg7[6]._nValue);
		break;
	case 0xf42a90d6:
		// getClientPosition, arg3
	{
		auto idx = memSymbolIdx(arg3[1]._nValue, arg3[2]._nValue, arg3[0]._nValue);
		bytes32		key(idx);
		eth_storageLoad(&key, &val32);
		memPosition *mp = (memPosition *)((void *)&val32);
		ret2[0]._nValue = mp->nLong;
		ret2[1]._nValue = mp->nShort;
	}
		break;
	}
}
