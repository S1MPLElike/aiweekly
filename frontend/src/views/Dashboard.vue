<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { clearAuth, logout } from '../services/auth'
import { getMonthlyRecordStats, type MonthlyRecordStats } from '../api/work'
import TodayWork from './TodayWork.vue'
import DayWorkDialog from './DayWorkDialog.vue'
import MonthlyLedger from './MonthlyLedger.vue'
import Settings from './Settings.vue'
import Profile from './Profile.vue'

const router = useRouter()

type ModuleType = 'yesterday' | 'today' | 'calendar' | 'settings' | 'profile'

const currentModule = ref<ModuleType>('today')

const moduleList = [
  { id: 'yesterday', label: '昨日工作', icon: '📅' },
  { id: 'today', label: '开启今日工作', icon: '✍️' },
  { id: 'calendar', label: '月台账', icon: '📆' },
  { id: 'settings', label: '设置', icon: '⚙️' },
  { id: 'profile', label: '我', icon: '👤' },
] as const
const currentModuleLabel = computed(() => {
  const module = moduleList.find(m => m.id === currentModule.value)
  return module?.label || ''
})

const today = ref(new Date())

const yesterdayDate = computed(() => {
  const yesterday = new Date(today.value)
  yesterday.setDate(yesterday.getDate() - 1)
  return `${yesterday.getFullYear()}-${String(yesterday.getMonth() + 1).padStart(2, '0')}-${String(yesterday.getDate()).padStart(2, '0')}`
})

function switchModule(moduleId: ModuleType) {
  currentModule.value = moduleId
}

const dialogVisible = ref(false)
const dialogDate = ref('')

function openDayDialog(dateStr: string) {
  dialogDate.value = dateStr
  dialogVisible.value = true
}

function closeDayDialog() {
  dialogVisible.value = false
}

// 日历相关逻辑
const registrationDate = ref(new Date(2026, 3, 1))

const now = new Date()
// 默认显示上个月，以便查看有工作记录的日期
let defaultYear = now.getFullYear()
let defaultMonth = now.getMonth() - 1
if (defaultMonth < 0) {
  defaultMonth = 11
  defaultYear--
}
const currentYear = ref(defaultYear)
const currentMonth = ref(defaultMonth)

const minYear = computed(() => registrationDate.value.getFullYear())
const minMonth = computed(() => registrationDate.value.getMonth())

const maxYear = computed(() => today.value.getFullYear() + 1)
const maxMonth = computed(() => today.value.getMonth())

const canPrevMonth = computed(() => {
  if (currentYear.value > minYear.value) return true
  if (currentYear.value === minYear.value && currentMonth.value > minMonth.value) return true
  return false
})

// 月度工作记录统计
const monthlyStats = ref<MonthlyRecordStats[]>([])
const isLoadingStats = ref(false)

function getRecordCount(dateStr: string): number {
  const stat = monthlyStats.value.find(s => s.date === dateStr)
  console.log(`查找日期 ${dateStr}, 找到数据:`, stat)
  return stat?.count || 0
}

function getDotColor(count: number): string {
  if (count >= 7) return 'red'
  if (count >= 4) return 'orange'
  if (count >= 1) return 'green'
  return ''
}

function getDotBackgroundColor(color: string): string {
  switch (color) {
    case 'red':
      return '#ff4d4f'
    case 'orange':
      return '#faad14'
    case 'green':
      return '#52c41a'
    default:
      return ''
  }
}

async function fetchMonthlyStats() {
  isLoadingStats.value = true
  try {
    const data = await getMonthlyRecordStats(currentYear.value, currentMonth.value + 1)
    monthlyStats.value = data
  } catch (error) {
    console.error('获取月度统计失败:', error)
  } finally {
    isLoadingStats.value = false
  }
}

const canNextMonth = computed(() => {
  if (currentYear.value < maxYear.value) return true
  if (currentYear.value === maxYear.value && currentMonth.value < maxMonth.value) return true
  return false
})

const firstDayOfMonth = computed(() => {
  const firstDay = new Date(currentYear.value, currentMonth.value, 1)
  let day = firstDay.getDay()
  return day === 0 ? 6 : day - 1
})

const daysInMonth = computed(() => {
  return new Date(currentYear.value, currentMonth.value + 1, 0).getDate()
})

const calendarDays = computed(() => {
  const days = []
  for (let i = 0; i < firstDayOfMonth.value; i++) {
    days.push(null)
  }
  for (let i = 1; i <= daysInMonth.value; i++) {
    const date = new Date(currentYear.value, currentMonth.value, i)
    const dateStr = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}-${String(i).padStart(2, '0')}`
    const recordCount = getRecordCount(dateStr)
    days.push({
      day: i,
      isToday: date.toDateString() === today.value.toDateString(),
      recordCount,
      dotColor: getDotColor(recordCount)
    })
  }
  return days
})

const monthNames = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']
const currentMonthName = computed(() => `${currentYear.value}年${monthNames[currentMonth.value]}`)

function prevMonth() {
  if (!canPrevMonth.value) return
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
  fetchMonthlyStats()
}

function nextMonth() {
  if (!canNextMonth.value) return
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
  fetchMonthlyStats()
}

function handleDateClick(day: { day: number }) {
  const dateStr = `${currentYear.value}-${String(currentMonth.value + 1).padStart(2, '0')}-${String(day.day).padStart(2, '0')}`
  openDayDialog(dateStr)
}

async function handleLogout() {
  try {
    await logout()
  } finally {
    clearAuth()
    router.push('/login')
  }
}

onMounted(() => {
  fetchMonthlyStats()
})
</script>

<template>
  <main :class="$style.dashboard">
    <aside :class="$style.sidebar" aria-label="模块导航">
      <div :class="$style.sidebarHeader">
        <h1 :class="$style.brand">冲击波智能工作周报</h1>
      </div>
      <nav :class="$style.nav" role="navigation" aria-label="功能模块">
        <ul :class="$style.navList" role="list">
          <li v-for="module in moduleList" :key="module.id">
            <button
              type="button"
              role="menuitem"
              :class="[$style.navButton, currentModule === module.id && $style.navButtonActive]"
              :aria-current="currentModule === module.id ? 'page' : undefined"
              @click="switchModule(module.id as ModuleType)"
            >
              <span :class="$style.navIcon" aria-hidden="true">{{ module.icon }}</span>
              <span :class="$style.navLabel">{{ module.label }}</span>
            </button>
          </li>
        </ul>
      </nav>

      <div v-if="currentModule !== 'calendar'" :class="$style.calendarPanel">
        <div :class="$style.calendarHeader">
          <button
            type="button"
            :class="[$style.calendarNavButton, !canPrevMonth && $style.calendarNavButtonDisabled]"
            :disabled="!canPrevMonth"
            @click="prevMonth"
            aria-label="上一个月"
          >
            ◀
          </button>
          <span :class="$style.calendarTitle">{{ currentMonthName }}</span>
          <button
            type="button"
            :class="[$style.calendarNavButton, !canNextMonth && $style.calendarNavButtonDisabled]"
            :disabled="!canNextMonth"
            @click="nextMonth"
            aria-label="下一个月"
          >
            ▶
          </button>
        </div>
        <div :class="$style.calendarWeekdays">
          <span :class="$style.weekday">一</span>
          <span :class="$style.weekday">二</span>
          <span :class="$style.weekday">三</span>
          <span :class="$style.weekday">四</span>
          <span :class="$style.weekday">五</span>
          <span :class="$style.weekday">六</span>
          <span :class="$style.weekday">日</span>
        </div>
        <div :class="$style.calendarGrid">
          <button
            v-for="(day, index) in calendarDays"
            :key="index"
            type="button"
            :class="[$style.calendarDay, day && day.isToday && $style.calendarDayToday]"
            :disabled="!day"
            @click="day && handleDateClick(day)"
            :aria-label="day ? `${currentYear}年${currentMonth + 1}月${day.day}日${day.recordCount ? `，${day.recordCount}条工作记录` : ''}` : undefined"
          >
            <span v-if="day">{{ day.day }}</span>
            <span
              v-if="day && day.dotColor"
              :class="$style.calendarDot"
              :style="{ backgroundColor: getDotBackgroundColor(day.dotColor) }"
            ></span>
          </button>
        </div>
      </div>

      <button
        type="button"
        :class="$style.logoutButton"
        @click="handleLogout"
        aria-label="登出"
      >
        <span :class="$style.logoutIcon" aria-hidden="true">🚪</span>
        <span :class="$style.logoutText">登出</span>
      </button>

      <div :class="$style.sidebarFooter">
        <p :class="$style.footerText">AI Weekly Assistant</p>
      </div>
    </aside>

    <section :class="$style.content" aria-label="模块内容区域">
      <template v-if="currentModule === 'today'">
        <TodayWork />
      </template>
      <template v-else-if="currentModule === 'yesterday'">
        <TodayWork :date="yesterdayDate" />
      </template>
      <template v-else-if="currentModule === 'calendar'">
        <MonthlyLedger :registrationDate="registrationDate" @dateClick="openDayDialog" />
      </template>
      <template v-else-if="currentModule === 'settings'">
        <Settings />
      </template>
      <template v-else-if="currentModule === 'profile'">
        <Profile />
      </template>
      <template v-else>
        <header :class="$style.contentHeader">
          <h2 :class="$style.contentTitle">{{ currentModuleLabel }}</h2>
        </header>
        <div :class="$style.contentBody">
          <div :class="$style.placeholder">
            <div :class="$style.placeholderIcon" aria-hidden="true">📋</div>
            <p :class="$style.placeholderText">{{ currentModuleLabel }}模块开发中...</p>
            <p :class="$style.placeholderHint">敬请期待更多功能</p>
          </div>
        </div>
      </template>
    </section>

    <DayWorkDialog
      :visible="dialogVisible"
      :date="dialogDate"
      @close="closeDayDialog"
    />
  </main>
</template>

<style module>
.dashboard {
  --sidebar-width: 20%;
  min-height: 100vh;
  display: grid;
  grid-template-columns: var(--sidebar-width) 1fr;
  background: var(--theme-bg);
  color: var(--theme-text-primary);
  font-family: 'Segoe UI', 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.sidebar {
  background: var(--theme-surface);
  border-right: 1px solid var(--theme-border);
  display: grid;
  grid-template-rows: auto auto 1fr auto auto;
  padding: 24px 16px;
  gap: 24px;
  box-shadow: 4px 0 12px rgb(66 55 20 / 4%);
}

.sidebarHeader {
  padding-bottom: 16px;
  border-bottom: 1px solid var(--theme-border);
}

.brand {
  font-size: 1.1rem;
  line-height: 1.4;
  margin: 0;
  color: var(--theme-text-primary);
}

.nav {
  display: flex;
  flex-direction: column;
}

.navList {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.navButton {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border: 1px solid transparent;
  border-radius: 12px;
  background: transparent;
  color: var(--theme-text-secondary);
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
}

.navButton:hover {
  background: var(--theme-surface-soft);
  color: var(--theme-text-primary);
  border-color: var(--theme-border);
}

.navButton:focus-visible {
  outline: 3px solid #ffda76;
  outline-offset: 2px;
}

.navButtonActive {
  background: var(--theme-accent);
  color: #463500;
  font-weight: 600;
  border-color: var(--theme-accent-strong);
}

.navButtonActive:hover {
  background: var(--theme-accent-strong);
  color: #463500;
}

.navIcon {
  font-size: 1.2rem;
  flex-shrink: 0;
}

.navLabel {
  flex: 1;
}

.calendarPanel {
  background: var(--theme-surface-soft);
  border-radius: 16px;
  padding: 16px;
  border: 1px solid var(--theme-border);
}

.calendarHeader {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.calendarNavButton {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 8px;
  background: var(--theme-surface);
  color: var(--theme-text-primary);
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.calendarNavButton:hover:not(:disabled) {
  background: var(--theme-accent);
}

.calendarNavButton:focus-visible {
  outline: 2px solid #ffda76;
  outline-offset: 2px;
}

.calendarNavButtonDisabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.calendarTitle {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--theme-text-primary);
}

.calendarWeekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: 8px;
}

.weekday {
  text-align: center;
  font-size: 0.75rem;
  color: var(--theme-text-secondary);
  padding: 4px 0;
}

.calendarGrid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.calendarDay {
  aspect-ratio: 1;
  border: none;
  border-radius: 8px;
  background: var(--theme-surface);
  color: var(--theme-text-primary);
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 2px;
}

.calendarDay:hover:not(:disabled) {
  background: var(--theme-accent);
}

.calendarDay:focus-visible {
  outline: 2px solid #ffda76;
  outline-offset: 1px;
}

.calendarDay:disabled {
  opacity: 0.3;
  cursor: default;
}

.calendarDayToday {
  background: var(--theme-accent);
  font-weight: 600;
}

.calendarDot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.logoutButton {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border: 1px solid var(--theme-border);
  border-radius: 12px;
  background: transparent;
  color: #c53030;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logoutButton:hover {
  background: #fff5f5;
  border-color: #fc8181;
}

.logoutButton:focus-visible {
  outline: 3px solid #fc8181;
  outline-offset: 2px;
}

.logoutIcon {
  font-size: 1rem;
}

.logoutText {
  font-weight: 500;
}

.sidebarFooter {
  padding-top: 16px;
  border-top: 1px solid var(--theme-border);
}

.footerText {
  margin: 0;
  font-size: 0.75rem;
  letter-spacing: 0.08em;
  color: var(--theme-text-secondary);
  text-transform: uppercase;
}

.content {
  display: flex;
  flex-direction: column;
  padding: 32px 40px;
}

.contentHeader {
  margin-bottom: 28px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--theme-border);
}

.contentTitle {
  margin: 0;
  font-size: 1.6rem;
  color: var(--theme-text-primary);
}

.contentBody {
  flex: 1;
  display: flex;
  align-items: stretch;
  overflow: hidden;
}

.placeholder {
  text-align: center;
  padding: 48px;
  background: var(--theme-surface);
  border: 1px dashed var(--theme-border);
  border-radius: 20px;
  max-width: 480px;
}

.placeholderIcon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.placeholderText {
  margin: 0 0 8px;
  font-size: 1.1rem;
  color: var(--theme-text-secondary);
}

.placeholderHint {
  margin: 0;
  font-size: 0.9rem;
  color: var(--theme-text-secondary);
  opacity: 0.7;
}

@media (max-width: 1024px) {
  .dashboard {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr;
  }

  .sidebar {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto auto;
    flex-direction: row;
    padding: 16px;
    gap: 16px;
    border-right: 0;
    border-bottom: 1px solid var(--theme-border);
    overflow-x: auto;
  }

  .sidebarHeader {
    padding-bottom: 0;
    border-bottom: 0;
    border-right: 1px solid var(--theme-border);
    padding-right: 16px;
  }

  .nav {
    flex-direction: row;
  }

  .navList {
    flex-direction: row;
    gap: 8px;
  }

  .navButton {
    padding: 10px 14px;
    white-space: nowrap;
  }

  .calendarPanel {
    display: none;
  }

  .sidebarFooter {
    padding-top: 0;
    border-top: 0;
    border-left: 1px solid var(--theme-border);
    padding-left: 16px;
    display: flex;
    align-items: center;
  }

  .content {
    padding: 24px 20px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .navButton,
  .calendarDay,
  .calendarNavButton {
    transition: none;
  }
}
</style>