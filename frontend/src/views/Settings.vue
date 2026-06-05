<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { themes, getStoredTheme, saveTheme, applyTheme } from '../services/theme'
import { getUserSetting, updateUserSetting, type UserSetting } from '../api/work'

const activeTab = ref('account')
const isLoadingSettings = ref(false)
const isSavingSettings = ref(false)

const accountForm = reactive({
  username: '',
  phone: '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const aiPreferences = reactive({
  tone: 'professional' as 'professional' | 'relaxed' | 'concise',
  includeWeekends: false,
  autoGenerate: true
})

const currentTheme = ref(getStoredTheme())

const toneOptions = [
  { value: 'professional' as const, label: '专业严谨', desc: '适合正式工作汇报' },
  { value: 'relaxed' as const, label: '轻松活泼', desc: '适合团队内部分享' },
  { value: 'concise' as const, label: '简洁概括', desc: '适合快速浏览' }
]

function selectTheme(themeId: string) {
  currentTheme.value = themeId
  saveTheme(themeId)
  applyTheme(themeId)
}

async function loadSettings() {
  isLoadingSettings.value = true
  try {
    const setting = await getUserSetting()
    aiPreferences.tone = setting.weekly_report_style
  } catch (error) {
    console.error('加载设置失败:', error)
  } finally {
    isLoadingSettings.value = false
  }
}

onMounted(() => {
  currentTheme.value = getStoredTheme()
  loadSettings()
})

function updateAccount() {
  console.log('更新账户信息:', accountForm)
}

async function savePreferences() {
  isSavingSettings.value = true
  try {
    await updateUserSetting(aiPreferences.tone)
    alert('偏好已保存！')
  } catch (error) {
    console.error('保存设置失败:', error)
    alert('保存失败，请重试')
  } finally {
    isSavingSettings.value = false
  }
}
</script>

<template>
  <div :class="$style.container">
    <div :class="$style.sidebar">
      <h2 :class="$style.sidebarTitle">设置</h2>
      <nav :class="$style.nav">
        <button
          type="button"
          :class="[$style.navItem, activeTab === 'account' && $style.navItemActive]"
          @click="activeTab = 'account'"
        >
          <svg :class="$style.navIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
          <span>账户信息</span>
        </button>
        <button
          type="button"
          :class="[$style.navItem, activeTab === 'ai' && $style.navItemActive]"
          @click="activeTab = 'ai'"
        >
          <svg :class="$style.navIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
          </svg>
          <span>AI周报偏好</span>
        </button>
        <button
          type="button"
          :class="[$style.navItem, activeTab === 'theme' && $style.navItemActive]"
          @click="activeTab = 'theme'"
        >
          <svg :class="$style.navIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3" />
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82v.09a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
          </svg>
          <span>主题外观</span>
        </button>
      </nav>
    </div>

    <div :class="$style.content">
      <div v-if="activeTab === 'account'" :class="$style.section">
        <div :class="$style.sectionHeader">
          <h3 :class="$style.sectionTitle">账户信息</h3>
          <p :class="$style.sectionDesc">管理您的账户基本信息</p>
        </div>

        <div :class="$style.formGroup">
          <label :class="$style.formLabel">用户名</label>
          <input
            v-model="accountForm.username"
            type="text"
            :class="$style.formInput"
            placeholder="请输入用户名"
          />
        </div>

        <div :class="$style.formGroup">
          <label :class="$style.formLabel">手机号</label>
          <input
            v-model="accountForm.phone"
            type="tel"
            :class="$style.formInput"
            placeholder="请输入手机号"
          />
        </div>

        <div :class="$style.divider"></div>

        <div :class="$style.sectionHeader">
          <h3 :class="$style.sectionTitle">修改密码</h3>
          <p :class="$style.sectionDesc">请输入当前密码和新密码</p>
        </div>

        <div :class="$style.formGroup">
          <label :class="$style.formLabel">当前密码</label>
          <input
            v-model="accountForm.currentPassword"
            type="password"
            :class="$style.formInput"
            placeholder="请输入当前密码"
          />
        </div>

        <div :class="$style.formGroup">
          <label :class="$style.formLabel">新密码</label>
          <input
            v-model="accountForm.newPassword"
            type="password"
            :class="$style.formInput"
            placeholder="请输入新密码"
          />
        </div>

        <div :class="$style.formGroup">
          <label :class="$style.formLabel">确认新密码</label>
          <input
            v-model="accountForm.confirmPassword"
            type="password"
            :class="$style.formInput"
            placeholder="请再次输入新密码"
          />
        </div>

        <button :class="$style.submitButton" @click="updateAccount">
          保存更改
        </button>
      </div>

      <div v-if="activeTab === 'ai'" :class="$style.section">
        <div :class="$style.sectionHeader">
          <h3 :class="$style.sectionTitle">AI周报偏好</h3>
          <p :class="$style.sectionDesc">自定义AI生成周报的风格和内容</p>
        </div>

        <div :class="$style.formGroup">
          <label :class="$style.formLabel">生成语气</label>
          <div :class="$style.toneOptions">
            <label
              v-for="option in toneOptions"
              :key="option.value"
              :class="[$style.toneOption, aiPreferences.tone === option.value && $style.toneOptionActive]"
            >
              <input
                type="radio"
                v-model="aiPreferences.tone"
                :value="option.value"
                :class="$style.radioInput"
              />
              <div :class="$style.toneContent">
                <span :class="$style.toneLabel">{{ option.label }}</span>
                <span :class="$style.toneDesc">{{ option.desc }}</span>
              </div>
            </label>
          </div>
        </div>

        <div :class="$style.formGroup">
          <div :class="$style.checkboxOption">
            <input
              v-model="aiPreferences.includeWeekends"
              type="checkbox"
              :class="$style.checkboxInput"
            />
            <label :class="$style.checkboxLabel">包含周末工作记录</label>
          </div>
        </div>

        <div :class="$style.formGroup">
          <div :class="$style.checkboxOption">
            <input
              v-model="aiPreferences.autoGenerate"
              type="checkbox"
              :class="$style.checkboxInput"
            />
            <label :class="$style.checkboxLabel">每周自动生成周报</label>
          </div>
        </div>

        <button :class="$style.submitButton" @click="savePreferences" :disabled="isSavingSettings">
          {{ isSavingSettings ? '保存中...' : '保存偏好' }}
        </button>
      </div>

      <div v-if="activeTab === 'theme'" :class="$style.section">
        <div :class="$style.sectionHeader">
          <h3 :class="$style.sectionTitle">主题外观</h3>
          <p :class="$style.sectionDesc">选择您喜欢的界面主题</p>
        </div>

        <div :class="$style.themeGrid">
          <button
            v-for="theme in themes"
            :key="theme.id"
            type="button"
            :class="[$style.themeCard, currentTheme === theme.id && $style.themeCardActive]"
            @click="selectTheme(theme.id)"
          >
            <div :class="$style.themePreview" :style="{ background: `linear-gradient(135deg, ${theme.primary} 0%, ${theme.accent} 100%)` }"></div>
            <span :class="$style.themeName">{{ theme.name }}</span>
            <svg v-if="currentTheme === theme.id" :class="$style.checkIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style module>
.container {
  display: flex;
  height: 100%;
  background: var(--theme-bg);
}

.sidebar {
  width: 220px;
  background: var(--theme-surface);
  padding: 24px 0;
  border-right: 1px solid var(--theme-border);
}

.sidebarTitle {
  padding: 0 24px 20px;
  margin: 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--theme-text-primary);
}

.nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 12px;
}

.navItem {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border: none;
  border-radius: 10px;
  background: transparent;
  font-size: 0.9375rem;
  color: var(--theme-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.navItem:hover {
  background: var(--theme-surface-soft);
  color: var(--theme-text-primary);
}

.navItemActive {
  background: var(--theme-primary);
  color: var(--theme-accent-text);
  font-weight: 600;
}

.navIcon {
  width: 18px;
  height: 18px;
}

.content {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
}

.section {
  max-width: 600px;
}

.sectionHeader {
  margin-bottom: 24px;
}

.sectionTitle {
  margin: 0 0 6px 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--theme-text-primary);
}

.sectionDesc {
  margin: 0;
  font-size: 0.9375rem;
  color: var(--theme-text-secondary);
}

.formGroup {
  margin-bottom: 20px;
}

.formLabel {
  display: block;
  margin-bottom: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--theme-text-primary);
}

.formInput {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--theme-border);
  border-radius: 10px;
  background: var(--theme-surface);
  font-size: 0.9375rem;
  color: var(--theme-text-primary);
  transition: border-color 0.2s ease;
  box-sizing: border-box;
}

.formInput:focus {
  outline: none;
  border-color: var(--theme-primary);
}

.formInput::placeholder {
  color: var(--theme-text-secondary);
  opacity: 0.6;
}

.divider {
  height: 1px;
  background: var(--theme-border);
  margin: 32px 0;
}

.submitButton {
  padding: 12px 28px;
  border: none;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-accent-strong) 100%);
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--theme-accent-text);
  cursor: pointer;
  transition: all 0.2s ease;
}

.submitButton:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(245, 205, 106, 0.3);
}

.toneOptions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.toneOption {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 2px solid var(--theme-border);
  border-radius: 10px;
  background: var(--theme-surface);
  cursor: pointer;
  transition: all 0.2s ease;
}

.toneOption:hover {
  border-color: var(--theme-primary);
}

.toneOptionActive {
  border-color: var(--theme-primary);
  background: var(--theme-surface-soft);
}

.radioInput {
  width: 18px;
  height: 18px;
  accent-color: var(--theme-primary);
}

.toneContent {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.toneLabel {
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--theme-text-primary);
}

.toneDesc {
  font-size: 0.8125rem;
  color: var(--theme-text-secondary);
}

.checkboxOption {
  display: flex;
  align-items: center;
  gap: 10px;
}

.checkboxInput {
  width: 18px;
  height: 18px;
  accent-color: var(--theme-primary);
}

.checkboxLabel {
  font-size: 0.9375rem;
  color: var(--theme-text-primary);
  cursor: pointer;
}

.themeGrid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

.themeCard {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 20px;
  border: 2px solid var(--theme-border);
  border-radius: 14px;
  background: var(--theme-surface);
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.themeCard:hover {
  border-color: var(--theme-primary);
  transform: translateY(-2px);
}

.themeCardActive {
  border-color: var(--theme-primary);
  background: var(--theme-surface-soft);
}

.themePreview {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

.themeName {
  font-size: 0.9375rem;
  font-weight: 500;
  color: #2f2a22;
}

.checkIcon {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 20px;
  height: 20px;
  color: var(--theme-primary);
}
</style>
