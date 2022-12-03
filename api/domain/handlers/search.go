package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/jjg-akers/simple-image-sharing-webapp/domain"
	"github.com/jjg-akers/simple-image-sharing-webapp/domain/imagemanager"
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
	switch err {
	case nil:
	case domain.ErrNotFound:
		log.Println("errnotfound")
		// w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		tpl.ExecuteTemplate(w, "index.html", images)
		return

	default:
		log.Println("search handler failed url search: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", images)
}
