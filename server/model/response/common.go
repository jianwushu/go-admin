package response

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"` // 状态码：0=成功，非0=失败
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

// PageResponse 分页响应结构
type PageResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Total int64       `json:"total"` // 总记录数
	Page  int         `json:"page"`  // 当前页
	Size  int         `json:"size"`  // 每页大小
}
