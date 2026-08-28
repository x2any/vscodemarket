// Constitution Principle III guard: no download actions in shipped code.
// Run as part of CI; fails the build if a literal slips in.
import { describe, it, expect } from 'vitest'
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join, extname } from 'node:path'

const ROOT = join(__dirname, '..', 'src')
const PATTERNS = [/\bwget\b/, /\bcurl\b/, /Invoke-WebRequest/, /\bfetch\s*\(/]
const EXT_ALLOW = new Set(['.ts', '.vue', '.js'])
const COMMENT_OK_LINES = /^(?:\s*(?:\/\/|\*|\#))/ // Comment-shaped lines may mention the terms educationally

function* walk(dir: string): Generator<string> {
  for (const e of readdirSync(dir)) {
    const p = join(dir, e)
    const s = statSync(p)
    if (s.isDirectory()) yield* walk(p)
    else if (EXT_ALLOW.has(extname(p))) yield p
  }
}

describe('zero-network-script principle', () => {
  it('frontend src contains no wget/curl/Invoke-WebRequest/fetch call', () => {
    const offenders: string[] = []
    for (const file of walk(ROOT)) {
      const lines = readFileSync(file, 'utf8').split('\n')
      lines.forEach((line, idx) => {
        if (COMMENT_OK_LINES.test(line)) return
        for (const re of PATTERNS) {
          if (re.test(line)) {
            offenders.push(`${file}:${idx + 1}: ${line.trim()}`)
          }
        }
      })
    }
    expect(offenders, offenders.join('\n')).toEqual([])
  })
})
