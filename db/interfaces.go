package db

import (
	"github.com/cnaize/mz-common/model"
	"time"
)

const (
	RequestsCleanPeriod = 8 * time.Second
)

type DB interface {
	SearchProvider
}

type SearchProvider interface {
	GetSearchRequestList(offset, count uint) (model.SearchRequestList, error)
	AddSearchRequest(request model.SearchRequest) error

	GetSearchResponseList(request model.SearchRequest, offset, count uint) (model.SearchResponseList, error)
	AddSearchResponseList(user model.User, request model.SearchRequest, responseList model.SearchResponseList) error

	IsSearchItemNotFound(err error) bool
}
