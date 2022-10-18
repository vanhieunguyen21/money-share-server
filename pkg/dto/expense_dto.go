package dto

import (
	"gorm.io/gorm"
	"log"
	"money_share/pkg/model"
	"money_share/pkg/util"
	"time"
)

type ExpenseDTO struct {
	ID           uint    `json:"id,omitempty"`
	Title        string  `json:"title,omitempty"`
	Description  string  `json:"description,omitempty"`
	Amount       float32 `json:"amount,omitempty"`
	PurchaseTime string  `json:"purchaseTime,omitempty"`
	Status       string  `json:"status,omitempty"` // pending, approved, denied
	MemberID     uint    `json:"memberID,omitempty"`
	GroupID      uint    `json:"groupID,omitempty"`
}

func (dto ExpenseDTO) MapToDomain() model.Expense {
	// Parse purchase time
	purchaseTime, err := time.Parse(util.DateTimeLayout, dto.PurchaseTime)
	if err != nil {
		log.Println(err)
	}

	return model.Expense{
		Model:        gorm.Model{ID: dto.ID},
		Title:        dto.Title,
		Description:  dto.Description,
		Amount:       dto.Amount,
		PurchaseTime: purchaseTime,
		Status:       dto.Status,
		MemberID:     dto.MemberID,
		GroupID:      dto.GroupID,
	}
}
