package sqlite

import (
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

func (db *DB) GetSearchRequest(request model.SearchRequest) (model.SearchRequest, error) {
	var res model.SearchRequest
	if err := db.db.First(&res, "text = ?", request.Text).Error; err != nil {
		return res, err
	}

	return res, nil
}

// TODO:
//  check if the user already sent response for the request
func (db *DB) GetSearchRequestList(user model.User, offset, count uint) (model.SearchRequestList, error) {
	var res model.SearchRequestList
	user, err := db.GetUser(user)
	if err != nil {
		return res, err
	}

	if err := db.db.Find(&res.Items).Offset(offset).Limit(count).Error; err != nil {
		return res, err
	}

	db.db.Model(&model.SearchRequest{}).Count(&res.AllItemsCount)

	return res, nil
}

func (db *DB) AddSearchRequest(request model.SearchRequest) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.db.Save(&request).Error
}

func (db *DB) GetSearchResponseList(request model.SearchRequest, offset, count uint) (model.SearchResponseList, error) {
	var res model.SearchResponseList
	request, err := db.GetSearchRequest(request)
	if err != nil {
		return res, err
	}

	if err := db.db.Model(&request).Related(&res.Items).Offset(offset).Limit(count).Error; err != nil {
		return res, err
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

func (db *DB) AddSearchResponseList(user model.User, request model.SearchRequest, responseList model.SearchResponseList) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	user, err := db.GetUser(user)
	if err != nil {
		return err
	}
	request, err = db.GetSearchRequest(request)
	if err != nil {
		return err
	}

	if !db.db.Model(&request).Related(&model.SearchResponse{}).Limit(1).RecordNotFound() {
		return nil
	}

	tx := db.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, r := range responseList.Items {
		// remove id to prevent creation media with core-side id
		r.Media.ID = 0

		r.UserID = user.ID
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
