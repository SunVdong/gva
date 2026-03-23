<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="场地">
          <el-select v-model="searchInfo.venueId" placeholder="全部" clearable style="width: 160px">
            <el-option v-for="s in siteOptions" :key="s.ID" :label="s.name" :value="s.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="预约日期">
          <el-date-picker v-model="searchInfo.reserveDate" type="date" value-format="YYYY-MM-DD" placeholder="日期" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="待核销" :value="0" />
            <el-option label="已核销" :value="1" />
            <el-option label="已取消" :value="2" />
            <el-option label="已过期" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="核销码">
          <el-input v-model="searchInfo.verifyCode" placeholder="核销码" clearable style="width: 140px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="onSubmit">查询</el-button>
          <el-button icon="Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table style="width: 100%" :data="tableData" row-key="ID">
        <el-table-column align="left" label="ID" prop="ID" width="70" />
        <el-table-column align="left" label="预约单号" prop="reservationNo" width="200" />
        <el-table-column align="left" label="场地" min-width="100">
          <template #default="{ row }">{{ siteName(row.venueId) }}</template>
        </el-table-column>
        <el-table-column align="left" label="预约日期" width="120">
          <template #default="{ row }">{{ formatDate(row.reserveDate) }}</template>
        </el-table-column>
        <el-table-column align="left" label="时段" width="120">
          <template #default="{ row }">{{ slotLabel(row.timeslotId) }}</template>
        </el-table-column>
        <el-table-column align="left" label="联系人" prop="contactName" width="100" />
        <el-table-column align="left" label="联系电话" prop="contactPhone" width="120" />
        <el-table-column align="left" label="人数" prop="contactCount" width="70" />
        <el-table-column align="left" label="核销码" prop="verifyCode" width="130" />
        <el-table-column align="left" label="核销时间" width="160">
          <template #default="{ row }">{{ formatDateTime(row.verifiedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'warning' : row.status === 1 ? 'success' : row.status === 2 ? 'info' : 'danger'">
              {{
                row.status === 0
                  ? '待核销'
                  : row.status === 1
                    ? '已核销'
                    : row.status === 2
                      ? '已取消'
                      : '已过期'
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button v-if="row.status === 0" type="primary" link icon="Check" @click="showQr(row)">二维码</el-button>
            <el-button v-if="row.status === 0" type="primary" link icon="Close" @click="cancelRow(row)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <!-- 二维码弹窗 -->
    <el-dialog v-model="qrDialogVisible" title="预约成功 - 核销二维码" width="400px" align-center>
      <div class="text-center">
        <vue-qr :text="qrCodeText" :size="240" class="mx-auto" />
        <p class="mt-2 font-mono text-sm">核销码：{{ currentReservation?.verifyCode }}</p>
        <p class="text-gray-500 text-xs">请出示此二维码供工作人员核销</p>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { getSiteList } from '@/plugin/camping/api/site'
import { getTimeSlotsByVenuePublic } from '@/plugin/camping/api/timeSlot'
import { getReservationList, cancelReservation } from '@/plugin/camping/api/reservation'
import vueQr from 'vue-qr/src/packages/vue-qr.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted } from 'vue'

defineOptions({ name: 'CampingReservation' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const siteOptions = ref([])
const allSlotsForList = ref([])
const qrDialogVisible = ref(false)
const currentReservation = ref(null)
const qrCodeText = ref('')

function siteName(id) {
  const s = siteOptions.value.find((x) => x.ID === id)
  return s ? s.name : id || '-'
}
function slotLabel(timeslotId) {
  const s = allSlotsForList.value.find((x) => x.ID === timeslotId)
  if (!s) return timeslotId || '-'
  return `${s.startTime?.slice(0, 5) || s.startTime}-${s.endTime?.slice(0, 5) || s.endTime}`
}
function formatDate(d) {
  if (!d) return '-'
  if (typeof d === 'string') return d.slice(0, 10)
  return d
}

function formatDateTime(d) {
  if (!d) return '-'
  const v = typeof d === 'string' ? d : ''
  if (!v) return '-'
  // 统一显示到分钟，例如 2025-01-02 15:04
  return v.replace('T', ' ').slice(0, 16)
}

const getTableData = async () => {
  const res = await getReservationList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = res.data.page || page.value
    pageSize.value = res.data.pageSize || pageSize.value
  }
}

const onSubmit = () => { page.value = 1; getTableData() }
const onReset = () => { searchInfo.value = {}; getTableData() }
const handleCurrentChange = (val) => { page.value = val; getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; getTableData() }

/** 生成核销页链接：{域名}{path}#/h5/verify?type=reservation&code=核销码 */
function getVerifyPageUrl(code) {
  const base = window.location.origin + (window.location.pathname || '').replace(/\/$/, '')
  return `${base}#/h5/verify?type=reservation&code=${encodeURIComponent(code || '')}`
}

const showQr = (row) => {
  currentReservation.value = row
  qrCodeText.value = getVerifyPageUrl(row.verifyCode || '')
  qrDialogVisible.value = true
}

const cancelRow = (row) => {
  ElMessageBox.confirm('确定取消该预约？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const res = await cancelReservation({ id: row.ID })
      if (res.code === 0) { ElMessage.success('已取消'); getTableData() }
    })
}

onMounted(async () => {
  const siteRes = await getSiteList({ page: 1, pageSize: 500 })
  if (siteRes.code === 0) {
    siteOptions.value = siteRes.data.list || []
    const all = []
    for (const v of siteOptions.value) {
      const slotRes = await getTimeSlotsByVenuePublic({ venueId: v.ID })
      if (slotRes.code === 0 && slotRes.data?.length) all.push(...slotRes.data)
    }
    allSlotsForList.value = all
  }
  getTableData()
})
</script>
