// Copyright (c) 2014-2018 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package transactionrecord

import (
	"strings"
	"unicode/utf8"

	"github.com/bitmark-inc/bitmarkd/account"
	"github.com/bitmark-inc/bitmarkd/currency"
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/util"
)

// supported currency sets
// code here will support all versions
var versions = []currency.Set{
	currency.MakeSet(),                                    // 0
	currency.MakeSet(currency.Bitcoin, currency.Litecoin), // 1
}

// currently supported block foundation version (used by proofer)
const (
	FoundationVersion = 1
)

// pack BaseData
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (baseData *OldBaseData) Pack(address *account.Account) (Packed, error) {
	if len(baseData.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	// prevent nil or zero account
	if nil == baseData.Owner || nil == address || baseData.Owner.IsZero() || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	err := baseData.Currency.ValidateAddress(baseData.PaymentAddress, address.IsTesting())
	if nil != err {
		return nil, err
	}

	// concatenate bytes
	message := createPacked(BaseDataTag)
	message.appendUint64(baseData.Currency.Uint64())
	message.appendString(baseData.PaymentAddress)
	message.appendAccount(baseData.Owner)
	message.appendUint64(baseData.Nonce)

	// signature
	err = address.CheckSignature(message, baseData.Signature)
	if nil != err {
		return message, err
	}
	// Signature Last
	return *message.appendBytes(baseData.Signature), nil
}

// pack AssetData
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last.
//
// Note: the metadata field consists of key value pairs each preceded
//       by its count (
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (assetData *AssetData) Pack(address *account.Account) (Packed, error) {
	if len(assetData.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}
	// prevent nil or zero account
	if nil == assetData.Registrant || nil == address || assetData.Registrant.IsZero() || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	if utf8.RuneCountInString(assetData.Name) < minNameLength {
		return nil, fault.ErrNameTooShort
	}
	if utf8.RuneCountInString(assetData.Name) > maxNameLength {
		return nil, fault.ErrNameTooLong
	}

	if utf8.RuneCountInString(assetData.Fingerprint) < minFingerprintLength {
		return nil, fault.ErrFingerprintTooShort
	}
	if utf8.RuneCountInString(assetData.Fingerprint) > maxFingerprintLength {
		return nil, fault.ErrFingerprintTooLong
	}

	if utf8.RuneCountInString(assetData.Metadata) > maxMetadataLength {
		return nil, fault.ErrMetadataTooLong
	}

	// check that metadata contains a vailid map:
	// i.e.  key1 <NUL> value1 <NUL> key2 <NUL> value2 <NUL> … keyN <NUL> valueN
	// Notes: 1: no NUL after last value
	//        2: no empty key or value is allowed
	if 0 != len(assetData.Metadata) {
		splitMetadata := strings.Split(assetData.Metadata, "\u0000")
		if 1 == len(splitMetadata)%2 {
			return nil, fault.ErrMetadataIsNotMap
		}
		for _, v := range splitMetadata {
			if 0 == len(v) {
				return nil, fault.ErrMetadataIsNotMap
			}
		}
	}

	// concatenate bytes
	message := createPacked(AssetDataTag)
	message.appendString(assetData.Name)
	message.appendString(assetData.Fingerprint)
	message.appendString(assetData.Metadata)
	message.appendAccount(assetData.Registrant)

	// signature
	err := address.CheckSignature(message, assetData.Signature)
	if nil != err {
		return message, err
	}
	// Signature Last
	return *message.appendBytes(assetData.Signature), nil
}

// pack BitmarkIssue
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (issue *BitmarkIssue) Pack(address *account.Account) (Packed, error) {
	if len(issue.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	// prevent nil or zero account
	if nil == issue.Owner || nil == address || issue.Owner.IsZero() || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	// concatenate bytes
	message := createPacked(BitmarkIssueTag)
	message.appendBytes(issue.AssetId[:])
	message.appendAccount(issue.Owner)
	message.appendUint64(issue.Nonce)

	// signature
	err := address.CheckSignature(message, issue.Signature)
	if nil != err {
		return message, err
	}

	// Signature Last
	return *message.appendBytes(issue.Signature), nil
}

// local function to pack BitmarkTransfer
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (transfer *BitmarkTransferUnratified) Pack(address *account.Account) (Packed, error) {
	if len(transfer.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	// Note: In this case Owner can be zero ⇒ bitmark is destroyed
	//       and no further transfers are allowed.
	//       theddress cannot be zero to prevent discovery of the
	//       corresponding private key being able to transfer all
	//       previously destroyed bitmarks to a new account.
	if nil == transfer.Owner || nil == address || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	testnet := address.IsTesting()

	// concatenate bytes
	message := createPacked(BitmarkTransferUnratifiedTag)
	message.appendBytes(transfer.Link[:])
	_, err := message.appendEscrow(transfer.Escrow, testnet)
	if nil != err {
		return nil, err
	}
	message.appendAccount(transfer.Owner)

	// signature
	err = address.CheckSignature(message, transfer.Signature)
	if nil != err {
		return message, err
	}

	// Signature Last
	return *message.appendBytes(transfer.Signature), nil
}

// local function to pack BitmarkTransferCountersigned
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (transfer *BitmarkTransferCountersigned) Pack(address *account.Account) (Packed, error) {
	if len(transfer.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	if len(transfer.Countersignature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	// Note: impossible to have 2 signature transfer to zero public key
	if nil == transfer.Owner || nil == address || transfer.Owner.IsZero() || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	testnet := address.IsTesting()

	// concatenate bytes
	message := createPacked(BitmarkTransferCountersignedTag)
	message.appendBytes(transfer.Link[:])
	_, err := message.appendEscrow(transfer.Escrow, testnet)
	if nil != err {
		return nil, err
	}
	message.appendAccount(transfer.Owner)

	// signature
	err = address.CheckSignature(message, transfer.Signature)
	if nil != err {
		return message, err
	}

	// add signature Signature
	message.appendBytes(transfer.Signature)

	err = transfer.Owner.CheckSignature(message, transfer.Countersignature)
	if nil != err {
		return message, err
	}

	// Countersignature Last
	return *message.appendBytes(transfer.Countersignature), nil
}

// pack BlockFoundation
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (foundation *BlockFoundation) Pack(address *account.Account) (Packed, error) {
	if len(foundation.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	// prevent nil or zero account
	if nil == foundation.Owner || nil == address || foundation.Owner.IsZero() || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	err := CheckPayments(foundation.Version, address.IsTesting(), foundation.Payments)
	if nil != err {
		return nil, err
	}
	packedPayments, err := foundation.Payments.Pack(address.IsTesting())
	if nil != err {
		return nil, err
	}

	// concatenate bytes
	message := createPacked(BlockFoundationTag)
	message.appendUint64(foundation.Version)
	message.appendBytes(packedPayments)
	message.appendAccount(foundation.Owner)
	message.appendUint64(foundation.Nonce)

	// signature
	err = address.CheckSignature(message, foundation.Signature)
	if nil != err {
		return message, err
	}
	// Signature Last
	return *message.appendBytes(foundation.Signature), nil
}

// pack BlockOwnerTransfer
//
// Pack Varint64(tag) followed by fields in order as struct above with
// signature last
//
// NOTE: returns the "unsigned" message on signature failure - for
//       debugging/testing
func (transfer *BlockOwnerTransfer) Pack(address *account.Account) (Packed, error) {
	if len(transfer.Signature) > maxSignatureLength {
		return nil, fault.ErrSignatureTooLong
	}

	// prevent nil or zero account
	if nil == transfer.Owner || nil == address || transfer.Owner.IsZero() || address.IsZero() {
		return nil, fault.ErrInvalidOwnerOrRegistrant
	}

	err := CheckPayments(transfer.Version, address.IsTesting(), transfer.Payments)
	if nil != err {
		return nil, err
	}

	packedPayments, err := transfer.Payments.Pack(address.IsTesting())
	if nil != err {
		return nil, err
	}

	testnet := address.IsTesting()

	// concatenate bytes
	message := createPacked(BlockOwnerTransferTag)
	message.appendBytes(transfer.Link[:])
	_, err = message.appendEscrow(transfer.Escrow, testnet)
	if nil != err {
		return nil, err
	}
	message.appendUint64(transfer.Version)
	message.appendBytes(packedPayments)
	message.appendAccount(transfer.Owner)

	// signature
	err = address.CheckSignature(message, transfer.Signature)
	if nil != err {
		return message, err
	}
	message.appendBytes(transfer.Signature)

	err = transfer.Owner.CheckSignature(message, transfer.Countersignature)
	if nil != err {
		return message, err
	}

	// Countersignature Last
	return *message.appendBytes(transfer.Countersignature), nil
}

// internal routines below here
// ----------------------------

// check all currency addresses for correct network and validity
func CheckPayments(version uint64, testnet bool, payments currency.Map) error {
	// validate version
	if version < 1 || version >= uint64(len(versions)) {
		return fault.ErrInvalidPaymentVersion
	}

	cs := currency.MakeSet()
	for currency, address := range payments {

		err := currency.ValidateAddress(address, testnet)
		if nil != err {
			return err
		}

		// if a duplicate currency value
		if cs.Add(currency) {
			return fault.ErrInvalidCurrencyAddress
		}
	}

	// validate the set of supplied currencies
	if versions[version] != cs {
		return fault.ErrInvalidCurrencyAddress
	}

	return nil
}

// create a new packed buffer
func createPacked(tag TagType) Packed {
	return util.ToVarint64(uint64(tag))
}

// append a single field to a buffer
//
// the field is prefixed by Varint64(length)
func (buffer *Packed) appendString(s string) *Packed {
	l := util.ToVarint64(uint64(len(s)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, s...)
	return buffer
}

// append an address to a buffer
//
// the field is prefixed by Varint64(length)
func (buffer *Packed) appendAccount(address *account.Account) *Packed {
	data := address.Bytes()
	l := util.ToVarint64(uint64(len(data)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, data...)
	return buffer
}

// append a bytes to a buffer
//
// the field is prefixed by Varint64(length)
func (buffer *Packed) appendBytes(data []byte) *Packed {
	l := util.ToVarint64(uint64(len(data)))
	*buffer = append(*buffer, l...)
	*buffer = append(*buffer, data...)
	return buffer
}

// append a Varint64 to buffer
func (buffer *Packed) appendUint64(value uint64) *Packed {
	valueBytes := util.ToVarint64(value)
	*buffer = append(*buffer, valueBytes...)
	return buffer
}

// append a Escrop[ payment to buffer
func (buffer *Packed) appendEscrow(escrow *Payment, testnet bool) (*Packed, error) {

	if nil == escrow {
		*buffer = append(*buffer, 0)
	} else {
		err := escrow.Currency.ValidateAddress(escrow.Address, testnet)
		if nil != err {
			return nil, err
		}
		*buffer = append(*buffer, 1)
		buffer.appendUint64(escrow.Currency.Uint64())
		buffer.appendString(escrow.Address)
		buffer.appendUint64(escrow.Amount)
	}
	return buffer, nil
}
