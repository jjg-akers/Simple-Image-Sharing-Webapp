package imagemanager

import (
	"context"
	"errors"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/imagestorage"
)

var ErrNotFound = errors.New("No images found in db for given tag")

// type Searcher interface {
// 	Search(ctx context.Context, tags []string) ([]*Image, error)
// }

type Uploader interface {
	// Upload(ctx context.Context, imageName, tag string) error
	Upload(ctx context.Context, image *imagestorage.ImageV1) error
}

// type SearcherUploader interface {
// 	Searcher
// 	Uploader
// }

type Retriever interface {
	Retrieve(ctx context.Context, tags []string) ([]*imagestorage.ImageV1, error)
}

type UploaderRetriever interface {
	Uploader
	Retriever
}

// func (mc MinIOClient) UploadImage(ctx context.Context, bucketName, imageName, filePath string) error {

// 	contentType := "application/jpg"

// 	n, err := mc.Client.FPutObject(ctx, bucketName, imageName, filePath, minio.PutObjectOptions{ContentType: contentType})
// 	if err != nil {
// 		log.Println("error uploading image: ", err)
// 		return err
// 		//log.Fatalln(err)
// 	}

// 	log.Printf("Successfully uploaded %s of size %d\n", imageName, n.Size)
// 	return nil
// }
