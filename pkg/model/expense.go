package model

import (
	"gorm.io/gorm"
	"time"
)

type Expense struct {
	gorm.Model
	Title        string    `gorm:"not null"`
	Description  string    ``
	Amount       float32   `gorm:"not null"`
	PurchaseTime time.Time `gorm:"not null"`
	Status       string    `gorm:"not null;default:pending"` // pending, approved, denied
	GroupID      uint      ``                                // Many-to-one relationship with Group entity
	MemberID     uint      ``                                // Many-to-one relationship with Member entity
}

// Hooks to update member and group expenses

func (e *Expense) AfterCreate(tx *gorm.DB) (err error) {
	return e.updateMemberAndGroupExpenses(tx)
}

func (e *Expense) AfterUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Status") && e.Status == "approved" {
		return e.updateMemberAndGroupExpenses(tx)
	}
	if tx.Statement.Changed("Amount") {
		return e.updateMemberAndGroupExpenses(tx)
	}
	return
}

func (e *Expense) AfterDelete(tx *gorm.DB) (err error) {
	if e.Status != "approved" {
		return e.updateMemberAndGroupExpenses(tx)
	}
	return
}

func (e *Expense) updateMemberAndGroupExpenses(tx *gorm.DB) (err error) {
	// Update member total expense
	err = tx.Model(&Member{}).Where("user_id = ? AND group_id = ?",
		e.MemberID, e.GroupID).Update("total_expense",
		tx.Model(&Expense{}).Select("SUM(amount)").Where("member_id = ? AND group_id = ? AND status='approved'",
			e.MemberID, e.GroupID)).Error
	if err != nil {
		return
	}

	// Update group total expense and average expense
	groupSpecificClause := "group_id = ? AND status='approved'"
	err = tx.Model(&Group{}).Where("id = ?", e.GroupID).Updates(map[string]interface{}{
		"total_expense": tx.Model(&Expense{}).Select("SUM(amount)").Where(groupSpecificClause, e.GroupID),
		"average_expense": tx.Model(&Expense{}).Select("SUM(amount)/COUNT(*)").Where(groupSpecificClause, e.GroupID),
	}).Error

	return
}
