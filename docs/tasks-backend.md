# 开发任务清单

## 首发阵容提醒与赛后补看摘要 — 后端开发

- [x] 1. 扩展 Notification 模型（增加 target_type, target_id）
- [x] 2. 新增 UserMatchEventLog 模型
- [x] 3. 新增 user_match_event_repo.go 去重仓储
- [x] 4. 扩展 favorite_repo / reminder_repo 查询目标用户方法
- [x] 5. 新增 Config 配置项
- [x] 6. 实现首发阵容提醒 service (lineup_alert_service.go)
- [x] 7. 实现首发阵容提醒 job (lineup_alert_scanner.go)
- [x] 8. 实现赛后摘要 service 方法 (ai_service.go GeneratePostMatchSummary)
- [x] 9. 实现赛后摘要 handler + router
- [x] 10. 实现赛后摘要自动生成 job
- [x] 11. 更新 main.go (AutoMigrate + 启动 jobs)
- [x] 12. 构建验证 ✅ `go build ./...` + `go vet ./...` 通过
