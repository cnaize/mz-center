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
	RawPath     string   `json:"-"` // lower path, used for search
	MediaRootID uint     `json:"-" gorm:"not_null"`
}

type MediaList struct {
	Items         []*Media `json:"items" form:"items"`
	AllItemsCount *uint    `json:"allItemsCount,omitempty" form:"allItemsCount"`
}

type MediaRoot struct {
	Base
	Dir        string          `json:"dir" gorm:"unique_index" form:"dir" binding:"required"`
	AccessType MediaAccessType `json:"accessType" form:"accessType" binding:"required"`
	MediaCount *uint           `json:"mediaCount,omitempty" form:"mediaCount"`
}

type MediaRootList struct {
	Items []*MediaRoot `json:"items"`
}

type MediaRequest struct {
	User      User   `json:"user"`
	Owner     User   `json:"owner" form:"owner" binding:"required,dive"`
	MediaID   uint   `json:"mediaID" form:"mediaID" binding:"required"`
	WebRTCKey string `json:"webRTCKey" form:"webRTCKey" binding:"required"`
}

type MediaRequestList struct {
	Items []*MediaRequest `json:"items" form:"items" binding:"required,dive"`
}

type MediaResponse struct {
	User      User   `json:"user" form:"user" binding:"required,dive"`
	Owner     User   `json:"owner"`
	Media     Media  `json:"media" form:"media" binding:"required,dive"`
	WebRTCKey string `json:"webRTCKey" form:"webRTCKey" binding:"required"`
}

type MediaResponseList struct {
	Items []*MediaResponse `json:"items" form:"items" binding:"required,dive"`
}
