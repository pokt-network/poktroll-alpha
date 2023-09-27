package types

import (
	"testing"

	"poktroll/testutil/sample"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgStakePortal_ValidateBasic(t *testing.T) {
	coins := sdk.NewCoin("stake", sdk.NewInt(int64(1)))
	tests := []struct {
		name string
		msg  MsgStakePortal
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgStakePortal{
				Address: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "missing stake amount",
			msg: MsgStakePortal{
				Address: sample.AccAddress(),
			},
			err: ErrNilStakeAmount,
		},
		{
			name: "valid address and non nil stake amount",
			msg: MsgStakePortal{
				Address:     sample.AccAddress(),
				StakeAmount: &coins,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
