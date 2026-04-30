package entity

// CodegenConfig 代码生成配置实体
type CodegenConfig struct {
	BaseModel
	TblName      string `json:"tableName" gorm:"column:table_name;size:128;not null;uniqueIndex"` // 数据库表名
	TableComment string `json:"tableComment" gorm:"size:256"`                                     // 表注释
	ClassName    string `json:"className" gorm:"size:128;not null"`                               // 大驼峰类名
	BusinessName string `json:"businessName" gorm:"size:128;not null"`                            // 业务名（小驼峰）
	FunctionName string `json:"functionName" gorm:"size:256"`                                     // 功能名（中文）
	ModuleName   string `json:"moduleName" gorm:"size:64;not null"`                               // 模块名
	PackageName  string `json:"packageName" gorm:"size:128;not null"`                             // 包名
	Author       string `json:"author" gorm:"size:64"`                                            // 作者
	Fields       string `json:"fields" gorm:"type:text"`                                          // 字段配置 JSON
}

// TableName 返回带前缀的表名
func (CodegenConfig) TableName() string {
	return TableName("codegen_config")
}
