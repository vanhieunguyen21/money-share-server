package dto

import (
	"gorm.io/gorm"
	"money_share/pkg/model"
)

type GroupDTO struct {
	ID              uint         `json:"id"`
	Name            string       `json:"name"`
	GroupImageUrl   string       `json:"groupImageUrl"`
	TotalExpense    float32      `json:"totalExpense"`
	AverageExpense  float32      `json:"averageExpense"`
	Members         []MemberDTO  `json:"members,omitempty"`
	Expenses        []ExpenseDTO `json:"expenses,omitempty"`
}

func (dto GroupDTO) MapToDomain() (model.Group, error) {
	// Parse members
	var members []model.Member
	for _, memberDTO := range dto.Members {
		member, err := memberDTO.MapToDomain()
		if err != nil {
			return model.Group{}, err
		}
		members = append(members, member)
	}
	// Parse expenses
	var expenses []model.Expense
	for _, expenseDTO := range dto.Expenses {
		expenses = append(expenses, expenseDTO.MapToDomain())
	}

	return model.Group{
		Model:           gorm.Model{ID: dto.ID},
		Name:            dto.Name,
		GroupImageUrl:   dto.GroupImageUrl,
		TotalExpense:    dto.TotalExpense,
		AverageExpense:  dto.AverageExpense,
		Members:         members,
		Expenses:        expenses,
	}, nil
}
