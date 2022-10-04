package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username     string    `gorm:"unique;not null"`
	Password     string    `gorm:"not null"`
	DisplayName  string    ``
	PhoneNumber  string    ``
	EmailAddress string    ``
	DateOfBirth  time.Time ``
	Members      []Member  `gorm:"constraint:OnDelete:SET NULL;"` // One to many with Member entity
}

func (u *User) HashPassword() (err error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err == nil {
		u.Password = string(hashedPwd)
	}
	return
}

func (u *User) ComparePassword(providedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPwd))
	return err == nil
}
