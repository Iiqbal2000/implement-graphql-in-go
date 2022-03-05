package platform

import (
	_ "embed"
	"encoding/json"
	"net/http"

	"github.com/Iiqbal2000/mygopher/internal/links"
	"github.com/Iiqbal2000/mygopher/internal/users"
	"github.com/graph-gophers/graphql-go"
)

//go:embed schema.graphql
var schema []byte

type Server struct {
	Schema  *graphql.Schema
	Loaders Loaders
}

func CreateServer(uSvc users.Service, lSvc links.Service) Server {
	resolver := &Resolver{
		UserSvc: uSvc,
		LinkSvc: lSvc,
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}

	g := Server{
		Schema:  graphql.MustParseSchema(readSchema(), resolver, opts...),
		Loaders: InitLoaders(uSvc),
	}

	return g
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type QueryIn struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	var params QueryIn

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := s.Loaders.Attach(r.Context())

	response := s.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	// fmt.Println(response)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func readSchema() string {
	return string(schema)
}
