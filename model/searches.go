package model

import "time"

type SearchRequest struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"-"`
}

type SearchRequestList struct {
	Items []SearchRequest `json:"items"`
}

type SearchResponse struct {
	Owner    User   `json:"-"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
}

type SearchResponseList struct {
	Request SearchRequest    `json:"-"`
	Items   []SearchResponse `json:"items"`
}
