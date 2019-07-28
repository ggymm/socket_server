/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50725
Source Host           : localhost:3306
Source Database       : socket_server

Target Server Type    : MYSQL
Target Server Version : 50725
File Encoding         : 65001

Date: 2019-06-23 17:15:42
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `message_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键，自动生成',
  `from_id` bigint(20) DEFAULT NULL COMMENT '来源用户ID',
  `to_id` bigint(20) DEFAULT NULL COMMENT '目标用户ID',
  `create_time` datetime DEFAULT NULL COMMENT '消息创建时间',
  `msg_type` int(11) DEFAULT NULL,
  `chat_type` int(11) DEFAULT NULL,
  `group_id` varchar(200) DEFAULT NULL,
  `content` text COMMENT '消息内容',
  `extras` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`message_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
