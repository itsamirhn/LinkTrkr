package models

import "net/http"

type Redirect struct {
	URL    string `json:"url"`
	ChatID int64  `json:"chat_id"`
}

type Click struct {
	Redirect *Redirect
	Request  *http.Request
}
