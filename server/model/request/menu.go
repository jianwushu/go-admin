package request

// MenuCreateRequest 创建菜单请求
type MenuCreateRequest struct {
	ParentID  int64  `json:"parentId"`                    // 父菜单ID
	Name      string `json:"name" binding:"required,max=64"` // 路由名称
	Path      string `json:"path" binding:"max=128"`      // 路由路径
	Component string `json:"component" binding:"max=128"` // 组件路径
	Icon      string `json:"icon" binding:"max=64"`       // 图标
	Type      int    `json:"type" binding:"oneof=0 1 2"`  // 类型：0=目录 1=菜单 2=按钮
	Sort      int    `json:"sort"`                         // 排序
	Visible   int    `json:"visible" binding:"oneof=0 1"`  // 是否可见：0=隐藏 1=显示
	Status    int    `json:"status" binding:"oneof=0 1"`   // 状态：0=禁用 1=正常
	Perms     string `json:"perms" binding:"max=128"`      // 权限标识
}

// MenuUpdateRequest 更新菜单请求
type MenuUpdateRequest struct {
	ID        int64  `json:"id" binding:"required"`          // 菜单ID
	ParentID  int64  `json:"parentId"`                       // 父菜单ID
	Name      string `json:"name" binding:"required,max=64"` // 路由名称
	Path      string `json:"path" binding:"max=128"`        // 路由路径
	Component string `json:"component" binding:"max=128"`   // 组件路径
	Icon      string `json:"icon" binding:"max=64"`         // 图标
	Type      int    `json:"type" binding:"oneof=0 1 2"`    // 类型：0=目录 1=菜单 2=按钮
	Sort      int    `json:"sort"`                           // 排序
	Visible   int    `json:"visible" binding:"oneof=0 1"`    // 是否可见：0=隐藏 1=显示
	Status    int    `json:"status" binding:"oneof=0 1"`     // 状态：0=禁用 1=正常
	Perms     string `json:"perms" binding:"max=128"`        // 权限标识
}
