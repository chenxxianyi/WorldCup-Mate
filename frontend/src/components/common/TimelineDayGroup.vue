<script setup lang="ts">
import TimelineMatchCard from '@/components/common/TimelineMatchCard.vue'

defineProps<{
  title: string
  matches: any[]
  isToday?: boolean
  isYesterday?: boolean
  isTomorrow?: boolean
  compact?: boolean
}>()
</script>

<template>
  <section class="timeline-day" :class="{ compact }">
    <div class="timeline-day-head">
      <span class="day-title">{{ title }}</span>
      <span v-if="isToday" class="tag live">今天</span>
      <span v-else-if="isYesterday" class="tag">昨天</span>
      <span v-else-if="isTomorrow" class="tag blue">明天</span>
      <span class="day-count">{{ matches.length }} 场</span>
    </div>

    <div class="timeline-list">
      <TimelineMatchCard
        v-for="(match, index) in matches"
        :key="match.id"
        :match="match"
        :first="index === 0"
        :last="index === matches.length - 1"
      />
    </div>
  </section>
</template>

<style scoped>
.timeline-day {
  position: relative;
  display: grid;
  gap: 6px;
  margin-top: 18px;
}

.timeline-day-head {
  position: sticky;
  top: 0;
  z-index: 2;
  display: flex;
  align-items: center;
  gap: 8px;
  min-height: 34px;
  padding: 4px 0 4px 54px;
  background: linear-gradient(180deg, var(--bg) 72%, color-mix(in srgb, var(--bg) 0%, transparent));
  font-size: 14px;
  font-weight: 800;
  color: var(--text);
}

.timeline-day-head::before {
  content: '';
  position: absolute;
  left: 26px;
  top: 50%;
  width: 13px;
  height: 13px;
  border: 3px solid var(--bg);
  border-radius: 999px;
  background: var(--primary);
  box-shadow:
    0 0 0 1px color-mix(in srgb, var(--primary) 22%, var(--line)),
    0 0 0 6px color-mix(in srgb, var(--primary) 8%, transparent);
  transform: translate(-50%, -50%);
}

.day-title {
  flex: 0 0 auto;
}

.day-count {
  margin-left: auto;
  color: var(--weak);
  font-size: 11px;
  font-weight: 750;
}

.timeline-list {
  display: grid;
  gap: 0;
}

.compact {
  margin-top: 14px;
}

.compact .timeline-day-head {
  min-height: 30px;
  padding-left: 48px;
  font-size: 13px;
}

.compact .timeline-day-head::before {
  left: 24px;
  width: 11px;
  height: 11px;
}
</style>
