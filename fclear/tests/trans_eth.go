package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func calcETH(v *big.Int) float64 {
	r := v.Div(v, big.NewInt(1e14))
	return float64(r.Int64()) / 10000.0
}

func main() {
	ipcPath := os.Getenv("HOME") + "/testebc/data1/geth.ipc"
	fmt.Println("IPC attach", ipcPath)
	client, err := ethclient.Dial(ipcPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	var accounts []*common.Address
	accounts, err = client.Accounts(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(accounts) == 0 {
		log.Fatal("no accounts")
	}

	fromAddress := accounts[0]
	//nonce, err := client.PendingNonceAt(context.Background(), *fromAddress)
	bal, err := client.BalanceAt(ctx, *fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Before trans balance:", calcETH(bal))

	value := big.NewInt(8000000000000000000) // in wei (8 eth)
	gasLimit := uint64(21000)                // in units
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x9bf382bea61312c51ad8d31d42a24ac4f704a648")
	var data []byte
	tx := ethereum.CallMsg{
		From:     *fromAddress,
		To:       &toAddress,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Value:    value,
		Data:     data,
	}

	t1 := time.Now()
	var t2 time.Time
	if hash, err := client.SignSendTransaction(context.Background(), &tx); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("tx sent: %s", hash.Hex())
		for {
			_, isPend, _ := client.TransactionByHash(context.Background(), *hash)
			t2 = time.Now()
			if !isPend || t2.Sub(t1) > 5*time.Second {
				if isPend {
					fmt.Println("... timeout")
				} else {
					fmt.Printf("... done, cost %d ms\n", t2.Sub(t1).Milliseconds())
				}
				break
			}
		}
	}
	if bal, err := client.BalanceAt(ctx, *fromAddress, nil); err == nil {
		fmt.Println("After trans balance:", calcETH(bal))
	}

}
