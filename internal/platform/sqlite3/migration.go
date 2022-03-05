package platform

import (
	"database/sql"
	"log"
)

var userTable = `CREATE TABLE IF NOT EXISTS Users(
	ID INTEGER PRIMARY KEY UNIQUE,
	Username VARCHAR (127) NOT NULL UNIQUE,
	Password VARCHAR (127) NOT NULL
)`

var linkTable = `CREATE TABLE IF NOT EXISTS Links(
	ID INT PRIMARY KEY UNIQUE,
	Title VARCHAR (255),
	Address VARCHAR (255),
	UserID INT,
	FOREIGN KEY (UserID) REFERENCES Users(ID)
)`

func Migrate(DbConn *sql.DB) {
	_, err := DbConn.Exec(userTable)
	if err != nil {
		log.Fatal("failure when migrating user table: ", err.Error())
	}

	_, err = DbConn.Exec(linkTable)
	if err != nil {
		log.Fatal("failure when migrating link table: ", err.Error())
	}
}
