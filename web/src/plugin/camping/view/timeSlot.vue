<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="场地">
          <el-select v-model="searchInfo.venueId" placeholder="全部" clearable style="width: 160px">
            <el-option v-for="v in venueOptions" :key="v.ID" :label="v.name" :value="v.ID" />
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
        <el-button type="primary" icon="Plus" @click="openDialog">新增时段</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="tableRef" style="width: 100%" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="场地" min-width="100">
          <template #default="{ row }">{{ venueName(row.venueId) }}</template>
        </el-table-column>
        <el-table-column align="left" label="开始时间" prop="startTime" width="100" />
        <el-table-column align="left" label="结束时间" prop="endTime" width="100" />
        <el-table-column align="left" label="可预约数" prop="capacity" width="100" />
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
    <el-drawer v-model="dialogVisible" destroy-on-close size="500" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增时段' : '编辑时段' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" :rules="rules" label-width="100px">
        <el-form-item label="场地" prop="venueId">
          <el-select v-model="formData.venueId" placeholder="请选择场地" style="width: 100%" filterable :disabled="type === 'update'">
            <el-option v-for="v in venueOptions" :key="v.ID" :label="v.name" :value="v.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间" prop="startTime">
          <el-time-select v-model="formData.startTime" start="00:00" step="00:30" end="23:30" placeholder="开始时间" style="width: 100%" />
        </el-form-item>
        <el-form-item label="结束时间" prop="endTime">
          <el-time-select v-model="formData.endTime" start="00:00" step="00:30" end="23:30" placeholder="结束时间" style="width: 100%" />
        </el-form-item>
        <el-form-item label="可预约数量" prop="capacity">
          <el-input-number v-model="formData.capacity" :min="1" style="width: 100%" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getSiteList } from '@/plugin/camping/api/site'
import { createTimeSlot, deleteTimeSlot, deleteTimeSlotByIds, updateTimeSlot, findTimeSlot, getTimeSlotList } from '@/plugin/camping/api/timeSlot'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({ name: 'CampingTimeSlot' })

const formRef = ref()
const dialogVisible = ref(false)
const type = ref('create')
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const multipleSelection = ref([])
const searchInfo = ref({})
const venueOptions = ref([])

const formData = ref({
  venueId: null,
  startTime: '',
  endTime: '',
  capacity: 1
})

const rules = reactive({
  venueId: [{ required: true, message: '请选择场地', trigger: 'change' }],
  startTime: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  endTime: [{ required: true, message: '请选择结束时间', trigger: 'change' }]
})

function venueName(id) {
  const v = venueOptions.value.find((x) => x.ID === id)
  return v ? v.name : id || '-'
}

const getTableData = async () => {
  const params = { page: page.value, pageSize: pageSize.value }
  if (searchInfo.value.venueId != null && searchInfo.value.venueId !== '') {
    params.venueId = searchInfo.value.venueId
  }
  const res = await getTimeSlotList(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = res.data.page || page.value
    pageSize.value = res.data.pageSize || pageSize.value
  }
}

const handleCurrentChange = (val) => { page.value = val; getTableData() }
const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
const handleSelectionChange = (val) => { multipleSelection.value = val }

const onSubmit = () => { page.value = 1; getTableData() }
const onReset = () => { searchInfo.value = {}; getTableData() }

const openDialog = () => {
  type.value = 'create'
  formData.value = { venueId: null, startTime: '', endTime: '', capacity: 1 }
  dialogVisible.value = true
}

const updateFunc = async (row) => {
  const res = await findTimeSlot({ id: row.ID })
  if (res.code === 0) {
    const d = res.data
    formData.value = {
      ID: d.ID,
      venueId: d.venueId,
      startTime: d.startTime?.slice(0, 5) || d.startTime || '',
      endTime: d.endTime?.slice(0, 5) || d.endTime || '',
      capacity: d.capacity ?? 1
    }
    type.value = 'update'
    dialogVisible.value = true
  }
}

const closeDialog = () => { dialogVisible.value = false }

const enterDialog = async () => {
  await formRef.value?.validate(async (valid) => {
    if (!valid) return
    const payload = {
      ...formData.value,
      startTime: formData.value.startTime?.length === 5 ? formData.value.startTime + ':00' : formData.value.startTime,
      endTime: formData.value.endTime?.length === 5 ? formData.value.endTime + ':00' : formData.value.endTime
    }
    let res
    if (type.value === 'create') res = await createTimeSlot(payload)
    else res = await updateTimeSlot(payload)
    if (res.code === 0) {
      ElMessage.success('操作成功')
      closeDialog()
      getTableData()
    }
  })
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除该时段？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const res = await deleteTimeSlot({ id: row.ID })
      if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
}

const onDelete = async () => {
  if (!multipleSelection.value.length) { ElMessage.warning('请选择要删除的数据'); return }
  ElMessageBox.confirm('确定删除所选时段？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
    .then(async () => {
      const ids = multipleSelection.value.map((r) => r.ID)
      const res = await deleteTimeSlotByIds(ids)
      if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
}

onMounted(async () => {
  const siteRes = await getSiteList({ page: 1, pageSize: 500 })
  if (siteRes.code === 0) venueOptions.value = siteRes.data.list || []
  getTableData()
})
</script>
