<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="活动名称">
          <el-input v-model="searchInfo.name" placeholder="名称" clearable />
        </el-form-item>
        <el-form-item label="显示状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 100px">
            <el-option label="显示" :value="1" />
            <el-option label="隐藏" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="onSubmit">查询</el-button>
          <el-button icon="Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="Plus" @click="openDialog">新增活动</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="活动名称" prop="name" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="活动地点" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">
            {{ row.address || '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="封面" width="80">
          <template #default="{ row }">
            <el-image
              v-if="row.coverImage"
              :src="getUrl(row.coverImage)"
              style="width:48px;height:48px"
              fit="cover"
              :preview-src-list="[getUrl(row.coverImage)]"
              :preview-teleported="true"
              :z-index="9999"
            />
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="活动时间" min-width="200">
          <template #default="{ row }">
            {{ formatDate(row.startTime) }} ~ {{ formatDate(row.endTime) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="价格" width="120">
          <template #default="{ row }">
            ¥{{ (row.price ?? 0).toFixed(2) }}
            <span class="text-gray-400 text-xs">/市¥{{ (row.marketPrice ?? 0).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="名额" width="110">
          <template #default="{ row }">
            {{ row.sold ?? 0 }}/{{ row.quota ?? 0 }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="显示" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '显示' : '隐藏' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建人" prop="createdBy" width="80" />
        <el-table-column align="left" label="更新人" prop="updatedBy" width="80" />
        <el-table-column align="left" label="创建时间" width="170">
          <template #default="{ row }">
            {{ row.CreatedAt ? formatDate(row.CreatedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" width="170">
          <template #default="{ row }">
            {{ row.UpdatedAt ? formatDate(row.UpdatedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="160">
          <template #default="{ row }">
            <el-button type="primary" link icon="Edit" @click="updateFunc(row)">编辑</el-button>
            <el-button type="primary" link icon="Delete" @click="deleteRow(row)">删除</el-button>
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

    <el-drawer v-model="dialogVisible" destroy-on-close size="720" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增限时活动' : '编辑限时活动' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" :rules="rules" label-width="100px">
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入活动名称" clearable />
        </el-form-item>
        <el-form-item label="活动地点" prop="address">
          <el-input v-model="formData.address" placeholder="请输入活动地点" clearable />
        </el-form-item>
        <el-form-item label="活动时间" prop="timeRange">
          <el-date-picker
            v-model="formData.timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            style="width: 100%"
          />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="实际价格" prop="price">
              <el-input-number v-model="formData.price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="市场价" prop="marketPrice">
              <el-input-number v-model="formData.marketPrice" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="总名额(人次)" prop="quota">
          <el-input-number v-model="formData.quota" :min="0" style="width: 100%" />
          <div v-if="type === 'update'" class="text-gray-500 text-xs mt-1">
            已占用 {{ formData.sold ?? 0 }}，名额不可小于已占用
          </div>
        </el-form-item>
        <el-form-item label="封面" prop="coverImage">
          <SelectImage v-model="formData.coverImage" :multiple="false" />
        </el-form-item>
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="群二维码" prop="groupQr">
              <SelectImage v-model="formData.groupQr" :multiple="false" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="客服二维码" prop="serviceQr">
              <SelectImage v-model="formData.serviceQr" :multiple="false" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="活动详情" prop="detail">
          <el-input v-model="formData.detail" type="textarea" :rows="5" placeholder="活动详情" />
        </el-form-item>
        <el-form-item label="显示状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">显示</el-radio>
            <el-radio :value="0">隐藏</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createActivity,
  deleteActivity,
  deleteActivityByIds,
  updateActivity,
  findActivity,
  getActivityList
} from '@/plugin/limitedActivity/api/activity'
import { getUrl } from '@/utils/image'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref } from 'vue'
import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({ name: 'LimitedActivityManage' })

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])
const dialogVisible = ref(false)
const type = ref('create')
const formRef = ref()

const formData = ref({
  name: '',
  address: '',
  detail: '',
  timeRange: null,
  marketPrice: 0,
  price: 0,
  quota: 0,
  sold: 0,
  coverImage: '',
  groupQr: '',
  serviceQr: '',
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  timeRange: [{ required: true, message: '请选择活动时间', trigger: 'change' }],
  price: [{ required: true, message: '请输入实际价格', trigger: 'blur' }],
  quota: [{ required: true, message: '请输入名额', trigger: 'blur' }]
}

function formatDate(d) {
  if (!d) return ''
  const t = typeof d === 'string' ? d : (d && d.toISOString ? d.toISOString() : '')
  return t ? t.slice(0, 19).replace('T', ' ') : ''
}

const getTableData = async () => {
  const res = await getActivityList({
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
const handleSelectionChange = (val) => { multipleSelection.value = val }

const openDialog = () => {
  type.value = 'create'
  formData.value = {
    name: '',
    address: '',
    detail: '',
    timeRange: null,
    marketPrice: 0,
    price: 0,
    quota: 0,
    sold: 0,
    coverImage: '',
    groupQr: '',
    serviceQr: '',
    status: 1
  }
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}

const updateFunc = async (row) => {
  const res = await findActivity({ id: row.ID })
  if (res.code !== 0 || !res.data) return
  const d = res.data
  type.value = 'update'
  formData.value = {
    ID: d.ID,
    name: d.name || '',
    address: d.address || '',
    detail: d.detail || '',
    timeRange: d.startTime && d.endTime ? [d.startTime, d.endTime] : null,
    marketPrice: d.marketPrice ?? 0,
    price: d.price ?? 0,
    quota: d.quota ?? 0,
    sold: d.sold ?? 0,
    coverImage: d.coverImage || '',
    groupQr: d.groupQr || '',
    serviceQr: d.serviceQr || '',
    status: d.status ?? 1
  }
  dialogVisible.value = true
}

const enterDialog = async () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) return
    const range = formData.value.timeRange || []
    const payload = {
      ID: formData.value.ID,
      name: formData.value.name,
      address: formData.value.address,
      detail: formData.value.detail,
      startTime: range[0],
      endTime: range[1],
      marketPrice: formData.value.marketPrice,
      price: formData.value.price,
      quota: formData.value.quota,
      coverImage: formData.value.coverImage,
      groupQr: formData.value.groupQr,
      serviceQr: formData.value.serviceQr,
      status: formData.value.status
    }
    let res
    if (type.value === 'create') {
      res = await createActivity(payload)
    } else {
      res = await updateActivity(payload)
    }
    if (res.code === 0) {
      ElMessage.success(res.msg || '操作成功')
      closeDialog()
      getTableData()
    }
  })
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除该活动吗？', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteActivity({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  }).catch(() => {})
}

const onDelete = () => {
  ElMessageBox.confirm('确定删除选中活动吗？', '提示', { type: 'warning' }).then(async () => {
    const ids = multipleSelection.value.map((i) => i.ID)
    const res = await deleteActivityByIds(ids)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  }).catch(() => {})
}

getTableData()
</script>
