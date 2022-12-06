package handlers

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/schema"
	"github.com/jjg-akers/simple-image-sharing-webapp/domain/imagemanager"
)

type UploadHandler struct {
	Decoder      *schema.Decoder
	ImageHandler imagemanager.ImageService
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// parse submitted file
	mf, fh, err := r.FormFile("nf")
	if err != nil {
		fmt.Println("err reading form: ", err)
	}
	defer mf.Close()

	// parse form fields
	imageMeta := &imagemanager.Meta{}

	if err = h.Decoder.Decode(imageMeta, r.PostForm); err != nil {
		log.Println("err decoding post form: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	image := &imagemanager.ImageV1{
		Meta: imageMeta,
	}

	//create hash for filename
	ext := path.Ext(fh.Filename)
	fileHash := sha1.New()
	io.Copy(fileHash, mf)
	fileName := fmt.Sprintf("%x", fileHash.Sum(nil)) + ext

	//reset
	mf.Seek(0, 0)

	image.File = mf
	image.Meta.Size = fh.Size
	image.Meta.FileName = fileName
	image.Meta.DateAdded = time.Now()

	//attempt upload
	if err := h.ImageHandler.Upload(r.Context(), image); err != nil {
		fmt.Println("index handler failed upload: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	relatedImages, err := h.ImageHandler.Retrieve(r.Context(), []string{imageMeta.Tag})
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", relatedImages)
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
