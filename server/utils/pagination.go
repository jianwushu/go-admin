package utils

import (
	"go-admin/model/request"

	"gorm.io/gorm"
)

// Paginate GORM 分页作用域
// 使用方式: db.Scopes(Paginate(pageReq)).Find(&list)
func Paginate(page request.PageRequest) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageNum := page.GetPage()
		pageSize := page.GetPageSize()
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// PaginateWithParams 使用原始参数的分页作用域
func PaginateWithParams(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
