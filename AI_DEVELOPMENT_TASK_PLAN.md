# WorldCup Mate V2 自用版 AI 看球助手任务计划

## 1. 使用说明

本文档基于 `AI_DEVELOPMENT_PLAN.md` 的自用版方案拆分任务。

自用版目标是快速做出个人可用的 AI 看球助手，不优先考虑商业化、公开注册、复杂后台运营、付费体系或大规模用户并发。

开发原则：

- 先跑通最小闭环，再扩展页面和能力。
- 先 Mock Provider，再真实模型。
- 先比赛详情页 AI 看点，再做 AI 首页和聊天。
- 后台 AI 管理、反馈统计、复杂 Token 报表全部后置。
- 每个阶段完成后运行必要验证，避免大步堆功能。

## 2. 阶段总览

| 阶段 | 名称 | 目标 | 优先级 |
| --- | --- | --- | --- |
| T0 | 基线与数据治理 | 让 AI 可读取可信中文数据 | P0 |
| T1 | 后端 AI 最小底座 | 配置、模型、Mock Provider、日志 | P0 |
| T2 | 比赛看点闭环 | 后端 match-insight + 前端详情页卡片 | P0 |
| T3 | 真实模型接入 | OpenAI-compatible Provider 和安全兜底 | P0 |
| T4 | 今日推荐与小组分析 | AI 首页和积分榜 AI 解读 | P1 |
| T5 | 聊天助手与分享文案 | 自用完整体验 | P1 |
| T6 | 轻量缓存、限流与文档 | 降低成本，方便维护 | P1 |
| T7 | 可选增强 | 观赛计划、后台查看、反馈等 | P2 |

## 3. T0 基线与数据治理

### T0.1 确认 V2 分支

目标：

- 后续改动都进入 `V2` 分支。

建议命令：

```powershell
git status
git switch -c V2
git push -u origin V2
```

如果已经存在远程 `V2`：

```powershell
git fetch origin
git switch V2
git pull origin V2
```

验收标准：

- 当前分支为 `V2`。
- `.env`、`backend/.env` 不被提交。
- 工作区状态明确。

### T0.2 基线检查

目标：

- 记录当前项目能否正常构建和测试。

执行：

```powershell
cd backend
go test ./...
```

```powershell
cd frontend
npm run build
```

验收标准：

- 后端测试结果有记录。
- 前端构建结果有记录。
- 已知失败项明确归因，避免和 AI 改动混在一起。

### T0.3 AI 核心数据排查

目标：

- 确认 AI 上下文会用到的数据是否存在乱码。

重点检查表：

```text
teams
cities
stadia
matches
groups
group_standings
```

重点字段：

```text
teams.name
teams.name_en
teams.fifa_code
teams.continent
cities.name
cities.name_en
stadia.name
stadia.name_en
matches.recommend_tag
matches.recommend_reason
matches.kickoff_time_utc
matches.status
```

验收标准：

- 明确哪些中文字段乱码。
- 明确是否有重复城市、球场、球队。
- 明确比赛数据是否能支撑 AI 看点。

### T0.4 最小数据清洗

目标：

- 修复 AI 第一阶段必读字段。

建议只先修这些：

- 球队中文名。
- 城市中文名。
- 球场中文名。
- 大洲字段。
- 比赛推荐标签。

不建议第一阶段就做：

- 大规模数据重构。
- 外键合并。
- 重写同步逻辑。
- 复杂数据后台。

验收标准：

- 比赛详情页和球队页无明显乱码。
- MatchContext 不读取乱码字段。
- 不破坏外键关系。

## 4. T1 后端 AI 最小底座

### T1.1 新增 AI 配置

目标：

- 后端支持 `mock` 和真实 Provider 切换。

修改文件：

```text
backend/internal/config/config.go
backend/.env.example
```

新增环境变量：

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

验收标准：

- 不配置 API Key 时默认使用 Mock。
- 所有 AI 配置来自环境变量。
- API Key 不写入代码。

### T1.2 新增 AI 数据模型

目标：

- 支持自用版聊天、生成内容和调用日志。

新增文件：

```text
backend/internal/models/ai_conversation.go
backend/internal/models/ai_message.go
backend/internal/models/ai_generated_content.go
backend/internal/models/ai_usage_log.go
```

修改文件：

```text
backend/cmd/server/main.go
```

验收标准：

- AutoMigrate 包含新模型。
- 启动后可创建表。
- 不新增第一版不需要的 `ai_feedbacks`。

### T1.3 新增 AI Repository

目标：

- 集中管理 AI 入库操作。

新增文件：

```text
backend/internal/repositories/ai_repo.go
```

第一版需要支持：

- 保存 generated content。
- 保存 usage log。
- 创建 conversation。
- 保存 user message。
- 保存 assistant message。
- 查询 conversation 列表。
- 查询 conversation detail。
- 删除 conversation。

验收标准：

- Repository 不调用模型。
- Repository 不拼 Prompt。
- 查询会话支持按用户过滤。

### T1.4 新增 Provider 接口

目标：

- 统一 Mock 和真实模型调用。

新增文件：

```text
backend/internal/ai/provider.go
```

建议结构：

```go
type Provider interface {
    Name() string
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}
```

验收标准：

- Provider 与业务场景解耦。
- ChatRequest 支持 system、user、temperature、max_tokens。
- ChatResponse 支持 content、token usage、model。

### T1.5 实现 Mock Provider

目标：

- 无真实 API Key 时也能开发和测试。

新增文件：

```text
backend/internal/ai/mock_provider.go
```

第一版 Mock 覆盖：

- match-insight。
- today-recommendations。
- group-analysis。
- explain。
- share-copy。
- chat。

验收标准：

- `AI_PROVIDER=mock` 可用。
- 不访问外部网络。
- 返回结构稳定。

### T1.6 新增基础 Prompt 和 Safety

目标：

- 固定 AI 人设和输出边界。

新增文件：

```text
backend/internal/ai/prompts.go
backend/internal/ai/safety.go
```

要求：

- 禁止投注建议。
- 禁止编造比分、伤病、新闻。
- 禁止确定性预测。
- 数据不足时明确提示。

验收标准：

- 所有 AI 接口共用系统 Prompt。
- Safety 可过滤明显博彩诱导词。

## 5. T2 比赛看点闭环

### T2.1 实现 MatchContext

目标：

- 为比赛 AI 看点提供可信上下文。

新增文件：

```text
backend/internal/ai/context_builder.go
```

上下文包含：

- 比赛基础信息。
- 主队、客队。
- 城市、球场。
- 小组。
- 开球时间。
- 比赛状态。
- 比分。
- 小组积分。
- 用户关注球队。
- 用户时区。

验收标准：

- 比赛不存在时返回错误。
- 无积分时明确标记。
- 不包含乱码字段。
- 不包含数据库不存在的事实。

### T2.2 实现 AI Service 最小流程

目标：

- 串起限流、上下文、Provider、解析和日志。

新增文件：

```text
backend/internal/services/ai_service.go
```

第一版流程：

1. 校验 match_id。
2. 构建 MatchContext。
3. 生成 Prompt。
4. 调用 Provider。
5. 解析结果。
6. 保存 generated content。
7. 保存 usage log。
8. 返回前端。

验收标准：

- Mock Provider 下完整跑通。
- 调用失败也写 usage log。
- 返回统一响应格式。

### T2.3 实现 match-insight Handler 和路由

目标：

- 后端暴露比赛 AI 看点接口。

新增或修改：

```text
backend/internal/handlers/ai_handler.go
backend/internal/routes/router.go
```

接口：

```text
POST /api/ai/match-insight
```

请求：

```json
{
  "match_id": 1
}
```

验收标准：

- 未登录也可以本地测试，或登录用户可用。
- 参数错误有友好提示。
- Mock 返回结构化看点。

### T2.4 前端新增 AI 类型和 API

目标：

- 前端能调用 match-insight。

新增文件：

```text
frontend/src/types/ai.ts
frontend/src/api/ai.ts
```

第一版类型：

- MatchInsight。
- MatchInsightRequest。
- MatchInsightResponse。

验收标准：

- 不直接使用模型供应商 API。
- 类型和后端返回字段一致。

### T2.5 新增 useAIStore 最小版

目标：

- 管理比赛 AI 看点状态。

新增文件：

```text
frontend/src/stores/useAIStore.ts
```

第一版 state：

- currentMatchInsight。
- loading。
- error。

第一版 action：

- generateMatchInsight。

验收标准：

- loading 和 error 正常。
- 接口失败不影响页面主体。

### T2.6 新增 MatchInsightCard

目标：

- 结构化展示比赛 AI 看点。

新增文件：

```text
frontend/src/components/ai/MatchInsightCard.vue
```

展示：

- 一句话总结。
- 看点指数。
- 推荐理由。
- 双方对比。
- 小白观赛重点。
- 出线影响。
- 是否值得熬夜。

验收标准：

- 移动端可读。
- 空状态可用。
- Loading 状态可用。

### T2.7 比赛详情页接入 AI 看点

目标：

- 形成第一个完整自用闭环。

修改文件：

```text
frontend/src/pages/user/MatchDetailPage.vue
```

功能：

- 增加 AI 看点区域。
- 点击按钮生成。
- 展示 MatchInsightCard。
- 支持重新生成。

验收标准：

- 比赛详情页能生成 AI 看点。
- AI 接口失败不影响比赛详情。
- `npm run build` 通过。
- `go test ./...` 通过。

## 6. T3 真实模型接入

### T3.1 实现 OpenAI-compatible Provider

目标：

- 支持真实模型调用。

新增文件：

```text
backend/internal/ai/openai_compatible.go
```

要求：

- 支持 `AI_BASE_URL`。
- 支持 `AI_API_KEY`。
- 支持 `AI_MODEL`。
- 支持超时。
- 支持 token usage 解析。
- 错误不暴露 API Key。

验收标准：

- Mock 和真实 Provider 可通过配置切换。
- 真实接口失败有友好错误。
- usage log 可记录失败原因。

### T3.2 实现 Parser 兜底

目标：

- 避免模型返回非标准 JSON 时前端崩掉。

新增文件：

```text
backend/internal/ai/parser.go
```

要求：

- 优先解析 JSON。
- JSON 解析失败时返回 Markdown 兜底。
- 必要字段缺失时填默认值。

验收标准：

- 模型输出轻微偏格式时仍可返回。
- 前端不会因为解析失败白屏。

### T3.3 验证真实模型调用

目标：

- 确认自用环境可以调用真实 AI。

验证内容：

- `.env` 配置真实 API。
- 调用 `/api/ai/match-insight`。
- 检查 usage log。
- 检查前端展示。

验收标准：

- 配置真实 API 后可用。
- API Key 不出现在前端构建产物。
- API Key 不出现在日志。

## 7. T4 今日推荐与小组分析

### T4.1 BuildTodayMatchesContext

目标：

- 为今日推荐提供候选比赛。

上下文包含：

- 指定日期比赛。
- 用户时区。
- 关注球队。
- 比赛阶段。
- 比赛重要度。
- 比赛状态。

验收标准：

- 当天无比赛时能找到下一比赛日。
- 只推荐数据库存在的比赛。

### T4.2 实现 today-recommendations API

目标：

- 生成“今天看什么”。

接口：

```text
POST /api/ai/today-recommendations
```

验收标准：

- 返回最多 3 场推荐。
- 返回“只看一场”的推荐。
- Mock 和真实 Provider 都可用。

### T4.3 BuildGroupContext

目标：

- 为小组出线分析提供上下文。

上下文包含：

- 小组球队。
- 当前积分。
- 已结束比赛。
- 未开始比赛。
- 最佳第三名规则。

验收标准：

- 小组无积分时明确说明。
- 不生成虚假概率。

### T4.4 实现 group-analysis API

目标：

- 生成小组出线解读。

接口：

```text
POST /api/ai/group-analysis
```

验收标准：

- 输入 group_id。
- 基于 `group_standings`。
- 输出小白可读。

### T4.5 新增 AI 首页

目标：

- 提供自用 AI 入口。

新增文件：

```text
frontend/src/pages/user/AIHomePage.vue
```

修改：

```text
frontend/src/router/index.ts
```

功能：

- 今日推荐。
- 快捷问题。
- 跳转聊天。
- 跳转比赛详情。

验收标准：

- `/ai` 可打开。
- 今日推荐可生成。
- 移动端可用。

### T4.6 积分榜页接入 AI 解读

目标：

- 在积分榜页查看小组出线分析。

修改文件：

```text
frontend/src/pages/user/StandingsPage.vue
```

验收标准：

- 每个小组可触发 AI 解读。
- 接口失败不影响积分榜。

## 8. T5 聊天助手与分享文案

### T5.1 实现 chat API

目标：

- 支持自用聊天。

接口：

```text
POST /api/ai/chat
GET /api/ai/conversations
GET /api/ai/conversations/:id
DELETE /api/ai/conversations/:id
```

验收标准：

- 登录用户保存聊天记录。
- 只带最近 N 条历史。
- 可按 context_type 构建上下文。
- 第一版非流式即可。

### T5.2 新增 AIChatPage

目标：

- 前端可聊天。

新增文件：

```text
frontend/src/pages/user/AIChatPage.vue
frontend/src/components/ai/AIChatPanel.vue
frontend/src/components/ai/AIMessageBubble.vue
frontend/src/components/ai/AIInputBox.vue
frontend/src/components/ai/AIThinking.vue
```

验收标准：

- 支持发送消息。
- 支持 Markdown 展示。
- 支持复制回答。
- Loading 状态清楚。

### T5.3 实现 explain API

目标：

- 解释足球规则和术语。

接口：

```text
POST /api/ai/explain
```

验收标准：

- 可解释越位、高位逼抢、补时、点球大战等。
- 涉及 2026 赛制时使用内置规则。

### T5.4 实现 share-copy API

目标：

- 生成社交分享文案。

接口：

```text
POST /api/ai/share-copy
```

验收标准：

- 支持平台、语气、长度。
- 不编造事实。
- 不涉及博彩。

### T5.5 新增分享文案页面

目标：

- 前端可选择比赛并生成文案。

新增文件：

```text
frontend/src/pages/user/AIShareCopyPage.vue
frontend/src/components/ai/ShareCopyCard.vue
```

验收标准：

- 可生成文案。
- 可复制。
- 失败提示友好。

## 9. T6 轻量缓存、限流与文档

### T6.1 登录用户限流

目标：

- 避免自用时误刷模型成本。

Redis Key：

```text
ai:limit:user:{userId}:{date}
```

配置：

```env
AI_DAILY_LIMIT_USER=50
```

验收标准：

- 达到上限返回友好提示。
- Redis 不可用时不导致核心页面崩溃。

### T6.2 轻量缓存

目标：

- 减少重复调用成本。

缓存 Key：

```text
ai:match_insight:{matchId}:{lang}
ai:today_recommendations:{date}:{timezone}:{userId}
ai:group_analysis:{groupId}
```

验收标准：

- 相同比赛短时间重复生成优先读缓存。
- 允许“重新生成”绕过缓存。

### T6.3 README 更新

目标：

- 方便自己下次启动和维护。

修改文件：

```text
README.md
backend/.env.example
```

新增内容：

- 自用版 AI 功能说明。
- Mock Provider 启动方式。
- 真实 Provider 配置方式。
- 常见问题。
- 数据不足时 AI 的行为。
- API Key 安全说明。

验收标准：

- 无 API Key 时可按 README 使用 Mock。
- 有 API Key 时可按 README 接入真实模型。

### T6.4 最终构建检查

执行：

```powershell
cd backend
go test ./...
```

```powershell
cd frontend
npm run build
```

验收标准：

- 后端测试通过或已记录已知问题。
- 前端构建通过。
- 文档已更新。

## 10. T7 可选增强

以下任务不是自用第一版必须项。

### T7.1 观赛计划

接口：

```text
POST /api/ai/viewing-plan
```

页面：

```text
/ai/viewing-plan
```

适合在今日推荐稳定后再做。

### T7.2 赛后总结

接口：

```text
POST /api/ai/post-match-summary
```

只有比赛状态为 `finished` 且有比分时启用。

### T7.3 AI 反馈

接口：

```text
POST /api/ai/feedback
```

自用版可暂缓，除非你想持续改 Prompt。

### T7.4 后台 AI 查看页

可选接口：

```text
GET /api/admin/ai/usage-logs
GET /api/admin/ai/generated-contents
```

可选页面：

```text
frontend/src/pages/admin/AdminAIPage.vue
```

自用版可以先通过数据库查看 usage log，不必马上做后台页面。

## 11. 推荐提交节奏

建议按最小闭环提交：

```text
chore: prepare v2 self-use ai baseline
fix: clean ai context data encoding
feat: add ai config and models
feat: add mock ai provider
feat: add match insight backend api
feat: integrate match insight into match detail
feat: add openai compatible provider
feat: add ai home recommendations
feat: add group ai analysis
feat: add ai chat and share copy
docs: update self-use ai setup guide
```

## 12. 最小闭环任务清单

如果只想最快看到效果，按这个顺序做：

1. T0.3 AI 核心数据排查。
2. T0.4 最小数据清洗。
3. T1.1 新增 AI 配置。
4. T1.2 新增 AI 数据模型。
5. T1.4 新增 Provider 接口。
6. T1.5 实现 Mock Provider。
7. T2.1 实现 MatchContext。
8. T2.2 实现 AI Service 最小流程。
9. T2.3 实现 match-insight Handler 和路由。
10. T2.4 前端新增 AI 类型和 API。
11. T2.5 新增 useAIStore 最小版。
12. T2.6 新增 MatchInsightCard。
13. T2.7 比赛详情页接入 AI 看点。

完成后，自用版 AI 模块已经有第一个稳定闭环。

## 13. 任务完成记录模板

```text
完成时间：

修改文件：

实现内容：

验证命令：

验证结果：

遗留问题：
```
