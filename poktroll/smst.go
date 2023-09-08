package poktroll

import (
	"crypto/sha256"

	"github.com/pokt-network/smt"
)

var SMST *smt.SMST

func init() {
	db, err := smt.NewKVStore("")
	if err != nil {
		panic(err)
	}
	defer db.Stop()
	SMST = smt.NewSparseMerkleSumTree(db, sha256.New())
}
