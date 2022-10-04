package dto

import (
	"gorm.io/gorm"
	"log"
	"money_share/pkg/model"
	"money_share/pkg/util"
	"time"
)

type ExpenseDTO struct {
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Amount       float32 `json:"amount"`
	PurchaseTime string  `json:"purchaseTime"`
	Status       string  `json:"status"` // pending, approved, denied
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
	}
}
