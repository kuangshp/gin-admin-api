-- ----------------------------
-- 后台账号表
-- ----------------------------
DROP TABLE IF EXISTS `sys_account`;
CREATE TABLE `sys_account` (
   `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
   `username` varchar(60) NOT NULL COMMENT '登录帐号',
   `email` varchar(60) DEFAULT null COMMENT '邮箱',
   `mobile` varchar(11) DEFAULT null COMMENT '手机号',
   `password` varchar(100) NOT NULL COMMENT '登录密码',
   `last_login_date` datetime default null COMMENT '最后一次登录时间',
   `last_login_ip` varchar(30) default null  COMMENT '最后一次登录ip',
   `status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '状态1是正常,2是禁用',
   `avatar` varchar(200) default null comment '头像',
   `is_admin` tinyint(4) NOT NULL DEFAULT 2 COMMENT '1是超级管理员，2是普通管理员',
   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
   `created_by` int(11) DEFAULT NULL COMMENT '创建人',
   `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
   UNIQUE KEY `uk_username` (`username`) USING BTREE,
   UNIQUE KEY `uk_email` (`email`) USING BTREE,
   UNIQUE KEY `uk_mobile` (`mobile`) USING BTREE,
   KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='后台账号表';


-- ----------------------------
-- 角色表
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
    `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
    `name` varchar(50) NOT NULL COMMENT '角色名称',
    `description` varchar(255) DEFAULT NULL COMMENT '描述',
    `status` tinyint(4) DEFAULT 1 COMMENT '状态1是正常,2是禁用',
    `sort` int(11) DEFAULT 1 COMMENT '排序',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
    `created_by` int(11) DEFAULT NULL COMMENT '创建人',
    `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
    UNIQUE KEY `uk_name` (`name`) USING BTREE,
    KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- ----------------------------
-- 账号角色表
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_role`;
CREATE TABLE `sys_account_role` (
    `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
    `account_id` int(11) NOT NULL COMMENT '关联到sys_account表主键id',
    `role_id` int(11) NOT NULL COMMENT '关联到sys_role表主键id',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
    `created_by` int(11) DEFAULT NULL COMMENT '创建人',
    `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
    UNIQUE KEY `uk_account_role` (`account_id`, `role_id`) USING BTREE,
    KEY `idx_account_id` (`account_id`) USING BTREE,
    KEY `idx_role_id` (`role_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账号和角色中间表';

-- ----------------------------
-- 资源表
-- ----------------------------
DROP TABLE IF EXISTS `sys_resources`;
CREATE TABLE `sys_resources` (
     `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
     `title` varchar(50) NOT NULL COMMENT '名称:按钮标题,或菜单标题',
     `url` varchar(100) NOT NULL COMMENT '按钮请求url,或菜单路由',
     `method` varchar(50) DEFAULT NULL COMMENT '接口的请求方式',
     `icon` varchar(100) DEFAULT NULL COMMENT '菜单小图标',
     `resources_type` tinyint(4) DEFAULT 1 COMMENT '类型:1表示目录,2表示菜单,3表示接口',
     `is_cache` tinyint(4) DEFAULT 1 COMMENT '是否缓存:1表示缓存:2不缓存',
     `is_hidden` tinyint(4) DEFAULT 1 COMMENT '是否隐藏:1表示不隐藏,2表示隐藏',
     `is_link` tinyint(4) DEFAULT 1 COMMENT '是否为外部链接:1表示不是,2表示是',
     `parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '上一级id，0=顶级',
     `sort` int(11) DEFAULT 1 COMMENT '菜单,或按钮排序',
     `status` tinyint(4) DEFAULT 1 COMMENT '状态1是正常,2是禁用',
     `description` varchar(200) DEFAULT NULL COMMENT '描述',
     `is_admin_have` tinyint(4) default 0 comment '是否超管独有,1表示是,0表示不是',
     `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
     `created_by` int(11) DEFAULT NULL COMMENT '创建人',
     `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
     KEY `idx_parent_id` (`parent_id`) USING BTREE,
     KEY `idx_resources_type` (`resources_type`) USING BTREE,
     KEY `idx_status` (`status`) USING BTREE
)  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资源表';

-- ----------------------------
-- 角色和资源中间表
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_resources`;
CREATE TABLE `sys_role_resources` (
  `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
  `resources_id` int(11) NOT NULL COMMENT '关联到sys_resources表主键id',
  `role_id` int(11) NOT NULL COMMENT '关联到sys_role表主键id',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  `created_by` int(11) DEFAULT NULL COMMENT '创建人',
  `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
  UNIQUE KEY `uk_role_resources` (`resources_id`, `role_id`) USING BTREE,
  KEY `idx_resources_id` (`resources_id`) USING BTREE,
  KEY `idx_role_id` (`role_id`) USING BTREE
)   ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色和资源中间表';

-- ----------------------------
-- 部门表（最终正确版）
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
    `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
    `name` varchar(50) NOT NULL COMMENT '部门名称',
    `parent_id` int(11) NOT NULL DEFAULT 0 COMMENT '上级部门id，0=顶级',
    `full_id` varchar(255) NOT NULL DEFAULT '' COMMENT '全层级ID 例：1,5,12',
    `full_name` varchar(255) NOT NULL DEFAULT '' COMMENT '全层级名称 例：总公司,技术部,后端组',
    `sort` int(11) DEFAULT 1 COMMENT '排序',
    `status` tinyint(4) DEFAULT 1 COMMENT '状态1是正常,2是禁用',
    `leader` varchar(50) DEFAULT NULL COMMENT '部门负责人',
    `phone` varchar(20) DEFAULT NULL COMMENT '部门联系电话',
    `email` varchar(60) DEFAULT NULL COMMENT '部门邮箱',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
    `created_by` int(11) DEFAULT NULL COMMENT '创建人',
    `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
    KEY `idx_parent_id` (`parent_id`) USING BTREE,
    KEY `idx_full_id` (`full_id`) USING BTREE,
    KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='部门表';

-- ----------------------------
-- 岗位表
-- ----------------------------
DROP TABLE IF EXISTS `sys_post`;
CREATE TABLE `sys_post` (
    `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
    `name` varchar(50) NOT NULL COMMENT '岗位名称',
    `code` varchar(60) NOT NULL COMMENT '岗位编码',
    `dept_id` int(11) NOT NULL COMMENT '所属部门id',
    `sort` int(11) DEFAULT 1 COMMENT '排序',
    `status` tinyint(4) DEFAULT 1 COMMENT '状态1是正常,2是禁用',
    `remark` varchar(255) DEFAULT NULL COMMENT '备注',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
    `created_by` int(11) DEFAULT NULL COMMENT '创建人',
    `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
    UNIQUE KEY `uk_code` (`code`) USING BTREE,
    KEY `idx_dept_id` (`dept_id`) USING BTREE,
    KEY `idx_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='岗位表';

-- ----------------------------
-- 账号与岗位关联表
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_post`;
CREATE TABLE `sys_account_post` (
    `id` int(11) NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
    `account_id` int(11) NOT NULL COMMENT '账号id',
    `post_id` int(11) NOT NULL COMMENT '岗位id',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
    `created_by` int(11) DEFAULT NULL COMMENT '创建人',
    `updated_by` int(11) DEFAULT NULL COMMENT '更新人',
    UNIQUE KEY `uk_account_post` (`account_id`, `post_id`) USING BTREE,
    KEY `idx_account_id` (`account_id`) USING BTREE,
    KEY `idx_post_id` (`post_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='账号岗位关联表';

-- ----------------------------
-- 角色数据权限表
-- ----------------------------
DROP TABLE IF EXISTS `sys_data_scope`;
CREATE TABLE `sys_data_scope` (
  `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键id',
  `role_id` int(11) NOT NULL COMMENT '角色ID',
  `data_type` tinyint(4) NOT NULL DEFAULT 4 COMMENT '数据范围：1全部 2本部门 3本部门及下级 4仅本人 5自定义部门',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY `uk_role_id` (`role_id`) USING BTREE,
  KEY `idx_data_type` (`data_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色数据权限表';

-- ----------------------------
-- 数据权限-自定义部门表
-- ----------------------------
DROP TABLE IF EXISTS `sys_data_scope_dept`;
CREATE TABLE `sys_data_scope_dept` (
   `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键id',
   `scope_id` int(11) NOT NULL COMMENT '数据权限id',
   `dept_id` int(11) NOT NULL COMMENT '可查看部门id',
   `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` datetime NULL DEFAULT NULL COMMENT '软删除时间',
   UNIQUE KEY `uk_scope_dept` (`scope_id`,`dept_id`) USING BTREE,
   KEY `idx_scope_id` (`scope_id`) USING BTREE,
   KEY `idx_dept_id` (`dept_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色数据权限-自定义部门关联';

