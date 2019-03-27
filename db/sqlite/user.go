package sqlite

import (
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
	"strings"
)

// TODO:
// Move it to Postgres

// WARNING: not thread safe
func (db *DB) GetUser(user model.User) (model.User, error) {
	var res model.User
	if err := db.db.First(&res, "LOWER(username) = ?", strings.ToLower(user.Username)).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) CreateUser(user model.User) error {
	db.Lock()
	defer db.Unlock()

	return db.db.Create(&user).Error
}

func (db *DB) IsUserItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
