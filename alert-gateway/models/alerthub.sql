CREATE DATABASE IF NOT EXISTS alerthub;

/*
 Navicat Premium Data Transfer

 Source Server         : doraemon
 Source Server Type    : MySQL
 Source Server Version : 50641
 Source Host           : 172.30.31.126:3306
 Source Schema         : doraemon

 Target Server Type    : MySQL
 Target Server Version : 50641
 File Encoding         : 65001

 Date: 04/07/2020 23:32:26
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for alert
-- ----------------------------
DROP TABLE IF EXISTS `alert`;
CREATE TABLE `alert` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rule_id` bigint(20) NOT NULL,
  `labels` varchar(4095) NOT NULL DEFAULT '',
  `value` double NOT NULL DEFAULT '0',
  `count` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `summary` varchar(1023) NOT NULL DEFAULT '',
  `description` varchar(1023) NOT NULL DEFAULT '',
  `hostname` varchar(255) NOT NULL DEFAULT '',
  `confirmed_by` varchar(1023) NOT NULL DEFAULT '',
  `fired_at` datetime NOT NULL,
  `confirmed_at` datetime DEFAULT NULL,
  `confirmed_before` datetime DEFAULT NULL,
  `resolved_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ruleid_labels_firedat` (`rule_id`,`labels`(255),`fired_at`),
  KEY `alert_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=2038 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `service_id` bigint(20) NOT NULL DEFAULT '0',
  `idc` varchar(255) NOT NULL DEFAULT '',
  `proto` varchar(255) NOT NULL DEFAULT '',
  `auto` varchar(255) NOT NULL DEFAULT '',
  `port` int(11) NOT NULL DEFAULT '0',
  `metric` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `user` varchar(1023) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for host
-- ----------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `mid` bigint(20) NOT NULL DEFAULT '0',
  `hostname` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `mid` (`mid`,`hostname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for maintain
-- ----------------------------
DROP TABLE IF EXISTS `maintain`;
CREATE TABLE `maintain` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `flag` tinyint(1) NOT NULL DEFAULT '0',
  `time_start` varchar(15) NOT NULL DEFAULT '',
  `time_end` varchar(15) NOT NULL DEFAULT '',
  `month` int(11) NOT NULL DEFAULT '0',
  `day_start` tinyint(4) NOT NULL DEFAULT '0',
  `day_end` tinyint(4) NOT NULL DEFAULT '0',
  `valid` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `maintain_valid_day_start_day_end_flag_time_start_time_end` (`valid`,`day_start`,`day_end`,`flag`,`time_start`,`time_end`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for manage
-- ----------------------------
DROP TABLE IF EXISTS `manage`;
CREATE TABLE `manage` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `servicename` varchar(255) NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `servicename` (`servicename`),
  KEY `manage_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for plan
-- ----------------------------
DROP TABLE IF EXISTS `plan`;
CREATE TABLE `plan` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rule_labels` varchar(255) NOT NULL DEFAULT '',
  `description` varchar(1023) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for plan_receiver
-- ----------------------------
DROP TABLE IF EXISTS `plan_receiver`;
CREATE TABLE `plan_receiver` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `plan_id` bigint(20) NOT NULL,
  `start_time` varchar(31) NOT NULL DEFAULT '',
  `end_time` varchar(31) NOT NULL DEFAULT '',
  `start` int(11) NOT NULL DEFAULT '0',
  `period` int(11) NOT NULL DEFAULT '0',
  `expression` varchar(1023) NOT NULL DEFAULT '',
  `reverse_polish_notation` varchar(1023) NOT NULL DEFAULT '',
  `user` varchar(1023) NOT NULL DEFAULT '',
  `group` varchar(1023) NOT NULL DEFAULT '',
  `duty_group` varchar(255) NOT NULL DEFAULT '',
  `method` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `plan_receiver_plan_id` (`plan_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for prom
-- ----------------------------
DROP TABLE IF EXISTS `prom`;
CREATE TABLE `prom` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(1023) NOT NULL DEFAULT '',
  `url` varchar(1023) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rule
-- ----------------------------
DROP TABLE IF EXISTS `rule`;
CREATE TABLE `rule` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `expr` varchar(1023) NOT NULL DEFAULT '',
  `op` varchar(31) NOT NULL DEFAULT '',
  `value` varchar(1023) NOT NULL DEFAULT '',
  `for` varchar(1023) NOT NULL DEFAULT '',
  `summary` varchar(1023) NOT NULL DEFAULT '',
  `description` varchar(1023) NOT NULL DEFAULT '',
  `prom_id` bigint(20) NOT NULL,
  `plan_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(1023) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO `alerthub`.`users`(`id`, `name`, `password`) VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e');

