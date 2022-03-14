package filter

import (
	"gorm.io/gorm"
	"strconv"
)

type (
	TransactionCollection struct {
		Page           int
		Size           int
		Account        string
		Type           string
		EventDateStart string
		EventDateEnd   string
	}
)

func (t *TransactionCollection) Filter() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Account != "" {
			account, _ := strconv.Atoi(t.Account)
			db.Where("account_id = ?", account)
		}

		if t.Type != "" {
			operationType, _ := strconv.Atoi(t.Type)
			db.Where("operation_type_id = ?", operationType)
		}

		if t.EventDateStart != "" && t.EventDateEnd != "" {
			db.Where(
				"TO_CHAR(event_date, 'YYYY-MM-DD HH24:MI:SS') BETWEEN ? AND ?",
				t.EventDateStart,
				t.EventDateEnd,
			)
		}

		return db
	}
}
