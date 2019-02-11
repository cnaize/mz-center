package sqlite

import (
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

// TODO:
// Move it to Postgres

func (db *DB) GetUser(user model.User) (model.User, error) {
	var res model.User
	if err := db.db.First(&res, "username = ?", user.Username).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (db *DB) IsUserItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
