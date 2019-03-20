package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

// WARNING: not thread safe
func (db *DB) GetMedia(media model.Media) (model.Media, error) {
	var res model.Media
	if err := db.db.First(&res, "name = ? AND ext = ? AND dir = ? AND core_side_id = ? AND media_root_id = ?",
		media.Name, media.Ext, media.Dir, media.CoreSideID, media.MediaRootID).
		Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) GetMediaRequestList(owner model.User) (model.MediaRequestList, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var res model.MediaRequestList
	owner, err := db.GetUser(owner)
	if err != nil {
		return res, err
	}

	if err := db.db.Find(&res.Items, "owner_id = ?", owner.ID).Error; err != nil {
		return res, err
	}

	for _, r := range res.Items {
		r.Owner = owner
		if err := db.db.Model(&r).Related(&r.User).Error; err != nil {
			return res, err
		}
		if err := db.db.Model(&r).Related(&r.Media).Error; err != nil {
			return res, err
		}
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

	// TODO:
	//  handle "protected" mode

	request.User = model.User{}
	request.Owner = model.User{}
	request.UserID = user.ID
	request.OwnerID = owner.ID

	if media, err := db.GetMedia(request.Media); err == nil {
		request.Media = model.Media{}
		request.MediaID = media.ID
	}

	return db.db.Create(&request).Error
}

func (db *DB) GetMediaResponseList(user model.User) (model.MediaResponseList, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var res model.MediaResponseList
	user, err := db.GetUser(user)
	if err != nil {
		return res, err
	}

	if err := db.db.Find(&res.Items, "user_id = ?", user.ID).Error; err != nil {
		return res, err
	}

	for _, r := range res.Items {
		r.User = user
		if err := db.db.First(&r.Owner, r.OwnerID).Error; err != nil {
			return res, err
		}
		if err := db.db.Model(&r).Related(&r.Media).Error; err != nil {
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
	media, err := db.GetMedia(response.Media)
	if err != nil {
		return err
	}

	// TODO:
	//  handle "protected" mode

	if response.Mode == model.MediaAccessTypePrivate {
		if err := db.db.First(&model.MediaRequest{},
			"user_id = ? AND owner_id = ? AND media_id = ? AND mode = ?",
			user.ID, owner.ID, media.ID, response.Mode).
			Error; err != nil {
			return fmt.Errorf("media request not found: user %s, owner %s, media id %d",
				user.Username, owner.Username, media.ID)
		}
	}

	// update media timestamp
	if err := db.db.Save(&response.Media).Error; err != nil {
		return err
	}

	response.User = model.User{}
	response.Owner = model.User{}
	response.Media = model.Media{}
	response.UserID = user.ID
	response.OwnerID = owner.ID
	response.MediaID = media.ID

	return db.db.Create(&response).Error
}

func (db *DB) IsMediaItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
