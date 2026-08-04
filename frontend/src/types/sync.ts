/** Backend sync state record (models.SyncState JSON). */
export interface SyncState {
  id: number
  provider: string
  resource: string
  status: string
  last_synced_at: string | null
  next_sync_at: string | null
  last_error: string
}

/** Result of a league sync run (services.LeagueSyncResult JSON). */
export interface LeagueSyncResult {
  provider: string
  resource: string
  reason: string
  total: number
  created: number
  updated: number
  skipped: number
  started_at: string
  finished_at: string
}
