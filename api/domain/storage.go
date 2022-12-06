package domain

import (
	"context"
	"io"
)

type ImageGetter interface {
	GetImage(ctx context.Context, metas []*Meta) ([]*ImageV1, error)
}

type ImageSetter interface {
	SetImage(ctx context.Context, filename string, size int64, file io.Reader) error
}

type ImmageRepo interface {
	ImageGetter
	ImageSetter
}
