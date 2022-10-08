package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Username        string    `gorm:"unique;not null"`
	Password        string    `gorm:"not null"`
	DisplayName     string    `gorm:"not null"`
	ProfileImageUrl string    ``
	PhoneNumber     string    ``
	EmailAddress    string    ``
	DateOfBirth     time.Time ``
	Members         []Member  `gorm:"constraint:OnDelete:SET NULL;"` // One to many with Member entity
}

func (u *User) ValidateUsername() (err error) {
	return ValidateUsername(u.Username)
}

func ValidateUsername(username string) (err error) {
	// Validate length
	if len(username) < 8 {
		err = errors.New("username must be at least 8 characters")
		return
	}
	if len(username) > 20 {
		err = errors.New("username must be at most 20 characters")
		return
	}
	// Validate characters used
	match, _ := regexp.MatchString("^[a-zA-Z0-9._]+$", username)
	if !match {
		err = errors.New("username must only contains alphabet characters, number and/or dot(.) and/or underscore(_)")
		return
	}

	// [^_.].*[^_.]: no _ or . at the beginning or at the end
	// (?!.*[_.]{2}): no __ or _. or ._ or .. inside
	match, _ = regexp.MatchString("^(?!.*[_.]{2})[^_.].*[^_.]$", username)
	if !match {
		err = errors.New("username must not contains no __ or _. or ._ or .. inside or _ or . at the beginning or at the end")
		return
	}

	return
}

func (u *User) ValidatePassword() (err error) {
	return ValidatePassword(u.Password)
}

func ValidatePassword(password string) (err error) {
	// Validate length
	if len(password) < 8 {
		err = errors.New("password must be at least 8 characters")
		return
	}
	if len(password) > 20 {
		err = errors.New("password must be at most 20 characters")
		return
	}

	return
}

func (u *User) TrimDisplayName() {
	u.DisplayName = strings.TrimSpace(u.DisplayName)
}

func (u *User) ValidateDisplayName()(err error) {
	return ValidateDisplayName(u.DisplayName)
}

func ValidateDisplayName(displayName string) (err error) {
	// Validate length
	if len(displayName) < 4 {
		err = errors.New("display name must be at least 4 characters")
		return
	}
	if len(displayName) > 32 {
		err = errors.New("password must be at most 32 characters")
		return
	}

	return
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

func (u *User) ValidateFields() (err error) {
	if err = u.ValidateUsername(); err != nil {
		return
	}
	if err = u.ValidatePassword(); err != nil {
		return
	}
	if err = u.ValidateDisplayName(); err != nil {
		return
	}
	return
}

func (u *User) ValidateNonNullFields() (err error) {
	// Validate username
	if u.Username != "" {
		err = u.ValidateUsername()
		if err != nil {
			return
		}
	}
	// Validate password
	if u.Password != "" {
		err = u.ValidatePassword()
		if err != nil {
			return
		}
	}
	// Validate display name
	if u.DisplayName != "" {
		err = u.ValidateDisplayName()
		if err != nil {
			return
		}
	}

	return
}
