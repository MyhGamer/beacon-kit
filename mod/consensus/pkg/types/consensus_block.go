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
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type ConsensusBlock[BeaconBlockT any] struct {
	blk BeaconBlockT

	// blkTime assigned by CometBFT to the beacon block
	blkTime math.U64

	// block proposer address assigned  by CometBFT to the beacon block
	blkProposerAddress []byte
}

// New creates a new SlotData instance.
func (b *ConsensusBlock[BeaconBlockT]) New(
	beaconBlock BeaconBlockT,
	blkTime time.Time,
	blkProposerAddress []byte,
) *ConsensusBlock[BeaconBlockT] {
	b = &ConsensusBlock[BeaconBlockT]{
		blk:                beaconBlock,
		blkTime:            math.U64(blkTime.Unix()),
		blkProposerAddress: blkProposerAddress,
	}
	return b
}

func (b *ConsensusBlock[BeaconBlockT]) GetBeaconBlock() BeaconBlockT {
	return b.blk
}

func (b *ConsensusBlock[_]) GetConsensusBlockTime() math.U64 {
	return b.blkTime
}

// TODO: harden the return type
func (b *ConsensusBlock[_]) GetConsensusProposerAddress() []byte {
	return b.blkProposerAddress
}