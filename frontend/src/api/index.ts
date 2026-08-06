import request from './request'

export * from './auth'
export * from './matches'
export * from './teams'
export * from './standings'
export * from './favorites'
export * from './reminders'
export * from './notifications'
export * from './sync'
export * from './admin'

// Homepage aggregation endpoint (DATA-10).
export interface HomeAggregate {
  upcoming_matches: unknown[]
  live_matches: unknown[]
  competitions: unknown[]
  sync_states: unknown[]
  synced_at: string
}

export function apiHomeAggregate() {
  return request.get('/home') as Promise<HomeAggregate>
}
