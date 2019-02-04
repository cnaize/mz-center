package sqlite

import (
	"github.com/cnaize/mz-common/model"
	"github.com/jinzhu/gorm"
)

func (db *DB) GetUser(user model.User) (model.User, error) {
	if err := db.db.First(&user, "username = ?", user.Username).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (db *DB) GetUserByXToken(xtoken string) (model.User, error) {
	var res model.User
	if err := db.db.First(&res, "x_token = ?", xtoken).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (db *DB) IsUserItemNotFound(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}
