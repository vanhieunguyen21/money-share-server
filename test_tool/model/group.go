package model

import (
	"money_share/pkg/model"
	"money_share/test_tool/util"
)

func GenerateRandomGroup() model.Group {
	return model.Group{
		Name:           util.RandomStringRange(4, 32),
		GroupImageUrl:  util.RandomStringRange(4, 32),
	}
}
