package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/handlers"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

func main() {
	fmt.Println("starting application")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	//create handler
	// get client
	minioClient, err := remotestorage.NewMinIOClient()
	if err != nil {
		log.Fatalln(err)
	}

	//make new bucket
	if err := minioClient.MakeNewBucket(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// upload image
	if err := minioClient.UploadImage(context.Background(), "mytestbucket"); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("successfulling made bucket and uploaded image")

	indexHandler := &handlers.IndexHandler{
		RemoteStore: minioClient,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go startServer(context.Background(), &wg, interrupt, indexHandler)
	wg.Wait()

}

func startServer(ctx context.Context, wg *sync.WaitGroup, interrupt chan os.Signal, index http.Handler) {

	// define handler func for "/"
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.Handle("/testfiles/", http.StripPrefix("/testfiles", http.FileServer(http.Dir("testfiles"))))

	http.Handle("/", index)
	// build.Build()

	// Start server -- listen at localhost, port 8080
	go func() {
		fmt.Println("starting server of 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	//remotestorage.NewMinIOClient()

	// fmt.Println("successfully gat new client, made bucket, uploaded image, got URL")

	// fmt.Println("Successfully generated presigned URL", signedURL)

	<-interrupt
	wg.Done()
}
