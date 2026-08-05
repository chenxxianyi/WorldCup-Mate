/**
 * LIVE-01: page-level live polling with adaptive intervals.
 *
 * Interval strategy:
 *   - any live match            -> 30s   (scores/minutes move fast)
 *   - kickoff within 60 minutes -> 60s   (pre-match updates)
 *   - otherwise                 -> 5 min (idle, cheap)
 *
 * Guarantees:
 *   - pauses while the tab is hidden (no background request storms) and
 *     refreshes immediately on visibility restore.
 *   - no overlapping requests: a slow refresh skips the next tick.
 *   - timers are cleaned up on component unmount.
 */
import { onBeforeUnmount } from 'vue'

export interface PollingStatus {
  hasLive: boolean
  /** Minutes until the nearest upcoming kickoff, or null when none. */
  nextKickoffInMinutes: number | null
}

export const LIVE_POLL_MS = 30_000
export const NEAR_POLL_MS = 60_000
export const IDLE_POLL_MS = 5 * 60_000

/** Pure interval decision — unit-testable (LIVE-01E). */
export function intervalFor(status: PollingStatus): number {
  if (status.hasLive) return LIVE_POLL_MS
  if (status.nextKickoffInMinutes != null && status.nextKickoffInMinutes <= 60) return NEAR_POLL_MS
  return IDLE_POLL_MS
}

/** Derive the polling status from the currently displayed matches. */
export function pollingStatusFrom(
  matches: { status: string; kickoff_time_utc?: string | null }[]
): PollingStatus {
  let hasLive = false
  let nextKickoffInMinutes: number | null = null
  const now = Date.now()
  for (const m of matches) {
    if (m.status === 'live') hasLive = true
    if ((m.status === 'scheduled' || m.status === 'upcoming') && m.kickoff_time_utc) {
      const diffMs = new Date(m.kickoff_time_utc).getTime() - now
      if (diffMs > 0 && (nextKickoffInMinutes == null || diffMs < nextKickoffInMinutes * 60_000)) {
        nextKickoffInMinutes = diffMs / 60_000
      }
    }
  }
  return { hasLive, nextKickoffInMinutes }
}

export function useLivePolling(
  status: () => PollingStatus,
  refresh: () => void | Promise<void>,
  enabled: () => boolean = () => true
) {
  let timer: ReturnType<typeof setTimeout> | null = null
  let visible = !document.hidden
  let running = false

  function clear() {
    if (timer !== null) {
      clearTimeout(timer)
      timer = null
    }
  }

  /** (Re)arm the next tick with the freshest status-derived interval. */
  function schedule() {
    clear()
    if (!visible || !enabled()) return
    timer = setTimeout(() => void tick(), intervalFor(status()))
  }

  async function tick() {
    if (running) {
      // A previous refresh is still in flight: skip this round, keep the
      // loop armed (no overlapping requests).
      schedule()
      return
    }
    running = true
    try {
      await refresh()
    } catch {
      // silent: polling failures must not surface errors on every tick
    } finally {
      running = false
    }
    // Re-evaluate the interval AFTER the refresh: a match that just went
    // live switches the loop from 60s/5min to 30s immediately.
    schedule()
  }

  function onVisibility() {
    visible = !document.hidden
    if (visible) {
      void tick() // refresh immediately when the tab becomes visible
    } else {
      clear() // fully pause in the background
    }
  }
  document.addEventListener('visibilitychange', onVisibility)
  schedule()

  onBeforeUnmount(() => {
    clear()
    document.removeEventListener('visibilitychange', onVisibility)
  })

  return { schedule }
}
