package db

import (
	"github.com/cnaize/mz-common/model"
	"time"
)

const (
	MaxRequestItemsCount  = 100
	MaxResponseItemsCount = 100
	RequestsCleanPeriod   = 8 * time.Second
)

type DB interface {
	SearchProvider
}

type SearchProvider interface {
	GetSearchRequestList(offset, count uint) (*model.SearchRequestList, error)
	AddSearchRequest(request *model.SearchRequest) (*model.SearchRequest, error)

	GetSearchResponseList(request *model.SearchRequest, offset, count uint) (*model.SearchResponseList, error)
	AddSearchResponseList(user *model.User, request *model.SearchRequest, response *model.SearchResponseList) error
}
