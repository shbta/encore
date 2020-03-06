package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	wasm "github.com/shbta/go-wasm"
)

var clearDef = `[{"inputs":[{"internalType":"uint32","name":"_multi","type":"uint32"},{"internalType":"string","name":"_name","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint16","name":"mem","type":"uint16"},{"indexed":false,"internalType":"uint16","name":"ric","type":"uint16"},{"indexed":false,"internalType":"bool","name":"isOff","type":"bool"},{"indexed":false,"internalType":"bool","name":"isBuy","type":"bool"}],"name":"Clear","type":"event"},{"inputs":[{"internalType":"uint32","name":"client","type":"uint32"},{"internalType":"uint32","name":"qty","type":"uint32"},{"internalType":"uint64","name":"price","type":"uint64"},{"internalType":"uint16","name":"symbol","type":"uint16"},{"internalType":"uint16","name":"member","type":"uint16"},{"internalType":"bool","name":"isOffset","type":"bool"},{"internalType":"bool","name":"isBuy","type":"bool"}],"name":"dealClearing","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"client","type":"uint32"},{"internalType":"uint16","name":"symbol","type":"uint16"},{"internalType":"uint16","name":"member","type":"uint16"}],"name":"getClientPosition","outputs":[{"internalType":"uint32","name":"nLong","type":"uint32"},{"internalType":"uint32","name":"nShort","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getMulti","outputs":[{"internalType":"uint32","name":"_multi","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"multi","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

func calcETH(v *big.Int) float64 {
	r := v.Div(v, big.NewInt(1e14))
	return float64(r.Int64()) / 10000.0
}

var emptyAddr = common.Address{}

var clearABI abi.ABI
var client *ethclient.Client
var accounts []*common.Address
var contractAddr common.Address
var ctx context.Context

func clearName() (ret string) {
	var clearBytes []byte
	if cBytes, err := clearABI.Pack("name"); err != nil {
		log.Fatal("Pack name", err)
	} else {
		clearBytes = cBytes
	}
	tx := ethereum.CallMsg{
		From: *accounts[0],
		To:   &contractAddr,
		Data: clearBytes,
	}
	var rr interface{}
	if res, err := client.CallContract(ctx, tx, nil); err != nil {
		log.Fatal(err)
	} else if err = clearABI.Unpack(&rr, "name", res); err != nil {
		log.Fatal("Unpack name() ", err)
		//} else {
		//	log.Println("result:", res[:32], res[32:64], res[64:])
	}
	if res, ok := rr.(string); !ok {
		log.Fatal("type of return mismatch")
	} else if !strings.HasPrefix(res, "SHFE Clear") {
		log.Fatal("contract address may be wrong")
	} else {
		ret = res
	}
	return
}

func getClientPosition(clt uint32, sym, member uint16) (nLong, nShort uint32) {
	var clearBytes []byte
	if cBytes, err := clearABI.Pack("getClientPosition", clt, sym, member); err != nil {
		log.Fatal(err)
	} else {
		clearBytes = cBytes
	}
	tx := ethereum.CallMsg{
		From: *accounts[0],
		To:   &contractAddr,
		Data: clearBytes,
	}
	var rr [2]interface{}
	if res, err := client.CallContract(ctx, tx, nil); err != nil {
		log.Fatal(err)
	} else if err = clearABI.Unpack(&rr, "getClientPosition", res); err != nil {
		log.Fatal(err, "res:", res)
	} else {
		nLong = rr[0].(uint32)
		nShort = rr[1].(uint32)
	}
	return
}

func runConstruct() error {
	var clearBytes []byte
	if clearABI.Constructor.Sig() != "()" {
		if cBytes, err := clearABI.Pack("", uint32(5),
			"SHFE Clear"); err != nil {
			return err
		} else {
			clearBytes = cBytes
		}
	}
	//log.Println("Constructor input", len(clearBytes), clearBytes)
	// cost about 41004 gas
	gasLimit := uint64(80000) // in units
	tx := ethereum.CallMsg{
		From:     *accounts[0],
		To:       &contractAddr,
		Gas:      gasLimit,
		GasPrice: big.NewInt(0),
		Data:     clearBytes,
	}
	t1 := time.Now()
	hash, err := client.SignSendTransaction(ctx, &tx)
	if err != nil {
		log.Println("Call constructor", err)
		return err
	}
	// wait for last TX commit
	tx1, _, err := client.TransactionByHash(ctx, *hash)
	if err != nil {
		log.Fatal("Constructor tx failed", err)
	}
	fmt.Printf("Constructor tx sent: %s", hash.Hex())
	if tr, err := bind.WaitMined(ctx, client, tx1); err != nil {
		fmt.Println("...timeout or error", err)
		return err
	} else {
		t2 := time.Now()
		ms := t2.Sub(t1).Milliseconds()
		nSec := float64(ms) / 1000.0
		sCode := "OK"
		if tr.Status == 0 {
			sCode = "Failed"
		}
		fmt.Printf("Constructor %s, cost %.3f seconds, %d gas used\n",
			sCode, nSec, tr.GasUsed)
	}
	return nil
}

func memSymIdx(symb, memb uint16, clt uint32) (ret uint64) {
	ret = uint64(memb) << 48
	ret |= uint64(symb) << 32
	ret |= uint64(clt)
	return
}

type memPosition struct {
	nLong  uint32
	nShort uint32
	p_l    uint64
	fee    uint64
}

var memPos = map[uint64]memPosition{}
var nShortInv, nLongInv int

func goClearing(clt, qty uint32, price uint64, sym, member uint16, isOff,
	isBuy bool) {
	idx := memSymIdx(sym, member, clt)
	mpo := memPos[idx]
	if isOff {
		if isBuy {
			if mpo.nShort < qty {
				nShortInv++
				//log.Println("no enough short position to offset")
				return
			}
			mpo.nShort -= qty
		} else {
			if mpo.nLong < qty {
				nLongInv++
				//log.Println("no enough long position to offset")
				return
			}
			mpo.nLong -= qty
		}
	} else {
		if isBuy {
			mpo.nLong += qty
		} else {
			mpo.nShort += qty
		}
	}
	memPos[idx] = mpo
}

func dealClearing(clt, qty uint32, price uint64, sym, member uint16, isOff,
	isBuy bool) (*common.Hash, error) {
	var clearBytes []byte
	if cBytes, err := clearABI.Pack("dealClearing", clt, qty, price, sym,
		member, isOff, isBuy); err != nil {
		return nil, err
	} else {
		clearBytes = cBytes
	}
	goClearing(clt, qty, price, sym, member, isOff, isBuy)
	// cost about 30469 gas
	gasLimit := uint64(180000) // in units
	//gasPrice, err := client.SuggestGasPrice(ctx)
	//if err != nil {
	//log.Fatal(err)
	//}
	tx := ethereum.CallMsg{
		From:     *accounts[0],
		To:       &contractAddr,
		Gas:      gasLimit,
		GasPrice: big.NewInt(0),
		Data:     clearBytes,
	}
	hash, err := client.SignSendTransaction(ctx, &tx)
	if err != nil {
		log.Fatal(err)
	}
	return hash, err
}

func deploy(code []byte, bValidate bool) (common.Address, error) {
	// cost about 30469 gas
	// cost about 1003264 gas for testfib DEBUG version
	// cost about 863740 gas for testfib Release version
	// cost about 1229104 gas for clear DEBUG version
	// cost about 1096008 gas for clear Release version
	if bValidate {
		var mod wasm.ValModule
		if err := mod.ReadValModule(code); err != nil {
			log.Println("ewasm ReadValModule", err)
			return emptyAddr, err
		}
		if err := mod.Validate(); err != nil {
			log.Println("ewasm Validate", err)
			return emptyAddr, err
		}
		ocLen := len(code)
		if ncLen := len(mod.Bytes()); ncLen < ocLen {
			code = mod.Bytes()
			log.Printf("ewasm contract stripped, old CodeLen %d stripped %d bytes", ocLen, ocLen-ncLen)
		}
	}
	log.Printf("ewasm contract length: %d\n", len(code))

	gasLimit := uint64(1500000) // in units
	tx := ethereum.CallMsg{
		From: *accounts[0],
		Gas:  gasLimit,
		Data: code,
	}
	hash, err := client.SignSendTransaction(ctx, &tx)
	if err != nil {
		log.Fatal(err)
	}
	tx1, _, err := client.TransactionByHash(ctx, *hash)
	if err != nil {
		log.Fatal("last tx failed", err)
	}
	var addr common.Address
	if tr, err := bind.WaitMined(ctx, client, tx1); err != nil {
		return emptyAddr, err
	} else if tr.Status == 0 || tr.ContractAddress == emptyAddr {
		return emptyAddr, fmt.Errorf("zero address")
	} else {
		addr = tr.ContractAddress
		fmt.Printf("Contract addr: %s, gas used: %d\n", addr.Hex(), tr.GasUsed)
		if cc, err := client.CodeAt(ctx, addr, nil); err != nil {
			return emptyAddr, err
		} else if len(cc) == 0 {
			return addr, fmt.Errorf("No code after deploy")
		}
	}
	contractAddr = addr
	if err := runConstruct(); err != nil {
		log.Fatal("run Constructor", err)
	}
	return addr, nil
}

func TimeMs2String(ms uint64) string {
	sec := int64(ms / 1000)
	ns := int64(ms%1000) * 1000000
	tt := time.Unix(sec, ns)
	return tt.Format("2006-01-02 15:04:05.000Z07:00")
}

// old addr w/out ewasm fib "0x6866423b57c92e666274eb8f982FA1438735Ef2B"
// old addr w/ ewasm fib "0x594668030104D245a4Ed6d785E15f66a8200B824"
// addr w/ ewasm fib "0xf5704f03B4e5833AF156B768aCf76Af6491E258D"
func main() {
	var count int
	var dataDir string
	var ctAddr string
	var dumpABI bool
	var bRawContract bool
	var codeDeploy string
	var abiPath string
	flag.IntVar(&count, "count", 1000, "number of contract calls")
	flag.StringVar(&dataDir, "data", "~/testebc", "Data directory for database")
	flag.StringVar(&ctAddr, "contract", "0xf5704f03B4e5833AF156B768aCf76Af6491E258D", "Address of Clearing contract")
	//flag.StringVar(&ctAddr, "contract", "0x6866423b57c92e666274eb8f982FA1438735Ef2B", "Address of Clearing contract")
	flag.BoolVar(&dumpABI, "dump", false, "dump clearABI")
	flag.BoolVar(&bRawContract, "raw", false, "deploy contract as RAW no strpping")
	flag.StringVar(&codeDeploy, "deploy", "", "code to deploy")
	flag.StringVar(&abiPath, "abi", "", "path of ABI file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: clear [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	if abiPath != "" {
		if ff, err := os.Open(abiPath); err != nil {
			fmt.Println("abi File", err)
		} else {
			if code, err := ioutil.ReadAll(ff); err != nil {
				ff.Close()
				log.Fatal(err)
			} else {
				clearDef = string(code)
			}
			ff.Close()
		}
	}
	if cABI, err := abi.JSON(strings.NewReader(clearDef)); err != nil {
		log.Fatal("abi.JSON", err)
	} else {
		clearABI = cABI
	}
	if dumpABI {
		{
			ab := clearABI.Constructor
			fmt.Printf("Constructor Method %s Id: %s, Sig: %s\n",
				ab.Name, common.ToHex(ab.ID()), ab.Sig())
		}
		for _, ab := range clearABI.Methods {
			fmt.Printf("Method %s Id: %s, Sig: %s\n", ab.Name,
				common.ToHex(ab.ID()), ab.Sig())
		}
		os.Exit(0)
	}
	var ipcPath string
	if len(dataDir) > 0 && dataDir[0] == '~' {
		ipcPath = os.Getenv("HOME") + dataDir[1:] + "/geth.ipc"
	} else {
		ipcPath = dataDir + "/geth.ipc"
	}
	fmt.Println("IPC attach", ipcPath)
	if clt, err := ethclient.Dial(ipcPath); err != nil {
		log.Fatal(err)
	} else {
		client = clt
	}

	contractAddr = common.HexToAddress(ctAddr)
	ctx = context.Background()
	if acct, err := client.Accounts(ctx); err != nil {
		log.Fatal(err)
	} else if len(acct) == 0 {
		log.Fatal("no accounts")
	} else {
		accounts = acct
	}
	if codeDeploy != "" {
		if ff, err := os.Open(codeDeploy); err != nil {
			log.Fatal("open", codeDeploy, err)
		} else if code, err := ioutil.ReadAll(ff); err != nil {
			ff.Close()
			log.Fatal(err)
		} else {
			ff.Close()
			if addr, err := deploy(code, !bRawContract); err != nil {
				if addr != emptyAddr {
					log.Fatal(addr.Hex(), " error:", err)
				} else {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("Contract deployed at %s\n", addr.Hex())
			}
		}
		os.Exit(0)
	}

	fromAddress := accounts[0]
	//nonce, err := client.PendingNonceAt(context.Background(), *fromAddress)
	if bal, err := client.BalanceAt(ctx, *fromAddress, nil); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Before dealClear balance:", calcETH(bal))
	}
	fmt.Println("Clear contract name:", clearName())

	var blockS uint64
	if hh, err := client.HeaderByNumber(ctx, nil); err == nil && hh.Number != nil {
		blockS = hh.Number.Uint64()
		fmt.Printf("block %d %s before call\n", blockS, TimeMs2String(hh.TimeMilli))
	}
	var clt uint32
	memb := uint16(101)
	// fetch 1000 .. 1003 memPosition
	for clt = 1000; clt < 1004; clt++ {
		idx := memSymIdx(1, memb, clt)
		mpo := memPos[idx]
		mpo.nLong, mpo.nShort = getClientPosition(clt, 1, memb)
		memPos[idx] = mpo
		if clt == 1003 {
			fmt.Printf("Before Clear client %d position: %d long %d short\n",
				clt, mpo.nLong, mpo.nShort)
		}
	}
	t1 := time.Now()
	rand.Seed(t1.Unix())
	var nSec float64
	var t2 time.Time
	var hash *common.Hash
	for i := 0; i < count; i++ {
		clt = uint32(1000 + i%4)
		qty := uint32(1 + rand.Int31n(200))
		prc := uint64(53000 + 100*rand.Int31n(30))
		bOff := (i & 1) == 0
		bBuy := (qty & 1) != 0
		if h, err := dealClearing(clt, qty, prc, 1, memb, bOff, bBuy); err != nil {
			log.Fatal("dealClearing failed", err)
		} else {
			hash = h
		}
	}
	// wait for last TX commit
	tx, _, err := client.TransactionByHash(ctx, *hash)
	if err != nil {
		log.Fatal("last tx failed", err)
	}
	fmt.Printf("last tx sent: %s", hash.Hex())
	if tr, err := bind.WaitMined(ctx, client, tx); err != nil {
		fmt.Println("...timeout or error", err)
	} else {
		t2 = time.Now()
		ms := t2.Sub(t1).Milliseconds()
		nSec = float64(ms) / 1000.0
		sCode := "OK"
		if tr.Status == 0 {
			sCode = "Failed"
		}
		fmt.Printf("... %s @Block %d, cost %d ms\n", sCode, tr.BlockNumber.Uint64(), ms)
		fmt.Printf("%d Gas used per call\n", tr.GasUsed)
		fmt.Printf("%.3f calls per second\n", float64(count)/nSec)
	}
	if hh, err := client.HeaderByNumber(ctx, nil); err == nil && hh.Number != nil {
		blockE := hh.Number.Uint64()
		fmt.Printf("block %d %s after call\n", blockE, TimeMs2String(hh.TimeMilli))
		fmt.Printf("mined %.2f blocks per second\n", float64(blockE-blockS)/nSec)
		nBlk := int(blockE - blockS)
		fmt.Printf("%d contract calls per block\n", 1000/nBlk)
	}
	if bal, err := client.BalanceAt(ctx, *fromAddress, nil); err == nil {
		fmt.Println("After dealClear balance:", calcETH(bal))
	}
	nL, nS := getClientPosition(clt, 1, memb)
	fmt.Printf("client %d position: %d long, %d short\n", clt, nL, nS)
	fmt.Printf("short not enough times: %d, long not enough times: %d\n",
		nShortInv, nLongInv)
	idx := memSymIdx(1, memb, clt)
	mpo := memPos[idx]
	if mpo.nLong != nL || mpo.nShort != nS {
		fmt.Printf("goClear diff, client %d position: %d long %d short\n", clt,
			mpo.nLong, mpo.nShort)
	}
}
