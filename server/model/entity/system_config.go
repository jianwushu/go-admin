package entity

import "time"

// SystemConfig 系统配置实体
type SystemConfig struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	ConfigKey   string    `json:"configKey" gorm:"column:config_key;size:128;uniqueIndex;not null"`
	ConfigValue string    `json:"configValue" gorm:"column:config_value;type:text"`
	ConfigType  string    `json:"configType" gorm:"column:config_type;size:32;default:text"` // text/image/json
	Remark      string    `json:"remark" gorm:"size:512"`
	CreatedBy   int64     `json:"createdBy" gorm:"column:created_by;default:0"`
	UpdatedBy   int64     `json:"updatedBy" gorm:"column:updated_by;default:0"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

// TableName 返回带前缀的表名
func (SystemConfig) TableName() string {
	return TableName("system_config")
}
