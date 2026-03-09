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
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 100px">
            <el-option label="待确认" :value="0" />
            <el-option label="已预约" :value="1" />
            <el-option label="已取消" :value="2" />
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
      <div class="gva-btn-list">
        <el-button type="primary" icon="Plus" @click="openReserveDialog">新增预约</el-button>
      </div>
      <el-table style="width: 100%" :data="tableData" row-key="ID">
        <el-table-column align="left" label="ID" prop="ID" width="70" />
        <el-table-column align="left" label="预约单号" prop="reservationNo" width="140" />
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
        <el-table-column align="left" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'warning' : row.status === 1 ? 'success' : 'info'">
              {{ row.status === 0 ? '待确认' : row.status === 1 ? '已预约' : '已取消' }}
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
    <!-- 新增预约 -->
    <el-drawer v-model="reserveDialogVisible" title="新增预约" size="500" destroy-on-close>
      <el-form ref="reserveFormRef" :model="reserveForm" label-position="top" :rules="reserveRules">
        <el-form-item label="场地" prop="venueId">
          <el-select v-model="reserveForm.venueId" placeholder="请选择场地" style="width: 100%" filterable @change="onVenueChange">
            <el-option v-for="s in siteOptions" :key="s.ID" :label="s.name" :value="s.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="预约日期" prop="reserveDate">
          <el-date-picker v-model="reserveForm.reserveDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" style="width: 100%" @change="onDateChange" />
        </el-form-item>
        <el-form-item label="时段" prop="timeslotId">
          <el-select v-model="reserveForm.timeslotId" placeholder="请先选择场地" style="width: 100%">
            <el-option
              v-for="s in slotOptions"
              :key="s.ID"
              :label="`${s.startTime?.slice(0,5) || s.startTime}-${s.endTime?.slice(0,5) || s.endTime}（可约${s.capacity}）`"
              :value="s.ID"
              :disabled="reservedSlotIds.includes(s.ID)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="联系人" prop="contactName">
          <el-input v-model="reserveForm.contactName" placeholder="联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contactPhone">
          <el-input v-model="reserveForm.contactPhone" placeholder="手机号" />
        </el-form-item>
        <el-form-item label="预约人数" prop="contactCount">
          <el-input-number v-model="reserveForm.contactCount" :min="1" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="reserveForm.remark" type="textarea" :rows="2" placeholder="选填" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitReserve">提交预约</el-button>
        </el-form-item>
      </el-form>
    </el-drawer>
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
import { createReservation, getReservationList, cancelReservation, getReservedSlotIdsPublic } from '@/plugin/camping/api/reservation'
import vueQr from 'vue-qr/src/packages/vue-qr.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({ name: 'CampingReservation' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const siteOptions = ref([])
const slotOptions = ref([])
const allSlotsForList = ref([])
const reservedSlotIds = ref([])
const reserveDialogVisible = ref(false)
const reserveFormRef = ref()
const reserveForm = reactive({
  venueId: null,
  reserveDate: '',
  timeslotId: null,
  contactName: '',
  contactPhone: '',
  contactCount: 1,
  remark: ''
})
const reserveRules = {
  venueId: [{ required: true, message: '请选择场地', trigger: 'change' }],
  reserveDate: [{ required: true, message: '请选择日期', trigger: 'change' }],
  timeslotId: [{ required: true, message: '请选择时段', trigger: 'change' }],
  contactName: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contactPhone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  contactCount: [{ required: true, message: '请输入人数', trigger: 'blur' }]
}
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

const onVenueChange = async () => {
  reserveForm.timeslotId = null
  slotOptions.value = []
  reservedSlotIds.value = []
  if (reserveForm.venueId) {
    const res = await getTimeSlotsByVenuePublic({ venueId: reserveForm.venueId })
    if (res.code === 0) slotOptions.value = res.data || []
  }
}

const onDateChange = async () => {
  reservedSlotIds.value = []
  if (reserveForm.venueId && reserveForm.reserveDate) {
    const res = await getReservedSlotIdsPublic({ venueId: reserveForm.venueId, reserveDate: reserveForm.reserveDate })
    if (res.code === 0) reservedSlotIds.value = res.data || []
  }
}

const openReserveDialog = () => {
  reserveForm.venueId = null
  reserveForm.reserveDate = ''
  reserveForm.timeslotId = null
  reserveForm.contactName = ''
  reserveForm.contactPhone = ''
  reserveForm.contactCount = 1
  reserveForm.remark = ''
  slotOptions.value = []
  reservedSlotIds.value = []
  reserveDialogVisible.value = true
}

const submitReserve = async () => {
  await reserveFormRef.value?.validate(async (valid) => {
    if (!valid) return
    const res = await createReservation({
      venueId: reserveForm.venueId,
      reserveDate: reserveForm.reserveDate,
      timeslotId: reserveForm.timeslotId,
      contactName: reserveForm.contactName,
      contactPhone: reserveForm.contactPhone,
      contactCount: reserveForm.contactCount,
      remark: reserveForm.remark
    })
    if (res.code === 0) {
      ElMessage.success('预约成功')
      reserveDialogVisible.value = false
      currentReservation.value = res.data
      qrCodeText.value = res.data.verifyCode || ''
      qrDialogVisible.value = true
      getTableData()
    } else {
      ElMessage.warning(res.msg || '预约失败')
    }
  })
}

const showQr = (row) => {
  currentReservation.value = row
  qrCodeText.value = row.verifyCode || ''
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
