import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useGeneration, debounce } from './useRequestGuard'

describe('useGeneration (LIVE-02)', () => {
  let gen: ReturnType<typeof useGeneration>

  beforeEach(() => {
    gen = useGeneration()
  })

  it('isCurrent is true for the newest claimed generation', () => {
    const g = gen.next()
    expect(gen.isCurrent(g)).toBe(true)
  })

  it('isCurrent is false after bump() (stale response)', () => {
    const g = gen.next()
    gen.bump()
    expect(gen.isCurrent(g)).toBe(false)
  })

  it('isCurrent is false after a newer next() claim', () => {
    const g1 = gen.next()
    const g2 = gen.next()
    expect(gen.isCurrent(g1)).toBe(false)
    expect(gen.isCurrent(g2)).toBe(true)
  })
})

describe('debounce (LIVE-02)', () => {
  beforeEach(() => {
    vi.useFakeTimers()
  })

  it('fires only once after rapid calls', () => {
    const fn = vi.fn()
    const d = debounce(fn, 300)
    d()
    d()
    d()
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(299)
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(1)
    expect(fn).toHaveBeenCalledTimes(1)
  })

  it('passes arguments through', () => {
    const fn = vi.fn()
    const d = debounce(fn, 100)
    d('a', 1)
    vi.advanceTimersByTime(100)
    expect(fn).toHaveBeenCalledWith('a', 1)
  })

  it('resets the timer on every call', () => {
    const fn = vi.fn()
    const d = debounce(fn, 300)
    d()
    vi.advanceTimersByTime(200)
    d()
    vi.advanceTimersByTime(200)
    expect(fn).not.toHaveBeenCalled()
    vi.advanceTimersByTime(100)
    expect(fn).toHaveBeenCalledTimes(1)
  })
})
