package config

// Server 服务器配置
type Server struct {
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"` // debug / release / test
}

// Database 数据库配置
type Database struct {
	Type     string   `mapstructure:"type" json:"type" yaml:"type"` // sqlite / mysql / postgres / oracle
	Sqlite   Sqlite   `mapstructure:"sqlite" json:"sqlite" yaml:"sqlite"`
	Mysql    Mysql    `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Postgres Postgres `mapstructure:"postgres" json:"postgres" yaml:"postgres"`
	Oracle   Oracle   `mapstructure:"oracle" json:"oracle" yaml:"oracle"`
}

type Sqlite struct {
	Path string `mapstructure:"path" json:"path" yaml:"path"`
}

type Mysql struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
}

type Postgres struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
}

type Oracle struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
}

// Redis Redis配置
type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

// JWT JWT配置
type JWT struct {
	Secret  string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Expire  int    `mapstructure:"expire" json:"expire" yaml:"expire"`    // Token 过期时间（秒）
	Refresh int    `mapstructure:"refresh" json:"refresh" yaml:"refresh"` // 刷新过期时间（秒）
}

// Log 日志配置
type Log struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`             // debug / info / warn / error
	Path       string `mapstructure:"path" json:"path" yaml:"path"`                // 日志文件路径
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`    // 单文件最大 MB
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"` // 最大备份数
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`       // 最大保留天数
}

// Codegen 代码生成配置
type Codegen struct {
	Output string `mapstructure:"output" json:"output" yaml:"output"` // 代码生成输出目录
}

// Config 总配置
type Config struct {
	Server      Server   `mapstructure:"server" json:"server" yaml:"server"`
	TablePrefix string   `mapstructure:"table_prefix" json:"table_prefix" yaml:"table_prefix"` // 数据库表前缀
	Database    Database `mapstructure:"database" json:"database" yaml:"database"`
	Redis       Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	JWT         JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Log         Log      `mapstructure:"log" json:"log" yaml:"log"`
	Codegen     Codegen  `mapstructure:"codegen" json:"codegen" yaml:"codegen"`
}
