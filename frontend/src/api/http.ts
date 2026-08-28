import { ElMessage } from 'element-plus'

const baseURL = '/api/v1'

export interface APIError {
  code: string
  message_zh: string
  message_en: string
}

export class HttpError extends Error {
  status: number
  body: APIError | null
  constructor(status: number, body: APIError | null, fallback: string) {
    super(body?.message_zh ?? fallback)
    this.status = status
    this.body = body
  }
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(baseURL + path, {
    headers: { 'Content-Type': 'application/json', ...(init.headers ?? {}) },
    ...init
  })
  if (!res.ok) {
    let body: APIError | null = null
    try {
      const j = await res.json()
      body = j?.error ?? null
    } catch { /* ignore parse failure */ }
    const err = new HttpError(res.status, body, `HTTP ${res.status}`)
    ElMessage.error(body?.message_zh ?? err.message)
    throw err
  }
  if (res.status === 202) return undefined as unknown as T
  return res.json() as Promise<T>
}

export const http = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) })
}
