package utils

import (
	"strings"
	"unicode"
)

// snakeToCamel 蛇形命名转大驼峰
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// CamelToSnake 大驼峰转蛇形命名
func CamelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// FirstLower 首字母小写
func FirstLower(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// FirstUpper 首字母大写
func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// MapSqliteTypeToGo SQLite 类型映射到 Go 类型
func MapSqliteTypeToGo(columnType string) string {
	ct := strings.ToUpper(columnType)
	switch {
	case strings.Contains(ct, "INT"):
		return "int64"
	case strings.Contains(ct, "REAL") || strings.Contains(ct, "FLOAT") || strings.Contains(ct, "DOUBLE"):
		return "float64"
	case strings.Contains(ct, "BOOL"):
		return "bool"
	case strings.Contains(ct, "BLOB"):
		return "[]byte"
	case strings.Contains(ct, "DATETIME") || strings.Contains(ct, "TIMESTAMP") || strings.Contains(ct, "DATE"):
		return "time.Time"
	default:
		return "string"
	}
}

// MapMysqlTypeToGo MySQL 类型映射到 Go 类型
func MapMysqlTypeToGo(columnType string) string {
	ct := strings.ToUpper(columnType)
	switch {
	case strings.Contains(ct, "BIGINT"):
		return "int64"
	case strings.Contains(ct, "INT") || strings.Contains(ct, "TINYINT") || strings.Contains(ct, "SMALLINT") || strings.Contains(ct, "MEDIUMINT"):
		return "int"
	case strings.Contains(ct, "FLOAT"):
		return "float32"
	case strings.Contains(ct, "DOUBLE") || strings.Contains(ct, "DECIMAL"):
		return "float64"
	case strings.Contains(ct, "BOOL"):
		return "bool"
	case strings.Contains(ct, "BLOB") || strings.Contains(ct, "BINARY"):
		return "[]byte"
	case strings.Contains(ct, "DATETIME") || strings.Contains(ct, "TIMESTAMP") || strings.Contains(ct, "DATE"):
		return "time.Time"
	default:
		return "string"
	}
}

// MapGoTypeToTs Go 类型映射到 TypeScript 类型
func MapGoTypeToTs(goType string) string {
	switch goType {
	case "int", "int32", "int64", "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	default:
		return "string"
	}
}

// GuessHtmlType 根据字段名和类型猜测表单控件类型
func GuessHtmlType(columnName, goType string) string {
	name := strings.ToLower(columnName)

	// 特殊字段名匹配
	switch {
	case strings.HasSuffix(name, "_id"):
		return "select"
	case name == "status":
		return "select"
	case name == "type" || name == "type_id":
		return "select"
	case strings.Contains(name, "remark") || strings.Contains(name, "desc") || strings.Contains(name, "content") || strings.Contains(name, "note"):
		return "textarea"
	case strings.Contains(name, "email"):
		return "input"
	case strings.Contains(name, "phone") || strings.Contains(name, "mobile"):
		return "input"
	case strings.Contains(name, "url") || strings.Contains(name, "link"):
		return "input"
	case strings.Contains(name, "image") || strings.Contains(name, "avatar") || strings.Contains(name, "icon") || strings.Contains(name, "logo"):
		return "upload"
	case strings.Contains(name, "date") && !strings.Contains(name, "datetime"):
		return "date"
	case strings.Contains(name, "time") || strings.Contains(name, "datetime"):
		return "datetime"
	case strings.Contains(name, "sort") || strings.Contains(name, "order"):
		return "input"
	case strings.Contains(name, "password") || strings.Contains(name, "pwd"):
		return "input"
	case goType == "bool":
		return "radio"
	case goType == "time.Time":
		return "datetime"
	default:
		return "input"
	}
}

// GuessQueryType 根据字段名猜测查询方式
func GuessQueryType(columnName string) string {
	name := strings.ToLower(columnName)
	switch {
	case strings.Contains(name, "status") || strings.HasSuffix(name, "_id") || strings.Contains(name, "type"):
		return "="
	case strings.Contains(name, "name") || strings.Contains(name, "title") || strings.Contains(name, "username") || strings.Contains(name, "nickname"):
		return "LIKE"
	case strings.Contains(name, "date") || strings.Contains(name, "time"):
		return "BETWEEN"
	default:
		return "="
	}
}

// MapSqliteTypeToGo 导出版本（供 repository 使用）
func MapSqliteTypeToGoExport(columnType string) string {
	return MapSqliteTypeToGo(columnType)
}

// MapMysqlTypeToGo 导出版本
func MapMysqlTypeToGoExport(columnType string) string {
	return MapMysqlTypeToGo(columnType)
}

// MapGoTypeToTs 导出版本
func MapGoTypeToTsExport(goType string) string {
	return MapGoTypeToTs(goType)
}

// GuessHtmlType 导出版本
func GuessHtmlTypeExport(columnName, goType string) string {
	return GuessHtmlType(columnName, goType)
}

// GuessQueryType 导出版本
func GuessQueryTypeExport(columnName string) string {
	return GuessQueryType(columnName)
}
