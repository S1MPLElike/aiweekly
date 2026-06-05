import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { getStoredRefreshToken, refreshToken, saveAuth, clearAuth, isAccessTokenValid, getStoredUser } from './services/auth'
import { initTheme } from './services/theme'
import './styles/theme.css'

export async function initAuth() {
  const user = getStoredUser()
  if (!user) return

  if (isAccessTokenValid()) {
    return
  }

  const refreshTokenStr = getStoredRefreshToken()
  if (!refreshTokenStr) {
    clearAuth()
    return
  }

  try {
    const tokens = await refreshToken(refreshTokenStr)
    saveAuth(tokens)
  } catch {
    clearAuth()
  }
}

async function bootstrap() {
  await initAuth()
  initTheme()

  const app = createApp(App)
  app.use(router)
  app.mount('#app')
}

bootstrap()
