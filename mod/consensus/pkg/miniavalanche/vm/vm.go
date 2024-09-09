package vm

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/snow/validators"
	"github.com/ava-labs/avalanchego/version"

	berablock "github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/block"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/miniavalanche/middleware"
)

var (
	_ block.ChainVM = (*VM)(nil)

	// mini-avalanche seems to distinguish from third party libs (e.g. github.com/shirou/gopsutils)
	errNotYetImplemented = errors.New("mini-avalanche: not yet implemented")
)

type VM struct {
	chainCtx *snow.Context

	// middleware interfaces with the bus to send/receive data from the EVM
	middleware middleware.VMMiddleware

	db    database.Database
	state chainState

	validators validators.Manager // exposed to consensus engine

	// in memory list of verified but not yet finalized blocks
	verifiedBlocks map[ids.ID]*StatefulBlock

	preferredBlkID ids.ID
	bb             *blockBuilder
}

func (vm *VM) Initialize(
	ctx context.Context,
	chainCtx *snow.Context,
	db database.Database,
	genesisBytes []byte,
	_ []byte,
	_ []byte,
	toEngine chan<- common.Message,
	_ []*common.Fx,
	_ common.AppSender,
) error {
	vm.chainCtx = chainCtx
	vm.db = db

	// parse genesis to retrieve its components
	genBlk, genVals, ethGen, err := parseGenesis(genesisBytes)
	if err != nil {
		return fmt.Errorf("failed initializing VM: %w", err)
	}

	// init validator set, static for now
	for _, val := range genVals {
		err = vm.validators.AddStaker(chainCtx.SubnetID, val.NodeID, nil, val.id, val.Weight)
		if err != nil {
			return fmt.Errorf("failed registration of validator %v: %w", val.id, err)
		}
	}

	state, err := newState(db, genBlk)
	if err != nil {
		return fmt.Errorf("failed initializing vm state: %w", err)
	}
	vm.state = state

	// initialize block verification stuff
	vm.verifiedBlocks = make(map[ids.ID]*StatefulBlock)

	// initialize block building stuff
	genBlkID := vm.state.GetLastAccepted()
	vm.bb = newBlockBuilder(toEngine, vm)

	// TODO: handle dynamic validator set
	// At this stage of hooking stuff up, we consider a static validators set
	// where validators (and especially the mapping validator -> NodeID) is
	// setup in Genesis. We don't even check data correspondence and assume
	// genesis is well formatted
	_, err = vm.middleware.InitGenesis(ctx, ethGen)
	if err != nil {
		return fmt.Errorf("failed initializing genesis in middleware: %w", err)
	}

	vm.preferredBlkID = genBlkID
	vm.bb = newBlockBuilder(toEngine, vm)
	return nil
}

func (vm *VM) SetState(_ context.Context, state snow.State) error {
	if state == snow.NormalOp {
		// NormalOp signals that both state-sync and bootstrapping are done.
		// Consensus is in sync with the network, so VM can start building blocks.
		vm.bb.Start()
	}
	return nil
}

func (vm *VM) Shutdown(context.Context) error {
	if vm.state == nil {
		// Shutdown may be called before VM in initialized
		// Nothing to do in this case
		return nil
	}

	vm.bb.Shutdown()
	return errors.Join(
		vm.state.Close(),
		vm.db.Close(),
	)
}

func (vm *VM) Version(context.Context) (string, error) {
	return vmVersion.String(), nil
}

func (vm *VM) CreateHandlers(context.Context) (map[string]http.Handler, error) {
	return nil, fmt.Errorf("createHandler: %w", errNotYetImplemented)
}

func (vm *VM) HealthCheck(context.Context) (interface{}, error) {
	return nil, fmt.Errorf("healthCheck: %w", errNotYetImplemented)
}

func (vm *VM) Connected(_ context.Context, _ ids.NodeID, _ *version.Application) error {
	return nil
}

func (vm *VM) Disconnected(_ context.Context, _ ids.NodeID) error {
	return nil
}

func (vm *VM) GetBlock(_ context.Context, blkID ids.ID) (snowman.Block, error) {
	return vm.getBlock(blkID)
}

func (vm *VM) getBlock(blkID ids.ID) (*StatefulBlock, error) {
	fullBlk, found := vm.verifiedBlocks[blkID]
	if found {
		return fullBlk, nil
	}

	switch blk, err := vm.state.GetBlock(blkID); err {
	case nil:
		return &StatefulBlock{
			StatelessBlock: blk,
			vm:             vm,
		}, nil
	case database.ErrNotFound:
		return nil, database.ErrNotFound
	default:
		return nil, fmt.Errorf("failed retrieving block %s: %w", blkID, err)
	}
}

func (vm *VM) ParseBlock(_ context.Context, blockBytes []byte) (snowman.Block, error) {
	blk, err := berablock.ParseStatelessBlock(blockBytes)
	if err != nil {
		return nil, err
	}

	return &StatefulBlock{
		StatelessBlock: blk,
		vm:             vm,
	}, nil
}

func (vm *VM) BuildBlock(ctx context.Context) (snowman.Block, error) {
	return vm.bb.BuildBlock(ctx)
}

func (vm *VM) SetPreference(_ context.Context, blkID ids.ID) error {
	vm.preferredBlkID = blkID
	return nil
}

func (vm *VM) LastAccepted(context.Context) (ids.ID, error) {
	return vm.state.GetLastAccepted(), nil
}

func (vm *VM) GetBlockIDAtHeight(_ context.Context, h uint64) (ids.ID, error) {
	switch blkID, err := vm.state.GetBlockID(h); err {
	case nil:
		return blkID, err
	case database.ErrNotFound:
		return ids.Empty, database.ErrNotFound
	default:
		return ids.Empty, fmt.Errorf("failed retrieving block ID at height %d: %w", h, err)
	}
}