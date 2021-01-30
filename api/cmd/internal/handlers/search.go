package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/imagemanager"
)

type SearchRequestParams struct {
	Tag []string `schema:"tag,required"`
}

type SearchHandler struct {
	// ImageManager imagemanager.Searcher
	ImageRetriever imagemanager.Retriever
	Decoder        *schema.Decoder
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//1. parse get params
	//2. run search
	//3. return template

	rp := &SearchRequestParams{}

	err := h.Decoder.Decode(rp, r.URL.Query())
	if err != nil {
		log.Println("Could not decode request params: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(rp.Tag) == 0 || strings.TrimSpace(rp.Tag[0]) == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	// get images
	images, err := h.ImageRetriever.Retrieve(r.Context(), rp.Tag)
	//urls, err := h.ImageManager.Search(r.Context(), rp.Tag)
	switch err {
	case nil:
	case imagemanager.ErrNotFound:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tpl.ExecuteTemplate(w, "index.gohtml", "")
	default:
		log.Println("search handler failed url search: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println("images from search: ", images)

	p := []string{}

	for _, image := range images {
		p = append(p, image.URI)
	}

	type pathserver struct {
		Paths []string
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl.ExecuteTemplate(w, "index.gohtml", pathserver{Paths: p})
}
