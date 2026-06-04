import request from './request'

export function apiGetSyncStatus() {
  return request.get('/sync/status') as Promise<any[]>
}
