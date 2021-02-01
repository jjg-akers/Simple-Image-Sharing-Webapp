package imagemanager

import (
	"context"
	"fmt"
	"log"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/imagestorage"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/meta"
)

var _ UploaderRetriever = &SQLMinIOImpl{}

type SQLMinIOImpl struct {
	Meta    meta.GetterSetter
	Storage imagestorage.GetterSetter
}

func (sm *SQLMinIOImpl) Retrieve(ctx context.Context, tags []string) ([]*imagestorage.ImageV1, error) {
	//Get(ctx context.Context, tags []string) ([]*Meta, error)
	imageMetas, err := sm.Meta.Get(ctx, tags)
	if err != nil {
		log.Println("error retreiving meta: ", err)
		return nil, err
	}

	//get images
	images, err := sm.Storage.Get(ctx, imageMetas)
	if err != nil {
		log.Println("error getting images: ", err)
		return nil, err
	}

	return images, nil
}

func (sm *SQLMinIOImpl) Upload(ctx context.Context, image *imagestorage.ImageV1) error {
	// load meta
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
