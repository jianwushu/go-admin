-- go-admin 建表脚本 (MySQL)
-- 使用 Go text/template 渲染，变量 {{.TablePrefix}} 来自 config.yaml

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 用户表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) NOT NULL,
    `password` VARCHAR(128) NOT NULL,
    `nickname` VARCHAR(64) DEFAULT NULL,
    `email` VARCHAR(128) DEFAULT NULL,
    `phone` VARCHAR(20) DEFAULT NULL,
    `status` INT DEFAULT 1,
    `dept_id` BIGINT DEFAULT NULL,
    `avatar` VARCHAR(256) DEFAULT NULL,
    `remark` VARCHAR(512) DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_username` (`username`),
    KEY `idx_user_dept_id` (`dept_id`),
    KEY `idx_user_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}role` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(64) NOT NULL,
    `code` VARCHAR(64) NOT NULL,
    `data_scope` INT DEFAULT 1,
    `sort` INT DEFAULT 0,
    `status` INT DEFAULT 1,
    `remark` VARCHAR(512) DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_role_code` (`code`),
    KEY `idx_role_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 菜单表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}menu` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `parent_id` BIGINT DEFAULT 0,
    `name` VARCHAR(64) NOT NULL,
    `i18n_key` VARCHAR(128) DEFAULT '',
    `path` VARCHAR(128) DEFAULT NULL,
    `component` VARCHAR(128) DEFAULT NULL,
    `icon` VARCHAR(64) DEFAULT NULL,
    `type` INT DEFAULT NULL,
    `sort` INT DEFAULT 0,
    `visible` INT DEFAULT 1,
    `status` INT DEFAULT 1,
    `perms` VARCHAR(128) DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_menu_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 部门表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}dept` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `parent_id` BIGINT DEFAULT 0,
    `name` VARCHAR(64) NOT NULL,
    `sort` INT DEFAULT 0,
    `status` INT DEFAULT 1,
    `leader` VARCHAR(64) DEFAULT NULL,
    `phone` VARCHAR(20) DEFAULT NULL,
    `email` VARCHAR(128) DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_dept_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}user_role` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT NOT NULL,
    `role_id` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_role` (`user_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色菜单关联表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}role_menu` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `role_id` BIGINT NOT NULL,
    `menu_id` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_role_menu` (`role_id`, `menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 角色部门关联表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}role_dept` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `role_id` BIGINT NOT NULL,
    `dept_id` BIGINT NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_role_dept` (`role_id`, `dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 操作日志表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}operation_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `module` VARCHAR(64) DEFAULT NULL,
    `action` VARCHAR(64) DEFAULT NULL,
    `method` VARCHAR(10) DEFAULT NULL,
    `url` VARCHAR(256) DEFAULT NULL,
    `ip` VARCHAR(64) DEFAULT NULL,
    `operator` VARCHAR(64) DEFAULT NULL,
    `request_param` TEXT DEFAULT NULL,
    `response_data` TEXT DEFAULT NULL,
    `status` INT DEFAULT NULL,
    `error_msg` VARCHAR(512) DEFAULT NULL,
    `duration` INT DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_oplog_created_at` (`created_at`),
    KEY `idx_oplog_operator` (`operator`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 登录日志表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}login_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(64) DEFAULT NULL,
    `ip` VARCHAR(64) DEFAULT NULL,
    `location` VARCHAR(128) DEFAULT NULL,
    `browser` VARCHAR(64) DEFAULT NULL,
    `os` VARCHAR(64) DEFAULT NULL,
    `status` INT DEFAULT NULL,
    `msg` VARCHAR(256) DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_loginlog_created_at` (`created_at`),
    KEY `idx_loginlog_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 代码生成配置表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}codegen_config` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `table_name` VARCHAR(128) NOT NULL,
    `table_comment` VARCHAR(256) DEFAULT NULL,
    `class_name` VARCHAR(128) NOT NULL,
    `business_name` VARCHAR(128) NOT NULL,
    `function_name` VARCHAR(256) DEFAULT NULL,
    `module_name` VARCHAR(64) NOT NULL,
    `package_name` VARCHAR(128) NOT NULL,
    `author` VARCHAR(64) DEFAULT NULL,
    `fields` TEXT DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_codegen_config_table_name` (`table_name`),
    KEY `idx_codegen_config_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ========== 第二阶段新增表 ==========

-- 定时任务表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}job` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(128) NOT NULL,
    `job_type` INT DEFAULT 1,
    `cron_expr` VARCHAR(64) NOT NULL,
    `func_name` VARCHAR(128) DEFAULT NULL,
    `http_url` VARCHAR(512) DEFAULT NULL,
    `http_method` VARCHAR(10) DEFAULT 'GET',
    `http_body` TEXT DEFAULT NULL,
    `status` INT DEFAULT 1,
    `remark` VARCHAR(512) DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_job_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 定时任务执行日志表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}job_log` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `job_id` BIGINT NOT NULL,
    `job_name` VARCHAR(128) DEFAULT NULL,
    `status` INT DEFAULT 1,
    `result` TEXT DEFAULT NULL,
    `error_msg` TEXT DEFAULT NULL,
    `duration` INT DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_job_log_job_id` (`job_id`),
    KEY `idx_job_log_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统配置表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}system_config` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `config_key` VARCHAR(128) NOT NULL,
    `config_value` TEXT DEFAULT NULL,
    `config_type` VARCHAR(32) DEFAULT 'text',
    `remark` VARCHAR(512) DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_system_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统公告表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}announcement` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `title` VARCHAR(256) NOT NULL,
    `content` TEXT DEFAULT NULL,
    `type` INT DEFAULT 1,
    `status` INT DEFAULT 1,
    `top` INT DEFAULT 0,
    `publish_by` BIGINT DEFAULT NULL,
    `publish_at` DATETIME DEFAULT NULL,
    `created_by` BIGINT DEFAULT 0,
    `updated_by` BIGINT DEFAULT 0,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_announcement_deleted_at` (`deleted_at`),
    KEY `idx_announcement_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 站内信消息表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}message` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `sender_id` BIGINT NOT NULL,
    `title` VARCHAR(256) DEFAULT NULL,
    `content` TEXT DEFAULT NULL,
    `msg_type` INT DEFAULT 1,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_message_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 站内信用户关联表
CREATE TABLE IF NOT EXISTS `{{.TablePrefix}}message_user` (
    `id` BIGINT NOT NULL AUTO_INCREMENT,
    `message_id` BIGINT NOT NULL,
    `receiver_id` BIGINT NOT NULL,
    `is_read` INT DEFAULT 0,
    `read_at` DATETIME DEFAULT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_message_user_receiver` (`receiver_id`),
    KEY `idx_message_user_message` (`message_id`),
    KEY `idx_message_user_read` (`receiver_id`, `is_read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
