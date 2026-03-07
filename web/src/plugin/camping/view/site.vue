<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="场地名称">
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
        <el-button type="primary" icon="Plus" @click="openDialog">新增场地</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="tableRef" style="width: 100%" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="场地名称" prop="name" min-width="120" />
        <el-table-column align="left" label="开放时间" prop="openTimeDesc" show-overflow-tooltip />
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
          <span class="text-lg">{{ type === 'create' ? '新增场地' : '编辑场地' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" :rules="rules" label-width="100px">
        <el-form-item label="场地名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入场地名称" clearable />
        </el-form-item>
        <el-form-item label="轮播图" prop="carouselImages">
          <SelectImage v-model="formData.carouselImages" :multiple="true" :max-update-count="10" />
        </el-form-item>
        <el-form-item label="场地介绍（富文本）" prop="introduction">
          <RichEdit v-model="formData.introduction" />
        </el-form-item>
        <el-form-item label="预约规则介绍（富文本）" prop="reserveRules">
          <RichEdit v-model="formData.reserveRules" />
        </el-form-item>
        <el-form-item label="开放时间说明" prop="openTimeDesc">
          <el-input v-model="formData.openTimeDesc" type="textarea" :rows="3" placeholder="如：周一至周日 8:00-18:00" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { createSite, deleteSite, deleteSiteByIds, updateSite, findSite, getSiteList } from '@/plugin/camping/api/site'
import RichEdit from '@/components/richtext/rich-edit.vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'

defineOptions({ name: 'CampingSite' })

const formRef = ref()
const tableRef = ref()
const dialogVisible = ref(false)
const type = ref('create')
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const multipleSelection = ref([])
const searchInfo = ref({})

const formData = ref({
  name: '',
  carouselImages: [],
  introduction: '',
  reserveRules: '',
  openTimeDesc: '',
  status: 1
})

const rules = reactive({
  name: [{ required: true, message: '请输入场地名称', trigger: 'blur' }]
})

function toCarouselForSubmit(v) {
  if (!v || !Array.isArray(v)) return []
  return v.map((i) => (typeof i === 'string' ? i : i?.url)).filter(Boolean)
}

const getTableData = async () => {
  const res = await getSiteList({
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

const onSubmit = () => {
  page.value = 1
  getTableData()
}
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const openDialog = () => {
  type.value = 'create'
  formData.value = {
    name: '',
    carouselImages: [],
    introduction: '',
    reserveRules: '',
    openTimeDesc: '',
    status: 1
  }
  dialogVisible.value = true
}

const updateFunc = async (row) => {
  const res = await findSite({ id: row.ID })
  if (res.code === 0) {
    const d = res.data
    formData.value = {
      ID: d.ID,
      name: d.name || '',
      carouselImages: Array.isArray(d.carouselImages) ? d.carouselImages.map((u) => (typeof u === 'string' ? { url: u } : u)) : [],
      introduction: d.introduction || '',
      reserveRules: d.reserveRules || '',
      openTimeDesc: d.openTimeDesc || '',
      status: d.status ?? 1
    }
    type.value = 'update'
    dialogVisible.value = true
  }
}

const closeDialog = () => {
  dialogVisible.value = false
}

const enterDialog = async () => {
  await formRef.value?.validate(async (valid) => {
    if (!valid) return
    const payload = {
      ...formData.value,
      carouselImages: toCarouselForSubmit(formData.value.carouselImages)
    }
    let res
    if (type.value === 'create') {
      res = await createSite(payload)
    } else {
      res = await updateSite(payload)
    }
    if (res.code === 0) {
      ElMessage.success('操作成功')
      closeDialog()
      getTableData()
    }
  })
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除该场地？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteSite({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}

const onDelete = async () => {
  if (!multipleSelection.value.length) {
    ElMessage.warning('请选择要删除的数据')
    return
  }
  ElMessageBox.confirm('确定删除所选场地？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = multipleSelection.value.map((r) => r.ID)
    const res = await deleteSiteByIds(ids)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}

getTableData()
</script>
