# API 接口规范

## 通用规范

### 基础路径

```
http://localhost:8080/api/v1/
```

### 请求头

| Header | 说明 | 示例 |
|--------|------|------|
| Authorization | JWT Token | Bearer eyJhbGciOiJIUzI1NiIs... |
| Content-Type | 请求体类型 | application/json |
| Accept-Language | 语言偏好 | zh-CN / en |

### 数据安全规范

**密码处理：**
- 存储：使用 bcrypt 加密存储，永不存储明文密码
- API 响应：密码字段永不出现在任何响应中

**敏感数据脱敏规则：**

| 字段类型 | 脱敏规则 | 示例 |
|----------|----------|------|
| 手机号 | 保留前3位和后4位，中间用 **** 替代 | 138****8000 |
| 邮箱 | 用户名部分保留首尾字符，中间用 *** 替代 | z***n@example.com |
| 身份证 | 保留前4位和后4位 | 3101**********1234 |
| 银行卡 | 仅显示后4位 | ****1234 |

**脱敏实现位置：** 后端 Service 层返回数据前统一脱敏，通过 `utils/desensitize.go` 工具函数处理。

### 统一响应格式

**成功响应：**
```json
{
  "code": 0,
  "msg": "操作成功",
  "data": { ... }
}
```

**分页响应：**
```json
{
  "code": 0,
  "msg": "查询成功",
  "data": [ ... ],
  "total": 100,
  "page": 1,
  "size": 10
}
```

**错误响应：**
```json
{
  "code": 2001,
  "msg": "参数校验失败",
  "data": null
}
```

### 分页参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码 |
| pageSize | int | 10 | 每页大小（最大100） |

---

## 接口清单

### 一、认证模块 `/api/v1/auth`

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | /auth/login | 用户登录 | 公开 |
| POST | /auth/logout | 用户登出 | 需认证 |
| POST | /auth/refresh | 刷新Token | 需认证 |
| GET | /auth/captcha | 获取验证码 | 公开 |

**登录请求：**
```json
POST /api/v1/auth/login
{
  "username": "admin",
  "password": "admin123",
  "captcha_id": "xxx",
  "captcha_code": "1234"
}
```

**登录响应：**
```json
{
  "code": 0,
  "msg": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

---

### 二、用户模块 `/api/v1/users`

| 方法 | 路径 | 说明 | 权限标识 |
|------|------|------|----------|
| GET | /users | 用户列表（分页） | system:user:list |
| GET | /users/:id | 用户详情 | system:user:list |
| POST | /users | 新增用户 | system:user:add |
| PUT | /users/:id | 修改用户 | system:user:edit |
| DELETE | /users/:id | 删除用户 | system:user:delete |
| PUT | /users/:id/status | 修改用户状态 | system:user:edit |
| PUT | /users/:id/password | 重置密码 | system:user:edit |
| GET | /user/info | 获取当前用户信息 | 需认证 |
| PUT | /user/info | 修改个人信息 | 需认证 |
| PUT | /user/password | 修改个人密码 | 需认证 |

**用户列表请求：**
```
GET /api/v1/users?page=1&pageSize=10&username=admin&status=1&deptId=1
```

**用户列表响应：**
```json
{
  "code": 0,
  "msg": "查询成功",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "nickname": "超级管理员",
      "email": "admin@example.com",
      "phone": "13800138000",
      "status": 1,
      "deptId": 1,
      "deptName": "总公司",
      "roles": [
        { "id": 1, "name": "超级管理员", "code": "admin" }
      ],
      "createdAt": "2024-01-01 00:00:00"
    }
  ],
  "total": 1,
  "page": 1,
  "size": 10
}
```

**新增用户请求：**
```json
POST /api/v1/users
{
  "username": "zhangsan",
  "password": "123456",
  "nickname": "张三",
  "email": "zhangsan@example.com",
  "phone": "13900139000",
  "status": 1,
  "deptId": 2,
  "roleIds": [2],
  "remark": "备注"
}
```

---

### 三、角色模块 `/api/v1/roles`

| 方法 | 路径 | 说明 | 权限标识 |
|------|------|------|----------|
| GET | /roles | 角色列表（分页） | system:role:list |
| GET | /roles/:id | 角色详情 | system:role:list |
| GET | /roles/all | 所有角色（下拉选择用） | 需认证 |
| POST | /roles | 新增角色 | system:role:add |
| PUT | /roles/:id | 修改角色 | system:role:edit |
| DELETE | /roles/:id | 删除角色 | system:role:delete |
| PUT | /roles/:id/status | 修改角色状态 | system:role:edit |

**新增角色请求（含菜单权限和数据权限）：**
```json
POST /api/v1/roles
{
  "name": "编辑员",
  "code": "editor",
  "dataScope": 2,
  "sort": 1,
  "status": 1,
  "menuIds": [1, 2, 3, 10, 11, 12, 13],
  "deptIds": [],
  "remark": "编辑员角色"
}
```

---

### 四、菜单模块 `/api/v1/menus`

| 方法 | 路径 | 说明 | 权限标识 |
|------|------|------|----------|
| GET | /menus | 菜单列表（树形） | system:menu:list |
| GET | /menus/:id | 菜单详情 | system:menu:list |
| GET | /menus/tree | 菜单树（角色授权用） | 需认证 |
| POST | /menus | 新增菜单 | system:menu:add |
| PUT | /menus/:id | 修改菜单 | system:menu:edit |
| DELETE | /menus/:id | 删除菜单 | system:menu:delete |

**菜单树响应：**
```json
{
  "code": 0,
  "msg": "查询成功",
  "data": [
    {
      "id": 1,
      "parentId": 0,
      "name": "系统管理",
      "path": "/system",
      "icon": "setting",
      "type": 0,
      "sort": 1,
      "children": [
        {
          "id": 2,
          "parentId": 1,
          "name": "用户管理",
          "path": "user",
          "component": "system/user/index",
          "icon": "user",
          "type": 1,
          "sort": 1,
          "children": [
            {
              "id": 10,
              "parentId": 2,
              "name": "用户查询",
              "perms": "system:user:list",
              "type": 2,
              "sort": 1
            }
          ]
        }
      ]
    }
  ]
}
```

---

### 五、部门模块 `/api/v1/depts`

| 方法 | 路径 | 说明 | 权限标识 |
|------|------|------|----------|
| GET | /depts | 部门列表（树形） | system:dept:list |
| GET | /depts/:id | 部门详情 | system:dept:list |
| GET | /depts/tree | 部门树（选择用） | 需认证 |
| POST | /depts | 新增部门 | system:dept:add |
| PUT | /depts/:id | 修改部门 | system:dept:edit |
| DELETE | /depts/:id | 删除部门 | system:dept:delete |

---

### 六、系统监控 `/api/v1/monitor`

| 方法 | 路径 | 说明 | 权限标识 |
|------|------|------|----------|
| GET | /monitor/server | 服务器信息 | monitor:server:list |
| GET | /monitor/operation-logs | 操作日志列表 | monitor:log:list |
| GET | /monitor/login-logs | 登录日志列表 | monitor:log:list |
| DELETE | /monitor/operation-logs | 清空操作日志 | monitor:log:delete |
| DELETE | /monitor/login-logs | 清空登录日志 | monitor:log:delete |

**服务器监控响应：**
```json
{
  "code": 0,
  "msg": "查询成功",
  "data": {
    "cpu": {
      "cores": 8,
      "usage": 23.5
    },
    "memory": {
      "total": 16384,
      "used": 8192,
      "usage": 50.0
    },
    "disk": {
      "total": 512000,
      "used": 256000,
      "usage": 50.0
    },
    "go": {
      "version": "1.22.0",
      "goroutines": 42,
      "allocMemory": 10240000,
      "gcCount": 5
    },
    "db": {
      "type": "sqlite",
      "status": "connected"
    },
    "redis": {
      "status": "connected",
      "usedMemory": "1.2M"
    }
  }
}
```

---

### 七、代码生成 `/api/v1/codegen`

| 方法 | 路径 | 说明 | 权限标识 |
|------|------|------|----------|
| GET | /codegen/tables | 数据库表列表 | codegen:list |
| GET | /codegen/tables/:name | 表字段信息 | codegen:list |
| POST | /codegen/preview | 预览生成代码 | codegen:preview |
| POST | /codegen/generate | 生成并下载代码 | codegen:generate |
| POST | /codegen/configs | 保存生成配置 | codegen:save |
| GET | /codegen/configs/:id | 获取生成配置 | codegen:list |

**代码生成请求：**
```json
POST /api/v1/codegen/generate
{
  "tableName": "sys_user",
  "tableComment": "用户表",
  "moduleName": "system",
  "businessName": "user",
  "functionName": "用户管理",
  "packageName": "go-admin",
  "fields": [
    {
      "columnName": "username",
      "columnType": "varchar",
      "goType": "string",
      "goField": "Username",
      "htmlType": "input",
      "isRequired": true,
      "isList": true,
      "isQuery": true,
      "queryType": "LIKE",
      "label": "用户名"
    }
  ]
}
```

---

## 错误码汇总

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | Token 已过期 |
| 1002 | Token 无效 |
| 1003 | 无操作权限 |
| 1004 | 账号已被禁用 |
| 1005 | 用户名或密码错误 |
| 1006 | 验证码错误 |
| 2001 | 参数校验失败 |
| 2002 | 请求方法不允许 |
| 3001 | 用户已存在 |
| 3002 | 角色已存在 |
| 3003 | 部门下存在子部门，无法删除 |
| 3004 | 角色下存在用户，无法删除 |
| 3005 | 菜单下存在子菜单，无法删除 |
| 5001 | 数据库错误 |
| 5002 | Redis 错误 |
| 5003 | 系统内部错误 |
