package poktroll

import (
	"github.com/pokt-network/smt"
)

var SMST *smt.SMST

func init() {
	db := smt.NewKVStore("")
	defer db.Stop()
	SMST = smt.NewSparseMerkleSumTree(db, sha256.New())
}
