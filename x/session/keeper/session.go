package keeper

import (
	"crypto"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "golang.org/x/crypto/sha3"

	apptypes "poktroll/x/application/types"
	srvstypes "poktroll/x/service/types"
	svctypes "poktroll/x/servicer/types"
	"poktroll/x/session/types"
	sharedtypes "poktroll/x/shared/types"
)

const (
	// TODO_REFACTOR: Move these constants into governance parameters
	NumSessionBlocks       uint64 = 4
	numServicersPerSession uint64 = 25
)

var (
	// TODO_REFACTOR: Move these constants into a shared repo
	SHA3HashLen = crypto.SHA3_256.Size()
)

func (k Keeper) GetSessionForApp(ctx sdk.Context, appAddress string, serviceId string, blockHeight uint64) (*sharedtypes.Session, error) {
	logger := k.Logger(ctx).With("module", types.ModuleName).With("method", "GetSessionForApp")
	logger.Info(fmt.Sprintf("About to get session for app address %s", appAddress))

	service := srvstypes.ServiceId{Id: serviceId}

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

	// filter servicers only if there is an overlap between the services the app & servicers both staked for
	servicerPointers = findMatchingServicers(app, servicerPointers, &srvstypes.ServiceId{Id: serviceId})

	session := sharedtypes.Session{
		// NB: These parameters are hydrated by the hydrator below
		// SessionId:,
		// SessionNumber:,
		// SessionBlockStartHeight:,
		// NumBlocksPerSession:,
		Service:     &service,
		Application: &app,
		Servicers:   servicerPointers,
	}

	// TODO_CLEANUP: Look at `utility/session.go` in the v1 repo for a cleaner hydration/implementation.
	sessionHydrator := &sessionHydrator{
		session:     &session,
		blockHeight: blockHeight,
	}
	if err := sessionHydrator.hydrateSessionMetadata(); err != nil {
		return nil, fmt.Errorf("failed to hydrate session Metadata: %w", err)
	}
	if err := sessionHydrator.hydrateSessionID(); err != nil {
		return nil, fmt.Errorf("failed to hydrate session ID: %w", err)
	}

	return &session, nil

}

// TODO_IMPLEMENT: Need to pseudo-randomly select only the relevant servicers
func findMatchingServicers(app apptypes.Application, servicers []*svctypes.Servicers, targetServiceId *srvstypes.ServiceId) []*svctypes.Servicers {
	matchingServicers := []*svctypes.Servicers{}

	serviceIDMap := make(map[string]struct{})
	for _, service := range app.Services {
		serviceIDMap[service.Id] = struct{}{}
	}

	for _, servicer := range servicers {
		for _, service := range servicer.Services {
			if service.Id != nil && targetServiceId.Id == service.Id.Id {
				if _, exists := serviceIDMap[service.Id.Id]; exists {
					matchingServicers = append(matchingServicers, servicer)
					break
				}
			}
		}
	}

	return matchingServicers
}

type sessionHydrator struct {
	// The session being hydrated and returned
	session *sharedtypes.Session

	// The height at which the request is being made to get session information
	blockHeight uint64

	// A redundant helper that maintains a hex decoded copy of `session.Id` used for session hydration
	sessionIdBz []byte
}

// hydrateSessionMetadata hydrates the height at which the session started, its number, and the number of blocks per session
func (s *sessionHydrator) hydrateSessionMetadata() error {
	numBlocksPerSession := NumSessionBlocks // TODO: Get from governance params
	numBlocksAheadOfSession := s.blockHeight % numBlocksPerSession
	s.session.NumBlocksPerSession = numBlocksPerSession
	s.session.SessionNumber = s.blockHeight / numBlocksPerSession
	s.session.SessionBlockStartHeight = s.blockHeight - numBlocksAheadOfSession
	return nil
}

// hydrateSessionID use both session and on-chain data to determine a unique session ID
func (s *sessionHydrator) hydrateSessionID() error {
	sessionHeightBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(sessionHeightBz, s.session.SessionBlockStartHeight)

	// TODO: Retrieve this from the blockchain
	prevBlockHash := hex.EncodeToString([]byte("prevHash"))
	prevHashBz, err := hex.DecodeString(prevBlockHash)
	if err != nil {
		return err
	}

	appPubKeyBz := []byte(s.session.Application.Address) // TODO: Use public key instead of address
	serviceIdBz := []byte(string(s.session.Service.Id))

	s.sessionIdBz = concat(sessionHeightBz, prevHashBz, serviceIdBz, appPubKeyBz)
	s.session.SessionId = GetHashStringFromBytes(s.sessionIdBz)

	return nil
}

func concat(b ...[]byte) (result []byte) {
	for _, bz := range b {
		result = append(result, bz...)
	}
	return result
}

// TODO_REFACTOR: Move into a shared library
func SHA3Hash(bz []byte) []byte {
	hasher := crypto.SHA3_256.New()
	hasher.Write(bz)
	return hasher.Sum(nil)
}

// GetHashStringFromBytes returns hex(SHA3Hash(bytesArgument)); typically used to compute a TransactionHash
func GetHashStringFromBytes(bytes []byte) string {
	return hex.EncodeToString(SHA3Hash(bytes))
}
