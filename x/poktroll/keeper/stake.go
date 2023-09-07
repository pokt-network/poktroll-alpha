package keeper

import (
	"fmt"
	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StakeActor(ctx sdk.Context, msg *types.MsgStake) error {
	// Find which actor we're dealing with
	var storePrefix string

	switch msg.GetActorType() {
	case types.ServicerPrefix: // TODO: Add other actor types here
		storePrefix = types.ServicerPrefix
	// case types.WatcherPrefix, types.PortalPrefix, types.ApplicationPrefix:
	// 	// Set store prefix to the actor type
	// 	storePrefix = msg.GetActorType()
	default:
		return fmt.Errorf("invalid actor type")
	}

	// TODO: sends coins to the staking module's pool!

	// Update store with new staking amount
	store := ctx.KVStore(k.storeKey)
	actorStore := prefix.NewStore(store, []byte(storePrefix))

	// Convert msg.Amount string to a Coin value
	amtToStake, err := parseCoins(msg.Amount)
	if err != nil {
		return err
	}

	// Convert the valAddr to bytes as it will be the key to store validator info
	addr := sdk.ValAddress(msg.GetCreator())
	byteKey := addr.Bytes()

	// Get existing validator data from store
	bz := actorStore.Get(byteKey)

	// Support only servicer staking for now.
	var servicer types.Servicer
	if bz != nil {
		// Deserialize the byte array into a Validator object
		k.cdc.Unmarshal(bz, &servicer)
		// Parse staker coins
		amtStaked, err := parseCoins(servicer.GetStakeInfo().GetCoinsStaked())
		if err != nil {
			return err
		}
		// Update staking amount
		amtStaked = amtStaked.Add(amtToStake)
		servicer.GetStakeInfo().CoinsStaked = amtStaked.String()
	} else {
		// Create a new Servicer object if not found
		servicer = types.Servicer{
			StakeInfo: &types.StakeInfo{
				Address:     addr.String(),
				CoinsStaked: amtToStake.String(),
			},
		}
	}

	// Serialize the Validator object back to bytes
	bz, err = k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the updated Validator object back to the store
	actorStore.Set(byteKey, bz)

	return nil
}

func parseCoins(coins string) (sdk.Coin, error) {
	amt, err := sdk.ParseCoinNormalized(coins)
	if err != nil {
		return sdk.Coin{}, err
	}
	return amt, nil
}

func (k Keeper) UnstakeActor(ctx sdk.Context, msg *types.MsgUnstake) error {
	// Find which actor we're dealing with
	var storePrefix string

	switch msg.GetActorType() {
	case types.ServicerPrefix: // TODO: Add other actor types here
		storePrefix = types.ServicerPrefix
	// case types.WatcherPrefix, types.PortalPrefix, types.ApplicationPrefix:
	// 	// Set store prefix to the actor type
	// 	storePrefix = msg.GetActorType()
	default:
		return fmt.Errorf("invalid actor type")
	}

	store := ctx.KVStore(k.storeKey)
	actorStore := prefix.NewStore(store, []byte(storePrefix))

	// Convert msg.Amount string to a Coin value
	amtToUnstake, err := parseCoins(msg.Amount)
	if err != nil {
		return err
	}

	// Convert the valAddr to bytes as it will be the key to store validator info
	addr := sdk.ValAddress(msg.GetCreator())
	byteKey := addr.Bytes()

	// Get existing validator data from the store
	bz := actorStore.Get(byteKey)

	logger := ctx.Logger()
	if bz == nil {
		logger.Info("servicer not found")
		return fmt.Errorf("servicer not found")
	}

	var servicer types.Servicer

	// Deserialize the byte array into a Validator object
	k.cdc.Unmarshal(bz, &servicer)

	amtStaked, err := parseCoins(servicer.GetStakeInfo().GetCoinsStaked())
	if err != nil {
		return err
	}
	// Update staking amount
	newStakeAmount := amtStaked.Sub(amtToUnstake)
	servicer.GetStakeInfo().CoinsStaked = newStakeAmount.String()
	// TODO: Add staked amount checks here
	// if err != nil {
	// 	// Error: trying to unstake more than currently staked
	// 	return fmt.Errorf("insufficient staking amount to unstake")
	// }

	// Serialize the Validator object back to bytes
	bz, err = k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the updated Validator object back to the store
	actorStore.Set(byteKey, bz)

	// TODO: Add logic to transfer the unstaked coins from the module's account to the validator's account
	// (This is usually done through the bank module)

	return nil

}
