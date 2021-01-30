package imagemanager

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/imagestorage"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/meta"

	"golang.org/x/sync/errgroup"
)

type SQLMinIOImpl struct {
	Meta    meta.GetterSetter
	Storage imagestorage.GetterSetter
}

func (sm *SQLMinIOImpl) Retrieve(ctx context.Context, tags []string) ([]*ImageV1, error) {
	//Get(ctx context.Context, tags []string) ([]*Meta, error)
	imageMetas, err := sm.Meta.Get(ctx, tags)
	if err != nil {
		log.Println("error retreiving meta: ", err)
		return nil, err
	}

	// get signed urls
	g, ctx := errgroup.WithContext(ctx)
	imageChan := make(chan *ImageV1)

	done := int32(len(imageMetas))
	for _, meta := range imageMetas {
		meta := meta
		g.Go(func() error {

			// use defer func to close chan after last out
			defer func() {
				if atomic.AddInt32(&done, -1) == 0 {
					close(imageChan)
				}
			}()

			signedUri, err := sm.Storage.Get(ctx, meta.FileName)
			if err != nil {
				return fmt.Errorf("Failed to get uri for file %s, err: %s", meta.FileName, err)
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:

				image := &ImageV1{
					Meta: meta,
					URI:  signedUri.String(),
				}

				imageChan <- image
			}
			return nil
		})
	}

	imagesToReturn := make([]*ImageV1, len(imageMetas))

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

func (sm *SQLMinIOImpl) Upload(ctx context.Context, image *ImageV1) error {
	// load meta
	//image.FileName, image.Tag, image.Title, image.Description, image.DateAdded)
	if err := sm.Meta.Set(ctx, image.Meta); err != nil {
		fmt.Println("failed upload insert")
		return err
	}

	if err := sm.Storage.Set(ctx, image.Meta.FileName, image.Meta.Size, image.File); err != nil {
		fmt.Println("failed upload to storage")
		return err
	}

	return nil
}
