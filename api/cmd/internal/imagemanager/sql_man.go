package imagemanager

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jjg-akers/simple-image-sharing-webapp/cmd/internal/db"
)

// runtime validation
var _ DBImageManager = &SQLDBManager{}

type SQLDBManager struct {
	DB *sql.DB
}

//client *remotestorage.UploaderRetriever
func NewSQLDBManager(db *sql.DB) *SQLDBManager {
	return &SQLDBManager{
		DB: db,
	}
}

//	Search(ctx context.Context, tags []string) ([]string, error)

//Search queries the DB for given tags
func (m *SQLDBManager) Search(ctx context.Context, tags []string) ([]string, error) {
	//select filename from DB where tag in(....)

	//build query
	query := buildSearchQuery(len(tags))

	//get params
	args, err := db.NewSQLParams(db.StringParam(tags))
	if err != nil {
		fmt.Println("err creating sql params: ", err)
		return nil, err
	}

	//query := fmt.Sprintf("SELECT `image_name` FROM `photoshare`.`images` WHERE `tag` IN(?)")

	rows, err := m.DB.QueryContext(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("could not query db: %s", err)
	}

	defer rows.Close()

	toReturn := make([]string, len(tags))

	for rows.Next() {
		var (
			fname string
		)

		err = rows.Scan(&fname)
		if err != nil {
			return nil, fmt.Errorf("error scanning results: %s", err)
		}

		toReturn = append(toReturn, fname)
	}

	return toReturn, nil
}

//Upload saves the provided information
func (m *SQLDBManager) Insert(ctx context.Context, imageName, tag string) error {

	query := fmt.Sprintf("INSERT INTO `photoshare`.`images` (`image_name`, `tag`, `date_added`) VALUES (?, ?, ?);")

	_, err := m.DB.ExecContext(ctx, query, imageName, tag, time.Now())
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

	log.Printf("successfully inserted image %s with tag %s.\n", imageName, tag)

	return nil
}

func buildSearchQuery(numerOfArgs int) string {

	sb := strings.Builder{}

	//fmt.Fprintf(&sb, "INSERT INTO photoshare.images (image_name, tag, date_added) VALUES (?, ?, ?);")

	sb.WriteString("SELECT image_name FROM photoshare.images WHERE tag IN(")

	sb.WriteString(strings.Repeat("?,", numerOfArgs-1))
	sb.WriteString("?);")

	return sb.String()

}

// CREATE TABLE `photoshare`.`images` (
//     id int(11) not null primary key auto_increment,
//     imageName varchar(255) not null,
//     title varchar(100),
//     tag varchar(100),
//     date_added datetime not null
// ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
