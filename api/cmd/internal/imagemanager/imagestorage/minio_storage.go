package imagestorage

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

var _ GetterSetter = &MinioStorage{}

// will implement the uploaderSearcher interface
type MinioStorage struct {
	Client *remotestorage.MinIOClient
}

//client *remotestorage.UploaderRetriever
func NewMinioStorage(client *remotestorage.MinIOClient) *MinioStorage {
	return &MinioStorage{
		Client: client,
	}

}

//Get gets a single signeduri
func (mm *MinioStorage) Get(ctx context.Context, filename string) (*url.URL, error) {

	url, err := mm.Client.Get(ctx, filename)
	if err != nil {
		fmt.Println("err Getting urls from minio image getter")
		return nil, err
	}

	return url, nil
}

//Set saves a image
func (mm *MinioStorage) Set(ctx context.Context, filename string, size int64, file io.Reader) error {
	if err := mm.Client.Upload(ctx, filename, file, size); err != nil {
		fmt.Println("err Uploading image from minio setter: ", err)
		return err
	}
	return nil
}
