package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jjg-akers/simple-image-sharing-webapp/dependencies/remotestorage"
	"github.com/jjg-akers/simple-image-sharing-webapp/domain/imagemanager"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("cmd/../../ui/templates/*.html"))
}

type IndexHandler struct {
	RemoteStore *remotestorage.MinIOClient
	DB          *sql.DB
	ImageGetter imagemanager.ImageGetter
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	images, err := getInitialImages(r.Context(), 10, h.DB, h.ImageGetter)
	if err != nil {
		log.Println("failed getting initial images. err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tpl.ExecuteTemplate(w, "index.html", images)
}

func getInitialImages(ctx context.Context, n int, db *sql.DB, imageGetter imagemanager.ImageGetter) ([]*imagemanager.ImageV1, error) {
	meta, err := selectRandom(ctx, 10, db)
	if err != nil {
		return nil, err
	}

	images, err := imageGetter.GetImage(ctx, meta)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func selectRandom(ctx context.Context, n int, db *sql.DB) ([]*imagemanager.Meta, error) {

	query := fmt.Sprintf(`SELECT image_name, tag, title, description 
	FROM photoshare.images AS r1 JOIN (SELECT CEIL(RAND() * (SELECT MAX(id) FROM photoshare.images)) AS id) AS r2 
	WHERE r1.id >= r2.id ORDER BY r1.id ASC LIMIT 6;`)

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not query db: %s", err)
	}

	defer rows.Close()

	var toReturn []*imagemanager.Meta

	for rows.Next() {
		m := imagemanager.Meta{}

		err = rows.Scan(&m.FileName, &m.Tag, &m.Title, &m.Description)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %s", err)
		}

		toReturn = append(toReturn, &m)
	}

	return toReturn, nil
}
