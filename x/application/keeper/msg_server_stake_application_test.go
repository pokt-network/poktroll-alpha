package keeper_test

import (
	"fmt"
	"poktroll/testutil/sample"
	"poktroll/x/application/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgApplicationStake(t *testing.T) {
	ms, ctx := setupMsgServer(t)

	coins := sdk.NewCoin("stake", sdk.NewInt(int64(1)))
	msgStakeApplication := &types.MsgStakeApplication{
		Address:     sample.AccAddress(),
		StakeAmount: &coins,
	}

	res, err := ms.StakeApplication(ctx, msgStakeApplication)
	require.Nil(t, err)
	fmt.Println(res)
}
