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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// DepositDataSize is the size of the SSZ encoding of a DepositData.
const DepositDataSize = 184 // 48 + 32 + 8 + 96

// Compile-time assertions to ensure Deposit implements necessary interfaces.
var (
	_ ssz.StaticObject                    = (*DepositData)(nil)
	_ constraints.SSZMarshallableRootable = (*DepositData)(nil)
)

// Deposit into the consensus layer from the deposit contract in the execution
// layer.
type DepositData struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials WithdrawalCredentials `json:"credentials"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
	// Signature of the deposit data.
	Signature crypto.BLSSignature `json:"signature"`
}

// NewDepositData creates a new DepositData instance.
func NewDepositData(
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
) *DepositData {
	return &DepositData{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
		Signature:   signature,
	}
}

// Empty creates an empty Deposit instance.
func (d *DepositData) Empty() *DepositData {
	return &DepositData{}
}

// New creates a new Deposit instance.
func (d *DepositData) New(
	pubkey crypto.BLSPubkey,
	credentials WithdrawalCredentials,
	amount math.Gwei,
	signature crypto.BLSSignature,
) *DepositData {
	return NewDepositData(
		pubkey, credentials, amount, signature,
	)
}

// VerifySignature verifies the deposit data and signature.
func (d *DepositData) VerifySignature(
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
func (d *DepositData) DefineSSZ(c *ssz.Codec) {
	ssz.DefineStaticBytes(c, &d.Pubkey)
	ssz.DefineStaticBytes(c, &d.Credentials)
	ssz.DefineUint64(c, &d.Amount)
	ssz.DefineStaticBytes(c, &d.Signature)
}

// MarshalSSZ marshals the DepositData object to SSZ format.
func (d *DepositData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(d))
	return buf, ssz.EncodeToBytes(buf, d)
}

// UnmarshalSSZ unmarshals the DepositData object from SSZ format.
func (d *DepositData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, d)
}

// SizeSSZ returns the SSZ encoded size of the DepositData object.
func (d *DepositData) SizeSSZ(*ssz.Sizer) uint32 {
	return DepositDataSize
}

// HashTreeRoot computes the Merkleization of the DepositData object.
func (d *DepositData) HashTreeRoot() common.Root {
	return ssz.HashSequential(d)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the DepositData object into a pre-allocated byte slice.
func (d *DepositData) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := d.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the DepositData object with a hasher.
func (d *DepositData) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Pubkey'
	hh.PutBytes(d.Pubkey[:])

	// Field (1) 'Credentials'
	hh.PutBytes(d.Credentials[:])

	// Field (2) 'Amount'
	hh.PutUint64(uint64(d.Amount))

	// Field (3) 'Signature'
	hh.PutBytes(d.Signature[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the DepositData object.
func (d *DepositData) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(d)
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// GetAmount returns the deposit amount in gwei.
func (d *DepositData) GetAmount() math.Gwei {
	return d.Amount
}

// GetPubkey returns the public key of the validator specified in the deposit.
func (d *DepositData) GetPubkey() crypto.BLSPubkey {
	return d.Pubkey
}

// GetSignature returns the signature of the deposit data.
func (d *DepositData) GetSignature() crypto.BLSSignature {
	return d.Signature
}

// GetWithdrawalCredentials returns the staking credentials of the deposit.
func (d *DepositData) GetWithdrawalCredentials() WithdrawalCredentials {
	return d.Credentials
}

// HasEth1WithdrawalCredentials returns true if the deposit has eth1 withdrawal
// credentials.
func (d *DepositData) HasEth1WithdrawalCredentials() bool {
	return d.Credentials[0] == EthSecp256k1CredentialPrefix
}
