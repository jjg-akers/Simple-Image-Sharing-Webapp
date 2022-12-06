package build

import (
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinIOStorage(ctx context.Context, config *RemoteStorageConfig) (*minio.Client, error) {

	// Initialize minio client object.
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.AccessKeySecret, ""),
		Secure: config.useSSL,
	})
	if err != nil {
		log.Println("failed to create new minio client")
		return nil, err
	}

	if err := initializeBucket(ctx, client, "testy-mctest-face", "us-east-1"); err != nil {
		return nil, fmt.Errorf("failed to create new bucket, err: %s", err)
	}

	return client, nil
}

func initializeBucket(ctx context.Context, client *minio.Client, bucketName, location string) error {

	//bucketName := "mytestbucket"
	//location := "us-east-1"

	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
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
