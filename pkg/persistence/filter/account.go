package filter

import (
	"gorm.io/gorm"
)

type (
	AccountCollection struct {
		Page     int
		Size     int
		Document string
	}
)

func (t *AccountCollection) Filter() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Document != "" {
			db.Where("document_number ILIKE ?", "%"+t.Document+"%")
		}

		return db
	}
}
