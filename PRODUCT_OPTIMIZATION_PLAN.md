# WorldCup Mate 产品优化方案

## 1. 目标

WorldCup Mate 当前已经具备世界杯赛程、球队、积分榜、收藏、提醒、个人中心、后台和数据同步能力。下一阶段的核心目标不是继续堆叠页面，而是把产品从“赛事信息展示”升级为“用户个人化世界杯观赛助手”。

核心闭环：

1. 用户选择关注球队。
2. 首页和赛程自动突出用户关心的比赛。
3. 用户可以方便设置提醒。
4. 比赛临近时通过站内通知或邮件提醒用户。
5. 用户回到产品查看赛程、详情、积分和后续比赛。

## 2. 优先级总览

| 优先级 | 方向 | 目标 |
| --- | --- | --- |
| P0 | 关注球队个性化首页 | 让关注行为产生明确价值，提高复访 |
| P0 | 球队详情页 | 补齐球队列表到球队信息的自然跳转 |
| P0 | 提醒能力升级 | 从单点提醒升级为可配置、可感知的提醒系统 |
| P0 | 数据可信度提示 | 降低用户对赛程和积分数据准确性的疑虑 |
| P1 | 后台管理补全 | 提升数据维护效率，减少手动修库 |
| P1 | 性能优化 | 保证移动端访问流畅 |
| P1 | 生产安全检查 | 避免默认账号、配置泄露和上传风险 |

## 3. P0 功能优化方案

### 3.1 关注球队个性化首页

#### 问题

当前用户可以关注球队，但关注后的价值感还不够明显。关注球队后，首页和赛程页没有足够突出“我关心的比赛”。

#### 建议功能

- 首页新增“我的关注赛程”模块。
- 展示关注球队的下一场比赛。
- 展示今日关注球队比赛。
- 赛程页新增“只看关注”筛选。
- 关注球队后，比赛卡片增加轻量标识。
- 未登录用户点击关注时，引导登录。

#### 前端改动

- `frontend/src/pages/user/HomePage.vue`
  - 增加关注球队下一场比赛模块。
  - 增加关注球队今日比赛模块。
- `frontend/src/pages/user/SchedulePage.vue`
  - 增加“只看关注”筛选项。
  - 搜索、日期、状态筛选需要和“只看关注”组合生效。
- `frontend/src/stores/useFavoriteStore.ts`
  - 确保关注球队和收藏比赛状态在刷新后能稳定恢复。

#### 后端改动

- 可新增接口：
  - `GET /api/user/followed-team-matches`
  - 返回当前用户关注球队相关比赛，支持 `status`、`date`、`limit` 参数。
- 或复用 `GET /api/matches`，增加 `followed_only=true` 参数。

#### 验收标准

- 已登录用户关注球队后，首页能看到关注球队的下一场比赛。
- 赛程页开启“只看关注”后，只展示关注球队相关比赛。
- 无关注球队时展示清晰空状态。
- 未登录用户点击关注或只看关注时，有登录引导。

### 3.2 球队详情页

#### 问题

球队页当前只有球队卡片列表。用户点击某支球队后，缺少球队详情、赛程和积分入口。

#### 建议功能

- 新增 `/teams/:id` 球队详情页。
- 展示球队名称、英文名、队徽、大洲、小组、关注状态。
- 展示该球队完整比赛列表。
- 展示该球队所在小组积分榜。
- 展示下一场比赛快捷入口。

#### 前端改动

- 新增页面：
  - `frontend/src/pages/user/TeamDetailPage.vue`
- 修改路由：
  - `frontend/src/router/index.ts`
  - 增加 `{ path: 'teams/:id', name: 'team-detail', component: ... }`
- 修改卡片：
  - `frontend/src/components/common/TeamCard.vue`
  - 点击卡片进入详情页，收藏按钮继续阻止冒泡。

#### 后端改动

- 当前已有：
  - `GET /api/teams/:id`
- 建议确认或补齐：
  - `GET /api/teams/:id/matches`
  - `GET /api/groups/:id/standings`

#### 验收标准

- 点击球队卡片可以进入球队详情。
- 球队详情页能展示该球队所有比赛。
- 球队详情页能展示所在小组积分。
- 关注/取消关注后，列表页和详情页状态一致。

### 3.3 提醒能力升级

#### 问题

当前比赛提醒已经存在，但提醒方式和反馈还比较基础。用户设置后不容易感知提醒是否真正生效。

#### 建议功能

- 支持多档提醒：
  - 赛前 1 天
  - 赛前 1 小时
  - 赛前 15 分钟
- 支持关注球队自动提醒开关。
- 站内通知增加未读红点。
- 邮箱未设置时，邮件提醒入口提示先配置邮箱。
- 提醒成功后给出明确反馈。

#### 前端改动

- `frontend/src/components/common/MatchCard.vue`
  - 提醒按钮点击后弹出提醒时间和渠道选择。
- `frontend/src/pages/user/MatchDetailPage.vue`
  - 补齐和赛程卡片一致的提醒配置能力。
- `frontend/src/pages/user/ProfilePage.vue`
  - 增加“关注球队自动提醒”开关。
  - 增加邮箱配置状态提示。
- `frontend/src/components/common/TopBar.vue`
  - 可增加通知未读红点。

#### 后端改动

- `backend/internal/models/reminder.go`
  - 确认是否支持同一用户、同一比赛、多提醒时间。
  - 如果不支持，需要调整唯一约束或提醒模型。
- `backend/internal/services/reminder_service.go`
  - 支持批量创建提醒。
  - 创建提醒时校验比赛是否已经开赛。
- 新增或扩展接口：
  - `POST /api/reminders/batch`
  - `GET /api/notifications/unread-count`

#### 验收标准

- 用户可以为同一场比赛选择不同提醒时间。
- 未配置邮箱时不能静默失败。
- 站内通知未读数量能正确展示。
- 比赛已经开始或结束时，不能创建无效提醒。

### 3.4 数据可信度提示

#### 问题

赛事产品对数据准确性要求高。当前用户无法直观看到数据来源、同步状态和更新时间。

#### 建议功能

- 赛程页显示最近同步时间。
- 积分榜页显示最近更新时间。
- 后台显示同步状态、同步结果和失败原因。
- 如果当前是演示数据，需要明确标记。

#### 前端改动

- `frontend/src/pages/user/SchedulePage.vue`
  - 顶部增加“更新时间”。
- `frontend/src/pages/user/StandingsPage.vue`
  - 增加积分更新时间。
- `frontend/src/pages/admin/AdminDashboard.vue`
  - 增加同步状态卡片。

#### 后端改动

- 当前已有：
  - `GET /api/sync/status`
- 建议扩展返回：
  - `last_success_at`
  - `last_failed_at`
  - `last_error`
  - `provider`
  - `is_demo_data`

#### 验收标准

- 用户能看到赛程数据最后更新时间。
- 后台能看到最近一次同步结果。
- 同步失败时后台能看到失败原因。

## 4. P1 技术产品化优化

### 4.1 后台管理补全

#### 问题

后台目前更接近数据看板，侧边栏中的球队管理、比赛管理、积分榜管理还没有完整独立页面。

#### 建议功能

- 球队管理：列表、搜索、编辑名称、队徽、大洲、小组。
- 比赛管理：列表、筛选、编辑比分、状态、时间、球场。
- 积分榜管理：重算积分、查看计算结果。
- 同步管理：手动同步、同步日志、失败重试。

#### 验收标准

- 管理员不需要直接操作数据库即可修正核心赛事数据。
- 同步失败后可以在后台看到失败原因并重试。

### 4.2 性能优化

#### 问题

前端生产构建存在大 chunk 提示。随着页面和功能增加，移动端加载会变慢。

#### 建议优化

- 路由级代码分割已经存在，继续检查大依赖是否被首页直接引入。
- 图片资源懒加载。
- 赛程长列表继续使用分页加载。
- 后续如果比赛或通知列表变长，可引入虚拟滚动。
- 对接口增加合理分页，避免一次性返回过多数据。

#### 验收标准

- 首屏加载稳定。
- 移动端滚动不卡顿。
- 赛程页和球队页不会一次性渲染过多节点。

### 4.3 生产安全检查

#### 问题

项目中存在默认管理员账号、JWT、SMTP、上传、CORS 等生产部署敏感点。

#### 建议优化

- 登录页生产环境不展示默认测试账号。
- 生产启动时强制检查 `JWT_SECRET` 不能是默认值。
- 管理员默认密码首次启动后强制修改。
- 上传头像限制文件类型和大小。
- CORS 根据环境收敛允许域名。
- `.env` 严禁提交仓库。
- 后端日志避免输出敏感配置。

#### 验收标准

- 生产环境无法使用默认弱密钥启动。
- 头像上传不能上传非图片文件。
- 登录页不会暴露测试账号。
- Git 仓库不包含 `.env` 和敏感密钥。

## 5. 建议实施排期

### 第一阶段：核心闭环增强

建议优先完成：

1. 球队详情页。
2. 关注球队个性化首页模块。
3. 赛程页“只看关注”筛选。
4. 提醒配置弹层优化。

目标：用户关注球队后，能在首页、赛程和提醒中获得连续体验。

### 第二阶段：数据可信与后台能力

建议完成：

1. 同步状态展示。
2. 数据更新时间展示。
3. 后台比赛管理。
4. 后台球队管理。
5. 同步日志和失败重试。

目标：减少数据维护成本，提高赛事数据可信度。

### 第三阶段：产品亮点与增长

建议完成：

1. 添加到日历。
2. 观赛清单。
3. 比赛分享卡片。
4. 预测玩法。
5. PWA 支持。

目标：让产品从工具变成有传播能力的世界杯助手。

## 6. 推荐下一步

建议下一步先开发“球队详情页 + 只看关注 + 首页关注赛程”这一组功能。

原因：

- 改动范围清晰。
- 和当前已有的球队、比赛、收藏能力高度复用。
- 用户价值明显。
- 不依赖复杂外部服务。
- 做完后产品闭环会立刻更完整。

## 7. Agent 开发任务拆分

本章节用于指导开发 Agent 按步骤逐步实现。每个任务尽量保持边界清晰，完成后都需要进行构建或测试验证。

### 7.1 开发顺序总览

建议按以下顺序推进：

1. 梳理现有接口和数据结构。
2. 开发球队详情页。
3. 开发赛程页“只看关注”。
4. 开发首页关注球队赛程模块。
5. 升级提醒配置弹层。
6. 增加通知未读数和红点。
7. 增加数据更新时间和同步状态展示。
8. 补全后台管理入口和基础页面。
9. 做性能优化。
10. 做生产安全检查。

### 7.2 任务 0：开发前基线检查（已完成）

#### 目标

在开发新功能前确认项目当前状态，避免把已有问题和新功能混在一起。

#### 建议操作

- 查看 Git 工作区改动。
- 确认前端构建是否通过。
- 确认后端测试是否通过。
- 梳理当前路由、API、store 和类型定义。

#### 重点文件

- `frontend/src/router/index.ts`
- `frontend/src/stores/useTeamStore.ts`
- `frontend/src/stores/useMatchStore.ts`
- `frontend/src/stores/useFavoriteStore.ts`
- `frontend/src/types/team.ts`
- `frontend/src/types/match.ts`
- `backend/internal/routes/routes.go`
- `backend/internal/handlers/team_handler.go`
- `backend/internal/repositories/team_repo.go`
- `backend/internal/repositories/match_repo.go`

#### 验收标准

- 明确当前可复用的接口。
- 明确哪些接口需要新增。
- 前端构建问题和后端测试问题有记录。

#### 完成记录

- 完成时间：2026-06-04
- 已确认可复用接口：
  - `GET /api/teams/:id`
  - `GET /api/teams/:id/matches`
  - `GET /api/groups/:id/standings`
- 已确认可复用前端模块：
  - `useTeamStore.fetchTeamDetail`
  - `apiGetTeamMatches`
  - `apiGetGroupStandings`
  - `MatchCard`
  - `StandingTable`
- 验证结果：
  - `npm run build` 通过。
  - `go test ./...` 通过。

### 7.3 任务 1：球队详情页（已完成）

#### 目标

用户点击球队卡片后，可以查看该球队详情、完整赛程和所在小组积分。

#### 前端任务

1. 新增页面文件：
   - `frontend/src/pages/user/TeamDetailPage.vue`
2. 修改路由：
   - `frontend/src/router/index.ts`
   - 增加 `/teams/:id`
3. 修改球队卡片：
   - `frontend/src/components/common/TeamCard.vue`
   - 卡片主体点击跳转详情页。
   - 收藏按钮点击时阻止冒泡。
4. 在详情页展示：
   - 队徽、中文名、英文名、FIFA code、大洲、小组。
   - 关注/取消关注按钮。
   - 下一场比赛。
   - 全部比赛列表。
   - 所在小组积分榜。

#### 后端任务

1. 确认是否已有：
   - `GET /api/teams/:id`
   - `GET /api/teams/:id/matches`
2. 如果缺少 `GET /api/teams/:id/matches`，新增接口。
3. 确保比赛返回包含：
   - 主队、客队、队徽、城市、球场、开球时间、状态、比分、小组。

#### 建议接口

```text
GET /api/teams/:id
GET /api/teams/:id/matches?page=1&page_size=20
GET /api/groups/:id/standings
```

#### 验收标准

- 从球队列表点击任意球队可进入详情页。
- 详情页刷新后数据仍能正常加载。
- 关注状态和列表页保持一致。
- 球队比赛为空时有空状态。
- `npm run build` 通过。
- `go test ./...` 通过。

#### 完成记录

- 完成时间：2026-06-04
- 新增页面：
  - `frontend/src/pages/user/TeamDetailPage.vue`
- 修改文件：
  - `frontend/src/router/index.ts`
  - `frontend/src/components/common/TeamCard.vue`
  - `frontend/src/types/team.ts`
- 实现内容：
  - 新增 `/teams/:id` 用户端路由。
  - 球队卡片点击进入球队详情页。
  - 收藏按钮阻止冒泡，不触发页面跳转。
  - 球队详情页展示队徽、名称、英文名、FIFA code、大洲、小组。
  - 球队详情页展示关注按钮、下一场比赛、球队完整赛程、所在小组积分。
  - 未登录点击关注会跳转登录页。
  - 球队不存在、加载失败、暂无赛程均有页面状态。
- 验证结果：
  - `npm run build` 通过。
  - `go test ./...` 通过。

### 7.4 任务 2：赛程页“只看关注”（已完成）

#### 目标

让用户可以在赛程页快速筛选自己关注球队相关比赛。

#### 前端任务

1. 修改：
   - `frontend/src/pages/user/SchedulePage.vue`
2. 在筛选项中增加：
   - `只看关注`
3. 筛选逻辑需要和现有条件组合：
   - 日期
   - 状态
   - 小组
   - 淘汰赛
   - 搜索关键词
4. 未登录用户点击“只看关注”时：
   - 引导登录。
   - 或展示登录提示空状态。
5. 无关注球队时：
   - 展示“还没有关注球队”的空状态。
   - 提供跳转球队页按钮。

#### 后端任务

优先复用前端已有数据过滤。如果比赛数据量变大，再新增后端过滤。

可选新增接口参数：

```text
GET /api/matches?followed_only=true
```

#### 验收标准

- 已关注球队后，赛程页只展示相关比赛。
- 搜索和只看关注可以同时生效。
- 切换筛选项后分页加载不会重复或错乱。
- 未登录和无关注场景有清晰提示。
- `npm run build` 通过。

#### 完成记录

- 完成时间：2026-06-04
- 修改文件：
  - `frontend/src/pages/user/SchedulePage.vue`
- 实现内容：
  - 新增“只看关注”筛选项。
  - 登录用户会先拉取关注球队 ID，再过滤主客队包含关注球队的比赛。
  - 未登录用户选择“只看关注”时展示登录引导。
  - 已登录但未关注球队时展示去关注球队的引导。
  - 保留最近三天默认展示、下滑提示、延迟加载、返回顶部能力。
  - 搜索关键词和“只看关注”可以组合生效。
  - 分页参数从 `pageSize` 修正为后端实际读取的 `page_size`。
- 验证结果：
  - `npm run build` 通过。

### 7.5 任务 3：首页关注球队赛程模块（已完成）

#### 目标

用户打开首页后，能第一时间看到自己关注球队的下一场比赛。

#### 前端任务

1. 修改：
   - `frontend/src/pages/user/HomePage.vue`
2. 新增模块：
   - `我的关注赛程`
3. 模块内容：
   - 最近一场关注球队比赛。
   - 今日关注球队比赛。
   - 跳转赛程页按钮。
4. 状态处理：
   - 未登录：提示登录后关注球队。
   - 已登录但无关注：提示去球队页关注。
   - 有关注但无比赛：展示暂无赛程。

#### 后端任务

可先复用：

```text
GET /api/matches/upcoming
GET /api/favorites/teams
```

如前端组合复杂，可新增：

```text
GET /api/user/followed-team-matches?limit=5
```

#### 验收标准

- 首页能展示关注球队下一场比赛。
- 点击比赛可以进入比赛详情。
- 点击球队可以进入球队详情。
- 未登录、无关注、无比赛都有对应状态。
- `npm run build` 通过。

#### 完成记录

- 完成时间：2026-06-04
- 修改文件：
  - `frontend/src/pages/user/HomePage.vue`
- 实现内容：
  - 首页新增“我的关注赛程”模块。
  - 登录用户会根据关注球队筛选未开始比赛。
  - 优先展示关注球队的下一场比赛。
  - 如果今天还有其他关注球队比赛，会展示轻量快捷入口。
  - 未登录时展示登录引导。
  - 已登录但无关注球队时展示去关注球队引导。
  - 有关注但暂无未开始比赛时展示查看全部赛程入口。
  - 首页侧栏“我的关注”球队卡片支持点击进入球队详情。
  - 保留下一场倒计时、赛事进度、今日比赛、热门推荐、积分速览能力。
- 验证结果：
  - `npm run build` 通过。

### 7.6 任务 4：提醒配置弹层升级（已完成）

#### 目标

让用户可以清楚选择提醒时间和提醒渠道，避免提醒配置不透明。

#### 前端任务

1. 修改：
   - `frontend/src/components/common/MatchCard.vue`
   - `frontend/src/pages/user/MatchDetailPage.vue`
2. 提醒弹层支持：
   - 赛前 1 天
   - 赛前 1 小时
   - 赛前 15 分钟
   - 站内通知
   - 邮件通知
3. 邮件提醒逻辑：
   - 如果用户未设置通知邮箱，提示去个人中心设置。
4. 已设置提醒后：
   - 展示已提醒状态。
   - 支持取消提醒。

#### 后端任务

1. 检查 `Reminder` 模型是否支持同一比赛多个提醒时间。
2. 如不支持，需要调整唯一约束。
3. 增加批量创建提醒接口：

```text
POST /api/reminders/batch
```

请求建议：

```json
{
  "match_id": 1,
  "minutes": [1440, 60, 15],
  "channel": "site"
}
```

4. 创建提醒时校验：
   - 比赛不存在。
   - 比赛已开始。
   - 提醒时间已经过去。
   - 邮件提醒但用户没有通知邮箱。

#### 验收标准

- 可以选择不同提醒时间。
- 无效提醒不会创建。
- 邮件提醒缺少邮箱时有明确提示。
- 取消提醒后状态同步更新。
- `go test ./...` 通过。
- `npm run build` 通过。

#### 完成记录

- 完成时间：2026-06-04
- 新增文件：
  - `frontend/src/components/common/ReminderControl.vue`
- 修改文件：
  - `frontend/src/api/reminders.ts`
  - `frontend/src/stores/useReminderStore.ts`
  - `frontend/src/components/common/MatchCard.vue`
  - `frontend/src/pages/user/MatchDetailPage.vue`
  - `backend/internal/services/reminder_service.go`
  - `backend/internal/handlers/reminder_handler.go`
  - `backend/internal/routes/router.go`
- 实现内容：
  - 新增 `POST /api/reminders/batch` 批量创建提醒接口。
  - 支持赛前 1 天、赛前 1 小时、赛前 15 分钟多选。
  - 支持站内通知和邮件通知渠道。
  - 后端拒绝已开始、已结束、已取消比赛的提醒创建。
  - 后端拒绝已经过期的提醒时间。
  - 后端会跳过同一用户、同一比赛、同一时间、同一渠道的重复提醒。
  - 前端新增通用 `ReminderControl`，赛程卡片和比赛详情页共用同一套提醒弹层。
  - 已设置提醒时再次点击会取消该比赛下的全部提醒。
  - 邮件提醒缺少邮箱目标时，前端会提示先设置通知邮箱。
- 验证结果：
  - `go test ./...` 通过。
  - `npm run build` 通过。

### 7.7 任务 5：通知未读数和红点（已完成）

#### 目标

让用户知道自己有未读提醒或系统通知。

#### 前端任务

1. 修改：
   - `frontend/src/components/common/TopBar.vue`
   - `frontend/src/pages/user/ProfilePage.vue`
2. 增加未读红点。
3. 个人中心增加通知列表入口。
4. 支持标记单条已读和全部已读。

#### 后端任务

1. 增加接口：

```text
GET /api/notifications/unread-count
```

2. 确认已有接口：

```text
GET /api/notifications
PUT /api/notifications/:id/read
PUT /api/notifications/read-all
```

#### 验收标准

- 有未读通知时顶部显示红点。
- 阅读后红点消失。
- 全部已读后未读数为 0。
- 未登录用户不请求通知接口。

#### 完成记录

- 完成时间：2026-06-04
- 新增文件：
  - `frontend/src/api/notifications.ts`
  - `frontend/src/stores/useNotificationStore.ts`
  - `frontend/src/components/common/NotificationList.vue`
- 修改文件：
  - `frontend/src/api/index.ts`
  - `frontend/src/components/common/TopBar.vue`
  - `frontend/src/pages/user/ProfilePage.vue`
  - `backend/internal/repositories/notification_repo.go`
  - `backend/internal/services/notification_service.go`
  - `backend/internal/handlers/notification_handler.go`
  - `backend/internal/routes/router.go`
- 实现内容：
  - 新增 `GET /api/notifications/unread-count`。
  - 顶部栏登录状态下自动拉取未读通知数量。
  - 个人中心入口显示未读红点。
  - 个人中心新增通知列表。
  - 支持点击单条通知标记已读。
  - 支持全部标记已读。
  - 未登录状态不请求通知接口。
- 验证结果：
  - `go test ./...` 通过。
  - `npm run build` 通过。

### 7.8 任务 6：数据更新时间和同步状态展示（已完成）

#### 目标

增强赛事数据可信度，让用户和管理员知道数据是否新鲜。

#### 前端任务

1. 修改：
   - `frontend/src/pages/user/SchedulePage.vue`
   - `frontend/src/pages/user/StandingsPage.vue`
   - `frontend/src/pages/admin/AdminDashboard.vue`
2. 赛程页展示：
   - 最近更新时间。
   - 数据来源。
3. 积分榜页展示：
   - 最近更新时间。
4. 后台展示：
   - 同步状态。
   - 最近成功时间。
   - 最近失败时间。
   - 最近失败原因。

#### 后端任务

1. 扩展：

```text
GET /api/sync/status
```

2. 返回字段建议：

```json
{
  "provider": "football-data",
  "enabled": true,
  "last_success_at": "2026-06-04T10:00:00Z",
  "last_failed_at": null,
  "last_error": "",
  "next_run_at": "2026-06-04T11:00:00Z",
  "is_demo_data": false
}
```

#### 验收标准

- 赛程页和积分榜页能看到更新时间。
- 后台能看到同步失败原因。
- 同步状态接口失败时，前端不影响主内容展示。

#### 完成记录

- 完成时间：2026-06-04
- 新增文件：
  - `frontend/src/api/sync.ts`
  - `frontend/src/components/common/SyncStatusBadge.vue`
- 修改文件：
  - `frontend/src/api/index.ts`
  - `frontend/src/pages/user/SchedulePage.vue`
  - `frontend/src/pages/user/StandingsPage.vue`
  - `frontend/src/pages/admin/AdminDashboard.vue`
- 实现内容：
  - 复用现有 `GET /api/sync/status`。
  - 赛程页顶部展示数据同步状态和最近更新时间。
  - 积分榜页顶部展示数据同步状态和最近更新时间。
  - 后台看板展示同步状态卡片。
  - 后台同步失败时展示最近失败原因。
  - 同步状态接口失败时显示“演示数据或未同步”，不影响主内容渲染。
  - 后台看板进入页面时加载比赛列表。
- 验证结果：
  - `npm run build` 通过。

### 7.9 任务 7：后台管理基础页面（已完成）

#### 目标

让管理员可以通过后台维护基础赛事数据，而不是直接操作数据库。

#### 前端任务

1. 修改：
   - `frontend/src/router/index.ts`
   - `frontend/src/layouts/AdminLayout.vue`
2. 新增页面：
   - `frontend/src/pages/admin/AdminTeamsPage.vue`
   - `frontend/src/pages/admin/AdminMatchesPage.vue`
   - `frontend/src/pages/admin/AdminStandingsPage.vue`
   - `frontend/src/pages/admin/AdminSyncPage.vue`
3. 后台侧边栏跳转到真实页面。
4. 每个页面先完成列表、搜索、分页。
5. 再逐步补充编辑能力。

#### 后端任务

确认或补齐：

```text
GET /api/admin/teams
PUT /api/admin/teams/:id
GET /api/admin/matches
PUT /api/admin/matches/:id
PUT /api/admin/matches/:id/score
PUT /api/admin/matches/:id/status
POST /api/admin/standings/recalculate
POST /api/admin/sync/matches
GET /api/admin/sync/logs
```

#### 验收标准

- 后台侧边栏每个菜单能进入独立页面。
- 球队和比赛支持搜索分页。
- 管理员可以修改比分和比赛状态。
- 修改比赛结果后积分榜可重算。

#### 完成记录

- 完成时间：2026-06-04
- 新增文件：
  - `frontend/src/api/admin.ts`
  - `frontend/src/pages/admin/AdminTeamsPage.vue`
  - `frontend/src/pages/admin/AdminMatchesPage.vue`
  - `frontend/src/pages/admin/AdminStandingsPage.vue`
  - `frontend/src/pages/admin/AdminSyncPage.vue`
- 修改文件：
  - `frontend/src/api/index.ts`
  - `frontend/src/router/index.ts`
  - `frontend/src/layouts/AdminLayout.vue`
  - `frontend/src/pages/admin/AdminDashboard.vue`
- 实现内容：
  - 后台侧边栏接入真实路由，并支持当前菜单高亮。
  - 新增球队管理、比赛管理、积分榜管理、同步管理基础页面。
  - 球队和比赛页面支持列表、搜索、筛选和分页。
  - 积分榜页面支持小组筛选和手动重算。
  - 同步页面展示同步状态，并支持手动触发比赛同步。
- 验证结果：
  - `npm run build` 通过。

### 7.10 任务 8：性能优化（已完成）

#### 目标

保证移动端加载和滚动体验稳定。

#### 前端任务

1. 检查构建产物 chunk 体积。
2. 避免首页直接引入后台或重型依赖。
3. 图片继续使用懒加载。
4. 长列表保持分页加载。
5. 后续如果列表继续增长，增加虚拟滚动。

#### 后端任务

1. 所有列表接口确认支持分页。
2. 常用查询增加必要索引：
   - 比赛时间。
   - 比赛状态。
   - 小组。
   - 主队、客队。
3. 搜索接口避免无边界大查询。

#### 验收标准

- `npm run build` 通过。
- 首屏不引入明显无关的大模块。
- 赛程页滚动加载稳定。
- 列表接口默认分页。

#### 完成记录

- 完成时间：2026-06-04
- 修改文件：
  - `frontend/src/main.ts`
  - `backend/internal/models/match.go`
  - `backend/internal/models/team.go`
  - `backend/internal/models/stadium.go`
  - `backend/internal/models/reminder.go`
  - `backend/internal/models/notification.go`
- 实现内容：
  - 移除入口文件中未使用的 Element Plus 全量注册和全量 CSS 引入。
  - 保留现有路由级懒加载，避免后台页面进入用户端首屏主包。
  - 为赛程列表常用筛选字段增加索引：开球时间、状态、阶段、小组、主队、客队、城市、球场。
  - 为球队列表常用筛选字段增加索引：名称、英文名、洲、小组。
  - 为提醒扫描和通知未读数查询增加索引。
  - 保留赛程页现有分页和渐进加载策略。
- 构建体积变化：
  - 首屏 `index` JS 从约 `1085.94 kB` 降到 `158.12 kB`。
  - 首屏 `index` CSS 从约 `372.62 kB` 降到 `13.90 kB`。
  - Vite 大 chunk 警告已消失。
- 验证结果：
  - `npm run build` 通过。
  - `go test ./...` 通过。

### 7.11 任务 9：生产安全检查（已完成）

#### 目标

降低上线风险，避免默认账号、弱密钥、上传和 CORS 问题。

#### 前端任务

1. 登录页生产环境隐藏测试账号提示。
2. 上传头像前校验：
   - 文件类型。
   - 文件大小。
3. 请求错误统一处理，不暴露敏感信息。

#### 后端任务

1. 生产环境启动时校验：
   - `JWT_SECRET` 不能是默认值。
   - 默认管理员密码不能继续使用。
2. 上传接口校验：
   - MIME 类型。
   - 文件扩展名。
   - 文件大小。
3. CORS 根据环境配置允许域名。
4. 日志避免输出敏感配置。
5. 确认 `.env` 不进入 Git。

#### 验收标准

- 生产环境弱密钥无法启动。
- 非图片文件不能作为头像上传。
- 登录页生产环境不展示默认账号。
- `.env` 没有被 Git 跟踪。

#### 完成记录

- 完成时间：2026-06-04
- 新增文件：
  - `.gitignore`
- 删除文件：
  - `backend/internal/database/seed.go.1026024177`
- 修改文件：
  - `frontend/src/pages/user/LoginPage.vue`
  - `frontend/src/pages/user/ProfilePage.vue`
  - `frontend/src/api/request.ts`
  - `backend/cmd/server/main.go`
  - `backend/internal/config/config.go`
  - `backend/internal/database/seed.go`
  - `backend/internal/middleware/cors.go`
  - `backend/internal/services/auth_service.go`
- 实现内容：
  - 登录页测试账号提示仅在开发环境显示，生产构建隐藏。
  - 前端统一把 5xx 和网络错误转成通用提示，避免直接暴露后端内部错误信息。
  - 头像上传前端限制 `JPG/PNG/GIF/WebP` 和 `5MB` 大小。
  - 头像上传后端同时校验扩展名、文件大小和真实 MIME 类型。
  - 生产环境启动时拒绝默认 `JWT_SECRET`、空 `CORS_ALLOWED_ORIGINS` 和长度不足的 JWT 密钥。
  - 生产环境启动时检查默认管理员密码，发现 `admin123456` 会拒绝启动。
  - 生产环境不再自动创建默认管理员账号。
  - CORS 在开发环境默认开放，在生产环境按 `CORS_ALLOWED_ORIGINS` 白名单放行。
  - 根目录 `.env` 已从 Git 索引移除，并通过 `.gitignore` 防止再次提交。
  - 删除一个已跟踪的临时 seed 备份文件，避免把默认账号信息重复留在仓库。
- 验证结果：
  - `go test ./...` 通过。
  - `npm run build` 通过。
  - `git ls-files -- .env backend/.env` 无输出。

## 8. Agent 执行规范

### 8.1 每个任务的推荐工作流

1. 先读相关文件，不直接改。
2. 明确现有接口和状态管理是否可复用。
3. 小范围修改。
4. 补齐必要的空状态、加载状态和错误状态。
5. 跑验证命令。
6. 总结改动、影响范围和未完成项。

### 8.2 每个任务完成后的验证命令

前端任务：

```bash
cd frontend
npm run build
```

后端任务：

```bash
cd backend
go test ./...
```

### 8.3 不建议一次性完成的内容

以下内容不要和 P0 功能混在一个任务中做：

- 大规模 UI 重构。
- 更换 UI 框架。
- 重写数据同步服务。
- 重写认证系统。
- 一次性补全全部后台 CRUD。
- 引入复杂预测玩法。

### 8.4 推荐给 Agent 的第一条开发指令

```text
请先实现 PRODUCT_OPTIMIZATION_PLAN.md 中的任务 1：球队详情页。
要求：
1. 先检查现有 team/match/standing 接口和 store。
2. 优先复用已有接口。
3. 新增 /teams/:id 前端路由和 TeamDetailPage.vue。
4. TeamCard 点击进入详情页，收藏按钮不触发跳转。
5. 详情页展示球队信息、关注按钮、该球队比赛、所在小组积分。
6. 完成后运行 npm run build；如果改了后端，运行 go test ./...。
7. 不要启动前后端服务。
```
