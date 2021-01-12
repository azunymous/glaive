package server

import "glaive/board"

type Board struct {
	Host   string `json:"host"`
	Images string `json:"images"`
}

type boardResponse struct {
	Status string       `json:"status"`
	No     string       `json:"no"`
	Thread board.Thread `json:"thread"`
	Type   string       `json:"type"`
}
