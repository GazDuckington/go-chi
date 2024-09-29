package model

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
}

// BeforeSave is a GORM hook to generate UUID before saving
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return nil
}

func (u *User) ComparePassword(plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
}
