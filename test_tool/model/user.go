package model

import (
	"money_share/pkg/model"
	util2 "money_share/test_tool/util"
)

func GenerateRandomUser() model.User {
	return model.User{
		Username:        util2.RandomStringRange(8, 20),
		Password:        util2.RandomStringRange(8, 20),
		DisplayName:     util2.RandomStringRange(4, 32),
		ProfileImageUrl: util2.RandomStringRange(4, 32),
		PhoneNumber:     util2.RandomNumericString(9),
		EmailAddress:    util2.RandomStringRange(8, 32),
		DateOfBirth:     util2.RandomDateTillNow(1900).UTC(),
	}
}
