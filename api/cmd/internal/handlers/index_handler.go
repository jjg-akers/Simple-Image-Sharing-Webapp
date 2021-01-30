package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/schema"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager/meta"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("cmd/templates/*.html"))
}

type IndexHandler struct {
	RemoteStore *remotestorage.MinIOClient
	// DB           *sql.DB
	Decoder *schema.Decoder
	// ImageManager imagemanager.SearcherUploader
	ImageHandler imagemanager.UploaderRetriever
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if r.Method == http.MethodPost {

	// 	fmt.Println("method post")

	// 	// parse submitted file
	// 	mf, fh, err := r.FormFile("nf")
	// 	if err != nil {
	// 		fmt.Println("err reading form: ", err)
	// 	}
	// 	defer mf.Close()

	// 	// parse form fields
	// 	imageMeta := &meta.Meta{}

	// 	if err = h.Decoder.Decode(imageMeta, r.PostForm); err != nil {
	// 		log.Println("err decoding post form: ", err)
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		return
	// 	}

	// 	image := &imagemanager.ImageV1{
	// 		Meta: imageMeta,
	// 	}

	// 	//fmt.Println("IMAGEVE: ", image)

	// 	//create hash for filename
	// 	ext := path.Ext(fh.Filename)
	// 	fileHash := sha1.New()
	// 	io.Copy(fileHash, mf)
	// 	fileName := fmt.Sprintf("%x", fileHash.Sum(nil)) + ext

	// 	//reset
	// 	mf.Seek(0, 0)

	// 	image.File = mf
	// 	image.Meta.Size = fh.Size
	// 	image.Meta.FileName = fileName
	// 	image.Meta.DateAdded = time.Now()

	// 	//attempt upload
	// 	if err := h.ImageHandler.Upload(r.Context(), image); err != nil {
	// 		fmt.Println("index handler failed upload: ", err)
	// 		w.WriteHeader(http.StatusInternalServerError)
	// 		return
	// 	}

	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }

	fmt.Println("method GET")

	type pathserver struct {
		Paths []*url.URL
	}

	// for now just throw sthe default images up
	images, err := getDefaultImages(r.Context(), h.RemoteStore)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	// load the paths to be served
	// ps := pathserver{
	// 	Paths: paths,
	// }

	w.Header().Set("Content-Type", "text/html")
	//w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "index.html", images)
}

func getDefaultImages(ctx context.Context, client *remotestorage.MinIOClient) ([]*imagemanager.ImageV1, error) {
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

	// imageMeta := &meta.Meta{}

	images := make([]*imagemanager.ImageV1, len(files))

	// paths := make([]*url.URL, len(files))

	for i, file := range files {
		imageName := strings.TrimSuffix(file, filepath.Ext(file))
		signedURL, err := client.NewPresignedGet(ctx, imageName)
		//fmt.Println("url: ", signedURL)
		if err != nil {
			return nil, fmt.Errorf("error getting signed url: %s", err)
		}

		imageMeta := &meta.Meta{
			Tag:         "default",
			Title:       imageName,
			Description: "stuff about the photo",
		}

		image := &imagemanager.ImageV1{
			Meta: imageMeta,
			URI:  signedURL.String(),
		}

		images[i] = image
		// paths[i] = signedURL
	}

	// return paths, nil
	return images, nil
}
