<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { getWorkRecords, createWorkRecord, generateDailyReport, type WorkRecord, type DailyReport, type DayWorkData } from '../api/work'
import WorkRecordDetail from './WorkRecordDetail.vue'

const props = defineProps<{
  date?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

function getLocalDateString(): string {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const currentDate = ref(props.date || getLocalDateString())
const isLoading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const records = ref<WorkRecord[]>([])
const dailyReport = ref<DailyReport | null>(null)
const dayStats = ref<{
  total_tasks: number
  total_hours: number
  earliest_start: number
  latest_end: number
  morning_hours: number
  afternoon_hours: number
  evening_hours: number
}>({
  total_tasks: 0,
  total_hours: 0,
  earliest_start: 0,
  latest_end: 0,
  morning_hours: 0,
  afternoon_hours: 0,
  evening_hours: 0,
})

const form = ref({
  title: '',
  content: '',
  startHour: 9,
  endHour: 11
})

const isSubmitting = ref(false)
const detailVisible = ref(false)
const selectedRecord = ref<WorkRecord | null>(null)

const hours = Array.from({ length: 24 }, (_, i) => i)

function sortRecords(records: WorkRecord[]): WorkRecord[] {
  return [...records].sort((a, b) => {
    if (a.start_hour !== b.start_hour) {
      return a.start_hour - b.start_hour
    }
    const timeA = new Date(a.created_at).getTime()
    const timeB = new Date(b.created_at).getTime()
    return timeA - timeB
  })
}

function calculateStats(recordsList: WorkRecord[]): typeof dayStats.value {
  const total_tasks = recordsList.length
  let total_hours = 0
  let earliest_start = 0
  let latest_end = 0
  let morning_hours = 0
  let afternoon_hours = 0
  let evening_hours = 0

  if (total_tasks > 0) {
    earliest_start = recordsList[0].start_hour
    latest_end = recordsList[0].end_hour
    for (const r of recordsList) {
      total_hours += r.end_hour - r.start_hour
      if (r.start_hour < earliest_start) {
        earliest_start = r.start_hour
      }
      if (r.end_hour > latest_end) {
        latest_end = r.end_hour
      }
      const start = r.start_hour
      const end = r.end_hour
      if (start < 12 && end > 6) {
        const mStart = Math.max(start, 6)
        const mEnd = Math.min(end, 12)
        if (mEnd > mStart) morning_hours += mEnd - mStart
      }
      if (start < 18 && end > 12) {
        const aStart = Math.max(start, 12)
        const aEnd = Math.min(end, 18)
        if (aEnd > aStart) afternoon_hours += aEnd - aStart
      }
      if (start < 24 && end > 18) {
        const eStart = Math.max(start, 18)
        const eEnd = Math.min(end, 24)
        if (eEnd > eStart) evening_hours += eEnd - eStart
      }
    }
  }

  return {
    total_tasks,
    total_hours,
    earliest_start,
    latest_end,
    morning_hours,
    afternoon_hours,
    evening_hours,
  }
}

async function fetchData() {
  isLoading.value = true
  errorMessage.value = ''
  try {
    const data = await getWorkRecords(currentDate.value)
    records.value = sortRecords(data.records || [])
    dailyReport.value = data.dailyReport
    dayStats.value = calculateStats(records.value)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '获取数据失败'
  } finally {
    isLoading.value = false
  }
}

async function handleSubmit() {
  if (!form.value.title.trim()) {
    errorMessage.value = '请填写工作标题'
    return
  }

  if (!form.value.content.trim()) {
    errorMessage.value = '请填写工作内容'
    return
  }

  if (form.value.endHour < form.value.startHour) {
    errorMessage.value = '结束时间必须大于或等于开始时间'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await createWorkRecord({
      title: form.value.title,
      content: form.value.content,
      recordDate: currentDate.value,
      startHour: form.value.startHour,
      endHour: form.value.endHour
    })
    successMessage.value = '添加成功'
    form.value.title = ''
    form.value.content = ''
    form.value.startHour = 9
    form.value.endHour = 11
    await fetchData()
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '添加失败'
  } finally {
    isSubmitting.value = false
  }
}

function handleDetail(record: WorkRecord) {
  selectedRecord.value = record
  detailVisible.value = true
}

function handleCloseDetail() {
  detailVisible.value = false
  selectedRecord.value = null
}

function handleRecordDeleted() {
  fetchData()
}

const isGeneratingReport = ref(false)

async function handleGenerateReport() {
  if (records.value.length === 0) {
    errorMessage.value = '请先添加工作记录'
    return
  }
  
  isGeneratingReport.value = true
  errorMessage.value = ''
  successMessage.value = ''
  
  try {
    const report = await generateDailyReport(currentDate.value)
    dailyReport.value = report
    successMessage.value = '日报生成成功'
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '生成日报失败'
  } finally {
    isGeneratingReport.value = false
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${year}年${month}月${day}日`
}

function formatHour(hour: number): string {
  return `${hour}:00`
}

onMounted(() => {
  fetchData()
})

// 监听日期变化
watch(() => props.date, (newDate) => {
  if (newDate) {
    currentDate.value = newDate
    fetchData()
  }
})
</script>

<template>
  <div :class="$style.container">
    <header :class="$style.header">
      <div :class="$style.headerLeft">
        <h2 :class="$style.title">{{ formatDate(currentDate) }}</h2>
        <p :class="$style.subtitle">工作回顾</p>
      </div>
      <div :class="$style.headerStats">
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ dayStats.total_tasks }}</span>
          <span :class="$style.statLabel">工作任务</span>
        </div>
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ dayStats.total_hours }}</span>
          <span :class="$style.statLabel">工作小时</span>
        </div>
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ formatHour(dayStats.earliest_start) }}</span>
          <span :class="$style.statLabel">最早开始</span>
        </div>
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ formatHour(dayStats.latest_end) }}</span>
          <span :class="$style.statLabel">最晚结束</span>
        </div>
      </div>
      <div :class="$style.messages">
        <span v-if="successMessage" :class="$style.success">{{ successMessage }}</span>
      </div>
    </header>

    <section :class="$style.content">
      <section :class="$style.formSection">
        <header :class="$style.formHeader">
          <h3 :class="$style.sectionTitle">添加工作记录</h3>
          <span v-if="errorMessage" :class="$style.formError">{{ errorMessage }}</span>
        </header>
        
        <form :class="$style.form" @submit.prevent="handleSubmit">
          <div :class="$style.formGrid">
            <div :class="$style.leftSection">
              <div :class="$style.field">
                <label :class="$style.label">工作标题</label>
                <input
                  v-model="form.title"
                  type="text"
                  :class="$style.input"
                  placeholder="请输入工作标题"
                />
              </div>
              
              <div :class="$style.field">
                <label :class="$style.label">工作内容</label>
                <textarea
                  v-model="form.content"
                  :class="$style.textarea"
                  placeholder="请输入工作内容"
                  rows="3"
                />
              </div>
            </div>

            <div :class="$style.rightSection">
              <div :class="$style.timeFields">
                <div :class="$style.field">
                  <label :class="$style.label">开始时间</label>
                  <select v-model="form.startHour" :class="$style.select">
                    <option v-for="h in hours" :key="h" :value="h">
                      {{ formatHour(h) }}
                    </option>
                  </select>
                </div>

                <div :class="$style.field">
                  <label :class="$style.label">结束时间</label>
                  <select v-model="form.endHour" :class="$style.select">
                    <option v-for="h in hours" :key="h" :value="h">
                      {{ formatHour(h) }}
                    </option>
                  </select>
                </div>
              </div>

              <div :class="$style.buttonContainer">
                <button
                  type="submit"
                  :class="$style.submitBtn"
                  :disabled="isSubmitting"
                >
                  {{ isSubmitting ? '添加中...' : '添加记录' }}
                </button>
              </div>
            </div>
          </div>
        </form>
      </section>

      <div :class="$style.mainArea">
        <section :class="$style.recordsSection">
          <header :class="$style.sectionHeader">
            <h3 :class="$style.sectionTitle">📝 工作记录</h3>
          </header>
          
          <div v-if="isLoading" :class="$style.loading">加载中...</div>
          <div v-else-if="records.length === 0" :class="$style.empty">
            <div :class="$style.emptyIcon">📋</div>
            <p :class="$style.emptyText">暂无工作记录</p>
            <p :class="$style.emptyHint">在上方添加您的第一条工作记录吧</p>
          </div>
          <div v-else :class="$style.recordsList">
            <div
              v-for="record in records"
              :key="record.id"
              :class="$style.recordCard"
            >
              <div :class="$style.recordTop">
                <h4 :class="$style.recordTitle">{{ record.title }}</h4>
                <p :class="$style.recordContent">{{ record.content }}</p>
              </div>
              <div :class="$style.recordBottom">
                <span :class="$style.recordTime">
                  <span :class="$style.timeIcon">⏰</span>
                  {{ formatHour(record.start_hour) }} - {{ formatHour(record.end_hour) }}
                </span>
                <button
                  :class="$style.detailBtn"
                  @click="handleDetail(record)"
                  aria-label="查看详情"
                >
                  <span :class="$style.detailIcon">🔍</span>
                  详细
                </button>
              </div>
            </div>
            <div v-if="records.length > 3" :class="$style.scrollIndicator">
              👈 滑动查看更多
            </div>
          </div>
        </section>

        <aside :class="$style.reportSection">
          <div :class="$style.reportHeader">
            <h3 :class="$style.reportTitle">📄 当日日报</h3>
          </div>
          <div v-if="dailyReport" :class="$style.reportContent">
            <h4 v-if="dailyReport.title" :class="$style.reportTitleText">{{ dailyReport.title }}</h4>
            <div :class="$style.reportBody">{{ dailyReport.content }}</div>
          </div>
          <div v-else :class="$style.noReport">
            <div :class="$style.noReportIcon">📝</div>
            <p :class="$style.noReportText">暂无日报</p>
            <p :class="$style.noReportHint">日报将在当天工作结束后自动生成</p>
          </div>
          <button :class="$style.generateBtn" @click="handleGenerateReport" :disabled="isGeneratingReport">
            <span :class="$style.generateIcon">✨</span>
            {{ isGeneratingReport ? '生成中...' : '生成日报' }}
          </button>
        </aside>
      </div>
    </section>

    <WorkRecordDetail
      :visible="detailVisible"
      :record="selectedRecord"
      @close="handleCloseDetail"
      @deleted="handleRecordDeleted"
    />
  </div>
</template>

<style module>
.container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #d4f4e8;
  overflow: hidden;
}

.header {
  padding: 12px 24px;
  background: linear-gradient(90deg, #4dd9a0 0%, #3bbf8c 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.headerLeft {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.title {
  margin: 0;
  font-size: 1.2rem;
  font-weight: 700;
  color: #2c5a47;
}

.subtitle {
  margin: 0;
  font-size: 0.8rem;
  color: #3a7a5f;
}

.headerStats {
  display: flex;
  gap: 8px;
}

.statCard {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 6px;
  backdrop-filter: blur(8px);
}

.statValue {
  font-size: 1.2rem;
  font-weight: 700;
  color: #2c5a47;
}

.statLabel {
  font-size: 0.7rem;
  color: #3a7a5f;
}

.messages {
  flex: 1;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.error {
  color: #c53030;
  font-size: 0.9rem;
  font-weight: 500;
}

.success {
  color: #2e7d32;
  font-size: 0.9rem;
  font-weight: 500;
}

.formError {
  color: #c53030;
  font-size: 0.85rem;
  font-weight: 500;
}

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px 24px;
  gap: 16px;
  overflow-y: auto;
  overflow-x: hidden;
  position: relative;
  min-height: 0;
}

.formSection {
  background: #d4f4e8;
  border-radius: 16px;
  padding: 20px 24px;
  border: 2px solid #a8e4c9;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.formHeader {
  margin-bottom: 16px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
}

.form {
  display: flex;
  flex-direction: column;
}

.formGrid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
  align-items: stretch;
}

.leftSection {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.rightSection {
  display: flex;
  flex-direction: column;
  gap: 16px;
  justify-content: flex-start;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.label {
  font-size: 0.88rem;
  color: #2c5a47;
  font-weight: 600;
}

.input,
.textarea,
.select {
  border: 2px solid #a8e4c9;
  border-radius: 10px;
  padding: 12px 16px;
  font-size: 0.95rem;
  background: #ffffff;
  transition: all 0.2s ease;
  color: #2c5a47;
}

.input:focus,
.textarea:focus,
.select:focus {
  outline: none;
  border-color: #4dd9a0;
  box-shadow: 0 0 0 3px rgba(77, 217, 160, 0.25);
}

.textarea {
  resize: vertical;
  min-height: 80px;
}

.timeFields {
  display: flex;
  gap: 12px;
}

.timeFields .field {
  flex: 1;
}

.select {
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='14' height='14' viewBox='0 0 14 14'%3E%3Cpath fill='%232c5a47' d='M7 10L2 4h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 14px center;
  padding-right: 38px;
}

.buttonContainer {
  display: flex;
  justify-content: stretch;
  align-items: flex-end;
  flex: 1;
}

.submitBtn {
  width: 100%;
  border: none;
  border-radius: 12px;
  padding: 16px 24px;
  background: linear-gradient(135deg, #4dd9a0 0%, #3bbf8c 100%);
  color: #1a3a2c;
  font-weight: 700;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(77, 217, 160, 0.3);
}

.submitBtn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(77, 217, 160, 0.4);
}

.submitBtn:focus-visible {
  outline: 3px solid #4dd9a0;
  outline-offset: 2px;
}

.submitBtn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.mainArea {
  display: flex;
  gap: 24px;
  min-height: 0;
  overflow: hidden;
}

.recordsSection {
  flex: 1;
  background: #d4f4e8;
  border-radius: 16px;
  padding: 20px;
  overflow: hidden;
  border: 2px solid #a8e4c9;
  display: flex;
  flex-direction: column;
}

.recordsSection::-webkit-scrollbar {
  display: none;
}

.sectionHeader {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid #a8e4c9;
}

.sectionTitle {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #2c5a47;
}

.loading {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
  color: #3a7a5f;
  font-size: 0.95rem;
}

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  text-align: center;
}

.emptyIcon {
  font-size: 3rem;
  margin-bottom: 12px;
}

.emptyText {
  margin: 0 0 8px;
  color: #2c5a47;
  font-size: 1rem;
  font-weight: 500;
}

.emptyHint {
  margin: 0;
  color: #3a7a5f;
  font-size: 0.9rem;
}

.recordsList {
  display: flex;
  gap: 16px;
  padding: 8px 0;
  scroll-snap-type: x mandatory;
  overflow-x: auto;
  overflow-y: hidden;
  flex: 1;
  min-height: 0;
  width: 100%;
}

.recordCard {
  width: 200px;
  height: 180px;
  background: #ffffff;
  border: 2px solid #a8e4c9;
  border-radius: 12px;
  padding: 12px 14px;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.04);
  scroll-snap-align: start;
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex-shrink: 0;
}

.recordCard:hover {
  border-color: #4dd9a0;
  box-shadow: 0 4px 12px rgba(77, 217, 160, 0.2);
}

.recordTop {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  overflow: hidden;
}

.recordBottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  flex-shrink: 0;
  border-top: 1px solid #a8e4c9;
  gap: 10px;
}

.detailBtn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  padding: 5px 10px;
  border: 1px solid #4dd9a0;
  background: rgba(77, 217, 160, 0.1);
  color: #2c5a47;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 70px;
}

.detailBtn:hover {
  background: #4dd9a0;
  color: #1a3a2c;
}

.detailBtn:focus-visible {
  outline: 2px solid #4dd9a0;
  outline-offset: 2px;
}

.detailIcon {
  font-size: 0.9rem;
}

.recordTitle {
  margin: 0;
  font-weight: 700;
  color: #2c5a47;
  font-size: 1rem;
}

.recordContent {
  margin: 0;
  color: #3a7a5f;
  font-size: 0.92rem;
  line-height: 1.6;
}

.recordTime {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  font-size: 0.75rem;
  color: #2c5a47;
  background: rgba(77, 217, 160, 0.1);
  padding: 5px 10px;
  border-radius: 20px;
  font-weight: 600;
  border: 1px solid #4dd9a0;
  flex: 1;
  min-width: 70px;
}

.timeIcon {
  font-size: 0.75rem;
}

.scrollIndicator {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px 24px;
  background: rgba(77, 217, 160, 0.15);
  border: 1px dashed #4dd9a0;
  border-radius: 12px;
  color: #2c5a47;
  font-size: 0.85rem;
  animation: pulse 2s ease-in-out infinite;
}

.reportSection {
  width: 320px;
  background: #d4f4e8;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  border: 2px solid #a8e4c9;
  overflow: hidden;
}

.reportHeader {
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid #a8e4c9;
}

.reportTitle {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 700;
  color: #2c5a47;
}

.reportContent {
  flex: 1;
  font-size: 0.92rem;
  color: #3a7a5f;
  line-height: 1.7;
  white-space: pre-wrap;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reportTitleText {
  margin: 0;
  font-size: 1rem;
  font-weight: 700;
  color: #2c5a47;
  border-bottom: 2px solid #a8e4c9;
  padding-bottom: 8px;
}

.reportBody {
  flex: 1;
}

.noReport {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  gap: 8px;
}

.noReportIcon {
  font-size: 2.5rem;
}

.noReportText {
  margin: 0;
  color: #3a7a5f;
  font-size: 0.95rem;
}

.noReportHint {
  margin: 0;
  color: #5c9a7f;
  font-size: 0.85rem;
}

.generateBtn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px 20px;
  border: none;
  background: linear-gradient(135deg, #4dd9a0 0%, #3bbf8c 100%);
  color: #1a3a2c;
  border-radius: 12px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 16px;
  box-shadow: 0 2px 8px rgba(77, 217, 160, 0.3);
}

.generateBtn:hover {
  background: linear-gradient(135deg, #3bbf8c 0%, #2d9a70 100%);
  box-shadow: 0 4px 12px rgba(77, 217, 160, 0.4);
  transform: translateY(-1px);
}

.generateBtn:focus-visible {
  outline: 2px solid #4dd9a0;
  outline-offset: 2px;
}

.generateIcon {
  font-size: 1rem;
}

@media (max-width: 1024px) {
  .mainArea {
    flex-direction: column;
  }

  .reportSection {
    width: 100%;
    min-height: 300px;
  }

  .formGrid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .header {
    padding: 16px 20px;
    flex-direction: column;
    gap: 12px;
  }

  .headerStats {
    width: 100%;
    justify-content: space-between;
  }

  .content {
    padding: 16px 20px;
    gap: 16px;
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
  }
}
</style>