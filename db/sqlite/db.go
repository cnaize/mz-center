package sqlite

import (
	"fmt"
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DB struct {
	db *gorm.DB
}

func New() (*DB, error) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("open failed: %+v", err)
	}

	if err := prepare(db); err != nil {
		return nil, fmt.Errorf("prepare failed: %+v", err)
	}

	return &DB{
		db: db,
	}, nil
}

func (db *DB) GetSearchRequestList(offset, count uint) (model.SearchRequestList, error) {
	var res model.SearchRequestList
	if err := db.db.Offset(offset).Limit(count).Find(&res.Items).Error; err != nil {
		return res, err
	}

	db.db.Model(&model.SearchRequest{}).Count(&res.AllItemsCount)

	return res, nil
}

func (db *DB) AddSearchRequest(request model.SearchRequest) error {
	return db.db.Create(&request).Error
}

func (db *DB) GetSearchResponseList(request model.SearchRequest, offset, count uint) (model.SearchResponseList, error) {
	var res model.SearchResponseList
	if err := db.db.Model(&request).Related(&res.Items).Offset(offset).Limit(count).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) AddSearchResponseList(user model.User, request model.SearchRequest, responseList model.SearchResponseList) error {
	for _, r := range responseList.Items {
		r.Owner = &user
		r.SearchRequestID = request.ID
	}

	if err := db.db.Create(responseList.Items).Error; err != nil {
		return err
	}

	return nil
}

func (db *DB) IsSearchItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

func prepare(db *gorm.DB) error {
	db.LogMode(true)

	if err := db.AutoMigrate(&model.SearchRequest{}).Error; err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.SearchResponse{}).Error; err != nil {
		return err
	}

	return nil
}
