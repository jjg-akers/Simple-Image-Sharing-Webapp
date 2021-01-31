package build

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/urfave/cli/v2"
)

type SQLDBConfig struct {
	UserName     string
	Password     string
	Host         string
	Port         string
	MaxOpenConns int
	MaxIdleConns int
}

func NewSQLDBConfig() *SQLDBConfig {
	return &SQLDBConfig{}
}

func LoadSQLDBConfig(config *SQLDBConfig) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "mysql-user",
			EnvVars:     []string{"MYSQL_USER"},
			Destination: &config.UserName,
		},
		&cli.StringFlag{
			Name:        "mysql-passwor",
			EnvVars:     []string{"MYSQL_PASSWORD"},
			Destination: &config.Password,
		},
		&cli.StringFlag{
			Name:        "mysql-host",
			EnvVars:     []string{"MYSQL_HOST"},
			Destination: &config.Host,
		},
		&cli.StringFlag{
			Name:        "mysql-port",
			EnvVars:     []string{"MYSQL_PORT"},
			Destination: &config.Port,
		},
		&cli.IntFlag{
			Name:        "mysql-max-open-conns",
			EnvVars:     []string{"MYSQL_MAX_OPEN_CONNS"},
			Destination: &config.MaxOpenConns,
			Value:       20,
		},
		&cli.IntFlag{
			Name:        "mysql-max-idle-conns",
			EnvVars:     []string{"MYSQL_MAX_IDLE_CONNS"},
			Destination: &config.MaxIdleConns,
			Value:       20,
		},
	}
}

func NewSQLDB(config *SQLDBConfig) (*sql.DB, error) {

	dsn := config.UserName + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("error opening DB: ", err)
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	if err := db.Ping(); err != nil {
		log.Println("error ping DB: ", err)
		return nil, err
	}

	return db, nil
}
