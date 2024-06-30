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
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/state/deneb"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
	"github.com/berachain/beacon-kit/mod/storage/pkg/sszdb"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput struct {
	depinject.In
	AvailabilityStore *AvailabilityStore
	ChainSpec         common.ChainSpec
	DepositStore      *DepositStore
	KVStore           *KVStore
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend(
	in StorageBackendInput,
) StorageBackend {
	return storage.NewBackend[
		*AvailabilityStore,
		*BeaconBlockBody,
		BeaconState,
		*BeaconStateMarshallable,
		*DepositStore,
	](
		in.ChainSpec,
		in.AvailabilityStore,
		in.KVStore,
		in.DepositStore,
	)
}

// KVStoreInput is the input for the ProvideKVStore function.
type KVStoreInput struct {
	depinject.In
	Environment appmodule.Environment
	AppOpts     servertypes.AppOptions
}

// ProvideKVStore is the depinject provider that returns a beacon KV store.
func ProvideKVStore(
	in KVStoreInput,
) (*KVStore, error) {
	cfg := sszdb.BackendConfig{
		Path: cast.ToString(in.AppOpts.Get(flags.FlagHome)) + "/data/sszdb.db",
	}
	backend, err := sszdb.NewBackend(cfg)
	if err != nil {
		return nil, err
	}
	stateObject := &deneb.BeaconState{}
	stateObject.EmptyState()
	szdb, err := sszdb.NewSchemaDb[*ExecutionPayloadHeader](backend, stateObject)
	if err != nil {
		return nil, err
	}

	payloadCodec := &encoding.
		SSZInterfaceCodec[*ExecutionPayloadHeader]{}
	return beacondb.New[
		*BeaconBlockHeader,
		*types.Eth1Data,
		*ExecutionPayloadHeader,
		*types.Fork,
		*types.Validator,
	](szdb, in.Environment.KVStoreService, payloadCodec), nil
}
