package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID `gorm:"PRIMARY_KEY;" json:"Id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.NewV4()
	return nil
}

func (User) TableName() string {
	return "users"
}

func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
	user := User{}
	err = tx.First(&user, "id = ?", u.Id).Error
	if err != nil {
		return nil
	}
	u.Id = user.Id
	return nil
}
