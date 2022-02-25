package mygopher

type User struct {
	ID       string
	Username string
	Password string
}

type UserIn struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserOut struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
