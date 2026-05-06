package initialize

import (
	"bytes"
	"fmt"
	"go-admin/config"
	"go-admin/global"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func DB() *gorm.DB {
	cfg := global.Config.Database

	switch cfg.Type {
	case "sqlite":
		return initSqlite(cfg.Sqlite)
	default:
		return initSqlite(cfg.Sqlite)
	}
}

func initSqlite(cfg config.Sqlite) *gorm.DB {
	// 确保 SQLite 数据库文件所在目录存在
	if err := ensureDir(cfg.Path); err != nil {
		panic(fmt.Sprintf("创建 SQLite 数据库目录失败: %v", err))
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Sprintf("连接 SQLite 数据库失败: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("获取数据库实例失败: %v", err))
	}

	// 连接池配置
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(time.Minute * 30) // 空闲连接最大生命周期

	// 执行建表和种子数据脚本
	initSQL(db)

	return db
}

// initSQL 执行建表和种子数据 SQL 脚本
func initSQL(db *gorm.DB) {
	// 检查是否需要初始化（通过检查 user 表是否存在来判断）
	tableName := fmt.Sprintf("%suser", global.Config.TablePrefix)
	if db.Migrator().HasTable(tableName) {
		global.Log.Info("数据库表已存在，跳过全量初始化", zap.String("table", tableName))
		// 增量迁移：检查并创建新增的表
		migrateNewTables(db)
		return
	}

	global.Log.Info("开始初始化数据库表结构和数据")

	// 执行建表脚本
	if err := executeSQLTemplate(db, "sql/init.sql"); err != nil {
		panic(fmt.Sprintf("执行建表脚本失败: %v", err))
	}
	global.Log.Info("建表脚本执行完成")

	// 执行种子数据脚本
	if err := executeSQLTemplate(db, "sql/seed.sql"); err != nil {
		panic(fmt.Sprintf("执行种子数据脚本失败: %v", err))
	}
	global.Log.Info("种子数据脚本执行完成")
}

// migrateNewTables 增量迁移：检查并创建新增的表
func migrateNewTables(db *gorm.DB) {
	newTables := []string{
		"job",
		"job_log",
		"system_config",
		"announcement",
		"message",
		"message_user",
	}
	for _, t := range newTables {
		fullName := fmt.Sprintf("%s%s", global.Config.TablePrefix, t)
		if !db.Migrator().HasTable(fullName) {
			global.Log.Info("发现新增表，开始创建", zap.String("table", fullName))
			if err := executeSQLTemplate(db, "sql/sqlite/init.sql"); err != nil {
				global.Log.Error("增量建表失败", zap.String("table", fullName), zap.Error(err))
			}
			// 一次执行整个 init.sql 即可，CREATE TABLE IF NOT EXISTS 会跳过已存在的表
			break
		}
	}
}

// executeSQLTemplate 读取 SQL 模板文件，渲染后执行
func executeSQLTemplate(db *gorm.DB, filePath string) error {
	// 读取 SQL 模板文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取 %s 失败: %v", filePath, err)
	}

	// 解析模板
	tmpl, err := template.New(filePath).Parse(string(content))
	if err != nil {
		return fmt.Errorf("解析 %s 模板失败: %v", filePath, err)
	}

	// 渲染模板
	var buf bytes.Buffer
	data := map[string]string{
		"TablePrefix": global.Config.TablePrefix,
	}
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染 %s 模板失败: %v", filePath, err)
	}

	// 执行 SQL
	if err := db.Exec(buf.String()).Error; err != nil {
		return fmt.Errorf("执行 %s 失败: %v", filePath, err)
	}

	return nil
}

// ensureDir 确保文件所在目录存在
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}
