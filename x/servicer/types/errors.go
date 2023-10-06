package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/servicer module sentinel errors
var (
	ErrNilStakeAmount               = sdkerrors.Register(ModuleName, 1, "Stake amount is nil")
	ErrEmptyStakeAmount             = sdkerrors.Register(ModuleName, 2, "Stake amount is empty")
	ErrStakeAmountMustBeHigher      = sdkerrors.Register(ModuleName, 3, "The stake amount for a previously staked servicer must be explicitly higher than the prior amount")
	ErrUnstakingNonExistentServicer = sdkerrors.Register(ModuleName, 4, "Could not unstake non-existent servicer")
	ErrNoServicesToStake            = sdkerrors.Register(ModuleName, 5, "Must stake for at least one service")
	ErrActiveFirstSession           = sdkerrors.Register(ModuleName, 6, "First session has not ended yet")
	ErrFutureSessionNumber          = sdkerrors.Register(ModuleName, 7, "Msg session number is in the future")
	ErrClaimSessionNumberNotEnded   = sdkerrors.Register(ModuleName, 8, "Claim session number has not ended")
	ErrProofAndClaimSignerMismatch  = sdkerrors.Register(ModuleName, 9, "Proof and claim signer mismatch")
	ErrSMSTRootHashMismatch         = sdkerrors.Register(ModuleName, 10, "SMST root hash mismatch")
	ErrEarlyProofSubmission         = sdkerrors.Register(ModuleName, 11, "Proof submitted too early")
	ErrLateProofSubmission          = sdkerrors.Register(ModuleName, 12, "Proof submitted too late")
	ErrEarlyClaimSubmission         = sdkerrors.Register(ModuleName, 13, "Claim submitted too early")
	ErrLateClaimSubmission          = sdkerrors.Register(ModuleName, 14, "Claim submitted too late")
	ErrInvalidProof                 = sdkerrors.Register(ModuleName, 15, "Failed to validate proof")
	ErrInvalidPath                  = sdkerrors.Register(ModuleName, 16, "Invalid path")
	ErrUnsupportedMultiSig          = sdkerrors.Register(ModuleName, 17, "Unsupported multi-sig")
)
