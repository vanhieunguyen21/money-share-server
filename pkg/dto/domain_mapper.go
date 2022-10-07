package dto

import (
	"money_share/pkg/model"
	"money_share/pkg/util"
)

func ExpenseToExpenseDTO(domain model.Expense) ExpenseDTO {
	// Convert purchase time to string
	purchaseTime := domain.PurchaseTime.Format(util.DateTimeLayout)

	return ExpenseDTO{
		ID:           domain.Model.ID,
		Title:        domain.Title,
		Description:  domain.Description,
		Amount:       domain.Amount,
		PurchaseTime: purchaseTime,
		Status:       domain.Status,
	}
}

func UserToUserDTO(domain model.User) UserDTO {
	// Convert dob to string
	dob := domain.DateOfBirth.Format(util.ShortDateLayout)

	return UserDTO{
		ID:              domain.Model.ID,
		Username:        domain.Username,
		DisplayName:     domain.DisplayName,
		ProfileImageUrl: domain.ProfileImageUrl,
		PhoneNumber:     domain.PhoneNumber,
		EmailAddress:    domain.EmailAddress,
		DateOfBirth:     dob,
	}
}

func GroupToGroupDTO(domain model.Group) GroupDTO {
	// Map members
	var members []MemberDTO
	for _, member := range domain.Members {
		members = append(members, MemberToMemberDTO(member))
	}

	// Map expenses
	var expenses []ExpenseDTO
	for _, expense := range domain.Expenses {
		expenses = append(expenses, ExpenseToExpenseDTO(expense))
	}

	return GroupDTO{
		ID:              domain.ID,
		GroupIdentifier: domain.GroupIdentifier,
		Name:            domain.Name,
		GroupImageUrl:   domain.GroupImageUrl,
		TotalExpense:    domain.TotalExpense,
		AverageExpense:  domain.AverageExpense,
		Members:         members,
		Expenses:        expenses,
	}
}

func MemberToMemberDTO(domain model.Member) MemberDTO {
	return MemberDTO{
		User:         UserToUserDTO(domain.User),
		Role:         domain.Role,
		TotalExpense: domain.TotalExpense,
	}
}
