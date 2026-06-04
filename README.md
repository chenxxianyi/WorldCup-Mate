# WorldCup Mate

WorldCup Mate 是一个面向 2026 世界杯的赛事助手与赛事数据管理项目。当前仓库包含 Vue 3 移动端/管理端前端、Go/Gin 后端 API、MySQL/Redis 数据层、football-data.org 同步能力，以及早期 H5 原型和设计系统页面。

## 当前代码状态

- `frontend/` 是 Vue 3 + TypeScript + Vite 应用，包含用户端页面和后台管理页面。
- `backend/` 是 Go 1.23 + Gin + GORM API 服务，包含认证、赛事、球队、积分榜、收藏、提醒、通知、后台管理、数据同步和 AI 助手模块。
- 后端启动入口为 `backend/cmd/server/main.go`，启动时会加载配置、连接 MySQL/Redis、执行 `AutoMigrate`、写入 seed 数据、注册路由、挂载 `/uploads` 静态目录，并启动提醒扫描和比赛同步后台任务。
- 根目录 `README.md` 已按当前代码重新整理；仓库中部分中文 seed 数据和前端提示文案仍存在历史编码异常，需要后续单独修复。
- 当前 `docker-compose.yml` 不能直接按生产配置启动后端：`APP_ENV=production` 时需要足够强的 `JWT_SECRET` 和 `CORS_ALLOWED_ORIGINS`，同时 MySQL 服务用户与后端 `MYSQL_DSN` 中的 `xxladmin` 用户不一致。

## 功能概览

### 用户端

- 首页：下一场比赛倒计时、今日比赛、推荐比赛、赛事进度、小组积分速览、关注球队。
- 赛程：按今日、明日、小组赛、淘汰赛、状态和关键词筛选比赛。
- 比赛详情：展示对阵、比分、状态、城市、球场、小组积分，并支持收藏和提醒。
- 球队：球队搜索、按小组/洲筛选、关注球队、查看球队详情与相关比赛。
- 积分榜：查看 Group A 到 Group L 的小组积分，以及最佳第三名排行。
- 我的：头像上传、关注球队、收藏比赛、比赛提醒、站内通知、提醒渠道、时区、主题和密码修改。
- 登录注册：基于 JWT 的用户登录、注册和本地 token 持久化。

### 管理端

- 后台看板：赛事、球队、用户、提醒、同步状态等统计信息。
- 球队管理：球队列表、新增、编辑、删除。
- 赛事管理：比赛列表、新增、编辑、删除、比分更新、状态更新、批量导入。
- 积分榜管理：积分榜列表、手动编辑、重新计算。
- 数据同步：查看同步状态，并手动触发 football-data.org 比赛同步。

### 后端能力

- 公开 API：比赛、球队、小组、积分榜、城市、球场、数据同步状态。
- 登录用户 API：资料维护、头像上传、关注球队、收藏比赛、比赛提醒、站内通知。
- 管理 API：球队、小组、城市、球场、比赛、比分状态、积分榜、用户状态和比赛导入。
- 积分计算：小组积分重算、最佳第三名计算、晋级状态标记。
- 数据同步：通过 football-data.org 同步 2026 世界杯比赛、球队、比分和状态。
- 提醒通知：定时扫描比赛提醒，生成站内通知，并可按配置发送邮件通知。
- AI 助手：支持 OpenAI-compatible Provider，提供比赛看点、今日推荐、小组解读、规则解释、分享文案和登录用户聊天会话。

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 前端 | Vue 3, TypeScript, Vite 6, Pinia, Vue Router, Axios |
| UI | Element Plus, Tailwind CSS 4, 自定义响应式样式 |
| 后端 | Go 1.23, Gin, GORM |
| 数据库 | MySQL 8, Redis 7 |
| 认证 | JWT, bcrypt |
| 外部服务 | football-data.org API, SMTP |
| 部署 | Docker, Docker Compose, Nginx |

## 目录结构

```text
.
├── frontend/                  # Vue 3 前端应用
│   ├── src/api/               # API 请求封装
│   ├── src/assets/            # 球队图标等静态资源
│   ├── src/components/        # 通用组件
│   ├── src/layouts/           # 用户端和管理端布局
│   ├── src/pages/             # 用户端页面和后台页面
│   ├── src/router/            # Vue Router 配置
│   ├── src/stores/            # Pinia 状态管理
│   ├── src/styles/            # 全局样式和主题变量
│   └── vite.config.ts         # Vite 配置和开发代理
├── backend/                   # Go 后端 API
│   ├── cmd/server/            # 后端启动入口
│   ├── internal/ai/           # AI Provider、Prompt、Parser、Safety、上下文构建
│   ├── internal/config/       # 环境变量配置
│   ├── internal/database/     # MySQL、Redis、seed 数据
│   ├── internal/handlers/     # Gin handlers
│   ├── internal/jobs/         # 比赛同步和提醒扫描任务
│   ├── internal/middleware/   # CORS、日志、JWT、管理员鉴权
│   ├── internal/models/       # GORM 模型
│   ├── internal/providers/    # 第三方数据源客户端
│   ├── internal/repositories/ # 数据访问层
│   ├── internal/routes/       # API 路由注册
│   ├── internal/services/     # 业务逻辑
│   └── internal/utils/        # 响应、JWT、密码、邮件等工具
├── h5-demo/                   # 静态 H5 原型和设计系统
├── docker-compose.yml         # MySQL、Redis、前端、后端编排
├── worldcup_mate.sql          # 数据库 SQL 文件
├── index.html                 # 根目录静态 H5 demo
├── WC26_Logo.webp             # 视觉资源
└── image-contact-sheet.jpg    # 视觉资源
```

## 本地开发

### 前置要求

- Node.js 20 或更高版本
- npm
- Go 1.23 或更高版本
- MySQL 8
- Redis 7
- 可选：football-data.org API Key
- 可选：SMTP 邮件账号

### 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认访问地址：

```text
http://localhost:5173
```

Vite 开发代理：

| 路径 | 代理目标 |
| --- | --- |
| `/api` | `http://localhost:8080` |
| `/uploads` | `http://localhost:8080` |

### 配置后端

```bash
cd backend
cp .env.example .env
```

Windows PowerShell：

```powershell
cd backend
Copy-Item .env.example .env
```

根据本地 MySQL 和 Redis 修改 `backend/.env`：

```env
APP_ENV=development
APP_PORT=8080
JWT_SECRET=please_change_me_to_a_long_random_string
MYSQL_DSN=xxladmin:XXLadmin_2021!@tcp(127.0.0.1:3310)/worldcup_mate?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
DATA_SYNC_ENABLED=false
DATA_SYNC_PROVIDER=football-data
DATA_SYNC_LIVE_INTERVAL_SECONDS=120
DATA_SYNC_IDLE_INTERVAL_MINUTES=30
DATA_SYNC_FULL_INTERVAL_HOURS=6
FOOTBALL_DATA_API_KEY=
FOOTBALL_DATA_BASE_URL=https://api.football-data.org/v4
AI_PROVIDER=openai
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=your-ai-api-key
AI_MODEL=gpt-4o-mini
AI_TIMEOUT_SECONDS=60
AI_DAILY_LIMIT_USER=50
AI_TEMPERATURE=0.7
AI_MAX_TOKENS=1200
AI_CACHE_ENABLED=true
```

AI 使用 OpenAI-compatible 服务。请设置 `AI_PROVIDER=openai` 或供应商名、`AI_BASE_URL`、`AI_API_KEY` 和 `AI_MODEL`；后端不会把 API Key 写入响应或 usage log。未配置 `AI_API_KEY` 时，AI 接口会返回明确的配置错误，不会生成假数据。

### 启动后端

```bash
cd backend
go mod download
go run ./cmd/server
```

开发环境默认监听：

```text
http://localhost:8080
```

启动流程：

1. 加载 `.env` 和系统环境变量。
2. 初始化 JWT、邮件、比赛同步和 AI Provider 配置。
3. 连接 MySQL 和 Redis。
4. 执行 GORM `AutoMigrate`。
5. 写入基础 seed 数据。
6. 更新球队中文名，并应用官方球场映射。
7. 注册 API 路由和 `/uploads` 静态资源目录。
8. 启动比赛提醒扫描和比赛数据同步后台任务。
9. 监听 `APP_PORT` 指定端口。

## 构建与检查

前端生产构建：

```bash
cd frontend
npm run build
```

后端测试：

```bash
cd backend
go test ./...
```

说明：README 更新时未重新修改业务代码。若构建失败，优先检查仓库中现有的中文编码异常、TypeScript 类型定义和本地数据库配置。

## Docker 部署

仓库提供了 `docker-compose.yml`，包含 MySQL、Redis、后端和前端服务：

```bash
docker compose up -d --build
```

默认端口：

| 服务 | 端口 |
| --- | --- |
| 前端 Nginx | `3000:80` |
| 后端 API | `8080:8080` |
| MySQL | `3306:3306` |
| Redis | `6379:6379` |

当前 Docker 配置需要先调整后再用于完整启动：

- `backend` 使用 `APP_ENV=production`，因此必须配置长度至少 32 位的 `JWT_SECRET`。
- 生产环境必须配置 `CORS_ALLOWED_ORIGINS`，例如 `http://localhost:3000`。
- MySQL 服务当前只创建 `root` 用户和 `worldcup_mate` 数据库，但后端 `MYSQL_DSN` 使用 `xxladmin` 用户。需要改为 root DSN，或为 MySQL 服务增加 `MYSQL_USER` 和 `MYSQL_PASSWORD`。

## 环境变量

| 变量 | 说明 |
| --- | --- |
| `APP_ENV` | 运行环境，默认 `development` |
| `APP_PORT` | 后端服务端口，默认 `8080` |
| `JWT_SECRET` | JWT 签名密钥，生产环境必须替换为强随机值 |
| `MYSQL_DSN` | MySQL 连接串 |
| `REDIS_ADDR` | Redis 地址 |
| `REDIS_PASSWORD` | Redis 密码 |
| `REDIS_DB` | Redis 数据库编号 |
| `DATA_SYNC_ENABLED` | 是否启用第三方比赛数据同步 |
| `DATA_SYNC_PROVIDER` | 数据同步供应商，当前为 `football-data` |
| `DATA_SYNC_LIVE_INTERVAL_SECONDS` | 直播窗口同步间隔 |
| `DATA_SYNC_IDLE_INTERVAL_MINUTES` | 非直播窗口同步间隔 |
| `DATA_SYNC_FULL_INTERVAL_HOURS` | 全量同步间隔 |
| `FOOTBALL_DATA_API_KEY` | football-data.org API Key |
| `FOOTBALL_DATA_BASE_URL` | football-data.org API Base URL |
| `AI_PROVIDER` | AI Provider，默认 `openai` |
| `AI_BASE_URL` | OpenAI-compatible API Base URL |
| `AI_API_KEY` | AI Provider API Key |
| `AI_MODEL` | AI 模型名称 |
| `AI_TIMEOUT_SECONDS` | AI 调用超时时间 |
| `AI_DAILY_LIMIT_USER` | 登录用户每日 AI 调用上限 |
| `AI_TEMPERATURE` | AI 生成温度 |
| `AI_MAX_TOKENS` | AI 单次最大输出 token |
| `AI_CACHE_ENABLED` | 是否启用 AI 缓存 |
| `SMTP_HOST` | SMTP 主机 |
| `SMTP_PORT` | SMTP 端口 |
| `SMTP_USERNAME` | SMTP 用户名 |
| `SMTP_PASSWORD` | SMTP 密码 |
| `SMTP_FROM` | 邮件发送方 |
| `CORS_ALLOWED_ORIGINS` | 生产环境允许的前端 Origin，多个值用英文逗号分隔 |

## API 摘要

所有业务接口默认挂载在 `/api` 下，响应结构由后端工具统一封装为 `{ code, message, data }`。

### 公开接口

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/matches`
- `GET /api/matches/today`
- `GET /api/matches/tomorrow`
- `GET /api/matches/upcoming`
- `GET /api/matches/live`
- `GET /api/matches/recommended`
- `GET /api/matches/progress`
- `GET /api/matches/:id`
- `GET /api/matches/by-team/:teamId`
- `GET /api/matches/by-group/:groupId`
- `GET /api/matches/by-stage/:stage`
- `GET /api/teams`
- `GET /api/teams/:id`
- `GET /api/teams/:id/matches`
- `GET /api/groups`
- `GET /api/groups/:id`
- `GET /api/groups/:id/standings`
- `GET /api/standings`
- `GET /api/standings/best-third`
- `GET /api/cities`
- `GET /api/stadiums`
- `GET /api/stadiums/:id`
- `GET /api/sync/status`
- `POST /api/ai/match-insight`
- `POST /api/ai/today-recommendations`
- `POST /api/ai/group-analysis`
- `POST /api/ai/explain`
- `POST /api/ai/share-copy`

### 登录后接口

- `GET /api/user/profile`
- `PUT /api/user/profile`
- `PUT /api/user/password`
- `POST /api/user/avatar`
- `GET /api/favorites/teams`
- `POST /api/favorites/teams/:teamId`
- `DELETE /api/favorites/teams/:teamId`
- `GET /api/favorites/matches`
- `POST /api/favorites/matches/:matchId`
- `DELETE /api/favorites/matches/:matchId`
- `GET /api/reminders`
- `POST /api/reminders`
- `POST /api/reminders/batch`
- `PUT /api/reminders/:id`
- `DELETE /api/reminders/:id`
- `GET /api/notifications`
- `GET /api/notifications/unread-count`
- `PUT /api/notifications/:id/read`
- `PUT /api/notifications/read-all`
- `POST /api/ai/chat`
- `GET /api/ai/conversations`
- `GET /api/ai/conversations/:id`
- `DELETE /api/ai/conversations/:id`

### 管理接口

- `POST /api/admin/login`
- `GET /api/admin/dashboard`
- `GET /api/admin/teams`
- `POST /api/admin/teams`
- `PUT /api/admin/teams/:id`
- `DELETE /api/admin/teams/:id`
- `GET /api/admin/groups`
- `POST /api/admin/groups`
- `PUT /api/admin/groups/:id`
- `GET /api/admin/cities`
- `POST /api/admin/cities`
- `PUT /api/admin/cities/:id`
- `DELETE /api/admin/cities/:id`
- `GET /api/admin/stadiums`
- `POST /api/admin/stadiums`
- `PUT /api/admin/stadiums/:id`
- `DELETE /api/admin/stadiums/:id`
- `GET /api/admin/matches`
- `POST /api/admin/matches`
- `PUT /api/admin/matches/:id`
- `DELETE /api/admin/matches/:id`
- `PUT /api/admin/matches/:id/score`
- `PUT /api/admin/matches/:id/status`
- `POST /api/admin/matches/import`
- `POST /api/admin/sync/matches`
- `GET /api/admin/standings`
- `POST /api/admin/standings/recalculate`
- `PUT /api/admin/standings/:id`
- `GET /api/admin/users`
- `PUT /api/admin/users/:id/status`

管理接口需要 JWT，并且用户角色必须为 `admin`。

## 数据初始化

`backend/internal/database/seed.go` 会在开发环境写入基础数据：

- 默认管理员：`admin@worldcup.local / admin123456`
- Group A 到 Group L
- 部分主办城市和球场
- 部分球队演示数据
- 当 `DATA_SYNC_ENABLED=false` 时，创建若干相对当前日期的演示比赛

生产环境不会创建默认管理员，并且会拒绝继续使用默认管理员密码。首次部署后请使用安全方式创建管理员账号，并替换 `JWT_SECRET`。

## 数据同步

比赛同步由 `backend/internal/services/sync_service.go` 和 `backend/internal/jobs/match_syncer.go` 实现。

- 供应商：football-data.org
- 赛事代码：`WC`
- 赛季：`2026`
- 同步内容：比赛、球队、小组、比分、状态、胜者、开球时间
- 同步后会按需要重算小组积分和最佳第三名

启用方式：

```env
DATA_SYNC_ENABLED=true
FOOTBALL_DATA_API_KEY=your_api_key
```

## 前端路由

| 路径 | 页面 |
| --- | --- |
| `/` | 首页 |
| `/schedule` | 赛程 |
| `/matches/:id` | 比赛详情 |
| `/teams` | 球队 |
| `/teams/:id` | 球队详情 |
| `/standings` | 积分榜 |
| `/ai` | AI 助手首页 |
| `/ai/chat` | AI 聊天助手 |
| `/ai/match/:id` | AI 比赛看点 |
| `/ai/share-copy` | AI 分享文案 |
| `/profile` | 我的 |
| `/login` | 登录/注册 |
| `/admin` | 后台看板 |
| `/admin/teams` | 球队管理 |
| `/admin/matches` | 赛事管理 |
| `/admin/standings` | 积分榜管理 |
| `/admin/sync` | 数据同步 |

## 开发备注

- 前端请求统一从 `frontend/src/api/request.ts` 发起，默认 `baseURL` 为 `/api`。
- 用户 token 存储在 `localStorage` 的 `wm-token`。
- 主题、时区、语言和默认提醒渠道也通过 `localStorage` 保存。
- 头像上传接口返回 `/uploads/...` 路径，开发环境通过 Vite 代理访问，Docker 环境需要确保 Nginx 也代理或暴露该路径。
- 后端开发环境 CORS 默认放开；生产环境必须通过 `CORS_ALLOWED_ORIGINS` 明确配置允许来源。
- 根目录 `.env` 当前不是后端示例配置；后端本地开发请使用 `backend/.env.example` 复制生成 `backend/.env`。
