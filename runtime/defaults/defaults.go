package defaults

import (
	"os"
)

func init() {
	initDefaultRootDirectory()
}

func initDefaultRootDirectory() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultRootDirectory = homeDir + "/.poktroll-alpha"
}

var (
	// DefaultRootDirectory is root directory for the pocket node is initialized in the init function to be the home directory + /.poktroll-alpha
	DefaultRootDirectory = ""

	// persistence
	// DefaultPersistencePostgresURL    = "postgres://postgres:postgres@pocket-db:5432/postgres"
	// DefaultPersistenceBlockStorePath = "/var/blockstore"

	// logger
	DefaultLoggerLevel  = "debug"
	DefaultLoggerFormat = "pretty"

	// blocks per session
	DefaultBlocksPerSession = int64(2)
)
