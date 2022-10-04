package repository

import "money_share/pkg/model"

type MemberRepository interface {
	GetByID(userID uint, groupID uint) (*model.Member, error)
	GetByGroup(groupID uint) ([]*model.Member, error)
	AddMemberToGroup(userID uint, groupID uint) error
	RemoveMemberFromGroup(userID uint, groupID uint) error
	IncreaseTotalExpense(userID uint, groupID uint, updateValue float32) error
}
