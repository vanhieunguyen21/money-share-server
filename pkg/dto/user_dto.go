package dto

import (
	"gorm.io/gorm"
	"log"
	"money_share/pkg/model"
	"money_share/pkg/util"
	"time"
)

type UserDTO struct {
	ID              uint   `json:"id,omitempty"`
	Username        string `json:"username,omitempty"`
	DisplayName     string `json:"displayName,omitempty"`
	ProfileImageUrl string `json:"profileImageUrl,omitempty"`
	PhoneNumber     string `json:"phoneNumber,omitempty"`
	EmailAddress    string `json:"emailAddress,omitempty"`
	DateOfBirth     string `json:"dateOfBirth,omitempty"`
}

func (dto UserDTO) MapToDomain() (model.User, error) {
	// Parse date of birth
	dob := time.Time{}
	var err error
	if dto.DateOfBirth != "" {
		dob, err = time.Parse(util.ShortDateLayout, dto.DateOfBirth)
		if err != nil {
			log.Println(err)
			return model.User{}, err
		}
	}

	return model.User{
		Model:           gorm.Model{ID: dto.ID},
		Username:        dto.Username,
		DisplayName:     dto.DisplayName,
		ProfileImageUrl: dto.ProfileImageUrl,
		PhoneNumber:     dto.PhoneNumber,
		EmailAddress:    dto.EmailAddress,
		DateOfBirth:     dob,
	}, nil
}
