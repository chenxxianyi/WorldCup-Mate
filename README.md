# WorldCup Mate

WorldCup Mate 是一个面向 2026 世界杯的赛程助手与赛事数据管理项目，已扩展支持五大联赛（英超、西甲、德甲、意甲、法甲）多赛事模式。项目包含面向用户的 Vue 3 前端、基于 Go/Gin 的后端 API、赛事数据同步能力，以及静态 H5 原型和设计系统页面。

## 当前状态说明

- `frontend/` 已实现用户端主要页面。当前生产构建存在两处 TypeScript 类型问题，见“构建与检查”。
- `backend/` 已实现启动入口、模型、路由、处理器、服务、仓库、数据同步、提醒和通知等模块。
- 后端入口位于 `backend/cmd/server/main.go`，启动时会读取配置、初始化 JWT/邮件/同步配置、连接 MySQL 和 Redis、执行 `AutoMigrate`、写入 seed 数据、注册静态上传目录、启动后台任务并监听 HTTP 端口。
- `backend/.gitignore` 中的 `server` 规则会让部分文件搜索命令默认忽略 `cmd/server` 目录；如需完整搜索可使用 `rg -u`。
- `docker-compose.yml` 中后端 `MYSQL_DSN` 使用 `xxladmin` 用户，但 MySQL 服务目前只配置了 `root` 密码和数据库名。直接 Docker 启动前需要让 MySQL 用户配置与 DSN 对齐。

## 功能概览

### 用户端

- 全局赛事切换：世界杯 / 英超 / 西甲 / 德甲 / 意甲 / 法甲（默认世界杯，行为不变）。
- 首页：下一场比赛倒计时、今日比赛、跨赛事"今日焦点"、热门推荐、赛事进度、小组积分速览 / 联赛积分榜 Top5、我的关注球队。
- 赛程页：世界杯按今日、明日、小组、淘汰赛、未开始状态和关键词筛选比赛；联赛按轮次（第 N 轮）与日期筛选。
- 比赛详情：展示对阵、比分、状态、城市/球场、小组积分（世界杯）或轮次（联赛），并支持设置比赛提醒。
- 球队页：国家队按小组或洲筛选、俱乐部按赛事筛选，支持关注球队。
- 积分榜页：世界杯查看 Group A 至 Group L 积分榜与最佳第三名榜；联赛查看 20 队单表（欧冠/欧战/降级区标记、主客场切换）。
- 我的页面：头像上传、关注球队（按赛事分组）、收藏比赛、比赛提醒、通知邮箱、提醒渠道、时区、主题和修改密码。
- 登录注册：基于 JWT 的用户登录、注册和本地 token 持久化。

### 后端能力

- 公开 API：比赛、球队、小组、积分榜、城市、球场和数据同步状态。
- 用户能力：资料维护、头像上传、关注球队、收藏比赛、比赛提醒、站内通知。
- 管理 API：球队、小组、城市、球场、比赛、比分状态、积分榜、用户状态和比赛导入。
- 积分计算：小组积分榜重算、最佳第三名计算、晋级状态标记。
- 数据同步：支持通过 football-data.org 同步 2026 世界杯与五大联赛比赛数据（`SYNC_COMPETITIONS` 配置，默认仅世界杯）。
- 联赛积分：联赛积分榜（TOTAL/HOME/AWAY）同步与按比赛重算，欧战/降级分区标记。
- 提醒通知：定时扫描比赛提醒，生成站内通知，并可按配置发送邮件通知。

### 原型与设计系统

- `index.html`：根目录静态 H5 demo。
- `h5-demo/index.html`：完整 H5 原型页面。
- `h5-demo/worldcup-mate-ds.html`：WorldCup Mate 设计系统页面。
- `WC26_Logo.webp`、`image-contact-sheet.jpg`：项目视觉资源。

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
├── frontend/                 # Vue 3 前端应用
│   ├── src/api/              # API 请求封装
│   ├── src/assets/           # 球队图标等静态资源
│   ├── src/components/       # 通用组件
│   ├── src/layouts/          # 用户端和后台布局
│   ├── src/pages/            # 用户端页面和后台页面
│   ├── src/router/           # Vue Router 配置
│   ├── src/stores/           # Pinia 状态管理
│   ├── src/styles/           # 全局样式和主题变量
│   └── vite.config.ts        # Vite 配置与代理
├── backend/                  # Go 后端 API
│   ├── cmd/server/           # 后端启动入口
│   ├── internal/config/      # 环境变量配置
│   ├── internal/database/    # MySQL、Redis、seed 数据
│   ├── internal/handlers/    # Gin handlers
│   ├── internal/jobs/        # 比赛同步和提醒扫描任务
│   ├── internal/middleware/  # CORS、日志、JWT、管理员鉴权
│   ├── internal/models/      # GORM 模型
│   ├── internal/providers/   # 第三方数据源客户端
│   ├── internal/repositories/# 数据访问层
│   ├── internal/routes/      # API 路由注册
│   ├── internal/services/    # 业务逻辑
│   └── internal/utils/       # 响应、JWT、密码、邮件等工具
├── h5-demo/                  # 静态 H5 原型和设计系统
├── docker-compose.yml        # MySQL、Redis、前端、后端编排
├── 前端开发方案.md
└── 后端开发方案.md
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

Vite 已配置代理：

- `/api` -> `http://localhost:8080`
- `/uploads` -> `http://localhost:8080`

### 配置后端环境变量

```bash
cd backend
cp .env.example .env
```

Windows PowerShell 可使用：

```powershell
cd backend
Copy-Item .env.example .env
```

根据本地数据库和 Redis 修改 `.env`：

```env
APP_ENV=development
APP_PORT=8080
JWT_SECRET=please_change_me_to_a_long_random_string
MYSQL_DSN=user:password@tcp(127.0.0.1:3306)/worldcup_mate?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=127.0.0.1:5001
REDIS_PASSWORD=
REDIS_DB=0
DATA_SYNC_ENABLED=false
FOOTBALL_DATA_API_KEY=

# 联赛同步（五大联赛）。格式：CODE:SEASON，逗号分隔。
# 留空 = 仅同步世界杯（默认，保持原有行为）。
# 示例：SYNC_COMPETITIONS=PL:2025,PD:2025,BL1:2025,SA:2025,FL1:2025
SYNC_COMPETITIONS=
LEAGUE_SYNC_INTERVAL_MINUTES=30
```

### 启动后端

```bash
cd backend
go mod download
go run ./cmd/server
```

后端启动入口会自动执行：

1. 读取配置：`config.Load()`
2. 初始化 JWT 与邮件配置
3. 配置 football-data.org 比赛同步参数
4. 初始化 MySQL 和 Redis
5. 执行 GORM `AutoMigrate`
6. 写入 seed 数据
7. 更新球队中文名称
8. 注册 API 路由与 `/uploads` 静态目录
9. 启动比赛提醒扫描和比赛数据同步任务
10. 监听 `APP_PORT`

## 构建与检查

前端生产构建：

```bash
cd frontend
npm run build
```

当前验证结果：`npm run build` 会因为既有类型问题失败。

- `src/data/mockUser.ts` 缺少 `User.notificationEmail`
- `src/pages/user/ProfilePage.vue` 使用了当前 `User` 类型中未声明的 `email`

后端包编译检查：

```bash
cd backend
go test ./...
```

## Docker 部署

项目提供了 `docker-compose.yml`，包含 MySQL、Redis、后端和前端服务。

```bash
docker compose up -d --build
```

默认端口：

| 服务 | 端口 |
| --- | --- |
| 前端 Nginx | `3000:80` |
| 后端 API | `8080:8080` |
| MySQL | `3306:3306` |
| Redis | `5001:6379` |

注意：当前 Docker 部署前需要先修正 `MYSQL_DSN` 与 MySQL 服务用户配置不一致的问题。

## 环境变量

| 变量 | 说明 |
| --- | --- |
| `APP_ENV` | 运行环境，默认 `development` |
| `APP_PORT` | 后端服务端口，默认 `8080` |
| `JWT_SECRET` | JWT 签名密钥，生产环境必须替换 |
| `MYSQL_DSN` | MySQL 连接串 |
| `REDIS_ADDR` | Redis 地址 |
| `REDIS_PASSWORD` | Redis 密码 |
| `REDIS_DB` | Redis 数据库编号 |
| `DATA_SYNC_ENABLED` | 是否启用第三方比赛数据同步 |
| `DATA_SYNC_PROVIDER` | 数据同步供应商，当前为 `football-data` |
| `DATA_SYNC_LIVE_INTERVAL_SECONDS` | 直播窗口同步间隔 |
| `DATA_SYNC_IDLE_INTERVAL_MINUTES` | 非直播窗口同步间隔 |
| `DATA_SYNC_FULL_INTERVAL_HOURS` | 全量同步间隔配置 |
| `FOOTBALL_DATA_API_KEY` | football-data.org API Key |
| `FOOTBALL_DATA_BASE_URL` | football-data.org API Base URL |
| `SMTP_HOST` | SMTP 主机 |
| `SMTP_PORT` | SMTP 端口 |
| `SMTP_USERNAME` | SMTP 用户名 |
| `SMTP_PASSWORD` | SMTP 密码 |
| `SMTP_FROM` | 邮件发送方 |

## API 摘要

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
- `GET /api/teams`
- `GET /api/teams/:id`
- `GET /api/groups`
- `GET /api/groups/:id/standings`
- `GET /api/standings`
- `GET /api/standings/best-third`
- `GET /api/cities`
- `GET /api/stadiums`
- `GET /api/sync/status`

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
- `PUT /api/reminders/:id`
- `DELETE /api/reminders/:id`
- `GET /api/notifications`
- `PUT /api/notifications/:id/read`
- `PUT /api/notifications/read-all`

### 管理接口

- `POST /api/admin/login`
- `GET /api/admin/dashboard`
- `/api/admin/teams`
- `/api/admin/groups`
- `/api/admin/cities`
- `/api/admin/stadiums`
- `/api/admin/matches`
- `/api/admin/standings`
- `/api/admin/users`
- `POST /api/admin/matches/import`
- `POST /api/admin/sync/matches`

管理接口需要 JWT，并且用户角色必须为 `admin`。

## 数据初始化

`backend/internal/database/seed.go` 中包含 seed 数据：

- 默认管理员：`admin@worldcup.local / admin123456`
- Group A 至 Group L
- 部分主办城市和球场
- 部分球队演示数据
- 当 `DATA_SYNC_ENABLED=false` 时，创建若干相对当前日期的演示比赛

生产环境首次部署后请立即修改默认管理员密码，并替换 `JWT_SECRET`。

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
| `/standings` | 积分榜 |
| `/profile` | 我的 |
| `/login` | 登录/注册 |
| `/admin` | 后台看板 |

## 开发备注

- 前端请求统一走 `frontend/src/api/request.ts`，响应拦截器要求后端返回 `{ code, message, data }` 结构。
- 用户 token 存储在 `localStorage` 的 `wm-token`。
- 主题、时区、语言和默认提醒渠道也通过 `localStorage` 保存。
- 头像上传接口返回 `/uploads/...` 路径，前端通过 Vite 或 Nginx 代理访问。
- 后台前端当前主要是 Dashboard 页面，完整后台 CRUD 前端页面仍可继续扩展。
