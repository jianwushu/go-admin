package service

import (
	"context"
	"fmt"
	"go-admin/global"
	"go-admin/model/entity"
	"go-admin/model/response"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// OnlineService 在线用户服务
type OnlineService struct{}

// NewOnlineService 创建在线用户服务实例
func NewOnlineService() *OnlineService {
	return &OnlineService{}
}

// GetOnlineUsers 获取在线用户列表
func (s *OnlineService) GetOnlineUsers(username string) ([]response.OnlineUser, error) {
	ctx := context.Background()
	prefix := global.Config.TablePrefix

	// 使用 SCAN 扫描 jwt:token:* 的 key 来获取在线用户
	// key 格式: jwt:token:{userId}:{uuid}
	var onlineUsers []response.OnlineUser
	var cursor uint64

	for {
		keys, nextCursor, err := global.Redis.Scan(ctx, cursor, tokenKeyPrefix+"*", 100).Result()
		if err != nil {
			global.Log.Warn("SCAN Redis keys 失败", zap.Error(err))
			break
		}

		for _, key := range keys {
			// 解析 key: jwt:token:{userId}:{uuid}
			parts := strings.Split(key, ":")
			if len(parts) < 4 {
				continue
			}
			userIDStr := parts[2]
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				continue
			}

			// 查询用户信息
			var user entity.User
			if err := global.DB.Table(prefix+"user").Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
				continue
			}

			// 如果指定了用户名过滤
			if username != "" && !strings.Contains(user.Username, username) {
				continue
			}

			// 查询部门名称
			deptName := ""
			if user.DeptID > 0 {
				var dept entity.Dept
				if err := global.DB.Table(prefix+"dept").Where("id = ? AND deleted_at IS NULL", user.DeptID).First(&dept).Error; err == nil {
					deptName = dept.Name
				}
			}

			// 获取登录时间
			loginTimeKey := fmt.Sprintf("%s%d", loginTimeKeyPrefix, userID)
			loginTimeStr, err := global.Redis.Get(ctx, loginTimeKey).Result()
			var loginTime int64
			var onlineDuration string
			if err == nil {
				loginTime, _ = strconv.ParseInt(loginTimeStr, 10, 64)
				onlineDuration = formatDuration(time.Now().Unix() - loginTime)
			}

			// 获取客户端IP（从 Token 中无法直接获取，使用空值）
			ip := ""

			onlineUsers = append(onlineUsers, response.OnlineUser{
				UserID:         user.ID,
				Username:       user.Username,
				Nickname:       user.Nickname,
				DeptName:       deptName,
				IP:             ip,
				LoginTime:      loginTime,
				OnlineDuration: onlineDuration,
			})
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return onlineUsers, nil
}

// ForceLogout 强制用户下线
func (s *OnlineService) ForceLogout(userID int64) error {
	authService := NewAuthService()
	return authService.ForceLogout(userID)
}

// formatDuration 格式化在线时长
func formatDuration(seconds int64) string {
	if seconds < 60 {
		return fmt.Sprintf("%d秒", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%d分%d秒", seconds/60, seconds%60)
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	if hours < 24 {
		return fmt.Sprintf("%d小时%d分", hours, minutes)
	}
	days := hours / 24
	hours = hours % 24
	return fmt.Sprintf("%d天%d小时", days, hours)
}
