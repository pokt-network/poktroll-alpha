package keeper

import (
	"context"
	"fmt"
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"poktroll/x/application/types"
)

func (k msgServer) DelegateToPortal(goCtx context.Context, msg *types.MsgDelegateToPortal) (*types.MsgDelegateToPortalResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx).With("method", "DelegateToPortal")
	if _, err := anyPkToSecp256k1(msg.PortalPubKey); err != nil {
		return nil, fmt.Errorf("error extracting secp256k1 public key from any type: %w", err)
	}
	logger.Info(fmt.Sprintf("About to delegate application %v to %+v", msg.Address, msg.PortalPubKey.Value))

	// Update the store
	if err := k.DelegatePortal(ctx, msg.Address, *msg.PortalPubKey); err != nil {
		logger.Error(fmt.Sprintf("could not update store with delegated portal for application: %v", msg.Address))
		return nil, err
	}
	logger.Info(fmt.Sprintf("Successfully delegated application %v to %+v", msg.Address, msg.PortalPubKey.Value))

	return &types.MsgDelegateToPortalResponse{}, nil
}

func anyPkToSecp256k1(pk *types1.Any) (*secp256k1.PubKey, error) {
	pubBz := pk.GetValue()
	key := new(secp256k1.PubKey)
	if err := key.Unmarshal(pubBz); err != nil {
		return nil, err
	}
	/**
	key, ok := pk.GetCachedValue().(*secp256k1.PubKey)
	if !ok {
		return nil, fmt.Errorf("cached value is not secp256k1 public key")
	}
	**/
	return key, nil
}
