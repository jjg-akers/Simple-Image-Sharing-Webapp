package build

import "github.com/urfave/cli/v2"

type RemoteStorageConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	useSSL          bool
}

func NewRemoteSotrageConfig() *RemoteStorageConfig {
	return &RemoteStorageConfig{}
}

func LoadStorageConfig(config *RemoteStorageConfig) []cli.Flag {
	//var flags []cli.Flag
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "remote-storage-endpoint",
			EnvVars:     []string{"REMOTE_STORAGE_ENDPOINT"},
			Destination: &config.Endpoint,
		},
		&cli.StringFlag{
			Name:        "remote-storage-access-key-id",
			EnvVars:     []string{"REMOTE_STORAGE_ACCESS_KEY_ID"},
			Destination: &config.AccessKeyID,
		},
		&cli.StringFlag{
			Name:        "remote-storage-access-key-secret",
			EnvVars:     []string{"REMOTE_STORAGE_ACCESS_KEY_SECRET"},
			Destination: &config.AccessKeySecret,
		},
		&cli.BoolFlag{
			Name:        "remote-storage-use-ssl",
			EnvVars:     []string{"REMOTE_STORAGE_USE_SSL"},
			Destination: &config.useSSL,
			Value:       false,
		},
	}
}
