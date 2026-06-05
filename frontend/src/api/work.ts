import { fetchWithAuth } from './client'

export interface WorkRecord {
  id: number
  user_id: number
  title: string
  content: string
  record_date: string
  start_hour: number
  end_hour: number
  created_at: string
  updated_at: string
}

export interface DailyReport {
  id: number
  user_id: number
  title: string
  content: string
  type: 'simple' | 'default'
  report_date: string
  created_at: string
}

export interface DayWorkData {
  date: string
  records: WorkRecord[]
  dailyReport: DailyReport | null
  stats: {
    total_tasks: number
    total_hours: number
    earliest_start: number
    latest_end: number
    morning_hours: number
    afternoon_hours: number
    evening_hours: number
  }
}

export async function getWorkRecords(date: string): Promise<DayWorkData> {
  const response = await fetchWithAuth(`/work/records?date=${date}`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('获取工作记录失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '获取工作记录失败')
  }

  return {
    date: result.data.date,
    records: result.data.records || [],
    dailyReport: result.data.daily_report || null,
    stats: result.data.stats || {
      total_tasks: 0,
      total_hours: 0,
      earliest_start: 0,
      latest_end: 0,
      morning_hours: 0,
      afternoon_hours: 0,
      evening_hours: 0,
    }
  }
}

export interface CreateWorkParams {
  title: string
  content: string
  recordDate: string
  startHour: number
  endHour: number
}

export async function createWorkRecord(params: CreateWorkParams): Promise<WorkRecord> {
  const response = await fetchWithAuth('/work/records', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: params.title,
      content: params.content,
      record_date: params.recordDate,
      start_hour: params.startHour,
      end_hour: params.endHour
    }),
  })

  if (!response.ok) {
    throw new Error('创建工作记录失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '创建工作记录失败')
  }

  return result.data
}

export interface UpdateWorkParams {
  title: string
  content: string
  start_hour: number
  end_hour: number
}

export async function updateWorkRecord(id: number, params: UpdateWorkParams): Promise<WorkRecord> {
  const response = await fetchWithAuth(`/work/records/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params),
  })

  if (!response.ok) {
    throw new Error('更新工作记录失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '更新工作记录失败')
  }

  return result.data
}

export async function deleteWorkRecord(id: number): Promise<void> {
  const response = await fetchWithAuth(`/work/records/${id}`, {
    method: 'DELETE',
  })

  if (!response.ok) {
    throw new Error('删除工作记录失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '删除工作记录失败')
  }
}

export interface MonthlyRecordStats {
  date: string
  count: number
}

export async function getMonthlyRecordStats(year: number, month: number): Promise<MonthlyRecordStats[]> {
  const response = await fetchWithAuth(`/work/records/monthly?year=${year}&month=${month}`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('获取月度工作记录统计失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '获取月度工作记录统计失败')
  }

  return result.data || []
}

export interface WeeklyReport {
	id: number
	user_id: number
	title: string
	summary: string
	content: string
	week_start: string
	week_end: string
	created_at: string
}

export interface UserSetting {
	id: number
	user_id: number
	weekly_report_style: 'professional' | 'relaxed' | 'concise'
	created_at: string
	updated_at: string
}

export interface MonthlySummary {
  total_hours: number
  total_tasks: number
  week_count: number
  report_count: number
  busy_days: string[]
}

export async function getWeeklyReport(weekStart: string, weekEnd: string): Promise<WeeklyReport | null> {
  const response = await fetchWithAuth(`/work/weekly-report?week_start=${weekStart}&week_end=${weekEnd}`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('获取周报失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '获取周报失败')
  }

  return result.data
}

export async function generateWeeklyReport(weekStart: string, weekEnd: string): Promise<WeeklyReport> {
  const response = await fetchWithAuth('/work/weekly-report/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      week_start: weekStart,
      week_end: weekEnd,
    }),
  })

  if (!response.ok) {
    throw new Error('生成周报失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '生成周报失败')
  }

  return result.data
}

export async function getMonthlySummary(year: number, month: number): Promise<MonthlySummary> {
  const response = await fetchWithAuth(`/work/monthly-summary?year=${year}&month=${month}`, {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('获取月度统计失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '获取月度统计失败')
  }

  return result.data
}

export async function generateDailyReport(reportDate: string): Promise<DailyReport> {
  const response = await fetchWithAuth('/work/daily-report/generate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ report_date: reportDate }),
  })

  if (!response.ok) {
    throw new Error('生成日报失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '生成日报失败')
  }

  return result.data
}

export async function getUserSetting(): Promise<UserSetting> {
  const response = await fetchWithAuth('/work/settings', {
    method: 'GET',
  })

  if (!response.ok) {
    throw new Error('获取用户设置失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '获取用户设置失败')
  }

  return result.data
}

export async function updateUserSetting(weeklyReportStyle: 'professional' | 'relaxed' | 'concise'): Promise<UserSetting> {
  const response = await fetchWithAuth('/work/settings', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ weekly_report_style: weeklyReportStyle }),
  })

  if (!response.ok) {
    throw new Error('更新用户设置失败')
  }

  const result = await response.json()
  if (result.code !== 0) {
    throw new Error(result.msg || '更新用户设置失败')
  }

  return result.data
}