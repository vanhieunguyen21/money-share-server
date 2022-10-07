package dto

import (
	"gorm.io/gorm"
	"log"
	"money_share/pkg/model"
	"money_share/pkg/util"
	"time"
)

type UserDTO struct {
	ID              uint   `json:"id"`
	Username        string `json:"username"`
	DisplayName     string `json:"displayName"`
	ProfileImageUrl string `json:"profileImageUrl"`
	PhoneNumber     string `json:"phoneNumber"`
	EmailAddress    string `json:"emailAddress"`
	DateOfBirth     string `json:"dateOfBirth"`
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
