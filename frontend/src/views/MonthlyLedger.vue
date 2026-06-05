<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import {
  getMonthlyRecordStats,
  getWeeklyReport,
  generateWeeklyReport,
  getMonthlySummary,
  type MonthlyRecordStats,
  type WeeklyReport,
  type MonthlySummary
} from '../api/work'

const props = defineProps<{
  registrationDate?: Date
  onDateClick?: (dateStr: string) => void
}>()

const emit = defineEmits<{
  (e: 'dateClick', dateStr: string): void
}>()

const currentYear = ref(new Date().getFullYear())
const currentMonth = ref(new Date().getMonth())

const today = new Date()
const registrationDate = props.registrationDate || new Date(2026, 3, 1)

const minYear = computed(() => registrationDate.getFullYear())
const minMonth = computed(() => registrationDate.getMonth())

const maxYear = computed(() => today.getFullYear() + 1)
const maxMonth = computed(() => today.getMonth())

const canPrevMonth = computed(() => {
  if (currentYear.value > minYear.value) return true
  if (currentYear.value === minYear.value && currentMonth.value > minMonth.value) return true
  return false
})

const canNextMonth = computed(() => {
  if (currentYear.value < maxYear.value) return true
  if (currentYear.value === maxYear.value && currentMonth.value < maxMonth.value) return true
  return false
})

const monthNames = ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月', '九月', '十月', '十一月', '十二月']
const currentMonthName = computed(() => `${currentYear.value}年${monthNames[currentMonth.value]}`)

const monthlyStats = ref<MonthlyRecordStats[]>([])
const monthlySummary = ref<MonthlySummary | null>(null)
const selectedWeeklyReport = ref<WeeklyReport | null>(null)
const selectedWeek = ref<{ weekStart: string; weekEnd: string } | null>(null)
const isGenerating = ref(false)

interface CalendarWeek {
  days: (CalendarDay | null)[]
  weekStart: string
  weekEnd: string
}

interface CalendarDay {
  day: number
  dateStr: string
  isCurrentMonth: boolean
  isToday: boolean
  recordCount: number
  dotColor: string
}

const calendarWeeks = computed<CalendarWeek[]>(() => {
  const weeks: CalendarWeek[] = []
  const firstDay = new Date(currentYear.value, currentMonth.value, 1)
  const lastDay = new Date(currentYear.value, currentMonth.value + 1, 0)

  let startDate = new Date(firstDay)
  const dayOfWeek = firstDay.getDay()
  const daysToMonday = (dayOfWeek + 6) % 7
  startDate.setDate(startDate.getDate() - daysToMonday)

  while (startDate <= lastDay || weeks.length === 0 || weeks[weeks.length - 1].days.length < 7) {
    const week: CalendarWeek = {
      days: [],
      weekStart: formatDate(startDate),
      weekEnd: ''
    }

    let weekEnd = new Date(startDate)
    weekEnd.setDate(weekEnd.getDate() + 6)
    week.weekEnd = formatDate(weekEnd)

    for (let i = 0; i < 7; i++) {
      const date = new Date(startDate)
      date.setDate(date.getDate() + i)

      const isCurrentMonth = date.getMonth() === currentMonth.value && date.getFullYear() === currentYear.value

      const dateStr = formatDate(date)
      const stat = monthlyStats.value.find(s => s.date === dateStr)
      const count = stat?.count || 0
      const dotColor = getDotColor(count)

      week.days.push({
        day: date.getDate(),
        dateStr,
        isCurrentMonth,
        isToday: date.toDateString() === today.toDateString(),
        recordCount: count,
        dotColor
      })
    }

    weeks.push(week)
    startDate.setDate(startDate.getDate() + 7)
  }

  return weeks
})

function formatDate(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function getDotColor(count: number): string {
  if (count >= 7) return '#ff4d4f'
  if (count >= 4) return '#faad14'
  if (count >= 1) return '#52c41a'
  return 'transparent'
}

function prevMonth() {
  if (!canPrevMonth.value) return
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

function nextMonth() {
  if (!canNextMonth.value) return
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

async function fetchData() {
  try {
    monthlyStats.value = await getMonthlyRecordStats(currentYear.value, currentMonth.value + 1)
    monthlySummary.value = await getMonthlySummary(currentYear.value, currentMonth.value + 1)
  } catch (err) {
    console.error('获取数据失败:', err)
  }
}

async function handleViewReport(weekStart: string, weekEnd: string) {
  try {
    const report = await getWeeklyReport(weekStart, weekEnd)
    selectedWeeklyReport.value = report
    selectedWeek.value = { weekStart, weekEnd }
  } catch (err) {
    console.error('获取周报失败:', err)
  }
}

async function handleGenerateReport(weekStart: string, weekEnd: string) {
  isGenerating.value = true
  try {
    const report = await generateWeeklyReport(weekStart, weekEnd)
    selectedWeeklyReport.value = report
    selectedWeek.value = { weekStart, weekEnd }
    monthlySummary.value = await getMonthlySummary(currentYear.value, currentMonth.value + 1)
  } catch (err) {
    console.error('生成周报失败:', err)
  } finally {
    isGenerating.value = false
  }
}

function handleDayClick(day: CalendarDay) {
  if (props.onDateClick) {
    props.onDateClick(day.dateStr)
  } else {
    emit('dateClick', day.dateStr)
  }
}

watch([currentYear, currentMonth], () => {
  fetchData()
  selectedWeeklyReport.value = null
  selectedWeek.value = null
})

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div :class="$style.container">
    <div :class="$style.header">
      <div :class="$style.headerLeft">
        <div :class="$style.monthNav">
          <button
            type="button"
            :class="[$style.navButton, !canPrevMonth && $style.navButtonDisabled]"
            :disabled="!canPrevMonth"
            @click="prevMonth"
            aria-label="上一月"
          >
            <svg :class="$style.chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <h2 :class="$style.monthTitle">{{ currentMonthName }}</h2>
          <button
            type="button"
            :class="[$style.navButton, !canNextMonth && $style.navButtonDisabled]"
            :disabled="!canNextMonth"
            @click="nextMonth"
            aria-label="下一月"
          >
            <svg :class="$style.chevron" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>

      <div :class="$style.headerRight">
        <div :class="$style.statsBar">
          <div :class="$style.statCard">
            <span :class="$style.statValue">{{ monthlySummary?.total_tasks || 0 }}</span>
            <span :class="$style.statLabel">工作任务</span>
          </div>
          <div :class="$style.statCard">
            <span :class="$style.statValue">{{ monthlySummary?.total_hours || 0 }}</span>
            <span :class="$style.statLabel">工作小时</span>
          </div>
          <div :class="$style.statCard">
            <span :class="$style.statValue">{{ monthlySummary?.week_count || 0 }}</span>
            <span :class="$style.statLabel">周数</span>
          </div>
          <div :class="$style.statCard">
            <span :class="$style.statValue">{{ monthlySummary?.report_count || 0 }}</span>
            <span :class="$style.statLabel">周报数</span>
          </div>
        </div>
      </div>
    </div>

    <div :class="$style.body">
      <div :class="$style.calendarPanel">
        <div :class="$style.calendarWrapper">
          <div :class="$style.weekdayHeader">
            <span :class="$style.weekday">一</span>
            <span :class="$style.weekday">二</span>
            <span :class="$style.weekday">三</span>
            <span :class="$style.weekday">四</span>
            <span :class="$style.weekday">五</span>
            <span :class="$style.weekday">六</span>
            <span :class="$style.weekday">日</span>
            <span :class="$style.weekday">周报</span>
          </div>

          <div v-for="(week, weekIndex) in calendarWeeks" :key="weekIndex" :class="$style.weekRow">
            <div
              v-for="(day, dayIndex) in week.days"
              :key="dayIndex"
              :class="[$style.dayCell, !day?.isCurrentMonth && $style.otherMonth]"
            >
              <button
                v-if="day"
                type="button"
                :class="[$style.dayButton, day.isToday && $style.today]"
                @click="handleDayClick(day)"
              >
                <span :class="$style.dayNum">{{ day.day }}</span>
                <span v-if="day.recordCount > 0" :class="$style.dot" :style="{ backgroundColor: day.dotColor }"></span>
              </button>
            </div>
            <div :class="$style.actionCell">
              <div :class="$style.actionButtons">
                <button
                  type="button"
                  :class="$style.viewButton"
                  @click="handleViewReport(week.weekStart, week.weekEnd)"
                >
                  <span>查看</span>
                </button>
                <button
                  type="button"
                  :class="[$style.generateButton, isGenerating && $style.buttonLoading]"
                  @click="handleGenerateReport(week.weekStart, week.weekEnd)"
                  :disabled="isGenerating"
                >
                  <span v-if="isGenerating" :class="$style.spinner"></span>
                  <span>{{ isGenerating ? '生成中' : '生成' }}</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div :class="$style.reportPanel">
        <div :class="$style.reportHeader">
          <h3 :class="$style.reportTitle">
            {{ selectedWeek ? `${selectedWeek.weekStart} ~ ${selectedWeek.weekEnd}` : '周报详情' }}
          </h3>
        </div>

        <div v-if="selectedWeeklyReport" :class="$style.reportContent">
          <div :class="$style.reportMeta">
            <span :class="$style.reportDate">{{ selectedWeeklyReport.week_start }} ~ {{ selectedWeeklyReport.week_end }}</span>
            <span :class="$style.reportTime">{{ new Date(selectedWeeklyReport.created_at).toLocaleString('zh-CN') }}</span>
          </div>
          <div v-if="selectedWeeklyReport.title" :class="$style.reportTitleText">{{ selectedWeeklyReport.title }}</div>
          <div v-if="selectedWeeklyReport.summary" :class="$style.reportSummary">{{ selectedWeeklyReport.summary }}</div>
          <pre :class="$style.reportText">{{ selectedWeeklyReport.content }}</pre>
        </div>

        <div v-else-if="selectedWeek" :class="$style.emptyState">
          <div :class="$style.emptyIcon">📝</div>
          <p :class="$style.emptyTitle">暂无周报</p>
          <p :class="$style.emptyDesc">点击左侧"生成"按钮，基于本周工作记录生成周报</p>
        </div>

        <div v-else :class="$style.emptyState">
          <div :class="$style.emptyIcon">📋</div>
          <p :class="$style.emptyTitle">选择一周</p>
          <p :class="$style.emptyDesc">从左侧日历中选择一周查看或生成周报</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style module>
.container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--theme-bg);
  border-radius: 20px;
  overflow: hidden;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 28px;
  background: linear-gradient(135deg, var(--theme-accent) 0%, var(--theme-accent-strong) 100%);
  box-shadow: 0 4px 12px rgba(245, 205, 106, 0.3);
}

.headerLeft {
  display: flex;
  align-items: center;
}

.headerRight {
  display: flex;
  align-items: center;
}

.monthNav {
  display: flex;
  align-items: center;
  gap: 16px;
}

.navButton {
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.3);
  color: var(--theme-accent-text);
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.navButton:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.5);
  transform: scale(1.05);
}

.navButtonDisabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.chevron {
  width: 18px;
  height: 18px;
}

.monthTitle {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--theme-accent-text);
}

.statsBar {
  display: flex;
  gap: 16px;
}

.statCard {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.4);
  border-radius: 12px;
  backdrop-filter: blur(8px);
}

.statValue {
  font-size: 1.5rem;
  font-weight: 700;
  color: #463500;
}

.statLabel {
  font-size: 0.75rem;
  color: #5c4a00;
}

.body {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 400px;
  overflow: hidden;
}

.calendarPanel {
  padding: 20px;
  overflow-y: auto;
}

.calendarWrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.weekdayHeader {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
  padding: 12px 16px;
  background: var(--theme-surface-soft);
  border-radius: 12px;
}

.weekday {
  text-align: center;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--theme-text-secondary);
}

.weekRow {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 8px;
}

.dayCell {
  aspect-ratio: 1;
}

.dayButton {
  width: 100%;
  height: 100%;
  border: 2px solid var(--theme-border);
  border-radius: 10px;
  background: var(--theme-surface);
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  transition: all 0.2s ease;
  position: relative;
}

.dayButton:hover {
  border-color: var(--theme-accent);
  background: var(--theme-surface-soft);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(245, 205, 106, 0.25);
}

.dayNum {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--theme-text-primary);
}

.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.today {
  background: linear-gradient(135deg, var(--theme-accent) 0%, var(--theme-accent-strong) 100%);
  border-color: var(--theme-accent-strong);
}

.today .dayNum {
  color: var(--theme-accent-text);
}

.otherMonth .dayButton {
  opacity: 0.4;
  cursor: default;
}

.otherMonth .dayButton:hover {
  transform: none;
  box-shadow: none;
}

.actionCell {
  display: flex;
  align-items: center;
  justify-content: center;
}

.actionButtons {
  display: flex;
  flex-direction: column;
  gap: 6px;
  width: 100%;
}

.viewButton {
  padding: 6px 8px;
  border: 1px solid var(--theme-border);
  border-radius: 8px;
  background: var(--theme-surface);
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--theme-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.viewButton:hover {
  border-color: var(--theme-accent);
  color: var(--theme-accent-text);
  background: var(--theme-surface-soft);
}

.generateButton {
  padding: 6px 8px;
  border: none;
  border-radius: 8px;
  background: linear-gradient(135deg, var(--theme-accent) 0%, var(--theme-accent-strong) 100%);
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--theme-accent-text);
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.generateButton:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 3px 8px rgba(245, 205, 106, 0.4);
}

.generateButton:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.spinner {
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.reportPanel {
  background: var(--theme-surface);
  border-left: 1px solid var(--theme-border);
  display: flex;
  flex-direction: column;
}

.reportHeader {
  padding: 20px 24px;
  border-bottom: 1px solid var(--theme-border);
  background: var(--theme-surface-soft);
}

.reportTitle {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--theme-text-primary);
}

.reportContent {
  flex: 1;
  padding: 20px 24px;
  overflow-y: auto;
}

.reportMeta {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--theme-border);
}

.reportDate {
  padding: 4px 10px;
  background: var(--theme-surface-soft);
  border-radius: 16px;
  font-size: 0.8125rem;
  color: var(--theme-text-primary);
}

.reportTime {
  padding: 4px 10px;
  background: rgba(245, 205, 106, 0.2);
  border-radius: 16px;
  font-size: 0.8125rem;
  color: var(--theme-accent-text);
}

.reportTitleText {
  margin: 0 0 12px 0;
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--theme-text-primary);
  padding-bottom: 8px;
  border-bottom: 2px solid var(--theme-border);
}

.reportSummary {
  margin: 0 0 16px 0;
  font-size: 0.9375rem;
  line-height: 1.6;
  color: var(--theme-text-primary);
  padding: 12px 16px;
  background: var(--theme-surface-soft);
  border-radius: 10px;
  border-left: 3px solid var(--theme-accent);
}

.reportText {
  margin: 0;
  font-family: inherit;
  font-size: 0.9375rem;
  line-height: 1.7;
  white-space: pre-wrap;
  color: var(--theme-text-primary);
}

.emptyState {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.emptyIcon {
  font-size: 4rem;
  margin-bottom: 16px;
}

.emptyTitle {
  margin: 0 0 8px 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--theme-text-primary);
}

.emptyDesc {
  margin: 0;
  font-size: 0.9375rem;
  color: var(--theme-text-secondary);
  max-width: 280px;
}

@media (max-width: 1024px) {
  .body {
    grid-template-columns: 1fr;
    grid-template-rows: 1fr 1fr;
  }

  .reportPanel {
    border-left: none;
    border-top: 1px solid var(--theme-border);
  }

  .statsBar {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .header {
    padding: 16px 20px;
  }

  .monthTitle {
    font-size: 1.25rem;
  }

  .calendarPanel {
    padding: 16px;
  }

  .weekRow {
    gap: 4px;
  }

  .weekdayHeader {
    gap: 4px;
    padding: 10px 12px;
  }

  .dayNum {
    font-size: 0.8125rem;
  }
}
</style>
