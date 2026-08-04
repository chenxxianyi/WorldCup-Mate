# WorldCup Mate 基线检查清单（BASELINE）

> 来源：《项目全面优化实施方案.md》BASE-01
> 目的：记录改造前的可运行状态，供后续每个里程碑对比回归。
> 建立日期：2026-08-04

## 1. 环境版本

| 组件 | 版本 | 备注 |
| --- | --- | --- |
| Node.js | v22.22.3 | 前端构建 |
| Go | go1.24.11 (windows/386) | 后端 |
| MySQL | 8.0.32 | 本机运行于 127.0.0.1:3310 |
| Redis | 7（目标版本） | 本机运行于 127.0.0.1:5001 |
| 构建工具 | npm（可用）、git 2.31.1 | |

> 注：本机未安装 `docker`、`make`、`rg`；`python3` 不可用但 `python` 可用。

## 2. 构建与测试基线（2026-08-04 实测）

| 检查项 | 命令 | 结果 |
| --- | --- | --- |
| 前端生产构建 | `cd frontend && npm run build` | ✅ 通过（✓ built，gzip 主包约 61-158 KB） |
| 后端编译 | `cd backend && go build ./...` | ✅ 通过 |
| 后端静态检查 | `cd backend && go vet ./...` | ✅ 通过 |
| 后端单元测试 | `cd backend && go test ./...` | ✅ 通过（internal/services、internal/utils） |

## 3. 数据库基线

- 结构基线文件：`docs/db-baseline.sql`（12 张表，来自 `worldcup_mate.sql` Navicat dump，MySQL 8.0.32）
- 多赛事扩展后新增表（AutoMigrate 生成）：`competitions`、`league_standings`
- `matches`/`teams` 表新增可空列：`competition_id`、`season`、`matchday`（matches）；`external_code`、`team_type`、`country`、`venue`（teams）

### 3.1 六项赛事数据量（待补录）

> 本机无 `mysql` 客户端，以下统计需在具备客户端或应用环境时补录并归档。

| 赛事 | 比赛数 | 球队数 | 积分榜行数 |
| --- | --- | --- | --- |
| WC 世界杯 | 待补 | 待补（dump 显示约 50 队） | 待补 |
| PL 英超 | 待补 | 待补 | 待补 |
| PD 西甲 | 待补 | 待补 | 待补 |
| BL1 德甲 | 待补 | 待补 | 待补 |
| SA 意甲 | 待补 | 待补 | 待补 |
| FL1 法甲 | 待补 | 待补 | 待补 |

补录 SQL 示例：

```sql
SELECT c.code, c.name, COUNT(DISTINCT m.id) AS matches
FROM competitions c
LEFT JOIN matches m ON m.competition_id = c.id
GROUP BY c.id ORDER BY c.sort_order;
```

## 4. 接口清单（2026-08-04 从 router.go 提取）

### 4.1 公开接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/auth/register | 注册 |
| POST | /api/auth/login | 登录 |
| GET | /api/matches | 比赛列表（可选 competitionId/season/matchday） |
| GET | /api/matches/today | 今日比赛（跨赛事） |
| GET | /api/matches/tomorrow | 明日比赛 |
| GET | /api/matches/upcoming | 未开始 |
| GET | /api/matches/live | 直播中 |
| GET | /api/matches/recommended | 热门推荐 |
| GET | /api/matches/progress | 赛事进度 |
| GET | /api/matches/:id | 比赛详情 |
| GET | /api/matches/by-team/:teamId | 球队比赛 |
| GET | /api/matches/by-group/:groupId | 小组比赛 |
| GET | /api/matches/by-stage/:stage | 阶段比赛 |
| GET | /api/teams | 球队列表（可选 teamType/country） |
| GET | /api/teams/:id | 球队详情 |
| GET | /api/teams/:id/matches | 球队比赛 |
| GET | /api/groups | 小组列表 |
| GET | /api/groups/:id | 小组详情 |
| GET | /api/groups/:id/standings | 小组积分榜 |
| GET | /api/standings | 全部积分榜 |
| GET | /api/standings/best-third | 最佳第三名 |
| GET | /api/cities | 城市列表 |
| GET | /api/stadiums | 球场列表 |
| GET | /api/stadiums/:id | 球场详情 |
| GET | /api/sync/status | 同步状态 |
| GET | /api/competitions | 赛事列表 |
| GET | /api/competitions/:code/standings | 联赛积分榜（season/type 可选） |

### 4.2 用户接口（JWT）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET/PUT | /api/user/profile | 个人资料 |
| PUT | /api/user/password | 修改密码 |
| POST | /api/user/avatar | 头像上传 |
| POST/DELETE | /api/favorites/teams/:teamId | 关注/取消关注球队 |
| GET | /api/favorites/teams | 关注列表 |
| POST/DELETE | /api/favorites/matches/:matchId | 收藏/取消收藏比赛 |
| GET | /api/favorites/matches | 收藏列表 |
| POST | /api/reminders/batch | 批量创建提醒 |
| POST | /api/reminders | 创建提醒 |
| GET | /api/reminders | 提醒列表 |
| PUT | /api/reminders/:id | 更新提醒 |
| DELETE | /api/reminders/:id | 删除提醒 |
| GET | /api/notifications | 通知列表 |
| GET | /api/notifications/unread-count | 未读数 |
| PUT | /api/notifications/:id/read | 单条已读 |
| PUT | /api/notifications/read-all | 全部已读 |

### 4.3 管理接口（JWT + admin）

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | /api/admin/login | 管理员登录 |
| GET | /api/admin/dashboard | 后台看板 |
| GET/POST/PUT | /api/admin/competitions | 赛事管理 |
| GET | /api/admin/teams | 球队列表 |
| POST/PUT/DELETE | /api/admin/teams/... | 球队管理 |
| GET/POST | /api/admin/groups | 分组管理 |
| GET/POST | /api/admin/cities | 城市管理 |
| GET/POST | /api/admin/stadiums | 球场管理 |
| GET/POST/PUT/DELETE | /api/admin/matches | 比赛管理 |
| GET | /api/admin/standings | 积分榜 |
| POST | /api/admin/standings/recalculate | 小组积分重算 |
| POST | /api/admin/standings/league/recalculate | 联赛积分重算 |
| PUT | /api/admin/standings/:id | 积分修正 |
| GET/PUT | /api/admin/users | 用户管理 |
| POST | /api/admin/sync/matches | 手动同步（?code= 指定赛事） |

## 5. 手动冒烟测试清单

### 5.1 用户端

- [ ] 首页：倒计时、今日比赛、今日焦点、热门推荐、赛事进度、积分速览均正常
- [ ] 赛事切换：六项赛事切换后各模块数据跟随变化
- [ ] 赛程页：世界杯小组/淘汰赛筛选；联赛轮次筛选
- [ ] 比赛详情：比分、状态、提醒设置
- [ ] 球队页/详情：搜索、关注、俱乐部视图
- [ ] 积分榜：世界杯小组榜 + 最佳第三名；联赛单表 + 分区标记
- [ ] 登录/注册/退出
- [ ] 我的：头像上传、关注分组、收藏、提醒、通知

### 5.2 管理端

- [ ] 管理员登录（admin@worldcup.local / admin123456）
- [ ] 后台看板、比赛/球队/积分榜/同步管理各页
- [ ] 联赛积分榜查看与重算
- [ ] 单赛事手动同步

### 5.3 数据与任务

- [ ] `GET /api/competitions` 返回 6 个赛事
- [ ] `/api/sync/status` 展示各赛事同步状态
- [ ] 联赛同步（需 FOOTBALL_DATA_API_KEY）后比赛/球队/积分榜落库

## 6. 已知基线问题（后续任务修复）

| 问题 | 对应任务 |
| --- | --- |
| 通知单条已读无用户归属校验 | SEC-01 |
| 管理端部分接口为占位实现 | ADM-02～ADM-06 |
| 所有接口 HTTP 状态码固定 200 | API-01 |
| 无 Request ID | API-02 |
| 前端 API 层大量 any | API-04 |
| 登录注册无限流、密码策略弱 | SEC-02、SEC-03 |
| Token 无刷新/撤销机制 | SEC-04 |
| 上传文件与构建产物被 Git 跟踪 | SEC-06 |
| 无健康检查、优雅关闭 | OBS-02、OBS-03 |
