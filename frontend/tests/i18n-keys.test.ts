// Vitest-style guard. Invoked manually or in CI; ensures zh-CN and en-US
// expose the same top-level key set so no UI string gets lost in translation.
import { describe, it, expect } from 'vitest'
import { zhCN } from '../src/i18n/zh-CN'
import { enUS } from '../src/i18n/en-US'

function flatKeys(obj: Record<string, unknown>, prefix = ''): string[] {
  const out: string[] = []
  for (const [k, v] of Object.entries(obj)) {
    const path = prefix ? `${prefix}.${k}` : k
    if (v && typeof v === 'object' && !Array.isArray(v)) {
      out.push(...flatKeys(v as Record<string, unknown>, path))
    } else {
      out.push(path)
    }
  }
  return out.sort()
}

describe('i18n parity', () => {
  it('zh-CN and en-US expose identical key sets', () => {
    const a = flatKeys(zhCN as unknown as Record<string, unknown>)
    const b = flatKeys(enUS as unknown as Record<string, unknown>)
    expect(b).toEqual(a)
  })
})
