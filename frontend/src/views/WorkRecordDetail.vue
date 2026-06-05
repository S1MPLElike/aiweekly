<script setup lang="ts">
import { ref } from 'vue'
import type { WorkRecord } from '../api/work'
import { deleteWorkRecord } from '../api/work'

const props = defineProps<{
  visible: boolean
  record: WorkRecord | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'deleted'): void
}>()

const isDeleting = ref(false)
const showConfirm = ref(false)

function handleClose() {
  emit('close')
}

function openConfirm() {
  showConfirm.value = true
}

function closeConfirm() {
  showConfirm.value = false
}

async function confirmDelete() {
  isDeleting.value = true
  try {
    if (props.record) {
      await deleteWorkRecord(props.record.id)
      emit('deleted')
      emit('close')
    }
  } catch (error) {
    console.error('删除失败:', error)
  } finally {
    isDeleting.value = false
    showConfirm.value = false
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const weekDay = weekDays[date.getDay()]
  return `${year}年${month}月${day}日 ${weekDay}`
}

function formatHour(hour: number): string {
  return `${hour}`
}

function formatDateTime(dateStr: string): string {
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible && record" :class="$style.overlay" @click.self="handleClose">
        <div :class="$style.modal">
          <header :class="$style.modalHeader">
            <h3 :class="$style.modalTitle">📋 工作记录详情</h3>
            <button :class="$style.closeBtn" @click="handleClose" aria-label="关闭">
              ✕
            </button>
          </header>
          
          <div :class="$style.modalBody">
            <div :class="$style.detailSection">
              <div :class="$style.rowGroup">
                <div :class="$style.infoCol">
                  <div :class="$style.detailRow">
                    <span :class="$style.detailLabel">记录日期</span>
                    <span :class="$style.detailValue">{{ formatDate(record.record_date) }}</span>
                  </div>
                  
                  <div :class="$style.detailRow">
                    <span :class="$style.detailLabel">工作标题</span>
                    <span :class="[$style.detailValue, $style.detailTitle]">{{ record.title }}</span>
                  </div>
                  
                  <div :class="$style.timeRow">
                    <div :class="$style.timeItem">
                      <span :class="$style.detailLabel">开始时间</span>
                      <span :class="$style.detailValue">{{ formatHour(record.start_hour) }}:00</span>
                    </div>
                    <div :class="$style.timeSeparator">—</div>
                    <div :class="$style.timeItem">
                      <span :class="$style.detailLabel">结束时间</span>
                      <span :class="$style.detailValue">{{ formatHour(record.end_hour) }}:00</span>
                    </div>
                  </div>
                  
                  <div :class="$style.detailRow">
                    <span :class="$style.detailLabel">工作时长</span>
                    <span :class="[$style.detailValue, $style.highlight]">{{ record.end_hour - record.start_hour }} 小时</span>
                  </div>
                </div>
                
                <div :class="$style.contentCol">
                  <div :class="$style.detailRow">
                    <span :class="$style.detailLabel">工作内容</span>
                    <p :class="$style.detailContent">{{ record.content }}</p>
                  </div>
                </div>
              </div>
              
              <div :class="$style.metaSection">
                <div :class="$style.detailRow">
                  <span :class="$style.detailLabel">创建时间</span>
                  <span :class="$style.detailValue">{{ formatDateTime(record.created_at) }}</span>
                </div>
                
                <div v-if="record.updated_at !== record.created_at" :class="$style.detailRow">
                  <span :class="$style.detailLabel">更新时间</span>
                  <span :class="$style.detailValue">{{ formatDateTime(record.updated_at) }}</span>
                </div>
              </div>
            </div>
          </div>
          
          <footer :class="$style.modalFooter">
            <button :class="$style.deleteButton" @click="openConfirm">
              <span :class="$style.deleteIcon">🗑️</span>
              删除
            </button>
            <button :class="$style.closeButton" @click="handleClose">
              关闭
            </button>
          </footer>
        </div>
      </div>
    </Transition>
    
    <Transition name="modal">
      <div v-if="showConfirm" :class="$style.confirmOverlay" @click.self="closeConfirm">
        <div :class="$style.confirmModal">
          <div :class="$style.confirmIcon">⚠️</div>
          <h4 :class="$style.confirmTitle">确认删除</h4>
          <p :class="$style.confirmText">确定要删除这条工作记录吗？此操作无法撤销。</p>
          <div :class="$style.confirmButtons">
            <button :class="$style.cancelButton" @click="closeConfirm">
              取消
            </button>
            <button 
              :class="$style.confirmButton" 
              @click="confirmDelete"
              :disabled="isDeleting"
            >
              {{ isDeleting ? '删除中...' : '确认删除' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style module>
.overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal {
  width: 85%;
  max-width: 800px;
  max-height: 85vh;
  background: var(--theme-surface);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  border: 2px solid var(--theme-border);
}

.modalHeader {
  padding: 20px 24px;
  background: linear-gradient(135deg, var(--theme-accent) 0%, var(--theme-accent-strong) 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modalTitle {
  margin: 0;
  font-size: 1.15rem;
  font-weight: 700;
  color: #463500;
}

.closeBtn {
  width: 36px;
  height: 36px;
  border: none;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 10px;
  cursor: pointer;
  font-size: 1.2rem;
  color: #766e60;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.closeBtn:hover {
  background: #ffffff;
  transform: scale(1.1);
}

.closeBtn:focus-visible {
  outline: 2px solid #463500;
  outline-offset: 2px;
}

.modalBody {
  padding: 24px;
  overflow-y: auto;
  max-height: calc(85vh - 140px);
}

.detailSection {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rowGroup {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 24px;
}

.infoCol {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.contentCol {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detailRow {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.metaSection {
  padding-top: 16px;
  border-top: 1px solid var(--border, #f2e8c8);
  display: flex;
  gap: 24px;
}

.detailLabel {
  font-size: 0.85rem;
  color: var(--theme-text-secondary);
  font-weight: 600;
}

.detailValue {
  font-size: 1rem;
  color: var(--theme-text-primary);
  font-weight: 500;
}

.detailTitle {
  font-size: 1.15rem;
  font-weight: 700;
}

.highlight {
  color: #4a90d9;
  font-weight: 700;
}

.detailContent {
  margin: 0;
  font-size: 1rem;
  color: var(--text-primary, #2f2a22);
  line-height: 1.7;
  background: var(--surface-soft, #fff8df);
  padding: 16px;
  border-radius: 12px;
  border: 1px solid var(--border, #f2e8c8);
  min-height: 120px;
  max-height: 200px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
}

.timeRow {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  gap: 12px;
}

.timeItem {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.timeItem:first-child {
  text-align: left;
}

.timeItem:last-child {
  text-align: right;
}

.timeSeparator {
  font-size: 1.2rem;
  color: var(--accent);
  font-weight: 300;
}

.modalFooter {
  padding: 16px 24px;
  border-top: 1px solid var(--theme-border);
  display: flex;
  justify-content: flex-end;
}

.deleteButton {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border: 1px solid #fc8181;
  background: rgba(252, 129, 129, 0.1);
  color: #c53030;
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.deleteButton:hover {
  background: #fff5f5;
  border-color: #fc8181;
}

.deleteButton:focus-visible {
  outline: 2px solid #fc8181;
  outline-offset: 2px;
}

.deleteIcon {
  font-size: 1rem;
}

.closeButton {
  padding: 10px 28px;
  border: none;
  background: var(--accent);
  color: var(--accent-text);
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(245, 205, 106, 0.3);
}

.closeButton:hover {
  background: var(--accent-hover);
  box-shadow: 0 4px 12px rgba(245, 205, 106, 0.4);
  transform: translateY(-1px);
}

.closeButton:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.confirmOverlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2100;
}

.confirmModal {
  width: 90%;
  max-width: 400px;
  background: var(--theme-surface);
  border-radius: 16px;
  padding: 28px;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.confirmIcon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.confirmTitle {
  margin: 0 0 12px;
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--theme-text-primary);
}

.confirmText {
  margin: 0 0 24px;
  color: var(--theme-text-secondary);
  font-size: 0.95rem;
  line-height: 1.6;
}

.confirmButtons {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.cancelButton {
  padding: 10px 24px;
  border: 1px solid var(--theme-border);
  background: var(--theme-surface);
  color: var(--theme-text-secondary);
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cancelButton:hover {
  background: var(--theme-surface-soft);
}

.cancelButton:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.confirmButton {
  padding: 10px 24px;
  border: none;
  background: #fc8181;
  color: #ffffff;
  border-radius: 10px;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.confirmButton:hover:not(:disabled) {
  background: #f56565;
}

.confirmButton:focus-visible {
  outline: 2px solid #fc8181;
  outline-offset: 2px;
}

.confirmButton:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.25s ease;
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.25s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: scale(0.95) translateY(10px);
  opacity: 0;
}
</style>