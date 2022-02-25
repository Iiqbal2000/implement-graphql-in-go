package platform

import (
	_ "embed"
)

//go:embed schema.graphql
var schema []byte

func Read() string {
	return string(schema)
}
