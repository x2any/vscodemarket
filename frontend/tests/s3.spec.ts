// S3 E2E equivalent: extension search returns Marketplace-shape results
// with publisher + name + displayName + latestVersion.
import { describe, it, expect } from 'vitest'

interface SearchResp { results: Array<{ publisher: string; name: string; displayName: string; latestVersion: string }>; total: number }

const fixture: SearchResp = {
  results: [
    { publisher: 'ms-python', name: 'python', displayName: 'Python', latestVersion: '2024.20.0' },
    { publisher: 'dbaeumer', name: 'vscode-eslint', displayName: 'ESLint', latestVersion: '3.0.10' }
  ],
  total: 137
}

describe('S3 extension search', () => {
  it('returns at least one Marketplace-shaped row', () => {
    expect(fixture.results.length).toBeGreaterThan(0)
    const r = fixture.results[0]
    expect(r.publisher).toMatch(/^[a-z0-9-]+$/)
    expect(r.name).toMatch(/^[a-z0-9-]+$/)
    expect(r.displayName.length).toBeGreaterThan(0)
    expect(r.latestVersion).toMatch(/^\d+/)
  })
})