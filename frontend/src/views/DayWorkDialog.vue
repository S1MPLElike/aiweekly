<script setup lang="ts">
import TodayWork from './TodayWork.vue'

const props = defineProps<{
  visible: boolean
  date: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

function handleClose() {
  emit('close')
}
</script>

<template>
  <Teleport to="body">
    <Transition name="dialog">
      <div v-if="visible" :class="[$style.overlay, 'dialog-container']" @click.self="handleClose">
        <div :class="$style.dialog">
          <div :class="$style.closeBar">
            <button :class="$style.closeBtn" @click="handleClose" aria-label="关闭">
              ✕
            </button>
          </div>
          <div :class="$style.dialogContent">
            <TodayWork :date="props.date" />
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
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: stretch;
  justify-content: flex-end;
  z-index: 1000;
}

.dialog {
  width: 80%;
  height: 100vh;
  background: var(--theme-bg);
  border-radius: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: -8px 0 32px rgba(0, 0, 0, 0.15);
  border-left: 1px solid var(--theme-border);
}

.dialogContent {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.closeBar {
  position: absolute;
  top: 16px;
  right: 24px;
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  padding: 0;
  background: transparent;
  border-bottom: none;
}

.closeBtn {
  width: 44px;
  height: 44px;
  border: none;
  background: #ffffff;
  border-radius: 12px;
  cursor: pointer;
  font-size: 1.4rem;
  color: #766e60;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.closeBtn:hover {
  background: #fff8df;
  color: #2f2a22;
  transform: scale(1.1);
}

.closeBtn:focus-visible {
  outline: 2px solid #f5cd6a;
  outline-offset: 2px;
}

.dialog-enter-active,
.dialog-leave-active {
  transition: opacity 0.3s ease;
}

.dialog-enter-active .dialog,
.dialog-leave-active .dialog {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-enter-from .dialog,
.dialog-leave-to .dialog {
  transform: translateX(100%);
}
</style>