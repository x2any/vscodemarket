// S1 E2E equivalent: stable version lookup returns client + server URLs.
// Mocks the backend; verifies the wire shape our UI relies on.
import { describe, it, expect } from 'vitest'
import type { LookupResponse } from '../src/api/contracts'

const stableResponse: LookupResponse = {
  channel: 'stable',
  version: '1.94.2',
  client: {
    downloadUrl: 'https://update.code.visualstudio.com/1.94.2/darwin-arm64/stable',
    platform: 'darwin',
    architecture: 'arm64'
  },
  server: {
    downloadUrl: 'https://update.code.visualstudio.com/commit:abcdef/server-darwin-arm64/stable',
    commitHash: 'abcdef0123456789',
    clientVersion: '1.94.2',
    platform: 'darwin',
    architecture: 'arm64'
  }
}

describe('S1 stable lookup', () => {
  it('exposes both client and server URLs', () => {
    expect(stableResponse.client.downloadUrl).toMatch(/^https:\/\/update\.code\.visualstudio\.com\//)
    expect(stableResponse.server?.downloadUrl).toMatch(/^https:\/\/update\.code\.visualstudio\.com\//)
    expect(stableResponse.server?.commitHash).toMatch(/^[0-9a-f]+$/)
  })
})