package repository

import (
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/response"
	"go-admin/utils"

	"gorm.io/gorm"
)

// CodegenRepository 代码生成仓储层
type CodegenRepository struct {
	BaseRepository
}

// NewCodegenRepository 创建代码生成仓储实例
func NewCodegenRepository() *CodegenRepository {
	return &CodegenRepository{
		BaseRepository: BaseRepository{DB: global.DB},
	}
}

// GetAllTables 获取所有用户表（排除系统内部表）
func (r *CodegenRepository) GetAllTables() ([]response.TableInfoResponse, error) {
	var tables []response.TableInfoResponse

	dbType := global.Config.Database.Type

	switch dbType {
	case "mysql":
		err := r.getMysqlTables(&tables)
		if err != nil {
			return nil, err
		}
	default:
		// 默认使用 SQLite 方式
		err := r.getSQLiteTables(&tables)
		if err != nil {
			return nil, err
		}
	}

	return tables, nil
}

// getSQLiteTables 获取 SQLite 表信息
func (r *CodegenRepository) getSQLiteTables(tables *[]response.TableInfoResponse) error {
	rows, err := r.DB.Raw(`
		SELECT name AS table_name, '' AS table_comment, 'sqlite' AS engine, '' AS create_time
		FROM sqlite_master
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%' AND name != 'codegen_config'
		ORDER BY name
	`).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var table response.TableInfoResponse
		if err := rows.Scan(&table.TableName, &table.TableComment, &table.Engine, &table.CreateTime); err != nil {
			return err
		}
		*tables = append(*tables, table)
	}
	return nil
}

// getMysqlTables 获取 MySQL 表信息
func (r *CodegenRepository) getMysqlTables(tables *[]response.TableInfoResponse) error {
	dbName := global.Config.Database.Mysql.DBName
	rows, err := r.DB.Raw(`
		SELECT table_name, table_comment, engine, create_time
		FROM information_schema.tables
		WHERE table_schema = ? AND table_type = 'BASE TABLE'
		ORDER BY create_time DESC
	`, dbName).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var table response.TableInfoResponse
		if err := rows.Scan(&table.TableName, &table.TableComment, &table.Engine, &table.CreateTime); err != nil {
			return err
		}
		*tables = append(*tables, table)
	}
	return nil
}

// GetColumnsByTableName 获取指定表的列信息
func (r *CodegenRepository) GetColumnsByTableName(tableName string) ([]response.ColumnInfoResponse, error) {
	var columns []response.ColumnInfoResponse

	dbType := global.Config.Database.Type

	switch dbType {
	case "mysql":
		err := r.getMysqlColumns(tableName, &columns)
		if err != nil {
			return nil, err
		}
	default:
		err := r.getSQLiteColumns(tableName, &columns)
		if err != nil {
			return nil, err
		}
	}

	return columns, nil
}

// getSQLiteColumns 获取 SQLite 列信息
func (r *CodegenRepository) getSQLiteColumns(tableName string, columns *[]response.ColumnInfoResponse) error {
	rows, err := r.DB.Raw("PRAGMA table_info(" + tableName + ")").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	idx := 0
	for rows.Next() {
		var cid int
		var columnName, columnType string
		var notNull int
		var dfltValue interface{}
		var pk int

		if err := rows.Scan(&cid, &columnName, &columnType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}

		goType := utils.MapSqliteTypeToGo(columnType)
		col := response.ColumnInfoResponse{
			ColumnName:    columnName,
			ColumnType:    columnType,
			ColumnComment: "",
			IsPk:          pk == 1,
			IsNull:        notNull == 0,
			GoType:        goType,
			GoField:       utils.SnakeToCamel(columnName),
			TsType:        utils.MapGoTypeToTs(goType),
			HtmlType:      utils.GuessHtmlType(columnName, goType),
			QueryType:     utils.GuessQueryType(columnName),
			IsList:        true,
			IsQuery:       false,
			IsRequired:    notNull == 1 && pk == 0,
			IsEdit:        pk == 0,
			Sort:          idx,
		}
		*columns = append(*columns, col)
		idx++
	}
	return nil
}

// getMysqlColumns 获取 MySQL 列信息
func (r *CodegenRepository) getMysqlColumns(tableName string, columns *[]response.ColumnInfoResponse) error {
	dbName := global.Config.Database.Mysql.DBName
	rows, err := r.DB.Raw(`
		SELECT column_name, column_type, column_comment, 
		       CASE WHEN column_key = 'PRI' THEN 1 ELSE 0 END AS is_pk,
		       CASE WHEN is_nullable = 'YES' THEN 1 ELSE 0 END AS is_nullable,
		       IFNULL(character_maximum_length, 0) AS max_length
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
		ORDER BY ordinal_position
	`, dbName, tableName).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	idx := 0
	for rows.Next() {
		var col response.ColumnInfoResponse
		var isPk, isNull int
		if err := rows.Scan(&col.ColumnName, &col.ColumnType, &col.ColumnComment, &isPk, &isNull, &col.MaxLength); err != nil {
			return err
		}
		col.IsPk = isPk == 1
		col.IsNull = isNull == 1
		col.GoType = utils.MapMysqlTypeToGo(col.ColumnType)
		col.GoField = utils.SnakeToCamel(col.ColumnName)
		col.TsType = utils.MapGoTypeToTs(col.GoType)
		col.HtmlType = utils.GuessHtmlType(col.ColumnName, col.GoType)
		col.QueryType = utils.GuessQueryType(col.ColumnName)
		col.IsList = true
		col.IsQuery = false
		col.IsRequired = !col.IsNull && !col.IsPk
		col.IsEdit = !col.IsPk
		col.Sort = idx
		*columns = append(*columns, col)
		idx++
	}
	return nil
}

// SaveConfig 保存代码生成配置
func (r *CodegenRepository) SaveConfig(config *entity.CodegenConfig) error {
	// 先查询是否已存在
	var existing entity.CodegenConfig
	err := r.DB.Where("table_name = ?", config.TblName).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		// 不存在，创建
		return r.DB.Create(config).Error
	}
	if err != nil {
		return err
	}
	// 已存在，更新
	config.ID = existing.ID
	return r.DB.Save(config).Error
}

// GetConfigByTableName 根据表名获取配置
func (r *CodegenRepository) GetConfigByTableName(tableName string) (*entity.CodegenConfig, error) {
	var config entity.CodegenConfig
	err := r.DB.Where("table_name = ?", tableName).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllConfigs 获取所有配置
func (r *CodegenRepository) GetAllConfigs() ([]entity.CodegenConfig, error) {
	var configs []entity.CodegenConfig
	err := r.DB.Order("id ASC").Find(&configs).Error
	return configs, err
}

// DeleteConfig 删除配置
func (r *CodegenRepository) DeleteConfig(id int64) error {
	return r.DB.Delete(&entity.CodegenConfig{}, id).Error
}
