package imagestorage

import (
	"context"
	"io"
	"net/url"
)

type Getter interface {
	Get(ctx context.Context, filename string) (*url.URL, error)
}

//tx, image.FileName, image.File, image.Size)
type Setter interface {
	Set(ctx context.Context, filename string, size int64, file io.Reader) error
}

type GetterSetter interface {
	Getter
	Setter
}
