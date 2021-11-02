package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Booking struct {
	Id        uuid.UUID `gorm:"PRIMARY_KEY;" json:"Id"`
	Status    string    `gorm:"default:'Booking'"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	b.Id = uuid.NewV4()
	return nil
}

func (Booking) TableName() string {
	return "booking"
}

func (b *Booking) AfterUpdate(tx *gorm.DB) (err error) {
	booking := Booking{}
	err = tx.First(&booking, "id = ?", b.Id).Error
	if err != nil {
		return nil
	}
	b.Id = booking.Id
	return nil
}
