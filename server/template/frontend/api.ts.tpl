import request from '@/utils/request'

// 获取{{.FunctionName}}列表
export function get{{.ClassName}}List(params: any) {
  return request({
    url: '/api/v1/{{.ModuleName}}/{{.BusinessName}}/list',
    method: 'get',
    params,
  })
}

// 根据ID获取{{.FunctionName}}
export function get{{.ClassName}}ById(id: number) {
  return request({
    url: `/api/v1/{{.ModuleName}}/{{.BusinessName}}/${id}`,
    method: 'get',
  })
}

// 创建{{.FunctionName}}
export function create{{.ClassName}}(data: any) {
  return request({
    url: '/api/v1/{{.ModuleName}}/{{.BusinessName}}',
    method: 'post',
    data,
  })
}

// 更新{{.FunctionName}}
export function update{{.ClassName}}(data: any) {
  return request({
    url: '/api/v1/{{.ModuleName}}/{{.BusinessName}}',
    method: 'put',
    data,
  })
}

// 删除{{.FunctionName}}
export function delete{{.ClassName}}(id: number) {
  return request({
    url: `/api/v1/{{.ModuleName}}/{{.BusinessName}}/${id}`,
    method: 'delete',
  })
}
