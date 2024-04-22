package db

import (
	"awesomeProject/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Config struct {
	Login    string
	Password string
	Port     int
}

var db *sql.DB = nil

func InitDb(config *config.DBConfig) (*sql.DB, error) {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		config.Host, config.Port, config.User, config.Pass, config.Name))

	if err != nil {
		log.Println(err)
		log.Println("can't connect do DB")

		return nil, err
	}

	ping := db.Ping()

	if ping != nil {
		db.Close()
		log.Println("can't ping DB")

		return nil, ping
	}

	log.Println("DB ping OK")
	return db, nil
}

func GetDB() *sql.DB {
	return db
}
