package imagemanager

import (
	"context"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/imagestorage"
)

type Uploader interface {
	// Upload(ctx context.Context, imageName, tag string) error
	Upload(ctx context.Context, image *imagestorage.ImageV1) error
}

type Retriever interface {
	Retrieve(ctx context.Context, tags []string) ([]*imagestorage.ImageV1, error)
}

type UploaderRetriever interface {
	Uploader
	Retriever
}
