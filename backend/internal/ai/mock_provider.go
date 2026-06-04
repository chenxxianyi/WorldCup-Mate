package ai

import (
	"context"
	"strings"
)

type MockProvider struct {
	model string
}

func NewMockProvider(model string) *MockProvider {
	if strings.TrimSpace(model) == "" {
		model = "mock-worldcup-mate"
	}
	return &MockProvider{model: model}
}

func (p *MockProvider) Name() string {
	return "mock"
}

func (p *MockProvider) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	prompt := strings.ToLower(req.UserPrompt)
	content := `{"answer":"我会基于已有赛程和规则回答，不编造数据库里没有的信息。","key_points":["优先看比赛时间和对阵","涉及预测时保持克制","不提供投注建议"]}`

	switch {
	case strings.Contains(prompt, "task:match_insight"):
		content = `{"summary":"这场适合看双方风格和小组形势，不必把它当成确定结果的预测。","watch_rating":4,"reasons":["两队信息已在赛程中确认，观赛门槛低。","小组赛阶段通常会影响后续出线节奏。","如果关注其中一队，这场能帮助判断状态和战术取向。"],"team_comparison":["主队更值得观察开场节奏。","客队的转换和定位球值得留意。"],"beginner_tips":["先看谁能持续控球。","留意上半场前 15 分钟的压迫强度。","比分未出现前，不要急着判断强弱。"],"qualification_impact":"如属于小组赛，本场结果会影响积分排序；具体影响以赛后积分榜为准。","should_stay_up":"如果时间不太晚，值得看；太晚的话可以看集锦和关键片段。","suitable_for":["关注参赛队的球迷","想了解小组形势的新球迷","喜欢赛前看重点的人"]}`
	case strings.Contains(prompt, "task:today_recommendations"):
		content = `{"date":"","timezone":"","recommendations":[{"match_id":0,"title":"今日重点比赛","kickoff_time":"","reason":"优先选择重要程度高、对阵清晰、适合轻松观看的比赛。","rating":4}],"only_one_match":null,"note":"这是 mock 结果；真实排序会基于数据库中的今日赛程。"}`
	case strings.Contains(prompt, "task:group_analysis"):
		content = `{"summary":"这个小组的出线判断应以当前积分榜为准，暂不做概率推断。","key_points":["先看积分，再看净胜球和进球数。","前两名通常最稳，第三名需要横向比较。","未赛比赛越多，结论越需要保守。"],"qualification_rules":"小组前两名晋级，部分成绩较好的第三名晋级；最终规则以赛事官方为准。","teams":[],"data_note":"mock 分析会等待 service 补入数据库球队和积分信息。"}`
	case strings.Contains(prompt, "task:share_copy"):
		content = `{"title":"今晚看球提醒","content":"这场比赛值得关注双方节奏和小组形势。约上朋友一起看，重点留意开局压迫、定位球和下半场换人。","tips":["不要写确定比分。","不要写未确认消息。","分享时带上开球时间更友好。"]}`
	case strings.Contains(prompt, "task:chat"):
		content = "我会按现有赛程和积分信息帮你梳理。你可以问我某场比赛值不值得看、某个小组形势，或者让它生成一段分享文案。"
	}

	return &ChatResponse{
		Content:          content,
		Model:            p.model,
		PromptTokens:     estimateTokens(req.SystemPrompt) + estimateTokens(req.UserPrompt),
		CompletionTokens: estimateTokens(content),
		TotalTokens:      estimateTokens(req.SystemPrompt) + estimateTokens(req.UserPrompt) + estimateTokens(content),
	}, nil
}

func estimateTokens(s string) int {
	n := len([]rune(s)) / 2
	if n < 1 {
		return 1
	}
	return n
}
