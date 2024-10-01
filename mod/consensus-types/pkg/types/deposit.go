// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"encoding/binary"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// DepositSize is the size of the SSZ encoding of a Deposit.
const DepositSize = 192 // 48 + 32 + 8 + 96 + 8

// Compile-time assertions to ensure Deposit implements necessary interfaces.
var (
	_ ssz.StaticObject                    = &Deposit[engineprimitives.Log]{}
	_ constraints.SSZMarshallableRootable = &Deposit[engineprimitives.Log]{}
)

// Deposit into the consensus layer from the deposit contract in the execution
// layer.
type Deposit[LogT interface {
	GetData() []byte
}] struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
	// Signature of the deposit data.
	Signature crypto.BLSSignature `json:"signature"`
	// Index of the deposit in the deposit contract.
	Index uint64 `json:"index"`
}

// NewDeposit creates a new Deposit instance.
func NewDeposit[LogT interface {
	GetData() []byte
}](
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
	index uint64,
) *Deposit[LogT] {
	return &Deposit[LogT]{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
		Signature:   signature,
		Index:       index,
	}
}

// Empty creates an empty Deposit instance.
func (d *Deposit[LogT]) Empty() *Deposit[LogT] {
	return &Deposit[LogT]{}
}

// New creates a new Deposit instance.
func (d *Deposit[LogT]) New(
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
	index uint64,
) *Deposit[LogT] {
	return NewDeposit[LogT](
		pubkey, credentials, amount, signature, index,
	)
}

// VerifySignature verifies the deposit data and signature.
func (d *Deposit[_]) VerifySignature(
	forkData *ForkData,
	domainType common.DomainType,
	signatureVerificationFn func(
		pubkey crypto.BLSPubkey, message []byte, signature crypto.BLSSignature,
	) error,
) error {
	return (&DepositMessage{
		Pubkey:      d.Pubkey,
		Credentials: d.Credentials,
		Amount:      d.Amount,
	}).VerifyCreateValidator(
		forkData, d.Signature,
		domainType, signatureVerificationFn,
	)
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// DefineSSZ defines the SSZ encoding for the Deposit object.
func (d *Deposit[_]) DefineSSZ(c *ssz.Codec) {
	ssz.DefineStaticBytes(c, &d.Pubkey)
	ssz.DefineStaticBytes(c, &d.Credentials)
	ssz.DefineUint64(c, &d.Amount)
	ssz.DefineStaticBytes(c, &d.Signature)
	ssz.DefineUint64(c, &d.Index)
}

// MarshalSSZ marshals the Deposit object to SSZ format.
func (d *Deposit[_]) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, d.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, d)
}

// UnmarshalSSZ unmarshals the Deposit object from SSZ format.
func (d *Deposit[_]) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, d)
}

// SizeSSZ returns the SSZ encoded size of the Deposit object.
func (d *Deposit[_]) SizeSSZ() uint32 {
	return DepositSize
}

// HashTreeRoot computes the Merkleization of the Deposit object.
func (d *Deposit[_]) HashTreeRoot() common.Root {
	return ssz.HashSequential(d)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the Deposit object into a pre-allocated byte slice.
func (d *Deposit[_]) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := d.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the Deposit object with a hasher.
func (d *Deposit[_]) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Pubkey'
	hh.PutBytes(d.Pubkey[:])

	// Field (1) 'Credentials'
	hh.PutBytes(d.Credentials[:])

	// Field (2) 'Amount'
	hh.PutUint64(uint64(d.Amount))

	// Field (3) 'Signature'
	hh.PutBytes(d.Signature[:])

	// Field (4) 'Index'
	hh.PutUint64(d.Index)

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the Deposit object.
func (d *Deposit[_]) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(d)
}

/* -------------------------------------------------------------------------- */
/*                                   EthLog                                   */
/* -------------------------------------------------------------------------- */

// UnmarshalLog unmarshals the Deposit object from an Ethereum log.
func (d *Deposit[LogT]) UnmarshalLog(log LogT) error {
	data := log.GetData()
	idx := binary.BigEndian.Uint64(data[152:160])
	d.Index = idx
	d.Pubkey = bytes.B48(data[208:256])
	amount := binary.BigEndian.Uint64(data[280:288])
	d.Amount = math.U64(amount)
	d.Credentials = WithdrawalCredentials(bytes.B32(data[288:320]))
	d.Signature = bytes.B96(data[352:448])
	return nil
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// GetAmount returns the deposit amount in gwei.
func (d *Deposit[_]) GetAmount() math.Gwei {
	return d.Amount
}

// GetPubkey returns the public key of the validator specified in the deposit.
func (d *Deposit[_]) GetPubkey() crypto.BLSPubkey {
	return d.Pubkey
}

// GetIndex returns the index of the deposit in the deposit contract.
func (d *Deposit[_]) GetIndex() math.U64 {
	return math.U64(d.Index)
}

// GetSignature returns the signature of the deposit data.
func (d *Deposit[_]) GetSignature() crypto.BLSSignature {
	return d.Signature
}

// GetWithdrawalCredentials returns the staking credentials of the deposit.
func (d *Deposit[_]) GetWithdrawalCredentials() WithdrawalCredentials {
	return d.Credentials
}
