import request from './request'
import type { ApiReminder, ReminderCreateInput, ReminderBatchInput } from '@/types/reminder'

export function apiCreateReminder(data: ReminderCreateInput) {
  return request.post('/reminders', data) as Promise<ApiReminder>
}

export function apiCreateReminderBatch(data: ReminderBatchInput) {
  return request.post('/reminders/batch', data) as Promise<ApiReminder[]>
}

export function apiListReminders() {
  return request.get('/reminders') as Promise<ApiReminder[]>
}

export function apiUpdateReminder(id: number, data: { remindBeforeMinutes?: number; channel?: string }) {
  return request.put(`/reminders/${id}`, data) as Promise<ApiReminder>
}

export function apiDeleteReminder(id: number) {
  return request.delete(`/reminders/${id}`) as Promise<{ deleted: boolean }>
}
