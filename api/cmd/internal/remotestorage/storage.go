package remotestorage

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	Client *minio.Client
}

func NewMinIOClient() (*MinIOClient, error) {
	//ctx := context.Background()
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
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

func (mc *MinIOClient) MakeNewBucket(ctx context.Context) error {

	// Make a new bucket called mymusic.
	bucketName := "mytestbucket"
	location := "us-east-1"

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

func (mc MinIOClient) UploadImage(ctx context.Context, bucketName string) error {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println(path)
	// Upload the zip file
	objectName := "Blackmore.jpg"
	filePath := "./testfiles/Blackmore.jpg"
	contentType := "application/jpg"

	// Upload the zip file with FPutObject
	n, err := mc.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Println("error uploading image: ", err)
		return err
		//log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
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
