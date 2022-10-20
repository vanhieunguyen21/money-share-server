package model

type Member struct {
	User         User    ``                      // Owned relationship with User entity
	Role         string  `gorm:"default:member"` // member/manager
	TotalExpense float32 `gorm:"default:0"`      //
	UserID       uint    `gorm:"primaryKey"`     // Many-to-one relationship with User entity
	GroupID      uint    `gorm:"primaryKey"`     // Many-to-one relationship with Group entity
	// One-to-many relationship with Expense entity
	Expenses []Expense `gorm:"foreignKey:MemberID,GroupID;references:UserID,GroupID;constraint:OnDelete:SET NULL;"`
}

// TODO: update average expense of group when new member is created