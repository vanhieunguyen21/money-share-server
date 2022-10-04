package repository

import (
	"gorm.io/gorm"
	"money_share/pkg/model"
)

type MemberRepositoryImpl struct {
	DB *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return MemberRepositoryImpl{db}
}

func (repository MemberRepositoryImpl) GetByID(userID uint, groupID uint) (*model.Member, error) {
	db := repository.DB
	member := &model.Member{}
	err := db.Where("user_id = ? AND group_id = ?", userID, groupID).Preload("User").First(member).Error
	return member, err
}

func (repository MemberRepositoryImpl) GetByGroup(groupID uint) ([]*model.Member, error) {
	db := repository.DB
	var members []*model.Member
	err := db.Where("group_id = ?", groupID).Preload("User").Preload("User").Find(&members).Error
	return members, err
}

func (repository MemberRepositoryImpl) AddMemberToGroup(userID uint, groupID uint) error {
	db := repository.DB
	member := &model.Member{
		UserID:  userID,
		GroupID: groupID,
	}
	err := db.Create(member).Error
	return err
}

func (repository MemberRepositoryImpl) RemoveMemberFromGroup(userID uint, groupID uint) error {
	db := repository.DB
	err := db.Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&model.Member{}).Error
	return err
}

func (repository MemberRepositoryImpl) IncreaseTotalExpense(userID uint, groupID uint, updateValue float32) error {
	db := repository.DB
	err := db.Transaction(func(tx *gorm.DB) error {
		member := &model.Member{}
		err := tx.Where("user_id = ? AND group_id = ?", userID, groupID).First(member).Error
		if err != nil {
			return err
		}
		newExpense := member.TotalExpense + updateValue
		err = tx.Model(member).Update("total_expense", newExpense).Error
		return err
	})
	return err
}
