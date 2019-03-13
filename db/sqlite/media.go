package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

func (db *DB) GetMediaRequestList(owner model.User) (model.MediaRequestList, error) {
	var res model.MediaRequestList
	if err := db.db.Find(&res.Items, "owner = ?", owner.Username).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) AddMediaRequest(request model.MediaRequest) error {
	if err := db.db.First(&model.MediaRequest{}, "user = ? AND owner = ? AND media_id = ?",
		request.User, request.Owner, request.MediaID).Error; err == nil {
		return fmt.Errorf("request already exists")
	}

	return db.db.Create(&request).Error
}

func (db *DB) GetMediaResponseList(user model.User) (model.MediaResponseList, error) {
	var res model.MediaResponseList
	if err := db.db.Find(&res.Items, "user = ?", user.Username).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) AddMediaResponse(response model.MediaResponse) error {
	if err := db.db.First(&model.MediaResponse{}, "user = ? AND owner = ? AND media_id = ?",
		response.User, response.Owner, response.MediaID).Error; err == nil {
		return fmt.Errorf("response already exists")
	}

	return db.db.Create(&response).Error
}

func (db *DB) IsMediaItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
