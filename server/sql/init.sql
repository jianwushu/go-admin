-- go-admin 建表脚本
-- 使用 Go text/template 渲染，变量 {{.TablePrefix}} 来自 config.yaml

-- 用户表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(128) NOT NULL,
    nickname VARCHAR(64),
    email VARCHAR(128),
    phone VARCHAR(20),
    status INTEGER DEFAULT 1,
    dept_id INTEGER,
    avatar VARCHAR(256),
    remark VARCHAR(512),
    created_by INTEGER DEFAULT 0,
    updated_by INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_username ON {{.TablePrefix}}user(username);
CREATE INDEX IF NOT EXISTS idx_user_dept_id ON {{.TablePrefix}}user(dept_id);
CREATE INDEX IF NOT EXISTS idx_user_deleted_at ON {{.TablePrefix}}user(deleted_at);

-- 角色表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(64) NOT NULL,
    code VARCHAR(64) NOT NULL,
    data_scope INTEGER DEFAULT 1,
    sort INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    remark VARCHAR(512),
    created_by INTEGER DEFAULT 0,
    updated_by INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_role_code ON {{.TablePrefix}}role(code);
CREATE INDEX IF NOT EXISTS idx_role_deleted_at ON {{.TablePrefix}}role(deleted_at);

-- 菜单表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}menu (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER DEFAULT 0,
    name VARCHAR(64) NOT NULL,
    i18n_key VARCHAR(128) DEFAULT '',
    path VARCHAR(128),
    component VARCHAR(128),
    icon VARCHAR(64),
    type INTEGER,
    sort INTEGER DEFAULT 0,
    visible INTEGER DEFAULT 1,
    status INTEGER DEFAULT 1,
    perms VARCHAR(128),
    created_by INTEGER DEFAULT 0,
    updated_by INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_menu_deleted_at ON {{.TablePrefix}}menu(deleted_at);

-- 部门表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}dept (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    parent_id INTEGER DEFAULT 0,
    name VARCHAR(64) NOT NULL,
    sort INTEGER DEFAULT 0,
    status INTEGER DEFAULT 1,
    leader VARCHAR(64),
    phone VARCHAR(20),
    email VARCHAR(128),
    created_by INTEGER DEFAULT 0,
    updated_by INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_dept_deleted_at ON {{.TablePrefix}}dept(deleted_at);

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}user_role (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_role ON {{.TablePrefix}}user_role(user_id, role_id);

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}role_menu (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    menu_id INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_role_menu ON {{.TablePrefix}}role_menu(role_id, menu_id);

-- 角色部门关联表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}role_dept (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    dept_id INTEGER NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_role_dept ON {{.TablePrefix}}role_dept(role_id, dept_id);

-- 操作日志表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}operation_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    module VARCHAR(64),
    action VARCHAR(64),
    method VARCHAR(10),
    url VARCHAR(256),
    ip VARCHAR(64),
    operator VARCHAR(64),
    request_param TEXT,
    response_data TEXT,
    status INTEGER,
    error_msg VARCHAR(512),
    duration INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_oplog_created_at ON {{.TablePrefix}}operation_log(created_at);
CREATE INDEX IF NOT EXISTS idx_oplog_operator ON {{.TablePrefix}}operation_log(operator);

-- 登录日志表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}login_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(64),
    ip VARCHAR(64),
    location VARCHAR(128),
    browser VARCHAR(64),
    os VARCHAR(64),
    status INTEGER,
    msg VARCHAR(256),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_loginlog_created_at ON {{.TablePrefix}}login_log(created_at);
CREATE INDEX IF NOT EXISTS idx_loginlog_username ON {{.TablePrefix}}login_log(username);

-- 代码生成配置表
CREATE TABLE IF NOT EXISTS {{.TablePrefix}}codegen_config (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    table_name VARCHAR(128) NOT NULL,
    table_comment VARCHAR(256),
    class_name VARCHAR(128) NOT NULL,
    business_name VARCHAR(128) NOT NULL,
    function_name VARCHAR(256),
    module_name VARCHAR(64) NOT NULL,
    package_name VARCHAR(128) NOT NULL,
    author VARCHAR(64),
    fields TEXT,
    created_by INTEGER DEFAULT 0,
    updated_by INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_codegen_config_table_name ON {{.TablePrefix}}codegen_config(table_name);
CREATE INDEX IF NOT EXISTS idx_codegen_config_deleted_at ON {{.TablePrefix}}codegen_config(deleted_at);
