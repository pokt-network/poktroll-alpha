package keeper

import (
	"fmt"
	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) UnstakeActor(ctx sdk.Context, msg *types.MsgUnstake) error {
	logger := ctx.Logger().With("function", "UnstakeActor")

	// TODO: Add other actor types by creating a map for actorType->prefix
	var storePrefix string
	switch msg.GetActorType() {
	case types.ServicerPrefix:
		storePrefix = types.ServicerPrefix
	default:
		return fmt.Errorf("invalid actor type")
	}

	// Convert msg.Amount string to a Coin value
	coinsToUnstake, err := parseCoins(msg.Amount)
	if err != nil {
		return err
	}

	// Convert the valAddr to bytes as it will be the key to store validator info
	addr := sdk.ValAddress(msg.GetCreator())
	byteKey := addr.Bytes()
	logger = logger.With("actor", addr.String())

	// Retrieve the store for the actor being unstake
	store := ctx.KVStore(k.storeKey)
	actorStore := prefix.NewStore(store, []byte(storePrefix))

	// Get existing validator data from the store
	bz := actorStore.Get(byteKey)

	if bz == nil {
		logger.Info("servicer not found")
		return fmt.Errorf("servicer not found")
	}

	// Deserialize the byte array into a Validator object
	var servicer types.Servicer
	k.cdc.Unmarshal(bz, &servicer)

	coinsStaked, err := parseCoins(servicer.GetStakeInfo().GetCoinsStaked())
	if err != nil {
		return err
	}
	// Update staking amount
	newStakeAmount := coinsStaked.Sub(coinsToUnstake)
	servicer.GetStakeInfo().CoinsStaked = newStakeAmount.String()
	// TODO: Add staked amount checks here when trying to unstake more than currently staked

	// Serialize the Validator object back to bytes
	bz, err = k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the updated Actor object back to the store
	actorStore.Set(byteKey, bz)

	// TODO: Add logic to transfer the unstaked coins from the staking pool to the servicer's addres; usually done through the bank module

	return nil
}
