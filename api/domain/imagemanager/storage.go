package imagemanager

import (
	"context"
	"io"
)



type ImageGetter interface {
	GetImage(ctx context.Context, metas []*Meta) ([]*ImageV1, error)
}

//tx, image.FileName, image.File, image.Size)
type ImageSetter interface {
	SetImage(ctx context.Context, filename string, size int64, file io.Reader) error
}

type ImmageRepo interface {
	ImageGetter
	ImageSetter
}
