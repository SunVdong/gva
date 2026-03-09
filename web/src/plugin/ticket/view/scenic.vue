<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="景区名称">
          <el-input v-model="searchInfo.name" placeholder="名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="全部" clearable style="width: 100px">
            <el-option label="启用" :value="1" />
            <el-option label="禁用" :value="0" />
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
        <el-button type="primary" icon="Plus" @click="openDialog">新增景区</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="tableRef" style="width: 100%" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="景区名称" prop="name" min-width="120" />
        <el-table-column align="left" label="轮播图" width="80">
          <template #default="{ row }">
            <el-image v-if="firstCarouselUrl(row)" :src="firstCarouselUrl(row)" style="width:48px;height:48px" fit="cover" />
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="可退订(小时)" prop="refundChangeHours" width="110" />
        <el-table-column align="left" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
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
    <el-drawer v-model="dialogVisible" destroy-on-close size="800" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增景区' : '编辑景区' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" :rules="rules" label-width="100px">
        <el-form-item label="景区名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入景区名称" clearable />
        </el-form-item>
        <el-form-item label="轮播图" prop="carouselImages">
          <SelectImage v-model="formData.carouselImages" :multiple="true" :max-update-count="10" />
        </el-form-item>
        <el-form-item label="可退订时间" prop="refundChangeHours">
          <el-input-number v-model="formData.refundChangeHours" :min="0" placeholder="游玩前多少小时可退订，0表示不可退订" style="width: 100%" />
          <div class="text-gray-500 text-xs mt-1">单位：小时。游玩日期前 N 小时内可退订，0 表示不可退订</div>
        </el-form-item>
        <el-form-item label="景区介绍" prop="description">
          <RichEdit v-model="formData.description" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="type === 'update' && formData.ID">
          <el-divider content-position="left">开放时间（按星期）</el-divider>
          <div class="flex justify-between items-center mb-2">
            <span class="text-gray-500 text-sm">设置每周开放时段，保存景区后生效</span>
            <el-button type="primary" size="small" @click="addOpenTimeRow">添加</el-button>
          </div>
          <el-table :data="openTimeList" border size="small">
            <el-table-column label="星期" width="120">
              <template #default="{ row }">
                <el-select v-model="row.weekDay" placeholder="星期" size="small" style="width:100%">
                  <el-option label="周一" :value="1" />
                  <el-option label="周二" :value="2" />
                  <el-option label="周三" :value="3" />
                  <el-option label="周四" :value="4" />
                  <el-option label="周五" :value="5" />
                  <el-option label="周六" :value="6" />
                  <el-option label="周日" :value="7" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="开放时间">
              <template #default="{ row }">
                <el-time-select v-model="row.openTime" start="00:00" step="00:30" end="23:30" placeholder="开始" size="small" style="width:48%" />
                <span class="px-1">-</span>
                <el-time-select v-model="row.closeTime" start="00:00" step="00:30" end="23:30" placeholder="结束" size="small" style="width:48%" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80">
              <template #default="{ $index }">
                <el-button type="danger" link size="small" @click="openTimeList.splice($index, 1)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-button v-if="openTimeList.length" type="primary" class="mt-2" @click="saveOpenTime">保存开放时间</el-button>
        </template>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { createScenic, deleteScenic, deleteScenicByIds, updateScenic, findScenic, getScenicList, getScenicOpenTimeByScenic, saveScenicOpenTime } from '@/plugin/ticket/api/scenic'
import RichEdit from '@/components/richtext/rich-edit.vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({ name: 'TicketScenic' })

const formRef = ref()
const dialogVisible = ref(false)
const type = ref('create')
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const multipleSelection = ref([])
const searchInfo = ref({})
const openTimeList = ref([])

const formData = ref({
  name: '',
  carouselImages: [],
  description: '',
  refundChangeHours: 0,
  status: 1
})

const rules = reactive({
  name: [{ required: true, message: '请输入景区名称', trigger: 'blur' }]
})

function firstCarouselUrl(row) {
  const v = row.carouselImages
  if (!v || !Array.isArray(v)) return ''
  const first = v[0]
  return typeof first === 'string' ? first : (first?.url || '')
}

function toCarouselForSubmit(v) {
  if (!v || !Array.isArray(v)) return []
  return v.map((i) => (typeof i === 'string' ? i : i?.url)).filter(Boolean)
}

const getTableData = async () => {
  const res = await getScenicList({
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
  formData.value = { name: '', carouselImages: [], description: '', refundChangeHours: 0, status: 1 }
  openTimeList.value = []
  dialogVisible.value = true
}

const updateFunc = async (row) => {
  const res = await findScenic({ id: row.ID })
  if (res.code === 0) {
    const d = res.data
    formData.value = {
      ID: d.ID,
      name: d.name || '',
      carouselImages: Array.isArray(d.carouselImages)
        ? d.carouselImages.map((u) => (typeof u === 'string' ? u : u?.url)).filter(Boolean)
        : [],
      description: d.description || '',
      refundChangeHours: d.refundChangeHours ?? 0,
      status: d.status ?? 1
    }
    type.value = 'update'
    dialogVisible.value = true
    const openRes = await getScenicOpenTimeByScenic({ scenicId: d.ID })
    if (openRes.code === 0 && openRes.data && openRes.data.length) {
      openTimeList.value = openRes.data.map((x) => ({
        weekDay: x.weekDay,
        openTime: (x.openTime || '').slice(0, 5) || '09:00',
        closeTime: (x.closeTime || '').slice(0, 5) || '18:00'
      }))
    } else {
      openTimeList.value = []
    }
  }
}

function addOpenTimeRow() {
  const used = openTimeList.value.map((r) => r.weekDay)
  if (used.length >= 7) { ElMessage.warning('周一至周日已全部添加'); return }
  let next = 1
  for (let d = 1; d <= 7; d++) { if (!used.includes(d)) { next = d; break } }
  openTimeList.value.push({ weekDay: next, openTime: '09:00', closeTime: '18:00' })
}

const saveOpenTime = async () => {
  const list = openTimeList.value
  const weekDays = list.map((x) => x.weekDay)
  const seen = new Set()
  const dup = weekDays.find((d) => { if (seen.has(d)) return true; seen.add(d); return false })
  if (dup !== undefined) {
    const dayNames = ['', '周一', '周二', '周三', '周四', '周五', '周六', '周日']
    ElMessage.warning(`星期不能重复，请修改「${dayNames[dup]}」`)
    return
  }
  const payload = {
    scenicId: formData.value.ID,
    list: list.map((x) => ({ scenicId: formData.value.ID, weekDay: x.weekDay, openTime: x.openTime, closeTime: x.closeTime }))
  }
  const res = await saveScenicOpenTime(payload)
  if (res.code === 0) ElMessage.success('开放时间已保存')
  else ElMessage.warning(res.msg || '保存失败')
}

const closeDialog = () => { dialogVisible.value = false }

const enterDialog = async () => {
  await formRef.value?.validate(async (valid) => {
    if (!valid) return
    const payload = { ...formData.value, carouselImages: toCarouselForSubmit(formData.value.carouselImages) }
    const res = type.value === 'create' ? await createScenic(payload) : await updateScenic(payload)
    if (res.code === 0) { ElMessage.success('操作成功'); closeDialog(); getTableData() }
  })
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除该景区？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }).then(async () => {
    const res = await deleteScenic({ id: row.ID })
    if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
  })
}

const onDelete = () => {
  if (!multipleSelection.value.length) { ElMessage.warning('请选择要删除的数据'); return }
  ElMessageBox.confirm('确定删除所选景区？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }).then(async () => {
    const ids = multipleSelection.value.map((r) => r.ID)
    const res = await deleteScenicByIds(ids)
    if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
  })
}

getTableData()
</script>
