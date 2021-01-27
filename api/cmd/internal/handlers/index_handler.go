package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	type pathserver struct {
		Paths []string
	}

	paths := []string{"testfiles/Blackmore.jpg", "testfiles/lightswitch wiring.jpg", "testfiles/test.png"}
	ps := pathserver{
		Paths: paths,
	}

	fmt.Println("executing indexHandler")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	tpl.ExecuteTemplate(w, "index.gohtml", ps)
}
