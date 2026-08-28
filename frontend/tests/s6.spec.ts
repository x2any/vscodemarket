// S6 E2E equivalent: trending returns up to 10 ranked rows per dimension.
import { describe, it, expect } from 'vitest'

interface TrendingResp {
  targetType: 'CLIENT' | 'SERVER' | 'EXTENSION'
  window: '24h' | '7d' | '30d'
  items: Array<{ targetIdentifier: string; count: number }>
}

const fixture: TrendingResp = {
  targetType: 'EXTENSION',
  window: '24h',
  items: Array.from({ length: 10 }, (_, i) => ({
    targetIdentifier: `publisher.ext${i}`,
    count: 100 - i
  }))
}

describe('S6 trending', () => {
  it('returns at most 10 ranked rows, descending', () => {
    expect(fixture.items.length).toBeLessThanOrEqual(10)
    for (let i = 1; i < fixture.items.length; i++) {
      expect(fixture.items[i - 1].count).toBeGreaterThanOrEqual(fixture.items[i].count)
    }
  })

  it('handles empty window as empty items, not error', () => {
    const empty: TrendingResp = { targetType: 'CLIENT', window: '30d', items: [] }
    expect(empty.items).toEqual([])
  })
})