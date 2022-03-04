package mygopher

// Error represents custom error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements error interface
func (e Error) Error() string {
	return e.Message
}

// Extensions implements ResolveError interface
// see: https://github.com/graph-gophers/graphql-go#custom-errors
func (e Error) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}
