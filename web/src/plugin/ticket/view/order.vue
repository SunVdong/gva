<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="订单号">
          <el-input v-model="searchInfo.orderNo" placeholder="订单号" clearable />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input-number v-model="searchInfo.userId" :min="0" placeholder="用户ID" clearable style="width: 120px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 110px">
            <el-option label="待支付" :value="0" />
            <el-option label="已支付" :value="1" />
            <el-option label="已退款" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="onSubmit">查询</el-button>
          <el-button icon="Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table :data="tableData" row-key="ID">
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="订单号" prop="orderNo" min-width="160" show-overflow-tooltip />
        <el-table-column align="left" label="用户ID" prop="userId" width="90" />
        <el-table-column align="left" label="订单金额" width="100">
          <template #default="{ row }">
            ¥{{ (row.totalAmount ?? 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="支付金额" width="100">
          <template #default="{ row }">
            ¥{{ (row.payAmount ?? 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="90">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="warning">待支付</el-tag>
            <el-tag v-else-if="row.status === 1" type="success">已支付</el-tag>
            <el-tag v-else-if="row.status === 2" type="info">已退款</el-tag>
            <el-tag v-else>未知</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="支付时间" width="170">
          <template #default="{ row }">
            {{ row.payTime ? formatDate(row.payTime) : '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建时间" width="170">
          <template #default="{ row }">
            {{ row.CreatedAt ? formatDate(row.CreatedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" width="100">
          <template #default="{ row }">
            <el-button type="primary" link @click="showDetail(row)">详情</el-button>
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
    <el-drawer v-model="detailVisible" destroy-on-close size="640" title="订单详情" :show-close="true">
      <template #header>
        <span class="text-lg">订单详情 · {{ detail.order?.orderNo || '' }}</span>
      </template>
      <div v-if="detail.order" class="space-y-4">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="订单号">{{ detail.order.orderNo }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ detail.order.userId }}</el-descriptions-item>
          <el-descriptions-item label="订单金额">¥{{ (detail.order.totalAmount ?? 0).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="支付金额">¥{{ (detail.order.payAmount ?? 0).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag v-if="detail.order.status === 0" type="warning">待支付</el-tag>
            <el-tag v-else-if="detail.order.status === 1" type="success">已支付</el-tag>
            <el-tag v-else-if="detail.order.status === 2" type="info">已退款</el-tag>
            <span v-else>未知</span>
          </el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ detail.order.payTime ? formatDate(detail.order.payTime) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ detail.order.CreatedAt ? formatDate(detail.order.CreatedAt) : '-' }}</el-descriptions-item>
        </el-descriptions>
        <div class="text-sm font-medium mt-4">订单明细</div>
        <el-table :data="detail.items" border size="small">
          <el-table-column label="门票名称" prop="skuName" min-width="120" />
          <el-table-column label="单价" width="90">
            <template #default="{ row }">¥{{ (row.price ?? 0).toFixed(2) }}</template>
          </el-table-column>
          <el-table-column label="数量" prop="quantity" width="70" />
          <el-table-column label="小计" width="90">
            <template #default="{ row }">¥{{ ((row.price ?? 0) * (row.quantity ?? 0)).toFixed(2) }}</template>
          </el-table-column>
          <el-table-column label="游玩日期" width="120">
            <template #default="{ row }">{{ row.visitDate ? formatVisitDate(row.visitDate) : '-' }}</template>
          </el-table-column>
        </el-table>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { getOrderList, findOrder } from '@/plugin/ticket/api/order'
import { ref } from 'vue'

defineOptions({ name: 'TicketOrder' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const detailVisible = ref(false)
const detail = ref({ order: null, items: [] })

function formatDate(d) {
  if (!d) return ''
  const t = typeof d === 'string' ? d : (d && d.toISOString ? d.toISOString() : '')
  return t ? t.slice(0, 19).replace('T', ' ') : ''
}

function formatVisitDate(d) {
  if (!d) return ''
  const t = typeof d === 'string' ? d : (d && d.toISOString ? d.toISOString() : '')
  return t ? t.slice(0, 10) : ''
}

const getTableData = async () => {
  const res = await getOrderList({
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

const showDetail = async (row) => {
  const res = await findOrder({ id: row.ID })
  if (res.code === 0 && res.data) {
    detail.value = { order: res.data.order, items: res.data.items || [] }
    detailVisible.value = true
  }
}

getTableData()
</script>
