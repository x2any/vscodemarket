// Fire-and-forget event tracking. Failures are silent by design:
// per FR-012 / SC-005, event tracking must NEVER block the user's main path.
import { http } from './http'

export type EventType = 'SEARCH' | 'DOWNLOAD'
export type TargetType = 'CLIENT' | 'SERVER' | 'EXTENSION'

export interface TrackPayload {
  eventType: EventType
  targetType: TargetType
  targetIdentifier: string
  platform?: string
  architecture?: string
  channel?: string
}

export function track(payload: TrackPayload): void {
  // Use sendBeacon when available so navigation isn't blocked; fall back
  // to fetch keepalive. We never await — promise is intentionally discarded.
  try {
    const body = JSON.stringify(payload)
    if (typeof navigator !== 'undefined' && typeof navigator.sendBeacon === 'function') {
      const blob = new Blob([body], { type: 'application/json' })
      navigator.sendBeacon('/api/v1/events', blob)
      return
    }
    void fetch('/api/v1/events', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body,
      keepalive: true
    })
  } catch {
    // Swallow — tracking failures must not surface to the user.
    void http // mark used; no further action
  }
}