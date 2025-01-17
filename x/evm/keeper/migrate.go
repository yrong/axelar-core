package keeper

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axelarnetwork/axelar-core/x/evm/types"
	"github.com/axelarnetwork/utils/slices"
)

// GetMigrationHandler returns the handler that performs in-place store migrations from v0.24 to v0.25
// The migration includes:
// - migrate contracts bytecode (CRUCIAL AND DO NOT DELETE) for all evm chains
// - set TransferLimit parameter
func GetMigrationHandler(k BaseKeeper, n types.Nexus, m types.MultisigKeeper) func(ctx sdk.Context) error {
	return func(ctx sdk.Context) error {
		// migrate contracts bytecode (CRUCIAL AND DO NOT DELETE) for all evm chains
		for _, chain := range slices.Filter(n.GetChains(ctx), types.IsEVMChain) {
			ck := k.ForChain(chain.Name).(chainKeeper)
			if err := migrateContractsBytecode(ctx, ck); err != nil {
				return err
			}
		}

		// set TransferLimit param
		for _, chain := range slices.Filter(n.GetChains(ctx), types.IsEVMChain) {
			ck := k.ForChain(chain.Name).(chainKeeper)
			if err := addTransferLimitParam(ctx, ck); err != nil {
				return err
			}
		}

		return nil
	}
}

func addTransferLimitParam(ctx sdk.Context, ck chainKeeper) error {
	subspace, ok := ck.getSubspace(ctx)
	if !ok {
		return fmt.Errorf("param subspace for chain %s should exist", ck.GetName())
	}

	subspace.Set(ctx, types.KeyTransferLimit, types.DefaultParams()[0].TransferLimit)

	return nil
}

// this function migrates the contracts bytecode to the latest for every existing
// EVM chain. It's crucial whenever contracts are changed between versions and
// DO NOT DELETE
func migrateContractsBytecode(ctx sdk.Context, ck chainKeeper) error {
	bzToken, err := hex.DecodeString(types.Token)
	if err != nil {
		return err
	}

	bzBurnable, err := hex.DecodeString(types.Burnable)
	if err != nil {
		return err
	}

	subspace, ok := ck.getSubspace(ctx)
	if !ok {
		return fmt.Errorf("param subspace for chain %s should exist", ck.GetName())
	}

	subspace.Set(ctx, types.KeyToken, bzToken)
	subspace.Set(ctx, types.KeyBurnable, bzBurnable)

	return nil
}
