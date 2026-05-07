# Go-Admin

一个基于 **Go + Vue 3** 的全栈后台管理系统，采用前后端分离架构，开箱即用。

## ✨ 功能特性

### 系统管理
- **用户管理** — 用户的增删改查、分页查询、状态启停、角色分配
- **角色管理** — 角色 CRUD、菜单权限分配、数据权限范围控制
- **菜单管理** — 树形菜单 CRUD，支持目录 / 菜单 / 按钮三种类型
- **部门管理** — 树形部门 CRUD，支持层级组织架构
- **系统配置** — 键值对形式的系统参数配置（标题、Logo、版权等）

### 系统监控
- **仪表盘** — 用户数、角色数、菜单数、部门数、今日登录等统计概览
- **服务器监控** — CPU、内存、磁盘、Go 运行时、数据库连接、Redis 状态
- **在线用户** — 查看当前在线用户，支持强制下线
- **操作日志** — 自动记录 API 操作日志（模块、动作、耗时、请求/响应）
- **登录日志** — 记录用户登录历史（IP、浏览器、操作系统、状态）

### 系统工具
- **代码生成** — 从数据库表自动生成前后端 CRUD 代码（Go + Vue），支持预览和 ZIP 下载
- **定时任务** — Cron 定时任务管理，支持函数任务和 HTTP 任务，含执行日志

### 通用能力
- **RBAC 权限控制** — 基于角色的访问控制，支持菜单权限和数据权限（全部 / 本部门 / 本部门及下级 / 自定义 / 仅本人）
- **JWT 认证** — Token 认证与自动刷新机制
- **国际化 (i18n)** — 支持中文 / 英文切换
- **主题切换** — 亮色 / 暗色主题
- **布局切换** — 侧边栏 / 顶部菜单布局
- **响应式设计** — 适配桌面端和移动端

## 🛠️ 技术栈

### 后端

| 技术 | 说明 |
|------|------|
| [Go](https://go.dev/) 1.25 | 编程语言 |
| [Gin](https://github.com/gin-gonic/gin) | Web 框架 |
| [GORM](https://gorm.io/) | ORM 框架 |
| [SQLite](https://www.sqlite.org/) / MySQL / PostgreSQL / Oracle | 数据库（默认 SQLite） |
| [Redis](https://redis.io/) | 缓存、在线用户管理 |
| [JWT](https://github.com/golang-jwt/jwt) | 身份认证 |
| [Viper](https://github.com/spf13/viper) | 配置管理 |
| [Zap](https://go.uber.org/zap) + Lumberjack | 日志系统 |
| [gopsutil](https://github.com/shirou/gopsutil) | 服务器监控 |

### 前端

| 技术 | 说明 |
|------|------|
| [Vue 3](https://vuejs.org/) | 前端框架 |
| [TypeScript](https://www.typescriptlang.org/) | 类型安全 |
| [Vite 8](https://vitejs.dev/) | 构建工具 |
| [Element Plus](https://element-plus.org/) | UI 组件库 |
| [Pinia](https://pinia.vuejs.org/) | 状态管理 |
| [Vue Router 4](https://router.vuejs.org/) | 路由管理 |
| [Axios](https://axios-http.com/) | HTTP 请求 |
| [Tailwind CSS 4](https://tailwindcss.com/) | 原子化 CSS |
| [vue-i18n](https://vue-i18n.intlify.dev/) | 国际化 |

## 📁 项目结构

```
go-admin/
├── server/                     # 后端服务
│   ├── main.go                 # 入口文件
│   ├── config/                 # 配置文件
│   │   ├── config.go           # 配置结构体
│   │   └── config.yaml         # 配置文件
│   ├── controller/             # 控制器层
│   ├── service/                # 业务逻辑层
│   ├── repository/             # 数据访问层
│   ├── model/                  # 数据模型
│   │   ├── request/            # 请求参数结构体
│   │   └── response/           # 响应数据结构体
│   ├── router/                 # 路由定义
│   ├── middleware/              # 中间件
│   │   ├── auth.go             # JWT 认证
│   │   ├── permission.go       # 权限校验
│   │   ├── data_scope.go       # 数据权限
│   │   ├── operation_log.go    # 操作日志
│   │   ├── cors.go             # 跨域处理
│   │   └── recovery.go         # 异常恢复
│   ├── initialize/             # 初始化
│   ├── cron/                   # 定时任务
│   ├── utils/                  # 工具函数
│   ├── template/               # 代码生成模板
│   │   ├── backend/            # 后端代码模板
│   │   └── frontend/           # 前端代码模板
│   ├── sql/                    # 数据库脚本
│   │   ├── init.sql            # 建表脚本
│   │   └── seed.sql            # 初始数据
│   └── global/                 # 全局变量
├── web/                        # 前端项目
│   ├── src/
│   │   ├── views/              # 页面组件
│   │   ├── layout/             # 布局组件
│   │   ├── router/             # 路由配置
│   │   ├── store/              # 状态管理
│   │   ├── hooks/              # 组合式函数
│   │   ├── i18n/               # 国际化
│   │   ├── utils/              # 工具函数
│   │   └── types/              # 类型定义
│   └── package.json
└── plans/                      # 项目设计文档
```

## 🚀 快速开始

### 环境要求

- **Go** >= 1.25
- **Node.js** >= 18
- **Redis** （可选，用于在线用户管理等功能）

### 后端启动

```bash
# 进入后端目录
cd server

# 安装依赖
go mod tidy

# 启动服务（默认端口 8080，使用 SQLite 无需额外安装数据库）
go run main.go
```

首次启动会自动：
- 创建 SQLite 数据库文件 `data/go-admin.db`
- 执行建表脚本和初始数据导入
- 创建默认管理员账号：`admin` / `admin123`

### 前端启动

```bash
# 进入前端目录
cd web

# 安装依赖
npm install

# 启动开发服务器（默认端口 5173）
npm run dev
```

访问 http://localhost:5173，使用 `admin` / `admin123` 登录。

### 生产构建

```bash
# 构建前端
cd web && npm run build

# 构建后端
cd server && go build -o go-admin main.go
```

## ⚙️ 配置说明

后端配置文件位于 [`server/config/config.yaml`](server/config/config.yaml)，主要配置项：

```yaml
server:
  port: 8080                    # 服务端口
  mode: debug                   # 运行模式：debug / release / test

table_prefix: xy_               # 数据库表前缀

database:
  type: sqlite                  # 数据库类型：sqlite / mysql / postgres / oracle
  sqlite:
    path: ./data/go-admin.db    # SQLite 文件路径

redis:
  addr: 127.0.0.1:6379
  password: "123456"
  db: 0

jwt:
  secret: go-admin-secret-key   # JWT 密钥（生产环境务必修改）
  expire: 7200                  # Token 过期时间（秒）
  refresh: 604800               # 刷新过期时间（秒）

log:
  level: info                   # 日志级别
  path: ./logs                  # 日志目录
```

## 📄 API 接口

所有接口统一前缀 `/api/v1`，主要模块：

| 模块 | 路径前缀 | 说明 |
|------|----------|------|
| 认证 | `/api/v1/auth` | 登录、登出、Token 刷新、用户信息 |
| 用户 | `/api/v1/system/user` | 用户 CRUD |
| 角色 | `/api/v1/system/role` | 角色 CRUD |
| 菜单 | `/api/v1/system/menu` | 菜单 CRUD |
| 部门 | `/api/v1/system/dept` | 部门 CRUD |
| 系统配置 | `/api/v1/system/config` | 系统配置管理 |
| 仪表盘 | `/api/v1/dashboard` | 统计数据 |
| 服务器监控 | `/api/v1/monitor/server` | 服务器信息 |
| 在线用户 | `/api/v1/monitor/online` | 在线用户管理 |
| 操作日志 | `/api/v1/monitor/operation-log` | 操作日志查询 |
| 登录日志 | `/api/v1/monitor/login-log` | 登录日志查询 |
| 代码生成 | `/api/v1/codegen` | 代码生成工具 |
| 定时任务 | `/api/v1/tool/job` | 定时任务管理 |

## 📊 数据库设计

系统包含以下核心数据表：

| 表名 | 说明 |
|------|------|
| `user` | 用户表 |
| `role` | 角色表 |
| `menu` | 菜单表（目录 / 菜单 / 按钮） |
| `dept` | 部门表 |
| `user_role` | 用户-角色关联表 |
| `role_menu` | 角色-菜单关联表 |
| `role_dept` | 角色-部门关联表（数据权限） |
| `operation_log` | 操作日志表 |
| `login_log` | 登录日志表 |
| `codegen_config` | 代码生成配置表 |
| `job` | 定时任务表 |
| `job_log` | 定时任务执行日志表 |
| `system_config` | 系统配置表 |

## 📝 设计文档

项目设计文档位于 [`plans/`](plans/) 目录：

- [项目总览](plans/00-项目总览.md)
- [后端架构设计](plans/01-后端架构设计.md)
- [前端架构设计](plans/02-前端架构设计.md)
- [数据库设计](plans/03-数据库设计.md)
- [RBAC 权限设计](plans/04-RBAC权限设计.md)
- [API 接口规范](plans/05-API接口规范.md)
- [代码生成设计](plans/06-代码生成设计.md)

## 📜 开源协议

MIT License
