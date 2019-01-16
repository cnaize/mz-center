package model

import "time"

type SearchRequest struct {
	Text           string    `json:"text"`
	TotalRespCount uint      `json:"totalRespCount"`
	CreatedAt      time.Time `json:"-"`
	Error          *Error    `json:"error,omitempty"`
}

type SearchRequestList struct {
	Items         []*SearchRequest `json:"items"`
	TotalReqCount uint             `json:"totalReqCount"`
	Error         *Error           `json:"error,omitempty"`
}

type SearchResponse struct {
	Owner    *User  `json:"owner"`
	Author   string `json:"author"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
	Error    *Error `json:"error,omitempty"`
}

type SearchResponseList struct {
	Request *SearchRequest    `json:"-"`
	Items   []*SearchResponse `json:"items"`
	Error   *Error            `json:"error,omitempty"`
}
