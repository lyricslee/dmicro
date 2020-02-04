-- --------------------------------------------------------
-- 主机:                           192.168.56.189
-- 服务器版本:                        5.7.29 - MySQL Community Server (GPL)
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  8.3.0.4694
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 导出  表 dm_passport.published 结构
CREATE TABLE IF NOT EXISTS `published` (
  `id` bigint(20) unsigned NOT NULL COMMENT '消息ID',
  `topic` varchar(256) NOT NULL DEFAULT '' COMMENT '消息主题名称',
  `name` varchar(256) NOT NULL DEFAULT '' COMMENT 'pb消息名称',
  `version` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息版本',
  `msg` varbinary(8192) NOT NULL DEFAULT '' COMMENT '消息以pb格式存储',
  `retries` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '重试次数',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态[0:未投递;1:投递成功;2:投递失败]',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- 数据导出被取消选择。


-- 导出  表 dm_passport.received 结构
CREATE TABLE IF NOT EXISTS `received` (
  `id` bigint(20) unsigned NOT NULL COMMENT '消息ID',
  `topic` varchar(256) NOT NULL DEFAULT '' COMMENT '消息主题名称',
  `name` varchar(256) NOT NULL DEFAULT '' COMMENT 'pb消息名称',
  `version` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '消息版本',
  `queue` varchar(256) NOT NULL DEFAULT '' COMMENT '队列名称',
  `msg` varbinary(8192) NOT NULL DEFAULT '' COMMENT '消息以pb格式存储',
  `retries` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '重试次数',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态[0:未投递;1:投递成功;2:投递失败]',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- 数据导出被取消选择。


-- 导出  表 dm_passport.user 结构
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL COMMENT '更新时间',
  `deleted_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机',
  `passwd` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- 数据导出被取消选择。


-- 导出  表 dm_passport.user_oauth 结构
CREATE TABLE IF NOT EXISTS `user_oauth` (
  `id` bigint(20) unsigned NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `oauth_type` varchar(50) NOT NULL,
  `oauth_id` bigint(20) NOT NULL,
  `access_token` varchar(50) NOT NULL,
  `expires_in` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

-- 数据导出被取消选择。
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
