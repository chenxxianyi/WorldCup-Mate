<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AIChatPanel from '@/components/ai/AIChatPanel.vue'
import PromptSuggestionCard from '@/components/ai/PromptSuggestionCard.vue'
import { useAIStore } from '@/stores/useAIStore'

const route = useRoute()
const ai = useAIStore()

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

onMounted(() => {
  const q = String(route.query.q || '').trim()
  if (q) send(q)
})
</script>

<template>
  <div class="chat-page">
    <div class="section-head">
      <div>
        <h2>AI 聊天助手</h2>
        <span>问赛程、球队、规则和观赛建议</span>
      </div>
      <button class="pill-btn" type="button" @click="ai.startNewConversation">新会话</button>
    </div>

    <div class="prompt-grid">
      <PromptSuggestionCard
        v-for="item in prompts"
        :key="item"
        :title="item"
        icon="chat"
        @select="send(item)"
      />
    </div>

    <AIChatPanel
      :messages="ai.chatMessages"
      :loading="ai.chatLoading"
      :error="ai.chatError"
      @send="send"
    />
  </div>
</template>

<style scoped>
.chat-page {
  display: grid;
  gap: 14px;
}

.section-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 20px;
}

.section-head span {
  display: block;
  margin-top: 5px;
  color: var(--muted);
  font-size: 13px;
}

.prompt-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

@media (min-width: 900px) {
  .prompt-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}
</style>
