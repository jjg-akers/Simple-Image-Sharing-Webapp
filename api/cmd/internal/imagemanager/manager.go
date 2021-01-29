package imagemanager

import (
	"context"
)

type Searcher interface {
	Search(ctx context.Context, tags []string) ([]string, error)
}

type Uploader interface {
	// Upload(ctx context.Context, imageName, tag string) error
	Upload(ctx context.Context, image *Image) error
}

type SearcherUploader interface {
	Searcher
	Uploader
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
