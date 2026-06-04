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
  apiSendAIChatStream,
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

const ACTIVE_CONVERSATION_KEY = 'wm-ai-active-conversation'

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

  function rememberActiveConversation(id?: number | null) {
    if (id && id > 0) {
      localStorage.setItem(ACTIVE_CONVERSATION_KEY, String(id))
    } else {
      localStorage.removeItem(ACTIVE_CONVERSATION_KEY)
    }
  }

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
    chatMessages.value = [
      ...chatMessages.value,
      { role: 'user', content: text },
      { role: 'assistant', content: '' },
    ]
    const assistantIndex = chatMessages.value.length - 1
    let conversationId = payload.conversation_id || activeConversation.value?.id || 0

    function updateAssistant(updater: (message: AIChatMessage) => AIChatMessage) {
      const next = [...chatMessages.value]
      const current = next[assistantIndex]
      if (!current || current.role !== 'assistant') return
      next[assistantIndex] = updater(current)
      chatMessages.value = next
    }

    try {
      const res = await apiSendAIChatStream({ ...payload, message: text }, (event) => {
        if (event.type === 'start' && event.conversation_id) {
          conversationId = event.conversation_id
          rememberActiveConversation(conversationId)
          if (!activeConversation.value || activeConversation.value.id !== conversationId) {
            activeConversation.value = {
              id: conversationId,
              title: text.slice(0, 24),
              context_type: payload.context_type || 'general',
              context_id: payload.context_id,
              last_message: '',
            }
          }
        }

        if (event.type === 'delta' && event.delta) {
          updateAssistant((message) => ({ ...message, content: message.content + event.delta }))
        }

        if (event.type === 'done' && event.message) {
          updateAssistant(() => event.message as AIChatMessage)
          if (activeConversation.value && event.conversation_id) {
            activeConversation.value.last_message = event.message.content
          }
        }
      })
      updateAssistant(() => res.message)
      if (!activeConversation.value || activeConversation.value.id !== res.conversation_id) {
        activeConversation.value = {
          id: res.conversation_id,
          title: text.slice(0, 24),
          context_type: payload.context_type || 'general',
          context_id: payload.context_id,
          last_message: res.message.content,
        }
      } else {
        activeConversation.value.last_message = res.message.content
      }
      rememberActiveConversation(res.conversation_id)
      return res
    } catch (err) {
      chatError.value = friendlyError(err)
      updateAssistant((message) => (
        message.content
          ? message
          : { ...message, content: 'AI 回复失败，请稍后再试。' }
      ))
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
    rememberActiveConversation(activeConversation.value.id)
    return activeConversation.value
  }

  async function restoreLatestConversation() {
    if (chatMessages.value.length > 0 || chatLoading.value) {
      return activeConversation.value
    }

    chatError.value = ''
    const savedID = Number(localStorage.getItem(ACTIVE_CONVERSATION_KEY) || 0)
    if (savedID > 0) {
      try {
        return await fetchConversation(savedID)
      } catch {
        rememberActiveConversation(null)
      }
    }

    const list = await fetchConversations()
    if (list.length === 0) {
      return null
    }
    return fetchConversation(list[0].id)
  }

  async function deleteConversation(id: number) {
    await apiDeleteAIConversation(id)
    conversations.value = conversations.value.filter((item) => item.id !== id)
    if (activeConversation.value?.id === id) {
      activeConversation.value = null
      chatMessages.value = []
      rememberActiveConversation(null)
    }
  }

  function startNewConversation() {
    activeConversation.value = null
    chatMessages.value = []
    chatError.value = ''
    rememberActiveConversation(null)
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
    restoreLatestConversation,
    deleteConversation,
    startNewConversation,
  }
})
