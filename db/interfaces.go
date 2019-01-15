package db

import "github.com/cnaize/mz-center/model"

const (
	MaxItemsCount = 100
)

type DB interface {
	SearchProvider
}

type SearchProvider interface {
	GetSearchRequestList(offset, count int) (model.SearchRequestList, error)
	AddSearchRequest(request model.SearchRequest) error

	GetSearchResponseList(request model.SearchRequest, offset, count int) (model.SearchResponseList, error)
	AddSearchResponseList(user model.User, request model.SearchRequest, response model.SearchResponseList) error
}
