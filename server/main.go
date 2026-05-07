package main

import (
	"fmt"
	"go-admin/cron"
	"go-admin/global"
	"go-admin/initialize"

	"go.uber.org/zap"
)

func main() {
	// 1. 初始化配置
	initialize.Viper()

	// 2. 初始化日志
	global.Log = initialize.Logger()
	global.Log.Info("日志系统初始化完成")

	// 3. 初始化数据库
	global.DB = initialize.DB()
	global.Log.Info("数据库初始化完成")

	// 4. 初始化 Redis
	global.Redis = initialize.Redis()
	global.Log.Info("Redis 初始化完成")

	// 5. 初始化定时任务管理器
	cron.Init()
	cron.InitJobs()
	global.Log.Info("定时任务管理器初始化完成")

	// 6. 初始化路由
	r := initialize.Router()
	global.Log.Info("路由初始化完成")

	// 7. 启动服务
	addr := fmt.Sprintf(":%d", global.Config.Server.Port)
	global.Log.Info("服务启动", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		global.Log.Fatal("服务启动失败", zap.Error(err))
	}
}
