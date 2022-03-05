package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Iiqbal2000/mygopher/internal/links"
	"github.com/Iiqbal2000/mygopher/internal/users"
	graphql "github.com/Iiqbal2000/mygopher/internal/platform/graphql"
	auth "github.com/Iiqbal2000/mygopher/internal/platform/auth"
	sqlite3 "github.com/Iiqbal2000/mygopher/internal/platform/sqlite3"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log := log.New(os.Stdout, "GOPHER: ", log.Lshortfile)

	db := sqlite3.Run()

	sqlite3.Migrate(db)

	userSvc := users.Service{
		Db:  db,
		Log: log,
	}

	linkSvc := links.Service{
		Db:      db,
		Log:     log,
		UserSvc: userSvc,
	}

	authSvc := auth.Service{
		UserSvc: userSvc,
	}

	graphqlServer := graphql.CreateServer(userSvc, linkSvc)
	
	http.Handle("/login", auth.Login(authSvc))
	http.Handle("/query", auth.Authorize(graphqlServer, authSvc))

	log.Print("Server is lintening at localhost:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
