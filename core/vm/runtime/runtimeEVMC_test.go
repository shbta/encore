// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
// +build evmc

package runtime

import (
	//"io/ioutil"
	"math/big"
	//"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

func TestDefaults(t *testing.T) {
	cfg := new(Config)
	setDefaults(cfg)

	if cfg.Difficulty == nil {
		t.Error("expected difficulty to be non nil")
	}

	if cfg.Time == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.GasLimit == 0 {
		t.Error("didn't expect gaslimit to be zero")
	}
	if cfg.GasPrice == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.Value == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.GetHashFn == nil {
		t.Error("expected time to be non nil")
	}
	if cfg.BlockNumber == nil {
		t.Error("expected block number to be non nil")
	}
}

func TestEVM(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("crashed with: %v", r)
		}
	}()

	code := []byte{
		byte(vm.DIFFICULTY),
		byte(vm.TIMESTAMP),
		byte(vm.GASLIMIT),
		byte(vm.PUSH1),
		byte(vm.ORIGIN),
		byte(vm.BLOCKHASH),
		byte(vm.COINBASE),
	}
	//println("test code:", hex.Dump(code))
	Execute(code, nil, nil)
}

var evmCode = common.Hex2Bytes("0061736d01000000010e0360027f7f0060000060027f7f0002130108657468657265756d0666696e6973680002030201010405017001010105030100020608017f01419088040b071102066d656d6f72790200046d61696e00010901000a0b010900418008410810000b0b0f01004180080b08000000000000000a")
var wasmCode = common.Hex2Bytes("0061736d010000000113046000017f60027f7f0060037f7f7f0060000002460308657468657265756d0f67657443616c6c4461746153697a65000008657468657265756d0666696e697368000108657468657265756d0c63616c6c44617461436f70790002030201030405017001010105030100020608017f0141a088040b071102066d656d6f72790200046d61696e00030a8a0301870302037f037e23808080800041106b220024808080800002400240024002400240108080808000220141034a0d004180888080004108108180808000410a21020c010b0240200141234a0d00410a210220014104470d01410041fe013a009388808000410041003a00878880800041808880800041201081808080000c040b2000410c6a41204104108280808000200028020c22014118742001410874418080fc07717220014108764180fe03712001411876727222024102490d014202210320014180808010460d020b2002417e6a2102420121034201210403402004200322057c2103200521042002417f6a22020d000c020b0b2002ad21030b410042003703848880800041002003423886200342288642808080808080c0ff0083842003421886428080808080e03f8320034208864280808080f01f838484200342088842808080f80f832003421888428080fc07838420034228884280fe038320034238888484843703988880800041808880800041201081808080000b200041106a2480808080000b0b2701004180080b20000000000000000a000000000000000000000000000000000000000000000000")

func TestExecute(t *testing.T) {
	code := evmCode
	ret, _, err := Execute(code, nil, nil)
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	num := new(big.Int).SetBytes(ret)
	if num.Cmp(big.NewInt(10)) != 0 {
		t.Error("Expected 10, got", num)
	}
}

func TestExecuteEWASM(t *testing.T) {
	/*
		var code []byte
		if ff, err := os.Open("test.wasm"); err != nil {
			t.Fatal("didn't expect error", err)
		} else {
			code, _ = ioutil.ReadAll(ff)
			ff.Close()
		}
		println("code:", common.ToHex(code))
	*/
	code := wasmCode
	println("test execute EWASM file Len:", len(code))

	ret, _, err := Execute(code, nil, nil)
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	num := new(big.Int).SetBytes(ret)
	if num.Cmp(big.NewInt(10)) != 0 {
		t.Error("Expected 10, got", num)
	}
}

var definition = `[{"constant":true,"inputs":[{"name":"n","type":"uint32"}],"name":"FibValue","outputs":[{"name":"res","type":"uint64"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor","signature":"constructor"}]`

func TestCall(t *testing.T) {
	state, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	address := common.HexToAddress("0x0a")
	code := wasmCode
	state.SetCode(address, code)
	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		t.Fatal(err)
	}
	cpurchase, err := abi.Pack("FibValue", uint32(50))
	if err != nil {
		t.Fatal(err)
	}

	ret, _, err := Call(address, cpurchase, &Config{State: state})
	if err != nil {
		t.Fatal("didn't expect error", err)
	}

	var num uint64
	if err := abi.Unpack(&num, "FibValue", ret); err != nil {
		t.Fatal("abi Unpack", err)
	}
	if num != 12586269025 {
		t.Error("Expected 12586269025, got", num)
	}
}

func TestCreate(t *testing.T) {
	var (
		statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		sender     = common.BytesToAddress([]byte("sender"))
	)

	code := wasmCode
	statedb.CreateAccount(sender)
	runtimeConfig := Config{
		Origin:      sender,
		State:       statedb,
		GasLimit:    10000000,
		Difficulty:  big.NewInt(0x200000),
		Time:        new(big.Int).SetUint64(0),
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).SetUint64(1),
		ChainConfig: &params.ChainConfig{
			ChainID:        big.NewInt(1),
			HomesteadBlock: new(big.Int),
			ByzantiumBlock: new(big.Int),
		},
		EVMConfig: vm.Config{},
	}
	code1, address, leftOverGas, err := Create(code, &runtimeConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ret code len:", len(code1))
	t.Log("ret addr:", common.ToHex(address.Bytes()))
	t.Log("gas used:", 10000000-leftOverGas)
}

func BenchmarkCall(b *testing.B) {
	code := wasmCode
	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		b.Fatal(err)
	}

	cpurchase, err := abi.Pack("FibValue", uint32(15))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 400; j++ {
			Execute(code, cpurchase, nil)
		}
	}
}

func benchmarkEVM_Create(bench *testing.B, code string) {
	var (
		statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		sender     = common.BytesToAddress([]byte("sender"))
		receiver   = common.BytesToAddress([]byte("receiver"))
	)

	statedb.CreateAccount(sender)
	statedb.SetCode(receiver, common.FromHex(code))
	runtimeConfig := Config{
		Origin:      sender,
		State:       statedb,
		GasLimit:    10000000,
		Difficulty:  big.NewInt(0x200000),
		Time:        new(big.Int).SetUint64(0),
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).SetUint64(1),
		ChainConfig: &params.ChainConfig{
			ChainID:             big.NewInt(1),
			HomesteadBlock:      new(big.Int),
			ByzantiumBlock:      new(big.Int),
			ConstantinopleBlock: new(big.Int),
			DAOForkBlock:        new(big.Int),
			DAOForkSupport:      false,
			EIP150Block:         new(big.Int),
			EIP155Block:         new(big.Int),
			EIP158Block:         new(big.Int),
		},
		EVMConfig: vm.Config{},
	}
	// Warm up the intpools and stuff
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		Call(receiver, []byte{}, &runtimeConfig)
	}
	bench.StopTimer()
}

func benchmarkEVM_CreateNew(bench *testing.B, code []byte) {
	var (
		statedb, _ = state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
		sender     = common.BytesToAddress([]byte("sender"))
	)

	statedb.CreateAccount(sender)
	runtimeConfig := Config{
		Origin:      sender,
		State:       statedb,
		GasLimit:    10000000,
		Difficulty:  big.NewInt(0x200000),
		Time:        new(big.Int).SetUint64(0),
		Coinbase:    common.Address{},
		BlockNumber: new(big.Int).SetUint64(1),
		ChainConfig: &params.ChainConfig{
			ChainID:        big.NewInt(1),
			HomesteadBlock: new(big.Int),
			ByzantiumBlock: new(big.Int),
		},
		EVMConfig: vm.Config{},
	}
	// Warm up the intpools and stuff
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		_, _, _, err := Create(code, &runtimeConfig)
		if err != nil {
			bench.Fatal(err)
		}
	}
	bench.StopTimer()
}

func BenchmarkEVM_CREATE_500(bench *testing.B) {
	// initcode size 500K, repeatedly calls CREATE and then modifies the mem contents
	//benchmarkEVM_Create(bench, "5b6207a120600080f0600152600056")
	benchmarkEVM_CreateNew(bench, wasmCode)
}
func BenchmarkEVM_CREATE2_500(bench *testing.B) {
	// initcode size 500K, repeatedly calls CREATE2 and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b586207a120600080f5600152600056")
}
func BenchmarkEVM_CREATE_1200(bench *testing.B) {
	// initcode size 1200K, repeatedly calls CREATE and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b62124f80600080f0600152600056")
}
func BenchmarkEVM_CREATE2_1200(bench *testing.B) {
	// initcode size 1200K, repeatedly calls CREATE2 and then modifies the mem contents
	benchmarkEVM_Create(bench, "5b5862124f80600080f5600152600056")
}
