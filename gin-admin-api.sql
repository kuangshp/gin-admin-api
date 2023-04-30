-- ----------------------------
-- 账号表
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
   `id` int NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
   `username` varchar(50) UNIQUE NOT NULL COMMENT '用户名',
   `password` varchar(100) NOT null COMMENT '密码',
   `status` tinyint(4) DEFAULT 1 COMMENT '状态1是正常,0是禁用',
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` timestamp(6) NULL DEFAULT NULL COMMENT '软删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='后台管理用户';

-- ----------------------------
-- 账号token表(使用redis存储token的时候需要使用)
-- ----------------------------
DROP TABLE IF EXISTS `account_token`;
CREATE TABLE `account_token` (
     `id` int NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '主键id',
     `account_id` int NOT NULL COMMENT '关联到账号表id',
     `username` varchar(50) DEFAULT NULL COMMENT '用户名',
     `token` varchar(200) NOT NULL COMMENT 'token',
     `expire_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'token过期时间',
     `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     `deleted_at` timestamp(6) NULL DEFAULT NULL COMMENT '软删除时间',
     UNIQUE KEY `account_id_token` (`account_id`,`token`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='账号token表';
