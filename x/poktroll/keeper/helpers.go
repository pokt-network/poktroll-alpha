package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func parseCoins(coins string) (sdk.Coin, error) {
	amt, err := sdk.ParseCoinNormalized(coins)
	if err != nil {
		return sdk.Coin{}, err
	}
	return amt, nil
}
