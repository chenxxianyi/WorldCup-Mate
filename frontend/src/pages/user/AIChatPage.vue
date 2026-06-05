<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AIChatPanel from '@/components/ai/AIChatPanel.vue'
import { useAIStore } from '@/stores/useAIStore'

const route = useRoute()
const router = useRouter()
const ai = useAIStore()
const historyOpen = ref(false)
const historyLoading = ref(false)
const historyError = ref('')

const conversationList = computed(() => ai.conversations)

const prompts = [
  '今晚如果只看一场，应该看哪场？',
  '什么是最佳第三名出线规则？',
  '小组赛最后一轮为什么经常很刺激？',
  '给我解释一下越位，别太专业。',
]

function send(message: string) {
  ai.sendChatMessage({
    conversation_id: ai.activeConversation?.id,
    message,
    context_type: 'general',
  }).catch(() => {})
}

function formatConversationTime(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleString('zh-CN', {
    month: 'numeric',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

async function openHistory() {
  historyOpen.value = true
  historyError.value = ''
  historyLoading.value = true
  try {
    await ai.fetchConversations()
  } catch (err: any) {
    historyError.value = err?.message || '读取历史会话失败'
  } finally {
    historyLoading.value = false
  }
}

async function selectConversation(id: number) {
  ai.stopChatGeneration()
  historyError.value = ''
  historyLoading.value = true
  try {
    await ai.fetchConversation(id)
    historyOpen.value = false
  } catch (err: any) {
    historyError.value = err?.message || '打开会话失败'
  } finally {
    historyLoading.value = false
  }
}

function newConversation() {
  ai.startNewConversation()
  historyOpen.value = false
}

function goBack() {
  const canGoBack = Boolean(window.history.state?.back)
  if (canGoBack) {
    router.back()
    return
  }
  router.push('/ai').catch(() => {})
}

onMounted(async () => {
  const q = String(route.query.q || '').trim()
  if (q) {
    const { q: _q, ...query } = route.query
    router.replace({ path: route.path, query }).catch(() => {})
    send(q)
    return
  }
  await ai.restoreLatestConversation().catch(() => {})
})
</script>

<template>
  <div class="chat-page">
    <button class="back-action" type="button" title="返回" aria-label="返回上一页" @click="goBack">
      <span class="material-symbols-outlined" aria-hidden="true">arrow_back</span>
    </button>

    <div class="section-head">
      <div class="head-title">
        <h2>AI 聊天助手</h2>
        <span>{{ ai.activeConversation?.title || '问赛程、球队、规则和观赛建议' }}</span>
      </div>
      <div class="head-actions">
        <button class="chat-action history-action" type="button" @click="openHistory">
          <span class="material-symbols-outlined">history</span>
          <span>历史</span>
        </button>
        <button class="chat-action new-action" type="button" @click="newConversation">
          <span class="material-symbols-outlined">add</span>
          <span>新会话</span>
        </button>
      </div>
    </div>

    <section v-if="historyOpen" class="history-panel" aria-label="历史会话">
      <div class="history-head">
        <strong>历史会话</strong>
        <button class="icon-action" type="button" title="关闭" @click="historyOpen = false">
          <span class="material-symbols-outlined">close</span>
        </button>
      </div>
      <p v-if="historyError" class="history-error">{{ historyError }}</p>
      <div v-if="historyLoading" class="history-state">正在加载...</div>
      <div v-else-if="conversationList.length === 0" class="history-state">还没有历史会话</div>
      <div v-else class="history-list">
        <button
          v-for="item in conversationList"
          :key="item.id"
          class="history-item"
          :class="{ active: ai.activeConversation?.id === item.id }"
          type="button"
          @click="selectConversation(item.id)"
        >
          <span class="history-title">{{ item.title || '新会话' }}</span>
          <span class="history-preview">{{ item.last_message || '暂无回复内容' }}</span>
          <span class="history-time">{{ formatConversationTime(item.updated_at || item.created_at) }}</span>
        </button>
      </div>
    </section>

    <section v-if="!ai.chatMessages.length && !historyOpen" class="quick-panel" aria-label="快捷问题">
      <div class="list-label">你可以这样问</div>
      <div class="quick-list">
        <button
          v-for="item in prompts"
          :key="item"
          class="quick-item"
          type="button"
          @click="send(item)"
        >
          <span class="material-symbols-outlined" aria-hidden="true">chat_bubble</span>
          <span>{{ item }}</span>
        </button>
      </div>
    </section>

    <AIChatPanel
      :messages="ai.chatMessages"
      :loading="ai.chatLoading"
      :error="ai.chatError"
      @send="send"
      @stop="ai.stopChatGeneration"
    />
  </div>
</template>

<style scoped>
.chat-page {
  display: grid;
  gap: 12px;
}

.section-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.head-title {
  min-width: 0;
}

.section-head h2 {
  margin: 0;
  font-size: 19px;
  line-height: 1.3;
}

.head-title > span {
  display: block;
  margin-top: 4px;
  color: var(--muted);
  font-size: 13px;
  line-height: 1.45;
}

.head-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.back-action {
  width: 28px;
  height: 28px;
  display: inline-grid;
  place-items: center;
  justify-self: start;
  margin-left: -12px;
  margin-bottom: -2px;
  border: 0;
  border-radius: 999px;
  color: var(--text);
  background: transparent;
  transition: color 160ms ease-out, background 160ms ease-out, border-color 160ms ease-out, transform 160ms ease-out;
}

.back-action .material-symbols-outlined {
  font-size: 20px;
}

.back-action:hover {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 7%, transparent);
}

.chat-action {
  min-width: 44px;
  height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 0 10px;
  border: 1px solid var(--line);
  border-radius: var(--radius-sm);
  color: var(--muted);
  background: var(--card);
  font-size: 13px;
  font-weight: 700;
  white-space: nowrap;
  transition: color 160ms ease-out, background 160ms ease-out, border-color 160ms ease-out, transform 160ms ease-out;
}

.chat-action .material-symbols-outlined {
  font-size: 19px;
}

.back-action:focus-visible,
.chat-action:focus-visible,
.history-item:focus-visible,
.quick-item:focus-visible {
  outline: 2px solid color-mix(in srgb, var(--blue) 55%, transparent);
  outline-offset: 2px;
}

.back-action:active,
.chat-action:active {
  transform: scale(0.97);
}

.history-action {
  color: var(--text);
}

.new-action {
  color: var(--blue);
  border-color: color-mix(in srgb, var(--blue) 32%, var(--line));
  background: color-mix(in srgb, var(--blue) 7%, var(--card));
}

.history-panel {
  display: grid;
  gap: 0;
  padding: 2px 0 8px;
  border-top: 1px solid var(--line);
  border-bottom: 1px solid var(--line);
}

.history-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  min-height: 44px;
}

.history-head strong {
  font-size: 14px;
}

.history-list {
  display: grid;
  max-height: 280px;
  overflow-y: auto;
}

.history-item {
  position: relative;
  width: 100%;
  display: grid;
  grid-template-columns: 22px minmax(0, 1fr) auto;
  align-items: center;
  column-gap: 10px;
  row-gap: 2px;
  min-height: 58px;
  padding: 9px 4px;
  border: 0;
  border-top: 1px solid var(--line);
  border-radius: 0;
  color: var(--text);
  background: transparent;
  text-align: left;
  transition: color 160ms ease-out, background 160ms ease-out;
}

.history-item::before {
  content: 'chat_bubble';
  grid-row: 1 / span 2;
  grid-column: 1;
  color: var(--muted);
  font-family: 'Material Symbols Outlined';
  font-size: 18px;
  line-height: 1;
  font-variation-settings: 'FILL' 0, 'wght' 300, 'GRAD' 0, 'opsz' 24;
}

.history-item:hover {
  color: var(--blue);
  background: color-mix(in srgb, var(--blue) 4%, transparent);
}

.history-item.active {
  color: var(--blue);
  background: color-mix(in srgb, var(--blue) 7%, transparent);
}

.history-item.active::before {
  color: var(--blue);
}

.history-title {
  display: block;
  grid-column: 2;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 700;
  line-height: 1.55;
}

.history-preview {
  display: -webkit-box;
  grid-column: 2 / 4;
  overflow: hidden;
  color: var(--muted);
  font-size: 12px;
  line-height: 1.45;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.history-time {
  display: block;
  grid-column: 3;
  align-self: start;
  color: var(--weak);
  font-size: 11px;
  line-height: 1.45;
  white-space: nowrap;
}

.history-state,
.history-error {
  margin: 0;
  padding: 12px 0;
  color: var(--muted);
  font-size: 13px;
  text-align: left;
}

.history-error {
  color: var(--blue);
}

.quick-panel {
  display: grid;
  gap: 6px;
  padding: 2px 0 4px;
}

.list-label {
  color: var(--weak);
  font-size: 12px;
  font-weight: 750;
  line-height: 1.4;
}

.quick-list {
  display: grid;
  border-top: 1px solid var(--line);
}

.quick-item {
  width: 100%;
  min-height: 46px;
  display: grid;
  grid-template-columns: 22px minmax(0, 1fr);
  align-items: center;
  gap: 10px;
  padding: 9px 4px;
  border: 0;
  border-bottom: 1px solid var(--line);
  color: var(--text);
  background: transparent;
  text-align: left;
  transition: color 160ms ease-out, background 160ms ease-out;
}

.quick-item:hover {
  color: var(--blue);
  background: color-mix(in srgb, var(--blue) 4%, transparent);
}

.quick-item .material-symbols-outlined {
  color: var(--muted);
  font-size: 18px;
}

.quick-item > span:not(.material-symbols-outlined) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 700;
}

@media (max-width: 520px) {
  .section-head {
    align-items: flex-start;
    flex-direction: row;
  }

  .head-actions {
    align-self: flex-start;
    gap: 4px;
  }

  .back-action {
    width: 28px;
    height: 28px;
  }

  .chat-action {
    width: 44px;
    height: 44px;
    padding: 0;
  }

  .chat-action > span:not(.material-symbols-outlined) {
    position: absolute;
    width: 1px;
    height: 1px;
    overflow: hidden;
    clip: rect(0 0 0 0);
    white-space: nowrap;
  }
}
</style>
