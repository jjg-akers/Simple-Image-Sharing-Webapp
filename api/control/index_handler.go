package control

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jjg-akers/simple-image-sharing-webapp/domain"
)

var tpl *template.Template

func init() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

type IndexHandler struct {
	ImageRetriever domain.Retriever
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	images, err := h.ImageRetriever.Random(r.Context(), 10)
	if err != nil {
		log.Println("failed getting initial images. err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", images)
}
