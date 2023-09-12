package keeper

import (
	"fmt"
	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) UnstakeActor(ctx sdk.Context, msg *types.MsgUnstake) error {
	logger := ctx.Logger().With("function", "UnstakeActor")

	// Do some basic parsing and validation on the message
	coinsToUnstake, err := parseCoins(msg.Amount)
	if err != nil {
		logger.Error("Error parsing coins", err.Error())
		return err
	}
	actorAddress := sdk.ValAddress(msg.GetCreator())

	// TODO: Add other actor types by creating a map for actorType->prefix
	switch msg.GetActorType() {
	case types.ServicerPrefix:
		return k.unstakeServicer(ctx, actorAddress, coinsToUnstake)
	default:
		return fmt.Errorf("invalid actor type")
	}

}

func (k Keeper) unstakeServicer(ctx sdk.Context, servicerAddress sdk.ValAddress, coinsToUnstake sdk.Coin) error {
	logger := ctx.Logger().With("servicer", servicerAddress.String())

	store := ctx.KVStore(k.storeKey)
	servicerStore := prefix.NewStore(store, []byte(types.ServicerPrefix))

	byteKey := servicerAddress.Bytes()
	bz := servicerStore.Get(byteKey)

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

	// Serialize the servicer object back to bytes
	bz, err = k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the updated Actor object back to the store
	servicerStore.Set(byteKey, bz)

	// TODO: Add logic to transfer the unstaked coins from the staking pool to the servicer's addres; usually done through the bank module
	return nil
}
