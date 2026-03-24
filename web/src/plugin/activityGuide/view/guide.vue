<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="ID">
          <el-input
            v-model.number="searchInfo.ID"
            placeholder="ID"
            clearable
            style="width: 120px"
          />
        </el-form-item>
        <el-form-item label="活动名称">
          <el-input
            v-model="searchInfo.name"
            placeholder="请输入活动名称"
            clearable
          />
        </el-form-item>
        <el-form-item label="显示状态">
          <el-select
            v-model="searchInfo.showStatus"
            placeholder="全部"
            clearable
            style="width: 100px"
          >
            <el-option label="显示" :value="true" />
            <el-option label="隐藏" :value="false" />
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
        <el-button type="primary" icon="Plus" @click="openDialog">新增</el-button>
        <el-button
          icon="Delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >
          删除
        </el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="创建时间" prop="CreatedAt" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.CreatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="封面" width="80">
          <template #default="scope">
            <el-image
              v-if="scope.row.coverImage"
              :src="getUrl(scope.row.coverImage)"
              fit="cover"
              class="w-12 h-12 rounded border"
              :preview-src-list="[getUrl(scope.row.coverImage)]"
              :z-index="9999"
              :preview-teleported="true"
            />
            <span v-else class="text-gray-400">—</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="活动名称" prop="name" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="简介" prop="summary" min-width="200" show-overflow-tooltip />
        <el-table-column align="left" label="显示状态" width="100">
          <template #default="scope">
            <el-switch
              v-model="scope.row.showStatus"
              :active-value="true"
              :inactive-value="false"
              @change="(val) => toggleShowStatus(scope.row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="160">
          <template #default="scope">
            <el-button type="primary" link icon="Edit" @click="updateGuideFunc(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="Delete" @click="deleteRow(scope.row)">删除</el-button>
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

    <el-drawer
      v-model="dialogFormVisible"
      destroy-on-close
      size="720"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增活动指南' : '编辑活动指南' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>

      <el-form ref="elFormRef" :model="formData" label-position="top" :rules="rule" label-width="100px">
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入活动名称" clearable maxlength="128" show-word-limit />
        </el-form-item>
        <el-form-item label="简介" prop="summary">
          <el-input
            v-model="formData.summary"
            type="textarea"
            :rows="3"
            placeholder="请输入简介"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="封面图" prop="coverImage">
          <SelectImage v-model="formData.coverImage" :multiple="false" />
        </el-form-item>
        <el-form-item label="介绍视频或图片" prop="media">
          <SelectMedia v-model="formData.media" />
          <div class="text-gray-500 text-xs mt-1">可上传多张图片或视频，用于活动介绍展示</div>
        </el-form-item>
        <el-form-item label="显示状态" prop="showStatus">
          <el-switch
            v-model="formData.showStatus"
            :active-value="true"
            :inactive-value="false"
            active-text="显示"
            inactive-text="隐藏"
          />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  createGuide,
  deleteGuide,
  deleteGuideByIds,
  updateGuide,
  findGuide,
  getGuideList
} from '@/plugin/activityGuide/api/guide'
import { getUrl } from '@/utils/image'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import SelectImage from '@/components/selectImage/selectImage.vue'
import SelectMedia from '@/components/selectMedia/selectMedia.vue'

defineOptions({ name: 'ActivityGuide' })

const rule = reactive({
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }]
})

const formData = ref({
  name: '',
  summary: '',
  coverImage: '',
  media: [],
  showStatus: true
})

const elFormRef = ref()
const elSearchFormRef = ref()

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const getTableData = async () => {
  const params = {
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  }
  if (params.showStatus === undefined || params.showStatus === null) delete params.showStatus
  if (params.ID === undefined || params.ID === null || params.ID === '' || Number.isNaN(params.ID)) delete params.ID
  const res = await getGuideList(params)
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total || 0
    page.value = res.data.page || 1
    pageSize.value = res.data.pageSize || 10
  }
}

getTableData()

const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const toggleShowStatus = async (row, val) => {
  try {
    const res = await updateGuide({
      ...row,
      showStatus: val
    })
    if (res.code === 0) {
      ElMessage.success(val ? '已设为显示' : '已设为隐藏')
    } else {
      row.showStatus = !val
      ElMessage.error(res.msg || '操作失败')
    }
  } catch (e) {
    row.showStatus = !val
    ElMessage.error('操作失败')
  }
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除该活动指南吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => deleteGuideFunc(row))
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除所选记录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    if (!multipleSelection.value.length) {
      ElMessage.warning('请选择要删除的数据')
      return
    }
    const IDs = multipleSelection.value.map((item) => item.ID)
    const res = await deleteGuideByIds({ IDs })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === IDs.length && page.value > 1) page.value--
      getTableData()
    }
  })
}

const type = ref('')
const updateGuideFunc = async (row) => {
  const res = await findGuide({ ID: row.ID })
  type.value = 'update'
  if (res.code === 0) {
    const data = res.data
    formData.value = {
      ...data,
      media: Array.isArray(data.media) ? data.media : data.media ? JSON.parse(data.media) : []
    }
    dialogFormVisible.value = true
  }
}

const deleteGuideFunc = async (row) => {
  const res = await deleteGuide({ ID: row.ID })
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) page.value--
    getTableData()
  }
}

const dialogFormVisible = ref(false)

const openDialog = () => {
  type.value = 'create'
  formData.value = {
    name: '',
    summary: '',
    coverImage: '',
    media: [],
    showStatus: true
  }
  dialogFormVisible.value = true
}

const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    name: '',
    summary: '',
    coverImage: '',
    media: [],
    showStatus: true
  }
}

const enterDialog = async () => {
  await elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    const payload = { ...formData.value }
    const fn = type.value === 'create' ? createGuide : updateGuide
    const res = await fn(payload)
    if (res.code === 0) {
      ElMessage.success(type.value === 'create' ? '创建成功' : '更新成功')
      closeDialog()
      getTableData()
    }
  })
}
</script>

<style scoped>
.w-12 {
  width: 3rem;
}
.h-12 {
  height: 3rem;
}

:deep(.el-image-viewer__wrapper) {
  z-index: 9999 !important;
}
</style>
