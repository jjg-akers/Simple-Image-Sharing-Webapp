package imagestorage

import (
	"context"
	"io"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/meta"
)

type Getter interface {
	Get(ctx context.Context, metas []*meta.Meta) ([]*ImageV1, error)
}

//tx, image.FileName, image.File, image.Size)
type Setter interface {
	Set(ctx context.Context, filename string, size int64, file io.Reader) error
}

type GetterSetter interface {
	Getter
	Setter
}
