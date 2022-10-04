package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"money_share/pkg/model"
)

type ExpenseRepositoryImpl struct {
	DB *gorm.DB
}

func NewExpenseRepository(db *gorm.DB) ExpenseRepository {
	return ExpenseRepositoryImpl{db}
}

func (repository ExpenseRepositoryImpl) GetById(expenseId uint) (*model.Expense, error) {
	db := repository.DB
	if expenseId <= 0 {
		return nil, errors.New("expenseId must be greater than 0")
	}
	var expense *model.Expense
	err := db.First(expense, expenseId).Error
	return expense, err
}

func (repository ExpenseRepositoryImpl) GetByGroup(groupId uint) ([]*model.Expense, error) {
	db := repository.DB
	if groupId <= 0 {
		return nil, errors.New("groupId must be greater than 0")
	}
	var expenses []*model.Expense
	err := db.Where("group_id = ?", groupId).Find(&expenses).Error
	return expenses, err
}

func (repository ExpenseRepositoryImpl) GetByMember(memberId uint, groupId uint) ([]*model.Expense, error) {
	db := repository.DB
	if memberId <= 0 || groupId <= 0 {
		return nil, errors.New("memberId and groupId must be greater than 0")
	}
	var expenses []*model.Expense
	err := db.Where("group_id = ? AND member_id = ?", groupId, memberId).Find(&expenses).Error
	return expenses, err
}

func (repository ExpenseRepositoryImpl) Create(expense *model.Expense) error {
	db := repository.DB
	// Validate fields
	if len(expense.Title) == 0 {
		return errors.New("title cannot be empty")
	}
	if expense.Amount < 0 {
		return errors.New("amount must be equal or greater than 0")
	}
	if expense.PurchaseTime.IsZero() {
		return errors.New("purchase time is not set")
	}
	if expense.MemberID <= 0 {
		return errors.New("memberId must be greater than 0")
	}

	err := db.Create(expense).Error
	return err
}

func (repository ExpenseRepositoryImpl) Update(expense *model.Expense) error {
	db := repository.DB
	// Validate fields
	if expense.ID <= 0 {
		return errors.New("expenseId must be greater than 0")
	}
	if expense.Amount < 0 {
		return errors.New("amount must be equal or greater than 0")
	}
	if len(expense.Status) > 0 && (expense.Status != "pending" && expense.Status != "accepted" && expense.Status != "denied") {
		return errors.New("invalid status, must be 'pending', 'accepted' or 'denied'")
	}

	updateExpense := &model.Expense{}
	updateExpense.ID = expense.ID

	err := db.Transaction(func(tx *gorm.DB) error {
		// Make sure record exists
		queryRs := db.First(updateExpense)
		if err := queryRs.Error; err != nil {
			return err
		}

		// Update fields
		// Omit forbidden fields
		err := queryRs.Omit("ID", "MemberID", "GroupID", clause.Associations).Updates(expense).Error

		return err
	})

	return err
}

func (repository ExpenseRepositoryImpl) Delete(expenseId uint) error {
	db := repository.DB
	return db.Delete(&model.Expense{}, expenseId).Error
}
