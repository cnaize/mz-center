package model

type MediaExt string

const (
	MediaExtUnknown MediaExt = "unknown"
	MediaExtMP3     MediaExt = "mp3"
)

var SupportedMediaExtList = []MediaExt{MediaExtMP3}

type MediaAccessType string

const (
	MediaAccessTypeUnknown   MediaAccessType = "unknown"
	MediaAccessTypePrivate   MediaAccessType = "private"
	MediaAccessTypeProtected MediaAccessType = "protected"
	MediaAccessTypePublic    MediaAccessType = "public"
)

type Media struct {
	Base
	Name        string   `json:"name" form:"name" binding:"required"`
	Ext         MediaExt `json:"ext" form:"ext" binding:"required"`
	Dir         string   `json:"dir" form:"dir" binding:"required"`
	Error       *Error   `json:"error,omitempty" sql:"-" form:"error"`
	RawPath     string   `json:"-"` // lower path, used for search
	MediaRootID uint     `json:"-" gorm:"not_null"`
}

type MediaList struct {
	Items         []*Media `json:"items" form:"items"`
	AllItemsCount *uint    `json:"allItemsCount,omitempty" form:"allItemsCount"`
	Error         *Error   `json:"error,omitempty" form:"error"`
}

type MediaRoot struct {
	Base
	Dir        string          `json:"dir" gorm:"unique_index" form:"dir" binding:"required"`
	AccessType MediaAccessType `json:"accessType" form:"accessType" binding:"required"`
	MediaCount *uint           `json:"mediaCount,omitempty" form:"mediaCount"`
	Error      *Error          `json:"error,omitempty" sql:"-" form:"error"`
}

type MediaRootList struct {
	Items []*MediaRoot `json:"items"`
	Error *Error       `json:"error,omitempty"`
}
