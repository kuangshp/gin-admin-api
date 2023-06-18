-- ----------------------------
-- 账号表
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
   `id` int NOT NULL AUTO_INCREMENT primary key COMMENT '主键id',
   `username` varchar(50) UNIQUE NOT NULL COMMENT '用户名',
   `password` varchar(100) NOT null COMMENT '密码',
   `name` varchar(10) DEFAULT NULL COMMENT '真实姓名',
   `mobile` varchar(11) DEFAULT NULL COMMENT '手机号码',
   `email` varchar(50) DEFAULT NULL COMMENT '邮箱地址',
   `avatar` varchar(200) DEFAULT NULL COMMENT '用户头像',
   `is_admin` tinyint(4) DEFAULT 0 COMMENT '是否为超级管理员:0否,1是',
   `status` tinyint(4) DEFAULT 1 COMMENT '状态1是正常,0是禁用',
   `last_login_ip` varchar(30) COMMENT '最后登录ip地址',
   `last_login_date` timestamp(6) COMMENT '最后登录时间',
   `salt` varchar(10) COMMENT '密码盐',
   `token` varchar(200) DEFAULT NULL COMMENT 'token',
   `expire_time` timestamp DEFAULT NULL COMMENT 'token过期时间',
   `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
   `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间',
   UNIQUE KEY `UK_username_deleted_at` (`username`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='后台管理用户';
