package repository

import (
	"go-admin/utils"

	"gorm.io/gorm"
)

// BaseRepository 基础仓储层，提供数据权限过滤等通用方法
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository 创建基础仓储实例
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

// ApplyDataScope 对已有查询应用数据权限过滤
// query: 已有的 GORM 查询
// scopeInfo: 数据权限信息（从 Context 获取）
// tableName: 表别名（如 "u"，空字符串表示默认表）
// deptColumn: 部门ID字段名（默认 "dept_id"）
func (r *BaseRepository) ApplyDataScope(query *gorm.DB, scopeInfo *utils.DataScopeInfo, tableName string, deptColumn string) *gorm.DB {
	return utils.ApplyDataScope(query, scopeInfo, tableName, deptColumn)
}
