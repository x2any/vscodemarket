export type Channel = 'stable' | 'insider'
export type Platform = 'windows' | 'linux' | 'darwin'
export type Architecture = 'x86_64' | 'arm64' | 'armv7'

export interface ClientPayload {
  downloadUrl: string
  platform: Platform
  architecture: Architecture
}

export interface ServerPayload {
  downloadUrl: string
  commitHash: string
  clientVersion: string
  platform: Platform
  architecture: Architecture
}

export interface LookupResponse {
  channel: Channel
  version: string
  commit?: string
  clients: ClientPayload[]
  servers?: ServerPayload[]
}

export interface LookupRequest {
  channel: Channel
  version: string
  platform?: Platform
  architecture?: Architecture
}