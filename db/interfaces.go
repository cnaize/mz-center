package db

import (
	"github.com/cnaize/mz-common/model"
)

type DB interface {
	UserProvider
	SearchProvider
}

type UserProvider interface {
	GetUser(user model.User) (model.User, error)
	CreateUser(user model.User) error

	IsUserItemNotFound(err error) bool
}

type SearchProvider interface {
	GetSearchRequestList(offset, count uint) (model.SearchRequestList, error)
	AddSearchRequest(request model.SearchRequest) error

	GetSearchResponseList(request model.SearchRequest, offset, count uint) (model.SearchResponseList, error)
	AddSearchResponseList(user model.User, request model.SearchRequest, responseList model.SearchResponseList) error

	IsSearchItemNotFound(err error) bool
}
