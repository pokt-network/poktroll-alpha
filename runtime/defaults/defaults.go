package defaults

import (
	"os"
)

func init() {
	initDefaultRootDirectory()
}

func initDefaultRootDirectory() {
	// use home directory + /.cmt-pokt as root directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultRootDirectory = homeDir + "/.cmt-pokt"
}

var (
	// DefaultRootDirectory is root directory for the pocket node is initialized in the init function to be the home directory + /.cmt-pokt
	DefaultRootDirectory = ""

	// persistence
	DefaultPersistencePostgresURL    = "postgres://postgres:postgres@pocket-db:5432/postgres"
	DefaultPersistenceBlockStorePath = "/var/blockstore"

	// logger
	DefaultLoggerLevel  = "debug"
	DefaultLoggerFormat = "pretty"

	// blocks per session
	DefaultBlocksPerSession = uint64(2)
)
