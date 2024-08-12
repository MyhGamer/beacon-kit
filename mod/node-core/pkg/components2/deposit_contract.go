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

package components

import (
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// BeaconDepositContractInput is the input for the beacon deposit contract
// for the dep inject framework.
type BeaconDepositContractInput struct {
	depinject.In
	ChainSpec    common.ChainSpec
	EngineClient *EngineClient
}

// ProvideBeaconDepositContract provides a beacon deposit contract through the
// dep inject framework.
func ProvideBeaconDepositContract[
	DepositT Deposit[DepositT, ForkDataT, WithdrawalCredentialsT],
	ForkDataT any,
	WithdrawalCredentialsT ~[32]byte,
](
	in BeaconDepositContractInput,
) (deposit.Contract[DepositT], error) {
	// Build the deposit contract.
	return deposit.NewWrappedBeaconDepositContract[
		DepositT,
		WithdrawalCredentialsT,
	](
		in.ChainSpec.DepositContractAddress(),
		in.EngineClient,
	)
}