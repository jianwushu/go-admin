-- go-admin 初始数据脚本
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
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (1, 0, '系统管理', '/system', '', 'setting', 0, 1, 1, 1, '');

-- 用户管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (2, 1, '用户管理', 'user', 'system/user/index', 'user', 1, 1, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (3, 2, '用户查询', '', '', '', 2, 0, 1, 1, 'system:user:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (4, 2, '用户新增', '', '', '', 2, 0, 1, 1, 'system:user:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (5, 2, '用户修改', '', '', '', 2, 0, 1, 1, 'system:user:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (6, 2, '用户删除', '', '', '', 2, 0, 1, 1, 'system:user:delete');

-- 角色管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (7, 1, '角色管理', 'role', 'system/role/index', 'peoples', 1, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (8, 7, '角色查询', '', '', '', 2, 0, 1, 1, 'system:role:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (9, 7, '角色新增', '', '', '', 2, 0, 1, 1, 'system:role:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (10, 7, '角色修改', '', '', '', 2, 0, 1, 1, 'system:role:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (11, 7, '角色删除', '', '', '', 2, 0, 1, 1, 'system:role:delete');

-- 菜单管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (12, 1, '菜单管理', 'menu', 'system/menu/index', 'tree-table', 1, 3, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (13, 12, '菜单查询', '', '', '', 2, 0, 1, 1, 'system:menu:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (14, 12, '菜单新增', '', '', '', 2, 0, 1, 1, 'system:menu:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (15, 12, '菜单修改', '', '', '', 2, 0, 1, 1, 'system:menu:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (16, 12, '菜单删除', '', '', '', 2, 0, 1, 1, 'system:menu:delete');

-- 部门管理菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (17, 1, '部门管理', 'dept', 'system/dept/index', 'tree', 1, 4, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (18, 17, '部门查询', '', '', '', 2, 0, 1, 1, 'system:dept:list');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (19, 17, '部门新增', '', '', '', 2, 0, 1, 1, 'system:dept:add');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (20, 17, '部门修改', '', '', '', 2, 0, 1, 1, 'system:dept:edit');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (21, 17, '部门删除', '', '', '', 2, 0, 1, 1, 'system:dept:delete');

-- 系统监控目录
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (22, 0, '系统监控', '/monitor', '', 'monitor', 0, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (23, 22, '服务器监控', 'server', 'monitor/server/index', 'server', 1, 1, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (24, 22, '操作日志', 'operation-log', 'monitor/operation-log/index', 'form', 1, 2, 1, 1, '');
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (25, 22, '登录日志', 'login-log', 'monitor/login-log/index', 'logininfor', 1, 3, 1, 1, '');

-- 系统工具目录
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (26, 0, '系统工具', '/tool', '', 'tool', 0, 3, 1, 1, '');

-- 代码生成菜单
INSERT OR IGNORE INTO {{.TablePrefix}}menu (id, parent_id, name, path, component, icon, type, sort, visible, status, perms) VALUES (27, 26, '代码生成', 'codegen', 'tool/codegen/index', 'code', 1, 1, 1, 1, '');

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

-- 超级管理员拥有所有部门权限
INSERT OR IGNORE INTO {{.TablePrefix}}role_dept (role_id, dept_id) VALUES (1, 1);
INSERT OR IGNORE INTO {{.TablePrefix}}role_dept (role_id, dept_id) VALUES (1, 2);
