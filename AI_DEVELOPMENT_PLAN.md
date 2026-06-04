# WorldCup Mate V2 自用版 AI 看球助手开发方案

## 1. 文档目标

本文档用于指导 WorldCup Mate 在现有项目基础上扩展一个“自用为主”的 V2 AI 看球助手模块。

自用版的核心目标不是商业化、不是开放给大量用户，也不是做复杂运营后台，而是让项目主人和少量可信用户在本地或私有部署环境里，更舒服地完成世界杯观赛决策、赛程理解、出线分析和分享文案生成。

## 2. 自用版定位

WorldCup Mate V2 自用版定位：

> 一个私人世界杯 AI 观赛助手，帮助自己更快判断看什么、怎么看、为什么值得看。

设计取舍：

- 优先个人体验，不追求公开用户增长。
- 优先本地可用，不优先复杂生产部署。
- 优先 Mock Provider 和 OpenAI-compatible Provider，不做多供应商管理平台。
- 优先比赛看点、今日推荐、小组分析、聊天问答，不做复杂运营后台。
- 优先低维护成本，不引入复杂 RAG、付费、积分、用户等级。
- 保留安全边界，避免编造比分、伤病、新闻和投注建议。

## 3. 当前项目基础

### 3.1 可复用能力

现有项目已经具备：

- 用户登录和 JWT。
- 赛程、今日比赛、明日比赛、比赛详情。
- 球队列表、球队详情。
- 小组积分榜和最佳第三名榜。
- 关注球队、收藏比赛。
- 比赛提醒、通知、邮件提醒。
- football-data.org 数据同步。
- MySQL、Redis、Docker Compose。
- 管理后台基础能力。

这些能力足够支撑自用 AI 模块，不需要重做业务底座。

### 3.2 必须先处理的问题

当前 `worldcup_mate.sql` 中部分中文字段存在乱码，例如球队名、城市名、球场名、大洲、推荐标签。

自用版可以接受功能简单，但不能接受 AI 上下文脏数据。因此开发前必须先做最小数据治理：

- 修复 `teams.name`、`teams.continent`。
- 修复 `cities.name`。
- 修复 `stadia.name`。
- 修复 `matches.recommend_tag`、`matches.recommend_reason`。
- 尽量保留 `name_en`、`fifa_code`、`external_id` 作为稳定匹配依据。
- 不强制做完复杂数据合并，但不能让 AI 读到明显乱码。

## 4. 自用版功能范围

### 4.1 第一版必须完成

自用第一版只做能形成闭环的功能：

- 比赛 AI 看点。
- 今日 AI 观赛推荐。
- AI 聊天助手。
- 小组出线形势 AI 解读。
- 足球规则和术语解释。
- 分享文案生成。
- Mock AI Provider。
- OpenAI-compatible Provider。
- AI 调用日志。
- 简单 Redis 限流或本地防刷。

### 4.2 第一版可以暂缓

以下内容自用版暂缓：

- 游客体系。
- 完整运营后台。
- AI 用户排行。
- Token 复杂统计。
- 缓存命中率统计。
- 多模型供应商管理。
- 流式输出。
- 语音问答。
- 实时陪看。
- 新闻聚合。
- 复杂 RAG。
- 多语言。
- 付费系统。
- Web Push。
- 公网生产级运维体系。

### 4.3 自用版建议顺序

第一阶段先做最小闭环：

1. 数据清洗。
2. 后端 AI 配置。
3. Mock Provider。
4. OpenAI-compatible Provider。
5. MatchContext。
6. `/api/ai/match-insight`。
7. 比赛详情页 AI 看点卡片。

完成后再扩展：

1. 今日推荐。
2. 小组分析。
3. 聊天助手。
4. 分享文案。
5. 观赛计划。

## 5. 前端设计

### 5.1 新增路由

自用版保留核心路由：

```text
/ai
/ai/chat
/ai/match/:id
/ai/share-copy
```

`/ai/viewing-plan` 可以第二阶段再加，也可以先集成在 `/ai` 页面里。

### 5.2 页面融合点

优先融合这些已有页面：

- 比赛详情页：加入 AI 看点卡片，这是第一优先级。
- 首页：加入“今天看什么”AI 推荐。
- 积分榜页：加入“小组出线 AI 解读”。
- 赛程页：比赛卡片增加 AI 看点入口。

自用版不强制在个人中心展示复杂 AI 历史，也不强制做完整 AI 内容管理页。

### 5.3 前端目录

建议新增：

```text
frontend/src/
├── pages/
│   └── user/
│       ├── AIHomePage.vue
│       ├── AIChatPage.vue
│       ├── AIMatchInsightPage.vue
│       └── AIShareCopyPage.vue
├── components/
│   └── ai/
│       ├── AIChatPanel.vue
│       ├── AIMessageBubble.vue
│       ├── AIInputBox.vue
│       ├── AIThinking.vue
│       ├── PromptSuggestionCard.vue
│       ├── MatchInsightCard.vue
│       └── ShareCopyCard.vue
├── api/
│   └── ai.ts
├── stores/
│   └── useAIStore.ts
└── types/
    └── ai.ts
```

## 6. 后端设计

### 6.1 新增目录

```text
backend/internal/
├── ai/
│   ├── provider.go
│   ├── mock_provider.go
│   ├── openai_compatible.go
│   ├── prompts.go
│   ├── context_builder.go
│   ├── safety.go
│   └── parser.go
├── handlers/
│   └── ai_handler.go
├── services/
│   └── ai_service.go
├── repositories/
│   └── ai_repo.go
└── models/
    ├── ai_conversation.go
    ├── ai_message.go
    ├── ai_generated_content.go
    └── ai_usage_log.go
```

自用版第一版不强制做 `ai_feedbacks`。如果需要点赞点踩，后续再加。

### 6.2 分层原则

- Handler 只负责参数和响应。
- Service 负责编排限流、上下文、Provider、日志。
- ContextBuilder 只负责从数据库组装事实。
- Provider 只负责调用模型。
- Safety 负责基础输出安全兜底。
- Parser 负责 JSON 解析和 Markdown 兜底。

## 7. 数据库设计

自用版保留 4 张 AI 表即可。

### 7.1 ai_conversations

保存聊天会话。

字段建议：

```text
id
user_id
title
context_type
context_id
last_message
created_at
updated_at
deleted_at
```

### 7.2 ai_messages

保存聊天消息。

字段建议：

```text
id
conversation_id
user_id
role
content
provider
model
prompt_tokens
completion_tokens
total_tokens
created_at
```

### 7.3 ai_generated_contents

保存比赛看点、今日推荐、小组分析、分享文案。

字段建议：

```text
id
user_id
type
target_type
target_id
content_json
content_markdown
provider
model
cache_key
created_at
updated_at
```

### 7.4 ai_usage_logs

保存调用日志，方便自查成本和错误。

字段建议：

```text
id
user_id
ip
endpoint
provider
model
prompt_tokens
completion_tokens
total_tokens
status
error_message
latency_ms
created_at
```

## 8. API 设计

### 8.1 第一版接口

```text
POST /api/ai/match-insight
POST /api/ai/today-recommendations
POST /api/ai/group-analysis
POST /api/ai/explain
POST /api/ai/share-copy
POST /api/ai/chat
GET /api/ai/conversations
GET /api/ai/conversations/:id
DELETE /api/ai/conversations/:id
```

### 8.2 暂缓接口

```text
POST /api/ai/viewing-plan
POST /api/ai/post-match-summary
POST /api/ai/feedback
GET /api/admin/ai/stats
GET /api/admin/ai/usage-logs
GET /api/admin/ai/feedback
GET /api/admin/ai/generated-contents
```

这些接口不是不要，而是自用第一版没必要抢先做。

### 8.3 响应格式

继续沿用项目已有格式：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 9. 核心功能细化

### 9.1 比赛 AI 看点

第一优先级。

输入：

```json
{
  "match_id": 1
}
```

输出：

```json
{
  "summary": "一句话看点",
  "watch_rating": 4,
  "reasons": [],
  "team_comparison": [],
  "beginner_tips": [],
  "qualification_impact": "",
  "should_stay_up": "",
  "suitable_for": []
}
```

要求：

- 只基于比赛、球队、小组、积分榜数据。
- 没有积分时明确说明。
- 不编造伤病、新闻、球员状态。
- 不输出投注建议。

### 9.2 今日 AI 推荐

第二优先级。

输入：

```json
{
  "date": "2026-06-11",
  "timezone": "Asia/Shanghai",
  "limit": 3
}
```

要求：

- 当天无比赛时返回下一比赛日。
- 已登录用户优先考虑关注球队。
- 推荐理由尽量短。
- 输出“如果只看一场，推荐哪场”。

### 9.3 小组出线分析

第三优先级。

输入：

```json
{
  "group_id": 1
}
```

要求：

- 严格基于 `group_standings`。
- 不生成概率。
- 小组赛未开始时输出赛前说明。
- 对小白解释前二和最佳第三名规则。

### 9.4 AI 聊天助手

第四优先级。

能力：

- 支持自然语言问赛程、球队、规则、观赛建议。
- 支持 match、team、group、general 四类上下文。
- 登录用户保存聊天记录。
- 自用版可以不支持游客。
- 第一版可以非流式返回。

### 9.5 分享文案

第五优先级。

能力：

- 生成朋友圈、微信群、小红书、微博、通用文案。
- 支持轻松、热血、专业、小白友好等语气。
- 支持 short、medium、long。

要求：

- 不编造事实。
- 不用侵权媒体内容。
- 不涉及博彩。

## 10. AI Provider

### 10.1 配置

```env
AI_PROVIDER=mock
AI_BASE_URL=https://api.openai.com/v1
AI_API_KEY=
AI_MODEL=
AI_TIMEOUT_SECONDS=60
AI_DAILY_LIMIT_USER=50
AI_TEMPERATURE=0.7
AI_MAX_TOKENS=1200
```

自用版可以不配置游客额度。

### 10.2 Mock Provider

必须实现。

用途：

- 本地调试。
- 无 API Key 时开发前端。
- 后端测试。

### 10.3 OpenAI-compatible Provider

必须实现。

要求：

- 支持自定义 Base URL。
- 支持自定义模型名。
- 支持超时。
- 错误不暴露 API Key。
- usage log 记录成功或失败。

## 11. 限流与缓存

### 11.1 自用版限流

自用版不需要复杂游客限流，建议只做登录用户每日限制：

```text
ai:limit:user:{userId}:{date}
```

默认：

```env
AI_DAILY_LIMIT_USER=50
```

如果只在本地使用，也可以先不启用强限流，但保留代码结构。

### 11.2 自用版缓存

建议先缓存这些内容：

```text
ai:match_insight:{matchId}:{lang}
ai:today_recommendations:{date}:{timezone}:{userId}
ai:group_analysis:{groupId}
```

缓存时间：

- 比赛看点：6 小时。
- 今日推荐：30 分钟。
- 小组分析：15 分钟。

## 12. Prompt 与安全规则

### 12.1 统一人设

> 你是一个懂球但不装专业的世界杯看球助手。你主要服务项目主人自用，要简洁、清楚、有判断力。你不能编造事实，不能提供投注建议。数据不足时要直接说明。

### 12.2 输出原则

- 先给结论。
- 再给 3 到 5 条理由。
- 尽量短，不写长篇论文。
- 对足球小白友好。
- 涉及预测必须克制。
- 涉及数据必须基于上下文。
- 不说“稳赢”“必胜”“100% 晋级”。

### 12.3 禁止内容

- 非法直播源。
- 真钱下注建议。
- 博彩诱导。
- 编造比分、赛程、伤病、新闻。
- 编造球员事件。
- 无根据的确定性预测。

## 13. 后台管理取舍

自用版第一版不做完整 AI 后台。

保留最小可观测能力：

- usage log 入库。
- 可以通过数据库查看调用记录。
- 后续需要时再做 `/admin/ai` 页面。

如果后续确实需要页面化查看，再补：

```text
GET /api/admin/ai/usage-logs
GET /api/admin/ai/generated-contents
```

## 14. 开发阶段计划

### 阶段 1：数据治理与基线

目标：

- 清理影响 AI 使用的乱码数据。
- 确认当前项目构建状态。

验收：

- 前端核心展示无乱码。
- AI 上下文无明显乱码。
- `go test ./...` 有结果。
- `npm run build` 有结果。

### 阶段 2：后端 AI 最小底座

目标：

- 能用 Mock Provider 生成比赛看点。

任务：

- 新增 AI 配置。
- 新增 AI 模型。
- 新增 Provider 接口。
- 实现 Mock Provider。
- 实现 usage log。
- 实现 MatchContext。
- 实现 `/api/ai/match-insight`。

验收：

- `AI_PROVIDER=mock` 可用。
- match-insight 可返回结构化结果。
- usage log 可入库。

### 阶段 3：真实模型接入

目标：

- 配置真实 API 后可正常调用。

任务：

- 实现 OpenAI-compatible Provider。
- 增加超时和错误处理。
- 增加 Safety 后处理。
- 增加 Markdown/JSON 解析兜底。

验收：

- Mock 和真实 Provider 可切换。
- API Key 不出现在前端或日志。
- 模型失败时有友好提示。

### 阶段 4：前端比赛看点闭环

目标：

- 用户在比赛详情页看到 AI 看点。

任务：

- 新增 AI 类型。
- 新增 AI API 封装。
- 新增 useAIStore。
- 新增 MatchInsightCard。
- 比赛详情页接入。

验收：

- 比赛详情页可生成 AI 看点。
- 接口失败不影响比赛基础信息。
- 前端构建通过。

### 阶段 5：扩展今日推荐和小组分析

目标：

- AI 能回答“今天看什么”和“小组怎么出线”。

任务：

- BuildTodayMatchesContext。
- BuildGroupContext。
- `/api/ai/today-recommendations`。
- `/api/ai/group-analysis`。
- AI 首页。
- 积分榜页 AI 解读入口。

验收：

- `/ai` 能展示今日推荐。
- 积分榜能展示小组分析。

### 阶段 6：聊天和分享文案

目标：

- 自用体验完整。

任务：

- `/api/ai/chat`。
- 会话和消息保存。
- AIChatPage。
- `/api/ai/share-copy`。
- AIShareCopyPage。

验收：

- 可以聊天。
- 可以生成并复制分享文案。
- 聊天记录可保存。

### 阶段 7：文档与轻量优化

目标：

- 自己后续维护方便。

任务：

- README 增加 AI 配置说明。
- 增加 Mock Provider 使用说明。
- 增加常见问题。
- 补充关键后端测试。
- 前端构建验证。

验收：

- 无 API Key 也能用 Mock 跑通。
- 配真实 API 后能调用。
- 文档能指导自己下次启动。

## 15. 自用版最终验收标准

完成后应满足：

- 后端能启动。
- 前端能构建。
- `AI_PROVIDER=mock` 可用。
- 真实 OpenAI-compatible API 可用。
- 比赛详情页能生成 AI 看点。
- `/ai` 能生成今日推荐。
- 积分榜能生成小组出线分析。
- `/ai/chat` 能聊天。
- 分享文案可生成并复制。
- AI 调用日志可入库。
- AI API Key 不暴露到前端。
- 数据不足时 AI 明确说明。

## 16. 推荐第一条开发任务

建议先实现最小闭环：

1. 修复 AI 上下文会读取的乱码字段。
2. 新增 AI 配置。
3. 新增 AI 数据模型。
4. 实现 Mock Provider。
5. 实现 MatchContext。
6. 实现 `/api/ai/match-insight`。
7. 在比赛详情页接入 MatchInsightCard。

这条路线最适合自用版，因为它最短、最稳，也最容易看到实际效果。
