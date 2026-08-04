import request from './request'
import type { SyncState } from '@/types/sync'

export function apiGetSyncStatus() {
  return request.get('/sync/status') as Promise<SyncState[]>
}
