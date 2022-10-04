package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	GroupIdentifier string    `gorm:"unique;not null"`
	Name            string    `gorm:"not null"`
	TotalExpense    float32   `gorm:"default:0"`
	AverageExpense  float32   `gorm:"default:0"`
	Members         []Member  `gorm:"constraint:OnDelete:CASCADE"` // One-to-many relationship with Member entity
	Expenses        []Expense `gorm:"constraint:OnDelete:CASCADE"` // One-to-many relationship with Expense entity
}
