package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Iiqbal2000/mygopher/internal/links"
	// graphql-platform
	graphqlP "github.com/Iiqbal2000/mygopher/internal/platform/graphql"
	// http-platform
	httpP "github.com/Iiqbal2000/mygopher/internal/platform/httphandler"
	authP "github.com/Iiqbal2000/mygopher/internal/platform/auth"
	"github.com/Iiqbal2000/mygopher/internal/storage"
	"github.com/Iiqbal2000/mygopher/internal/users"
	"github.com/graph-gophers/graphql-go"
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
		Db:  db,
		Log: log,
		UserSvc: userSvc,
	}

	auth := authP.Auth{
		UserSvc: userSvc,
	}

	resolver := &graphqlP.Resolver{
		UserSvc: userSvc,
		LinkSvc: linkSvc,
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	// read schema
	s := graphqlP.Read()

	schema := graphql.MustParseSchema(s, resolver, opts...)

	graphqlHandler := &httpP.GraphqlHandler{
		Schema: schema,
		Loaders: graphqlP.InitLoaders(userSvc),
	}

	http.Handle("/login", httpP.Login(auth))
	http.Handle("/query", httpP.AuthorizeMiddleware(graphqlHandler, auth))

	log.Print("Server is lintening at localhost:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
