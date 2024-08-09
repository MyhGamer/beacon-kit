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
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dablob "github.com/berachain/beacon-kit/mod/da/pkg/blob"
	"github.com/berachain/beacon-kit/mod/da/pkg/da"
	"github.com/berachain/beacon-kit/mod/da/pkg/kzg"
	"github.com/berachain/beacon-kit/mod/depinject"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app/components/metrics"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	gokzg4844 "github.com/crate-crypto/go-kzg-4844"
)

// BlobProofVerifierInput is the input for the
// dep inject framework.
type BlobProofVerifierInput struct {
	depinject.In
	Config           *config.Config
	JSONTrustedSetup *gokzg4844.JSONTrustedSetup
}

// ProvideBlobProofVerifier is a function that provides the module to the
// application.
func ProvideBlobProofVerifier(
	in BlobProofVerifierInput,
) (kzg.BlobProofVerifier, error) {
	return kzg.NewBlobProofVerifier(
		in.Config.KZG.Implementation,
		in.JSONTrustedSetup,
	)
}

// BlobVerifierInput is the input for the BlobVerifier.
type BlobVerifierInput struct {
	depinject.In
	BlobProofVerifier kzg.BlobProofVerifier
	TelemetrySink     *metrics.TelemetrySink
}

// ProvideBlobVerifier is a function that provides the BlobVerifier to the
// depinject framework.
func ProvideBlobVerifier(in BlobVerifierInput) *BlobVerifier {
	return dablob.NewVerifier(in.BlobProofVerifier, in.TelemetrySink)
}

// BlobProcessorIn is the input for the BlobProcessor.
type BlobProcessorIn struct {
	depinject.In

	BlobVerifier  *BlobVerifier
	ChainSpec     common.ChainSpec
	Logger        *Logger
	TelemetrySink *metrics.TelemetrySink
}

// ProvideBlobProcessor is a function that provides the BlobProcessor to the
// depinject framework.
func ProvideBlobProcessor(in BlobProcessorIn) *BlobProcessor {
	return dablob.NewProcessor[
		*AvailabilityStore,
		*BeaconBlockBody,
	](
		in.Logger.With("service", "blob-processor"),
		in.ChainSpec,
		in.BlobVerifier,
		types.BlockBodyKZGOffset,
		in.TelemetrySink,
	)
}

// DAServiceIn is the input for the BlobService.
type DAServiceIn struct {
	depinject.In

	AvailabilityStore *AvailabilityStore
	SidecarsBroker    *SidecarsBroker
	BlobProcessor     *BlobProcessor
	Logger            *Logger
}

// ProvideDAService is a function that provides the BlobService to the
// depinject framework.
func ProvideDAService(in DAServiceIn) *DAService {
	return da.NewService[
		*AvailabilityStore,
		*BeaconBlockBody,
		*BlobSidecars,
		*SidecarsBroker,
		*ExecutionPayload,
	](
		in.AvailabilityStore,
		in.BlobProcessor,
		in.SidecarsBroker,
		in.Logger.With("service", "da"),
	)
}