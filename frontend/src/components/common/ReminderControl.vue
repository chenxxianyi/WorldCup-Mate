<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useReminderStore } from '@/stores/useReminderStore'
import { useSettingStore } from '@/stores/useSettingStore'

const props = withDefaults(defineProps<{
  matchId: number
  mode?: 'icon' | 'pill'
}>(), {
  mode: 'icon',
})

const router = useRouter()
const auth = useAuthStore()
const reminder = useReminderStore()
const settings = useSettingStore()

const open = ref(false)
const saving = ref(false)
const error = ref('')
const selectedMinutes = ref<number[]>([60])
const selectedChannel = ref(settings.defaultReminderChannel || 'site')

const minuteOptions = [
  { label: '赛前 1 天', value: 1440 },
  { label: '赛前 1 小时', value: 60 },
  { label: '赛前 15 分钟', value: 15 },
]

const hasReminder = computed(() => reminder.hasReminder(props.matchId))
const hasEmailTarget = computed(() => !auth.user || Boolean(auth.user.notificationEmail || auth.user.email))
const currentReminders = computed(() => reminder.remindersForMatch(props.matchId))
const currentSummary = computed(() => {
  if (!currentReminders.value.length) return ''
  return currentReminders.value
    .map((item) => minuteOptions.find((option) => option.value === item.remind_before_minutes)?.label || `${item.remind_before_minutes} 分钟`)
    .join('、')
})

function toggleMinute(value: number) {
  if (selectedMinutes.value.includes(value)) {
    selectedMinutes.value = selectedMinutes.value.filter((item) => item !== value)
  } else {
    selectedMinutes.value = [...selectedMinutes.value, value]
  }
}

async function onTrigger() {
  error.value = ''
  if (!auth.isLoggedIn) {
    router.push('/login')
    return
  }
  if (hasReminder.value) {
    await reminder.removeRemindersByMatch(props.matchId)
    open.value = false
    return
  }
  selectedChannel.value = settings.defaultReminderChannel || 'site'
  open.value = !open.value
}

async function createReminders() {
  error.value = ''
  if (!selectedMinutes.value.length) {
    error.value = '请选择提醒时间'
    return
  }
  if (selectedChannel.value === 'email' && !hasEmailTarget.value) {
    error.value = '请先在个人中心设置通知邮箱'
    return
  }
  saving.value = true
  try {
    const minutes = [...selectedMinutes.value].sort((a, b) => b - a)
    await reminder.createReminderBatch(props.matchId, minutes, selectedChannel.value)
    settings.setDefaultReminderChannel(selectedChannel.value)
    open.value = false
  } catch (err: any) {
    error.value = err?.message || '提醒创建失败'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div
    class="reminder-control"
    @click.stop
  >
    <button
      v-if="mode === 'icon'"
      class="icon-trigger"
      :class="{ active: hasReminder }"
      title="比赛提醒"
      @click="onTrigger"
    >
      <span
        class="material-symbols-outlined"
        :style="hasReminder ? 'font-variation-settings: \'FILL\' 1' : ''"
      >notifications</span>
    </button>
    <button
      v-else
      class="pill-btn primary reminder-pill"
      :class="{ active: hasReminder }"
      @click="onTrigger"
    >
      <span
        class="material-symbols-outlined"
        :style="hasReminder ? 'font-variation-settings: \'FILL\' 1' : ''"
      >notifications</span>
      {{ hasReminder ? '取消提醒' : '设置提醒' }}
    </button>

    <div
      v-if="open && !hasReminder"
      class="reminder-popover"
    >
      <div class="popover-title">
        提醒时间
      </div>
      <div class="option-grid">
        <button
          v-for="option in minuteOptions"
          :key="option.value"
          class="choice"
          :class="{ active: selectedMinutes.includes(option.value) }"
          @click="toggleMinute(option.value)"
        >
          {{ option.label }}
        </button>
      </div>

      <div class="popover-title">
        提醒渠道
      </div>
      <div class="option-grid two">
        <button
          class="choice"
          :class="{ active: selectedChannel === 'site' }"
          @click="selectedChannel = 'site'"
        >
          站内通知
        </button>
        <button
          class="choice"
          :class="{ active: selectedChannel === 'email' }"
          @click="selectedChannel = 'email'"
        >
          邮件通知
        </button>
      </div>
      <p
        v-if="selectedChannel === 'email' && !hasEmailTarget"
        class="hint-text"
      >
        请先在个人中心设置通知邮箱
      </p>

      <p
        v-if="error"
        class="error-text"
      >
        {{ error }}
      </p>

      <div class="popover-actions">
        <button
          class="text-btn"
          @click="open = false"
        >
          取消
        </button>
        <button
          class="save-btn"
          :disabled="saving"
          @click="createReminders"
        >
          {{ saving ? '创建中...' : '确认' }}
        </button>
      </div>
    </div>

    <div
      v-if="hasReminder && mode === 'pill' && currentSummary"
      class="summary"
    >
      {{ currentSummary }}
    </div>
  </div>
</template>

<style scoped>
.reminder-control {
  position: relative;
  display: inline-grid;
  gap: 6px;
}

.icon-trigger {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 999px;
  color: var(--weak);
  background: transparent;
  cursor: pointer;
}

.icon-trigger.active {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.icon-trigger .material-symbols-outlined {
  font-size: 22px;
}

.reminder-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.reminder-pill .material-symbols-outlined {
  font-size: 18px;
}

.reminder-popover {
  position: absolute;
  right: 0;
  bottom: calc(100% + 8px);
  z-index: 120;
  width: min(260px, calc(100vw - 44px));
  padding: 12px;
  border: 1px solid var(--line);
  border-radius: 14px;
  background: var(--card);
  box-shadow: 0 14px 34px rgba(15, 23, 42, 0.18);
}

.popover-title {
  margin-bottom: 8px;
  color: var(--muted);
  font-size: 12px;
  font-weight: 700;
}

.popover-title:not(:first-child) {
  margin-top: 12px;
}

.option-grid {
  display: grid;
  gap: 8px;
}

.option-grid.two {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.choice {
  min-height: 34px;
  padding: 0 10px;
  border: 1px solid var(--line);
  border-radius: 10px;
  color: var(--text);
  background: var(--card-soft);
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.choice.active {
  color: #fff;
  border-color: var(--primary);
  background: var(--primary);
}

.error-text {
  margin: 10px 0 0;
  color: var(--primary);
  font-size: 12px;
}

.hint-text {
  margin: 8px 0 0;
  color: var(--muted);
  font-size: 12px;
}

.popover-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.text-btn,
.save-btn {
  min-height: 32px;
  padding: 0 12px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
}

.text-btn {
  border: 0;
  color: var(--muted);
  background: transparent;
}

.save-btn {
  border: 0;
  color: #fff;
  background: var(--primary);
}

.save-btn:disabled {
  opacity: 0.6;
  cursor: wait;
}

.summary {
  color: var(--muted);
  font-size: 12px;
}
</style>
