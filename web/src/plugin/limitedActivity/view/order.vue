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
        <el-form-item label="联系电话">
          <el-input v-model="searchInfo.contactPhone" placeholder="联系电话" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 130px">
            <el-option label="待支付" :value="0" />
            <el-option label="待核销" :value="1" />
            <el-option label="已核销" :value="2" />
            <el-option label="已关闭" :value="5" />
            <el-option label="已退款" :value="6" />
            <el-option label="退款中" :value="7" />
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
        <el-table-column align="left" label="活动" prop="activityName" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="联系人" prop="contactName" width="110" show-overflow-tooltip />
        <el-table-column align="left" label="联系电话" prop="contactPhone" width="130" show-overflow-tooltip />
        <el-table-column align="left" label="用户ID" prop="userId" width="90" />
        <el-table-column align="left" label="人次" prop="quantity" width="80" />
        <el-table-column align="left" label="核销进度" width="110">
          <template #default="{ row }">
            {{ row.verifiedTimes ?? 0 }}/{{ row.totalUseTimes ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="支付金额" width="100">
          <template #default="{ row }">
            ¥{{ (row.payAmount ?? 0).toFixed(2) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 0" type="warning">待支付</el-tag>
            <el-tag v-else-if="row.status === 1 && (row.verifiedTimes || 0) > 0" type="warning">核销中</el-tag>
            <el-tag v-else-if="row.status === 1" type="primary">待核销</el-tag>
            <el-tag v-else-if="row.status === 2" type="success">已核销</el-tag>
            <el-tag v-else-if="row.status === 5" type="info">已关闭</el-tag>
            <el-tag v-else-if="row.status === 6" type="success">已退款</el-tag>
            <el-tag v-else-if="row.status === 7" type="warning">退款中</el-tag>
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
        <el-table-column align="left" label="操作" fixed="right" width="140">
          <template #default="{ row }">
            <el-button type="primary" link @click="showDetail(row)">详情</el-button>
            <el-popconfirm
              v-if="canShowRefund(row)"
              :title="refundConfirmTitle(row)"
              @confirm="handleRefund(row)"
            >
              <template #reference>
                <el-button type="warning" link>退款</el-button>
              </template>
            </el-popconfirm>
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
          <el-descriptions-item label="活动">{{ detail.order.activityName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="用户ID">{{ detail.order.userId }}</el-descriptions-item>
          <el-descriptions-item label="联系人">{{ detail.order.contactName || '-' }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ detail.order.contactPhone || '-' }}</el-descriptions-item>
          <el-descriptions-item label="人次">{{ detail.order.quantity ?? '-' }}</el-descriptions-item>
          <el-descriptions-item label="单价">¥{{ (detail.order.unitPrice ?? 0).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="支付金额">¥{{ (detail.order.payAmount ?? 0).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag v-if="detail.order.status === 0" type="warning">待支付</el-tag>
            <el-tag v-else-if="detail.order.status === 1 && (detail.order.verifiedTimes || 0) > 0" type="warning">核销中</el-tag>
            <el-tag v-else-if="detail.order.status === 1" type="primary">待核销</el-tag>
            <el-tag v-else-if="detail.order.status === 2" type="success">已核销</el-tag>
            <el-tag v-else-if="detail.order.status === 5" type="info">已关闭</el-tag>
            <el-tag v-else-if="detail.order.status === 6" type="success">已退款</el-tag>
            <el-tag v-else-if="detail.order.status === 7" type="warning">退款中</el-tag>
            <span v-else>未知</span>
          </el-descriptions-item>
          <el-descriptions-item label="核销进度">
            已核销 {{ detail.order.verifiedTimes || 0 }}/{{ detail.order.totalUseTimes || 0 }} 次，剩余 {{ detail.remainingTimes ?? 0 }} 次
          </el-descriptions-item>
          <el-descriptions-item v-if="detail.canRefund" label="预计退款">
            ¥{{ (detail.refundAmount ?? 0).toFixed(2) }}（按未核销比例）
          </el-descriptions-item>
          <el-descriptions-item v-if="detail.order.status === 6" label="实退金额">
            ¥{{ (detail.order.refundAmount ?? 0).toFixed(2) }}
          </el-descriptions-item>
          <el-descriptions-item label="支付时间">{{ detail.order.payTime ? formatDate(detail.order.payTime) : '-' }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ detail.order.CreatedAt ? formatDate(detail.order.CreatedAt) : '-' }}</el-descriptions-item>
        </el-descriptions>

        <div v-if="detail.verifyRecords && detail.verifyRecords.length" class="mt-4">
          <div class="text-sm font-medium mb-2">核销记录</div>
          <el-table :data="detail.verifyRecords" border size="small">
            <el-table-column label="次序" prop="verifyNo" width="70" />
            <el-table-column label="核销时间" min-width="170">
              <template #default="{ row }">{{ row.verifiedAt ? formatDate(row.verifiedAt) : '-' }}</template>
            </el-table-column>
            <el-table-column label="备注" prop="remark" min-width="120">
              <template #default="{ row }">{{ row.remark || '-' }}</template>
            </el-table-column>
          </el-table>
        </div>

        <div v-if="detail.canRefund" class="mt-4">
          <el-popconfirm :title="refundConfirmTitle(detail.order)" @confirm="handleRefund(detail.order)">
            <template #reference>
              <el-button type="warning">按未核销比例退款</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { findOrder, getOrderList, refundOrder } from '@/plugin/limitedActivity/api/order'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

defineOptions({ name: 'LimitedActivityOrder' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const detailVisible = ref(false)
const detail = ref({ order: null, verifyRecords: [], remainingTimes: 0, canRefund: false, refundAmount: 0 })

function formatDate(d) {
  if (!d) return ''
  const t = typeof d === 'string' ? d : (d && d.toISOString ? d.toISOString() : '')
  return t ? t.slice(0, 19).replace('T', ' ') : ''
}

function canShowRefund(row) {
  if (!row || row.status !== 1 || !row.wxTransactionId) return false
  const total = row.totalUseTimes || 0
  const verified = row.verifiedTimes || 0
  return total - verified > 0
}

function refundConfirmTitle(row) {
  if (!row) return '确定退款吗？'
  const total = row.totalUseTimes || 1
  const remaining = Math.max(0, total - (row.verifiedTimes || 0))
  const amount = ((row.payAmount || 0) * remaining / total).toFixed(2)
  return `按未核销 ${remaining}/${total} 人次退款约 ¥${amount}，确定继续？`
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
    detail.value = {
      order: res.data.order,
      verifyRecords: res.data.verifyRecords || [],
      remainingTimes: res.data.remainingTimes ?? 0,
      canRefund: !!res.data.canRefund,
      refundAmount: res.data.refundAmount ?? 0
    }
    detailVisible.value = true
  }
}

const handleRefund = async (row) => {
  const res = await refundOrder({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success(res.msg || '退款成功或已受理')
    if (detailVisible.value && detail.value.order?.ID === row.ID) {
      await showDetail(row)
    }
    await getTableData()
  }
}

getTableData()
</script>
