package repostory

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"money_share/pkg/model"
	"money_share/pkg/repository"
	"money_share/test_tool/database"
	testmodel "money_share/test_tool/model"
	"testing"
)

type GroupRepositoryTestSuite struct {
	suite.Suite
	DB               *gorm.DB
	GroupRepository  repository.GroupRepository
	UserRepository   repository.UserRepository
	MemberRepository repository.MemberRepository
	PrePopulateUser  model.User
	PrePopulateGroup model.Group
}

func (suite *GroupRepositoryTestSuite) SetupTest() {
	// Establish connection to database
	db, err := database.Connect()
	suite.NoError(err)
	suite.NotNil(db, "database connection should be available")
	suite.DB = db
	suite.UserRepository = repository.NewUserRepository(db)
	suite.GroupRepository = repository.NewGroupRepository(db)
	suite.MemberRepository = repository.NewMemberRepository(db)

	// Pre-populate database with a user
	suite.PrePopulateUser = testmodel.GenerateRandomUser()
	err = suite.UserRepository.Create(&suite.PrePopulateUser)
	suite.NoError(err)
	suite.Greater(suite.PrePopulateUser.ID, uint(0), "created user must have id greater than 0")

	// Pre-populate database with a group and test group creation here as well
	suite.PrePopulateGroup = testmodel.GenerateRandomGroup()
	err = suite.GroupRepository.Create(&suite.PrePopulateGroup, suite.PrePopulateUser.ID)
	suite.NoError(err)
	suite.Greater(suite.PrePopulateGroup.ID, uint(0), "created group must have id greater than 0")

	// Validate manager of created group
	member, err := suite.MemberRepository.GetByID(suite.PrePopulateUser.ID, suite.PrePopulateGroup.ID)
	suite.NoError(err)
	suite.Equal(member.UserID, suite.PrePopulateUser.ID)
	suite.Equal(member.GroupID, suite.PrePopulateGroup.ID)
	suite.Equal(member.Role, "manager", "role should be manager")
}

func (suite *GroupRepositoryTestSuite) TestGetByID(){
	group, err := suite.GroupRepository.GetById(suite.PrePopulateGroup.ID)
	suite.NoError(err)
	// Compare fields
	suite.Equal(suite.PrePopulateGroup.ID, group.ID, "should have the same ID")
	suite.Equal(suite.PrePopulateGroup.Name, group.Name, "should have the same name")
	suite.Equal(suite.PrePopulateGroup.GroupImageUrl, group.GroupImageUrl, "should have the same image url")
	suite.Equal(suite.PrePopulateGroup.TotalExpense, group.TotalExpense, "should have the same total expense")
	suite.Equal(suite.PrePopulateGroup.AverageExpense, group.AverageExpense, "should have the same average expense")
	// Check member
	suite.Equal(len(group.Members), 1, "should have one member")
	member := group.Members[0]
	// Compare member fields
	suite.Equal(suite.PrePopulateUser.ID, member.User.ID, "should have the same ID")
	suite.Equal(suite.PrePopulateUser.Username, member.User.Username, "should have the same username")
	suite.Equal("manager", member.Role, "should have role as manager")

	// Error case: group not existed
	group, err = suite.GroupRepository.GetById(suite.PrePopulateGroup.ID + 100)
	suite.Error(err)
	suite.ErrorIs(err, gorm.ErrRecordNotFound)
}

func TestGroupRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GroupRepositoryTestSuite))
}
