package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"money_share/pkg/model"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepositoryImpl{db}
}

func (repository UserRepositoryImpl) GetById(userId uint) (*model.User, error) {
	db := repository.DB
	var user = &model.User{}
	err := db.First(user, userId).Error
	return user, err
}

func (repository UserRepositoryImpl) GetByUsername(username string) (*model.User, error) {
	db := repository.DB
	var user = &model.User{}
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

func (repository UserRepositoryImpl) CheckUsernameAvailability(username string) (bool, error) {
	db := repository.DB
	var recordFound int64
	err := db.Model(&model.User{}).Where("username = ?", username).Count(&recordFound).Error
	if err != nil {
		return false, err
	}
	return recordFound == 0, nil
}

func (repository UserRepositoryImpl) ValidateUsernameAndUserID(username string, userID uint) (bool, error) {
	db := repository.DB
	var recordFound int64
	err := db.Model(&model.User{}).Where("id = ? AND username = ?", userID, username).Count(&recordFound).Error
	if err != nil {
		return false, err
	}
	return recordFound != 0, nil
}

func (repository UserRepositoryImpl) Create(user *model.User) (*model.User, error) {
	// Skip all associations before inserting
	result := repository.DB.Omit(clause.Associations).Create(user)
	return user, result.Error
}

func (repository UserRepositoryImpl) Update(userID uint, updateMap map[string]interface{}) (*model.User, error) {
	db := repository.DB
	if userID == 0 {
		return &model.User{}, errors.New("user ID not provided")
	}

	user := &model.User{}
	user.ID = userID
	updateRs := db.Model(user).Clauses(clause.Returning{}).Updates(updateMap)
	return user, updateRs.Error
}

func (repository UserRepositoryImpl) Delete(userId uint) error {
	if userId == 0 {
		return errors.New("user ID not provided")
	}

	return repository.DB.Delete(&model.User{}, userId).Error
}
