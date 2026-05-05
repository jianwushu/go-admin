package response

// MenuResponse 菜单响应
type MenuResponse struct {
	ID        int64          `json:"id"`                 // 菜单ID
	ParentID  int64          `json:"parentId"`           // 父菜单ID
	Name      string         `json:"name"`               // 路由名称
	I18nKey   string         `json:"i18nKey"`            // 国际化翻译key
	Path      string         `json:"path"`               // 路由路径
	Component string         `json:"component"`          // 组件路径
	Icon      string         `json:"icon"`               // 图标
	Type      int            `json:"type"`               // 类型：0=目录 1=菜单 2=按钮
	Sort      int            `json:"sort"`               // 排序
	Visible   int            `json:"visible"`            // 是否可见
	Status    int            `json:"status"`             // 状态：0=禁用 1=正常
	Perms     string         `json:"perms"`              // 权限标识
	Children  []MenuResponse `json:"children,omitempty"` // 子菜单
}
