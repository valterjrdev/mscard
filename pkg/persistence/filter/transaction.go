package filter

import "gorm.io/gorm"

type (
	TransactionCollection struct {
		Page           int
		Size           int
		Account        string
		EventDateStart string
		EventDateEnd   string
	}
)

func (t *TransactionCollection) Filter() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if t.Account != "" {
			db.Where("account_id = ?", t.Account)
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
