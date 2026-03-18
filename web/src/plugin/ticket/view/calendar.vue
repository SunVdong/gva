<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" class="demo-form-inline">
        <el-form-item label="门票 SKU">
          <el-select
            v-model="currentSkuId"
            placeholder="请先选择门票 SKU"
            filterable
            style="width: 280px"
            @change="onSkuChange"
          >
            <el-option
              v-for="s in skuOptions"
              :key="s.ID"
              :label="`${s.name} (¥${s.price}) - ${productName(s.productId)}`"
              :value="s.ID"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" :disabled="!currentSkuId" @click="loadCalendar">查询</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div v-if="currentSkuId" class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="Plus" @click="openBatchSet">批量设置库存</el-button>
      </div>
      <el-table :data="calendarList" border size="small">
        <el-table-column label="游玩日期" width="240">
          <template #default="{ row }">
            {{ formatDate(row.visitDate) }}
          </template>
        </el-table-column>
        <el-table-column label="库存" prop="stock" width="100" />
        <el-table-column label="已售" prop="sold" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '可售' : '关闭' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="160">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="editOne(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, prev, pager, next"
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <el-dialog v-model="batchDialogVisible" title="批量设置库存" width="500" destroy-on-close>
      <el-form label-width="100px">
        <el-form-item label="开始日期">
          <el-date-picker v-model="batchStart" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker v-model="batchEnd" type="date" value-format="YYYY-MM-DD" style="width:100%" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="batchStock" :min="0" style="width:100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="batchStatus">
            <el-radio :value="1">可售</el-radio>
            <el-radio :value="0">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBatch">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="oneDialogVisible" title="编辑单日库存" width="400" destroy-on-close>
      <el-form label-width="100px">
        <el-form-item label="日期">{{ formatDate(editRow.visitDate) }}</el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="editRow.stock" :min="0" style="width:100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="editRow.status">
            <el-radio :value="1">可售</el-radio>
            <el-radio :value="0">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="oneDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitOne">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { getProductList, getSkuList, getCalendarBySku, setCalendar } from '@/plugin/ticket/api/product'
import { getScenicList } from '@/plugin/ticket/api/scenic'
import { ElMessage } from 'element-plus'
import { ref, onMounted } from 'vue'

defineOptions({ name: 'TicketCalendar' })

const currentSkuId = ref(null)
const dateRange = ref([])
const skuOptions = ref([])
const productMap = ref({})
const scenicMap = ref({})
const calendarList = ref([])
const page = ref(1)
const total = ref(0)
const pageSize = ref(31)
const batchDialogVisible = ref(false)
const batchStart = ref('')
const batchEnd = ref('')
const batchStock = ref(0)
const batchStatus = ref(1)
const oneDialogVisible = ref(false)
const editRow = ref({})

function productName(productId) {
  const p = productMap.value[productId]
  if (!p) return ''
  const s = scenicMap.value[p.scenicId]
  return s ? s.name : ''
}

function formatDate(d) {
  if (!d) return ''
  const s = typeof d === 'string'
    ? d
    : (d.year && d.month ? `${d.year}-${String(d.month).padStart(2, '0')}-${String(d.day || d.date || 1).padStart(2, '0')}` : '')
  if (s) return s.slice(0, 10)
  if (d instanceof Date) return d.toISOString().slice(0, 10)
  return String(d).slice(0, 10)
}

const loadSkuOptions = async () => {
  const [productRes, skuRes] = await Promise.all([
    getProductList({ page: 1, pageSize: 500 }),
    getSkuList({ page: 1, pageSize: 500 })
  ])
  if (productRes.code === 0 && productRes.data.list) {
    productRes.data.list.forEach((p) => { productMap.value[p.ID] = p })
  }
  const scenicRes = await getScenicList({ page: 1, pageSize: 500 })
  if (scenicRes.code === 0 && scenicRes.data.list) {
    scenicRes.data.list.forEach((s) => { scenicMap.value[s.ID] = s })
  }
  if (skuRes.code === 0 && skuRes.data.list) {
    skuOptions.value = skuRes.data.list
  }
}

const onSkuChange = () => {
  if (!currentSkuId.value || !dateRange.value || dateRange.value.length !== 2) return
  loadCalendar()
}

const loadCalendar = async () => {
  if (!currentSkuId.value) return
  const [start, end] = dateRange.value && dateRange.value.length === 2 ? dateRange.value : []
  const res = await getCalendarBySku({
    skuId: currentSkuId.value,
    startDate: start || '',
    endDate: end || '',
    page: page.value,
    pageSize: pageSize.value
  })
  if (res.code === 0) {
    calendarList.value = res.data.list || []
    total.value = res.data.total || 0
  }
}

const handlePageChange = (val) => {
  page.value = val
  loadCalendar()
}

const openBatchSet = () => {
  batchStart.value = ''
  batchEnd.value = ''
  batchStock.value = 0
  batchStatus.value = 1
  batchDialogVisible.value = true
}

const submitBatch = async () => {
  if (!batchStart.value || !batchEnd.value) {
    ElMessage.warning('请选择开始和结束日期')
    return
  }
  const list = []
  const start = new Date(batchStart.value)
  const end = new Date(batchEnd.value)
  for (let d = new Date(start); d <= end; d.setDate(d.getDate() + 1)) {
    list.push({
      skuId: currentSkuId.value,
      visitDate: d.toISOString().slice(0, 10),
      stock: batchStock.value,
      status: batchStatus.value
    })
  }
  const res = await setCalendar({ list })
  if (res.code === 0) {
    ElMessage.success('设置成功')
    batchDialogVisible.value = false
    loadCalendar()
  }
}

const editOne = (row) => {
  editRow.value = { ...row, visitDate: row.visitDate }
  oneDialogVisible.value = true
}

const submitOne = async () => {
  const d = formatDate(editRow.value.visitDate)
  const res = await setCalendar({
    list: [{
      skuId: currentSkuId.value,
      visitDate: d,
      stock: editRow.value.stock,
      status: editRow.value.status
    }]
  })
  if (res.code === 0) {
    ElMessage.success('已更新')
    oneDialogVisible.value = false
    loadCalendar()
  }
}

onMounted(() => {
  loadSkuOptions()
  const today = new Date()
  const next = new Date(today)
  next.setDate(next.getDate() + 30)
  dateRange.value = [today.toISOString().slice(0, 10), next.toISOString().slice(0, 10)]
})
</script>
