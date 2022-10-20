package repository

import "money_share/pkg/model"

type GroupRepository interface {
	GetById(groupId uint) (*model.Group, error)
	GetByUser(memberId uint) ([]*model.Group, error)
	Create(group *model.Group, creatorID uint) error
	Update(group *model.Group) error
	Delete(groupId uint) error
	GetMemberRole(memberID uint, groupID uint) (string, error)
}
