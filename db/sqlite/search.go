package sqlite

import (
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

func (db *DB) GetSearchRequestList(offset, count uint) (model.SearchRequestList, error) {
	var res model.SearchRequestList
	if err := db.db.Offset(offset).Limit(count).Find(&res.Items).Error; err != nil {
		return res, err
	}

	db.db.Model(&model.SearchRequest{}).Count(&res.AllItemsCount)

	return res, nil
}

func (db *DB) AddSearchRequest(request model.SearchRequest) error {
	return db.db.Save(&request).Error
}

func (db *DB) GetSearchResponseList(request model.SearchRequest, offset, count uint) (model.SearchResponseList, error) {
	var req model.SearchRequest
	var res model.SearchResponseList
	if err := db.db.First(&req, "text = ?", request.Text).Error; err != nil {
		return res, err
	}

	if err := db.db.Model(&req).Related(&res.Items).Offset(offset).Limit(count).Error; err != nil {
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

	var req model.SearchRequest
	if err := db.db.First(&req, "text = ?", request.Text).Error; err != nil {
		return err
	}

	if !db.db.Model(&req).Related(&model.SearchResponse{}).Limit(1).RecordNotFound() {
		return nil
	}

	tx := db.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, r := range responseList.Items {
		r.UserID = user.ID
		r.SearchRequestID = req.ID

		if err := tx.Create(&r).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (db *DB) IsSearchItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
