package filter

import (
	"gorm.io/gorm"
	"strconv"
)

type (
	TransactionCollection struct {
		Page            int
		Size            int
		Account         string
		Operation       string
		CreateDateStart string
		CreateDateEnd   string
	}
)

func (t *TransactionCollection) Filter() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Account != "" {
			account, _ := strconv.Atoi(t.Account)
			db.Where("account_id = ?", account)
		}

		if t.Operation != "" {
			operation, _ := strconv.Atoi(t.Operation)
			db.Where("operation_id = ?", operation)
		}

		if t.CreateDateStart != "" && t.CreateDateEnd != "" {
			db.Where(
				"TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') BETWEEN ? AND ?",
				t.CreateDateStart,
				t.CreateDateEnd,
			)
		}

		return db
	}
}
