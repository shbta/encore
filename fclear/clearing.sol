pragma solidity >=0.4.17 <0.7.0;

contract FuturesClearing {
struct MemPosition {
	uint32	nLong;
	uint32	nShort;
	int128	p_l;	// fixed128x8
	uint64	fee;	// comissions and regFee ...
}

struct DealRept {
	uint32	client;
	uint32	qty;
	uint64	price;	// fixed64x8
	uint16	symbol;
	uint16	member;
	bool	isOffset;
	bool	isBuy;
}

  address public owner;
  uint32 public multi;	// multi for amount
  string public	name;
  mapping (uint64 => MemPosition)	memPos;

  event Clear(uint16 mem, uint16 ric, bool isOff, bool isBuy);
  modifier isOwner() {
    require(owner == msg.sender, 'only owner call');
    _;
  }

  constructor(uint32 _multi, string memory _name) public {
	owner = msg.sender;
	multi = _multi;
	name = _name;
  }

	function memSymbolIdx(uint16 symb, uint16 memb, uint32 client) internal pure returns (uint64 msIdx) {
		msIdx = uint64(memb) << 48;
		msIdx += uint64(symb) << 32;
		msIdx += uint64(client);
	}

  function getMulti() public view returns (uint32 _multi) {
    return multi;
  }

  function dealClearing(uint32 client,uint32 qty,uint64 price,uint16 symbol,
	uint16 member,bool isOffset,bool isBuy) public {
		uint64 _idx = memSymbolIdx(symbol, member, client);
		emit Clear(member, symbol, isOffset, isBuy);
		if (isOffset) {
			if (isBuy) {
				assert(memPos[_idx].nShort >= qty);
				memPos[_idx].nShort -= qty;
			} else {
				assert(memPos[_idx].nShort >= qty);
				memPos[_idx].nLong -= qty;
			}
		} else {
			if (isBuy) memPos[_idx].nLong += qty; else memPos[_idx].nShort += qty;
		}
		int pl = price * qty * multi;
  }

  function getClientPosition(uint32 client, uint16 symbol, uint16 member) public view returns (uint32 nLong, uint32 nShort) {
	uint64 _idx = memSymbolIdx(symbol, member, client);
	nLong = memPos[_idx].nLong;
	nShort = memPos[_idx].nShort;
  }

}
