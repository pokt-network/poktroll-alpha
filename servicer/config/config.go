package config

type ServicerConfig struct {
	BlocksPerSession int64
}

func DefaultConfig() ServicerConfig {
	return ServicerConfig{
		// TECHDEBT: I have no idea what this value should be.
		BlocksPerSession: 10,
	}
}
