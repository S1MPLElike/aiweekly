const API_BASE = 'http://localhost:8082/api'

export interface User {
  id: number
  username: string
  phone: string
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  user: User
}

export async function refreshToken(refreshToken: string): Promise<AuthTokens> {
  const response = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ token: refreshToken }),
  })

  if (!response.ok) {
    throw new Error('Refresh token expired or invalid')
  }

  return response.json()
}

export async function logout(): Promise<void> {
  await fetch(`${API_BASE}/auth/logout`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
  })
}

const REFRESH_TOKEN_KEY = 'refresh_token'
const ACCESS_TOKEN_KEY = 'access_token'
const USER_KEY = 'auth_user'
const EXPIRY_KEY = 'token_expiry'

export function saveAuth(tokens: AuthTokens): void {
  localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refresh_token)
  localStorage.setItem(ACCESS_TOKEN_KEY, tokens.access_token)
  localStorage.setItem(USER_KEY, JSON.stringify(tokens.user))
  
  const expiry = parseJwtExpiry(tokens.access_token)
  localStorage.setItem(EXPIRY_KEY, expiry.toString())
}

function parseJwtExpiry(token: string): number {
  try {
    const payload = token.split('.')[1]
    const decoded = JSON.parse(atob(payload))
    return decoded.exp || Date.now() + 7200000
  } catch {
    return Date.now() + 7200000
  }
}

export function isAccessTokenValid(): boolean {
  const expiryStr = localStorage.getItem(EXPIRY_KEY)
  if (!expiryStr) return false
  
  const expiry = parseInt(expiryStr, 10)
  return Date.now() < expiry * 1000
}

export function getStoredUser(): User | null {
  const userStr = localStorage.getItem(USER_KEY)
  if (!userStr) return null
  try {
    return JSON.parse(userStr)
  } catch {
    return null
  }
}

export function getStoredRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

export function getStoredAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

export function clearAuth(): void {
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  localStorage.removeItem(EXPIRY_KEY)
}
