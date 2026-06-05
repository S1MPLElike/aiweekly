export type ThemeId = 'cream' | 'tiffany' | 'mint'

interface Theme {
  id: ThemeId
  name: string
  primary: string
  bg: string
  surface: string
  surfaceSoft: string
  textPrimary: string
  textSecondary: string
  border: string
  accent: string
  accentStrong: string
  accentText: string
}

export const themes: Theme[] = [
  {
    id: 'cream',
    name: '浅黄色',
    primary: '#f5cd6a',
    bg: '#fffcf3',
    surface: '#ffffff',
    surfaceSoft: '#fff8df',
    textPrimary: '#2f2a22',
    textSecondary: '#766e60',
    border: '#f2e8c8',
    accent: '#f5cd6a',
    accentStrong: '#e8b947',
    accentText: '#463500'
  },
  {
    id: 'tiffany',
    name: '蒂芙尼蓝',
    primary: '#4ecdc4',
    bg: '#f0fdf9',
    surface: '#ffffff',
    surfaceSoft: '#ccfbef',
    textPrimary: '#0f172a',
    textSecondary: '#475569',
    border: '#a7f3d0',
    accent: '#4ecdc4',
    accentStrong: '#2dd4bf',
    accentText: '#0284c7'
  },
  {
    id: 'mint',
    name: '浅绿色',
    primary: '#86efac',
    bg: '#f0fdf4',
    surface: '#ffffff',
    surfaceSoft: '#dcfce7',
    textPrimary: '#14532d',
    textSecondary: '#457a58',
    border: '#bbf7d0',
    accent: '#86efac',
    accentStrong: '#4ade80',
    accentText: '#166534'
  }
]

const THEME_KEY = 'aiweekly_theme'

export function getStoredTheme(): ThemeId {
  const stored = localStorage.getItem(THEME_KEY) as ThemeId | null
  if (stored && themes.some(t => t.id === stored)) {
    return stored
  }
  return 'cream'
}

export function saveTheme(themeId: ThemeId): void {
  localStorage.setItem(THEME_KEY, themeId)
}

export function applyTheme(themeId: ThemeId): void {
  const theme = themes.find(t => t.id === themeId)
  if (!theme) return

  document.documentElement.style.setProperty('--theme-primary', theme.primary)
  document.documentElement.style.setProperty('--theme-bg', theme.bg)
  document.documentElement.style.setProperty('--theme-surface', theme.surface)
  document.documentElement.style.setProperty('--theme-surface-soft', theme.surfaceSoft)
  document.documentElement.style.setProperty('--theme-text-primary', theme.textPrimary)
  document.documentElement.style.setProperty('--theme-text-secondary', theme.textSecondary)
  document.documentElement.style.setProperty('--theme-border', theme.border)
  document.documentElement.style.setProperty('--theme-accent', theme.accent)
  document.documentElement.style.setProperty('--theme-accent-strong', theme.accentStrong)
  document.documentElement.style.setProperty('--theme-accent-text', theme.accentText)
}

export function initTheme(): void {
  const themeId = getStoredTheme()
  applyTheme(themeId)
}
