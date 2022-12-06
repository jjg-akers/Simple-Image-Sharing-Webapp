package remotestorage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"sync/atomic"
	"time"

	//"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"

	domain "github.com/jjg-akers/simple-image-sharing-webapp/domain/imagemanager"
	"github.com/minio/minio-go/v7"
	"golang.org/x/sync/errgroup"
)

var _ domain.ImmageRepo = &MinioStorage{}

// will implement the uploaderSearcher interface
type MinioStorage struct {
	// Client *MinIOClient
	Client *minio.Client
}

//client *remotestorage.UploaderRetriever
func NewMinioStorage(client *minio.Client) *MinioStorage {
	return &MinioStorage{
		Client: client,
	}

}

//Get gets a []Images
func (mm *MinioStorage) GetImage(ctx context.Context, metas []*domain.Meta) ([]*domain.ImageV1, error) {
	// get signed urls
	g, ctx := errgroup.WithContext(ctx)
	imageChan := make(chan *domain.ImageV1)

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

			signedURI, err := mm.get(ctx, meta.FileName)
			if err != nil {
				return fmt.Errorf("failed to get uri for file %s, err: %s", meta.FileName, err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:

				image := &domain.ImageV1{
					Meta: meta,
					URI:  signedURI.String(),
					// URI: signedURI,
				}

				imageChan <- image
			}
			return nil
		})
	}

	imagesToReturn := make([]*domain.ImageV1, len(metas))

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

func (mm *MinioStorage) get(ctx context.Context, file string) (*url.URL, error) {

	signedURL, err := mm.newPresignedGet(ctx, file)
	if err != nil {
		return nil, fmt.Errorf("error getting signed url: %s", err)
	}

	return signedURL, nil
}

func (mc *MinioStorage) newPresignedGet(ctx context.Context, objectName string) (*url.URL, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")

	bucketName := "testy-mctest-face"

	// Generates a presigned url which expires in a day.
	presignedURL, err := mc.Client.PresignedGetObject(ctx, bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println("error generating presignedGet url: ", err)
		return nil, err
	}

	return presignedURL, nil
}

//Set saves a image
func (mm *MinioStorage) SetImage(ctx context.Context, filename string, size int64, file io.Reader) error {
	if err := mm.upload(ctx, filename, file, size); err != nil {
		fmt.Println("err Uploading image from minio setter: ", err)
		return err
	}
	return nil
}

func (mm *MinioStorage) upload(ctx context.Context, imageName string, reader io.Reader, size int64) error {
	contentType := "application/jpg"
	_, err := mm.Client.PutObject(ctx, "testy-mctest-face", imageName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println("failed to put file: ", err)
		return err
	}

	//fmt.Printf("succesffully put file. location: %s, size: %d\n", info.Location, info.Size)
	return nil
}
