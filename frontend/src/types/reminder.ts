import type { ApiMatch } from './match'

/** Backend reminder record (models.Reminder JSON). */
export interface ApiReminder {
  id: number
  user_id: number
  match_id: number
  remind_before_minutes: number
  remind_at: string
  channel: string
  status: string
  created_at: string
  match?: ApiMatch
}

export interface ReminderCreateInput {
  matchId: number
  remindBeforeMinutes?: number
  channel?: string
}

export interface ReminderBatchInput {
  match_id: number
  minutes: number[]
  channel?: string
}
