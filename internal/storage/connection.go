package storage

import (
	"database/sql"
	"log"
	
	_ "github.com/mattn/go-sqlite3"
)

const dbfile = "./mygopher.db"

func Run() *sql.DB {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal("failure when opening db connection: ", err.Error())
	}

	return db
}