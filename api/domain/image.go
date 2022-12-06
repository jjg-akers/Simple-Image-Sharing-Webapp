package domain

import (
	"context"
	"io"
)

type Uploader interface {
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

type ImageV1 struct {
	Meta *Meta     `json:"Meta"`
	URI  string    `json:"url"`
	File io.Reader `json:"File"`
}
