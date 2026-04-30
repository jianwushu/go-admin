package initialize

import (
	"fmt"
	"go-admin/global"

	"github.com/spf13/viper"
)

func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("config/config.yaml")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("读取配置文件失败: %v", err))
	}

	if err := v.Unmarshal(&global.Config); err != nil {
		panic(fmt.Sprintf("解析配置文件失败: %v", err))
	}

	return v
}
