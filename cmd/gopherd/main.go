package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Iiqbal2000/mygopher/internal/links"
	// graphql-platform
	graphqlP "github.com/Iiqbal2000/mygopher/internal/platform/graphql"
	// http-platform
	authP "github.com/Iiqbal2000/mygopher/internal/platform/auth"
	"github.com/Iiqbal2000/mygopher/internal/storage"
	"github.com/Iiqbal2000/mygopher/internal/users"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log := log.New(os.Stdout, "GOPHER: ", log.Lshortfile)

	db := storage.Run()

	storage.Migrate(db)

	userSvc := users.UserService{
		Db:  db,
		Log: log,
	}

	linkSvc := links.LinkService{
		Db:      db,
		Log:     log,
		UserSvc: userSvc,
	}

	auth := authP.Auth{
		UserSvc: userSvc,
	}

	graphqlServer := graphqlP.NewGraphQl(userSvc, linkSvc)

	http.Handle("/login", authP.Login(auth))
	http.Handle("/query", authP.AuthorizeMiddleware(graphqlServer, auth))

	log.Print("Server is lintening at localhost:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
