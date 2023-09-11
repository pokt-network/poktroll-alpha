package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"poktroll/x/poktroll/types"
)

func (k Keeper) StakeActor(ctx sdk.Context, msg *types.MsgStake) error {
	logger := ctx.Logger().With("function", "StakeActor")

	// TODO: Add other actor types by creating a map for actorType->prefix
	var storePrefix string
	switch msg.GetActorType() {
	case types.ServicerPrefix:
		storePrefix = types.ServicerPrefix
	default:
		return fmt.Errorf("invalid actor type")
	}

	// Convert msg.Amount string to a Coin value
	coinsToStake, err := parseCoins(msg.Amount)
	if err != nil {
		return err
	}

	// Convert the valAddr to bytes as it will be the key to store validator info
	addr := sdk.ValAddress(msg.GetCreator())
	byteKey := addr.Bytes()
	logger = logger.With("actor", addr.String())

	// Update store with new staking amount
	store := ctx.KVStore(k.storeKey)
	actorStore := prefix.NewStore(store, []byte(storePrefix))

	// Get existing validator data from store
	bz := actorStore.Get(byteKey)

	// Support only servicer staking for now.
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
				Address:     addr.String(),
				CoinsStaked: coinsToStake.String(),
			},
		}
	}

	// Serialize the Servicer object back to bytes
	bz, err = k.cdc.Marshal(&servicer)
	if err != nil {
		return err
	}

	// Save the updated Actor object back to the store
	actorStore.Set(byteKey, bz)

	// TODO: sends coins to the staking module's pool!

	return nil
}
