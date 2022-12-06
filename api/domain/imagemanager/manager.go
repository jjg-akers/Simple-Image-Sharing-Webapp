package imagemanager

import (
	"context"
)

type Uploader interface {
	// Upload(ctx context.Context, imageName, tag string) error
	Upload(ctx context.Context, image *ImageV1) error
}

type Retriever interface {
	Retrieve(ctx context.Context, tags []string) ([]*ImageV1, error)
	Random(ctx context.Context, n int) ([]*ImageV1, error)
}

type ImageService interface {
	Uploader
	Retriever
}
