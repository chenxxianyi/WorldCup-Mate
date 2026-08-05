/**
 * LIVE-02: stale-response protection for async page loads.
 *
 * A generation counter lets async flows discard responses that arrive after
 * the user has already switched away (competition change, route change).
 * Usage:
 *
 *   const gen = useGeneration()
 *   async function load() {
 *     const g = gen.next()
 *     const data = await fetch()
 *     if (!gen.isCurrent(g)) return   // a newer load started: drop this
 *     apply(data)
 *   }
 *   watch(competition, () => gen.bump())  // invalidate in-flight loads
 */
export function useGeneration() {
  let generation = 0

  /** Claim the current generation for a new async flow. */
  function next(): number {
    generation += 1
    return generation
  }

  /** True when `g` is still the newest claimed generation. */
  function isCurrent(g: number): boolean {
    return g === generation
  }

  /** Invalidate every in-flight flow (call on switch/route change). */
  function bump(): void {
    generation += 1
  }

  return { next, isCurrent, bump }
}

/**
 * Debounce helper: returns a function that delays `fn` until `wait` ms
 * passed without another call (LIVE-02D, search inputs).
 */
export function debounce<A extends unknown[]>(fn: (...args: A) => void, wait = 300) {
  let timer: ReturnType<typeof setTimeout> | null = null
  return (...args: A) => {
    if (timer !== null) clearTimeout(timer)
    timer = setTimeout(() => {
      timer = null
      fn(...args)
    }, wait)
  }
}
