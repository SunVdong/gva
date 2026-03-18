<template>
  <div class="h5-verify">
    <h1 class="app-title">乐享江岛核销系统</h1>
    <!-- 无核销码 -->
    <div v-if="!codeFromUrl" class="tip">
      <p>请扫描预约二维码进入核销页</p>
    </div>

    <!-- 身份验证（8 位数字），验证通过后 1 个月内不再要求输入 -->
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

    <!-- 预约详情 + 核销 -->
    <div v-else class="content">
      <template v-if="loadError">
        <p class="error">{{ loadError }}</p>
      </template>
      <template v-else-if="detail">
        <h2>预约详情</h2>
        <ul class="detail-list">
          <li><span>预约单号</span>{{ detail.reservationNo }}</li>
          <li><span>场地</span>{{ venueName || '—' }}</li>
          <li><span>预约日期</span>{{ formatDate(detail.reserveDate) }}</li>
          <li><span>时段</span>{{ timeslotLabel || '—' }}</li>
          <li><span>联系人</span>{{ detail.contactName }}</li>
          <li><span>联系电话</span>{{ detail.contactPhone }}</li>
          <li><span>预约人数</span>{{ detail.contactCount }} 人</li>
          <li><span>状态</span>
            <span :class="['status', statusClass]">{{ statusText }}</span>
          </li>
        </ul>
        <div v-if="detail.status === 0" class="actions">
          <button class="btn primary" :disabled="verifyLoading" @click="doVerify">确认核销</button>
        </div>
        <div v-else-if="detail.status === 1" class="msg success">该预约已核销</div>
        <div v-else class="msg info">{{ detail.status === 2 ? '该预约已取消' : '该预约已过期' }}</div>
      </template>
      <div v-else-if="loading" class="loading">加载中…</div>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { ref, computed, onMounted } from 'vue'
import { validateRedeemCode } from '@/api/sysParams'
import { getReservationByVerifyCodePublic, verifyReservationByCodePublic } from '@/plugin/camping/api/reservation'
import { getSiteDetailPublic } from '@/plugin/camping/api/site'
import { getTimeSlotsByVenuePublic } from '@/plugin/camping/api/timeSlot'

defineOptions({ name: 'H5Verify' })

const STORAGE_KEY = 'camping_staff_verified_until'
const ONE_MONTH_MS = 30 * 24 * 60 * 60 * 1000

const route = useRoute()
const codeFromUrl = ref('')
const redeemInput = ref('')
const passedAuth = ref(false)
const authError = ref('')
const loading = ref(false)
const loadError = ref('')
const detail = ref(null)
const venueName = ref('')
const timeslotLabel = ref('')
const verifyLoading = ref(false)

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
  if (!codeFromUrl.value) return
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

async function doVerify() {
  if (!codeFromUrl.value || !detail.value || detail.value.status !== 0) return
  verifyLoading.value = true
  try {
    const res = await verifyReservationByCodePublic({ code: codeFromUrl.value })
    if (res.code === 0) {
      if (res.data) detail.value = res.data
      else detail.value = { ...detail.value, status: 1 }
    } else {
      loadError.value = res.msg || '核销失败'
    }
  } catch (_) {
    loadError.value = '核销请求失败，请重试'
  } finally {
    verifyLoading.value = false
  }
}

onMounted(() => {
  codeFromUrl.value = (route.query.code || '').trim()
  checkSavedAuth()
  if (codeFromUrl.value && passedAuth.value) loadReservation()
})
</script>

<style scoped>
.h5-verify {
  min-height: 100vh;
  padding: 24px 16px;
  box-sizing: border-box;
  background: #f5f5f5;
  font-size: 14px;
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
