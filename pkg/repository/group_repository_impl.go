package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"money_share/pkg/model"
)

type GroupRepositoryImpl struct {
	DB *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return GroupRepositoryImpl{db}
}

func (repository GroupRepositoryImpl) GetById(groupId uint) (*model.Group, error) {
	db := repository.DB
	if groupId <= 0 {
		return &model.Group{}, errors.New("groupId must be greater than 0")
	}

	group := &model.Group{}
	group.ID = groupId
	err := db.Preload("Members").Preload("Members.User").Preload("Expenses").Find(group).Error
	return group, err
}

func (repository GroupRepositoryImpl) GetByUser(userId uint) ([]*model.Group, error) {
	db := repository.DB
	if userId <= 0 {
		return nil, errors.New("userId must be greater than 0")
	}
	var groups []*model.Group
	err := db.Where("id in (?)",
		db.Table("members").Select("group_id").Where("user_id = ?", userId)).Find(&groups).Error
	return groups, err
}

func (repository GroupRepositoryImpl) Create(group *model.Group, username string) error {
	db := repository.DB
	if len(username) == 0 {
		return errors.New("username cannot be empty")
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		// Validate creator existence
		creator := &model.User{}
		if err := tx.Where("username = ?", username).First(creator).Error; err != nil {
			return err
		}
		// Create group
		if err := tx.Omit(clause.Associations).Create(group).Error; err != nil {
			return err
		}
		// Create member object from creator
		member := &model.Member{
			TotalExpense: 0,
			UserID:       creator.ID,
			GroupID:      group.ID,
			Role:         "manager",
		}
		if err := tx.Model(group).Association("Members").Append(member); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (repository GroupRepositoryImpl) Update(group *model.Group) error {
	db := repository.DB
	if group.ID == 0 {
		return errors.New("group ID not provided")
	}
	updateGroup := &model.Group{}
	updateGroup.ID = group.ID

	err := db.Transaction(func(tx *gorm.DB) error {
		// Make sure record exists
		queryRs := db.First(&updateGroup)
		if err := queryRs.Error; err != nil {
			return err
		}

		// Update fields
		err := queryRs.Omit(clause.Associations).Updates(group).Error
		return err
	})

	return err
}

func (repository GroupRepositoryImpl) Delete(groupId uint) error {
	db := repository.DB
	if groupId == 0 {
		return errors.New("group ID not provided")
	}

	return db.Delete(&model.Group{}, groupId).Error
}

func (repository GroupRepositoryImpl) GetMemberRole(memberID uint, groupID uint) (role string, err error) {
	db := repository.DB
	if memberID <= 0 || groupID <= 0 {
		err = errors.New("member ID or group ID must be greater than 0")
		return
	}

	result := struct {
		Role string
	}{}
	err = db.Model(&model.Member{}).Where("user_id = ? AND group_id = ?", memberID, groupID).Find(&result).Error
	if err != nil {
		return
	}
	role = result.Role
	return
}
