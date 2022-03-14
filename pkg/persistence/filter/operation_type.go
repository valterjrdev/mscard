package filter

import (
	"gorm.io/gorm"
	"strconv"
)

type (
	OperationTypeCollection struct {
		Page        int
		Size        int
		Description string
		Negative    string
	}
)

func (t *OperationTypeCollection) Filter() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Negative != "" {
			negative, _ := strconv.ParseBool(t.Negative)
			db.Where("negative = ?", negative)
		}

		if t.Description != "" {
			db.Where("description ILIKE ?", "%"+t.Description+"%")
		}

		return db
	}
}
