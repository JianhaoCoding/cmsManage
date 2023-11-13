/*
 Navicat MySQL Data Transfer

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 80100
 Source Host           : localhost:3306
 Source Schema         : cms

 Target Server Type    : MySQL
 Target Server Version : 80100
 File Encoding         : 65001

 Date: 13/11/2023 10:11:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for adminer
-- ----------------------------
DROP TABLE IF EXISTS `adminer`;
CREATE TABLE `adminer` (
  `adminer_id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '管理员名称',
  `password` varchar(64) NOT NULL DEFAULT '' COMMENT '密码',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '昵称，用于前台显示',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '管理员邮箱',
  `mobile` char(11) NOT NULL DEFAULT '' COMMENT '管理员手机号',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '管理员组ID',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态 默认1 禁用 2启用 99删除',
  `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `last_ip` char(15) NOT NULL DEFAULT '' COMMENT '最后登陆IP',
  `last_time` int unsigned NOT NULL DEFAULT '0' COMMENT '最后登陆时间',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `update_time` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`adminer_id`),
  UNIQUE KEY `idx_adminer_name` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='后台管理员';

-- ----------------------------
-- Table structure for adminer_auth
-- ----------------------------
DROP TABLE IF EXISTS `adminer_auth`;
CREATE TABLE `adminer_auth` (
  `auth_id` int unsigned NOT NULL AUTO_INCREMENT,
  `auth_name` varchar(20) NOT NULL DEFAULT '' COMMENT '权限名称',
  `auth_con_act` varchar(100) NOT NULL DEFAULT '' COMMENT '权限对应的控制器和方法',
  `auth_roles` varchar(50) NOT NULL DEFAULT '' COMMENT '权限对应的角色',
  `auth_group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '权限组ID',
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '状态 默认1禁用  2启用 99删除',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '权限备注',
  `prent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `is_menu_show` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '是否导航展示1否，2是',
  `sort` int unsigned NOT NULL DEFAULT '0' COMMENT '排序字段',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `create_adminer` int unsigned NOT NULL DEFAULT '0' COMMENT '添加用户',
  `last_time` int unsigned NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  `last_adminer` int unsigned NOT NULL DEFAULT '0' COMMENT '最后更新用户',
  PRIMARY KEY (`auth_id`),
  KEY `idx_aca` (`auth_con_act`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='权限列表';

-- ----------------------------
-- Table structure for adminer_group
-- ----------------------------
DROP TABLE IF EXISTS `adminer_group`;
CREATE TABLE `adminer_group` (
  `group_id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT '组名',
  `role` varchar(1000) NOT NULL DEFAULT '' COMMENT '用户组权限',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '组备注',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `create_adminer` int unsigned NOT NULL DEFAULT '0' COMMENT '添加管理员',
  PRIMARY KEY (`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员组';

-- ----------------------------
-- Table structure for adminer_group_auth
-- ----------------------------
DROP TABLE IF EXISTS `adminer_group_auth`;
CREATE TABLE `adminer_group_auth` (
  `adminer_group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '组ID',
  `auth_id` int unsigned NOT NULL DEFAULT '0' COMMENT '权限ID',
  `create_time` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `create_adminer` int unsigned NOT NULL DEFAULT '0' COMMENT '添加管理员',
  UNIQUE KEY `idx_agid_auid` (`adminer_group_id`,`auth_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户组权限关联表';

-- ----------------------------
-- Table structure for adminer_log
-- ----------------------------
DROP TABLE IF EXISTS `adminer_log`;
CREATE TABLE `adminer_log` (
  `log_id` int unsigned NOT NULL AUTO_INCREMENT,
  `obj_id` int unsigned NOT NULL DEFAULT '0' COMMENT '操作对象ID',
  `obj_type` varchar(50) NOT NULL DEFAULT '' COMMENT '操作对象类型',
  `operate_type` varchar(50) NOT NULL,
  `sql` text NOT NULL COMMENT '操作SQL',
  `operate_time` int unsigned NOT NULL DEFAULT '0' COMMENT '操作时间',
  `operate_ip` char(15) NOT NULL DEFAULT '' COMMENT '操作IP',
  `operate_adminer` int unsigned NOT NULL DEFAULT '0' COMMENT '操作用户ID',
  PRIMARY KEY (`log_id`),
  KEY `idx_oid_ot_opy` (`obj_id`,`obj_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='管理员日志表';

SET FOREIGN_KEY_CHECKS = 1;
