<template>
  <div class="h5-verify" data-page="h5-verify">
    <h1 class="app-title">乐享江岛核销系统</h1>
    <!-- 无核销码 -->
    <div v-if="!codeFromUrl" class="tip">
      <p>请扫描二维码进入核销页</p>
    </div>

    <!-- 身份验证 -->
    <div v-else-if="!passedAuth" class="auth-box">
      <h2>工作人员验证</h2>
      <p class="hint">请输入核销密码</p>
      <input
        v-model="redeemInput"
        type="password"
        placeholder="核销密码"
        class="input"
      />
      <p v-if="authError" class="error">{{ authError }}</p>
      <button class="btn primary" :disabled="loading" @click="doAuth">验证</button>
    </div>

    <!-- 内容：按 type 区分预约 / 门票订单 -->
    <div v-else class="content">
      <template v-if="loadError">
        <p class="error">{{ loadError }}</p>
      </template>

      <!-- 预约详情 -->
      <template v-else-if="typeFromUrl === 'reservation' && detail">
        <h2>预约详情</h2>
        <ul class="detail-list">
          <li><span>预约单号</span>{{ detail.reservationNo }}</li>
          <li><span>场地</span>{{ venueName || '—' }}</li>
          <li><span>预约日期</span>{{ formatDate(detail.reserveDate) }}</li>
          <li><span>时段</span>{{ timeslotLabel || '—' }}</li>
          <li><span>联系人</span>{{ detail.contactName }}</li>
          <li><span>联系电话</span>{{ detail.contactPhone }}</li>
          <li><span>预约人数</span>{{ detail.contactCount }} 人</li>
          <li>
            <span>状态</span>
            <span :class="['status', statusClass]">{{ statusText }}</span>
          </li>
        </ul>
        <div v-if="detail.status === 0" class="actions">
          <button class="btn primary" :disabled="verifyLoading" @click="doVerify">确认核销</button>
        </div>
        <div v-else-if="detail.status === 1" class="msg success">该预约已核销</div>
        <div v-else class="msg info">{{ detail.status === 2 ? '该预约已取消' : '该预约已过期' }}</div>
      </template>

      <!-- 门票订单详情 -->
      <template v-else-if="typeFromUrl === 'ticket' && ticketOrder">
        <h2>门票订单</h2>
        <ul class="detail-list">
          <li><span>订单号</span>{{ ticketOrder.orderNo }}</li>
          <li><span>预定人</span>{{ ticketOrder.bookerName }}</li>
          <li><span>手机号</span>{{ ticketOrder.bookerPhone }}</li>
          <li><span>支付金额</span>{{ ticketOrder.payAmount }}</li>
          <li>
            <span>状态</span>
            <span :class="['status', ticketStatusClass(ticketOrder.status)]">{{ ticketStatusText(ticketOrder.status) }}</span>
          </li>
        </ul>
        <ul class="detail-list" v-if="ticketItems && ticketItems.length">
          <li v-for="it in ticketItems" :key="it.ID">
            <span>{{ it.productName }} {{ (it.visitDate || '').slice(0, 10) }}</span>
            <span>x {{ it.quantity }}</span>
          </li>
        </ul>
        <div v-if="ticketOrder.status === 1" class="actions">
          <button class="btn primary" :disabled="verifyLoading" @click="doVerify">确认核销</button>
        </div>
        <div v-else-if="ticketOrder.status === 2" class="msg success">该订单已核销</div>
        <div v-else class="msg info">该订单当前状态不支持核销</div>
      </template>

      <div v-else-if="loading" class="loading">加载中…</div>
      <div v-else class="loading">暂无数据</div>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { ref, computed, onMounted, watch } from 'vue'
import { validateRedeemCode } from '@/api/sysParams'
import { getReservationByVerifyCodePublic, verifyReservationByCodePublic } from '@/plugin/camping/api/reservation'
import { getTicketOrderByCodePublic, verifyTicketOrderByCodePublic } from '@/plugin/ticket/api/order'
import { getSiteDetailPublic } from '@/plugin/camping/api/site'
import { getTimeSlotsByVenuePublic } from '@/plugin/camping/api/timeSlot'

defineOptions({ name: 'H5Verify' })

const STORAGE_KEY = 'camping_staff_verified_until'
const ONE_MONTH_MS = 30 * 24 * 60 * 60 * 1000

const route = useRoute()
// 同步从 URL 取参，避免首屏拿不到 code/type 只显示进度条
const codeFromUrl = ref((route.query?.code && String(route.query.code).trim()) || '')
const typeFromUrl = ref((route.query?.type && String(route.query.type).trim()) || 'reservation')
const redeemInput = ref('')
const passedAuth = ref(false)
const authError = ref('')
const loading = ref(false)
const loadError = ref('')
const detail = ref(null)
const ticketOrder = ref(null)
const ticketItems = ref([])
const venueName = ref('')
const timeslotLabel = ref('')
const verifyLoading = ref(false)

function ticketStatusText(status) {
  const m = {
    0: '待支付',
    1: '待核销',
    2: '已核销',
    3: '已取消',
    4: '已过期',
    5: '已关闭'
  }
  return m[status] ?? ''
}

function ticketStatusClass(status) {
  const s = Number(status)
  if (s === 0 || s === 1) return 'pending'
  if (s === 2) return 'done'
  if (s === 3 || s === 5) return 'cancel'
  if (s === 4) return 'expired'
  return ''
}

const statusText = computed(() => {
  if (!detail.value) return ''
  const m = { 0: '待核销', 1: '已核销', 2: '已取消', 3: '已过期' }
  return m[detail.value.status] ?? ''
})

const statusClass = computed(() => {
  if (!detail.value) return ''
  const c = { 0: 'pending', 1: 'done', 2: 'cancel', 3: 'expired' }
  return c[detail.value.status] ?? ''
})

function formatDate(d) {
  if (!d) return '-'
  if (typeof d === 'string') return d.slice(0, 10)
  return d
}

function checkSavedAuth() {
  try {
    const until = parseInt(localStorage.getItem(STORAGE_KEY), 10)
    if (until && Date.now() < until) passedAuth.value = true
  } catch (_) {}
}

async function doAuth() {
  const code = (redeemInput.value || '').trim()
  if (!code) {
    authError.value = '请输入核销密码'
    return
  }
  authError.value = ''
  loading.value = true
  try {
    const res = await validateRedeemCode({ code })
    if (res.code === 0 && res.data && res.data.valid) {
      localStorage.setItem(STORAGE_KEY, String(Date.now() + ONE_MONTH_MS))
      passedAuth.value = true
      await loadReservation()
    } else {
      authError.value = '密码错误，请重试'
    }
  } catch (_) {
    authError.value = '验证失败，请重试'
  } finally {
    loading.value = false
  }
}

async function loadReservation() {
  if (!codeFromUrl.value || typeFromUrl.value !== 'reservation') return
  loadError.value = ''
  loading.value = true
  detail.value = null
  venueName.value = ''
  timeslotLabel.value = ''
  try {
    const res = await getReservationByVerifyCodePublic({ code: codeFromUrl.value })
    if (res.code === 0 && res.data) {
      detail.value = res.data
      if (res.data.venueId) {
        try {
          const siteRes = await getSiteDetailPublic({ id: res.data.venueId })
          if (siteRes.code === 0 && siteRes.data) {
            venueName.value = siteRes.data.name || siteRes.data.Name || ''
          }
        } catch (_) {}
        // 加载该场地时段以解析预约时段显示
        if (res.data.timeslotId) {
          try {
            const slotRes = await getTimeSlotsByVenuePublic({ venueId: res.data.venueId })
            if (slotRes.code === 0 && Array.isArray(slotRes.data)) {
              const s = slotRes.data.find((x) => x.ID === res.data.timeslotId)
              if (s) {
                const start = (s.startTime || '').slice(0, 5)
                const end = (s.endTime || '').slice(0, 5)
                timeslotLabel.value = start && end ? `${start}-${end}` : ''
              }
            }
          } catch (_) {}
        }
      }
    } else {
      loadError.value = res.msg || '预约不存在或已失效'
    }
  } catch (_) {
    loadError.value = '加载失败，请检查网络'
  } finally {
    loading.value = false
  }
}

async function loadTicketOrder() {
  if (!codeFromUrl.value || typeFromUrl.value !== 'ticket') return
  loadError.value = ''
  loading.value = true
  ticketOrder.value = null
  ticketItems.value = []
  try {
    const res = await getTicketOrderByCodePublic({ code: codeFromUrl.value })
    if (res.code === 0 && res.data) {
      ticketOrder.value = res.data.order
      ticketItems.value = res.data.items || []
    } else {
      loadError.value = res.msg || '订单不存在或已失效'
    }
  } catch (_) {
    loadError.value = '加载失败，请检查网络'
  } finally {
    loading.value = false
  }
}

async function doVerify() {
  if (!codeFromUrl.value) return
  verifyLoading.value = true
  try {
    if (typeFromUrl.value === 'reservation') {
      if (!detail.value || detail.value.status !== 0) return
      const res = await verifyReservationByCodePublic({ code: codeFromUrl.value })
      if (res.code === 0) {
        if (res.data) detail.value = res.data
        else detail.value = { ...detail.value, status: 1 }
      } else {
        loadError.value = res.msg || '核销失败'
      }
    } else if (typeFromUrl.value === 'ticket') {
      if (!ticketOrder.value || ticketOrder.value.status !== 1) return
      const res = await verifyTicketOrderByCodePublic({ code: codeFromUrl.value })
      if (res.code === 0 && res.data) {
        ticketOrder.value = res.data.order || ticketOrder.value
        ticketItems.value = res.data.items || ticketItems.value
      } else if (res.code === 0 && !res.data) {
        ticketOrder.value = { ...ticketOrder.value, status: 2 }
      } else {
        loadError.value = res.msg || '核销失败'
      }
    }
  } catch (_) {
    loadError.value = '核销请求失败，请重试'
  } finally {
    verifyLoading.value = false
  }
}

function applyQueryFromRoute() {
  const q = route.query || {}
  codeFromUrl.value = (q.code && String(q.code).trim()) || ''
  typeFromUrl.value = (q.type && String(q.type).trim()) || 'reservation'
}

watch(() => route.query, () => applyQueryFromRoute(), { deep: true })

onMounted(() => {
  applyQueryFromRoute()
  checkSavedAuth()
  if (codeFromUrl.value && passedAuth.value) {
    if (typeFromUrl.value === 'ticket') {
      loadTicketOrder()
    } else {
      loadReservation()
    }
  }
})
</script>

<style scoped>
.h5-verify {
  min-height: 100vh;
  height: 100%;
  padding: 24px 16px;
  box-sizing: border-box;
  background: #f5f5f5;
  font-size: 14px;
  overflow: auto;
}
.app-title {
  text-align: center;
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 16px;
}
.tip,
.auth-box,
.content {
  max-width: 400px;
  margin: 0 auto;
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}
.tip p,
.hint {
  color: #666;
  text-align: center;
  margin: 0 0 16px;
}
.auth-box h2,
.content h2 {
  margin: 0 0 16px;
  font-size: 18px;
  text-align: center;
}
.input {
  width: 100%;
  padding: 12px 16px;
  margin-bottom: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  box-sizing: border-box;
}
.error {
  color: #f56c6c;
  font-size: 12px;
  margin: 0 0 12px;
}
.btn {
  width: 100%;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
}
.btn.primary {
  background: #409eff;
  color: #fff;
}
.btn.primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.detail-list {
  list-style: none;
  padding: 0;
  margin: 0 0 20px;
}
.detail-list li {
  display: flex;
  padding: 10px 0;
  border-bottom: 1px solid #eee;
}
.detail-list li span:first-child {
  width: 80px;
  color: #666;
  flex-shrink: 0;
}
.status.pending { color: #e6a23c; }
.status.done { color: #67c23a; }
.status.cancel { color: #909399; }
.status.expired { color: #f56c6c; }
.actions { margin-top: 20px; }
.msg { text-align: center; padding: 12px; border-radius: 8px; }
.msg.success { background: #f0f9eb; color: #67c23a; }
.msg.info { background: #f4f4f5; color: #909399; }
.loading { text-align: center; padding: 24px; color: #666; }
</style>
