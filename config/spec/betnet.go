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

package spec

import (
	"github.com/berachain/beacon-kit/chain-spec/chain"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// BetnetChainSpec is the ChainSpec for the localnet.
func BetnetChainSpec() (chain.Spec[
	common.DomainType,
	math.Epoch,
	common.ExecutionAddress,
	math.Slot,
	any,
], error) {
	betnetSpec := BaseSpec()

	betnetSpec.DepositEth1ChainID = BetnetEth1ChainID

	betnetSpec.EVMInflationAddress = common.NewExecutionAddressFromHex(
		"0x289274787bAF083C15A45a174b7a8e44F0720660",
	)
	betnetSpec.EVMInflationPerBlock = 2.5e9

	betnetSpec.ValidatorSetCap = 5
	betnetSpec.MaxValidatorsPerWithdrawalsSweepPostUpgrade = 3

	betnetSpec.SlotsPerEpoch = 64

	return chain.NewChainSpec(betnetSpec)
}
