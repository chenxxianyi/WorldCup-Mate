export interface MatchInsightRequest {
  match_id: number
  force_refresh?: boolean
}

export interface MatchInsight {
  summary: string
  watch_rating: number
  reasons: string[]
  team_comparison: string[]
  beginner_tips: string[]
  qualification_impact: string
  should_stay_up: string
  suitable_for: string[]
  generated_at?: string
}

export interface TodayRecommendationRequest {
  date?: string
  timezone?: string
  limit?: number
  force_refresh?: boolean
}

export interface TodayRecommendedMatch {
  match_id: number
  title: string
  kickoff_time?: string
  reason: string
  rating?: number
}

export interface TodayRecommendations {
  date: string
  timezone: string
  recommendations: TodayRecommendedMatch[]
  only_one_match?: TodayRecommendedMatch | null
  note?: string
}

export interface GroupAnalysisRequest {
  group_id: number
  force_refresh?: boolean
}

export interface GroupAnalysisTeam {
  team_id?: number
  team_name: string
  status?: string
  note?: string
}

export interface GroupAnalysis {
  summary: string
  key_points: string[]
  qualification_rules: string
  teams: GroupAnalysisTeam[]
  data_note?: string
  generated_at?: string
}

export interface ExplainRequest {
  question: string
  context_type?: AIContextType
  context_id?: number
}

export interface ExplainResult {
  answer: string
  key_points?: string[]
}

export interface ShareCopyRequest {
  match_id: number
  platform: ShareCopyPlatform
  tone: ShareCopyTone
  length: ShareCopyLength
}

export interface ShareCopyResult {
  title?: string
  content: string
  tips?: string[]
}

export type ShareCopyPlatform = 'wechat' | 'group' | 'xiaohongshu' | 'weibo' | 'general'
export type ShareCopyTone = 'relaxed' | 'passionate' | 'professional' | 'beginner'
export type ShareCopyLength = 'short' | 'medium' | 'long'
export type AIContextType = 'general' | 'match' | 'team' | 'group'

export interface AIChatRequest {
  conversation_id?: number | null
  message: string
  context_type?: AIContextType
  context_id?: number
}

export interface AIChatMessage {
  id?: number
  role: 'user' | 'assistant' | 'system'
  content: string
  created_at?: string
}

export interface AIChatResponse {
  conversation_id: number
  message: AIChatMessage
}

export interface AIChatStreamEvent {
  type: 'start' | 'delta' | 'done' | 'error'
  conversation_id?: number
  delta?: string
  message?: AIChatMessage
}

export interface AIConversation {
  id: number
  title: string
  context_type: AIContextType
  context_id?: number
  last_message?: string
  messages?: AIChatMessage[]
  created_at?: string
  updated_at?: string
}
