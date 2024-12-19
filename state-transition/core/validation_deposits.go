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

package core

import (
	"bytes"
	"fmt"

	"github.com/berachain/beacon-kit/config/spec"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

func (sp *StateProcessor[
	_, BeaconStateT, _, _,
]) validateGenesisDeposits(
	st BeaconStateT,
	deposits []*ctypes.Deposit,
) error {
	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// Bartio does not properly validate deposits index
		// We skip checks for backward compatibility
		return nil

	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID:
		// Boonet inherited the bug from Bartio and it may have added some
		// validators before we activate the fork. So we skip all validations
		// but the validator set cap.
		//#nosec:G701 // can't overflow.
		if uint64(len(deposits)) > sp.cs.ValidatorSetCap() {
			return fmt.Errorf("validator set cap %d, deposits count %d: %w",
				sp.cs.ValidatorSetCap(),
				len(deposits),
				ErrValSetCapExceeded,
			)
		}
		return nil

	default:
		// TODO: improve error handling by distinguishing
		// ErrNotFound from other kind of errors
		if _, err := st.GetEth1DepositIndex(); err == nil {
			// there should not be Eth1DepositIndex stored before
			// genesis first deposit
			return errors.Wrap(
				ErrDepositMismatch,
				"Eth1DepositIndex should be unset at genesis",
			)
		}
		if len(deposits) == 0 {
			// there should be at least a validator in genesis
			return errors.Wrap(
				ErrDepositsLengthMismatch,
				"at least one validator should be in genesis",
			)
		}
		for i, deposit := range deposits {
			// deposit indices should be contiguous
			if deposit.GetIndex() != math.U64(i) {
				return errors.Wrapf(
					ErrDepositIndexOutOfOrder,
					"genesis deposit index: %d, expected index: %d",
					deposit.GetIndex().Unwrap(), i,
				)
			}
		}

		// BeaconKit enforces a cap on the validator set size.
		// If genesis deposits breaches the cap we return an error.
		//#nosec:G701 // can't overflow.
		if uint64(len(deposits)) > sp.cs.ValidatorSetCap() {
			return fmt.Errorf("validator set cap %d, deposits count %d: %w",
				sp.cs.ValidatorSetCap(),
				len(deposits),
				ErrValSetCapExceeded,
			)
		}
		return nil
	}
}

func (sp *StateProcessor[
	_, BeaconStateT, _, _,
]) validateNonGenesisDeposits(
	st BeaconStateT,
	blkDeposits []*ctypes.Deposit,
	blkDepositRoot common.Root,
) (common.Root, error) {
	slot, err := st.GetSlot()
	if err != nil {
		return common.Root{}, fmt.Errorf(
			"failed loading slot while processing deposits: %w", err,
		)
	}
	switch {
	case sp.cs.DepositEth1ChainID() == spec.BartioChainID:
		// Bartio does not properly validate deposits index
		// We skip checks for backward compatibility
		return common.Root{}, nil

	case sp.cs.DepositEth1ChainID() == spec.BoonetEth1ChainID &&
		slot < math.U64(spec.BoonetFork2Height):
		// Boonet inherited the bug from Bartio and it may have added some
		// validators before we activate the fork. So we skip validation
		// before fork activation
		return common.Root{}, nil

	default:
		var depositIndex uint64
		depositIndex, err = st.GetEth1DepositIndex()
		if err != nil {
			return common.Root{}, err
		}
		depositIndex += uint64(len(blkDeposits) + 1)

		var deposits []*ctypes.Deposit
		deposits, err = sp.ds.GetDepositsByIndex(0, depositIndex)
		if err != nil {
			return common.Root{}, err
		}

		newRoot := ctypes.Deposits(deposits).HashTreeRoot()
		if !bytes.Equal(blkDepositRoot[:], newRoot[:]) {
			return common.Root{}, ErrDepositsRootMismatch
		}
		return newRoot, nil
	}
}
