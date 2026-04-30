package request

// PageRequest 分页请求参数
type PageRequest struct {
	Page     int `json:"page" form:"page"`         // 页码，默认1
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小，默认10
}

// GetPage 获取页码（带默认值）
func (p *PageRequest) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetPageSize 获取每页大小（带默认值和上限）
func (p *PageRequest) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	if p.PageSize > 100 {
		return 100
	}
	return p.PageSize
}

// GetOffset 获取分页偏移量
func (p *PageRequest) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPageSize()
}
