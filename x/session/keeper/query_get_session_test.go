package keeper_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	testkeeper "poktroll/testutil/keeper"
	"poktroll/testutil/mocks"
	apptypes "poktroll/x/application/types"
	svctypes "poktroll/x/servicer/types"
	"poktroll/x/session/types"

	// simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetSessionQuery(t *testing.T) {
	ctrl := gomock.NewController(t)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	accs := simtypes.RandomAccounts(r, 2)
	application := apptypes.Application{
		Address: accs[0].Address.String(),
		Stake:   &sdk.Coin{Denom: "stake", Amount: sdk.NewInt(1)},
	}
	servicer := svctypes.Servicers{
		Address: accs[1].Address.String(),
		Stake:   &sdk.Coin{Denom: "stake", Amount: sdk.NewInt(1)},
	}

	mockApplicationKeeper := mocks.NewMockApplicationKeeper(ctrl)
	mockApplicationKeeper.EXPECT().GetApplication(gomock.Any(), gomock.Any()).Return(application, true).AnyTimes()

	mockServicerKeeper := mocks.NewMockServicerKeeper(ctrl)
	mockServicerKeeper.EXPECT().GetAllServicers(gomock.Any()).Return([]svctypes.Servicers{servicer}).AnyTimes()

	sessionKeeper, ctx := testkeeper.SessionKeeperWithMocks(t, mockApplicationKeeper, mockServicerKeeper)
	wctx := sdk.WrapSDKContext(ctx)

	req := &types.QueryGetSessionRequest{
		AppAddress: application.Address,
		// TODO_TEST: Need to improve the tests by adding ServiceId and BlockHeight
	}

	response, err := sessionKeeper.GetSession(wctx, req)
	require.NoError(t, err)
	fmt.Println(response)
}
