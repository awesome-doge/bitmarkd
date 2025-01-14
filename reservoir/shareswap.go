// SPDX-License-Identifier: ISC
// Copyright (c) 2014-2020 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package reservoir

import (
	"time"

	"github.com/bitmark-inc/bitmarkd/blockheader"
	"github.com/bitmark-inc/bitmarkd/constants"
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/merkle"
	"github.com/bitmark-inc/bitmarkd/ownership"
	"github.com/bitmark-inc/bitmarkd/pay"
	"github.com/bitmark-inc/bitmarkd/storage"
	"github.com/bitmark-inc/bitmarkd/transactionrecord"
	"github.com/bitmark-inc/logger"
)

// SwapInfo - result returned by store share
type SwapInfo struct {
	RemainingOne uint64
	RemainingTwo uint64
	Id           pay.PayId
	TxId         merkle.Digest
	Packed       []byte
	Payments     []transactionrecord.PaymentAlternative
}

// returned data from verifySwap
type verifiedSwapInfo struct {
	balanceOne          uint64
	balanceTwo          uint64
	txId                merkle.Digest
	packed              []byte
	issueTxId           merkle.Digest
	transferBlockNumber uint64
	issueBlockNumber    uint64
}

// storeSwap - verify and store a swap request
func storeSwap(swap *transactionrecord.ShareSwap, shareQuantityHandle storage.Handle, shareHandle storage.Handle, ownerDataHandle storage.Handle, blockOwnerPaymentHandle storage.Handle) (*SwapInfo, bool, error) {
	if nil == shareQuantityHandle || nil == shareHandle || nil == ownerDataHandle || nil == blockOwnerPaymentHandle {
		return nil, false, fault.NilPointer
	}

	globalData.Lock()
	defer globalData.Unlock()

	verifyResult, duplicate, err := verifySwap(swap, shareQuantityHandle, shareHandle, ownerDataHandle)
	if err != nil {
		return nil, false, err
	}

	// compute pay id
	packedSwap := verifyResult.packed
	payId := pay.NewPayId([][]byte{packedSwap})

	txId := verifyResult.txId

	payments := getPayments(verifyResult.transferBlockNumber, verifyResult.issueBlockNumber, nil, blockOwnerPaymentHandle)

	spendKeyOne := makeSpendKey(swap.OwnerOne, swap.ShareIdOne)
	spendKeyTwo := makeSpendKey(swap.OwnerTwo, swap.ShareIdTwo)

	spendOne := globalData.spend[spendKeyOne]
	spendTwo := globalData.spend[spendKeyTwo]

	result := &SwapInfo{
		RemainingOne: verifyResult.balanceOne - spendOne,
		RemainingTwo: verifyResult.balanceTwo - spendTwo,
		Id:           payId,
		TxId:         txId,
		Packed:       packedSwap,
		Payments:     payments,
	}

	// if already seen just return pay id and previous payments if present
	entry, ok := globalData.pendingTransactions[payId]
	if ok {
		if nil != entry.payments {
			result.Payments = entry.payments
		} else {
			// this would mean that reservoir data is corrupt
			logger.Panicf("storeSwap: failed to get current payment data for: %s  payid: %s", txId, payId)
		}
		return result, true, nil
	}

	// if duplicates were detected, but different duplicates were present
	// then it is an error
	if duplicate {
		return nil, true, fault.TransactionAlreadyExists
	}

	swapItem := &transactionData{
		txId:        txId,
		transaction: swap,
		packed:      packedSwap,
	}

	// already received the payment for the swap
	// approve the swap immediately if payment is ok
	detail, ok := globalData.orphanPayments[payId]
	if ok || globalData.autoVerify {
		if acceptablePayment(detail, payments) {
			globalData.verifiedTransactions[payId] = swapItem
			globalData.verifiedIndex[txId] = payId
			delete(globalData.pendingTransactions, payId)
			delete(globalData.pendingIndex, txId)
			delete(globalData.orphanPayments, payId)

			globalData.spend[spendKeyOne] += swap.QuantityOne
			globalData.spend[spendKeyTwo] += swap.QuantityTwo
			result.RemainingOne -= swap.QuantityOne
			result.RemainingTwo -= swap.QuantityTwo
			return result, false, nil
		}
	}

	// waiting for the payment to come
	payment := &transactionPaymentData{
		payId:     payId,
		tx:        swapItem,
		payments:  payments,
		expiresAt: time.Now().Add(constants.ReservoirTimeout),
	}

	globalData.pendingTransactions[payId] = payment
	globalData.pendingIndex[txId] = payId
	globalData.spend[spendKeyOne] += swap.QuantityOne
	globalData.spend[spendKeyTwo] += swap.QuantityTwo
	result.RemainingOne -= swap.QuantityOne
	result.RemainingTwo -= swap.QuantityTwo

	return result, false, nil
}

// CheckSwapBalances - check sufficient balance on both accounts to be able to execute a swap request
func CheckSwapBalances(trx storage.Transaction, swap *transactionrecord.ShareSwap, shareQuantityHandle storage.Handle) (uint64, uint64, error) {
	if nil == shareQuantityHandle {
		return 0, 0, fault.NilPointer
	}

	// check incoming quantity
	if 0 == swap.QuantityOne || 0 == swap.QuantityTwo {
		return 0, 0, fault.ShareQuantityTooSmall
	}

	oKeyOne := append(swap.OwnerOne.Bytes(), swap.ShareIdOne[:]...)
	var balanceOne uint64
	var ok bool
	if nil == trx {
		balanceOne, ok = shareQuantityHandle.GetN(oKeyOne)
	} else {
		balanceOne, ok = trx.GetN(shareQuantityHandle, oKeyOne)
	}

	// check if sufficient funds
	if !ok || balanceOne < swap.QuantityOne {
		return 0, 0, fault.InsufficientShares
	}

	oKeyTwo := append(swap.OwnerTwo.Bytes(), swap.ShareIdTwo[:]...)
	var balanceTwo uint64
	if nil == trx {
		balanceTwo, ok = shareQuantityHandle.GetN(oKeyTwo)
	} else {
		balanceTwo, ok = trx.GetN(shareQuantityHandle, oKeyTwo)
	}

	// check if sufficient funds
	if !ok || balanceTwo < swap.QuantityTwo {
		return 0, 0, fault.InsufficientShares
	}

	return balanceOne, balanceTwo, nil
}

// verify that a swap is ok
// ensure lock is held before calling
func verifySwap(swap *transactionrecord.ShareSwap, shareQuantityHandle storage.Handle, shareHandle storage.Handle, ownerDataHandle storage.Handle) (*verifiedSwapInfo, bool, error) {
	if nil == shareQuantityHandle || nil == shareHandle || nil == ownerDataHandle {
		return nil, false, fault.NilPointer
	}

	height := blockheader.Height()
	if swap.BeforeBlock <= height {
		return nil, false, fault.RecordHasExpired
	}

	balanceOne, balanceTwo, err := CheckSwapBalances(nil, swap, shareQuantityHandle)
	if nil != err {
		return nil, false, err
	}

	// pack swap and check signature
	packedSwap, err := swap.Pack(swap.OwnerOne)
	if nil != err {
		return nil, false, err
	}

	// transfer identifier and check for duplicate
	txId := packedSwap.MakeLink()

	// check for double spend
	_, okP := globalData.pendingIndex[txId]
	_, okV := globalData.verifiedIndex[txId]

	duplicate := false
	if okP {
		// if both then it is a possible duplicate
		// (depends on later pay id check)
		duplicate = true
	}

	// a single verified transfer fails the whole block
	if okV {
		return nil, false, fault.TransactionAlreadyExists
	}
	// a single confirmed transfer fails the whole block
	if storage.Pool.Transactions.Has(txId[:]) {
		return nil, false, fault.TransactionAlreadyExists
	}

	// log.Infof("share one: %x", swap.ShareOne)
	// log.Infof("share two: %x", swap.ShareTwo)

	// strip off the total leaving just normal ownerdata layout
	// ***** FIX THIS: only Share One for owner data to determine payment?
	// ***** FIX THIS: should share two's owner dat be used for double charge?
	// the owner data is under tx id of share record
	_ /*totalValue*/, shareTxId := shareHandle.GetNB(swap.ShareIdOne[:])
	if nil == shareTxId {
		return nil, false, fault.DoubleTransferAttempt
	}
	ownerData, err := ownership.GetOwnerDataB(nil, shareTxId, ownerDataHandle)
	if nil != err {
		return nil, false, fault.DoubleTransferAttempt
	}
	// log.Infof("ownerData: %x", ownerData)

	result := &verifiedSwapInfo{
		balanceOne:          balanceOne,
		balanceTwo:          balanceTwo,
		txId:                txId,
		packed:              packedSwap,
		issueTxId:           ownerData.IssueTxId(),
		transferBlockNumber: ownerData.TransferBlockNumber(),
		issueBlockNumber:    ownerData.IssueBlockNumber(),
	}
	return result, duplicate, nil
}
