import request from './request'
import type {
  AIChatRequest,
  AIChatResponse,
  AIConversation,
  ExplainRequest,
  ExplainResult,
  GroupAnalysis,
  GroupAnalysisRequest,
  MatchInsight,
  MatchInsightRequest,
  ShareCopyRequest,
  ShareCopyResult,
  TodayRecommendationRequest,
  TodayRecommendations,
} from '@/types/ai'

const aiRequestConfig = {
  timeout: 70000,
}

export function apiGenerateMatchInsight(payload: MatchInsightRequest) {
  return request.post('/ai/match-insight', payload, aiRequestConfig) as Promise<MatchInsight>
}

export function apiGenerateTodayRecommendations(payload: TodayRecommendationRequest) {
  return request.post('/ai/today-recommendations', payload, aiRequestConfig) as Promise<TodayRecommendations>
}

export function apiGenerateGroupAnalysis(payload: GroupAnalysisRequest) {
  return request.post('/ai/group-analysis', payload, aiRequestConfig) as Promise<GroupAnalysis>
}

export function apiExplainFootball(payload: ExplainRequest) {
  return request.post('/ai/explain', payload, aiRequestConfig) as Promise<ExplainResult>
}

export function apiGenerateShareCopy(payload: ShareCopyRequest) {
  return request.post('/ai/share-copy', payload, aiRequestConfig) as Promise<ShareCopyResult>
}

export function apiSendAIChat(payload: AIChatRequest) {
  return request.post('/ai/chat', payload, aiRequestConfig) as Promise<AIChatResponse>
}

export function apiListAIConversations() {
  return request.get('/ai/conversations') as Promise<AIConversation[]>
}

export function apiGetAIConversation(id: number) {
  return request.get(`/ai/conversations/${id}`) as Promise<AIConversation>
}

export function apiDeleteAIConversation(id: number) {
  return request.delete(`/ai/conversations/${id}`) as Promise<void>
}
