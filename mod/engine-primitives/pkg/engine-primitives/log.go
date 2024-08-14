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

package engineprimitives

import (
	"encoding/json"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Log represents a contract log event. These events are generated by the LOG
// opcode and
// stored/indexed by the node.
type Log struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address common.ExecutionAddress `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []common.ExecutionHash `json:"topics"  gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data []byte `json:"data"    gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"      rlp:"-"`
	// hash of the transaction
	TxHash common.ExecutionHash `json:"transactionHash"  rlp:"-" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" rlp:"-"`
	// hash of the block in which the transaction was included
	BlockHash common.ExecutionHash `json:"blockHash"        rlp:"-"`
	// index of the log in the block
	Index uint `json:"logIndex"         rlp:"-"`

	// The Removed field is true if this log was reverted due to a chain
	// reorganisation. You must pay attention to this field if you receive logs
	// through a filter query.
	Removed bool `json:"removed" rlp:"-"`
}

// MarshalJSON marshals as JSON.
func (l Log) MarshalJSON() ([]byte, error) {
	type Log struct {
		Address     common.ExecutionAddress `json:"address"          gencodec:"required"`
		Topics      []common.ExecutionHash  `json:"topics"           gencodec:"required"`
		Data        bytes.Bytes             `json:"data"             gencodec:"required"`
		BlockNumber math.U64                `json:"blockNumber"                          rlp:"-"`
		TxHash      common.ExecutionHash    `json:"transactionHash"  gencodec:"required" rlp:"-"`
		TxIndex     math.U64                `json:"transactionIndex"                     rlp:"-"`
		BlockHash   common.ExecutionHash    `json:"blockHash"                            rlp:"-"`
		Index       math.U64                `json:"logIndex"                             rlp:"-"`
		Removed     bool                    `json:"removed"                              rlp:"-"`
	}
	var enc Log
	enc.Address = l.Address
	enc.Topics = l.Topics
	enc.Data = l.Data
	enc.BlockNumber = math.U64(l.BlockNumber)
	enc.TxHash = l.TxHash
	enc.TxIndex = math.U64(l.TxIndex)
	enc.BlockHash = l.BlockHash
	enc.Index = math.U64(l.Index)
	enc.Removed = l.Removed
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (l *Log) UnmarshalJSON(input []byte) error {
	type Log struct {
		Address     *common.ExecutionAddress `json:"address"          gencodec:"required"`
		Topics      []common.ExecutionHash   `json:"topics"           gencodec:"required"`
		Data        *bytes.Bytes             `json:"data"             gencodec:"required"`
		BlockNumber *math.U64                `json:"blockNumber"                          rlp:"-"`
		TxHash      *common.ExecutionHash    `json:"transactionHash"  gencodec:"required" rlp:"-"`
		TxIndex     *math.U64                `json:"transactionIndex"                     rlp:"-"`
		BlockHash   *common.ExecutionHash    `json:"blockHash"                            rlp:"-"`
		Index       *math.U64                `json:"logIndex"                             rlp:"-"`
		Removed     *bool                    `json:"removed"                              rlp:"-"`
	}
	var dec Log
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Address == nil {
		return errors.New("missing required field 'address' for Log")
	}
	l.Address = *dec.Address
	if dec.Topics == nil {
		return errors.New("missing required field 'topics' for Log")
	}
	l.Topics = dec.Topics
	if dec.Data == nil {
		return errors.New("missing required field 'data' for Log")
	}
	l.Data = *dec.Data
	if dec.BlockNumber != nil {
		l.BlockNumber = uint64(*dec.BlockNumber)
	}
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionHash' for Log")
	}
	l.TxHash = *dec.TxHash
	if dec.TxIndex != nil {
		l.TxIndex = uint(*dec.TxIndex)
	}
	if dec.BlockHash != nil {
		l.BlockHash = *dec.BlockHash
	}
	if dec.Index != nil {
		l.Index = uint(*dec.Index)
	}
	if dec.Removed != nil {
		l.Removed = *dec.Removed
	}
	return nil
}
