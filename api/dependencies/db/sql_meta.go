package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jjg-akers/simple-image-sharing-webapp/domain"
	"github.com/jjg-akers/simple-image-sharing-webapp/domain/imagemanager"
)

// runtime validation
var _ imagemanager.MetaRepo = &SQLGetterSetter{}

type SQLGetterSetter struct {
	DB *sql.DB
}

func NewSQLDBManager(db *sql.DB) *SQLGetterSetter {
	return &SQLGetterSetter{
		DB: db,
	}
}

func (gs *SQLGetterSetter) GetMeta(ctx context.Context, tags []string) ([]*imagemanager.Meta, error) {
	//select filename from DB where tag in(....)
	//build query
	query, _ := NewQuery(Tags(tags))

	//get params
	args, err := NewSQLParams(StringParam(tags))
	if err != nil {
		fmt.Println("err creating sql params: ", err)
		return nil, err
	}

	rows, err := gs.DB.QueryContext(ctx, query, args.Params...)
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

	if len(toReturn) == 0 {
		log.Println("no results")
		return nil, domain.ErrNotFound
	}

	return toReturn, nil

}

func (gs *SQLGetterSetter) SetMeta(ctx context.Context, meta *imagemanager.Meta) error {
	query := fmt.Sprintf("INSERT INTO `photoshare`.`images` (`image_name`, `tag`, `title`, `description`, `date_added`) VALUES (?, ?, ?, ?, ?);")

	_, err := gs.DB.ExecContext(ctx, query, meta.FileName, meta.Tag, meta.Title, meta.Description, meta.DateAdded)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // duplicate insert, run as update?
				//query = "UPDATE photoshare.images (imageName, tag, date_added) VALUES (?, ?, ?)"
				// just log for now
				log.Println("got duplicate insert error: ", err)
				return nil
			default:
				return err
			}
		}
		return err
	}

	log.Printf("successfully inserted image %s with tag %s.\n", meta.FileName, meta.Tag)

	return nil
}

func buildSearchQuery(numerOfArgs int) string {

	sb := strings.Builder{}

	sb.WriteString("SELECT image_name, tag, title, description FROM photoshare.images WHERE tag IN(")

	sb.WriteString(strings.Repeat("?,", numerOfArgs-1))
	sb.WriteString("?);")

	return sb.String()

}

type Option func(sb *strings.Builder) error

func Tags(tag []string) Option {
	return func(sb *strings.Builder) error {
		sb.WriteString(" tag IN(")
		sb.WriteString(strings.Repeat("?,", len(tag)-1))
		sb.WriteString("?)")

		return nil
	}
}

func NewQuery(opts ...Option) (string, error) {

	sb := strings.Builder{}

	sb.WriteString("SELECT image_name, tag, title, description FROM photoshare.images WHERE")

	for _, opt := range opts {
		err := opt(&sb)
		if err != nil {
			return "", err
		}

		sb.WriteString(" and")
	}

	query := strings.TrimSuffix(sb.String(), " and") + ";"

	return query, nil
}
