package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

// WARNING: not thread safe
func (db *DB) GetSearchRequest(request model.SearchRequest) (model.SearchRequest, error) {
	var res model.SearchRequest
	if err := db.db.First(&res, "text = ? AND mode = ?", request.Text, request.Mode).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) GetSearchRequestList(offset, count uint) (model.SearchRequestList, error) {
	db.Lock()
	defer db.Unlock()

	var res model.SearchRequestList
	query := db.db.Joins("LEFT JOIN search_responses ON search_responses.search_request_id = search_requests.id").
		Where("search_responses.search_request_id IS NULL").
		Joins("LEFT JOIN users ON users.id = search_responses.user_id").
		Where("users.id IS NULL").
		Offset(offset).
		Limit(count)

	// TODO:
	//  handle "protected" mode

	var private model.SearchRequestList
	if err := query.Where("search_requests.mode = ?", model.MediaAccessTypePrivate).
		Find(&private.Items).
		Error; err != nil {
		return res, err
	}

	var public model.SearchRequestList
	if uint(len(private.Items)) < count {
		if err := query.Where("search_requests.mode = ?", model.MediaAccessTypePublic).
			Find(&public.Items).
			Error; err != nil {
			return res, err
		}
	}

	res.Items = append(private.Items, public.Items...)
	if len(res.Items) == 0 {
		return res, gorm.ErrRecordNotFound
	}

	return res, nil
}

func (db *DB) AddSearchRequest(user model.User, request model.SearchRequest) error {
	db.Lock()
	defer db.Unlock()

	if request.Mode == model.MediaAccessTypeProtected || request.Mode == model.MediaAccessTypePrivate {
		user, err := db.GetUser(user)
		if err != nil {
			return err
		}

		request.UserID = user.ID
	}

	if request, err := db.GetSearchRequest(request); err == nil {
		return db.db.Save(&request).Error
	}

	return db.db.Create(&request).Error
}

func (db *DB) GetSearchResponseList(user model.User, request model.SearchRequest, offset, count uint) (model.SearchResponseList, error) {
	db.Lock()
	defer db.Unlock()

	var res model.SearchResponseList
	request, err := db.GetSearchRequest(request)
	if err != nil {
		return res, err
	}

	query := db.db.Joins("INNER JOIN search_requests ON search_requests.id = search_responses.search_request_id").
		Where("search_requests.mode = ?", request.Mode).
		Offset(offset).
		Limit(count)

	// TODO:
	//  handle "protected" mode

	if request.Mode == model.MediaAccessTypePrivate {
		user, err := db.GetUser(user)
		if err != nil {
			return res, err
		}

		query = query.Where("search_requests.user_id = ?", user.ID)
	}

	if err := query.Find(&res.Items).Error; err != nil {
		if !db.IsSearchItemNotFound(err) {
			return res, err
		}
	}

	for _, r := range res.Items {
		if err := db.db.Model(&r).Related(&r.Owner).Error; err != nil {
			return res, err
		}
		if err := db.db.Model(&r).Related(&r.Media).Error; err != nil {
			return res, err
		}
	}

	return res, nil
}

func (db *DB) AddSearchResponseList(owner model.User, request model.SearchRequest, responseList model.SearchResponseList) error {
	db.Lock()
	defer db.Unlock()

	owner, err := db.GetUser(owner)
	if err != nil {
		return err
	}
	request, err = db.GetSearchRequest(request)
	if err != nil {
		return err
	}

	if db.db.Joins("INNER JOIN search_requests ON search_requests.id = search_responses.search_request_id").
		Where("search_requests.mode = ?", request.Mode).
		Where("search_responses.user_id = ?", owner.ID).
		First(&model.SearchResponse{}).
		RowsAffected > 0 {
		return nil
	}

	// TODO:
	//  handle "protected" mode

	if request.Mode == model.MediaAccessTypePrivate {
		if request.UserID != owner.ID {
			return fmt.Errorf("user and owner mismatch in \"private\" mode: %d != %d (%s)",
				request.UserID, owner.ID, owner.Username)
		}
	}

	tx := db.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, r := range responseList.Items {
		r.UserID = owner.ID
		r.SearchRequestID = request.ID

		if err := tx.Create(&r).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (db *DB) IsSearchItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
