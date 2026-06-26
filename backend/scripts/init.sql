CREATE DATABASE IF NOT EXISTS micro_admin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE micro_admin;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS sys_user_role;
DROP TABLE IF EXISTS sys_role_menu;
DROP TABLE IF EXISTS file_info;
DROP TABLE IF EXISTS device;
DROP TABLE IF EXISTS device_type;
DROP TABLE IF EXISTS sys_user;
DROP TABLE IF EXISTS sys_role;
DROP TABLE IF EXISTS sys_menu;
DROP TABLE IF EXISTS sys_dept;

CREATE TABLE sys_dept (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  parent_id BIGINT NOT NULL DEFAULT 0,
  name VARCHAR(64) NOT NULL,
  sort INT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sys_user (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  password VARCHAR(100) NOT NULL,
  nickname VARCHAR(64) NOT NULL DEFAULT '',
  email VARCHAR(100) NOT NULL DEFAULT '',
  phone VARCHAR(32) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1,
  dept_id BIGINT DEFAULT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,
  INDEX idx_user_dept(dept_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sys_role (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL UNIQUE,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sys_menu (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  parent_id BIGINT NOT NULL DEFAULT 0,
  name VARCHAR(64) NOT NULL,
  path VARCHAR(128) NOT NULL DEFAULT '',
  component VARCHAR(128) NOT NULL DEFAULT '',
  permission VARCHAR(128) NOT NULL DEFAULT '',
  icon VARCHAR(64) NOT NULL DEFAULT '',
  type TINYINT NOT NULL DEFAULT 0 COMMENT '0 menu, 1 button',
  sort INT NOT NULL DEFAULT 0,
  status TINYINT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_menu_parent(parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sys_user_role (
  user_id BIGINT NOT NULL,
  role_id BIGINT NOT NULL,
  PRIMARY KEY(user_id, role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sys_role_menu (
  role_id BIGINT NOT NULL,
  menu_id BIGINT NOT NULL,
  PRIMARY KEY(role_id, menu_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE device_type (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL UNIQUE,
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE device (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL,
  code VARCHAR(100) NOT NULL UNIQUE,
  type_id BIGINT NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'offline',
  location VARCHAR(128) NOT NULL DEFAULT '',
  remark VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_device_type(type_id),
  INDEX idx_device_status(status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE file_info (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  original_name VARCHAR(255) NOT NULL,
  stored_name VARCHAR(255) NOT NULL,
  path VARCHAR(500) NOT NULL,
  size BIGINT NOT NULL DEFAULT 0,
  content_type VARCHAR(128) NOT NULL DEFAULT '',
  uploader_id BIGINT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME DEFAULT NULL,
  INDEX idx_file_uploader(uploader_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO sys_dept(id, parent_id, name, sort, status) VALUES
(1, 0, '总公司', 1, 1),
(2, 1, '研发部', 1, 1),
(3, 1, '运维部', 2, 1),
(4, 1, '教学实验中心', 3, 1);

INSERT INTO sys_role(id, name, code, remark, status) VALUES
(1, '管理员', 'admin', '系统内置管理员角色', 1),
(2, '运维人员', 'ops', '设备和文件日常管理', 1);

INSERT INTO sys_user(id, username, password, nickname, email, phone, status, dept_id) VALUES
(1, 'admin', '$2a$10$mT3YZaXbZIbJfxBwkjDlYO32DM.cSG4uR3Z06k0PLdbxH1zvspZB.', '系统管理员', 'admin@example.com', '13800000000', 1, 2);

INSERT INTO sys_user_role(user_id, role_id) VALUES (1, 1);

INSERT INTO sys_menu(id, parent_id, name, path, component, permission, icon, type, sort) VALUES
(1, 0, '首页', '/dashboard', 'Dashboard', 'dashboard:view', 'House', 0, 1),
(2, 0, '系统管理', '/system', 'Layout', '', 'Setting', 0, 10),
(3, 2, '用户管理', '/system/users', 'UserManage', 'sys:user:list', 'User', 0, 11),
(4, 3, '新增用户', '', '', 'sys:user:create', '', 1, 1),
(5, 3, '修改用户', '', '', 'sys:user:update', '', 1, 2),
(6, 3, '删除用户', '', '', 'sys:user:delete', '', 1, 3),
(7, 3, '重置密码', '', '', 'sys:user:password', '', 1, 4),
(8, 3, '分配角色', '', '', 'sys:user:roles', '', 1, 5),
(9, 2, '角色管理', '/system/roles', 'RoleManage', 'sys:role:list', 'Avatar', 0, 12),
(10, 9, '新增角色', '', '', 'sys:role:create', '', 1, 1),
(11, 9, '修改角色', '', '', 'sys:role:update', '', 1, 2),
(12, 9, '删除角色', '', '', 'sys:role:delete', '', 1, 3),
(13, 9, '分配菜单', '', '', 'sys:role:menus', '', 1, 4),
(14, 2, '菜单管理', '/system/menus', 'MenuManage', 'sys:menu:list', 'Menu', 0, 13),
(15, 14, '新增菜单', '', '', 'sys:menu:create', '', 1, 1),
(16, 14, '修改菜单', '', '', 'sys:menu:update', '', 1, 2),
(17, 14, '删除菜单', '', '', 'sys:menu:delete', '', 1, 3),
(18, 2, '部门管理', '/system/depts', 'DeptManage', 'sys:dept:list', 'OfficeBuilding', 0, 14),
(19, 18, '新增部门', '', '', 'sys:dept:create', '', 1, 1),
(20, 18, '修改部门', '', '', 'sys:dept:update', '', 1, 2),
(21, 18, '删除部门', '', '', 'sys:dept:delete', '', 1, 3),
(22, 0, '设备管理', '/devices', 'DeviceManage', 'device:list', 'Monitor', 0, 20),
(23, 22, '新增设备', '', '', 'device:create', '', 1, 1),
(24, 22, '修改设备', '', '', 'device:update', '', 1, 2),
(25, 22, '删除设备', '', '', 'device:delete', '', 1, 3),
(26, 22, '设备统计', '', '', 'device:statistics', '', 1, 4),
(27, 22, '设备类型查询', '', '', 'device:type:list', '', 1, 5),
(28, 22, '新增设备类型', '', '', 'device:type:create', '', 1, 6),
(29, 22, '修改设备类型', '', '', 'device:type:update', '', 1, 7),
(30, 22, '删除设备类型', '', '', 'device:type:delete', '', 1, 8),
(31, 0, '文件管理', '/files', 'FileManage', 'FolderOpened', 0, 30),
(32, 31, '文件查询', '', '', 'file:list', '', 1, 1),
(33, 31, '文件上传', '', '', 'file:upload', '', 1, 2),
(34, 31, '文件下载', '', '', 'file:download', '', 1, 3),
(35, 31, '文件删除', '', '', 'file:delete', '', 1, 4);

INSERT INTO sys_role_menu(role_id, menu_id) SELECT 1, id FROM sys_menu;
INSERT INTO sys_role_menu(role_id, menu_id)
SELECT 2, id FROM sys_menu WHERE permission IN ('dashboard:view','device:list','device:create','device:update','device:statistics','device:type:list','file:list','file:upload','file:download');

INSERT INTO device_type(id, name, code, remark) VALUES
(1, '服务器', 'server', '机房服务器'),
(2, '交换机', 'switch', '网络交换设备'),
(3, '摄像头', 'camera', '安防摄像设备'),
(4, '传感器', 'sensor', '环境监测传感器');

INSERT INTO device(name, code, type_id, status, location, remark) VALUES
('认证服务器', 'DEV-0001', 1, 'online', 'A 栋机房', '核心认证节点'),
('数据库服务器', 'DEV-0002', 1, 'online', 'A 栋机房', 'MySQL 主库'),
('楼层交换机', 'DEV-0003', 2, 'offline', 'B 栋 2 层', '待巡检'),
('门禁摄像头', 'DEV-0004', 3, 'repair', '校门口', '镜头维修'),
('温湿度传感器', 'DEV-0005', 4, 'online', '实验室 301', '环境监测');

SET FOREIGN_KEY_CHECKS = 1;
