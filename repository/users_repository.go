package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/bhuvana-chinnadurai/users-service/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db: db}
}

func (u *Users) Create(user *model.User) (uuid.UUID, error) {
	user.Id = uuid.New()
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	result := u.db.Create(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return uuid.Nil, model.ErrDuplicate
		}
		return uuid.Nil, fmt.Errorf("error while creating the user: %s", result.Error.Error())
	}
	return user.Id, result.Error
}

func (u *Users) Get(id uuid.UUID) (*model.User, error) {

	user := &model.User{}
	result := u.db.First(&user).Where("id = ?", id.String())
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("error while finding the user: %s", result.Error.Error())
	}
	return user, nil
}

func (u *Users) Update(user *model.User) (uuid.UUID, error) {
	user.UpdatedAt = time.Now()

	result := u.db.Model(&user).Omit("Id,CreatedAt,Country").Updates(user)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			return uuid.Nil, model.ErrDuplicate
		}
		return uuid.Nil, fmt.Errorf("error while updating the user: %s", result.Error.Error())
	}
	return user.Id, nil
}

func (u *Users) Delete(id uuid.UUID) error {
	result := u.db.Where("id = ?", id.String()).Delete(&model.User{})
	if result.Error != nil {
		return fmt.Errorf("error while deleting the user: %s", result.Error.Error())
	}
	return nil
}

func (u *Users) GetAll(filter model.Filter, pagination model.Pagination) ([]model.User, error) {
	offset := pagination.Page - 1*pagination.Limit

	rows, err := u.db.Model(&model.User{}).Where(filter).Limit(pagination.Limit).Offset(offset).Rows()
	if err != nil {
		return nil, fmt.Errorf("error while querying users: %s", err.Error())
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := u.db.ScanRows(rows, &user)
		if err != nil {
			return nil, fmt.Errorf("error while scanning the next user: %s", err.Error())
		}
		users = append(users, user)
	}

	return users, nil
}
