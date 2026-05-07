package response

// JobResponse 定时任务响应
type JobResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	JobType    int    `json:"jobType"`
	CronExpr   string `json:"cronExpr"`
	FuncName   string `json:"funcName"`
	HTTPUrl    string `json:"httpUrl"`
	HTTPMethod string `json:"httpMethod"`
	HTTPBody   string `json:"httpBody"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
	CreatedBy  int64  `json:"createdBy"`
	UpdatedBy  int64  `json:"updatedBy"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// JobLogResponse 任务执行日志响应
type JobLogResponse struct {
	ID        int64  `json:"id"`
	JobID     int64  `json:"jobId"`
	JobName   string `json:"jobName"`
	Status    int    `json:"status"`
	Result    string `json:"result"`
	ErrorMsg  string `json:"errorMsg"`
	Duration  int64  `json:"duration"`
	CreatedAt string `json:"createdAt"`
}
