// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64,!arm64,notwork

package btc

var (
	//pSM2 p256Curve
	pBTC Curve
)

func initP256Arch() {
	// Use pure Go implementation.
	// p256.go not work yet, using default golang impl
	pBTC = BTCgo()
}
