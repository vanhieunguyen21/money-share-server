package request

import "money_share/pkg/dto"

type ExpenseCreationRequest struct {
	Expense  dto.ExpenseDTO `json:"expense"`
	GroupID  uint           `json:"groupID"`
	MemberID uint           `json:"memberID"`
}
