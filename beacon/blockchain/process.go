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
	"time"

	"github.com/berachain/beacon-kit/consensus-types/types"
	consensusTypes "github.com/berachain/beacon-kit/consensus/types"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
)

// ProcessGenesisData processes the genesis state and initializes the beacon
// state.
func (s *Service[
	_, _, _, _, _, _, _, _, _, _, _, _, GenesisT, _, _, _,
]) ProcessGenesisData(
	ctx context.Context,
	bytes []byte,
) (transition.ValidatorUpdates, error) {
	genesisData := *new(GenesisT)
	if err := json.Unmarshal(bytes, &genesisData); err != nil {
		s.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}
	return s.stateProcessor.InitializePreminedBeaconStateFromEth1(
		s.storageBackend.StateFromContext(ctx),
		genesisData.GetDeposits(),
		genesisData.GetExecutionPayloadHeader(),
		genesisData.GetForkVersion(),
	)
}

func (s *Service[
	_, _, ConsensusBlockT, BeaconBlockT, _, _, _, _, _, _, _, _, GenesisT, ConsensusSidecarsT, _, _,
]) ProcessProposal(
	ctx context.Context,
	blk consensusTypes.ConsensusBlock[*types.BeaconBlock],
	cs *consensusTypes.ConsensusSidecars[*datypes.BlobSidecars, *types.BeaconBlockHeader],
) (*cmtabci.ProcessProposalResponse, error) {
	s.logger.Info("!!!!!!!! ProcessProposal called")
	s.logger.Info("blk", blk)
	s.logger.Info("sidecars", cs)

	// verifySidecars
	//
	sidecars := cs.GetSidecars()
	if !sidecars.IsNil() && sidecars.Len() > 0 {
		s.logger.Info("Received incoming blob sidecars")

		// TODO: Remove this once we cleanup generics
		var sidecarI interface{} = cs
		data, ok := sidecarI.(ConsensusSidecarsT)
		if !ok {
			panic("could not convert sidecar to ConsensusSidecarsT")
		}

		// Verify the blobs and ensure they match the local state.
		if err := s.blobProcessor.VerifySidecars(data); err != nil {
			s.logger.Error(
				"rejecting incoming blob sidecars",
				"reason", err,
			)
			return nil, err
		}

		s.logger.Info(
			"Blob sidecars verification succeeded - accepting incoming blob sidecars",
			"num_blobs",
			sidecars.Len(),
		)
	}

	err := s.VerifyIncomingBlock(ctx, blk)
	if err != nil {
		s.logger.Error("failed to verify incoming block", "error", err)
		return nil, err
	}

	return nil, nil

}

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service[
	_, _, ConsensusBlockT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) ProcessBeaconBlock(
	ctx context.Context,
	blk consensusTypes.ConsensusBlock[*types.BeaconBlock],
) (transition.ValidatorUpdates, error) {
	beaconBlk := blk.GetBeaconBlock()

	// If the block is nil, exit early.
	if beaconBlk.IsNil() {
		return nil, ErrNilBlk
	}

	st := s.storageBackend.StateFromContext(ctx)
	valUpdates, err := s.executeStateTransition(ctx, st, blk)
	if err != nil {
		return nil, err
	}

	// If the blobs needed to process the block are not available, we
	// return an error. It is safe to use the slot off of the beacon block
	// since it has been verified as correct already.
	if !s.storageBackend.AvailabilityStore().IsDataAvailable(
		ctx, beaconBlk.GetSlot(), beaconBlk.GetBody(),
	) {
		return nil, ErrDataNotAvailable
	}

	// fetch and store the deposit for the block
	blockNum := beaconBlk.GetBody().GetExecutionPayload().GetNumber()
	s.depositFetcher(ctx, blockNum)

	// store the finalized block in the KVStore.
	slot := beaconBlk.GetSlot()
	if err = s.blockStore.Set(beaconBlk); err != nil {
		s.logger.Error(
			"failed to store block", "slot", slot, "error", err,
		)
	}

	// prune the availability and deposit store
	err = s.processPruning(beaconBlk)
	if err != nil {
		s.logger.Error("failed to processPruning", "error", err)
	}

	go s.sendPostBlockFCU(ctx, st, blk)

	return valUpdates.CanonicalSort(), nil
}

// executeStateTransition runs the stf.
func (s *Service[
	_, _, ConsensusBlockT, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _,
]) executeStateTransition(
	ctx context.Context,
	st BeaconStateT,
	blk consensusTypes.ConsensusBlock[*types.BeaconBlock],
) (transition.ValidatorUpdates, error) {
	startTime := time.Now()
	defer s.metrics.measureStateTransitionDuration(startTime)
	valUpdates, err := s.stateProcessor.Transition(
		&transition.Context{
			Context: ctx,

			// We set `OptimisticEngine` to true since this is called during
			// FinalizeBlock. We want to assume the payload is valid. If it
			// ends up not being valid later, the node will simply AppHash,
			// which is completely fine. This means we were syncing from a
			// bad peer, and we would likely AppHash anyways.
			OptimisticEngine: true,

			// When we are NOT synced to the tip, process proposal
			// does NOT get called and thus we must ensure that
			// NewPayload is called to get the execution
			// client the payload.
			//
			// When we are synced to the tip, we can skip the
			// NewPayload call since we already gave our execution client
			// the payload in process proposal.
			//
			// In both cases the payload was already accepted by a majority
			// of validators in their process proposal call and thus
			// the "verification aspect" of this NewPayload call is
			// actually irrelevant at this point.
			SkipPayloadVerification: false,

			ProposerAddress: blk.GetProposerAddress(),
			ConsensusTime:   blk.GetConsensusTime(),
		},
		st,
		*blk.GetBeaconBlock(),
	)
	return valUpdates, err
}
