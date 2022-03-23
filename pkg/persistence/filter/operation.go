package filter

import (
	"gorm.io/gorm"
	"strconv"
)

type (
	OperationCollection struct {
		Page        int
		Size        int
		Description string
		Debit       string
	}
)

func (t *OperationCollection) Filter() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Debit != "" {
			negative, _ := strconv.ParseBool(t.Debit)
			db.Where("debit = ?", negative)
		}

		if t.Description != "" {
			db.Where("description ILIKE ?", "%"+t.Description+"%")
		}

		return db
	}
}
