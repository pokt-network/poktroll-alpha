package keeper

import (
	"fmt"
	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StakeActor(ctx sdk.Context, msg *types.MsgStake) error {
	logger := ctx.Logger().With("function", "StakeActor")

	// Do some basic parsing and validation on the message
	coinsToStake, err := parseCoins(msg.Amount)
	if err != nil {
		logger.Error("Error parsing coins", err.Error())
		return err
	}
	actorAddress := sdk.ValAddress(msg.GetCreator())

	// TODO: Add other actor types by creating a map for actorType->prefix
	switch msg.GetActorType() {
	case types.ServicerPrefix:
		return k.stakeServicer(ctx, actorAddress, coinsToStake)
	default:
		return fmt.Errorf("invalid actor type")
	}

	// TODO: sends coins to the staking module's pool!

}

func (k Keeper) stakeServicer(ctx sdk.Context, servicerAddress sdk.ValAddress, coinsToStake sdk.Coin) error {
	logger := ctx.Logger().With("servicer", servicerAddress.String())

	store := ctx.KVStore(k.storeKey)
	servicerStore := prefix.NewStore(store, []byte(types.ServicerPrefix))

	byteKey := servicerAddress.Bytes()
	bz := servicerStore.Get(byteKey)

	var servicer types.Servicer
	if bz != nil {
		k.cdc.Unmarshal(bz, &servicer)
		coinsStaked, err := parseCoins(servicer.GetStakeInfo().GetCoinsStaked())
		if err != nil {
			return err
		}
		// Update staking amount
		coinsStaked = coinsStaked.Add(coinsToStake)
		servicer.GetStakeInfo().CoinsStaked = coinsStaked.String()
	} else {
		// Create a new Servicer object if not found
		servicer = types.Servicer{
			StakeInfo: &types.StakeInfo{
				Address:     servicerAddress.String(),
				CoinsStaked: coinsToStake.String(),
			},
		}
		logger.Info(fmt.Sprintf("Registered new servicer %s", servicerAddress.String()))
	}

	// Serialize the Servicer object back to bytes
	bz, err := k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the Servicer object bytes to the store
	servicerStore.Set(byteKey, bz)
	return nil
}
