package entity

import (
	"go-admin/global"
	"time"

	"gorm.io/gorm"
)

// BaseModel 通用基础模型，所有实体嵌入此结构
type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedBy int64          `json:"createdBy" gorm:"column:created_by;default:0"` // 创建者ID（用于数据权限-仅本人）
	UpdatedBy int64          `json:"updatedBy" gorm:"column:updated_by;default:0"` // 更新者ID
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 返回带前缀的表名，子结构体重写此方法时调用
func TableName(name string) string {
	prefix := global.Config.TablePrefix
	return prefix + name
}
