/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80018
 Source Host           : localhost:3306
 Source Schema         : user_srv

 Target Server Type    : MySQL
 Target Server Version : 80018
 File Encoding         : 65001

 Date: 16/03/2022 21:10:59
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `username` varchar(50) NOT NULL COMMENT '账号',
  `password` varchar(200) NOT NULL COMMENT '账号密码',
  `address` varchar(100) DEFAULT NULL COMMENT '地址',
  `avatar` varchar(100) DEFAULT NULL COMMENT '头像',
  `desc` varchar(100) DEFAULT NULL COMMENT '描述',
  `gender` varchar(10) DEFAULT NULL COMMENT '性别',
  `birth_day` datetime DEFAULT NULL COMMENT '出生年月',
  `role_id` int(11) DEFAULT NULL COMMENT '角色ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of account
-- ----------------------------
BEGIN;
INSERT INTO `account` VALUES (1, '2021-11-21 10:35:14.395', '2021-11-21 10:35:14.395', NULL, 'admin', '$2a$10$5zQ6meB74pXW3Ak.fD0Wl.KV9hwe1S4q5OiqYIircm1vxr8iHnZ1u', '', '', '', '', NULL, 0);
COMMIT;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
  `mobile` varchar(11) NOT NULL COMMENT '手机号码',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nick_name` varchar(50) DEFAULT NULL COMMENT '昵称',
  `birthday` datetime DEFAULT NULL COMMENT '生日',
  `gender` varchar(6) DEFAULT 'male' COMMENT 'female表示女,male表示男',
  `role` int(11) DEFAULT '1' COMMENT '1表示普通用户,2表示管理员',
  PRIMARY KEY (`id`),
  UNIQUE KEY `mobile` (`mobile`),
  KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
BEGIN;
INSERT INTO `user` VALUES (1, '2021-10-31 17:22:26.000', '2021-10-31 17:22:31.000', NULL, '18170601666', '122', NULL, '2021-10-30 21:49:48', 'male', 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
