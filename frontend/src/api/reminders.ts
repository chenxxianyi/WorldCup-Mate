import request from './request'

export function apiCreateReminder(data: { matchId: number; remindBeforeMinutes?: number; channel?: string }) {
  return request.post('/reminders', data) as Promise<any>
}

export function apiListReminders() {
  return request.get('/reminders') as Promise<any[]>
}

export function apiUpdateReminder(id: number, data: { remindBeforeMinutes?: number; channel?: string }) {
  return request.put(`/reminders/${id}`, data) as Promise<any>
}

export function apiDeleteReminder(id: number) {
  return request.delete(`/reminders/${id}`) as Promise<any>
}
