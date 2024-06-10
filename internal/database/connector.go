package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // enforce loading PG driver
	"github.com/pro-posal/webserver/config"
)

type DBConnector struct {
	Conn *sql.DB
}

func Connect() *DBConnector {
	connstr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Name,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.SSLMode,
	)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Could not connect to the database! %v", err)
	}

	return &DBConnector{
		Conn: db,
	}
}

func TestConnect() *DBConnector {
	connstr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		config.TestConfig.TestDatabase.UserName,
		config.TestConfig.TestDatabase.Password,
		config.TestConfig.TestDatabase.DbName,
		config.TestConfig.TestDatabase.Host,
		config.TestConfig.TestDatabase.Port,
		config.TestConfig.TestDatabase.SSLMode,
	)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Could not connect to the database! %v", err)
	}

	return &DBConnector{
		Conn: db,
	}
}
