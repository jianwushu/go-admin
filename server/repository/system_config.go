package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/request"

	"gorm.io/gorm"
)

// SystemConfigRepository 系统配置仓储层
type SystemConfigRepository struct {
	BaseRepository
}

// NewSystemConfigRepository 创建系统配置仓储实例
func NewSystemConfigRepository() *SystemConfigRepository {
	return &SystemConfigRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// FindAll 查询所有配置
func (r *SystemConfigRepository) FindAll() ([]entity.SystemConfig, error) {
	var configs []entity.SystemConfig
	if err := r.DB.Order("id ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// FindWithPage 分页查询配置
func (r *SystemConfigRepository) FindWithPage(req request.SystemConfigListRequest) ([]entity.SystemConfig, int64, error) {
	var configs []entity.SystemConfig
	var total int64

	query := r.DB.Model(&entity.SystemConfig{})

	// 条件过滤
	if req.ConfigKey != "" {
		query = query.Where("config_key LIKE ?", "%"+req.ConfigKey+"%")
	}
	if req.ConfigType != "" {
		query = query.Where("config_type = ?", req.ConfigType)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := req.GetOffset()
	if err := query.Offset(offset).Limit(req.GetPageSize()).Order("id ASC").Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

// FindByKey 根据配置键查询
func (r *SystemConfigRepository) FindByKey(key string) (*entity.SystemConfig, error) {
	var config entity.SystemConfig
	if err := r.DB.Where("config_key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// FindByKeys 根据多个配置键批量查询
func (r *SystemConfigRepository) FindByKeys(keys []string) ([]entity.SystemConfig, error) {
	var configs []entity.SystemConfig
	if err := r.DB.Where("config_key IN ?", keys).Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// Update 更新配置值
func (r *SystemConfigRepository) Update(id int64, configValue string, updatedBy int64) error {
	return r.DB.Model(&entity.SystemConfig{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"config_value": configValue,
			"updated_by":   updatedBy,
		}).Error
}

// BatchUpdate 批量更新配置值
func (r *SystemConfigRepository) BatchUpdate(updates []map[string]interface{}) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range updates {
			if err := tx.Model(&entity.SystemConfig{}).Where("id = ?", item["id"]).
				Updates(map[string]interface{}{
					"config_value": item["config_value"],
					"updated_by":   item["updated_by"],
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
