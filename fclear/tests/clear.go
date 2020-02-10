package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var clearDef = `[{"inputs":[{"internalType":"uint32","name":"_multi","type":"uint32"},{"internalType":"string","name":"_name","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint16","name":"mem","type":"uint16"},{"indexed":false,"internalType":"uint16","name":"ric","type":"uint16"},{"indexed":false,"internalType":"bool","name":"isOff","type":"bool"},{"indexed":false,"internalType":"bool","name":"isBuy","type":"bool"}],"name":"Clear","type":"event"},{"inputs":[{"internalType":"uint32","name":"client","type":"uint32"},{"internalType":"uint32","name":"qty","type":"uint32"},{"internalType":"uint64","name":"price","type":"uint64"},{"internalType":"uint16","name":"symbol","type":"uint16"},{"internalType":"uint16","name":"member","type":"uint16"},{"internalType":"bool","name":"isOffset","type":"bool"},{"internalType":"bool","name":"isBuy","type":"bool"}],"name":"dealClearing","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"client","type":"uint32"},{"internalType":"uint16","name":"symbol","type":"uint16"},{"internalType":"uint16","name":"member","type":"uint16"}],"name":"getClientPosition","outputs":[{"internalType":"uint32","name":"nLong","type":"uint32"},{"internalType":"uint32","name":"nShort","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getMulti","outputs":[{"internalType":"uint32","name":"_multi","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"multi","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

func calcETH(v *big.Int) float64 {
	r := v.Div(v, big.NewInt(1e14))
	return float64(r.Int64()) / 10000.0
}

var clearABI abi.ABI
var client *ethclient.Client
var accounts []*common.Address
var contractAddr common.Address
var ctx context.Context

func dealClearing(clt, qty uint32, price uint64, sym, member uint16, isOff, isBuy bool) (*common.Hash, error) {
	var clearBytes []byte
	if cBytes, err := clearABI.Pack("dealClearing", clt, qty, price, sym, member, isOff, isBuy); err != nil {
		return nil, err
	} else {
		clearBytes = cBytes
	}
	// cost about 22680 gas
	gasLimit := uint64(32000) // in units
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	tx := ethereum.CallMsg{
		From:     *accounts[0],
		To:       &contractAddr,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     clearBytes,
	}
	hash, err := client.SignSendTransaction(ctx, &tx)
	if err != nil {
		log.Fatal(err)
	}
	return hash, err
}

func TimeMs2String(ms uint64) string {
	sec := int64(ms / 1000)
	ns := int64(ms%1000) * 1000000
	tt := time.Unix(sec, ns)
	return tt.Format("2006-01-02 15:04:05.000Z07:00")
}

func main() {
	var count int
	flag.IntVar(&count, "count", 1000, "number of contract calls")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: clear [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	ipcPath := os.Getenv("HOME") + "/testebc/data1/geth.ipc"
	fmt.Println("IPC attach", ipcPath)
	if clt, err := ethclient.Dial(ipcPath); err != nil {
		log.Fatal(err)
	} else {
		client = clt
	}
	if cABI, err := abi.JSON(strings.NewReader(clearDef)); err != nil {
		log.Fatal("abi.JSON", err)
	} else {
		clearABI = cABI
	}

	ctx = context.Background()
	if acct, err := client.Accounts(ctx); err != nil {
		log.Fatal(err)
	} else if len(acct) == 0 {
		log.Fatal("no accounts")
	} else {
		accounts = acct
	}

	fromAddress := accounts[0]
	//nonce, err := client.PendingNonceAt(context.Background(), *fromAddress)
	if bal, err := client.BalanceAt(ctx, *fromAddress, nil); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Before dealClear balance:", calcETH(bal))
	}

	//toAddress := common.HexToAddress("0x9bf382bea61312c51ad8d31d42a24ac4f704a648")

	var blockS uint64
	if hh, err := client.HeaderByNumber(ctx, nil); err == nil && hh.Number != nil {
		blockS = hh.Number.Uint64()
		fmt.Printf("block %d %s before call\n", blockS, TimeMs2String(hh.TimeMilli))
	}
	t1 := time.Now()
	rand.Seed(t1.Unix())
	var nSec float64
	var t2 time.Time
	var hash *common.Hash
	for i := 0; i < count; i++ {
		clt := uint32(1000 + i%4)
		qty := uint32(1 + rand.Int31n(200))
		prc := uint64(53000 + 100*rand.Int31n(30))
		memb := uint16(101)
		bOff := false
		bBuy := (qty & 1) != 0
		if h, err := dealClearing(clt, qty, prc, 1, memb, bOff, bBuy); err != nil {
			log.Fatal("dealClearing failed", err)
		} else {
			hash = h
		}
	}
	// wait for last TX commit
	for {
		_, isPend, _ := client.TransactionByHash(context.Background(), *hash)
		t2 = time.Now()
		if !isPend || t2.Sub(t1) > 30*time.Second {
			fmt.Printf("last tx sent: %s", hash.Hex())
			if isPend {
				fmt.Println("... timeout")
			} else {
				ms := t2.Sub(t1).Milliseconds()
				nSec = float64(ms) / 1000.0
				fmt.Printf("... done, cost %d ms\n", ms)
				if tr, err := client.TransactionReceipt(ctx, *hash); err == nil {
					fmt.Printf("%d Gas used per call\n", tr.GasUsed)
				}
				fmt.Printf("%.3f calls per second\n", float64(count)/nSec)
			}
			break
		}
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

}
