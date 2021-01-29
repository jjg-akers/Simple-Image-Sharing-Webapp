package imagemanager

import (
	"context"
	"fmt"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

var _ SearcherUploader = &MinioImageManager{}

// will implement the uploaderSearcher interface
type MinioImageManager struct {
	Client *remotestorage.MinIOClient
}

//client *remotestorage.UploaderRetriever
func NewMinioManager(client *remotestorage.MinIOClient) *MinioImageManager {
	return &MinioImageManager{
		Client: client,
	}

}

func (mm *MinioImageManager) Search(ctx context.Context, tags []string) ([]string, error) {

	urls, err := mm.Client.Get(ctx, tags)
	if err != nil {
		fmt.Println("err Getting urls from minio image manager")
		return nil, err
	}

	return urls, nil
}

func (mm *MinioImageManager) Upload(ctx context.Context, image *Image) error {
	if err := mm.Client.Upload(ctx, image.Name, image.File, image.Size); err != nil {
		fmt.Println("err Uploading image from minioImageManager: ", err)
		return err
	}
	return nil
}
