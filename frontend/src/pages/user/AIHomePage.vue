<script setup lang="ts">
import { useRouter } from 'vue-router'
import PromptSuggestionCard from '@/components/ai/PromptSuggestionCard.vue'
import TodayRecommendationCard from '@/components/ai/TodayRecommendationCard.vue'
import { useAIStore } from '@/stores/useAIStore'
import { useSettingStore } from '@/stores/useSettingStore'

const router = useRouter()
const ai = useAIStore()
const settings = useSettingStore()

function generate(forceRefresh = false) {
  ai.generateTodayRecommendations(settings.timezone, forceRefresh).catch(() => {})
}

function openChat(prompt?: string) {
  router.push({ path: '/ai/chat', query: prompt ? { q: prompt } : undefined })
}
</script>

<template>
  <div class="ai-home">
    <section class="hero-card">
      <span class="eyebrow">AI 助手</span>
      <h2>今天看球，更省心</h2>
      <p>帮你挑比赛、读小组形势、解释规则，也能顺手写一段分享文案。</p>
      <div class="hero-actions">
        <button class="pill-btn primary" type="button" @click="generate(false)">生成今日推荐</button>
        <button class="pill-btn" type="button" @click="router.push('/ai/chat')">打开聊天</button>
      </div>
    </section>

    <TodayRecommendationCard
      :data="ai.todayRecommendations"
      :loading="ai.todayLoading"
      :error="ai.todayError"
      @generate="generate(false)"
      @refresh="generate(true)"
    />

    <section class="section">
      <div class="section-head">
        <h2>快捷入口</h2>
        <span>常用看球问题</span>
      </div>
      <div class="prompt-grid">
        <PromptSuggestionCard
          title="今晚只看一场"
          description="让 AI 从今日赛程里挑重点。"
          icon="sports_soccer"
          @select="generate(false)"
        />
        <PromptSuggestionCard
          title="解释一条规则"
          description="越位、点球大战、最佳第三名都能问。"
          icon="help"
          @select="openChat('什么是越位？用小白能懂的话解释一下')"
        />
        <PromptSuggestionCard
          title="写分享文案"
          description="生成朋友圈、群聊或微博文案。"
          icon="edit_note"
          @select="router.push('/ai/share-copy')"
        />
        <PromptSuggestionCard
          title="分析小组形势"
          description="去积分榜查看每组 AI 解读。"
          icon="table_chart"
          @select="router.push('/standings')"
        />
      </div>
    </section>
  </div>
</template>

<style scoped>
.ai-home {
  display: grid;
  gap: 14px;
}

.hero-card {
  display: grid;
  gap: 10px;
  padding: 10px 2px 4px;
}

.eyebrow {
  color: var(--primary);
  font-size: 12px;
  font-weight: 800;
}

.hero-card h2 {
  margin: 0;
  max-width: 680px;
  font-size: clamp(25px, 5vw, 34px);
  line-height: 1.16;
  letter-spacing: 0;
}

.hero-card p {
  max-width: 560px;
  margin: 0;
  color: var(--muted);
  line-height: 1.55;
  font-size: 15px;
}

.hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-top: 2px;
}

.section {
  display: grid;
  gap: 10px;
  padding-top: 2px;
}

.section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.section-head h2 {
  margin: 0;
  font-size: 17px;
}

.section-head span {
  color: var(--muted);
  font-size: 13px;
}

.prompt-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

@media (min-width: 520px) {
  .prompt-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (min-width: 768px) {
  .prompt-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}
</style>
