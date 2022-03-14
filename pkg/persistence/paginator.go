package persistence

import "gorm.io/gorm"

const (
	PaginatorPageDefault = 1
	PaginatorSizeDefault = 10
)

func Paginator(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = PaginatorPageDefault
		}

		if size == 0 || size > 100 {
			size = PaginatorSizeDefault
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
