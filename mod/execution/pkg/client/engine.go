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

package client

import (
	"context"
	"fmt"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/errors"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/execution/pkg/client/ethclient"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
)

// NewPayload calls the engine_newPayloadVX method via JSON-RPC.
func (s *EngineClient[ExecutionPayloadT]) NewPayload(
	ctx context.Context,
	payload ExecutionPayloadT,
	versionedHashes []common.ExecutionHash,
	parentBeaconBlockRoot *primitives.Root,
) (*common.ExecutionHash, error) {
	var (
		startTime    = time.Now()
		dctx, cancel = context.WithTimeoutCause(
			ctx, s.cfg.RPCTimeout, engineerrors.ErrEngineAPITimeout,
		)
	)
	defer s.metrics.measureNewPayloadDuration(startTime)
	defer cancel()

	// Call the appropriate RPC method based on the payload version.
	result, err := s.Eth1Client.NewPayload(
		dctx,
		payload,
		versionedHashes,
		parentBeaconBlockRoot,
	)
	if err != nil {
		if errors.Is(err, engineerrors.ErrEngineAPITimeout) {
			s.metrics.incrementNewPayloadTimeout()
		}
		return nil, err
	} else if result == nil {
		return nil, engineerrors.ErrNilPayloadStatus
	}

	// This case is only true when the payload is invalid, so
	// `processPayloadStatusResult` below will return an error.
	if validationErr := result.ValidationError; validationErr != nil {
		s.logger.Error(
			"Got a validation error in newPayload",
			"err",
			errors.New(*validationErr),
		)
	}

	return processPayloadStatusResult(result)
}

// ForkchoiceUpdated calls the engine_forkchoiceUpdatedV1 method via JSON-RPC.
func (s *EngineClient[ExecutionPayloadT]) ForkchoiceUpdated(
	ctx context.Context,
	state *engineprimitives.ForkchoiceStateV1,
	attrs engineprimitives.PayloadAttributer,
	forkVersion uint32,
) (*engineprimitives.PayloadID, *common.ExecutionHash, error) {
	var (
		startTime    = time.Now()
		dctx, cancel = context.WithTimeoutCause(
			ctx, s.cfg.RPCTimeout, engineerrors.ErrEngineAPITimeout,
		)
	)
	defer s.metrics.measureForkchoiceUpdateDuration(startTime)
	defer cancel()

	if attrs == nil || attrs.IsNil() {
		return nil, nil, engineerrors.ErrNilPayloadAttributes
	}
	// If the suggested fee recipient is not set, log a warning.
	if attrs.GetSuggestedFeeRecipient() == (common.ZeroAddress) {
		s.logger.Warn(
			"suggested fee recipient is not configured 🔆",
			"fee-recipent", common.DisplayBytes(
				common.ZeroAddress[:]).TerminalString(),
		)
	}

	result, err := s.Eth1Client.ForkchoiceUpdated(dctx, state, attrs, forkVersion)

	if err != nil {
		if errors.Is(err, engineerrors.ErrEngineAPITimeout) {
			s.metrics.incrementForkchoiceUpdateTimeout()
		}
		return nil, nil, s.handleRPCError(err)
	} else if result == nil {
		return nil, nil, engineerrors.ErrNilForkchoiceResponse
	}

	latestValidHash, err := processPayloadStatusResult((&result.PayloadStatus))
	if err != nil {
		return nil, latestValidHash, err
	}
	return result.PayloadID, latestValidHash, nil
}

// GetPayload retrieves the execution data and blobs bundle using the
// engine_getPayloadVX method via JSON-RPC.
func (s *EngineClient[ExecutionPayloadT]) GetPayload(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion uint32,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	var (
		startTime    = time.Now()
		dctx, cancel = context.WithTimeoutCause(
			ctx,
			s.cfg.RPCTimeout,
			engineerrors.ErrEngineAPITimeout,
		)
	)
	defer s.metrics.measureGetPayloadDuration(startTime)
	defer cancel()

	fmt.Println("GETPAYLOAD")

	result, err := s.Eth1Client.GetPayload(dctx, payloadID, forkVersion)
	if err != nil {
		if errors.Is(err, engineerrors.ErrEngineAPITimeout) {
			s.metrics.incrementGetPayloadTimeout()
		}
		return result, s.handleRPCError(err)
	}

	fmt.Println("POST GETPAYLOAD")

	if result == nil {
		return result, engineerrors.ErrNilExecutionPayloadEnvelope
	}

	if result.GetBlobsBundle() == nil && forkVersion >= version.Deneb {
		return result, engineerrors.ErrNilBlobsBundle
	}

	return result, nil
}

// ExchangeCapabilities calls the engine_exchangeCapabilities method via
// JSON-RPC.
func (s *EngineClient[ExecutionPayloadT]) ExchangeCapabilities(
	ctx context.Context,
) ([]string, error) {
	result, err := s.Eth1Client.ExchangeCapabilities(
		ctx, ethclient.BeaconKitSupportedCapabilities(),
	)
	if err != nil {
		return nil, err
	}

	// Capture and log the capabilities that the execution client has.
	for _, capability := range result {
		s.logger.Info("exchanged capability", "capability", capability)
		s.capabilities[capability] = struct{}{}
	}

	// Log the capabilities that the execution client does not have.
	for _, capability := range ethclient.BeaconKitSupportedCapabilities() {
		if _, exists := s.capabilities[capability]; !exists {
			s.logger.Warn(
				"your execution client may require an update 🚸",
				"unsupported_capability", capability,
			)
		}
	}

	return result, nil
}
