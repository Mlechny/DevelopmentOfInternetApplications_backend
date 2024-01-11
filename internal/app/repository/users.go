package repository

import (
	"errors"
	"gorm.io/gorm"
	"web-service/internal/app/ds"
)

func (r *Repository) AddUser(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{}
	if err := r.db.Where("login = ?", login).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserById(uuid string) (*ds.User, error) {
	user := &ds.User{}
	if err := r.db.Where("uuid = ?", uuid).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
