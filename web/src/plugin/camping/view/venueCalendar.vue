<template>
  <div class="p-4">
    <el-card>
      <template #header>
        <div class="flex flex-wrap items-center gap-4">
          <span class="font-medium">场地日历</span>
          <el-select v-model="currentVenueId" placeholder="请选择场地" clearable style="width: 200px" @change="onVenueChange">
            <el-option v-for="v in venueOptions" :key="v.ID" :label="v.name" :value="v.ID" />
          </el-select>
          <el-date-picker
            v-model="currentMonth"
            type="month"
            placeholder="选择月份"
            value-format="YYYY-MM"
            style="width: 140px"
            @change="loadCalendar"
          />
          <el-button type="primary" :disabled="!currentVenueId" @click="loadCalendar">刷新</el-button>
        </div>
      </template>
      <div v-if="!currentVenueId" class="text-gray-500 py-8 text-center">请先选择场地</div>
      <div v-else class="venue-calendar">
        <div class="calendar-header flex border-b pb-2 mb-2">
          <div v-for="w in weekLabels" :key="w" class="cell text-center text-gray-500 text-sm font-medium">{{ w }}</div>
        </div>
        <div class="calendar-body">
          <div
            v-for="(cell, idx) in monthCells"
            :key="idx"
            class="cell day-cell"
            :class="{
              'other-month': !cell.isCurrentMonth,
              'open': cell.isCurrentMonth && cell.status === 1,
              'closed': cell.isCurrentMonth && cell.status === 0,
              'no-data': cell.isCurrentMonth && cell.status === undefined
            }"
            @click="cell.isCurrentMonth ? toggleDay(cell) : null"
          >
            <span class="day-num">{{ cell.day }}</span>
            <span v-if="cell.isCurrentMonth" class="day-status">
              {{ cell.status === 1 ? '可约' : cell.status === 0 ? '关闭' : '可约' }}
            </span>
          </div>
        </div>
        <div class="mt-4 flex gap-6 text-sm text-gray-500">
          <span><span class="inline-block w-4 h-4 rounded bg-green-100 border border-green-400 mr-1" />可约</span>
          <span><span class="inline-block w-4 h-4 rounded bg-gray-100 border border-gray-400 mr-1" />关闭</span>
          <span>（无记录默认可约，点击日期切换）</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { getSiteList, getVenueCalendarByVenue, setVenueCalendar } from '@/plugin/camping/api/site'
import { ElMessage } from 'element-plus'
import { ref, computed, onMounted } from 'vue'

defineOptions({ name: 'CampingVenueCalendar' })

const weekLabels = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
const venueOptions = ref([])
const currentVenueId = ref(null)
const currentMonth = ref('')
const calendarList = ref([])

const monthCells = computed(() => {
  if (!currentMonth.value) return []
  const [y, m] = currentMonth.value.split('-').map(Number)
  const first = new Date(y, m - 1, 1)
  const last = new Date(y, m, 0)
  const startDay = first.getDay()
  const daysInMonth = last.getDate()
  const map = new Map()
  calendarList.value.forEach((item) => {
    const d = item.date
    let key = ''
    if (typeof d === 'string') key = d.slice(0, 10)
    else if (d && typeof d === 'object' && d.date) key = String(d.date).slice(0, 10)
    if (key) map.set(key, item.status)
  })
  const cells = []
  for (let i = 0; i < startDay; i++) {
    const prev = new Date(y, m - 1, -startDay + i + 1)
    cells.push({
      day: prev.getDate(),
      dateStr: formatDateStr(prev),
      isCurrentMonth: false,
      status: undefined
    })
  }
  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = `${y}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    cells.push({
      day: d,
      dateStr,
      isCurrentMonth: true,
      status: map.get(dateStr)
    })
  }
  const rest = 42 - cells.length
  for (let i = 0; i < rest; i++) {
    cells.push({
      day: i + 1,
      dateStr: '',
      isCurrentMonth: false,
      status: undefined
    })
  }
  return cells
})

function formatDateStr(d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function getMonthRange(ym) {
  const [y, m] = ym.split('-').map(Number)
  const start = `${y}-${String(m).padStart(2, '0')}-01`
  const lastDay = new Date(y, m, 0).getDate()
  const end = `${y}-${String(m).padStart(2, '0')}-${String(lastDay).padStart(2, '0')}`
  return { start, end }
}

const loadCalendar = async () => {
  if (!currentVenueId.value || !currentMonth.value) return
  const { start, end } = getMonthRange(currentMonth.value)
  const res = await getVenueCalendarByVenue({
    venueId: currentVenueId.value,
    start,
    end
  })
  if (res.code === 0) {
    calendarList.value = res.data || []
  }
}

const onVenueChange = () => {
  if (currentVenueId.value && !currentMonth.value) {
    const now = new Date()
    currentMonth.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
  }
  loadCalendar()
}

const toggleDay = async (cell) => {
  if (!cell.dateStr) return
  const nextStatus = cell.status === 0 ? 1 : 0
  const res = await setVenueCalendar({
    venueId: currentVenueId.value,
    date: cell.dateStr,
    status: nextStatus
  })
  if (res.code === 0) {
    const item = calendarList.value.find((x) => {
      const d = x.date
      const str = typeof d === 'string' ? d.slice(0, 10) : (d && (d.date ? String(d.date) : d))?.slice?.(0, 10)
      return str === cell.dateStr
    })
    if (item) item.status = nextStatus
    else calendarList.value.push({ date: cell.dateStr, status: nextStatus })
    ElMessage.success(nextStatus === 1 ? '已设为可预约' : '已设为关闭')
  } else {
    ElMessage.warning(res.msg || '设置失败')
  }
}

onMounted(async () => {
  const res = await getSiteList({ page: 1, pageSize: 500 })
  if (res.code === 0) venueOptions.value = res.data?.list || []
  const now = new Date()
  currentMonth.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
})
</script>

<style scoped>
.venue-calendar .calendar-header,
.venue-calendar .calendar-body {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
}
.venue-calendar .cell {
  min-height: 32px;
  padding: 4px;
}
.venue-calendar .day-cell {
  min-height: 64px;
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
  margin: 2px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.venue-calendar .day-cell.other-month {
  background: var(--el-fill-color-light);
  color: var(--el-text-color-placeholder);
  cursor: default;
}
.venue-calendar .day-cell.open {
  background: #f0fdf4;
  border-color: #86efac;
}
.venue-calendar .day-cell.closed {
  background: #f5f5f5;
  border-color: #d4d4d4;
  color: #737373;
}
.venue-calendar .day-cell.no-data {
  background: #f0fdf4;
  border-color: #bbf7d0;
}
.venue-calendar .day-num {
  font-weight: 600;
}
.venue-calendar .day-status {
  font-size: 12px;
  margin-top: 2px;
}
</style>
