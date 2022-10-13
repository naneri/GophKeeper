package user

import (
	"errors"
	"gorm.io/gorm"
)

type DatabaseRepository struct {
	Db *gorm.DB
}

func (d DatabaseRepository) GetByLogin(login string) (User, error) {
	var user User

	result := d.Db.Where("login = ?", login).Find(&user)

	return user, result.Error
}

func (d DatabaseRepository) Get(id uint32) (User, error) {
	panic("implement me")
}

func (d DatabaseRepository) Store(login, password string) (id uint32, err error) {
	var existingUser User
	checkUserResult := d.Db.Where("login = ?", login).Find(&existingUser)

	if checkUserResult.RowsAffected > 0 {
		return 0, errors.New("user with this login already exists")
	}

	user := User{
		Login:    login,
		Password: password,
	}

	result := d.Db.Create(&user)

	return user.Id, result.Error
}
