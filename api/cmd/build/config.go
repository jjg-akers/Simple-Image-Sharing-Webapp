package build

import (
	"github.com/urfave/cli/v2"
)

type AppConfig struct {
	DBConfig      *SQLDBConfig
	StorageConfig *RemoteStorageConfig
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		DBConfig:      NewSQLDBConfig(),
		StorageConfig: NewRemoteSotrageConfig(),
	}
}

func LoadAppConfig(config *AppConfig) []cli.Flag {
	var flags []cli.Flag

	flags = append(flags, LoadSQLDBConfig(config.DBConfig)...)
	flags = append(flags, LoadStorageConfig(config.StorageConfig)...)

	return flags
}
