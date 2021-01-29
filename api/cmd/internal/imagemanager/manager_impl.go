package imagemanager

import (
	"context"
	"fmt"
)

var _ SearcherUploader = &ImageManager{}

type ImageManager struct {
	StorageManager SearcherUploader
	DBManager      DBImageManager
}

// func NewImageManager(db SearcherUploader, storageMan remotestorage.UploaderGetter) *ImageManager {
// 	return &ImageManager{
// 		StorageManager
// 		DBManager: db,
// 	}
// }

func (im *ImageManager) Search(ctx context.Context, tags []string) ([]string, error) {
	// get tags from db
	files, err := im.DBManager.Search(ctx, tags)
	if err != nil {
		return nil, err
	}

	//return signed urls
	urls, err := im.StorageManager.Search(ctx, files)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (im *ImageManager) Upload(ctx context.Context, image *Image) error {
	// attempt SQL insert of meta{
	if err := im.DBManager.Insert(ctx, image.Name, image.Tag); err != nil {
		fmt.Println("failed upload insert")
		return err
	}

	//upload to storage
	//im.StorageManager.Upload()

	if err := im.StorageManager.Upload(ctx, image); err != nil {

		fmt.Println("Image manager failed uploading: ", err)
		//w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
