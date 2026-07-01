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
(1, 0, 'Head Office', 1, 1),
(2, 1, 'R&D Department', 1, 1),
(3, 1, 'Operations Department', 2, 1),
(4, 1, 'Teaching Lab Center', 3, 1);

INSERT INTO sys_role(id, name, code, remark, status) VALUES
(1, 'Administrator', 'admin', 'Built-in system administrator role', 1),
(2, 'Operator', 'ops', 'Daily device and file management role', 1);

INSERT INTO sys_user(id, username, password, nickname, email, phone, status, dept_id) VALUES
(1, 'admin', '$2a$10$mT3YZaXbZIbJfxBwkjDlYO32DM.cSG4uR3Z06k0PLdbxH1zvspZB.', 'System Admin', 'admin@example.com', '13800000000', 1, 2);

INSERT INTO sys_user_role(user_id, role_id) VALUES (1, 1);

INSERT INTO sys_menu(id, parent_id, name, path, component, permission, icon, type, sort) VALUES
(1, 0, 'Dashboard', '/dashboard', 'Dashboard', 'dashboard:view', 'House', 0, 1),
(2, 0, 'System Management', '/system', 'Layout', '', 'Setting', 0, 10),
(3, 2, 'User Management', '/system/users', 'UserManage', 'sys:user:list', 'User', 0, 11),
(4, 3, 'Create User', '', '', 'sys:user:create', '', 1, 1),
(5, 3, 'Update User', '', '', 'sys:user:update', '', 1, 2),
(6, 3, 'Delete User', '', '', 'sys:user:delete', '', 1, 3),
(7, 3, 'Reset Password', '', '', 'sys:user:password', '', 1, 4),
(8, 3, 'Assign Roles', '', '', 'sys:user:roles', '', 1, 5),
(9, 2, 'Role Management', '/system/roles', 'RoleManage', 'sys:role:list', 'Avatar', 0, 12),
(10, 9, 'Create Role', '', '', 'sys:role:create', '', 1, 1),
(11, 9, 'Update Role', '', '', 'sys:role:update', '', 1, 2),
(12, 9, 'Delete Role', '', '', 'sys:role:delete', '', 1, 3),
(13, 9, 'Assign Menus', '', '', 'sys:role:menus', '', 1, 4),
(14, 2, 'Menu Management', '/system/menus', 'MenuManage', 'sys:menu:list', 'Menu', 0, 13),
(15, 14, 'Create Menu', '', '', 'sys:menu:create', '', 1, 1),
(16, 14, 'Update Menu', '', '', 'sys:menu:update', '', 1, 2),
(17, 14, 'Delete Menu', '', '', 'sys:menu:delete', '', 1, 3),
(18, 2, 'Department Management', '/system/depts', 'DeptManage', 'sys:dept:list', 'OfficeBuilding', 0, 14),
(19, 18, 'Create Department', '', '', 'sys:dept:create', '', 1, 1),
(20, 18, 'Update Department', '', '', 'sys:dept:update', '', 1, 2),
(21, 18, 'Delete Department', '', '', 'sys:dept:delete', '', 1, 3),
(22, 0, 'Device Management', '/devices', 'DeviceManage', 'device:list', 'Monitor', 0, 20),
(23, 22, 'Create Device', '', '', 'device:create', '', 1, 1),
(24, 22, 'Update Device', '', '', 'device:update', '', 1, 2),
(25, 22, 'Delete Device', '', '', 'device:delete', '', 1, 3),
(26, 22, 'Device Statistics', '', '', 'device:statistics', '', 1, 4),
(27, 22, 'Device Type List', '', '', 'device:type:list', '', 1, 5),
(28, 22, 'Create Device Type', '', '', 'device:type:create', '', 1, 6),
(29, 22, 'Update Device Type', '', '', 'device:type:update', '', 1, 7),
(30, 22, 'Delete Device Type', '', '', 'device:type:delete', '', 1, 8),
(31, 0, 'File Management', '/files', 'FileManage', 'file:list', 'FolderOpened', 0, 30),
(32, 31, 'File List', '', '', 'file:list', '', 1, 1),
(33, 31, 'Upload File', '', '', 'file:upload', '', 1, 2),
(34, 31, 'Download File', '', '', 'file:download', '', 1, 3),
(35, 31, 'Delete File', '', '', 'file:delete', '', 1, 4);

INSERT INTO sys_role_menu(role_id, menu_id) SELECT 1, id FROM sys_menu;
INSERT INTO sys_role_menu(role_id, menu_id)
SELECT 2, id FROM sys_menu
WHERE permission IN ('dashboard:view','device:list','device:create','device:update','device:statistics','device:type:list','file:list','file:upload','file:download');

INSERT INTO device_type(id, name, code, remark) VALUES
(1, 'Server', 'server', 'Server and virtual machine assets'),
(2, 'Switch', 'switch', 'Network switching devices'),
(3, 'Camera', 'camera', 'Security camera devices'),
(4, 'Sensor', 'sensor', 'Environment monitoring sensors');

INSERT INTO device(name, code, type_id, status, location, remark) VALUES
('Auth Server', 'DEV-0001', 1, 'online', 'Room A-101', 'Authentication node'),
('Database Server', 'DEV-0002', 1, 'online', 'Room A-101', 'MySQL primary database'),
('Floor Switch', 'DEV-0003', 2, 'offline', 'Building B Floor 2', 'Waiting for inspection'),
('Gate Camera', 'DEV-0004', 3, 'repair', 'Main Gate', 'Lens maintenance'),
('Temperature Sensor', 'DEV-0005', 4, 'online', 'Lab 301', 'Environment monitoring');

SET FOREIGN_KEY_CHECKS = 1;
