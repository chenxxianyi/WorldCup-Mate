export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

/** Standard API error (API-03): stable code + user-facing message. */
export class ApiError extends Error {
  code: number
  requestId?: string

  constructor(code: number, message: string, requestId?: string) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.requestId = requestId
  }
}

/** True when the request was aborted by the caller (no UI error needed). */
export function isCancel(err: unknown): boolean {
  return (
    typeof err === 'object' &&
    err !== null &&
    'code' in err &&
    (err as { code?: string }).code === 'ERR_CANCELED'
  )
}
