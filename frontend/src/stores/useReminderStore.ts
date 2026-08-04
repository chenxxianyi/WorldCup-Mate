import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
  apiCreateReminder,
  apiCreateReminderBatch,
  apiDeleteReminder,
  apiListReminders,
} from '@/api/reminders'
import { normalizeMatch, type Match } from '@/types/match'
import type { ApiReminder } from '@/types/reminder'

export interface ReminderItem {
  id: number
  match_id: number
  remind_before_minutes: number
  channel: string
  match: Match | null
}

function normalizeReminder(r: ApiReminder, fallbackChannel = 'site'): ReminderItem {
  return {
    id: r.id,
    match_id: r.match_id,
    remind_before_minutes: r.remind_before_minutes,
    channel: r.channel || fallbackChannel,
    match: r.match ? normalizeMatch(r.match) : null,
  }
}

export const useReminderStore = defineStore('reminder', () => {
  const reminders = ref<ReminderItem[]>([])
  const matchReminderIds = ref<Set<number>>(new Set())
  const count = computed(() => reminders.value.length)

  function refreshMatchReminderIds() {
    matchReminderIds.value = new Set(reminders.value.map((r) => r.match_id))
  }

  async function fetchReminders() {
    try {
      const res = await apiListReminders()
      reminders.value = res.map((r) => normalizeReminder(r))
      refreshMatchReminderIds()
    } catch {
      reminders.value = []
      refreshMatchReminderIds()
    }
  }

  function hasReminder(matchId: number) {
    return matchReminderIds.value.has(matchId)
  }

  function remindersForMatch(matchId: number) {
    return reminders.value.filter((r) => r.match_id === matchId)
  }

  async function removeRemindersByMatch(matchId: number) {
    const existing = remindersForMatch(matchId)
    if (!existing.length) return true

    reminders.value = reminders.value.filter((r) => r.match_id !== matchId)
    refreshMatchReminderIds()
    try {
      await Promise.all(existing.map((item) => apiDeleteReminder(item.id)))
      return true
    } catch {
      reminders.value.push(...existing)
      refreshMatchReminderIds()
      return false
    }
  }

  async function createReminderBatch(matchId: number, minutes: number[], channel = 'site') {
    const res = await apiCreateReminderBatch({ match_id: matchId, minutes, channel })
    const created = res.map((r) => normalizeReminder(r, channel))
    reminders.value.push(...created)
    refreshMatchReminderIds()
    return created
  }

  async function toggleReminder(matchId: number, minutesBefore = 30, channel = 'site') {
    if (hasReminder(matchId)) {
      return removeRemindersByMatch(matchId)
    }

    try {
      const res = await apiCreateReminder({ matchId, remindBeforeMinutes: minutesBefore, channel })
      reminders.value.push(normalizeReminder(res, channel))
      refreshMatchReminderIds()
      return true
    } catch {
      return false
    }
  }

  return {
    reminders,
    count,
    hasReminder,
    remindersForMatch,
    fetchReminders,
    removeRemindersByMatch,
    createReminderBatch,
    toggleReminder,
  }
})
