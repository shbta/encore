package main

import (
	"context"
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
	gasLimit := uint64(2100000) // in units
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

func main() {
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

	t1 := time.Now()
	rand.Seed(t1.Unix())
	var t2 time.Time
	var hash *common.Hash
	for i := 0; i < 10; i++ {
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
		if !isPend || t2.Sub(t1) > 5*time.Second {
			fmt.Printf("last tx sent: %s", hash.Hex())
			if isPend {
				fmt.Println("... timeout")
			} else {
				fmt.Printf("... done, cost %d ms\n", t2.Sub(t1).Milliseconds())
			}
			break
		}
	}
	if bal, err := client.BalanceAt(ctx, *fromAddress, nil); err == nil {
		fmt.Println("After dealClear balance:", calcETH(bal))
	}

}
