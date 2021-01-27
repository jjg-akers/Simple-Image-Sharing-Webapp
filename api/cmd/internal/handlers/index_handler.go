package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

type IndexHandler struct {
	RemoteStore *remotestorage.MinIOClient
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("executing indexHandler")

	type pathserver struct {
		Paths []*url.URL
	}

	signedURL, err := h.RemoteStore.NewPresignedGet(r.Context(), "mytestbucket", "Blackmore.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	//paths := []string{"testfiles/Blackmore.jpg", "testfiles/lightswitch wiring.jpg", "testfiles/test.png"}
	ps := pathserver{
		Paths: []*url.URL{signedURL},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	tpl.ExecuteTemplate(w, "index.gohtml", ps)
}
