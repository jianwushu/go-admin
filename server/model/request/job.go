package request

// JobListRequest 定时任务列表查询参数
type JobListRequest struct {
	PageRequest
	Name     string `json:"name" form:"name"`         // 任务名称
	JobType  int    `json:"jobType" form:"jobType"`    // 任务类型
	Status   int    `json:"status" form:"status"`      // 状态
}

// JobCreateRequest 创建定时任务请求
type JobCreateRequest struct {
	Name       string `json:"name" binding:"required,max=128"`       // 任务名称
	JobType    int    `json:"jobType" binding:"required,oneof=1 2"`  // 1=内置函数 2=HTTP请求
	CronExpr   string `json:"cronExpr" binding:"required,max=64"`    // cron表达式
	FuncName   string `json:"funcName" binding:"max=128"`            // 内置函数名
	HTTPUrl    string `json:"httpUrl" binding:"max=512"`             // HTTP地址
	HTTPMethod string `json:"httpMethod" binding:"max=10"`           // GET/POST
	HTTPBody   string `json:"httpBody"`                              // 请求体
	Status     int    `json:"status"`                                // 1=启用 0=禁用
	Remark     string `json:"remark" binding:"max=512"`              // 备注
}

// JobUpdateRequest 更新定时任务请求
type JobUpdateRequest struct {
	ID         int64  `json:"id" binding:"required"`                  // 任务ID
	Name       string `json:"name" binding:"required,max=128"`       // 任务名称
	JobType    int    `json:"jobType" binding:"required,oneof=1 2"`  // 1=内置函数 2=HTTP请求
	CronExpr   string `json:"cronExpr" binding:"required,max=64"`    // cron表达式
	FuncName   string `json:"funcName" binding:"max=128"`            // 内置函数名
	HTTPUrl    string `json:"httpUrl" binding:"max=512"`             // HTTP地址
	HTTPMethod string `json:"httpMethod" binding:"max=10"`           // GET/POST
	HTTPBody   string `json:"httpBody"`                              // 请求体
	Status     int    `json:"status"`                                // 1=启用 0=禁用
	Remark     string `json:"remark" binding:"max=512"`              // 备注
}

// JobChangeStatusRequest 修改任务状态请求
type JobChangeStatusRequest struct {
	ID     int64 `json:"id" binding:"required"`      // 任务ID
	Status int   `json:"status" binding:"required"`  // 1=启用 0=禁用
}

// JobLogListRequest 任务日志列表查询参数
type JobLogListRequest struct {
	PageRequest
	JobID  int64  `json:"jobId" form:"jobId"`   // 任务ID
	Status int    `json:"status" form:"status"` // 状态
}
