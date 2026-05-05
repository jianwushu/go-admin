package response

// TableInfoResponse 数据库表信息响应
type TableInfoResponse struct {
	TableName    string `json:"tableName"`    // 表名
	TableComment string `json:"tableComment"` // 表注释
	Engine       string `json:"engine"`       // 存储引擎
	CreateTime   string `json:"createTime"`   // 创建时间
}

// ColumnInfoResponse 数据库列信息响应
type ColumnInfoResponse struct {
	ColumnName    string `json:"columnName"`    // 列名
	ColumnType    string `json:"columnType"`    // 列类型
	ColumnComment string `json:"columnComment"` // 列注释
	IsPk          bool   `json:"isPk"`          // 是否主键
	IsNull        bool   `json:"isNull"`        // 是否可空
	MaxLength     int    `json:"maxLength"`     // 最大长度
	GoType        string `json:"goType"`        // Go 类型
	GoField       string `json:"goField"`       // Go 字段名
	TsType        string `json:"tsType"`        // TS 类型
	HtmlType      string `json:"htmlType"`      // 表单类型
	QueryType     string `json:"queryType"`     // 查询方式
	IsList        bool   `json:"isList"`        // 是否列表显示
	IsQuery       bool   `json:"isQuery"`       // 是否查询条件
	IsRequired    bool   `json:"isRequired"`    // 是否必填
	IsEdit        bool   `json:"isEdit"`        // 是否可编辑
	DictType      string `json:"dictType"`      // 字典类型
	Sort          int    `json:"sort"`          // 排序
}

// CodegenPreviewResponse 代码预览响应
type CodegenPreviewResponse struct {
	Files []CodegenFile `json:"files"` // 生成的文件列表
}

// CodegenFile 代码文件
type CodegenFile struct {
	FileName string `json:"fileName"` // 文件名
	FilePath string `json:"filePath"` // 文件路径
	Content  string `json:"content"`  // 文件内容
}

// CodegenConfigResponse 代码生成配置响应
type CodegenConfigResponse struct {
	ID           int64          `json:"id"`           // 配置ID
	TableName    string         `json:"tableName"`    // 表名
	TableComment string         `json:"tableComment"` // 表注释
	ClassName    string         `json:"className"`    // 类名
	BusinessName string         `json:"businessName"` // 业务名
	FunctionName string         `json:"functionName"` // 功能名
	ModuleName   string         `json:"moduleName"`   // 模块名
	PackageName  string         `json:"packageName"`  // 包名
	Author       string         `json:"author"`       // 作者
	Fields       []ColumnConfig `json:"fields"`       // 字段配置
	CreatedAt    string         `json:"createdAt"`    // 创建时间
}

// ColumnConfig 字段配置（用于配置响应）
type ColumnConfig struct {
	ColumnName string `json:"columnName"` // 数据库列名
	ColumnType string `json:"columnType"` // 数据库类型
	GoType     string `json:"goType"`     // Go 类型
	GoField    string `json:"goField"`    // Go 字段名
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
	MaxLength  int    `json:"maxLength"`  // 最大长度（用于 VARCHAR）
	Sort       int    `json:"sort"`       // 排序
}
