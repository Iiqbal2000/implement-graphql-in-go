package platform

import (
	"encoding/json"
	"net/http"

	graphqlP "github.com/Iiqbal2000/mygopher/internal/platform/graphql"
	"github.com/graph-gophers/graphql-go"
)

type GraphqlHandler struct {
	Schema  *graphql.Schema
	Loaders graphqlP.Loaders
}

func (h *GraphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := h.Loaders.Attach(r.Context())

	response := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
