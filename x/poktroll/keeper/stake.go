package keeper

import (
	"fmt"
	"poktroll/x/poktroll/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) StakeActor(ctx sdk.Context, msg *types.MsgStake) error {
	// TODO: sends coins to the staking module's pool!

	// Update store with new staking amount
	store := ctx.KVStore(k.storeKey)
	stakerStore := prefix.NewStore(store, []byte(types.StakerPrefix))

	// Convert msg.Amount string to a Coin value
	amtToStake, err := parseCoins(msg.Amount)
	if err != nil {
		return err
	}

	// Convert the valAddr to bytes as it will be the key to store validator info
	addr := sdk.ValAddress(msg.GetCreator())
	byteKey := addr.Bytes()

	// Get existing validator data from store
	bz := stakerStore.Get(byteKey)

	var staker types.Staker
	if bz != nil {
		// Deserialize the byte array into a Validator object
		k.cdc.Unmarshal(bz, &staker)
		// Parse staker coins
		amtStaked, err := parseCoins(staker.CoinsStaked)
		if err != nil {
			return err
		}
		// Update staking amount
		amtStaked = amtStaked.Add(amtToStake)
		staker.CoinsStaked = amtStaked.String()
	} else {
		// Create a new Staker object if not found
		staker = types.Staker{
			Addr:        addr.String(),
			CoinsStaked: amtToStake.String(),
		}
	}

	// Serialize the Validator object back to bytes
	bz, err = k.cdc.Marshal(&staker)
	if err != nil {
		return err
	}

	// Save the updated Validator object back to the store
	stakerStore.Set(byteKey, bz)

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

	store := ctx.KVStore(k.storeKey)
	servStore := prefix.NewStore(store, []byte(types.StakerPrefix))

	// Convert msg.Amount string to a Coin value
	amtToUnstake, err := parseCoins(msg.Amount)
	if err != nil {
		return err
	}

	// Convert the valAddr to bytes as it will be the key to store validator info
	addr := sdk.ValAddress(msg.GetCreator())
	byteKey := addr.Bytes()

	// Get existing validator data from the store
	bz := servStore.Get(byteKey)

	logger := ctx.Logger()
	if bz == nil {
		logger.Info("staker not found")
		return fmt.Errorf("staker not found")
	}

	var staker types.Staker

	// Deserialize the byte array into a Validator object
	k.cdc.Unmarshal(bz, &staker)

	amtStaked, err := parseCoins(staker.CoinsStaked)
	if err != nil {
		return err
	}
	// Update staking amount
	newStakeAmount := amtStaked.Sub(amtToUnstake)
	staker.CoinsStaked = newStakeAmount.String()
	// TODO: Add staked amount checks here
	// if err != nil {
	// 	// Error: trying to unstake more than currently staked
	// 	return fmt.Errorf("insufficient staking amount to unstake")
	// }

	// Serialize the Validator object back to bytes
	bz, err = k.cdc.Marshal(&staker)
	if err != nil {
		return err
	}

	// Save the updated Validator object back to the store
	servStore.Set(byteKey, bz)

	// TODO: Add logic to transfer the unstaked coins from the module's account to the validator's account
	// (This is usually done through the bank module)

	return nil

}
