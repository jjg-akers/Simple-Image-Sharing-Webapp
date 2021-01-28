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

	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/remotestorage"
)

var tpl *template.Template

func init() {
	// wd, _ := os.Getwd()
	// log.Println(wd)
	tpl = template.Must(template.ParseGlob("cmd/templates/*.gohtml"))
}

type IndexHandler struct {
	RemoteStore *remotestorage.MinIOClient
	DB          *sql.DB
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//fmt.Println("executing indexHandler")

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
	w.WriteHeader(http.StatusOK)
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
		signedURL, err := client.NewPresignedGet(ctx, "testy-mctest-face", imageName)
		if err != nil {
			return nil, fmt.Errorf("error getting signed url: %s", err)
		}

		paths[i] = signedURL
	}

	return paths, nil
}
