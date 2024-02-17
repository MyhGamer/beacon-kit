// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package staking

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/itsdevbear/bolaris/runtime/modules/beacon/keeper/store"
)

// ProcessDeposit processes a deposit log from the execution layer
// and puts the deposit to the beacon state.
func (s *Service) ProcessDeposit(
	ctx context.Context,
	validatorPubkey []byte,
	withdrawalCredentials []byte,
	amount uint64,
	nonce uint64,
) error {
	beaconState := s.BeaconState(ctx)
	expectedNonce := beaconState.GetStakingNonce()
	// We may receive the same deposit log twice from the execution layer, just ignore it.
	if nonce < expectedNonce {
		return nil
	}
	// The deposit log does not come in order.
	if nonce != expectedNonce {
		return errors.Wrapf(
			ErrInvalidNonce, "expected nonce %d, got %d", expectedNonce, nonce,
		)
	}
	// Increase the staking nonce before AddDeposit,
	// so that it is ready to accept the nexr deposit,
	// even when AddDeposit returns an error.
	beaconState.SetStakingNonce(expectedNonce + 1)
	deposit := store.NewDeposit(
		validatorPubkey,
		amount,
		withdrawalCredentials,
	)
	// Add the deposit to the queue.
	err := beaconState.AddDeposit(deposit)
	if err != nil {
		return err
	}
	s.Logger().Info("delegating from execution layer",
		"validatorPubkey", validatorPubkey, "amount", amount, "nonce", nonce)
	return nil
}
