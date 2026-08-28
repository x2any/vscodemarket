// Constitution Principle III guard: no download actions in shipped code.
//
// Scope (intentionally narrow):
//   - Audits src/components/ + src/pages/ + src/i18n/ only.
//   - Skips src/api/ because that layer's `fetch` calls talk to OUR
//     backend, which is the legitimate frontend ↔ backend channel.
//   - Allows comment-shaped lines and the `forbidden` literal used inside
//     the no-download-literals test fixture itself.
//
// Run as part of CI; fails the build if a literal slips into a non-API
// file. If a new `src/<area>/` directory appears, extend SCAN_DIRS below.
import { describe, it, expect } from 'vitest'
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join, extname } from 'node:path'

const ROOT = join(__dirname, '..', 'src')
const SCAN_DIRS = ['components', 'pages', 'i18n']
const PATTERNS = [/\bwget\b/, /\bcurl\b/, /Invoke-WebRequest/]
const EXT_ALLOW = new Set(['.ts', '.vue', '.js'])
const COMMENT_OK_LINES = /^(?:\s*(?:\/\/|\*|\#))/ // Comment-shaped lines may mention the terms educationally

function* walk(dir: string): Generator<string> {
  if (!exists(dir)) return
  for (const e of readdirSync(dir)) {
    const p = join(dir, e)
    const s = statSync(p)
    if (s.isDirectory()) yield* walk(p)
    else if (EXT_ALLOW.has(extname(p))) yield p
  }
}

function exists(p: string) {
  try { statSync(p); return true } catch { return false }
}

describe('zero-network-script principle', () => {
  it('shipped UI contains no wget/curl/Invoke-WebRequest', () => {
    const offenders: string[] = []
    for (const sub of SCAN_DIRS) {
      for (const file of walk(join(ROOT, sub))) {
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
    }
    expect(offenders, offenders.join('\n')).toEqual([])
  })
})
