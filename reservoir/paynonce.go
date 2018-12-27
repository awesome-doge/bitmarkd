// Copyright (c) 2014-2018 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package reservoir

import (
	"encoding/hex"

	"github.com/bitmark-inc/bitmarkd/blockheader"
	"github.com/bitmark-inc/bitmarkd/fault"
)

// type to represent a payment nonce
// Note: no reversal is required for this
type PayNonce [8]byte

// create a random pay nonce
func NewPayNonce() PayNonce {
	_, digest, _, _ := blockheader.Get()
	nonce := PayNonce{}
	copy(nonce[:], digest[:])
	return nonce
}

// get a previous paynonce
func PayNonceFromBlock(number uint64) PayNonce {
	nonce := PayNonce{}
	digest, err := blockheader.DigestForBlock(number)
	if nil != err {
		return nonce
	}
	copy(nonce[:], digest[:])
	return nonce
}

// convert a binary pay nonce to big endian hex string for use by the fmt package (for %s)
func (paynonce PayNonce) String() string {
	return hex.EncodeToString(paynonce[:])
}

// convert a binary pay nonce to big endian hex string for use by the fmt package (for %#v)
func (paynonce PayNonce) GoString() string {
	return "<paynonce:" + hex.EncodeToString(paynonce[:]) + ">"
}

// convert pay nonce to big endian hex text
func (paynonce PayNonce) MarshalText() ([]byte, error) {
	size := hex.EncodedLen(len(paynonce))
	buffer := make([]byte, size)
	hex.Encode(buffer, paynonce[:])
	return buffer, nil
}

// convert little endian hex text into a pay nonce
func (paynonce *PayNonce) UnmarshalText(s []byte) error {
	if len(*paynonce) != hex.DecodedLen(len(s)) {
		return fault.ErrNotAPayNonce
	}
	byteCount, err := hex.Decode(paynonce[:], s)
	if nil != err {
		return err
	}
	if len(paynonce) != byteCount {
		return fault.ErrNotAPayNonce
	}
	return nil
}
