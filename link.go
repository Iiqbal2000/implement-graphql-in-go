package mygopher

type Link struct {
	ID      string
	Title   string
	Address string
	UserID  string
}

type LinkIn struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	UserID  string `json:"user_id"`
}

type LinkOut struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	UserID  string `json:"user_id"`
}
