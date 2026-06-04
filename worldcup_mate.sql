/*
 Navicat Premium Dump SQL

 Source Server         : dev
 Source Server Type    : MySQL
 Source Server Version : 80032 (8.0.32)
 Source Host           : localhost:3310
 Source Schema         : worldcup_mate

 Target Server Type    : MySQL
 Target Server Version : 80032 (8.0.32)
 File Encoding         : 65001

 Date: 04/06/2026 15:31:36
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cities
-- ----------------------------
DROP TABLE IF EXISTS `cities`;
CREATE TABLE `cities`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name_en` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `country` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `timezone` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cities_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_cities_name_en`(`name_en` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cities
-- ----------------------------
INSERT INTO `cities` VALUES (1, '墨西哥城', 'Mexico City', 'Mexico', 'America/Mexico_City', '', '2026-06-02 11:53:03.304', '2026-06-04 15:21:00.340', NULL);
INSERT INTO `cities` VALUES (2, '纽约', 'New York', 'USA', 'America/New_York', '', '2026-06-02 11:53:03.308', '2026-06-02 11:53:03.308', NULL);
INSERT INTO `cities` VALUES (3, '洛杉矶', 'Los Angeles', 'USA', 'America/Los_Angeles', '', '2026-06-02 11:53:03.313', '2026-06-04 15:21:00.470', NULL);
INSERT INTO `cities` VALUES (4, '达拉斯', 'Dallas', 'USA', 'America/Chicago', '', '2026-06-02 11:53:03.317', '2026-06-04 15:21:00.534', NULL);
INSERT INTO `cities` VALUES (5, '西雅图', 'Seattle', 'USA', 'America/Los_Angeles', '', '2026-06-02 11:53:03.321', '2026-06-04 15:21:00.387', NULL);
INSERT INTO `cities` VALUES (6, '迈阿密', 'Miami', 'USA', 'America/New_York', '', '2026-06-02 11:53:03.326', '2026-06-04 15:21:00.573', NULL);
INSERT INTO `cities` VALUES (7, '多伦多', 'Toronto', 'Canada', 'America/Toronto', '', '2026-06-02 11:53:03.331', '2026-06-04 15:21:00.172', NULL);
INSERT INTO `cities` VALUES (8, '温哥华', 'Vancouver', 'Canada', 'America/Vancouver', '', '2026-06-02 11:53:03.336', '2026-06-04 15:21:00.430', NULL);
INSERT INTO `cities` VALUES (9, 'TBD', 'TBD', '', 'UTC', '', '2026-06-02 15:55:43.020', '2026-06-02 15:55:43.020', NULL);
INSERT INTO `cities` VALUES (10, '瓜达拉哈拉', 'Guadalajara', 'Mexico', 'America/Mexico_City', '', '2026-06-04 10:01:09.693', '2026-06-04 15:20:59.739', NULL);
INSERT INTO `cities` VALUES (11, '堪萨斯城', 'Kansas City', 'USA', 'America/Chicago', '', '2026-06-04 10:01:09.705', '2026-06-04 15:21:00.510', NULL);
INSERT INTO `cities` VALUES (12, '亚特兰大', 'Atlanta', 'USA', 'America/New_York', '', '2026-06-04 10:01:09.716', '2026-06-04 15:21:00.554', NULL);
INSERT INTO `cities` VALUES (13, '蒙特雷', 'Monterrey', 'Mexico', 'America/Monterrey', '', '2026-06-04 10:01:09.729', '2026-06-04 15:21:00.007', NULL);
INSERT INTO `cities` VALUES (14, '旧金山湾区', 'San Francisco Bay Area', 'USA', 'America/Los_Angeles', '', '2026-06-04 10:01:09.740', '2026-06-04 15:21:00.130', NULL);
INSERT INTO `cities` VALUES (15, '新泽西', 'New Jersey', 'USA', 'America/New_York', '', '2026-06-04 10:01:09.764', '2026-06-04 15:21:00.594', NULL);
INSERT INTO `cities` VALUES (16, '休斯敦', 'Houston', 'USA', 'America/Chicago', '', '2026-06-04 10:01:09.823', '2026-06-04 15:21:00.277', NULL);
INSERT INTO `cities` VALUES (17, '波士顿', 'Boston', 'USA', 'America/New_York', '', '2026-06-04 10:01:09.895', '2026-06-04 15:21:00.450', NULL);
INSERT INTO `cities` VALUES (18, '费城', 'Philadelphia', 'USA', 'America/New_York', '', '2026-06-04 10:01:10.087', '2026-06-04 15:21:00.299', NULL);
INSERT INTO `cities` VALUES (19, '墨西哥城', 'Mexico City', 'Mexico', 'America/Mexico_City', '', '2026-06-04 10:04:12.124', '2026-06-04 10:04:12.124', NULL);
INSERT INTO `cities` VALUES (20, '洛杉矶', 'Los Angeles', 'USA', 'America/Los_Angeles', '', '2026-06-04 10:04:12.131', '2026-06-04 10:04:12.131', NULL);
INSERT INTO `cities` VALUES (21, '达拉斯', 'Dallas', 'USA', 'America/Chicago', '', '2026-06-04 10:04:12.135', '2026-06-04 10:04:12.135', NULL);
INSERT INTO `cities` VALUES (22, '西雅图', 'Seattle', 'USA', 'America/Los_Angeles', '', '2026-06-04 10:04:12.139', '2026-06-04 10:04:12.139', NULL);
INSERT INTO `cities` VALUES (23, '迈阿密', 'Miami', 'USA', 'America/New_York', '', '2026-06-04 10:04:12.144', '2026-06-04 10:04:12.144', NULL);
INSERT INTO `cities` VALUES (24, '多伦多', 'Toronto', 'Canada', 'America/Toronto', '', '2026-06-04 10:04:12.149', '2026-06-04 10:04:12.149', NULL);
INSERT INTO `cities` VALUES (25, '温哥华', 'Vancouver', 'Canada', 'America/Vancouver', '', '2026-06-04 10:04:12.153', '2026-06-04 10:04:12.153', NULL);

-- ----------------------------
-- Table structure for group_standings
-- ----------------------------
DROP TABLE IF EXISTS `group_standings`;
CREATE TABLE `group_standings`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `group_id` bigint UNSIGNED NOT NULL,
  `team_id` bigint UNSIGNED NOT NULL,
  `played` bigint NULL DEFAULT 0,
  `won` bigint NULL DEFAULT 0,
  `drawn` bigint NULL DEFAULT 0,
  `lost` bigint NULL DEFAULT 0,
  `goals_for` bigint NULL DEFAULT 0,
  `goals_against` bigint NULL DEFAULT 0,
  `goal_difference` bigint NULL DEFAULT 0,
  `points` bigint NULL DEFAULT 0,
  `rank` bigint NULL DEFAULT 0,
  `qualification_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'unknown',
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_group_team`(`group_id` ASC, `team_id` ASC) USING BTREE,
  INDEX `fk_group_standings_team`(`team_id` ASC) USING BTREE,
  CONSTRAINT `fk_group_standings_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_group_standings_team` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of group_standings
-- ----------------------------

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `stage` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'group',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_groups_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of groups
-- ----------------------------
INSERT INTO `groups` VALUES (1, 'Group A', 'group', '2026-06-02 11:53:03.253', '2026-06-02 11:53:03.253');
INSERT INTO `groups` VALUES (2, 'Group B', 'group', '2026-06-02 11:53:03.257', '2026-06-02 11:53:03.257');
INSERT INTO `groups` VALUES (3, 'Group C', 'group', '2026-06-02 11:53:03.261', '2026-06-02 11:53:03.261');
INSERT INTO `groups` VALUES (4, 'Group D', 'group', '2026-06-02 11:53:03.266', '2026-06-02 11:53:03.266');
INSERT INTO `groups` VALUES (5, 'Group E', 'group', '2026-06-02 11:53:03.270', '2026-06-02 11:53:03.270');
INSERT INTO `groups` VALUES (6, 'Group F', 'group', '2026-06-02 11:53:03.274', '2026-06-02 11:53:03.274');
INSERT INTO `groups` VALUES (7, 'Group G', 'group', '2026-06-02 11:53:03.278', '2026-06-02 11:53:03.278');
INSERT INTO `groups` VALUES (8, 'Group H', 'group', '2026-06-02 11:53:03.282', '2026-06-02 11:53:03.282');
INSERT INTO `groups` VALUES (9, 'Group I', 'group', '2026-06-02 11:53:03.286', '2026-06-02 11:53:03.286');
INSERT INTO `groups` VALUES (10, 'Group J', 'group', '2026-06-02 11:53:03.290', '2026-06-02 11:53:03.290');
INSERT INTO `groups` VALUES (11, 'Group K', 'group', '2026-06-02 11:53:03.294', '2026-06-02 11:53:03.294');
INSERT INTO `groups` VALUES (12, 'Group L', 'group', '2026-06-02 11:53:03.298', '2026-06-02 11:53:03.298');
INSERT INTO `groups` VALUES (13, 'Tbd', 'group', '2026-06-02 16:17:26.348', '2026-06-02 16:17:26.348');

-- ----------------------------
-- Table structure for matches
-- ----------------------------
DROP TABLE IF EXISTS `matches`;
CREATE TABLE `matches`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `match_no` bigint NOT NULL,
  `home_team_id` bigint UNSIGNED NULL DEFAULT NULL,
  `away_team_id` bigint UNSIGNED NULL DEFAULT NULL,
  `group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `stage` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'group',
  `stadium_id` bigint UNSIGNED NULL DEFAULT NULL,
  `city_id` bigint UNSIGNED NULL DEFAULT NULL,
  `kickoff_time_utc` datetime(3) NOT NULL,
  `home_score` bigint NULL DEFAULT NULL,
  `away_score` bigint NULL DEFAULT NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'scheduled',
  `winner_team_id` bigint UNSIGNED NULL DEFAULT NULL,
  `importance_level` bigint NULL DEFAULT 0,
  `recommend_tag` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `recommend_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `status_detail` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `live_minute` bigint NULL DEFAULT NULL,
  `external_provider` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `external_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `last_synced_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_matches_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_matches_external_provider`(`external_provider` ASC) USING BTREE,
  INDEX `idx_matches_external_id`(`external_id` ASC) USING BTREE,
  INDEX `idx_matches_group_id`(`group_id` ASC) USING BTREE,
  INDEX `idx_matches_city_id`(`city_id` ASC) USING BTREE,
  INDEX `idx_matches_status`(`status` ASC) USING BTREE,
  INDEX `idx_matches_home_team_id`(`home_team_id` ASC) USING BTREE,
  INDEX `idx_matches_away_team_id`(`away_team_id` ASC) USING BTREE,
  INDEX `idx_matches_stage`(`stage` ASC) USING BTREE,
  INDEX `idx_matches_stadium_id`(`stadium_id` ASC) USING BTREE,
  INDEX `idx_matches_kickoff_time_utc`(`kickoff_time_utc` ASC) USING BTREE,
  CONSTRAINT `fk_matches_away_team` FOREIGN KEY (`away_team_id`) REFERENCES `teams` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_matches_city` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_matches_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_matches_home_team` FOREIGN KEY (`home_team_id`) REFERENCES `teams` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_matches_stadium` FOREIGN KEY (`stadium_id`) REFERENCES `stadia` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 109 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of matches
-- ----------------------------
INSERT INTO `matches` VALUES (1, 1, 1, 2, 1, 'group', 1, 1, '2026-06-02 09:00:00.000', NULL, NULL, 'scheduled', NULL, 3, '揭幕战', '', '2026-06-02 13:51:13.312', '2026-06-02 13:51:13.312', '2026-06-03 08:56:46.702', NULL, NULL, NULL, NULL, NULL);
INSERT INTO `matches` VALUES (2, 2, 5, 6, 2, 'group', 2, 2, '2026-06-02 21:00:00.000', NULL, NULL, 'scheduled', NULL, 3, '焦点大战', '', '2026-06-02 13:51:13.318', '2026-06-02 13:51:13.318', '2026-06-03 08:56:46.702', NULL, NULL, NULL, NULL, NULL);
INSERT INTO `matches` VALUES (3, 3, 9, 10, 3, 'group', 3, 3, '2026-06-03 18:00:00.000', NULL, NULL, 'scheduled', NULL, 3, '热门比赛', '', '2026-06-02 13:51:13.322', '2026-06-02 13:51:13.322', '2026-06-03 08:56:46.702', NULL, NULL, NULL, NULL, NULL);
INSERT INTO `matches` VALUES (4, 4, 13, 14, 4, 'group', 5, 5, '2026-06-04 03:00:00.000', NULL, NULL, 'scheduled', NULL, 2, '', '', '2026-06-02 13:51:13.326', '2026-06-02 13:51:13.326', '2026-06-03 08:56:46.702', NULL, NULL, NULL, NULL, NULL);
INSERT INTO `matches` VALUES (5, 537327, 1, 3, 1, 'group', 21, 1, '2026-06-12 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 15:55:43.031', '2026-06-04 15:20:58.432', NULL, 'TIMED', NULL, 'football-data', '537327', '2026-06-04 15:20:58.431');
INSERT INTO `matches` VALUES (6, 537357, 15, 13, 6, 'group', 19, 4, '2026-06-15 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 15:55:43.082', '2026-06-04 15:20:58.630', NULL, 'TIMED', NULL, 'football-data', '537357', '2026-06-04 15:20:58.629');
INSERT INTO `matches` VALUES (7, 537411, 11, 8, 12, 'group', 22, 17, '2026-06-24 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 15:55:43.243', '2026-06-04 15:20:59.382', NULL, 'TIMED', NULL, 'football-data', '537411', '2026-06-04 15:20:59.380');
INSERT INTO `matches` VALUES (8, 537367, 4, 16, 7, 'group', 16, 8, '2026-06-27 11:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 15:55:43.294', '2026-06-04 15:20:59.789', NULL, 'TIMED', NULL, 'football-data', '537367', '2026-06-04 15:20:59.788');
INSERT INTO `matches` VALUES (9, 537328, 417, 418, 1, 'group', 11, 10, '2026-06-12 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.413', '2026-06-04 15:20:58.455', NULL, 'TIMED', NULL, 'football-data', '537328', '2026-06-04 15:20:58.453');
INSERT INTO `matches` VALUES (10, 537333, 2, 419, 2, 'group', 23, 7, '2026-06-13 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.426', '2026-06-04 15:20:58.479', NULL, 'TIMED', NULL, 'football-data', '537333', '2026-06-04 15:20:58.478');
INSERT INTO `matches` VALUES (11, 537345, 7, 420, 4, 'group', 24, 3, '2026-06-13 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.439', '2026-06-04 15:20:58.502', NULL, 'TIMED', NULL, 'football-data', '537345', '2026-06-04 15:20:58.500');
INSERT INTO `matches` VALUES (12, 537334, 421, 422, 2, 'group', 15, 14, '2026-06-14 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.452', '2026-06-04 15:20:58.523', NULL, 'TIMED', NULL, 'football-data', '537334', '2026-06-04 15:20:58.523');
INSERT INTO `matches` VALUES (13, 537339, 9, 423, 3, 'group', 17, 15, '2026-06-14 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.467', '2026-06-04 15:20:58.545', NULL, 'TIMED', NULL, 'football-data', '537339', '2026-06-04 15:20:58.544');
INSERT INTO `matches` VALUES (14, 537340, 424, 425, 3, 'group', 22, 17, '2026-06-14 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.480', '2026-06-04 15:20:58.565', NULL, 'TIMED', NULL, 'football-data', '537340', '2026-06-04 15:20:58.564');
INSERT INTO `matches` VALUES (15, 537346, 426, 427, 4, 'group', 16, 8, '2026-06-14 12:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.494', '2026-06-04 15:20:58.587', NULL, 'TIMED', NULL, 'football-data', '537346', '2026-06-04 15:20:58.585');
INSERT INTO `matches` VALUES (16, 537351, 14, 428, 5, 'group', 20, 16, '2026-06-15 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.508', '2026-06-04 15:20:58.608', NULL, 'TIMED', NULL, 'football-data', '537351', '2026-06-04 15:20:58.607');
INSERT INTO `matches` VALUES (17, 537352, 429, 430, 5, 'group', 25, 18, '2026-06-15 07:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.534', '2026-06-04 15:20:58.651', NULL, 'TIMED', NULL, 'football-data', '537352', '2026-06-04 15:20:58.649');
INSERT INTO `matches` VALUES (18, 537358, 431, 432, 6, 'group', 14, 13, '2026-06-15 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.549', '2026-06-04 15:20:58.672', NULL, 'TIMED', NULL, 'football-data', '537358', '2026-06-04 15:20:58.671');
INSERT INTO `matches` VALUES (19, 537369, 10, 433, 8, 'group', 13, 12, '2026-06-16 00:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.563', '2026-06-04 15:20:58.692', NULL, 'TIMED', NULL, 'football-data', '537369', '2026-06-04 15:20:58.691');
INSERT INTO `matches` VALUES (20, 537363, 16, 434, 7, 'group', 10, 5, '2026-06-16 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.577', '2026-06-04 15:20:58.713', NULL, 'TIMED', NULL, 'football-data', '537363', '2026-06-04 15:20:58.712');
INSERT INTO `matches` VALUES (21, 537370, 435, 436, 8, 'group', 18, 6, '2026-06-16 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.590', '2026-06-04 15:20:58.735', NULL, 'TIMED', NULL, 'football-data', '537370', '2026-06-04 15:20:58.734');
INSERT INTO `matches` VALUES (22, 537364, 437, 4, 7, 'group', 24, 3, '2026-06-16 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.605', '2026-06-04 15:20:58.755', NULL, 'TIMED', NULL, 'football-data', '537364', '2026-06-04 15:20:58.754');
INSERT INTO `matches` VALUES (23, 537391, 6, 438, 9, 'group', 17, 15, '2026-06-17 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.619', '2026-06-04 15:20:58.775', NULL, 'TIMED', NULL, 'football-data', '537391', '2026-06-04 15:20:58.775');
INSERT INTO `matches` VALUES (24, 537392, 439, 440, 9, 'group', 22, 17, '2026-06-17 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.633', '2026-06-04 15:20:58.796', NULL, 'TIMED', NULL, 'football-data', '537392', '2026-06-04 15:20:58.795');
INSERT INTO `matches` VALUES (25, 537397, 5, 441, 10, 'group', 12, 11, '2026-06-17 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.647', '2026-06-04 15:20:58.818', NULL, 'TIMED', NULL, 'football-data', '537397', '2026-06-04 15:20:58.816');
INSERT INTO `matches` VALUES (26, 537398, 442, 443, 10, 'group', 15, 14, '2026-06-17 12:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.661', '2026-06-04 15:20:58.840', NULL, 'TIMED', NULL, 'football-data', '537398', '2026-06-04 15:20:58.839');
INSERT INTO `matches` VALUES (27, 537403, 12, 444, 11, 'group', 20, 16, '2026-06-18 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.675', '2026-06-04 15:20:58.860', NULL, 'TIMED', NULL, 'football-data', '537403', '2026-06-04 15:20:58.859');
INSERT INTO `matches` VALUES (28, 537409, 11, 445, 12, 'group', 19, 4, '2026-06-18 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.688', '2026-06-04 15:20:58.880', NULL, 'TIMED', NULL, 'football-data', '537409', '2026-06-04 15:20:58.879');
INSERT INTO `matches` VALUES (29, 537410, 8, 446, 12, 'group', 23, 7, '2026-06-18 07:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.701', '2026-06-04 15:20:58.905', NULL, 'TIMED', NULL, 'football-data', '537410', '2026-06-04 15:20:58.905');
INSERT INTO `matches` VALUES (30, 537404, 447, 448, 11, 'group', 21, 1, '2026-06-18 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.713', '2026-06-04 15:20:58.925', NULL, 'TIMED', NULL, 'football-data', '537404', '2026-06-04 15:20:58.924');
INSERT INTO `matches` VALUES (31, 537329, 418, 3, 1, 'group', 13, 12, '2026-06-19 00:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.726', '2026-06-04 15:20:58.945', NULL, 'TIMED', NULL, 'football-data', '537329', '2026-06-04 15:20:58.944');
INSERT INTO `matches` VALUES (32, 537335, 422, 419, 2, 'group', 24, 3, '2026-06-19 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.739', '2026-06-04 15:20:58.967', NULL, 'TIMED', NULL, 'football-data', '537335', '2026-06-04 15:20:58.966');
INSERT INTO `matches` VALUES (33, 537336, 2, 421, 2, 'group', 16, 8, '2026-06-19 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.751', '2026-06-04 15:20:58.987', NULL, 'TIMED', NULL, 'football-data', '537336', '2026-06-04 15:20:58.986');
INSERT INTO `matches` VALUES (34, 537330, 1, 417, 1, 'group', 11, 10, '2026-06-19 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.764', '2026-06-04 15:20:59.008', NULL, 'TIMED', NULL, 'football-data', '537330', '2026-06-04 15:20:59.007');
INSERT INTO `matches` VALUES (35, 537348, 7, 426, 4, 'group', 10, 5, '2026-06-20 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.776', '2026-06-04 15:20:59.030', NULL, 'TIMED', NULL, 'football-data', '537348', '2026-06-04 15:20:59.029');
INSERT INTO `matches` VALUES (36, 537342, 425, 423, 3, 'group', 22, 17, '2026-06-20 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.790', '2026-06-04 15:20:59.052', NULL, 'TIMED', NULL, 'football-data', '537342', '2026-06-04 15:20:59.050');
INSERT INTO `matches` VALUES (37, 537341, 9, 424, 3, 'group', 25, 18, '2026-06-20 08:30:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.802', '2026-06-04 15:20:59.074', NULL, 'TIMED', NULL, 'football-data', '537341', '2026-06-04 15:20:59.073');
INSERT INTO `matches` VALUES (38, 537347, 427, 420, 4, 'group', 15, 14, '2026-06-20 11:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.814', '2026-06-04 15:20:59.095', NULL, 'TIMED', NULL, 'football-data', '537347', '2026-06-04 15:20:59.094');
INSERT INTO `matches` VALUES (39, 537359, 15, 431, 6, 'group', 20, 16, '2026-06-21 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.827', '2026-06-04 15:20:59.115', NULL, 'TIMED', NULL, 'football-data', '537359', '2026-06-04 15:20:59.115');
INSERT INTO `matches` VALUES (40, 537353, 14, 429, 5, 'group', 23, 7, '2026-06-21 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.839', '2026-06-04 15:20:59.135', NULL, 'TIMED', NULL, 'football-data', '537353', '2026-06-04 15:20:59.135');
INSERT INTO `matches` VALUES (41, 537354, 430, 428, 5, 'group', 12, 11, '2026-06-21 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.852', '2026-06-04 15:20:59.155', NULL, 'TIMED', NULL, 'football-data', '537354', '2026-06-04 15:20:59.154');
INSERT INTO `matches` VALUES (42, 537360, 432, 13, 6, 'group', 14, 13, '2026-06-21 12:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.864', '2026-06-04 15:20:59.175', NULL, 'TIMED', NULL, 'football-data', '537360', '2026-06-04 15:20:59.174');
INSERT INTO `matches` VALUES (43, 537371, 10, 435, 8, 'group', 13, 12, '2026-06-22 00:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.876', '2026-06-04 15:20:59.195', NULL, 'TIMED', NULL, 'football-data', '537371', '2026-06-04 15:20:59.195');
INSERT INTO `matches` VALUES (44, 537365, 16, 437, 7, 'group', 24, 3, '2026-06-22 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.890', '2026-06-04 15:20:59.215', NULL, 'TIMED', NULL, 'football-data', '537365', '2026-06-04 15:20:59.215');
INSERT INTO `matches` VALUES (45, 537372, 436, 433, 8, 'group', 18, 6, '2026-06-22 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.903', '2026-06-04 15:20:59.235', NULL, 'TIMED', NULL, 'football-data', '537372', '2026-06-04 15:20:59.234');
INSERT INTO `matches` VALUES (46, 537366, 4, 434, 7, 'group', 16, 8, '2026-06-22 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.916', '2026-06-04 15:20:59.256', NULL, 'TIMED', NULL, 'football-data', '537366', '2026-06-04 15:20:59.255');
INSERT INTO `matches` VALUES (47, 537399, 5, 442, 10, 'group', 19, 4, '2026-06-23 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.931', '2026-06-04 15:20:59.276', NULL, 'TIMED', NULL, 'football-data', '537399', '2026-06-04 15:20:59.275');
INSERT INTO `matches` VALUES (48, 537393, 6, 439, 9, 'group', 25, 18, '2026-06-23 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.946', '2026-06-04 15:20:59.299', NULL, 'TIMED', NULL, 'football-data', '537393', '2026-06-04 15:20:59.299');
INSERT INTO `matches` VALUES (49, 537394, 440, 438, 9, 'group', 17, 15, '2026-06-23 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.960', '2026-06-04 15:20:59.320', NULL, 'TIMED', NULL, 'football-data', '537394', '2026-06-04 15:20:59.320');
INSERT INTO `matches` VALUES (50, 537400, 443, 441, 10, 'group', 15, 14, '2026-06-23 11:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.975', '2026-06-04 15:20:59.342', NULL, 'TIMED', NULL, 'football-data', '537400', '2026-06-04 15:20:59.341');
INSERT INTO `matches` VALUES (51, 537405, 12, 447, 11, 'group', 20, 16, '2026-06-24 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:25.989', '2026-06-04 15:20:59.362', NULL, 'TIMED', NULL, 'football-data', '537405', '2026-06-04 15:20:59.362');
INSERT INTO `matches` VALUES (52, 537412, 446, 445, 12, 'group', 23, 7, '2026-06-24 07:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.018', '2026-06-04 15:20:59.402', NULL, 'TIMED', NULL, 'football-data', '537412', '2026-06-04 15:20:59.400');
INSERT INTO `matches` VALUES (53, 537406, 448, 444, 11, 'group', 11, 10, '2026-06-24 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.031', '2026-06-04 15:20:59.422', NULL, 'TIMED', NULL, 'football-data', '537406', '2026-06-04 15:20:59.421');
INSERT INTO `matches` VALUES (54, 537337, 422, 2, 2, 'group', 16, 8, '2026-06-25 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.045', '2026-06-04 15:20:59.443', NULL, 'TIMED', NULL, 'football-data', '537337', '2026-06-04 15:20:59.441');
INSERT INTO `matches` VALUES (55, 537338, 419, 421, 2, 'group', 10, 5, '2026-06-25 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.057', '2026-06-04 15:20:59.464', NULL, 'TIMED', NULL, 'football-data', '537338', '2026-06-04 15:20:59.463');
INSERT INTO `matches` VALUES (56, 537344, 423, 424, 3, 'group', 13, 12, '2026-06-25 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.071', '2026-06-04 15:20:59.485', NULL, 'TIMED', NULL, 'football-data', '537344', '2026-06-04 15:20:59.484');
INSERT INTO `matches` VALUES (57, 537343, 425, 9, 3, 'group', 18, 6, '2026-06-25 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.086', '2026-06-04 15:20:59.506', NULL, 'TIMED', NULL, 'football-data', '537343', '2026-06-04 15:20:59.506');
INSERT INTO `matches` VALUES (58, 537331, 418, 1, 1, 'group', 21, 1, '2026-06-25 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.099', '2026-06-04 15:20:59.528', NULL, 'TIMED', NULL, 'football-data', '537331', '2026-06-04 15:20:59.528');
INSERT INTO `matches` VALUES (59, 537332, 3, 417, 1, 'group', 14, 13, '2026-06-25 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.111', '2026-06-04 15:20:59.548', NULL, 'TIMED', NULL, 'football-data', '537332', '2026-06-04 15:20:59.547');
INSERT INTO `matches` VALUES (60, 537355, 430, 14, 5, 'group', 17, 15, '2026-06-26 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.126', '2026-06-04 15:20:59.568', NULL, 'TIMED', NULL, 'football-data', '537355', '2026-06-04 15:20:59.567');
INSERT INTO `matches` VALUES (61, 537356, 428, 429, 5, 'group', 25, 18, '2026-06-26 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.140', '2026-06-04 15:20:59.590', NULL, 'TIMED', NULL, 'football-data', '537356', '2026-06-04 15:20:59.587');
INSERT INTO `matches` VALUES (62, 537361, 432, 15, 6, 'group', 12, 11, '2026-06-26 07:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.153', '2026-06-04 15:20:59.614', NULL, 'TIMED', NULL, 'football-data', '537361', '2026-06-04 15:20:59.613');
INSERT INTO `matches` VALUES (63, 537362, 13, 431, 6, 'group', 19, 4, '2026-06-26 07:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.165', '2026-06-04 15:20:59.637', NULL, 'TIMED', NULL, 'football-data', '537362', '2026-06-04 15:20:59.636');
INSERT INTO `matches` VALUES (64, 537349, 427, 7, 4, 'group', 24, 3, '2026-06-26 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.177', '2026-06-04 15:20:59.661', NULL, 'TIMED', NULL, 'football-data', '537349', '2026-06-04 15:20:59.659');
INSERT INTO `matches` VALUES (65, 537350, 420, 426, 4, 'group', 15, 14, '2026-06-26 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.190', '2026-06-04 15:20:59.684', NULL, 'TIMED', NULL, 'football-data', '537350', '2026-06-04 15:20:59.683');
INSERT INTO `matches` VALUES (66, 537395, 440, 6, 9, 'group', 22, 17, '2026-06-27 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.202', '2026-06-04 15:20:59.705', NULL, 'TIMED', NULL, 'football-data', '537395', '2026-06-04 15:20:59.704');
INSERT INTO `matches` VALUES (67, 537396, 438, 439, 9, 'group', 23, 7, '2026-06-27 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.214', '2026-06-04 15:20:59.726', NULL, 'TIMED', NULL, 'football-data', '537396', '2026-06-04 15:20:59.724');
INSERT INTO `matches` VALUES (68, 537373, 436, 10, 8, 'group', 11, 10, '2026-06-27 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.227', '2026-06-04 15:20:59.747', NULL, 'TIMED', NULL, 'football-data', '537373', '2026-06-04 15:20:59.746');
INSERT INTO `matches` VALUES (69, 537374, 433, 435, 8, 'group', 20, 16, '2026-06-27 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.241', '2026-06-04 15:20:59.767', NULL, 'TIMED', NULL, 'football-data', '537374', '2026-06-04 15:20:59.766');
INSERT INTO `matches` VALUES (70, 537368, 434, 437, 7, 'group', 10, 5, '2026-06-27 11:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.267', '2026-06-04 15:20:59.809', NULL, 'TIMED', NULL, 'football-data', '537368', '2026-06-04 15:20:59.808');
INSERT INTO `matches` VALUES (71, 537413, 446, 11, 12, 'group', 17, 15, '2026-06-28 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.281', '2026-06-04 15:20:59.831', NULL, 'TIMED', NULL, 'football-data', '537413', '2026-06-04 15:20:59.830');
INSERT INTO `matches` VALUES (72, 537414, 445, 8, 12, 'group', 25, 18, '2026-06-28 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.293', '2026-06-04 15:20:59.852', NULL, 'TIMED', NULL, 'football-data', '537414', '2026-06-04 15:20:59.851');
INSERT INTO `matches` VALUES (73, 537407, 448, 12, 11, 'group', 18, 6, '2026-06-28 07:30:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.307', '2026-06-04 15:20:59.872', NULL, 'TIMED', NULL, 'football-data', '537407', '2026-06-04 15:20:59.871');
INSERT INTO `matches` VALUES (74, 537408, 444, 447, 11, 'group', 13, 12, '2026-06-28 07:30:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.319', '2026-06-04 15:20:59.893', NULL, 'TIMED', NULL, 'football-data', '537408', '2026-06-04 15:20:59.892');
INSERT INTO `matches` VALUES (75, 537401, 443, 5, 10, 'group', 19, 4, '2026-06-28 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.331', '2026-06-04 15:20:59.914', NULL, 'TIMED', NULL, 'football-data', '537401', '2026-06-04 15:20:59.913');
INSERT INTO `matches` VALUES (76, 537402, 441, 442, 10, 'group', 12, 11, '2026-06-28 10:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.345', '2026-06-04 15:20:59.934', NULL, 'TIMED', NULL, 'football-data', '537402', '2026-06-04 15:20:59.933');
INSERT INTO `matches` VALUES (77, 537417, 449, 449, NULL, 'round_of_32', 24, 3, '2026-06-29 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.361', '2026-06-04 15:20:59.954', NULL, 'TIMED', NULL, 'football-data', '537417', '2026-06-04 15:20:59.953');
INSERT INTO `matches` VALUES (78, 537423, 449, 449, NULL, 'round_of_32', 20, 16, '2026-06-30 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.373', '2026-06-04 15:20:59.975', NULL, 'TIMED', NULL, 'football-data', '537423', '2026-06-04 15:20:59.974');
INSERT INTO `matches` VALUES (79, 537415, 449, 449, NULL, 'round_of_32', 22, 17, '2026-06-30 04:30:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.397', '2026-06-04 15:20:59.995', NULL, 'TIMED', NULL, 'football-data', '537415', '2026-06-04 15:20:59.994');
INSERT INTO `matches` VALUES (80, 537418, 449, 449, NULL, 'round_of_32', 14, 13, '2026-06-30 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.423', '2026-06-04 15:21:00.015', NULL, 'TIMED', NULL, 'football-data', '537418', '2026-06-04 15:21:00.015');
INSERT INTO `matches` VALUES (81, 537424, 449, 449, NULL, 'round_of_32', 19, 4, '2026-07-01 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.443', '2026-06-04 15:21:00.036', NULL, 'TIMED', NULL, 'football-data', '537424', '2026-06-04 15:21:00.035');
INSERT INTO `matches` VALUES (82, 537416, 449, 449, NULL, 'round_of_32', 17, 15, '2026-07-01 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.465', '2026-06-04 15:21:00.057', NULL, 'TIMED', NULL, 'football-data', '537416', '2026-06-04 15:21:00.056');
INSERT INTO `matches` VALUES (83, 537425, 449, 449, NULL, 'round_of_32', 21, 1, '2026-07-01 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.487', '2026-06-04 15:21:00.077', NULL, 'TIMED', NULL, 'football-data', '537425', '2026-06-04 15:21:00.075');
INSERT INTO `matches` VALUES (84, 537426, 449, 449, NULL, 'round_of_32', 13, 12, '2026-07-02 00:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.508', '2026-06-04 15:21:00.095', NULL, 'TIMED', NULL, 'football-data', '537426', '2026-06-04 15:21:00.094');
INSERT INTO `matches` VALUES (85, 537422, 449, 449, NULL, 'round_of_32', 10, 5, '2026-07-02 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.531', '2026-06-04 15:21:00.115', NULL, 'TIMED', NULL, 'football-data', '537422', '2026-06-04 15:21:00.114');
INSERT INTO `matches` VALUES (86, 537421, 449, 449, NULL, 'round_of_32', 15, 14, '2026-07-02 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.553', '2026-06-04 15:21:00.138', NULL, 'TIMED', NULL, 'football-data', '537421', '2026-06-04 15:21:00.137');
INSERT INTO `matches` VALUES (87, 537420, 449, 449, NULL, 'round_of_32', 24, 3, '2026-07-03 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.576', '2026-06-04 15:21:00.160', NULL, 'TIMED', NULL, 'football-data', '537420', '2026-06-04 15:21:00.159');
INSERT INTO `matches` VALUES (88, 537419, 449, 449, NULL, 'round_of_32', 23, 7, '2026-07-03 07:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.599', '2026-06-04 15:21:00.180', NULL, 'TIMED', NULL, 'football-data', '537419', '2026-06-04 15:21:00.179');
INSERT INTO `matches` VALUES (89, 537429, 449, 449, NULL, 'round_of_32', 16, 8, '2026-07-03 11:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.618', '2026-06-04 15:21:00.201', NULL, 'TIMED', NULL, 'football-data', '537429', '2026-06-04 15:21:00.201');
INSERT INTO `matches` VALUES (90, 537428, 449, 449, NULL, 'round_of_32', 19, 4, '2026-07-04 02:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.632', '2026-06-04 15:21:00.225', NULL, 'TIMED', NULL, 'football-data', '537428', '2026-06-04 15:21:00.223');
INSERT INTO `matches` VALUES (91, 537427, 449, 449, NULL, 'round_of_32', 18, 6, '2026-07-04 06:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.645', '2026-06-04 15:21:00.245', NULL, 'TIMED', NULL, 'football-data', '537427', '2026-06-04 15:21:00.244');
INSERT INTO `matches` VALUES (92, 537430, 449, 449, NULL, 'round_of_32', 12, 11, '2026-07-04 09:30:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.657', '2026-06-04 15:21:00.265', NULL, 'TIMED', NULL, 'football-data', '537430', '2026-06-04 15:21:00.264');
INSERT INTO `matches` VALUES (93, 537376, 449, 449, NULL, 'round_of_16', 20, 16, '2026-07-05 01:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.669', '2026-06-04 15:21:00.285', NULL, 'TIMED', NULL, 'football-data', '537376', '2026-06-04 15:21:00.284');
INSERT INTO `matches` VALUES (94, 537375, 449, 449, NULL, 'round_of_16', 25, 18, '2026-07-05 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.682', '2026-06-04 15:21:00.307', NULL, 'TIMED', NULL, 'football-data', '537375', '2026-06-04 15:21:00.307');
INSERT INTO `matches` VALUES (95, 537377, 449, 449, NULL, 'round_of_16', 17, 15, '2026-07-06 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.695', '2026-06-04 15:21:00.327', NULL, 'TIMED', NULL, 'football-data', '537377', '2026-06-04 15:21:00.326');
INSERT INTO `matches` VALUES (96, 537378, 449, 449, NULL, 'round_of_16', 21, 1, '2026-07-06 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.706', '2026-06-04 15:21:00.348', NULL, 'TIMED', NULL, 'football-data', '537378', '2026-06-04 15:21:00.348');
INSERT INTO `matches` VALUES (97, 537379, 449, 449, NULL, 'round_of_16', 19, 4, '2026-07-07 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.719', '2026-06-04 15:21:00.375', NULL, 'TIMED', NULL, 'football-data', '537379', '2026-06-04 15:21:00.373');
INSERT INTO `matches` VALUES (98, 537380, 449, 449, NULL, 'round_of_16', 10, 5, '2026-07-07 08:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.736', '2026-06-04 15:21:00.395', NULL, 'TIMED', NULL, 'football-data', '537380', '2026-06-04 15:21:00.394');
INSERT INTO `matches` VALUES (99, 537381, 449, 449, NULL, 'round_of_16', 13, 12, '2026-07-08 00:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.748', '2026-06-04 15:21:00.416', NULL, 'TIMED', NULL, 'football-data', '537381', '2026-06-04 15:21:00.415');
INSERT INTO `matches` VALUES (100, 537382, 449, 449, NULL, 'round_of_16', 16, 8, '2026-07-08 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.761', '2026-06-04 15:21:00.438', NULL, 'TIMED', NULL, 'football-data', '537382', '2026-06-04 15:21:00.438');
INSERT INTO `matches` VALUES (101, 537383, 449, 449, NULL, 'quarter_final', 22, 17, '2026-07-10 04:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.774', '2026-06-04 15:21:00.458', NULL, 'TIMED', NULL, 'football-data', '537383', '2026-06-04 15:21:00.456');
INSERT INTO `matches` VALUES (102, 537384, 449, 449, NULL, 'quarter_final', 24, 3, '2026-07-11 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.787', '2026-06-04 15:21:00.478', NULL, 'TIMED', NULL, 'football-data', '537384', '2026-06-04 15:21:00.477');
INSERT INTO `matches` VALUES (103, 537385, 449, 449, NULL, 'quarter_final', 18, 6, '2026-07-12 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.799', '2026-06-04 15:21:00.498', NULL, 'TIMED', NULL, 'football-data', '537385', '2026-06-04 15:21:00.497');
INSERT INTO `matches` VALUES (104, 537386, 449, 449, NULL, 'quarter_final', 12, 11, '2026-07-12 09:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.810', '2026-06-04 15:21:00.518', NULL, 'TIMED', NULL, 'football-data', '537386', '2026-06-04 15:21:00.517');
INSERT INTO `matches` VALUES (105, 537387, 449, 449, NULL, 'semi_final', 19, 4, '2026-07-15 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.824', '2026-06-04 15:21:00.542', NULL, 'TIMED', NULL, 'football-data', '537387', '2026-06-04 15:21:00.541');
INSERT INTO `matches` VALUES (106, 537388, 449, 449, NULL, 'semi_final', 13, 12, '2026-07-16 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.837', '2026-06-04 15:21:00.562', NULL, 'TIMED', NULL, 'football-data', '537388', '2026-06-04 15:21:00.561');
INSERT INTO `matches` VALUES (107, 537389, 449, 449, NULL, 'third_place', 18, 6, '2026-07-19 05:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.849', '2026-06-04 15:21:00.582', NULL, 'TIMED', NULL, 'football-data', '537389', '2026-06-04 15:21:00.581');
INSERT INTO `matches` VALUES (108, 537390, 449, 449, NULL, 'final', 17, 15, '2026-07-20 03:00:00.000', NULL, NULL, 'scheduled', NULL, 0, '', '', '2026-06-02 16:17:26.861', '2026-06-04 15:21:00.604', NULL, 'TIMED', NULL, 'football-data', '537390', '2026-06-04 15:21:00.603');

-- ----------------------------
-- Table structure for notifications
-- ----------------------------
DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'reminder',
  `is_read` tinyint(1) NULL DEFAULT 0,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_notifications_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_notifications_user_read`(`user_id` ASC, `is_read` ASC) USING BTREE,
  INDEX `idx_notifications_created_at`(`created_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of notifications
-- ----------------------------

-- ----------------------------
-- Table structure for reminders
-- ----------------------------
DROP TABLE IF EXISTS `reminders`;
CREATE TABLE `reminders`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `match_id` bigint UNSIGNED NOT NULL,
  `remind_before_minutes` bigint NULL DEFAULT 30,
  `remind_at` datetime(3) NOT NULL,
  `channel` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'site',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'pending',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_reminders_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_reminders_user_id`(`user_id` ASC) USING BTREE,
  INDEX `idx_reminders_match_id`(`match_id` ASC) USING BTREE,
  INDEX `idx_reminders_remind_at`(`remind_at` ASC) USING BTREE,
  INDEX `idx_reminders_status`(`status` ASC) USING BTREE,
  CONSTRAINT `fk_reminders_match` FOREIGN KEY (`match_id`) REFERENCES `matches` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_reminders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of reminders
-- ----------------------------
INSERT INTO `reminders` VALUES (1, 2, 5, 30, '2026-06-12 02:30:00.000', 'site', 'pending', '2026-06-03 17:22:22.596', '2026-06-03 17:22:22.596', '2026-06-03 17:22:25.822');
INSERT INTO `reminders` VALUES (2, 2, 9, 30, '2026-06-12 09:30:00.000', 'site', 'pending', '2026-06-03 17:22:26.576', '2026-06-03 17:22:26.576', NULL);

-- ----------------------------
-- Table structure for stadia
-- ----------------------------
DROP TABLE IF EXISTS `stadia`;
CREATE TABLE `stadia`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name_en` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `city_id` bigint UNSIGNED NULL DEFAULT NULL,
  `capacity` bigint NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_stadia_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_stadia_name_en`(`name_en` ASC) USING BTREE,
  INDEX `idx_stadia_city_id`(`city_id` ASC) USING BTREE,
  CONSTRAINT `fk_stadia_city` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 26 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of stadia
-- ----------------------------
INSERT INTO `stadia` VALUES (1, 'Estadio Azteca', 'Estadio Azteca', 1, 87000, '', '2026-06-02 11:53:03.341', '2026-06-02 11:53:03.341', NULL);
INSERT INTO `stadia` VALUES (2, 'MetLife Stadium', 'MetLife Stadium', 2, 82500, '', '2026-06-02 11:53:03.346', '2026-06-02 11:53:03.346', NULL);
INSERT INTO `stadia` VALUES (3, 'SoFi Stadium', 'SoFi Stadium', 3, 70000, '', '2026-06-02 11:53:03.350', '2026-06-02 11:53:03.350', NULL);
INSERT INTO `stadia` VALUES (4, 'AT&T Stadium', 'AT&T Stadium', 4, 80000, '', '2026-06-02 11:53:03.356', '2026-06-02 11:53:03.356', NULL);
INSERT INTO `stadia` VALUES (5, 'Lumen Field', 'Lumen Field', 5, 69000, '', '2026-06-02 11:53:03.361', '2026-06-02 11:53:03.361', NULL);
INSERT INTO `stadia` VALUES (6, 'Hard Rock Stadium', 'Hard Rock Stadium', 6, 65000, '', '2026-06-02 11:53:03.366', '2026-06-02 11:53:03.366', NULL);
INSERT INTO `stadia` VALUES (7, 'BMO Field', 'BMO Field', 7, 30000, '', '2026-06-02 11:53:03.370', '2026-06-02 11:53:03.370', NULL);
INSERT INTO `stadia` VALUES (8, 'BC Place', 'BC Place', 8, 54000, '', '2026-06-02 11:53:03.376', '2026-06-02 11:53:03.376', NULL);
INSERT INTO `stadia` VALUES (9, 'TBD', 'TBD', 9, 0, '', '2026-06-02 15:55:43.026', '2026-06-02 15:55:43.026', NULL);
INSERT INTO `stadia` VALUES (10, '西雅图体育场', 'Seattle Stadium', 5, 0, '', '2026-06-04 10:01:09.670', '2026-06-04 15:21:00.391', NULL);
INSERT INTO `stadia` VALUES (11, '瓜达拉哈拉体育场', 'Guadalajara Stadium', 10, 0, '', '2026-06-04 10:01:09.697', '2026-06-04 15:20:59.743', NULL);
INSERT INTO `stadia` VALUES (12, '堪萨斯城体育场', 'Kansas City Stadium', 11, 0, '', '2026-06-04 10:01:09.709', '2026-06-04 15:21:00.514', NULL);
INSERT INTO `stadia` VALUES (13, '亚特兰大体育场', 'Atlanta Stadium', 12, 0, '', '2026-06-04 10:01:09.720', '2026-06-04 15:21:00.558', NULL);
INSERT INTO `stadia` VALUES (14, '蒙特雷体育场', 'Monterrey Stadium', 13, 0, '', '2026-06-04 10:01:09.733', '2026-06-04 15:21:00.011', NULL);
INSERT INTO `stadia` VALUES (15, '旧金山湾区体育场', 'San Francisco Bay Area Stadium', 14, 0, '', '2026-06-04 10:01:09.743', '2026-06-04 15:21:00.134', NULL);
INSERT INTO `stadia` VALUES (16, '温哥华 BC Place', 'BC Place Vancouver', 8, 0, '', '2026-06-04 10:01:09.756', '2026-06-04 15:21:00.434', NULL);
INSERT INTO `stadia` VALUES (17, '纽约/新泽西体育场', 'New York/New Jersey Stadium', 15, 0, '', '2026-06-04 10:01:09.768', '2026-06-04 15:21:00.598', NULL);
INSERT INTO `stadia` VALUES (18, '迈阿密体育场', 'Miami Stadium', 6, 0, '', '2026-06-04 10:01:09.790', '2026-06-04 15:21:00.578', NULL);
INSERT INTO `stadia` VALUES (19, '达拉斯体育场', 'Dallas Stadium', 4, 0, '', '2026-06-04 10:01:09.814', '2026-06-04 15:21:00.538', NULL);
INSERT INTO `stadia` VALUES (20, '休斯敦体育场', 'Houston Stadium', 16, 0, '', '2026-06-04 10:01:09.827', '2026-06-04 15:21:00.281', NULL);
INSERT INTO `stadia` VALUES (21, '墨西哥城体育场', 'Mexico City Stadium', 1, 0, '', '2026-06-04 10:01:09.875', '2026-06-04 15:21:00.344', NULL);
INSERT INTO `stadia` VALUES (22, '波士顿体育场', 'Boston Stadium', 17, 0, '', '2026-06-04 10:01:09.899', '2026-06-04 15:21:00.454', NULL);
INSERT INTO `stadia` VALUES (23, '多伦多体育场', 'Toronto Stadium', 7, 0, '', '2026-06-04 10:01:09.923', '2026-06-04 15:21:00.176', NULL);
INSERT INTO `stadia` VALUES (24, '洛杉矶体育场', 'Los Angeles Stadium', 3, 0, '', '2026-06-04 10:01:09.996', '2026-06-04 15:21:00.474', NULL);
INSERT INTO `stadia` VALUES (25, '费城体育场', 'Philadelphia Stadium', 18, 0, '', '2026-06-04 10:01:10.091', '2026-06-04 15:21:00.303', NULL);

-- ----------------------------
-- Table structure for sync_states
-- ----------------------------
DROP TABLE IF EXISTS `sync_states`;
CREATE TABLE `sync_states`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `provider` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `resource` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'idle',
  `last_synced_at` datetime(3) NULL DEFAULT NULL,
  `next_sync_at` datetime(3) NULL DEFAULT NULL,
  `last_error` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sync_state`(`provider` ASC, `resource` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sync_states
-- ----------------------------
INSERT INTO `sync_states` VALUES (1, 'football-data', 'matches', 'success', '2026-06-04 15:21:00.616', '2026-06-04 15:51:00.616', '', '2026-06-02 15:55:39.291', '2026-06-04 15:21:00.617');

-- ----------------------------
-- Table structure for teams
-- ----------------------------
DROP TABLE IF EXISTS `teams`;
CREATE TABLE `teams`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `name_en` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `fifa_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `flag_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `continent` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `group_id` bigint UNSIGNED NULL DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `coach` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_teams_fifa_code`(`fifa_code` ASC) USING BTREE,
  INDEX `idx_teams_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `idx_teams_name`(`name` ASC) USING BTREE,
  INDEX `idx_teams_name_en`(`name_en` ASC) USING BTREE,
  INDEX `idx_teams_continent`(`continent` ASC) USING BTREE,
  INDEX `idx_teams_group_id`(`group_id` ASC) USING BTREE,
  CONSTRAINT `fk_teams_group` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 451 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of teams
-- ----------------------------
INSERT INTO `teams` VALUES (1, '墨西哥', 'Mexico', 'MEX', '🇲🇽', '北美洲', 1, '', '', '2026-06-02 11:53:03.382', '2026-06-04 15:20:59.516', NULL);
INSERT INTO `teams` VALUES (2, '加拿大', 'Canada', 'CAN', '🇨🇦', '北美洲', 2, '', '', '2026-06-02 11:53:03.388', '2026-06-04 15:20:59.430', NULL);
INSERT INTO `teams` VALUES (3, '南非', 'South Africa', 'RSA', '🇿🇦', '非洲', 1, '', '', '2026-06-02 11:53:03.393', '2026-06-04 15:20:59.533', NULL);
INSERT INTO `teams` VALUES (4, '新西兰', 'New Zealand', 'NZL', '🇳🇿', '大洋洲', 7, '', '', '2026-06-02 11:53:03.398', '2026-06-04 15:20:59.772', NULL);
INSERT INTO `teams` VALUES (5, '阿根廷', 'Argentina', 'ARG', '🇦🇷', '南美洲', 10, '', '', '2026-06-02 11:53:03.403', '2026-06-04 15:20:59.902', NULL);
INSERT INTO `teams` VALUES (6, '法国', 'France', 'FRA', '🇫🇷', '欧洲', 9, '', '', '2026-06-02 11:53:03.409', '2026-06-04 15:20:59.693', NULL);
INSERT INTO `teams` VALUES (7, '美国', 'USA', 'USA', '🇺🇸', '北美洲', 4, '', '', '2026-06-02 11:53:03.413', '2026-06-04 15:20:59.646', NULL);
INSERT INTO `teams` VALUES (8, '加纳', 'Ghana', 'GHA', '🇬🇭', '非洲', 12, '', '', '2026-06-02 11:53:03.419', '2026-06-04 15:20:59.840', NULL);
INSERT INTO `teams` VALUES (9, '巴西', 'Brazil', 'BRA', '🇧🇷', '南美洲', 3, '', '', '2026-06-02 11:53:03.424', '2026-06-04 15:20:59.494', NULL);
INSERT INTO `teams` VALUES (10, '西班牙', 'Spain', 'ESP', '🇪🇸', '欧洲', 8, '', '', '2026-06-02 11:53:03.429', '2026-06-04 15:20:59.735', NULL);
INSERT INTO `teams` VALUES (11, '英格兰', 'England', 'ENG', '🏴󠁧󠁢󠁥󠁮󠁧󠁿', '欧洲', 12, '', '', '2026-06-02 11:53:03.433', '2026-06-04 15:20:59.818', NULL);
INSERT INTO `teams` VALUES (12, '葡萄牙', 'Portugal', 'POR', '🇵🇹', '欧洲', 11, '', '', '2026-06-02 11:53:03.439', '2026-06-04 15:20:59.860', NULL);
INSERT INTO `teams` VALUES (13, '日本', 'Japan', 'JPN', '🇯🇵', '亚洲', 6, '', '', '2026-06-02 11:53:03.443', '2026-06-04 15:20:59.619', NULL);
INSERT INTO `teams` VALUES (14, '德国', 'Germany', 'GER', '🇩🇪', '欧洲', 5, '', '', '2026-06-02 11:53:03.449', '2026-06-04 15:20:59.556', NULL);
INSERT INTO `teams` VALUES (15, '荷兰', 'Netherlands', 'NED', '🇳🇱', '欧洲', 6, '', '', '2026-06-02 11:53:03.454', '2026-06-04 15:20:59.600', NULL);
INSERT INTO `teams` VALUES (16, '比利时', 'Belgium', 'BEL', '🇧🇪', '欧洲', 7, '', '', '2026-06-02 11:53:03.461', '2026-06-04 15:20:59.776', NULL);
INSERT INTO `teams` VALUES (417, '韩国', 'Korea Republic', 'KOR', 'https://crests.football-data.org/772.png', '亚洲', 1, '', '', '2026-06-02 16:17:25.404', '2026-06-04 15:20:59.536', NULL);
INSERT INTO `teams` VALUES (418, '捷克', 'Czechia', 'CZE', 'https://crests.football-data.org/798.svg', '欧洲', 1, '', '', '2026-06-02 16:17:25.408', '2026-06-04 15:20:59.511', NULL);
INSERT INTO `teams` VALUES (419, '波黑', 'Bosnia-H.', 'BIH', 'https://crests.football-data.org/bosnia.svg', '欧洲', 2, '', '', '2026-06-02 16:17:25.422', '2026-06-04 15:20:59.448', NULL);
INSERT INTO `teams` VALUES (420, '巴拉圭', 'Paraguay', 'PAR', 'https://crests.football-data.org/761.svg', '南美洲', 4, '', '', '2026-06-02 16:17:25.435', '2026-06-04 15:20:59.666', NULL);
INSERT INTO `teams` VALUES (421, '卡塔尔', 'Qatar', 'QAT', 'https://crests.football-data.org/8030.svg', '亚洲', 2, '', '', '2026-06-02 16:17:25.443', '2026-06-04 15:20:59.452', NULL);
INSERT INTO `teams` VALUES (422, '瑞士', 'Switzerland', 'SUI', 'https://crests.football-data.org/788.svg', '欧洲', 2, '', '', '2026-06-02 16:17:25.447', '2026-06-04 15:20:59.426', NULL);
INSERT INTO `teams` VALUES (423, '摩洛哥', 'Morocco', 'MAR', 'https://crests.football-data.org/morocco.svg', '非洲', 3, '', '', '2026-06-02 16:17:25.462', '2026-06-04 15:20:59.470', NULL);
INSERT INTO `teams` VALUES (424, '海地', 'Haiti', 'HAI', 'https://crests.football-data.org/haiti.svg', '北美洲', 3, '', '', '2026-06-02 16:17:25.471', '2026-06-04 15:20:59.473', NULL);
INSERT INTO `teams` VALUES (425, '苏格兰', 'Scotland', 'SCO', 'https://crests.football-data.org/814.svg', '欧洲', 3, '', '', '2026-06-02 16:17:25.474', '2026-06-04 15:20:59.490', NULL);
INSERT INTO `teams` VALUES (426, '澳大利亚', 'Australia', 'AUS', 'https://crests.football-data.org/779.svg', '大洋洲', 4, '', '', '2026-06-02 16:17:25.485', '2026-06-04 15:20:59.670', NULL);
INSERT INTO `teams` VALUES (427, '土耳其', 'Turkey', 'TUR', 'https://crests.football-data.org/803.svg', '欧洲', 4, '', '', '2026-06-02 16:17:25.488', '2026-06-04 15:20:59.643', NULL);
INSERT INTO `teams` VALUES (428, '库拉索', 'Curaçao', 'CUW', 'https://crests.football-data.org/curacao.svg', '北美洲', 5, '', '', '2026-06-02 16:17:25.502', '2026-06-04 15:20:59.573', NULL);
INSERT INTO `teams` VALUES (429, '科特迪瓦', 'Ivory Coast', 'CIV', 'https://crests.football-data.org/787.svg', '非洲', 5, '', '', '2026-06-02 16:17:25.525', '2026-06-04 15:20:59.576', NULL);
INSERT INTO `teams` VALUES (430, '厄瓜多尔', 'Ecuador', 'ECU', 'https://crests.football-data.org/791.svg', '南美洲', 5, '', '', '2026-06-02 16:17:25.528', '2026-06-04 15:20:59.553', NULL);
INSERT INTO `teams` VALUES (431, '瑞典', 'Sweden', 'SWE', 'https://crests.football-data.org/792.svg', '欧洲', 6, '', '', '2026-06-02 16:17:25.538', '2026-06-04 15:20:59.623', NULL);
INSERT INTO `teams` VALUES (432, '突尼斯', 'Tunisia', 'TUN', 'https://crests.football-data.org/tunisia.svg', '非洲', 6, '', '', '2026-06-02 16:17:25.544', '2026-06-04 15:20:59.595', NULL);
INSERT INTO `teams` VALUES (433, '佛得角', 'Cape Verde', 'CPV', 'https://crests.football-data.org/cape_verde.svg', '非洲', 8, '', '', '2026-06-02 16:17:25.559', '2026-06-04 15:20:59.751', NULL);
INSERT INTO `teams` VALUES (434, '埃及', 'Egypt', 'EGY', 'https://crests.football-data.org/825.svg', '非洲', 7, '', '', '2026-06-02 16:17:25.571', '2026-06-04 15:20:59.792', NULL);
INSERT INTO `teams` VALUES (435, '沙特阿拉伯', 'Saudi Arabia', 'KSA', 'https://crests.football-data.org/saudi_arabia.svg', '亚洲', 8, '', '', '2026-06-02 16:17:25.581', '2026-06-04 15:20:59.755', NULL);
INSERT INTO `teams` VALUES (436, '乌拉圭', 'Uruguay', 'URY', 'https://crests.football-data.org/758.svg', '', 8, '', '', '2026-06-02 16:17:25.584', '2026-06-04 15:20:59.731', NULL);
INSERT INTO `teams` VALUES (437, '伊朗', 'Iran', 'IRN', 'https://crests.football-data.org/iran.svg', '亚洲', 7, '', '', '2026-06-02 16:17:25.595', '2026-06-04 15:20:59.797', NULL);
INSERT INTO `teams` VALUES (438, '塞内加尔', 'Senegal', 'SEN', 'https://crests.football-data.org/senegal.svg', '非洲', 9, '', '', '2026-06-02 16:17:25.613', '2026-06-04 15:20:59.710', NULL);
INSERT INTO `teams` VALUES (439, '伊拉克', 'Iraq', 'IRQ', 'https://crests.football-data.org/iraq.svg', '亚洲', 9, '', '', '2026-06-02 16:17:25.624', '2026-06-04 15:20:59.713', NULL);
INSERT INTO `teams` VALUES (440, '挪威', 'Norway', 'NOR', 'https://crests.football-data.org/813.svg', '欧洲', 9, '', '', '2026-06-02 16:17:25.627', '2026-06-04 15:20:59.689', NULL);
INSERT INTO `teams` VALUES (441, '阿尔及利亚', 'Algeria', 'ALG', 'https://crests.football-data.org/algeria.svg', '非洲', 10, '', '', '2026-06-02 16:17:25.641', '2026-06-04 15:20:59.919', NULL);
INSERT INTO `teams` VALUES (442, '奥地利', 'Austria', 'AUT', 'https://crests.football-data.org/816.svg', '欧洲', 10, '', '', '2026-06-02 16:17:25.652', '2026-06-04 15:20:59.922', NULL);
INSERT INTO `teams` VALUES (443, '约旦', 'Jordan', 'JOR', 'https://crests.football-data.org/8049.png', '亚洲', 10, '', '', '2026-06-02 16:17:25.656', '2026-06-04 15:20:59.898', NULL);
INSERT INTO `teams` VALUES (444, '刚果民主共和国', 'Congo DR', 'COD', 'https://crests.football-data.org/congo_dr.svg', '非洲', 11, '', '', '2026-06-02 16:17:25.669', '2026-06-04 15:20:59.876', NULL);
INSERT INTO `teams` VALUES (445, '克罗地亚', 'Croatia', 'CRO', 'https://crests.football-data.org/799.svg', '欧洲', 12, '', '', '2026-06-02 16:17:25.682', '2026-06-04 15:20:59.836', NULL);
INSERT INTO `teams` VALUES (446, '巴拿马', 'Panama', 'PAN', 'https://crests.football-data.org/panama.svg', '北美洲', 12, '', '', '2026-06-02 16:17:25.695', '2026-06-04 15:20:59.814', NULL);
INSERT INTO `teams` VALUES (447, '乌兹别克斯坦', 'Uzbekistan', 'UZB', 'https://crests.football-data.org/8070.png', '亚洲', 11, '', '', '2026-06-02 16:17:25.704', '2026-06-04 15:20:59.881', NULL);
INSERT INTO `teams` VALUES (448, '哥伦比亚', 'Colombia', 'COL', 'https://crests.football-data.org/818.svg', '南美洲', 11, '', '', '2026-06-02 16:17:25.708', '2026-06-04 15:20:59.856', NULL);
INSERT INTO `teams` VALUES (449, 'TBD', 'TBD', 'TBD', '', '', 13, '', '', '2026-06-02 16:17:26.352', '2026-06-04 15:21:00.590', NULL);
INSERT INTO `teams` VALUES (450, '库拉索', 'Curaçao', 'CUR', 'https://crests.football-data.org/curacao.svg', '北美洲', 5, '', '', '2026-06-02 16:47:31.875', '2026-06-04 13:22:24.772', NULL);

-- ----------------------------
-- Table structure for user_favorite_matches
-- ----------------------------
DROP TABLE IF EXISTS `user_favorite_matches`;
CREATE TABLE `user_favorite_matches`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `match_id` bigint UNSIGNED NOT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_user_match`(`user_id` ASC, `match_id` ASC) USING BTREE,
  INDEX `fk_user_favorite_matches_match`(`match_id` ASC) USING BTREE,
  CONSTRAINT `fk_user_favorite_matches_match` FOREIGN KEY (`match_id`) REFERENCES `matches` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_favorite_matches
-- ----------------------------

-- ----------------------------
-- Table structure for user_favorite_teams
-- ----------------------------
DROP TABLE IF EXISTS `user_favorite_teams`;
CREATE TABLE `user_favorite_teams`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL,
  `team_id` bigint UNSIGNED NOT NULL,
  `created_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_user_team`(`user_id` ASC, `team_id` ASC) USING BTREE,
  INDEX `fk_user_favorite_teams_team`(`team_id` ASC) USING BTREE,
  CONSTRAINT `fk_user_favorite_teams_team` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_favorite_teams
-- ----------------------------
INSERT INTO `user_favorite_teams` VALUES (2, 2, 6, '2026-06-03 17:22:11.818');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `timezone` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'Asia/Shanghai',
  `language` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'zh-CN',
  `role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT 'user',
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  `notification_email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_users_username`(`username` ASC) USING BTREE,
  UNIQUE INDEX `idx_users_email`(`email` ASC) USING BTREE,
  INDEX `idx_users_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'admin', 'admin@worldcup.local', '$2a$10$fxj9oFDca1SYsIJHy4gEtuKbL4j9P8CtNzvO2vw6.WoVt9Efi26/m', '', 'Asia/Shanghai', 'zh-CN', 'admin', '2026-06-02 11:53:03.246', '2026-06-02 11:53:03.246', NULL, NULL);
INSERT INTO `users` VALUES (2, 'cxy', 'cxy@qq.com', '$2a$10$UBS5bRdEG1hYf9q2DFFyBuj.pWDuaDt.60SQIF.je8f2SWKOxFKB2', '/uploads/avatars/2_1780471765242.png', 'Asia/Shanghai', 'zh-CN', 'user', '2026-06-02 17:26:27.048', '2026-06-04 08:33:47.321', NULL, '1270426208@qq.com');

SET FOREIGN_KEY_CHECKS = 1;
