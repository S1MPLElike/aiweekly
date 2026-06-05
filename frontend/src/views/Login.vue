<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { saveAuth } from '../services/auth'

type AuthMode = 'login' | 'register'

const router = useRouter()

const mode = ref<AuthMode>('login')
const globalMessage = ref('')
const isSubmitting = ref(false)

const loginForm = reactive({
  phone: '',
  password: '',
  remember: true,
})

const registerForm = reactive({
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
})

const loginError = reactive({
  phone: '',
  password: '',
})

const registerError = reactive({
  username: '',
  phone: '',
  password: '',
  confirmPassword: '',
})

const modeTitle = computed(() =>
  mode.value === 'login' ? '欢迎回来' : '创建你的账号',
)

const modeSubTitle = computed(() =>
  mode.value === 'login'
    ? '登录后即可继续记录工作并生成 AI 周报。'
    : '注册后即可开始你的第一条工作记录。',
)

function switchMode(nextMode: AuthMode) {
  mode.value = nextMode
  globalMessage.value = ''
  clearErrors()
}

function clearErrors() {
  loginError.phone = ''
  loginError.password = ''
  registerError.username = ''
  registerError.phone = ''
  registerError.password = ''
  registerError.confirmPassword = ''
}

function validateLogin() {
  clearErrors()
  let valid = true

  if (!loginForm.phone.trim()) {
    loginError.phone = '请输入手机号'
    valid = false
  }

  if (loginForm.phone.trim() && loginForm.phone.trim().length !== 11) {
    loginError.phone = '手机号必须为11位'
    valid = false
  }

  if (!loginForm.password.trim()) {
    loginError.password = '请输入密码'
    valid = false
  }

  return valid
}

function validateRegister() {
  clearErrors()
  let valid = true

  if (!registerForm.username.trim()) {
    registerError.username = '请输入昵称'
    valid = false
  }

  if (!registerForm.phone.trim()) {
    registerError.phone = '请输入手机号'
    valid = false
  }

  if (registerForm.phone.trim() && registerForm.phone.trim().length !== 11) {
    registerError.phone = '手机号必须为11位'
    valid = false
  }

  if (registerForm.password.length < 6) {
    registerError.password = '密码至少 6 位'
    valid = false
  }

  if (registerForm.confirmPassword !== registerForm.password) {
    registerError.confirmPassword = '两次密码输入不一致'
    valid = false
  }

  return valid
}

async function submitLogin() {
  if (!validateLogin()) return
  isSubmitting.value = true
  globalMessage.value = ''

  try {
    const response = await fetch('http://localhost:8082/api/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: JSON.stringify(loginForm),
    })

    const data = await response.json()
    if (response.ok && data.code === 0) {
      saveAuth({
        access_token: data.data.access_token,
        refresh_token: data.data.refresh_token,
        user: data.data.user
      })
      globalMessage.value = '登录成功！'
      router.push('/dashboard')
    } else {
      globalMessage.value = data.msg || '登录失败'
    }
  } catch (error) {
    globalMessage.value = '网络错误，请稍后重试'
    console.error('Login error:', error)
  } finally {
    isSubmitting.value = false
  }
}

async function submitRegister() {
  if (!validateRegister()) return
  isSubmitting.value = true
  globalMessage.value = ''

  try {
    const registerData = {
      username: registerForm.username,
      phone: registerForm.phone,
      password: registerForm.password
    }

    const response = await fetch('http://localhost:8082/api/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(registerData),
    })

    const data = await response.json()
    if (response.ok) {
      globalMessage.value = '注册成功！'
      setTimeout(() => {
        switchMode('login')
      }, 500)
    } else {
      globalMessage.value = data.msg || '注册失败'
    }
  } catch (error) {
    globalMessage.value = '网络错误，请稍后重试'
    console.error('Register error:', error)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <main :class="$style.page">
    <section :class="$style.visualPanel" aria-label="图片展示区域">
      <header :class="$style.visualHeader">
        <h1 :class="$style.brand">冲击波智能工作周报</h1>
      </header>
      <div :class="$style.visualBody">
        <div :class="$style.visualInner" role="img" aria-label="登录页预留图片容器">
          <div :class="$style.visualImage"></div>
        </div>
      </div>
    </section>

    <section :class="$style.formPanel" aria-label="用户登录与注册">
      <div :class="$style.formCard">
        <header :class="$style.formHeader">
          <p :class="$style.kicker">AI Weekly Assistant</p>
          <h2 :class="$style.title">{{ modeTitle }}</h2>
          <p :class="$style.subtitle">{{ modeSubTitle }}</p>
        </header>

        <div :class="$style.tabs" role="tablist" aria-label="登录注册切换">
          <button
            type="button"
            role="tab"
            :aria-selected="mode === 'login'"
            :class="[$style.tabButton, mode === 'login' && $style.tabActive]"
            @click="switchMode('login')"
          >
            登录
          </button>
          <button
            type="button"
            role="tab"
            :aria-selected="mode === 'register'"
            :class="[$style.tabButton, mode === 'register' && $style.tabActive]"
            @click="switchMode('register')"
          >
            注册
          </button>
        </div>

        <form
          v-if="mode === 'login'"
          :class="$style.formBody"
          @submit.prevent="submitLogin"
        >
          <label :class="$style.field">
            <span :class="$style.label">手机</span>
            <input
              v-model="loginForm.phone"
              :class="$style.input"
              type="tel"
              name="login-phone"
              autocomplete="tel"
              placeholder="请输入手机号"
            />
            <small v-if="loginError.phone" :class="$style.errorText">{{ loginError.phone }}</small>
          </label>

          <label :class="$style.field">
            <span :class="$style.label">密码</span>
            <input
              v-model="loginForm.password"
              :class="$style.input"
              type="password"
              name="login-password"
              autocomplete="current-password"
              placeholder="请输入密码"
            />
            <small v-if="loginError.password" :class="$style.errorText">{{ loginError.password }}</small>
          </label>

          <label :class="$style.remember">
            <input v-model="loginForm.remember" type="checkbox" />
            <span>7 天内记住登录状态</span>
          </label>

          <button :class="$style.primaryButton" :disabled="isSubmitting" type="submit">
            {{ isSubmitting ? '登录中...' : '登录' }}
          </button>
        </form>

        <form
          v-else
          :class="$style.formBody"
          @submit.prevent="submitRegister"
        >
          <label :class="$style.field">
            <span :class="$style.label">昵称</span>
            <input
              v-model="registerForm.username"
              :class="$style.input"
              type="text"
              name="register-username"
              autocomplete="username"
              placeholder="请输入昵称"
            />
            <small v-if="registerError.username" :class="$style.errorText">{{ registerError.username }}</small>
          </label>

          <label :class="$style.field">
            <span :class="$style.label">手机</span>
            <input
              v-model="registerForm.phone"
              :class="$style.input"
              type="tel"
              name="register-phone"
              autocomplete="tel"
              placeholder="请输入手机号"
            />
            <small v-if="registerError.phone" :class="$style.errorText">{{ registerError.phone }}</small>
          </label>

          <label :class="$style.field">
            <span :class="$style.label">密码</span>
            <input
              v-model="registerForm.password"
              :class="$style.input"
              type="password"
              name="register-password"
              autocomplete="new-password"
              placeholder="至少 6 位密码"
            />
            <small v-if="registerError.password" :class="$style.errorText">{{ registerError.password }}</small>
          </label>

          <label :class="$style.field">
            <span :class="$style.label">确认密码</span>
            <input
              v-model="registerForm.confirmPassword"
              :class="$style.input"
              type="password"
              name="register-confirm-password"
              autocomplete="new-password"
              placeholder="请再次输入密码"
            />
            <small v-if="registerError.confirmPassword" :class="$style.errorText">{{ registerError.confirmPassword }}</small>
          </label>

          <button :class="$style.primaryButton" :disabled="isSubmitting" type="submit">
            {{ isSubmitting ? '提交中...' : '注册账号' }}
          </button>
        </form>

        <p v-if="globalMessage" :class="$style.statusMessage" aria-live="polite">
          {{ globalMessage }}
        </p>
      </div>
    </section>
  </main>
</template>

<style module>
.page {
  --bg: #fffcf3;
  --surface: #ffffff;
  --surface-soft: #fff8df;
  --text-primary: #2f2a22;
  --text-secondary: #766e60;
  --border: #f2e8c8;
  --accent: #f5cd6a;
  --accent-strong: #e8b947;
  --error: #ca5d49;
  min-height: 100vh;
  display: grid;
  grid-template-columns: 2fr 1fr;
  background: var(--bg);
  color: var(--text-primary);
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.visualPanel {
  position: relative;
  padding: 48px;
  background:
    radial-gradient(circle at 20% 20%, #fff1bf 0%, transparent 45%),
    radial-gradient(circle at 80% 80%, #ffe8a8 0%, transparent 40%),
    linear-gradient(135deg, #fffef9 0%, #fff7dc 100%);
  border-right: 1px solid var(--border);
  display: grid;
  grid-template-rows: auto 1fr;
  gap: 18px;
}

.visualHeader {
  text-align: left;
}

.visualBody {
  display: grid;
  place-items: center;
}

.visualInner {
  width: min(520px, 92%);
  min-height: 360px;
  border: 1px dashed #e7d59b;
  border-radius: 24px;
  background-color: rgb(255 255 255 / 62%);
  overflow: hidden;
  padding: 0;
}

.brand {
  font-size: 2rem;
  line-height: 1.3;
  margin: 0;
}

.visualImage {
  width: 100%;
  height: 100%;
  min-height: 360px;
  border-radius: 24px;
  background-color: #f5ecd2;
  background-image: url('/auth/login.jpg');
  background-position: center;
  background-size: cover;
  background-repeat: no-repeat;
}

.formPanel {
  display: grid;
  place-items: center;
  padding: 32px 24px;
}

.formCard {
  width: min(420px, 100%);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 28px 24px;
  box-shadow: 0 16px 28px rgb(66 55 20 / 8%);
}

.formHeader {
  margin-bottom: 20px;
}

.kicker {
  margin: 0;
  font-size: 0.8rem;
  letter-spacing: 0.08em;
  color: var(--text-secondary);
  text-transform: uppercase;
}

.title {
  margin: 8px 0 4px;
  font-size: 1.6rem;
}

.subtitle {
  margin: 0;
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 16px;
}

.tabButton {
  border: 0;
  background: #fffdf5;
  color: #796f5f;
  height: 42px;
  font-size: 0.95rem;
  cursor: pointer;
}

.tabButton:focus-visible,
.input:focus-visible,
.primaryButton:focus-visible {
  outline: 3px solid #ffda76;
  outline-offset: 2px;
}

.tabActive {
  background: var(--surface-soft);
  color: #4d402a;
  font-weight: 600;
}

.formBody {
  display: grid;
  gap: 14px;
}

.field {
  display: grid;
  gap: 7px;
}

.label {
  font-size: 0.9rem;
}

.input {
  border: 1px solid #ebddb3;
  border-radius: 10px;
  height: 42px;
  padding: 0 12px;
  font-size: 0.95rem;
  background: #fffef8;
  transition: border-color 0.2s ease;
}

.input:hover {
  border-color: #e3cb85;
}

.remember {
  display: inline-flex;
  gap: 8px;
  align-items: center;
  color: var(--text-secondary);
  font-size: 0.88rem;
}

.primaryButton {
  border: 0;
  border-radius: 12px;
  height: 44px;
  background: linear-gradient(135deg, var(--accent) 0%, #f6d47c 100%);
  color: #463500;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.18s ease, filter 0.18s ease;
}

.primaryButton:hover {
  transform: translateY(-1px);
  filter: brightness(0.98);
}

.primaryButton:disabled {
  cursor: not-allowed;
  opacity: 0.7;
  transform: none;
}

.errorText {
  color: var(--error);
  font-size: 0.82rem;
}

.statusMessage {
  margin-top: 16px;
  font-size: 0.9rem;
  color: #6a560f;
  background: #fff7db;
  border: 1px solid #f1dfaa;
  border-radius: 10px;
  padding: 10px 12px;
}

@media (max-width: 1080px) {
  .page {
    grid-template-columns: 1fr;
  }

  .visualPanel {
    min-height: 240px;
    border-right: 0;
    border-bottom: 1px solid var(--border);
    padding: 28px 20px;
  }

  .brand {
    font-size: 1.5rem;
  }

  .visualInner {
    min-height: 240px;
  }

  .visualImage {
    min-height: 240px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .input,
  .primaryButton {
    transition: none;
  }
}
</style>
