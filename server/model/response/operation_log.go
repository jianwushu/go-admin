package response

// OperationLogResponse 操作日志响应
type OperationLogResponse struct {
	ID           int64  `json:"id"`
	Module       string `json:"module"`
	Action       string `json:"action"`
	Method       string `json:"method"`
	URL          string `json:"url"`
	IP           string `json:"ip"`
	Operator     string `json:"operator"`
	RequestParam string `json:"requestParam"`
	ResponseData string `json:"responseData"`
	Status       int    `json:"status"`
	ErrorMsg     string `json:"errorMsg"`
	Duration     int64  `json:"duration"`
	CreatedAt    string `json:"createdAt"`
}
