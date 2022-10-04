package repository

import "money_share/pkg/model"

type ExpenseRepository interface {
	GetById(expenseId uint) (*model.Expense, error)
	GetByGroup(groupId uint) ([]*model.Expense, error)
	GetByMember(memberId uint, groupId uint) ([]*model.Expense, error)
	Create(expense *model.Expense) error
	Update(expense *model.Expense) error
	Delete(expenseId uint) error
}
