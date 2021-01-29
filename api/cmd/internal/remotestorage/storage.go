package remotestorage

import (
	"context"
	"io"
)

//Define the interfaces a remote storage implementation must satisfy

type Uploader interface {
	Upload(ctx context.Context, imageName string, reader io.Reader, size int64) error
}

type Getter interface {
	Get(ctx context.Context, files []string) ([]string, error)
}

type UploaderGetter interface {
	Uploader
	Getter
}
