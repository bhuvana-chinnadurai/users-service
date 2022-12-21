package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `gorm:"primaryKey" json:"id"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	Nickname     string    `json:"nickname"`
	PasswordHash string    `gorm:"not null"`
	Email        string    `gorm:"index:idx_email_country,unique;not null;" json:"email"`
	Country      string    `gorm:"index:idx_email_country;not null" json:"country"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Filter struct {
	Country string `json:"country"`
}
type Pagination struct {
	Limit int `json:"int"`
	Page  int `json:"page"`
}

var ErrDuplicate = fmt.Errorf("user with the email already exists")
var ErrNotFound = fmt.Errorf("user with the id does not exist")

type Cursor interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}
