package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/schema"
	"github.com/urfave/cli/v2"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/build"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/handlers"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/imagestorage"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/meta"
)

type PhotoShareApp struct {
	config *build.AppConfig
}

func NewApp() *PhotoShareApp {
	return &PhotoShareApp{
		config: build.NewAppConfig(),
	}
}

func main() {

	fmt.Println("starting application")

	api := newPhotoShareApp()

	//build cli app
	var app = cli.NewApp()
	app.Usage = "Allow users to upload/ view photos in webapp"

	// load ENV into config
	app.Flags = build.LoadAppConfig(api.config)

	app.Action = api.startAPI

	if err := app.Run(os.Args); err != nil {
		log.Println("Error Running application: ", err)
	}
}

type photoShareApp struct {
	config *build.AppConfig
}

func newPhotoShareApp() *photoShareApp {
	return &photoShareApp{
		config: build.NewAppConfig(),
	}
}

func (api *photoShareApp) startAPI(cliCtx *cli.Context) error {

	// build components
	db, err := build.NewSQLDB(api.config.DBConfig)
	if err != nil {
		return fmt.Errorf("Failed to build SQL DB, err: %s", err)
	}

	minioClient, err := build.NewMinIOStorage(api.config.StorageConfig)
	if err != nil {
		return fmt.Errorf("Failed to build Minio client, err: %s", err)
	}

	imager := &imagemanager.SQLMinIOImpl{
		Meta:    meta.NewSQLDBManager(db),
		Storage: imagestorage.NewMinioStorage(minioClient),
	}

	//set up handlers
	indexHandler := &handlers.IndexHandler{
		RemoteStore: minioClient,
		DB:          db,
		ImageGetter: imagestorage.NewMinioStorage(minioClient),
	}

	searchHandler := &handlers.SearchHandler{
		ImageRetriever: imager,
		Decoder:        schema.NewDecoder(),
	}

	uploadHandler := &handlers.UploadHandler{
		Decoder:      schema.NewDecoder(),
		ImageHandler: imager,
	}

	// initialize a bucket and put some random photos in it
	if err := minioClient.MakeNewBucket(cliCtx.Context, "testy-mctest-face", "us-east-1"); err != nil {
		return fmt.Errorf("Failed to create new bucket, err: %s", err)
	}

	if err := uploadStockImages(cliCtx.Context, imager); err != nil {
		return fmt.Errorf("failed to upload stock images")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go startServer(cliCtx.Context, &wg, interrupt, indexHandler, searchHandler, uploadHandler)
	wg.Wait()

	return nil
}

func startServer(ctx context.Context, wg *sync.WaitGroup, interrupt chan os.Signal, index, search, upload http.Handler) {

	http.Handle("/favicon.ico", http.NotFoundHandler())
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/cmd/static"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/ui/static"))))

	http.Handle("/upload", upload)
	http.Handle("/search", search)
	http.Handle("/", index)

	// Start server -- listen at localhost, port 8080
	go func() {
		fmt.Println("starting server of 80")
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	<-interrupt
	wg.Done()
}

// function is called on application start up  just
// to get some photos into the db for demonstration
func uploadStockImages(ctx context.Context, imageUploader imagemanager.Uploader) error {
	// upload some default images
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("coule not get wd: ", err)
	}

	dir, err := os.Open(filepath.Join(wd, "cmd/testfiles"))
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}

	fileInfo, err := dir.Readdir(-1)
	if err != nil {
		log.Fatalf("failed reading directory: %s", err)
	}

	defer dir.Close()

	for _, file := range fileInfo {
		fileName := file.Name()

		var tag string

		switch {
		case strings.Contains(fileName, "cat"):
			tag = "cat food"
		case strings.Contains(fileName, "dog"):
			tag = "dog food"
		default:
			tag = "unknown"
		}

		fileReader, err := os.Open(filepath.Join(wd, "cmd/testfiles", fileName))
		if err != nil {
			log.Println("failed to open file. err: ", err)
			return err
		}

		image := &imagestorage.ImageV1{
			Meta: &meta.Meta{
				FileName:    fileName,
				Tag:         tag,
				Title:       fileName,
				Description: "Your mom goes to college",
				Size:        file.Size(),
				DateAdded:   time.Now(),
			},
			File: fileReader,
		}

		if err = imageUploader.Upload(ctx, image); err != nil {
			log.Println("failed uploading image. err: ", err)
			return err
		}
	}

	return nil
}
