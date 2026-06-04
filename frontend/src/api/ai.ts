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

export function apiGenerateMatchInsight(payload: MatchInsightRequest) {
  return request.post('/ai/match-insight', payload) as Promise<MatchInsight>
}

export function apiGenerateTodayRecommendations(payload: TodayRecommendationRequest) {
  return request.post('/ai/today-recommendations', payload) as Promise<TodayRecommendations>
}

export function apiGenerateGroupAnalysis(payload: GroupAnalysisRequest) {
  return request.post('/ai/group-analysis', payload) as Promise<GroupAnalysis>
}

export function apiExplainFootball(payload: ExplainRequest) {
  return request.post('/ai/explain', payload) as Promise<ExplainResult>
}

export function apiGenerateShareCopy(payload: ShareCopyRequest) {
  return request.post('/ai/share-copy', payload) as Promise<ShareCopyResult>
}

export function apiSendAIChat(payload: AIChatRequest) {
  return request.post('/ai/chat', payload) as Promise<AIChatResponse>
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
