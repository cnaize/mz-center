package sqlite

import (
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

// TODO:
// Move it to Postgres

// WARNING: not thread safe
func (db *DB) GetUser(user model.User) (model.User, error) {
	var res model.User
	if err := db.db.First(&res, "username = ?", user.Username).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) CreateUser(user model.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.db.Create(&user).Error
}

func (db *DB) IsUserItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
