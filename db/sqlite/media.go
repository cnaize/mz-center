package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

func (db *DB) GetMediaRequestList(owner model.User) (model.MediaRequestList, error) {
	var res model.MediaRequestList
	owner, err := db.GetUser(owner)
	if err != nil {
		return res, err
	}

	if err := db.db.Find(&res.Items, "owner_id = ?", owner.ID).Error; err != nil {
		return res, err
	}

	for _, r := range res.Items {
		if r.OwnerID != owner.ID {
			return res, fmt.Errorf("request owner id and owner id mismatch: %d != %d", r.OwnerID, owner.ID)
		}

		if err := db.db.First(&r.User, r.UserID).Error; err != nil {
			return res, err
		}

		r.Owner = owner
	}

	return res, nil
}

func (db *DB) AddMediaRequest(user model.User, request model.MediaRequest) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	user, err := db.GetUser(user)
	if err != nil {
		return err
	}
	owner, err := db.GetUser(request.Owner)
	if err != nil {
		return err
	}

	request.User = model.User{}
	request.Owner = model.User{}
	request.UserID = user.ID
	request.OwnerID = owner.ID

	return db.db.Create(&request).Error
}

func (db *DB) GetMediaResponseList(user model.User) (model.MediaResponseList, error) {
	var res model.MediaResponseList
	user, err := db.GetUser(user)
	if err != nil {
		return res, err
	}

	if err := db.db.Find(&res.Items, "user_id = ?", user.ID).Error; err != nil {
		return res, err
	}

	for _, r := range res.Items {
		if r.UserID != user.ID {
			return res, fmt.Errorf("response user id and user id mismatch: %d != %d", r.UserID, user.ID)
		}

		if err := db.db.First(&r.Owner, r.OwnerID).Error; err != nil {
			return res, err
		}
		if err := db.db.First(&r.Media, r.MediaID).Error; err != nil {
			return res, err
		}

		r.User = user
	}

	return res, nil
}

func (db *DB) AddMediaResponse(owner model.User, response model.MediaResponse) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	owner, err := db.GetUser(owner)
	if err != nil {
		return err
	}
	user, err := db.GetUser(response.User)
	if err != nil {
		return err
	}

	// remove id to prevent creation media with core-side id
	response.Media.ID = 0

	response.User = model.User{}
	response.Owner = model.User{}
	response.UserID = user.ID
	response.OwnerID = owner.ID

	return db.db.Create(&response).Error
}

func (db *DB) IsMediaItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
