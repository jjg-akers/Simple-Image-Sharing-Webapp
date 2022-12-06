package domain

import (
	"context"
	"fmt"
	"log"
)

var _ ImageService = &SQLMinIOImpl{}

type SQLMinIOImpl struct {
	Meta    MetaRepo
	Storage ImmageRepo
}

func (sm *SQLMinIOImpl) Retrieve(ctx context.Context, tags []string) ([]*ImageV1, error) {
	//Get(ctx context.Context, tags []string) ([]*Meta, error)
	imageMetas, err := sm.Meta.GetMeta(ctx, tags)
	if err != nil {
		log.Println("error retreiving meta: ", err)
		return nil, err
	}

	//get images
	images, err := sm.Storage.GetImage(ctx, imageMetas)
	if err != nil {
		log.Println("error getting images: ", err)
		return nil, err
	}

	return images, nil
}

func (sm *SQLMinIOImpl) Random(ctx context.Context, n int) ([]*ImageV1, error) {
	meta, err := sm.Meta.GetRandom(ctx, n)
	if err != nil {
		log.Println("error retreiving meta: ", err)
		return nil, err
	}

	images, err := sm.Storage.GetImage(ctx, meta)
	if err != nil {
		return nil, err
	}

	return images, nil
	// imageMetas, err := sm.Meta.GetMeta(ctx, tags)
	// if err != nil {
	// 	log.Println("error retreiving meta: ", err)
	// 	return nil, err
	// }

	// //get images
	// images, err := sm.Storage.GetImage(ctx, imageMetas)
	// if err != nil {
	// 	log.Println("error getting images: ", err)
	// 	return nil, err
	// }

	// return images, nil
}

func (sm *SQLMinIOImpl) Upload(ctx context.Context, image *ImageV1) error {
	// load meta
	if err := sm.Meta.SetMeta(ctx, image.Meta); err != nil {
		fmt.Println("failed upload insert")
		return err
	}

	if err := sm.Storage.SetImage(ctx, image.Meta.FileName, image.Meta.Size, image.File); err != nil {
		fmt.Println("failed upload to storage")
		return err
	}

	return nil
}
