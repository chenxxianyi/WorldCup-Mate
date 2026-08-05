// @vitest-environment happy-dom
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

// Collect unmount hooks so we can prove timers/listeners are cleaned up
// without mounting a real component.
const unmountFns: (() => void)[] = []
vi.mock('vue', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue')>()
  return { ...actual, onBeforeUnmount: (fn: () => void) => { unmountFns.push(fn) } }
})

import { useLivePolling, LIVE_POLL_MS, IDLE_POLL_MS } from './useLivePolling'

describe('useLivePolling lifecycle (LIVE-02E)', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    unmountFns.length = 0
    Object.defineProperty(document, 'hidden', { configurable: true, value: false })
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('does not tick while disabled; arms after data arrives + schedule()', async () => {
    const refresh = vi.fn()
    let isEnabled = false
    const polling = useLivePolling(
      () => ({ hasLive: false, nextKickoffInMinutes: null }),
      refresh,
      () => isEnabled,
    )
    await vi.advanceTimersByTimeAsync(IDLE_POLL_MS * 3)
    expect(refresh).not.toHaveBeenCalled()
    // Mirrors the pages' load().finally { polling.schedule() } behaviour:
    // data has arrived so enabled() now passes.
    isEnabled = true
    polling.schedule()
    await vi.advanceTimersByTimeAsync(IDLE_POLL_MS - 1)
    expect(refresh).not.toHaveBeenCalled()
    await vi.advanceTimersByTimeAsync(1)
    expect(refresh).toHaveBeenCalledTimes(1)
  })

  it('adapts the interval when a match goes live', async () => {
    let status = { hasLive: false, nextKickoffInMinutes: null }
    const refresh = vi.fn()
    useLivePolling(() => status, refresh, () => true)
    // The match goes live while the first idle tick is still pending.
    status = { hasLive: true, nextKickoffInMinutes: null }
    await vi.advanceTimersByTimeAsync(IDLE_POLL_MS) // t=300s: first tick fires
    expect(refresh).toHaveBeenCalledTimes(1)
    // After the tick, the loop re-arms with the fresh (live) status: the
    // next tick must come after LIVE_POLL_MS, not IDLE_POLL_MS.
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS) // t=330s
    expect(refresh).toHaveBeenCalledTimes(2)
  })

  it('pauses while the tab is hidden and refreshes on restore', async () => {
    const refresh = vi.fn()
    useLivePolling(() => ({ hasLive: true, nextKickoffInMinutes: null }), refresh, () => true)
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS)
    expect(refresh).toHaveBeenCalledTimes(1)

    // Hide the tab: no further ticks.
    Object.defineProperty(document, 'hidden', { configurable: true, value: true })
    document.dispatchEvent(new Event('visibilitychange'))
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS * 5)
    expect(refresh).toHaveBeenCalledTimes(1)

    // Restore: immediate refresh + timer re-armed.
    Object.defineProperty(document, 'hidden', { configurable: true, value: false })
    document.dispatchEvent(new Event('visibilitychange'))
    expect(refresh).toHaveBeenCalledTimes(2)
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS)
    expect(refresh).toHaveBeenCalledTimes(3)
  })

  it('cleans up timers and listeners on unmount', async () => {
    const refresh = vi.fn()
    useLivePolling(() => ({ hasLive: true, nextKickoffInMinutes: null }), refresh, () => true)
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS)
    expect(refresh).toHaveBeenCalledTimes(1)
    expect(unmountFns.length).toBe(1)

    unmountFns.forEach((fn) => fn())
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS * 5)
    expect(refresh).toHaveBeenCalledTimes(1) // no more ticks after unmount
    expect(vi.getTimerCount()).toBe(0)

    // Listener removed: visibility restore must not tick either.
    Object.defineProperty(document, 'hidden', { configurable: true, value: true })
    document.dispatchEvent(new Event('visibilitychange'))
    Object.defineProperty(document, 'hidden', { configurable: true, value: false })
    document.dispatchEvent(new Event('visibilitychange'))
    expect(refresh).toHaveBeenCalledTimes(1)
  })

  it('skips overlapping ticks while a refresh is in flight', async () => {
    let resolveRefresh!: () => void
    const refresh = vi.fn(() => new Promise<void>((resolve) => { resolveRefresh = resolve }))
    const polling = useLivePolling(() => ({ hasLive: true, nextKickoffInMinutes: null }), refresh, () => true)
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS)
    expect(refresh).toHaveBeenCalledTimes(1)
    // While the first refresh is still pending, a page-side schedule()
    // (load().finally) arms the next tick — it must NOT fire another
    // request while `running` is still true.
    polling.schedule()
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS * 3)
    expect(refresh).toHaveBeenCalledTimes(1)
    resolveRefresh()
    await vi.advanceTimersByTimeAsync(LIVE_POLL_MS)
    expect(refresh).toHaveBeenCalledTimes(2) // loop re-armed and ticked again
  })
})
