package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

// WARNING: not thread safe
func (db *DB) GetMedia(media model.Media) (model.Media, error) {
	var res model.Media
	if err := db.db.Where("name = ?", media.Name).
		Where("ext = ?", media.Ext).
		Where("dir = ?", media.Dir).
		Where("core_side_id = ?", media.CoreSideID).
		Where("media_root_id = ?", media.MediaRootID).
		First(&res).
		Error; err != nil {
		return res, err
	}

	return res, nil
}

// WARNING: not thread safe
func (db *DB) GetMediaRequest(request model.MediaRequest) (model.MediaRequest, error) {
	var res model.MediaRequest
	user, err := db.GetUser(request.User)
	if err != nil {
		return res, err
	}
	owner, err := db.GetUser(request.Owner)
	if err != nil {
		return res, err
	}
	media, err := db.GetMedia(request.Media)
	if err != nil {
		return res, err
	}

	if err := db.db.Where("mode = ?", request.Mode).
		Where("user_id = ?", user.ID).
		Where("owner_id = ?", owner.ID).
		Where("media_id = ?", media.ID).
		First(&res).
		Error; err != nil {
		return res, err
	}

	res.User = user
	res.Owner = owner
	res.Media = media

	return res, nil
}

func (db *DB) GetMediaRequestList(owner model.User) (model.MediaRequestList, error) {
	db.Lock()
	defer db.Unlock()

	var res model.MediaRequestList
	owner, err := db.GetUser(owner)
	if err != nil {
		return res, err
	}

	if err := db.db.Joins("LEFT JOIN media_responses ON media_responses.media_request_id = media_requests.id").
		Where("media_responses.id IS NULL").
		Where("media_requests.owner_id = ?", owner.ID).
		Find(&res.Items).
		Error; err != nil {
		return res, err
	}

	if len(res.Items) == 0 {
		return res, gorm.ErrRecordNotFound
	}

	// clear tokens
	owner.Token = ""

	for _, rq := range res.Items {
		rq.Owner = owner
		if err := db.db.Model(&rq).Related(&rq.User).Error; err != nil {
			return res, err
		}
		if err := db.db.Model(&rq).Related(&rq.Media).Error; err != nil {
			return res, err
		}

		// clear tokens
		rq.User.Token = ""
	}

	return res, nil
}

func (db *DB) AddMediaRequest(request model.MediaRequest) error {
	db.Lock()
	defer db.Unlock()

	if _, err := db.GetMediaRequest(request); err == nil {
		return nil
	}

	user, err := db.GetUser(request.User)
	if err != nil {
		return err
	}
	owner, err := db.GetUser(request.Owner)
	if err != nil {
		return err
	}
	media, err := db.GetMedia(request.Media)
	if err != nil {
		return err
	}

	// TODO:
	//  handle "protected" mode

	if request.Mode == model.MediaAccessTypePrivate && user.ID != owner.ID {
		return fmt.Errorf("media owner %s: access denied for user %s", owner.Username, user.Username)
	}

	// update media timestamp
	if err := db.db.Save(&media).Error; err != nil {
		return err
	}

	request.User = model.User{}
	request.Owner = model.User{}
	request.Media = model.Media{}
	request.UserID = user.ID
	request.OwnerID = owner.ID
	request.MediaID = media.ID

	return db.db.Create(&request).Error
}

func (db *DB) GetMediaResponseList(user model.User) (model.MediaResponseList, error) {
	db.Lock()
	defer db.Unlock()

	var res model.MediaResponseList
	user, err := db.GetUser(user)
	if err != nil {
		return res, err
	}
	var requestList model.MediaRequestList
	if err := db.db.Joins("INNER JOIN media_responses ON media_responses.media_request_id = media_requests.id").
		Where("media_requests.user_id = ?", user.ID).
		Find(&requestList.Items).
		Error; err != nil {
		return res, err
	}

	// clear tokens
	user.Token = ""

	for _, rq := range requestList.Items {
		rq.User = user
		if err := db.db.First(&rq.Owner, rq.OwnerID).Error; err != nil {
			return res, err
		}
		if err := db.db.Model(&rq).Related(&rq.Media).Error; err != nil {
			return res, err
		}

		// clear tokens
		rq.Owner.Token = ""

		var rs model.MediaResponse
		if err := db.db.Model(&rq).Related(&rs).Error; err != nil {
			return res, err
		}
		rs.MediaRequest = *rq

		res.Items = append(res.Items, &rs)
	}

	return res, nil
}

func (db *DB) AddMediaResponse(response model.MediaResponse) error {
	db.Lock()
	defer db.Unlock()

	request, err := db.GetMediaRequest(response.MediaRequest)
	if err != nil {
		return err
	}

	if err := db.db.First(&model.MediaResponse{}, "media_request_id = ?", request.ID).Error; err == nil {
		return nil
	}

	request.User = model.User{}
	request.Owner = model.User{}

	response.MediaRequest = request
	response.MediaRequestID = request.ID

	return db.db.Create(&response).Error
}

func (db *DB) IsMediaItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
