import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
  apiDeleteAIConversation,
  apiGenerateGroupAnalysis,
  apiGenerateMatchInsight,
  apiGenerateShareCopy,
  apiGenerateTodayRecommendations,
  apiGetAIConversation,
  apiListAIConversations,
  apiSendAIChat,
} from '@/api/ai'
import type {
  AIChatMessage,
  AIChatRequest,
  AIConversation,
  GroupAnalysis,
  MatchInsight,
  ShareCopyLength,
  ShareCopyPlatform,
  ShareCopyResult,
  ShareCopyTone,
  TodayRecommendations,
} from '@/types/ai'

function friendlyError(err: unknown) {
  return err instanceof Error && err.message
    ? err.message
    : 'AI 服务暂时不可用，请稍后再试'
}

export const useAIStore = defineStore('ai', () => {
  const currentMatchInsight = ref<MatchInsight | null>(null)
  const matchInsightLoading = ref(false)
  const matchInsightError = ref('')

  const todayRecommendations = ref<TodayRecommendations | null>(null)
  const todayLoading = ref(false)
  const todayError = ref('')

  const groupAnalysisMap = ref<Record<number, GroupAnalysis>>({})
  const groupLoadingMap = ref<Record<number, boolean>>({})
  const groupErrorMap = ref<Record<number, string>>({})

  const conversations = ref<AIConversation[]>([])
  const activeConversation = ref<AIConversation | null>(null)
  const chatMessages = ref<AIChatMessage[]>([])
  const chatLoading = ref(false)
  const chatError = ref('')

  const shareCopyResult = ref<ShareCopyResult | null>(null)
  const shareCopyLoading = ref(false)
  const shareCopyError = ref('')

  const hasMatchInsight = computed(() => !!currentMatchInsight.value)

  async function generateMatchInsight(matchId: number, forceRefresh = false) {
    matchInsightLoading.value = true
    matchInsightError.value = ''
    try {
      currentMatchInsight.value = await apiGenerateMatchInsight({
        match_id: matchId,
        force_refresh: forceRefresh,
      })
      return currentMatchInsight.value
    } catch (err) {
      matchInsightError.value = friendlyError(err)
      throw err
    } finally {
      matchInsightLoading.value = false
    }
  }

  function clearMatchInsight() {
    currentMatchInsight.value = null
    matchInsightError.value = ''
  }

  async function generateTodayRecommendations(timezone: string, forceRefresh = false) {
    todayLoading.value = true
    todayError.value = ''
    try {
      todayRecommendations.value = await apiGenerateTodayRecommendations({
        timezone,
        limit: 3,
        force_refresh: forceRefresh,
      })
      return todayRecommendations.value
    } catch (err) {
      todayError.value = friendlyError(err)
      throw err
    } finally {
      todayLoading.value = false
    }
  }

  async function generateGroupAnalysis(groupId: number, forceRefresh = false) {
    groupLoadingMap.value = { ...groupLoadingMap.value, [groupId]: true }
    groupErrorMap.value = { ...groupErrorMap.value, [groupId]: '' }
    try {
      const result = await apiGenerateGroupAnalysis({
        group_id: groupId,
        force_refresh: forceRefresh,
      })
      groupAnalysisMap.value = { ...groupAnalysisMap.value, [groupId]: result }
      return result
    } catch (err) {
      groupErrorMap.value = { ...groupErrorMap.value, [groupId]: friendlyError(err) }
      throw err
    } finally {
      groupLoadingMap.value = { ...groupLoadingMap.value, [groupId]: false }
    }
  }

  async function generateShareCopy(
    matchId: number,
    platform: ShareCopyPlatform,
    tone: ShareCopyTone,
    length: ShareCopyLength,
  ) {
    shareCopyLoading.value = true
    shareCopyError.value = ''
    try {
      shareCopyResult.value = await apiGenerateShareCopy({
        match_id: matchId,
        platform,
        tone,
        length,
      })
      return shareCopyResult.value
    } catch (err) {
      shareCopyError.value = friendlyError(err)
      throw err
    } finally {
      shareCopyLoading.value = false
    }
  }

  function clearShareCopy() {
    shareCopyResult.value = null
    shareCopyError.value = ''
  }

  async function sendChatMessage(payload: AIChatRequest) {
    const text = payload.message.trim()
    if (!text) return null

    chatLoading.value = true
    chatError.value = ''
    chatMessages.value = [...chatMessages.value, { role: 'user', content: text }]
    try {
      const res = await apiSendAIChat({ ...payload, message: text })
      chatMessages.value = [...chatMessages.value, res.message]
      if (!activeConversation.value || activeConversation.value.id !== res.conversation_id) {
        activeConversation.value = {
          id: res.conversation_id,
          title: text.slice(0, 24),
          context_type: payload.context_type || 'general',
          context_id: payload.context_id,
          last_message: res.message.content,
        }
      }
      return res
    } catch (err) {
      chatError.value = friendlyError(err)
      throw err
    } finally {
      chatLoading.value = false
    }
  }

  async function fetchConversations() {
    conversations.value = await apiListAIConversations()
    return conversations.value
  }

  async function fetchConversation(id: number) {
    activeConversation.value = await apiGetAIConversation(id)
    chatMessages.value = activeConversation.value.messages || []
    return activeConversation.value
  }

  async function deleteConversation(id: number) {
    await apiDeleteAIConversation(id)
    conversations.value = conversations.value.filter((item) => item.id !== id)
    if (activeConversation.value?.id === id) {
      activeConversation.value = null
      chatMessages.value = []
    }
  }

  function startNewConversation() {
    activeConversation.value = null
    chatMessages.value = []
    chatError.value = ''
  }

  return {
    currentMatchInsight,
    matchInsightLoading,
    matchInsightError,
    hasMatchInsight,
    todayRecommendations,
    todayLoading,
    todayError,
    groupAnalysisMap,
    groupLoadingMap,
    groupErrorMap,
    conversations,
    activeConversation,
    chatMessages,
    chatLoading,
    chatError,
    shareCopyResult,
    shareCopyLoading,
    shareCopyError,
    generateMatchInsight,
    clearMatchInsight,
    generateTodayRecommendations,
    generateGroupAnalysis,
    generateShareCopy,
    clearShareCopy,
    sendChatMessage,
    fetchConversations,
    fetchConversation,
    deleteConversation,
    startNewConversation,
  }
})
