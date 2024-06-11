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

package builder

import (
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/payload/pkg/cache"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// PayloadBuilder is used to build payloads on the
// execution client.
type PayloadBuilder[
	BeaconStateT BeaconState[ExecutionPayloadHeaderT, WithdrawalT],
	ExecutionPayloadT interface {
		IsNil() bool
		Empty(uint32) ExecutionPayloadT
		GetBlockHash() common.ExecutionHash
		GetFeeRecipient() common.ExecutionAddress
		GetParentHash() common.ExecutionHash
	},
	ExecutionPayloadHeaderT interface {
		GetBlockHash() common.ExecutionHash
		GetParentHash() common.ExecutionHash
	},
	WithdrawalT any,
] struct {
	// cfg holds the configuration settings for the PayloadBuilder.
	cfg *Config
	// chainSpec holds the chain specifications for the PayloadBuilder.
	chainSpec primitives.ChainSpec
	// logger is used for logging within the PayloadBuilder.
	logger log.Logger[any]
	// ee is the execution engine.
	ee ExecutionEngine[ExecutionPayloadT]
	// pc is the payload ID cache, it is used to store
	// "in-flight" payloads that are being built on
	// the execution client.
	pc *cache.PayloadIDCache[
		engineprimitives.PayloadID, [32]byte, math.Slot,
	]
}

// New creates a new service.
func New[
	BeaconStateT BeaconState[ExecutionPayloadHeaderT, WithdrawalT],
	ExecutionPayloadT interface {
		IsNil() bool
		Empty(uint32) ExecutionPayloadT
		GetBlockHash() common.ExecutionHash
		GetParentHash() common.ExecutionHash
		GetFeeRecipient() common.ExecutionAddress
	},
	ExecutionPayloadHeaderT interface {
		GetBlockHash() common.ExecutionHash
		GetParentHash() common.ExecutionHash
	},
	WithdrawalT any,
](
	cfg *Config,
	chainSpec primitives.ChainSpec,
	logger log.Logger[any],
	ee ExecutionEngine[ExecutionPayloadT],
	pc *cache.PayloadIDCache[
		engineprimitives.PayloadID, [32]byte, math.Slot,
	],
) *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
] {
	return &PayloadBuilder[
		BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
	]{
		cfg:       cfg,
		chainSpec: chainSpec,
		logger:    logger,
		ee:        ee,
		pc:        pc,
	}
}

// Enabled returns true if the payload builder is enabled.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT, WithdrawalT,
]) Enabled() bool {
	return pb.cfg.Enabled
}
