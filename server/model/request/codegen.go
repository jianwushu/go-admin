package request

// TableInfoRequest 表信息查询请求
type TableInfoRequest struct {
	PageRequest
	TableName string `json:"tableName" form:"tableName"` // 表名（模糊查询）
}

// ColumnConfig 字段配置
type ColumnConfig struct {
	ColumnName string `json:"columnName"` // 数据库列名
	ColumnType string `json:"columnType"` // 数据库类型
	GoType     string `json:"goType"`     // Go 类型
	GoField    string `json:"goField"`    // Go 字段名（驼峰）
	TsType     string `json:"tsType"`     // TS 类型
	Label      string `json:"label"`      // 显示标签
	HtmlType   string `json:"htmlType"`   // 表单类型
	QueryType  string `json:"queryType"`  // 查询方式
	IsList     bool   `json:"isList"`     // 是否列表显示
	IsQuery    bool   `json:"isQuery"`    // 是否查询条件
	IsRequired bool   `json:"isRequired"` // 是否必填
	IsEdit     bool   `json:"isEdit"`     // 是否可编辑
	DictType   string `json:"dictType"`   // 字典类型
	Comment    string `json:"comment"`    // 字段注释
	IsPk       bool   `json:"isPk"`       // 是否主键
	IsNull     bool   `json:"isNull"`     // 是否可空
	Sort       int    `json:"sort"`       // 排序
}

// CodegenPreviewRequest 代码预览请求
type CodegenPreviewRequest struct {
	TableName    string         `json:"tableName" binding:"required"`    // 数据库表名
	TableComment string         `json:"tableComment"`                    // 表注释
	ClassName    string         `json:"className" binding:"required"`    // 大驼峰类名
	BusinessName string         `json:"businessName" binding:"required"` // 业务名
	FunctionName string         `json:"functionName"`                    // 功能名
	ModuleName   string         `json:"moduleName" binding:"required"`   // 模块名
	PackageName  string         `json:"packageName" binding:"required"`  // 包名
	Author       string         `json:"author"`                          // 作者
	Fields       []ColumnConfig `json:"fields" binding:"required"`       // 字段配置列表
}

// CodegenGenerateRequest 代码生成请求（与预览请求相同）
type CodegenGenerateRequest = CodegenPreviewRequest

// CodegenSaveRequest 保存代码生成配置请求
type CodegenSaveRequest struct {
	TableName    string         `json:"tableName" binding:"required"`    // 数据库表名
	TableComment string         `json:"tableComment"`                    // 表注释
	ClassName    string         `json:"className" binding:"required"`    // 大驼峰类名
	BusinessName string         `json:"businessName" binding:"required"` // 业务名
	FunctionName string         `json:"functionName"`                    // 功能名
	ModuleName   string         `json:"moduleName" binding:"required"`   // 模块名
	PackageName  string         `json:"packageName" binding:"required"`  // 包名
	Author       string         `json:"author"`                          // 作者
	Fields       []ColumnConfig `json:"fields" binding:"required"`       // 字段配置列表
}
