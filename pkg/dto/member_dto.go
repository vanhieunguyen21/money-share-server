package dto

import "money_share/pkg/model"

type MemberDTO struct {
	User         UserDTO `json:"user"`
	TotalExpense float32 `json:"totalExpense"`
	Role         string  `json:"role"`
}

func (dto MemberDTO) MapToDomain() (model.Member, error) {
	user, err := dto.User.MapToDomain()
	if err != nil {
		return model.Member{}, err
	}
	return model.Member{
		User:         user,
		TotalExpense: dto.TotalExpense,
		Role:         dto.Role,
	}, nil
}
