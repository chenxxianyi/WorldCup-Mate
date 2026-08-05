<script setup lang="ts">
import { onBeforeUnmount, watch } from 'vue'
import { type CompetitionCode } from '@/data/leagueTheme'
import { useLeagueThemeStore } from '@/stores/useLeagueThemeStore'
import ThemeIcon from './ThemeIcon.vue'

const theme = useLeagueThemeStore()

function close() {
  theme.competitionDialogOpen = false
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') close()
}

watch(() => theme.competitionDialogOpen, (open) => {
  document.body.style.overflow = open ? 'hidden' : ''
  if (open) {
    document.addEventListener('keydown', onKeydown)
  } else {
    document.removeEventListener('keydown', onKeydown)
  }
})

onBeforeUnmount(() => {
  document.body.style.overflow = ''
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div
    v-if="theme.competitionDialogOpen"
    id="competition-dialog"
    class="dialog-layer"
  >
    <button
      class="dialog-backdrop"
      type="button"
      aria-label="关闭赛事选择"
      @click="close"
    />
    <section
      class="competition-sheet"
      role="dialog"
      aria-modal="true"
      aria-labelledby="competition-title"
    >
      <div
        class="sheet-grabber"
        aria-hidden="true"
      />
      <div class="sheet-heading">
        <div>
          <p class="eyebrow">
            CHOOSE YOUR WORLD
          </p><h2 id="competition-title">
            切换赛事
          </h2>
        </div>
        <button
          class="text-button"
          type="button"
          @click="close"
        >
          完成
        </button>
      </div>
      <div class="competition-grid">
        <button
          v-for="code in theme.competitionCodes"
          :key="code"
          class="competition-option"
          :class="{ active: code === theme.currentCode }"
          type="button"
          :data-preview="theme.themeFor(code).slug"
          :aria-pressed="code === theme.currentCode"
          @click="theme.setCompetition(code as CompetitionCode)"
        >
          <span
            v-if="code === theme.currentCode"
            class="option-check"
          ><ThemeIcon name="check" /></span>
          <span class="option-mark">{{ theme.themeFor(code).mark }}</span>
          <span><strong>{{ theme.themeFor(code).name }}</strong><small>{{ theme.themeFor(code).en }} · {{ theme.themeFor(code).season }}</small></span>
        </button>
      </div>
    </section>
  </div>
</template>
