<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="标题">
          <el-input v-model="searchInfo.title" placeholder="标题" clearable />
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
        <el-button type="primary" icon="Plus" @click="openDialog">新增 Banner</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="标题" prop="title" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="轮播图" width="90">
          <template #default="{ row }">
            <el-image
              v-if="row.image"
              :src="getUrl(row.image)"
              style="width:48px;height:48px"
              fit="cover"
              :preview-src-list="[getUrl(row.image)]"
            />
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="详情长图" width="90">
          <template #default="{ row }">
            <el-image
              v-if="row.detailImage"
              :src="getUrl(row.detailImage)"
              style="width:48px;height:48px"
              fit="cover"
              :preview-src-list="[getUrl(row.detailImage)]"
            />
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="排序" prop="sort" width="80" />
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

    <el-drawer v-model="dialogVisible" destroy-on-close size="560" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增 Banner' : '编辑 Banner' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" :rules="rules" label-width="100px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="formData.title" placeholder="后台识别用标题（小程序可不展示）" clearable />
        </el-form-item>
        <el-form-item label="轮播图" prop="image">
          <SelectImage v-model="formData.image" :multiple="false" />
        </el-form-item>
        <el-form-item label="详情长图" prop="detailImage">
          <SelectImage v-model="formData.detailImage" :multiple="false" />
        </el-form-item>
        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="0" style="width: 100%" />
          <div class="text-gray-500 text-xs mt-1">数值越小越靠前</div>
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
  createBanner,
  deleteBanner,
  deleteBannerByIds,
  updateBanner,
  findBanner,
  getBannerList
} from '@/plugin/limitedActivity/api/banner'
import { getUrl } from '@/utils/image'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref } from 'vue'
import SelectImage from '@/components/selectImage/selectImage.vue'

defineOptions({ name: 'LimitedActivityBanner' })

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
  title: '',
  image: '',
  detailImage: '',
  sort: 0,
  status: 1
})

const rules = {
  image: [{ required: true, message: '请选择轮播图', trigger: 'change' }],
  detailImage: [{ required: true, message: '请选择详情长图', trigger: 'change' }]
}

function formatDate(d) {
  if (!d) return ''
  const t = typeof d === 'string' ? d : (d && d.toISOString ? d.toISOString() : '')
  return t ? t.slice(0, 19).replace('T', ' ') : ''
}

const getTableData = async () => {
  const res = await getBannerList({
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
    title: '',
    image: '',
    detailImage: '',
    sort: 0,
    status: 1
  }
  dialogVisible.value = true
}

const closeDialog = () => {
  dialogVisible.value = false
}

const updateFunc = async (row) => {
  const res = await findBanner({ id: row.ID })
  if (res.code !== 0 || !res.data) return
  const d = res.data
  type.value = 'update'
  formData.value = {
    ID: d.ID,
    title: d.title || '',
    image: d.image || '',
    detailImage: d.detailImage || '',
    sort: d.sort ?? 0,
    status: d.status ?? 1
  }
  dialogVisible.value = true
}

const enterDialog = async () => {
  formRef.value?.validate(async (valid) => {
    if (!valid) return
    const payload = {
      ID: formData.value.ID,
      title: formData.value.title,
      image: formData.value.image,
      detailImage: formData.value.detailImage,
      sort: formData.value.sort,
      status: formData.value.status
    }
    let res
    if (type.value === 'create') {
      res = await createBanner(payload)
    } else {
      res = await updateBanner(payload)
    }
    if (res.code === 0) {
      ElMessage.success(res.msg || '操作成功')
      closeDialog()
      getTableData()
    }
  })
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除该 Banner 吗？', '提示', { type: 'warning' }).then(async () => {
    const res = await deleteBanner({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  }).catch(() => {})
}

const onDelete = () => {
  ElMessageBox.confirm('确定删除选中 Banner 吗？', '提示', { type: 'warning' }).then(async () => {
    const ids = multipleSelection.value.map((i) => i.ID)
    const res = await deleteBannerByIds(ids)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  }).catch(() => {})
}

getTableData()
</script>
