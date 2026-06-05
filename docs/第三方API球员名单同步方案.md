# WorldCup Mate 第三方 API 球员名单同步方案

## 1. 目标

为球队详情页提供真实球员大名单数据，通过第三方 API 同步球员姓名、号码、照片、位置，并写入本地 `players` 表。

本方案明确要求：

- 不使用演示球员数据。
- 不生成兜底球员数据。
- 第三方 API 没有返回的数据不由系统伪造。
- 球员名单为空时接口返回空数组 `[]`，前端展示“球员大名单暂未同步”。

## 2. 当前状态

当前已完成：

- 后端已有 `Player` 模型：`backend/internal/models/player.go`
- 后端已有查询接口：`GET /api/teams/:id/players`
- 前端球队详情页已接入该接口
- 当前后端不会 seed 演示球员
- 当前 seed 只清理历史 `source = demo_seed` 的旧演示记录

当前缺口：

- 没有真实球员数据源
- 没有第三方球队 ID 与本地球队 ID 的映射
- 没有球员同步 provider
- 没有球员 upsert 同步逻辑
- 没有后台手动同步入口

## 3. 推荐数据源

### 3.1 首选：API-Football / API-SPORTS

推荐使用 API-Football 作为第一版球员名单数据源。

理由：

- 支持按球队 ID 获取 squad。
- 字段覆盖球员姓名、号码、位置。
- 球员照片可通过官方 media URL 获取。
- 文档建议 `players/squads` 每周同步一次，适合阵容名单这类低频变化数据。

接口：

```http
GET https://v3.football.api-sports.io/players/squads?team={external_team_id}
Header: x-apisports-key: {API_FOOTBALL_KEY}
```

官方文档：

```text
https://www.api-football.com/documentation-v3#operation/get-players-squads
```

### 3.2 备选：Sportmonks

Sportmonks 可作为后续增强方案，适合需要更多球员资料、赛季阵容、球员详情时使用。

官方文档：

```text
https://docs.sportmonks.com/football/endpoints-and-entities/endpoints/team-squads
```

### 3.3 不优先：football-data.org

当前项目已经使用 football-data.org 同步比赛数据，但它的球队 squad 数据通常不包含球员照片，不满足本次“照片”需求，因此不作为球员大名单首选。

## 4. 数据字段映射

第三方 API 返回字段到本地 `players` 表的映射：

| 第三方字段 | 本地字段 | 说明 |
| --- | --- | --- |
| `player.id` | `source_player_id` | 第三方球员 ID，用于 upsert |
| `player.name` | `name` | 展示名称 |
| `player.name` | `name_en` | 第一版可与 `name` 相同 |
| `player.number` | `shirt_number` | 球衣号码 |
| `player.position` | `position` | 原始位置 |
| position 映射结果 | `position_label` | 仅基于真实 position 做标准化，不创造球员 |
| `player.photo` 或 media URL | `photo_url` | 真实照片 URL |
| `api-football` | `source` | 数据来源 |

API-Football 照片 URL 可按球员 ID 生成：

```text
https://media.api-sports.io/football/players/{player_id}.png
```

如果第三方没有球员照片，则 `photo_url` 保持空，不写入默认图片。

## 5. 数据库设计

### 5.1 扩展 players 表

当前 `Player` 模型已有核心字段。建议新增：

```go
ExternalTeamID string    `gorm:"size:100;index" json:"external_team_id"`
IsActive       bool      `gorm:"default:true;index" json:"is_active"`
LastSyncedAt   time.Time `json:"last_synced_at"`
```

建议增加唯一索引：

```text
team_id + source + source_player_id
```

原因：

- 球衣号码会变化，不能用号码做唯一键。
- 球员姓名可能重名或格式变化，不能只用姓名做唯一键。
- 第三方球员 ID 最适合做稳定同步键。

### 5.2 新增第三方球队映射表

新增模型：

```text
backend/internal/models/external_team_mapping.go
```

建议结构：

```go
type ExternalTeamMapping struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TeamID           uint      `gorm:"uniqueIndex:idx_provider_team;not null" json:"team_id"`
	Team             Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Provider         string    `gorm:"uniqueIndex:idx_provider_team;size:50;not null" json:"provider"`
	ExternalTeamID   string    `gorm:"size:100;index;not null" json:"external_team_id"`
	ExternalTeamName string    `gorm:"size:100" json:"external_team_name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
```

用途：

- 本地 `teams.id` 与 API-Football 的 team id 不一定一致。
- 同步球员前必须先通过映射表找到第三方 team id。
- 映射表可以由后台手动维护，也可以后续做自动搜索辅助。

## 6. 配置设计

在 `backend/.env` 增加：

```env
PLAYER_SYNC_ENABLED=true
PLAYER_SYNC_PROVIDER=api-football
API_FOOTBALL_KEY=your_api_key
API_FOOTBALL_BASE_URL=https://v3.football.api-sports.io
PLAYER_SYNC_INTERVAL_HOURS=168
```

在 `backend/internal/config/config.go` 增加字段：

```go
PlayerSyncEnabled       bool
PlayerSyncProvider      string
APIFootballKey          string
APIFootballBaseURL      string
PlayerSyncIntervalHours int
```

## 7. 后端模块设计

### 7.1 Provider 层

新增目录：

```text
backend/internal/providers/apifootball
```

文件：

```text
client.go
types.go
```

核心方法：

```go
func NewClient(baseURL, apiKey string) *Client
func (c *Client) TeamSquad(ctx context.Context, externalTeamID string) (*SquadResponse, error)
```

Provider 层只负责：

- 请求第三方 API
- 解析响应
- 返回原始 provider 数据结构

Provider 层不负责：

- 写数据库
- 处理本地 team id
- 做业务 upsert

### 7.2 Repository 层

新增：

```text
backend/internal/repositories/player_sync_repo.go
backend/internal/repositories/external_team_mapping_repo.go
```

建议方法：

```go
func GetExternalTeamMapping(teamID uint, provider string) (*models.ExternalTeamMapping, error)
func ListExternalTeamMappings(provider string) ([]models.ExternalTeamMapping, error)
func UpsertPlayer(player *models.Player) error
func MarkMissingPlayersInactive(teamID uint, source string, activeSourcePlayerIDs []string) error
```

### 7.3 Service 层

新增：

```text
backend/internal/services/player_sync_service.go
```

核心方法：

```go
func ConfigurePlayerSync(cfg PlayerSyncConfig)
func SyncTeamPlayers(ctx context.Context, teamID uint, reason string) (*PlayerSyncResult, error)
func SyncAllMappedTeamPlayers(ctx context.Context, reason string) (*PlayerSyncResult, error)
```

同步流程：

1. 检查 `PLAYER_SYNC_ENABLED`。
2. 根据 `teamID + provider` 查询 `ExternalTeamMapping`。
3. 用 `external_team_id` 请求 API-Football squad。
4. 转换球员字段。
5. 使用 `team_id + source + source_player_id` upsert。
6. 本次 API 没返回但本地已有的同 source 球员，标记 `is_active = false`。
7. 写入同步结果和错误日志。

### 7.4 Handler 层

新增管理员接口：

```http
POST /api/admin/teams/:id/sync-players
POST /api/admin/sync/players
GET  /api/admin/teams/:id/player-mapping
PUT  /api/admin/teams/:id/player-mapping
```

建议放在：

```text
backend/internal/handlers/player_sync_handler.go
```

说明：

- 同步接口必须放 admin 路由，不对普通用户开放。
- 用户端仍然只访问 `GET /api/teams/:id/players`。

## 8. 查询接口调整

当前接口：

```http
GET /api/teams/:id/players
```

建议调整查询条件：

```sql
WHERE team_id = ? AND is_active = true
ORDER BY shirt_number ASC, id ASC
```

如果暂时不加 `is_active` 字段，则保持当前查询即可。

返回示例：

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "team_id": 1,
      "name": "真实球员姓名",
      "name_en": "Real Player Name",
      "shirt_number": 10,
      "position": "Attacker",
      "position_label": "前锋",
      "photo_url": "https://media.api-sports.io/football/players/123.png",
      "club": "",
      "source": "api-football",
      "source_player_id": "123"
    }
  ]
}
```

## 9. 位置标准化

API-Football 可能返回：

- `Goalkeeper`
- `Defender`
- `Midfielder`
- `Attacker`

本地可以标准化：

| Provider position | position | position_label |
| --- | --- | --- |
| `Goalkeeper` | `GK` | `门将` |
| `Defender` | `DF` | `后卫` |
| `Midfielder` | `MF` | `中场` |
| `Attacker` | `FW` | `前锋` |

注意：

- 这是对真实 API 字段做标准化，不是生成兜底球员。
- 如果 API 没有 position，则本地 position 保持空。

## 10. 同步状态和日志

可以复用现有 `SyncState`，也可以新增资源类型：

```text
resource = players
provider = api-football
```

建议记录：

- provider
- resource
- status：running / success / failed
- last_error
- last_synced_at
- next_sync_at

同步结果结构：

```go
type PlayerSyncResult struct {
	Provider   string    `json:"provider"`
	Resource   string    `json:"resource"`
	Reason     string    `json:"reason"`
	Teams      int       `json:"teams"`
	Total      int       `json:"total"`
	Created    int       `json:"created"`
	Updated    int       `json:"updated"`
	Deactivated int      `json:"deactivated"`
	Skipped    int       `json:"skipped"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
}
```

## 11. 定时任务

新增：

```text
backend/internal/jobs/player_syncer.go
```

策略：

- 默认每 168 小时同步一次，也就是每周一次。
- 项目启动后不要立即同步全部球队，避免启动时阻塞和触发限流。
- 可延迟 1-5 分钟后启动定时任务。
- 同步全部球队时按映射表逐队同步。
- 每队之间加短暂 sleep，降低第三方 API 限流风险。

## 12. 管理后台建议

第一版可以先只做接口，不做复杂 UI。

后续后台页面可以增加：

- 球队第三方 ID 映射维护
- 单队同步按钮
- 全量同步按钮
- 最近同步时间
- 最近同步错误
- 本地球员数量

## 13. 开发步骤

### P0：单队同步闭环

1. 增加 `.env` 配置。
2. 扩展 `Player` 模型字段。
3. 新增 `ExternalTeamMapping` 模型并加入 `AutoMigrate`。
4. 新增 API-Football provider client。
5. 新增单队球员同步 service。
6. 新增 admin 单队同步接口。
7. 手动写入一条墨西哥映射。
8. 调用 `POST /api/admin/teams/1/sync-players`。
9. 验证 `GET /api/teams/1/players` 返回真实球员。

### P1：全量同步

1. 新增映射列表查询。
2. 新增全量同步 service。
3. 新增 `POST /api/admin/sync/players`。
4. 复用 `SyncState` 记录同步状态。
5. 增加错误统计和跳过统计。

### P2：定时同步和后台维护

1. 新增 `player_syncer.go` 定时任务。
2. 后台增加球队映射维护。
3. 后台增加同步按钮和同步状态展示。
4. 增加同步失败重试策略。

## 14. 验收标准

功能验收：

- 可以为某一支球队配置第三方 team id。
- 管理员可以触发单队球员同步。
- 同步后 `players` 表出现真实球员数据。
- 用户端 `GET /api/teams/:id/players` 返回真实球员。
- 球队详情页展示球员姓名、号码、照片、位置。
- 没有真实数据时仍返回空数组，不生成演示数据。

数据验收：

- `source = api-football`
- `source_player_id` 不能为空。
- 同一个球员重复同步不会重复插入。
- 第三方未返回照片时 `photo_url` 为空。
- 第三方未返回号码时 `shirt_number` 为空或 0，不自动补号码。

稳定性验收：

- API key 未配置时同步接口返回明确错误。
- 映射不存在时返回明确错误。
- 第三方 API 失败时不影响用户端球队详情接口。
- `go test ./...` 通过。
- `npm.cmd run build` 通过。

## 15. 风险和注意事项

- 第三方 API 可能有额度限制，全量同步需要限速。
- 2026 世界杯正式名单未公布前，第三方 squad 可能不是最终名单。
- 国家队球衣号码可能随赛事变化，不能用号码做唯一键。
- 球员照片 URL 可能失效，前端只展示真实 URL，失败时不生成假照片数据。
- 球队映射是关键工作，映射错误会导致同步到错误球队。

## 16. 推荐落地顺序

建议先完成 P0。

第一支验证球队可以选截图中的墨西哥：

```text
local team: MEX
local team_id: 1
provider: api-football
external_team_id: 需要通过 API-Football 查询确认
```

墨西哥同步成功后，再扩展法国、阿根廷、巴西等球队，最后做全量同步和后台维护。

## 开发完成记录

更新时间：2026-06-05

已根据本方案完成第一版第三方 API 球员名单同步能力：

- 已新增 API-Football provider：`backend/internal/providers/apifootball/client.go`、`types.go`。
- 已扩展 `Player` 模型，支持 `external_team_id`、`is_active`、`last_synced_at`。
- 已新增第三方球队映射模型：`backend/internal/models/external_team_mapping.go`。
- 已在 `AutoMigrate` 中加入 `ExternalTeamMapping`。
- 已新增球队映射 repository：`backend/internal/repositories/external_team_mapping_repo.go`。
- 已扩展 player repository，支持按 active 球员查询、按 `team_id + source + source_player_id` upsert、停用旧球员。
- 已新增球员同步 service：`backend/internal/services/player_sync_service.go`。
- 已新增管理员接口：`backend/internal/handlers/player_sync_handler.go`。
- 已新增球员定时同步任务：`backend/internal/jobs/player_syncer.go`。
- 已在后台路由注册映射维护、单队同步、全量同步接口。
- 已在 `.env.example` 增加球员同步相关配置。

当前实现仍然遵守“不使用演示数据、不生成兜底数据”的要求：

- 只有第三方 API 返回的真实球员会写入 `players` 表。
- 第三方未返回照片时，`photo_url` 保持为空。
- 第三方未返回位置时，`position` 保持为空。
- 同步没有拿到任何有效球员时，不会批量停用本地已有球员，避免第三方临时异常清空名单。
- 用户端 `GET /api/teams/:id/players` 只返回 `is_active = true` 的球员。

## 实际使用步骤

1. 在 `backend/.env` 配置：

```env
PLAYER_SYNC_ENABLED=true
PLAYER_SYNC_PROVIDER=api-football
API_FOOTBALL_KEY=your_api_key
API_FOOTBALL_BASE_URL=https://v3.football.api-sports.io
PLAYER_SYNC_INTERVAL_HOURS=168
```

2. 重启后端服务，使配置和数据库迁移生效。

3. 为本地球队配置第三方球队 ID：

```http
PUT /api/admin/teams/:id/player-mapping
```

请求体示例：

```json
{
  "provider": "api-football",
  "external_team_id": "463",
  "external_team_name": "Mexico"
}
```

4. 手动同步单支球队：

```http
POST /api/admin/teams/:id/sync-players
```

5. 手动同步全部已配置映射的球队：

```http
POST /api/admin/sync/players
```

6. 用户端查看球队球员名单：

```http
GET /api/teams/:id/players
```

## 验证结果

- `go test ./...` 通过。
- `go build ./...` 通过。
- `npm.cmd run build` 通过。
