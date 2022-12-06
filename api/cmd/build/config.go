package build

import (
	"github.com/urfave/cli/v2"
)

type AppConfig struct {
	Dev bool
	Version string
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
	 flags := []cli.Flag {
		&cli.BoolFlag{
			Name:        "dev",
			EnvVars:     []string{"DEV"},
			Destination: &config.Dev,
			Value: false,
		},
		&cli.StringFlag{
			Name:        "version",
			EnvVars:     []string{"VERSION"},
			Destination: &config.Version,
			Value: "V!1",
		},
	 }

	flags = append(flags, LoadSQLDBConfig(config.DBConfig)...)
	flags = append(flags, LoadStorageConfig(config.StorageConfig)...)

	return flags
}
