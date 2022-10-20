package repostory

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"money_share/pkg/model"
	"money_share/pkg/repository"
	"money_share/test_tool/database"
	test_model "money_share/test_tool/model"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	DB              *gorm.DB
	UserRepository  repository.UserRepository
	PrePopulateUser model.User
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// Establish connection to DB
	db, err := database.Connect()
	suite.NoError(err)
	suite.NotNil(db)
	suite.DB = db
	suite.UserRepository = repository.NewUserRepository(db)

	// Pre-populate with a dummy user and test user creation here as well
	suite.PrePopulateUser = test_model.GenerateRandomUser()

	err = suite.UserRepository.Create(&suite.PrePopulateUser)
	suite.NoError(err)
	suite.Greater(suite.PrePopulateUser.ID, uint(0), "created user must have id greater than 0")
}

func (suite *UserRepositoryTestSuite) TestGetByID() {
	// Get user from database
	actualUser, err := suite.UserRepository.GetById(suite.PrePopulateUser.ID)
	suite.NoError(err)
	// Compare fields
	suite.Equal(suite.PrePopulateUser.ID, actualUser.ID, "should have the same ID")
	suite.Equal(suite.PrePopulateUser.Username, actualUser.Username, "should have the same username")
	suite.Equal(suite.PrePopulateUser.DisplayName, actualUser.DisplayName, "should have the same display name")
	suite.Equal(suite.PrePopulateUser.ProfileImageUrl, actualUser.ProfileImageUrl, "should have the same profile image URL")
	suite.Equal(suite.PrePopulateUser.PhoneNumber, actualUser.PhoneNumber, "should have the same phone number")
	suite.Equal(suite.PrePopulateUser.EmailAddress, actualUser.EmailAddress, "should have the same email address")
	suite.Equal(suite.PrePopulateUser.DateOfBirth.UTC(), actualUser.DateOfBirth.UTC(), "should have the same date of birth")

	// Error case
	_, err = suite.UserRepository.GetById(suite.PrePopulateUser.ID + 100)
	suite.Error(err, "should return an error")
	suite.ErrorIs(err, gorm.ErrRecordNotFound, "error should be `record not found`")
}

func (suite *UserRepositoryTestSuite) TestGetByUsername() {
	// Get actual user from database
	actualUser, err := suite.UserRepository.GetByUsername(suite.PrePopulateUser.Username)
	suite.NoError(err)
	// Compare fields
	suite.Equal(suite.PrePopulateUser.ID, actualUser.ID, "should have the same ID")
	suite.Equal(suite.PrePopulateUser.Username, actualUser.Username, "should have the same username")
	suite.Equal(suite.PrePopulateUser.DisplayName, actualUser.DisplayName, "should have the same display name")
	suite.Equal(suite.PrePopulateUser.ProfileImageUrl, actualUser.ProfileImageUrl, "should have the same profile image URL")
	suite.Equal(suite.PrePopulateUser.PhoneNumber, actualUser.PhoneNumber, "should have the same phone number")
	suite.Equal(suite.PrePopulateUser.EmailAddress, actualUser.EmailAddress, "should have the same email address")
	suite.Equal(suite.PrePopulateUser.DateOfBirth.UTC(), actualUser.DateOfBirth.UTC(), "should have the same date of birth")

	// Error case
	_, err = suite.UserRepository.GetByUsername("no_such_username")
	suite.Error(err, "should return an error")
	suite.ErrorIs(err, gorm.ErrRecordNotFound, "error should be `record not found`")
}

func (suite *UserRepositoryTestSuite) TestCheckUsernameAvailability() {
	// Normal case
	available, err := suite.UserRepository.CheckUsernameAvailability("no_such_username")
	suite.NoError(err, "should not return error")
	suite.True(available, "should be available")

	// Error case
	available, err = suite.UserRepository.CheckUsernameAvailability(suite.PrePopulateUser.Username)
	suite.NoError(err, "should not return error")
	suite.False(available, "should not be available")
}

func (suite *UserRepositoryTestSuite) TestValidateUsernameAndUserID() {
	// Normal case
	validated, err := suite.UserRepository.ValidateUsernameAndUserID(suite.PrePopulateUser.Username, suite.PrePopulateUser.ID)
	suite.NoError(err, "should not return error")
	suite.True(validated, "should be validated")

	// Error cases
	validated, err = suite.UserRepository.ValidateUsernameAndUserID("no_such_username", suite.PrePopulateUser.ID)
	suite.NoError(err, "should not return error")
	suite.False(validated, "should not be validated")

	validated, err = suite.UserRepository.ValidateUsernameAndUserID(suite.PrePopulateUser.Username, suite.PrePopulateUser.ID+100)
	suite.NoError(err, "should not return error")
	suite.False(validated, "should not be validated")

	validated, err = suite.UserRepository.ValidateUsernameAndUserID("no_such_username", suite.PrePopulateUser.ID+100)
	suite.NoError(err, "should not return error")
	suite.False(validated, "should not be validated")
}

func (suite *UserRepositoryTestSuite) TestUpdate() {
	// Create a new user
	newUser := test_model.GenerateRandomUser()
	_ = suite.UserRepository.Create(&newUser)

	// Create new user update object
	updateUser := test_model.GenerateRandomUser()
	// Create an update map
	updateMap := make(map[string]interface{})
	updateMap["Username"] = updateUser.Username
	updateMap["Password"] = updateUser.Password
	updateMap["DisplayName"] = updateUser.DisplayName
	updateMap["ProfileImageUrl"] = updateUser.ProfileImageUrl
	updateMap["PhoneNumber"] = updateUser.PhoneNumber
	updateMap["EmailAddress"] = updateUser.EmailAddress
	updateMap["DateOfBirth"] = updateUser.DateOfBirth
	// Update the user
	updatedUser, err := suite.UserRepository.Update(newUser.ID, updateMap)
	suite.NoError(err, "should not return error")
	// Compare returned user
	suite.Equal(newUser.ID, updatedUser.ID, "should have the same ID")
	suite.Equal(updateUser.Username, updatedUser.Username, "should have the same username")
	suite.Equal(updateUser.DisplayName, updatedUser.DisplayName, "should have the same display name")
	suite.Equal(updateUser.ProfileImageUrl, updatedUser.ProfileImageUrl, "should have the same profile image URL")
	suite.Equal(updateUser.PhoneNumber, updatedUser.PhoneNumber, "should have the same phone number")
	suite.Equal(updateUser.EmailAddress, updatedUser.EmailAddress, "should have the same email address")
	suite.Equal(updateUser.DateOfBirth.UTC(), updatedUser.DateOfBirth.UTC(), "should have the same date of birth")

	// Compare user in database
	actualUser, err := suite.UserRepository.GetById(newUser.ID)
	suite.NoError(err, "should not return error")
	// Compare fields
	suite.Equal(updatedUser.ID, actualUser.ID, "should have the same ID")
	suite.Equal(updatedUser.Username, actualUser.Username, "should have the same username")
	suite.Equal(updatedUser.DisplayName, actualUser.DisplayName, "should have the same display name")
	suite.Equal(updatedUser.ProfileImageUrl, actualUser.ProfileImageUrl, "should have the same profile image URL")
	suite.Equal(updatedUser.PhoneNumber, actualUser.PhoneNumber, "should have the same phone number")
	suite.Equal(updatedUser.EmailAddress, actualUser.EmailAddress, "should have the same email address")
	suite.Equal(updatedUser.DateOfBirth.UTC(), actualUser.DateOfBirth.UTC(), "should have the same date of birth")
}

func (suite *UserRepositoryTestSuite) TestDelete() {
	// Create a new user in database
	newUser := test_model.GenerateRandomUser()
	err := suite.UserRepository.Create(&newUser)
	suite.NoError(err, "should not return error")

	// Validate user exists in database
	savedUser, err := suite.UserRepository.GetById(newUser.ID)
	suite.NoError(err, "should not return error")
	suite.NotEmpty(savedUser, "should return an user")
	suite.Equal(newUser.ID, savedUser.ID, "should have the same ID")

	// Delete user from database
	err = suite.UserRepository.Delete(savedUser.ID)
	suite.NoError(err, "should not return error")

	// Try to retrieve deleted user
	_, err = suite.UserRepository.GetById(savedUser.ID)
	suite.Error(err, "should return an error")
	suite.ErrorIs(err, gorm.ErrRecordNotFound, "error should be record not found")

	// Error case: try to delete non-exist user
	err = suite.UserRepository.Delete(savedUser.ID + 100)
	suite.Error(err, "should return an error")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
