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

	"github.com/gorilla/schema"
	"github.com/urfave/cli/v2"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/build"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/handlers"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"
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

	// initialize a bucket and put some random photos in it
	if err := minioClient.MakeNewBucket(cliCtx.Context, "testy-mctest-face", "us-east-1"); err != nil {
		return fmt.Errorf("Failed to create new bucket, err: %s", err)
	}

	// upload some default images
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln("coule not get wd: ", err)
	}
	dir, err := os.Open(filepath.Join(wd, "cmd/testfiles"))
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	files, err := dir.Readdirnames(-1)
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	dir.Close()

	for _, file := range files {
		path := filepath.Join(wd, "cmd/testfiles", file)
		imageName := strings.TrimSuffix(file, filepath.Ext(file))

		if err := minioClient.UploadImageFromFile(cliCtx.Context, "testy-mctest-face", imageName, path); err != nil {
			return fmt.Errorf("Failed to upload new image, err: %s", err)
		}
	}

	//objectName := "Blackmore.jpg"
	//filePath := "cmd/testfiles/Blackmore.jpg"
	// if err := minioClient.UploadImage(context.Background(), "mytestbucket"); err != nil {
	// 	return fmt.Errorf("Failed to upload new image, err: %s", err)
	// }

	imager := &imagemanager.ImageManager{
		StorageManager: imagemanager.NewMinioManager(minioClient),
		DBManager:      imagemanager.NewSQLDBManager(db),
	}

	//set up handlers
	indexHandler := &handlers.IndexHandler{
		RemoteStore: minioClient,
		DB:          db,
		// ImageManager: imagemanager.NewSQLManager(db),
		ImageManager: imager,
	}

	searchHandler := &handlers.SearchHandler{
		ImageManager: imager,
		Decoder:      schema.NewDecoder(),
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go startServer(cliCtx.Context, &wg, interrupt, indexHandler, searchHandler)
	wg.Wait()

	return nil
}

func startServer(ctx context.Context, wg *sync.WaitGroup, interrupt chan os.Signal, index, search http.Handler) {

	// define handler func for "/"
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.Handle("/testfiles/", http.StripPrefix("/testfiles", http.FileServer(http.Dir("testfiles"))))

	http.Handle("/search", search)

	http.Handle("/", index)
	// build.Build()

	// Start server -- listen at localhost, port 8080
	go func() {
		fmt.Println("starting server of 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-interrupt
	wg.Done()
}
