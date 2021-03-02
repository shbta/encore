// Copyright 2017 The go-ethereum Authors
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

// +build amd64,!cgo arm64,!cgo

package sm2crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/kjx98/go-fastecdsa/sm2"
)

// Ecrecover returns the uncompressed public key that created the given signature.
func Ecrecover(hash, sig []byte) ([]byte, error) {
	if len(sig) != SignatureLength {
		return nil, fmt.Errorf("signature length %d error", len(sig))
	}
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])
	v := uint(sig[64])
	msg := new(big.Int).SetBytes(hash)
	px, py, err := sm2.SM2asm().Recover(r, s, msg, v)
	if err != nil {
		return nil, err
	}
	return sm2.Marshal(S256(), px, py), nil
}

// SigToPub returns the public key that created the given signature.
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	x, y := sm2.Unmarshal(S256(), s)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

// Sign calculates an SM2 signature.
//
// This function is susceptible to chosen plaintext attacks that can leak
// information about the private key that is used for signing. Callers must
// be aware that the given digest cannot be chosen by an adversery. Common
// solution is to hash any input before calculating the signature.
//
// The produced signature is in the [R || S || V] format where V is 0 or 1.
func Sign(digestHash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(digestHash) != DigestLength {
		return nil, fmt.Errorf("hash is required to be exactly %d bytes (%d)", DigestLength, len(digestHash))
	}
	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
	defer zeroBytes(seckey)
	priv := new(big.Int).SetBytes(seckey)
	msg := new(big.Int).SetBytes(digestHash)
	r, s, v, err := sm2.SM2asm().Sign(rand.Reader, msg, priv)
	if err != nil {
		return nil, err
	}
	sig = make([]byte, 65)
	rB := r.Bytes()
	if len(rB) <= 32 {
		copy(sig[32-len(rB):32], rB)
	} else {
		copy(sig[:32], rB)
	}
	rB = s.Bytes()
	if len(rB) <= 32 {
		copy(sig[64-len(rB):], rB)
	} else {
		copy(sig[32:], rB)
	}
	sig[64] = byte(v)
	return sig, nil
}

// VerifySignature checks that the given public key created signature over digest.
// The public key should be in compressed (33 bytes) or uncompressed (65 bytes) format.
// The signature should have the 64 byte [R || S] format.
func VerifySignature(pubkey, digestHash, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}
	px, py := sm2.Unmarshal(S256(), pubkey)
	if px == nil || py == nil {
		return false
	}
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	msg := new(big.Int).SetBytes(digestHash)
	return sm2.SM2asm().Verify(r, s, msg, px, py)
}

// DecompressPubkey parses a public key in the 33-byte compressed format.
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	if len(pubkey) != 33 {
		return nil, fmt.Errorf("invalid public key")
	}
	x, y := sm2.Unmarshal(S256(), pubkey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: x, Y: y, Curve: S256()}, nil
}

// CompressPubkey encodes a public key to the 33-byte compressed format.
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	addr := make([]byte, 33)
	xb := pubkey.X.Bytes()
	if len(xb) < 32 {
		copy(addr[33-len(xb):], xb)
	} else {
		copy(addr[1:], xb)
	}
	if pubkey.Y.Bit(0) == 0 {
		addr[0] = 2
	} else {
		addr[0] = 3
	}
	return addr
}

// S256 returns an instance of the sm2 curve.
func S256() elliptic.Curve {
	return sm2.SM2()
}
