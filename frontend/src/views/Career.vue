<script setup lang="ts">
import { ref, computed } from 'vue'

interface Career {
  id: number
  company: string
  startDate: string
  endDate: string | null // null表示当前公司
}

const careers = ref<Career[]>([
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
])

function calculateWorkDays(startDate: string, endDate: string | null): number {
  const start = new Date(startDate)
  const end = endDate ? new Date(endDate) : new Date()
  
  const diffTime = Math.abs(end.getTime() - start.getTime())
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  return diffDays
}

const totalCareerDays = computed(() => {
  return careers.value.reduce((total, career) => {
    return total + calculateWorkDays(career.startDate, career.endDate)
  }, 0)
})
</script>

<template>
  <div :class="$style.container">
    <div :class="$style.header">
      <h1 :class="$style.title">职业生涯</h1>
      <p :class="$style.subtitle">记录您的职业发展历程</p>
    </div>

    <div :class="$style.summaryGrid">
      <div :class="$style.summaryCard">
        <div :class="$style.summaryIcon">📊</div>
        <div :class="$style.summaryInfo">
          <span :class="$style.summaryValue">{{ totalCareerDays }}</span>
          <span :class="$style.summaryLabel">累计工作天数</span>
        </div>
      </div>
      <div :class="$style.summaryCard">
        <div :class="$style.summaryIcon">🏢</div>
        <div :class="$style.summaryInfo">
          <span :class="$style.summaryValue">{{ careers.length }}</span>
          <span :class="$style.summaryLabel">经历公司数量</span>
        </div>
      </div>
    </div>

    <div :class="$style.careerTimeline">
      <div :class="$style.timelineTitle">工作经历</div>
      <div :class="$style.timelineLine"></div>
      <div :class="$style.careerList">
        <div
          v-for="(career, index) in careers"
