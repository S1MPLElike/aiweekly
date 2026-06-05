<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { logout, clearAuth } from '../services/auth'
import { getMonthlySummary, type MonthlySummary } from '../api/work'

interface UserInfo {
  username: string
  phone: string
  email?: string
  joinDate: string
}

interface Career {
  id: number
  company: string
  startDate: string
  endDate: string | null
}

interface CareerForm {
  id: number | null
  company: string
  startDate: string
  isCurrent: boolean
  endDate: string
}

const userInfo = ref<UserInfo>({
  username: '',
  phone: '',
  joinDate: ''
})

const monthlySummary = ref<MonthlySummary | null>(null)
const totalWeeks = ref(0)
const totalReports = ref(0)

const showCareerModal = ref(false)
const showDeleteModal = ref(false)
const editingCareerId = ref<number | null>(null)
const deletingCareerId = ref<number | null>(null)

const careerForm = ref<CareerForm>({
  id: null,
  company: '',
  startDate: '',
  isCurrent: false,
  endDate: ''
})

const careers = ref<Career[]>([])

function loadCareers() {
  try {
    const saved = localStorage.getItem('user_careers')
    if (saved) {
      careers.value = JSON.parse(saved)
    } else {
      careers.value = [
        {
          id: 1,
          company: '科技有限公司',
          startDate: '2026-04-27',
          endDate: null
        },
        {
          id: 2,
          company: '互联网创新企业',
          startDate: '2024-01-15',
          endDate: '2026-04-20'
        },
        {
          id: 3,
          company: '创业工作室',
          startDate: '2022-03-01',
          endDate: '2023-12-31'
        }
      ]
    }
  } catch (err) {
    console.error('加载职业生涯数据失败:', err)
    careers.value = []
  }
}

function saveCareers() {
  try {
    localStorage.setItem('user_careers', JSON.stringify(careers.value))
  } catch (err) {
    console.error('保存职业生涯数据失败:', err)
  }
}

function calculateWorkDays(startDate: string, endDate: string | null): number {
  try {
    const start = new Date(startDate)
    const end = endDate ? new Date(endDate) : new Date()
    
    const diffTime = Math.abs(end.getTime() - start.getTime())
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    
    return diffDays
  } catch (err) {
    return 0
  }
}

const totalCareerDays = computed(() => {
  if (!careers.value || !Array.isArray(careers.value)) return 0
  return careers.value.reduce((total, career) => {
    if (!career) return total
    return total + calculateWorkDays(career.startDate, career.endDate)
  }, 0)
})

function openAddCareer() {
  careerForm.value = {
    id: null,
    company: '',
    startDate: '',
    isCurrent: false,
    endDate: ''
  }
  editingCareerId.value = null
  showCareerModal.value = true
}

function openEditCareer(career: Career) {
  if (!career) return
  careerForm.value = {
    id: career.id,
    company: career.company || '',
    startDate: career.startDate || '',
    isCurrent: career.endDate === null,
    endDate: career.endDate || ''
  }
  editingCareerId.value = career.id
  showCareerModal.value = true
}

function openDeleteCareer(careerId: number) {
  deletingCareerId.value = careerId
  showDeleteModal.value = true
}

function closeCareerModal() {
  showCareerModal.value = false
  editingCareerId.value = null
}

function closeDeleteModal() {
  showDeleteModal.value = false
  deletingCareerId.value = null
}

function handleSaveCareer() {
  if (!careerForm.value.company || !careerForm.value.startDate) {
    return
  }

  if (careerForm.value.id) {
    const index = careers.value.findIndex(c => c && c.id === careerForm.value.id)
    if (index !== -1) {
      careers.value[index] = {
        id: careerForm.value.id,
        company: careerForm.value.company,
        startDate: careerForm.value.startDate,
        endDate: careerForm.value.isCurrent ? null : careerForm.value.endDate || null
      }
    }
  } else {
    const newId = careers.value.length > 0 
      ? Math.max(...careers.value.map(c => c ? c.id : 0), 0) + 1
      : 1
    careers.value.unshift({
      id: newId,
      company: careerForm.value.company,
      startDate: careerForm.value.startDate,
      endDate: careerForm.value.isCurrent ? null : careerForm.value.endDate || null
    })
  }

  saveCareers()
  closeCareerModal()
}

function handleDeleteCareer() {
  if (deletingCareerId.value !== null) {
    careers.value = careers.value.filter(c => c && c.id !== deletingCareerId.value)
    saveCareers()
  }
  closeDeleteModal()
}

async function loadUserInfo() {
  try {
    userInfo.value = {
      username: '测试用户',
      phone: '00000000001',
      joinDate: '2026-04-27'
    }
    
    const now = new Date()
    const summary = await getMonthlySummary(now.getFullYear(), now.getMonth() + 1)
    monthlySummary.value = summary
    
    totalWeeks.value = 8
    totalReports.value = 6
  } catch (err) {
    console.error('加载用户信息失败:', err)
  }
}

function handleLogout() {
  logout()
  clearAuth()
  window.location.href = '/login'
}

onMounted(() => {
  loadUserInfo()
  loadCareers()
})
</script>

<template>
  <div :class="$style.container">
    <div :class="$style.profileCard">
      <div :class="$style.avatarSection">
        <div :class="$style.avatar">
          <svg :class="$style.avatarIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
          </svg>
        </div>
        <h2 :class="$style.username">{{ userInfo.username }}</h2>
        <p :class="$style.userPhone">{{ userInfo.phone }}</p>
      </div>

      <div :class="$style.divider"></div>

      <div :class="$style.statsGrid">
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ monthlySummary?.total_tasks || 0 }}</span>
          <span :class="$style.statLabel">本月任务</span>
        </div>
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ monthlySummary?.total_hours || 0 }}</span>
          <span :class="$style.statLabel">本月工时</span>
        </div>
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ monthlySummary?.week_count || 0 }}</span>
          <span :class="$style.statLabel">本月周数</span>
        </div>
        <div :class="$style.statCard">
          <span :class="$style.statValue">{{ monthlySummary?.report_count || 0 }}</span>
          <span :class="$style.statLabel">本月周报</span>
        </div>
      </div>
    </div>

    <div :class="$style.mainLayout">
      <div :class="$style.leftPanel">
        <div :class="$style.sectionCard">
          <div :class="$style.sectionHeader">
            <h3 :class="$style.sectionTitle">累计统计</h3>
          </div>
          <div :class="$style.summaryStats">
            <div :class="$style.summaryItem">
              <div :class="$style.summaryIcon">📅</div>
              <div :class="$style.summaryInfo">
                <span :class="$style.summaryValue">{{ totalWeeks }} 周</span>
                <span :class="$style.summaryLabel">累计工作周</span>
              </div>
            </div>
            <div :class="$style.summaryItem">
              <div :class="$style.summaryIcon">📝</div>
              <div :class="$style.summaryInfo">
                <span :class="$style.summaryValue">{{ totalReports }} 份</span>
                <span :class="$style.summaryLabel">已生成周报</span>
              </div>
            </div>
          </div>
        </div>

        <div :class="$style.sectionCard">
          <div :class="$style.sectionHeader">
            <h3 :class="$style.sectionTitle">账户信息</h3>
          </div>
          <div :class="$style.infoList">
            <div :class="$style.infoItem">
              <span :class="$style.infoLabel">用户名</span>
              <span :class="$style.infoValue">{{ userInfo.username }}</span>
            </div>
            <div :class="$style.infoItem">
              <span :class="$style.infoLabel">手机号</span>
              <span :class="$style.infoValue">{{ userInfo.phone }}</span>
            </div>
            <div :class="$style.infoItem">
              <span :class="$style.infoLabel">注册日期</span>
              <span :class="$style.infoValue">{{ userInfo.joinDate }}</span>
            </div>
          </div>
        </div>

        <div :class="$style.sectionCard">
          <button :class="$style.logoutButton" @click="handleLogout">
            <svg :class="$style.logoutIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
              <polyline points="16 17 21 12 16 7" />
              <line x1="21" y1="12" x2="9" y2="12" />
            </svg>
            <span>退出登录</span>
          </button>
        </div>
      </div>

      <div :class="$style.rightPanel">
        <div :class="$style.sectionCard">
          <div :class="$style.sectionHeader">
            <h3 :class="$style.sectionTitle">职业生涯</h3>
            <button :class="$style.addCareerButton" @click="openAddCareer">
              <svg :class="$style.addIcon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19" />
                <line x1="5" y1="12" x2="19" y2="12" />
              </svg>
              新增经历
            </button>
          </div>
          <div :class="$style.careerSummary">
            <div :class="$style.careerSummaryItem">
              <span :class="$style.careerSummaryValue">{{ totalCareerDays }}</span>
              <span :class="$style.careerSummaryLabel">累计工作天数</span>
            </div>
            <div :class="$style.careerSummaryItem">
              <span :class="$style.careerSummaryValue">{{ careers.length }}</span>
              <span :class="$style.careerSummaryLabel">经历公司数量</span>
            </div>
          </div>
          <div v-if="careers && careers.length > 0" :class="$style.careerList">
            <div
              v-for="career in careers"
              :key="career?.id"
              :class="[$style.careerItem, career && career.endDate === null && $style.careerItemCurrent]"
            >
              <div v-if="career" :class="$style.careerInfo">
                <div :class="$style.careerCompany">
                  {{ career.company }}
                  <span v-if="career.endDate === null" :class="$style.currentBadge">当前</span>
                </div>
                <div :class="$style.careerDate">
                  {{ career.startDate }}
                  <span v-if="career.endDate"> - {{ career.endDate }}</span>
                  <span v-else> - 至今</span>
                </div>
              </div>
              <div v-if="career" :class="$style.careerActions">
                <div :class="$style.careerDays">
                  <span :class="$style.careerDaysValue">{{ calculateWorkDays(career.startDate, career.endDate) }}</span>
                  <span :class="$style.careerDaysLabel">天</span>
                </div>
                <div :class="$style.actionButtons">
                  <button :class="$style.editButton" @click="openEditCareer(career)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                      <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                    </svg>
                  </button>
                  <button :class="$style.deleteButton" @click="openDeleteCareer(career.id)">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3 6 5 6 21 6" />
                      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                      <line x1="10" y1="11" x2="10" y2="17" />
                      <line x1="14" y1="11" x2="14" y2="17" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div v-else :class="$style.emptyCareer">
            <p>暂无工作经历</p>
            <button :class="$style.addFirstButton" @click="openAddCareer">添加第一条经历</button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showCareerModal" :class="$style.modalOverlay" @click="closeCareerModal">
      <div :class="$style.modal" @click.stop>
        <div :class="$style.modalHeader">
          <h3 :class="$style.modalTitle">{{ editingCareerId ? '编辑工作经历' : '新增工作经历' }}</h3>
          <button :class="$style.closeButton" @click="closeCareerModal">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
        <div :class="$style.modalBody">
          <div :class="$style.formGroup">
            <label :class="$style.formLabel">公司名称</label>
            <input
              v-model="careerForm.company"
              :class="$style.formInput"
              type="text"
              placeholder="请输入公司名称"
            />
          </div>
          <div :class="$style.formRow">
            <div :class="$style.formGroup">
              <label :class="$style.formLabel">入职日期</label>
              <input
                v-model="careerForm.startDate"
                :class="$style.formInput"
                type="date"
              />
            </div>
          </div>
          <div :class="$style.formGroup">
            <label :class="$style.formCheckbox">
              <input
                v-model="careerForm.isCurrent"
                type="checkbox"
              />
              <span>当前公司</span>
            </label>
          </div>
          <div v-if="!careerForm.isCurrent" :class="$style.formGroup">
            <label :class="$style.formLabel">离职日期</label>
            <input
              v-model="careerForm.endDate"
              :class="$style.formInput"
              type="date"
            />
          </div>
        </div>
        <div :class="$style.modalFooter">
          <button :class="$style.cancelButton" @click="closeCareerModal">取消</button>
          <button :class="$style.saveButton" @click="handleSaveCareer">保存</button>
        </div>
      </div>
    </div>

    <div v-if="showDeleteModal" :class="$style.modalOverlay" @click="closeDeleteModal">
      <div :class="$style.modal" @click.stop>
        <div :class="$style.modalHeader">
          <h3 :class="$style.modalTitle">确认删除</h3>
          <button :class="$style.closeButton" @click="closeDeleteModal">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </div>
        <div :class="$style.modalBody">
          <p :class="$style.deleteConfirmText">确定要删除这条工作经历吗？此操作无法撤销。</p>
        </div>
        <div :class="$style.modalFooter">
          <button :class="$style.cancelButton" @click="closeDeleteModal">取消</button>
          <button :class="$style.deleteConfirmButton" @click="handleDeleteCareer">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style module>
.container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 32px;
  max-width: 1200px;
  margin: 0 auto;
  background: var(--theme-bg);
  min-height: 100vh;
}

.mainLayout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.leftPanel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rightPanel {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.profileCard {
  background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-accent-strong) 100%);
  border-radius: 20px;
  padding: 32px;
  text-align: center;
  box-shadow: 0 8px 24px rgba(245, 205, 106, 0.3);
}

.avatarSection {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.avatarIcon {
  width: 40px;
  height: 40px;
  color: var(--theme-accent-strong);
}

.username {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--theme-accent-text);
}

.userPhone {
  margin: 0;
  font-size: 0.9375rem;
  color: rgba(70, 53, 0, 0.7);
}

.divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.4);
  margin: 24px 0;
}

.statsGrid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.statCard {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.statValue {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--theme-accent-text);
}

.statLabel {
  font-size: 0.75rem;
  color: rgba(70, 53, 0, 0.8);
}

.sectionCard {
  background: var(--theme-surface);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.sectionHeader {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.sectionTitle {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--theme-text-primary);
}

.addCareerButton {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: none;
  border-radius: 10px;
  background: var(--theme-accent);
  color: var(--theme-accent-text);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.addCareerButton:hover {
  background: var(--theme-accent-strong);
  transform: translateY(-1px);
}

.addIcon {
  width: 18px;
  height: 18px;
}

.careerSummary {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.careerSummaryItem {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 16px;
  background: var(--theme-surface-soft);
  border-radius: 12px;
}

.careerSummaryValue {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--theme-text-primary);
}

.careerSummaryLabel {
  font-size: 0.8125rem;
  color: var(--theme-text-secondary);
}

.careerList {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.careerItem {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--theme-surface-soft);
  border-radius: 12px;
  border: 2px solid var(--theme-border);
  transition: all 0.2s ease;
}

.careerItem:hover {
  border-color: var(--theme-accent);
  transform: translateY(-2px);
}

.careerItemCurrent {
  border-color: var(--theme-accent);
  background: linear-gradient(135deg, var(--theme-surface-soft) 0%, rgba(245, 205, 106, 0.1) 100%);
}

.careerInfo {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.careerCompany {
  font-size: 1rem;
  font-weight: 600;
  color: var(--theme-text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.currentBadge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  background: var(--theme-accent);
  color: var(--theme-accent-text);
}

.careerDate {
  font-size: 0.8125rem;
  color: var(--theme-text-secondary);
}

.careerActions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.careerDays {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.careerDaysValue {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--theme-accent-text);
}

.careerDaysLabel {
  font-size: 0.8125rem;
  color: var(--theme-text-secondary);
}

.actionButtons {
  display: flex;
  gap: 8px;
}

.editButton {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: var(--theme-surface);
  color: var(--theme-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.editButton:hover {
  background: var(--theme-accent);
  color: var(--theme-accent-text);
  transform: scale(1.1);
}

.editButton svg {
  width: 18px;
  height: 18px;
}

.deleteButton {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: var(--theme-surface);
  color: #ff4d4f;
  cursor: pointer;
  transition: all 0.2s ease;
}

.deleteButton:hover {
  background: rgba(255, 77, 79, 0.1);
  transform: scale(1.1);
}

.deleteButton svg {
  width: 18px;
  height: 18px;
}

.emptyCareer {
  text-align: center;
  padding: 40px;
  color: var(--theme-text-secondary);
}

.addFirstButton {
  margin-top: 16px;
  padding: 10px 24px;
  border: none;
  border-radius: 10px;
  background: var(--theme-accent);
  color: var(--theme-accent-text);
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.addFirstButton:hover {
  background: var(--theme-accent-strong);
}

.summaryStats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.summaryItem {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--theme-surface-soft);
  border-radius: 12px;
}

.summaryIcon {
  font-size: 1.75rem;
}

.summaryInfo {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.summaryValue {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--theme-text-primary);
}

.summaryLabel {
  font-size: 0.8125rem;
  color: var(--theme-text-secondary);
}

.infoList {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.infoItem {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid var(--theme-border);
}

.infoItem:last-child {
  border-bottom: none;
}

.infoLabel {
  font-size: 0.9375rem;
  color: var(--theme-text-secondary);
}

.infoValue {
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--theme-text-primary);
}

.logoutButton {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
  padding: 14px;
  border: 2px solid #ff4d4f;
  border-radius: 12px;
  background: transparent;
  font-size: 0.9375rem;
  font-weight: 500;
  color: #ff4d4f;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logoutButton:hover {
  background: rgba(255, 77, 79, 0.1);
}

.logoutIcon {
  width: 18px;
  height: 18px;
}

.modalOverlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal {
  background: var(--theme-surface);
  border-radius: 16px;
  width: 100%;
  max-width: 480px;
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.2);
  animation: modalSlideIn 0.3s ease;
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modalHeader {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--theme-border);
}

.modalTitle {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--theme-text-primary);
}

.closeButton {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: var(--theme-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.closeButton:hover {
  background: var(--theme-surface-soft);
  color: var(--theme-text-primary);
}

.closeButton svg {
  width: 20px;
  height: 20px;
}

.modalBody {
  padding: 24px;
}

.formGroup {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.formRow {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.formLabel {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--theme-text-primary);
}

.formInput {
  padding: 12px 16px;
  border: 2px solid var(--theme-border);
  border-radius: 10px;
  background: var(--theme-surface);
  color: var(--theme-text-primary);
  font-size: 0.9375rem;
  transition: all 0.2s ease;
}

.formInput:focus {
  outline: none;
  border-color: var(--theme-accent);
}

.formCheckbox {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9375rem;
  color: var(--theme-text-primary);
  cursor: pointer;
}

.formCheckbox input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: var(--theme-accent);
}

.deleteConfirmText {
  margin: 0;
  font-size: 0.9375rem;
  color: var(--theme-text-secondary);
  text-align: center;
}

.modalFooter {
  display: flex;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid var(--theme-border);
}

.cancelButton {
  flex: 1;
  padding: 12px 24px;
  border: 2px solid var(--theme-border);
  border-radius: 10px;
  background: transparent;
  color: var(--theme-text-primary);
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancelButton:hover {
  background: var(--theme-surface-soft);
}

.saveButton {
  flex: 1;
  padding: 12px 24px;
  border: none;
  border-radius: 10px;
  background: var(--theme-accent);
  color: var(--theme-accent-text);
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.saveButton:hover {
  background: var(--theme-accent-strong);
  transform: translateY(-1px);
}

.deleteConfirmButton {
  flex: 1;
  padding: 12px 24px;
  border: none;
  border-radius: 10px;
  background: #ff4d4f;
  color: white;
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.deleteConfirmButton:hover {
  background: #ff3333;
  transform: translateY(-1px);
}

@media (max-width: 900px) {
  .mainLayout {
    grid-template-columns: 1fr;
  }

  .container {
    max-width: 500px;
    padding: 20px;
  }

  .formRow {
    grid-template-columns: 1fr;
  }
}
</style>
