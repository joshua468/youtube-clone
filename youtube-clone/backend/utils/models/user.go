package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint32          `gorm:"primary_key;auto_increment" json:"id"`
	Username  string          `gorm:"size:255;not null;unique" json:"username" validate:"required,min=3,max=50"`
	Email     string          `gorm:"size:100;not null;unique" json:"email" validate:"required,email"`
	Password  string          `gorm:"size:100;not null;" json:"password" validate:"required,min=8"`
	CreatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deletedAt,omitempty"`
}

func (user *User) Beforesave(tx *gorm.DB) error {
	if len(user.Password) > 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return nil

}
