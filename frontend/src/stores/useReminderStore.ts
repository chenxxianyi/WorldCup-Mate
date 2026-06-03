import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  apiCreateReminder,
  apiListReminders,
  apiDeleteReminder,
} from '@/api/reminders'
import { normalizeMatch, type Match } from '@/types/match'

export interface ReminderItem {
  id: number
  match_id: number
  remind_before_minutes: number
  channel: string
  match: Match | null
}

export const useReminderStore = defineStore('reminder', () => {
  const reminders = ref<ReminderItem[]>([])
  const matchReminderIds = ref<Set<number>>(new Set())
  const count = computed(() => reminders.value.length)

  async function fetchReminders() {
    try {
      const res = await apiListReminders() as any[]
      reminders.value = res.map((r: any) => ({
        id: r.id,
        match_id: r.match_id,
        remind_before_minutes: r.remind_before_minutes,
        channel: r.channel || 'site',
        match: r.match ? normalizeMatch(r.match) : null,
      }))
      matchReminderIds.value = new Set(reminders.value.map((r) => r.match_id))
    } catch {
      reminders.value = []
      matchReminderIds.value = new Set()
    }
  }

  function hasReminder(matchId: number) {
    return matchReminderIds.value.has(matchId)
  }

  async function toggleReminder(matchId: number, minutesBefore = 30, channel = 'site') {
    const existing = reminders.value.find((r) => r.match_id === matchId)
    if (existing) {
      reminders.value = reminders.value.filter((r) => r.match_id !== matchId)
      matchReminderIds.value.delete(matchId)
      try {
        await apiDeleteReminder(existing.id)
      } catch {
        reminders.value.push(existing)
        matchReminderIds.value.add(matchId)
      }
    } else {
      try {
        const res = await apiCreateReminder({ matchId, remindBeforeMinutes: minutesBefore, channel }) as any
        reminders.value.push({
          id: res.id,
          match_id: matchId,
          remind_before_minutes: minutesBefore,
          channel,
          match: null,
        })
        matchReminderIds.value.add(matchId)
      } catch {}
    }
  }

  return { reminders, count, hasReminder, fetchReminders, toggleReminder }
})
