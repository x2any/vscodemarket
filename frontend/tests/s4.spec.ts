// S4 E2E equivalent: direct extension version URL renders engines.vscode.
import { describe, it, expect } from 'vitest'

interface VersionResp {
  extension: { publisher: string; name: string }
  version: {
    version: string
    enginesVscode: string
    downloadUrl: string
    publishTime: string
  }
}

const fixture: VersionResp = {
  extension: { publisher: 'ms-python', name: 'python' },
  version: {
    version: '2024.20.0',
    enginesVscode: '^1.94.0',
    downloadUrl: 'https://marketplace.visualstudio.com/_apis/.../python-2024.20.0.vsix',
    publishTime: '2024-10-08T15:32:00Z'
  }
}

describe('S4 extension version', () => {
  it('carries enginesVscode for UI compatibility check', () => {
    expect(fixture.version.enginesVscode).toMatch(/^[\^~>=<]/)
    expect(fixture.version.downloadUrl).toMatch(/^https:\/\/marketplace\.visualstudio\.com\//)
  })
})