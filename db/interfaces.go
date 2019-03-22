package db

import (
	"github.com/cnaize/mz-common/model"
)

type DB interface {
	UserProvider
	MediaProvider
	SearchProvider
}

type UserProvider interface {
	GetUser(user model.User) (model.User, error)
	CreateUser(user model.User) error

	IsUserItemNotFound(err error) bool
}

type MediaProvider interface {
	GetMedia(media model.Media) (model.Media, error)

	GetMediaRequest(request model.MediaRequest) (model.MediaRequest, error)
	GetMediaRequestList(owner model.User) (model.MediaRequestList, error)
	AddMediaRequest(request model.MediaRequest) error

	GetMediaResponseList(user model.User) (model.MediaResponseList, error)
	AddMediaResponse(response model.MediaResponse) error

	IsMediaItemNotFound(err error) bool
}

type SearchProvider interface {
	GetSearchRequest(request model.SearchRequest) (model.SearchRequest, error)
	GetSearchRequestList(offset, count uint) (model.SearchRequestList, error)
	AddSearchRequest(user model.User, request model.SearchRequest) error

	GetSearchResponseList(user model.User, request model.SearchRequest, offset, count uint) (model.SearchResponseList, error)
	AddSearchResponseList(owner model.User, request model.SearchRequest, responseList model.SearchResponseList) error

	IsSearchItemNotFound(err error) bool
}
