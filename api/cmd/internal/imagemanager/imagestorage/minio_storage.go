package imagestorage

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	//"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/meta"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"

	"golang.org/x/sync/errgroup"
)

type ImageV1 struct {
	Meta *meta.Meta `json:"Meta"`
	URI  string     `json:"url"`
	File io.Reader  `json:"File"`
}

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

//Get gets a []Images
func (mm *MinioStorage) Get(ctx context.Context, metas []*meta.Meta) ([]*ImageV1, error) {
	// get signed urls
	g, ctx := errgroup.WithContext(ctx)
	imageChan := make(chan *ImageV1)

	done := int32(len(metas))
	for _, meta := range metas {
		meta := meta
		g.Go(func() error {

			// use defer func to close chan after last out
			defer func() {
				if atomic.AddInt32(&done, -1) == 0 {
					close(imageChan)
				}
			}()

			signedURI, err := mm.Client.Get(ctx, meta.FileName)
			if err != nil {
				return fmt.Errorf("Failed to get uri for file %s, err: %s", meta.FileName, err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:

				image := &ImageV1{
					Meta: meta,
					URI:  signedURI.String(),
					// URI: signedURI,
				}

				imageChan <- image
			}
			return nil
		})
	}

	imagesToReturn := make([]*ImageV1, len(metas))

	g.Go(func() error {
		i := 0
		for im := range imageChan {
			imagesToReturn[i] = im
			i++
		}
		return nil
	})

	return imagesToReturn, g.Wait()
}

//Set saves a image
func (mm *MinioStorage) Set(ctx context.Context, filename string, size int64, file io.Reader) error {
	if err := mm.Client.Upload(ctx, filename, file, size); err != nil {
		fmt.Println("err Uploading image from minio setter: ", err)
		return err
	}
	return nil
}
