SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS `task`;
CREATE TABLE `task` (
  `Id` int NOT NULL AUTO_INCREMENT,
  `task_group_id` int NOT NULL DEFAULT '0' COMMENT '任务组ID',
  `start_at` datetime(6) NOT NULL COMMENT '开始时间',
  `run_speed` int NOT NULL COMMENT '运行耗时',
  `client_ip` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '客户端IP',
  `progress` int NOT NULL COMMENT '进度0-100',
  `status` tinyint NOT NULL COMMENT '状态',
  `create_at` datetime(6) NOT NULL COMMENT '任务创建时间',
  `caption` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务组标题',
  `job_name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '实现Job的特性名称（客户端识别哪个实现类）',
  `run_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '实际执行时间',
  `client_id` bigint NOT NULL DEFAULT '0' COMMENT '客户端ID',
  `client_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '客户端名称',
  `scheduler_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '调度时间',
  PRIMARY KEY (`Id`) USING BTREE,
  KEY `group_id_status` (`task_group_id`,`status`,`create_at`,`Id`) USING BTREE,
  KEY `task_group_id` (`create_at`,`task_group_id`) USING BTREE,
  KEY `start_at` (`start_at`,`status`) USING BTREE,
  KEY `create_at` (`status`,`create_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1737 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for task_group
-- ----------------------------
DROP TABLE IF EXISTS `task_group`;
CREATE TABLE `task_group` (
  `Id` int NOT NULL AUTO_INCREMENT,
  `caption` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '任务组标题',
  `job_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '实现Job的特性名称（客户端识别哪个实现类）',
  `start_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '开始时间',
  `next_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '下次执行时间',
  `task_id` bigint NOT NULL DEFAULT '0' COMMENT '任务ID',
  `activate_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '活动时间',
  `last_run_at` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '最后一次完成时间',
  `run_speed_avg` bigint NOT NULL DEFAULT '0' COMMENT '运行平均耗时',
  `run_count` int NOT NULL DEFAULT '0' COMMENT '运行次数',
  `is_enable` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否开启',
  `interval_ms` bigint NOT NULL DEFAULT '1000' COMMENT '时间间隔',
  `cron` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '时间定时器表达式',
  `Data` varchar(2048) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '动态参数',
  PRIMARY KEY (`Id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
