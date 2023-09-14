package keeper

import (
	"fmt"
	svctypes "poktroll/x/servicer/types"
	"poktroll/x/session/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetSessionForApp(ctx sdk.Context, appAddress string) (*types.Session, error) {
	logger := k.Logger(ctx).With("module", types.ModuleName).With("method", "GetSessionForApp")
	logger.Info(fmt.Sprintf("About to get session for app address %s", appAddress))

	app, found := k.appKeeper.GetApplication(ctx, appAddress)
	if !found {
		logger.Error(fmt.Sprintf("App not found for address %s", appAddress))
		return nil, types.ErrFindApp
	}
	logger.Info(fmt.Sprintf("App found for address %s: %v", appAddress, app))

	servicers := k.svcKeeper.GetAllServicers(ctx)
	if len(servicers) == 0 {
		logger.Error("Error retrieving servicers: none found")
		return nil, types.ErrNoServicersFound
	}
	logger.Info(fmt.Sprintf("Servicers found: %v", servicers))

	// INVESTIGATE: The `Session` protobuf expects pointers but the `GetAllServicers` keep methods returns values. Look into cosmos to figure out the best path here.
	servicerPointers := make([]*svctypes.Servicers, len(servicers))
	for i, servicer := range servicers {
		servicerPointers[i] = &servicer
	}

	session := types.Session{
		Application: &app,
		Servicers:   servicerPointers,
	}
	return &session, nil

}
