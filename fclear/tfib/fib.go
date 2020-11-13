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

var fibDef = `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[{"internalType":"uint32","name":"n","type":"uint32"}],"name":"FibValue","outputs":[{"internalType":"uint64","name":"res","type":"uint64"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

func calcETH(v *big.Int) float64 {
	r := v.Div(v, big.NewInt(1e14))
	return float64(r.Int64()) / 10000.0
}

var emptyAddr = common.Address{}

var fibABI abi.ABI
var client *ethclient.Client
var accounts []*common.Address
var contractAddr common.Address
var ctx context.Context

func FibValue(fn uint) (ret uint64) {
	var fibBytes []byte
	if cBytes, err := fibABI.Pack("FibValue", uint32(fn)); err != nil {
		log.Fatal("Pack FibValue", err)
	} else {
		fibBytes = cBytes
	}
	tx := ethereum.CallMsg{
		From: *accounts[0],
		To:   &contractAddr,
		Data: fibBytes,
	}
	var rr interface{}
	if res, err := client.CallContract(ctx, tx, nil); err != nil {
		log.Fatal(err)
	} else if err = fibABI.UnpackIntoInterface(&rr, "FibValue", res); err != nil {
		log.Fatal("Unpack FibValue() ", err)
	}
	if res, ok := rr.(uint64); !ok {
		log.Fatal("type of return mismatch")
	} else {
		ret = res
	}
	return
}

func runConstruct() error {
	var fibBytes []byte
	if fibABI.Constructor.Sig != "()" {
		if cBytes, err := fibABI.Pack("", uint32(5),
			"SHFE Clear"); err != nil {
			return err
		} else {
			fibBytes = cBytes
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
		Data:     fibBytes,
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

func deploy(code []byte, bValidate bool) (common.Address, error) {
	// cost about 30469 gas
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

// addr of ewasm fib "0x1fbAeA3a0Cd2f9B5B01DfDBB85C68cf29afb2F55"
// addr of evm fib "0x594668030104D245a4Ed6d785E15f66a8200B824"
func main() {
	var count int
	var dataDir string
	var ctAddr string
	var dumpABI bool
	var codeDeploy string
	var abiPath string
	var bRawContract bool
	var fibN uint
	flag.IntVar(&count, "count", 1000, "number of contract calls")
	flag.UintVar(&fibN, "fib", 50, "number of fibonacci seq")
	flag.StringVar(&dataDir, "data", "~/testebc", "Data directory for database")
	flag.StringVar(&ctAddr, "contract", "0xf5704f03B4e5833AF156B768aCf76Af6491E258D", "Address of Clearing contract")
	//flag.StringVar(&ctAddr, "contract", "0x6866423b57c92e666274eb8f982FA1438735Ef2B", "Address of Clearing contract")
	flag.BoolVar(&dumpABI, "dump", false, "dump clearABI")
	flag.StringVar(&codeDeploy, "deploy", "", "code to deploy")
	flag.StringVar(&abiPath, "abi", "", "path of ABI file")
	flag.BoolVar(&bRawContract, "raw", false, "deploy contract as RAW no strpping")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: tfib [options]\n")
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
				fibDef = string(code)
			}
			ff.Close()
		}
	}
	if cABI, err := abi.JSON(strings.NewReader(fibDef)); err != nil {
		log.Fatal("abi.JSON", err)
	} else {
		fibABI = cABI
	}
	if dumpABI {
		{
			ab := fibABI.Constructor
			fmt.Printf("Constructor Method %s Id: %s, Sig: %s\n",
				ab.Name, common.Bytes2Hex(ab.ID), ab.Sig)
		}
		for _, ab := range fibABI.Methods {
			fmt.Printf("Method %s Id: %s, Sig: %s\n", ab.Name,
				common.Bytes2Hex(ab.ID), ab.Sig)
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
		fmt.Println("Before Fibonacci balance:", calcETH(bal))
	}

	var blockS uint64
	if hh, err := client.HeaderByNumber(ctx, nil); err == nil && hh.Number != nil {
		blockS = hh.Number.Uint64()
		fmt.Printf("block %d %s before call\n", blockS, TimeMs2String(hh.TimeMilli))
	}
	t1 := time.Now()
	rand.Seed(t1.Unix())
	var nSec float64
	var t2 time.Time
	var fibRes uint64
	for i := 0; i < count; i++ {
		fibRes = FibValue(fibN)
	}
	{
		t2 = time.Now()
		ms := t2.Sub(t1).Milliseconds()
		nSec = float64(ms) / 1000.0
		fmt.Printf("Fib(%d) is %d ... cost %d ms\n", fibN, fibRes, ms)
		fmt.Printf("%.3f calls per second\n", float64(count)/nSec)
	}
	if hh, err := client.HeaderByNumber(ctx, nil); err == nil && hh.Number != nil {
		blockE := hh.Number.Uint64()
		fmt.Printf("block %d %s after call\n", blockE, TimeMs2String(hh.TimeMilli))
		fmt.Printf("mined %.2f blocks per second\n", float64(blockE-blockS)/nSec)
		nBlk := int(blockE - blockS)
		if nBlk > 0 {
			fmt.Printf("%d contract calls per block\n", 1000/nBlk)
		}
	}
	if bal, err := client.BalanceAt(ctx, *fromAddress, nil); err == nil {
		fmt.Println("After Fibonacci balance:", calcETH(bal))
	}
}
