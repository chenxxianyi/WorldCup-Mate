import request from './request'
import type {
  AIChatRequest,
  AIChatResponse,
  AIChatStreamEvent,
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

export async function apiSendAIChatStream(
  payload: AIChatRequest,
  onEvent: (event: AIChatStreamEvent) => void,
  signal?: AbortSignal,
): Promise<AIChatResponse> {
  const token = localStorage.getItem('wm-token') || sessionStorage.getItem('wm-token')
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    Accept: 'text/event-stream',
  }
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch('/api/ai/chat/stream', {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
    signal,
  })
  if (!res.ok) {
    throw new Error('AI 服务暂时不可用，请稍后再试')
  }
  if (!res.body) {
    throw new Error('当前浏览器不支持流式响应')
  }

  const reader = res.body.getReader()
  const decoder = new TextDecoder()
  let buffer = ''
  let finalResponse: AIChatResponse | null = null

  function handleBlock(block: string) {
    const data = block
      .split(/\r?\n/)
      .filter((line) => line.startsWith('data:'))
      .map((line) => line.slice(5).trim())
      .join('\n')
    if (!data) return

    const event = JSON.parse(data) as AIChatStreamEvent
    onEvent(event)
    if (event.type === 'error') {
      throw new Error(event.delta || 'AI 服务暂时不可用，请稍后再试')
    }
    if (event.type === 'done' && event.message) {
      finalResponse = {
        conversation_id: event.conversation_id || 0,
        message: event.message,
      }
    }
  }

  while (true) {
    const { value, done } = await reader.read()
    buffer += decoder.decode(value || new Uint8Array(), { stream: !done })
    const blocks = buffer.split(/\r?\n\r?\n/)
    buffer = blocks.pop() || ''
    for (const block of blocks) {
      handleBlock(block)
    }
    if (done) break
  }
  if (buffer.trim()) {
    handleBlock(buffer)
  }
  if (!finalResponse) {
    throw new Error('AI 回复中断，请稍后重试')
  }
  return finalResponse as AIChatResponse
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
