package models

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Theatre{},
		&Screen{},
		&Seat{},
		&Movie{},
		&Show{},
		&Booking{},
		&BookingSeat{},
		&Payment{},
	)
}
