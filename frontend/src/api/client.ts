const API_BASE = 'http://localhost:8082/api'

import { getStoredAccessToken, getStoredRefreshToken, isAccessTokenValid, saveAuth, clearAuth, refreshToken } from '../services/auth'

let isRefreshing = false
let refreshPromise: Promise<void> | null = null

export async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
  const token = getStoredAccessToken()
  
  if (!token || !isAccessTokenValid()) {
    await ensureValidToken()
  }

  const headers = new Headers(options.headers)
  const accessToken = getStoredAccessToken()
  if (accessToken) {
    headers.set('Authorization', `Bearer ${accessToken}`)
  }

  const response = await fetch(`${API_BASE}${url}`, {
    ...options,
    headers,
    credentials: 'include',
  })

  if (response.status === 401) {
    await ensureValidToken(true)
    
    const newAccessToken = getStoredAccessToken()
    if (newAccessToken) {
      headers.set('Authorization', `Bearer ${newAccessToken}`)
      return fetch(`${API_BASE}${url}`, {
        ...options,
        headers,
        credentials: 'include',
      })
    }
  }

  return response
}

async function ensureValidToken(forceRefresh = false): Promise<void> {
  if (!forceRefresh && isAccessTokenValid()) {
    return
  }

  if (isRefreshing) {
    await refreshPromise
    return
  }

  isRefreshing = true
  refreshPromise = performRefresh()
  await refreshPromise
  isRefreshing = false
  refreshPromise = null
}

async function performRefresh(): Promise<void> {
  const refreshTokenStr = getStoredRefreshToken()
  if (!refreshTokenStr) {
    clearAuth()
    window.location.href = '/login'
    return
  }

  try {
    const tokens = await refreshToken(refreshTokenStr)
    saveAuth(tokens)
  } catch {
    clearAuth()
    window.location.href = '/login'
  }
}
