-- go-admin 初始数据脚本 (SQLite)
-- 管理员密码：admin123（bcrypt 加密）

-- 默认部门
INSERT OR IGNORE INTO {{.TablePrefix}}dept (id, parent_id, name, sort, status) VALUES (1, 0, '总公司', 0, 1);
INSERT OR IGNORE INTO {{.TablePrefix}}dept (id, parent_id, name, sort, status) VALUES (2, 1, '默认部门', 0, 1);

-- 默认角色
INSERT OR IGNORE INTO {{.TablePrefix}}role (id, name, code, data_scope, sort, status, remark) VALUES (1, '超级管理员', 'admin', 1, 0, 1, '拥有所有权限');
INSERT OR IGNORE INTO {{.TablePrefix}}role (id, name, code, data_scope, sort, status, remark) VALUES (2, '普通用户', 'user', 4, 1, 1, '仅本人数据');

-- 默认管理员（密码：admin123）
INSERT OR IGNORE INTO {{.TablePrefix}}user (id, username, password, nickname, status, dept_id)
VALUES (1, 'admin', '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', '超级管理员', 1, 1);

-- 管理员角色关联
INSERT OR IGNORE INTO {{.TablePrefix}}user_role (user_id, role_id) VALUES (1, 1);

-- 默认菜单
-- 系统管理目录
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (1, 0, '系统管理', 'menu.system', '/system', '', 'setting', 0, 1, 1, 1, '');

-- 用户管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (2, 1, '用户管理', 'menu.user', 'user', 'system/user/index', 'user', 1, 1, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (3, 2, '用户查询', 'menu.userQuery', '', '', '', 2, 0, 1, 1, 'system:user:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (4, 2, '用户新增', 'menu.userAdd', '', '', '', 2, 0, 1, 1, 'system:user:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (5, 2, '用户修改', 'menu.userEdit', '', '', '', 2, 0, 1, 1, 'system:user:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (6, 2, '用户删除', 'menu.userDelete', '', '', '', 2, 0, 1, 1, 'system:user:delete');

-- 角色管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (7, 1, '角色管理', 'menu.role', 'role', 'system/role/index', 'peoples', 1, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (8, 7, '角色查询', 'menu.roleQuery', '', '', '', 2, 0, 1, 1, 'system:role:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (9, 7, '角色新增', 'menu.roleAdd', '', '', '', 2, 0, 1, 1, 'system:role:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (10, 7, '角色修改', 'menu.roleEdit', '', '', '', 2, 0, 1, 1, 'system:role:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (11, 7, '角色删除', 'menu.roleDelete', '', '', '', 2, 0, 1, 1, 'system:role:delete');

-- 菜单管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (12, 1, '菜单管理', 'menu.menu', 'menu', 'system/menu/index', 'tree-table', 1, 3, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (13, 12, '菜单查询', 'menu.menuQuery', '', '', '', 2, 0, 1, 1, 'system:menu:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (14, 12, '菜单新增', 'menu.menuAdd', '', '', '', 2, 0, 1, 1, 'system:menu:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (15, 12, '菜单修改', 'menu.menuEdit', '', '', '', 2, 0, 1, 1, 'system:menu:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (16, 12, '菜单删除', 'menu.menuDelete', '', '', '', 2, 0, 1, 1, 'system:menu:delete');

-- 部门管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (17, 1, '部门管理', 'menu.dept', 'dept', 'system/dept/index', 'tree', 1, 4, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (18, 17, '部门查询', 'menu.deptQuery', '', '', '', 2, 0, 1, 1, 'system:dept:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (19, 17, '部门新增', 'menu.deptAdd', '', '', '', 2, 0, 1, 1, 'system:dept:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (20, 17, '部门修改', 'menu.deptEdit', '', '', '', 2, 0, 1, 1, 'system:dept:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (21, 17, '部门删除', 'menu.deptDelete', '', '', '', 2, 0, 1, 1, 'system:dept:delete');

-- 系统监控目录
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (22, 0, '系统监控', 'menu.monitor', '/monitor', '', 'monitor', 0, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (23, 22, '服务器监控', 'menu.server', 'server', 'monitor/server/index', 'server', 1, 1, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (24, 22, '操作日志', 'menu.operationLog', 'operation-log', 'monitor/operation-log/index', 'form', 1, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (25, 22, '登录日志', 'menu.loginLog', 'login-log', 'monitor/login-log/index', 'logininfor', 1, 3, 1, 1, '');

-- 在线用户（第二阶段新增）
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (28, 22, '在线用户', 'menu.onlineUser', 'online', 'monitor/online/index', 'online', 1, 4, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (29, 28, '在线用户查询', 'menu.onlineUserQuery', '', '', '', 2, 0, 1, 1, 'monitor:online:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (30, 28, '强制下线', 'menu.onlineUserForceLogout', '', '', '', 2, 0, 1, 1, 'monitor:online:forceLogout');

-- 系统工具目录
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (26, 0, '系统工具', 'menu.tool', '/tool', '', 'tool', 0, 3, 1, 1, '');

-- 代码生成菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (27, 26, '代码生成', 'menu.codegen', 'codegen', 'tool/codegen/index', 'code', 1, 1, 1, 1, '');

-- 定时任务（第二阶段新增）
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (31, 26, '定时任务', 'menu.job', 'job', 'tool/job/index', 'job', 1, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (32, 31, '任务查询', 'menu.jobQuery', '', '', '', 2, 0, 1, 1, 'tool:job:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (33, 31, '任务新增', 'menu.jobAdd', '', '', '', 2, 0, 1, 1, 'tool:job:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (34, 31, '任务修改', 'menu.jobEdit', '', '', '', 2, 0, 1, 1, 'tool:job:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (35, 31, '任务删除', 'menu.jobDelete', '', '', '', 2, 0, 1, 1, 'tool:job:delete');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (36, 31, '任务日志', 'menu.jobLog', 'job-log', 'tool/job/log', 'log', 1, 3, 1, 1, '');

-- 系统配置（第二阶段新增）
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (37, 1, '系统配置', 'menu.systemConfig', 'config', 'system/config/index', 'edit', 1, 5, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (38, 37, '配置查询', 'menu.configQuery', '', '', '', 2, 0, 1, 1, 'system:config:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (39, 37, '配置修改', 'menu.configEdit', '', '', '', 2, 0, 1, 1, 'system:config:edit');

-- 系统公告（第二阶段新增）
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (40, 1, '系统公告', 'menu.announcement', 'announcement', 'system/announcement/index', 'message', 1, 6, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (41, 40, '公告查询', 'menu.announcementQuery', '', '', '', 2, 0, 1, 1, 'system:announcement:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (42, 40, '公告新增', 'menu.announcementAdd', '', '', '', 2, 0, 1, 1, 'system:announcement:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (43, 40, '公告修改', 'menu.announcementEdit', '', '', '', 2, 0, 1, 1, 'system:announcement:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (44, 40, '公告删除', 'menu.announcementDelete', '', '', '', 2, 0, 1, 1, 'system:announcement:delete');

-- 站内信（第二阶段新增）
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (45, 0, '站内信', 'menu.message', '/message', '', 'email', 0, 4, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (46, 45, '我的消息', 'menu.myMessage', 'inbox', 'message/inbox/index', 'message', 1, 1, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, i18n_key, path, component, icon, type, sort, visible, status, perms) VALUES (47, 45, '发送消息', 'menu.sendMessage', 'send', 'message/send/index', 'edit', 1, 2, 1, 1, '');

-- 超级管理员角色拥有所有菜单权限
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 1);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 2);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 3);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 4);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 5);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 6);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 7);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 8);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 9);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 10);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 11);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 12);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 13);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 14);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 15);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 16);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 17);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 18);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 19);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 20);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 21);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 22);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 23);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 24);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 25);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 26);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 27);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 28);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 29);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 30);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 31);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 32);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 33);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 34);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 35);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 36);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 37);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 38);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 39);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 40);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 41);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 42);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 43);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 44);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 45);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 46);
INSERT OR IGNORE INTO {{.TablePrefix}}role_menu (role_id, menu_id) VALUES (1, 47);

-- 超级管理员拥有所有部门权限
INSERT OR IGNORE INTO {{.TablePrefix}}role_dept (role_id, dept_id) VALUES (1, 1);
INSERT OR IGNORE INTO {{.TablePrefix}}role_dept (role_id, dept_id) VALUES (1, 2);

-- 系统配置初始数据
INSERT OR IGNORE INTO {{.TablePrefix}}system_config (config_key, config_value, config_type, remark) VALUES ('system_title', 'Go-Admin', 'text', '系统标题');
INSERT OR IGNORE INTO {{.TablePrefix}}system_config (config_key, config_value, config_type, remark) VALUES ('system_copyright', '© 2026 Go-Admin. All rights reserved.', 'text', '版权信息');
INSERT OR IGNORE INTO {{.TablePrefix}}system_config (config_key, config_value, config_type, remark) VALUES ('system_logo', '', 'image', '系统Logo');
INSERT OR IGNORE INTO {{.TablePrefix}}system_config (config_key, config_value, config_type, remark) VALUES ('system_favicon', '', 'image', 'Favicon');
