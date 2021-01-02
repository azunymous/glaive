package domain

import "glaive/board"

type userResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
	Error    string `json:"error"`
	Token    string `json:"token"`
}

type boardResponse struct {
	Status string       `json:"status"`
	No     string       `json:"no"`
	Thread board.Thread `json:"thread"`
	Type   string       `json:"type"`
}
