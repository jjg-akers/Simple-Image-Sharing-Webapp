package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("cmd/templates/*.gohtml"))
}

type IndexHandler struct {
	RemoteStore  *remotestorage.MinIOClient
	DB           *sql.DB
	ImageManager imagemanager.SearcherUploader
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		fmt.Println("method post")

		mf, fh, err := r.FormFile("nf")
		if err != nil {
			fmt.Println("err reading form: ", err)
		}
		defer mf.Close()

		//get filename
		fmt.Println("filename: ", fh.Filename)

		//create image
		im := &imagemanager.Image{
			Name:      fh.Filename,
			File:      mf,
			Tag:       "test_tage",
			Size:      fh.Size,
			DateAdded: time.Now(),
		}

		//attempt upload

		if err := h.ImageManager.Upload(r.Context(), im); err != nil {
			fmt.Println("index handler failed upload: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//attempt sql insert
		// err = h.ImageManager.Upload(r.Context(), fh.Filename, "test_tag")
		// if err != nil {
		// 	fmt.Println("index handler err uploading: ", err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// } else {
		// 	fmt.Println("successful upload")
		// }

		//try to upload it
		//h.RemoteStore.UploadImage(r.Context(), "testy-mctest-face", fh.Filename)
		// contentType := "application/jpg"
		// info, err := h.RemoteStore.Client.PutObject(r.Context(), "testy-mctest-face", fh.Filename, mf, fh.Size, minio.PutObjectOptions{ContentType: contentType})
		// if err != nil {
		// 	fmt.Println("failed to put file: ", err)
		// }

		// fmt.Printf("succesffully put file. location: %s, size: %d\n", info.Location, info.Size)

		//w.WriteHeader(http.StatusCreated)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		//tpl.ExecuteTemplate(w, "index.gohtml", "")
		//return
		// contentType := "application/jpg"

		// n, err := mc.Client.FPutObject(ctx, bucketName, imageName, filePath, minio.PutObjectOptions{ContentType: contentType})
	}

	fmt.Println("method GET")

	type pathserver struct {
		Paths []*url.URL
	}

	// for now just throw sthe default images up
	paths, err := getDefaultImages(r.Context(), h.RemoteStore)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// load the paths to be served
	ps := pathserver{
		Paths: paths,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "index.gohtml", ps)
}

func getDefaultImages(ctx context.Context, client *remotestorage.MinIOClient) ([]*url.URL, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("coule not get wd: ", err)
	}
	dir, err := os.Open(filepath.Join(wd, "cmd/testfiles"))
	if err != nil {
		return nil, fmt.Errorf("failed opening directory: %s", err)
	}
	files, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed opening directory: %s", err)
	}
	dir.Close()

	paths := make([]*url.URL, len(files))

	for i, file := range files {
		imageName := strings.TrimSuffix(file, filepath.Ext(file))
		signedURL, err := client.NewPresignedGet(ctx, imageName)
		//fmt.Println("url: ", signedURL)
		if err != nil {
			return nil, fmt.Errorf("error getting signed url: %s", err)
		}

		paths[i] = signedURL
	}

	return paths, nil
}
