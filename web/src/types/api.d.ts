/** API 统一响应结构 */
export interface ApiResponse<T = unknown> {
  /** 状态码：0=成功，非0=失败 */
  code: number
  /** 消息 */
  msg: string
  /** 数据 */
  data: T
}

/** 分页请求参数 */
export interface PageQuery {
  /** 当前页码 */
  page: number
  /** 每页条数 */
  pageSize: number
}

/** 分页响应结构（对应后端 PageResponse） */
export interface PageResponse<T = unknown> {
  /** 状态码：0=成功，非0=失败 */
  code: number
  /** 消息 */
  msg: string
  /** 数据列表 */
  data: T[]
  /** 总条数 */
  total: number
  /** 当前页码 */
  page: number
  /** 每页条数 */
  size: number
}

/** 分页响应数据（嵌套在 data 中的结构） */
export interface PageResult<T = unknown> {
  /** 数据列表 */
  list: T[]
  /** 总条数 */
  total: number
  /** 当前页码 */
  page: number
  /** 每页条数 */
  pageSize: number
}

/** 登录请求 */
export interface LoginRequest {
  username: string
  password: string
  captchaCode?: string
  captchaKey?: string
}

/** 登录响应 */
export interface LoginResult {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

/** 用户信息 */
export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  phone: string
  status: number
  deptId: number
  deptName: string
  roles: string[]
  permissions: string[]
}

/** 菜单项 */
export interface MenuItem {
  id: number
  parentId: number
  name: string
  path: string
  component: string
  icon: string
  type: number
  sort: number
  visible: number
  status: number
  perms: string
  children?: MenuItem[]
}

/** 角色简要信息 */
export interface RoleBrief {
  id: number
  name: string
  code: string
}

/** 用户管理 - 用户列表项 */
export interface UserItem {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  status: number
  deptId: number
  deptName: string
  avatar: string
  remark: string
  roles: RoleBrief[]
  createdAt: string
}

/** 用户管理 - 创建/更新请求 */
export interface UserFormData {
  id?: number
  username: string
  password?: string
  nickname: string
  email: string
  phone: string
  status: number
  deptId: number
  avatar: string
  remark: string
  roleIds: number[]
}

/** 用户管理 - 列表查询参数 */
export interface UserListParams extends PageQuery {
  username?: string
  status?: number
  deptId?: number
  phone?: string
}

/** 角色信息 */
export interface RoleInfo {
  id: number
  name: string
  code: string
  dataScope: number
  sort: number
  status: number
  remark: string
  menuIds: number[]
  deptIds: number[]
  createdAt: string
}

/** 角色管理 - 创建/更新请求 */
export interface RoleFormData {
  id?: number
  name: string
  code: string
  dataScope: number
  sort: number
  status: number
  remark: string
  menuIds: number[]
  deptIds: number[]
}

/** 角色管理 - 列表查询参数 */
export interface RoleListParams extends PageQuery {
  name?: string
  code?: string
  status?: number
}

/** 部门信息 */
export interface DeptInfo {
  id: number
  parentId: number
  name: string
  sort: number
  leader: string
  phone: string
  email: string
  status: number
  children?: DeptInfo[]
}

/** 部门管理 - 创建/更新请求 */
export interface DeptFormData {
  id?: number
  parentId: number
  name: string
  sort: number
  status: number
  leader: string
  phone: string
  email: string
}

/** 菜单管理 - 创建/更新请求 */
export interface MenuFormData {
  id?: number
  parentId: number
  name: string
  path: string
  component: string
  icon: string
  type: number
  sort: number
  visible: number
  status: number
  perms: string
}
