package remotestorage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	Client *minio.Client
}

func NewMinIOClient(endpoint, accessKeyID, accessKeySecret string, useSSL bool) (*MinIOClient, error) {
	// endpoint := "localhost:9000"
	// accessKeyID := "minioadmin"
	// secretAccessKey := "minioadmin"
	// useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, accessKeySecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println("failed to create new client")
		return nil, err
	}

	return &MinIOClient{
		Client: minioClient,
	}, nil
}

func (mc *MinIOClient) MakeNewBucket(ctx context.Context, bucketName, location string) error {

	//bucketName := "mytestbucket"
	//location := "us-east-1"

	err := mc.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := mc.Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
			return nil
		} else {
			log.Println("error creating buckdt: ", err)
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
		return nil
	}
}

func (mc MinIOClient) UploadImage(ctx context.Context, bucketName, imageName, filePath string) error {

	contentType := "application/jpg"

	n, err := mc.Client.FPutObject(ctx, bucketName, imageName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error uploading image: ", err)
		return err
		//log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", imageName, n.Size)
	return nil
}

func (mc *MinIOClient) NewPresignedGet(ctx context.Context, bucketName, objectName string) (*url.URL, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"Blackmore.jpg\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := mc.Client.PresignedGetObject(ctx, bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println("error generating presignedGet url: ", err)
		return nil, err
	}

	return presignedURL, nil
}
