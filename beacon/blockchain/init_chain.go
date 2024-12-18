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

package blockchain

import (
	"context"
	"encoding/json"

	"github.com/berachain/beacon-kit/primitives/transition"
)

// ProcessGenesisData processes the genesis state and initializes the beacon
// state.
func (s *Service[
	_, _, _, _, _, GenesisT, _, _, _,
]) ProcessGenesisData(
	ctx context.Context,
	bytes []byte,
) (transition.ValidatorUpdates, error) {
	genesisData := *new(GenesisT)
	if err := json.Unmarshal(bytes, &genesisData); err != nil {
		s.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	// Store the genesis deposits.
	genesisDepositDatas := genesisData.GetDepositDatas()
	genesisExecutionPayloadHeader := genesisData.GetExecutionPayloadHeader()
	if err := s.depositStore.EnqueueDepositDatas(genesisDepositDatas); err != nil {
		s.logger.Error("Failed to store genesis deposits", "error", err)
		return nil, err
	}

	// Get the genesis deposits, with their proofs.
	genesisDeposits, genesisDepositsRoot, err := s.depositStore.GetDepositsByIndex(
		0, uint64(len(genesisDepositDatas)),
	)
	if err != nil {
		s.logger.Error("Failed to retrieve genesis deposits with proofs", "error", err)
		return nil, err
	}

	return s.stateProcessor.InitializePreminedBeaconStateFromEth1(
		s.storageBackend.StateFromContext(ctx),
		genesisDeposits,
		genesisDepositsRoot,
		genesisExecutionPayloadHeader,
		genesisData.GetForkVersion(),
	)
}
