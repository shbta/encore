package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"os"
)

func main() {

	ipcPath := os.Getenv("HOME") + "/testebc/data1/geth.ipc"
	fmt.Println("IPC attach", ipcPath)
	client, err := rpc.Dial(ipcPath)
	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}

	var accounts []string
	err = client.Call(&accounts, "eth_accounts")
	var result string
	for i, acct := range accounts {
		//var result hexutil.Big
		err = client.Call(&result, "eth_getBalance", acct, "latest")

		if err != nil {
			fmt.Println("client.Call err", err)
			return
		}
		if bi, err := hexutil.DecodeBig(result); err != nil {
			fmt.Println("decodeBig err", err)
			continue
		} else {
			bi.Div(bi, big.NewInt(1e14))
			fmt.Printf("account[%d]: %s\nbalance[%d]: %f ETH\n", i, acct, i,
				float64(bi.Int64())/10000)
		}
	}
}
