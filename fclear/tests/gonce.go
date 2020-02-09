package main

import (

	//基础库
	"context"
	"fmt"
	"log"
	"os"
	//第三方库
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getethnonce(fromaddress string) uint64 {
	ipcPath := os.Getenv("HOME") + "/testebc/data1/geth.ipc"
	fmt.Println("IPC attach", ipcPath)
	client, err := ethclient.Dial(ipcPath)
	if err != nil {
		log.Print("xxxxxxx:", err)
	}
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(fromaddress),
	}
	nonce, err := client.PendingNonceAt(context.Background(), fromAccDef.Address)
	if err != nil {
		log.Print("xxxxxxx:", err)
	}
	return nonce
}

func main() {
	faddr := "0x8a99c8b23686d1a079d3db4702a6ed40bc6b156f"
	once := getethnonce(faddr)
	fmt.Println("Once of", faddr, "is", once)
}
