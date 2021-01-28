package build

import (
	"log"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

func NewMinIOStorage(config *RemoteStorageConfig) (*remotestorage.MinIOClient, error) {

	client, err := remotestorage.NewMinIOClient(config.Endpoint, config.AccessKeyID, config.AccessKeySecret, config.useSSL)
	if err != nil {
		log.Println("err building minio: ", err)
		return nil, err
	}

	return client, nil
}
