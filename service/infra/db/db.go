package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	User     string
	Pass     string
	Endpoint string
	Name     string
}

func connectDB(config *Config) (*sql.DB, error) {
	connectStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s",
		config.User,
		config.Pass,
		config.Endpoint,
		"3306",
		config.Name,
		"utf8",
	)
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
