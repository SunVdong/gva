<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="昵称">
          <el-input v-model="searchInfo.nickname" placeholder="昵称" clearable />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="searchInfo.phone" placeholder="手机号" clearable />
        </el-form-item>
        <el-form-item label="OpenID">
          <el-input v-model="searchInfo.openid" placeholder="微信 openid" clearable />
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
        <el-table-column align="left" label="头像" width="70">
          <template #default="{ row }">
            <el-avatar v-if="row.avatarUrl" :src="row.avatarUrl" :size="40" />
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="昵称" prop="nickname" min-width="100" show-overflow-tooltip />
        <el-table-column align="left" label="手机号" prop="phone" width="120" />
        <el-table-column align="left" label="OpenID" prop="openid" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="创建时间" width="170">
          <template #default="{ row }">
            {{ row.CreatedAt ? formatDate(row.CreatedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" width="100">
          <template #default="{ row }">
            <el-button type="primary" link icon="Edit" @click="updateFunc(row)">编辑</el-button>
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
    <el-drawer v-model="dialogVisible" destroy-on-close size="500" title="编辑用户" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">编辑用户</span>
          <div>
            <el-button type="primary" @click="enterDialog">保存</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" label-width="80px">
        <el-form-item label="昵称">
          <el-input v-model="formData.nickname" placeholder="昵称" clearable />
        </el-form-item>
        <el-form-item label="头像">
          <el-input v-model="formData.avatar" placeholder="头像 URL" clearable />
          <el-avatar v-if="formData.avatar" :src="formData.avatar" :size="60" class="mt-2" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="formData.phone" placeholder="手机号" clearable />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getUserList, findUser, updateUser } from '@/plugin/ticket/api/user'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

defineOptions({ name: 'TicketUser' })

const formRef = ref()
const dialogVisible = ref(false)
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const formData = ref({
  id: 0,
  nickname: '',
  avatar: '',
  phone: ''
})

function formatDate(d) {
  if (!d) return ''
  const t = typeof d === 'string' ? d : (d && d.toISOString ? d.toISOString() : '')
  return t ? t.slice(0, 19).replace('T', ' ') : ''
}

const getTableData = async () => {
  const res = await getUserList({
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

const updateFunc = async (row) => {
  const res = await findUser({ id: row.ID })
  if (res.code === 0) {
    const d = res.data
    formData.value = {
      id: d.ID,
      nickname: d.nickname || '',
      avatar: d.avatarUrl || '',
      phone: d.phone || ''
    }
    dialogVisible.value = true
  }
}

const closeDialog = () => { dialogVisible.value = false }

const enterDialog = async () => {
  const res = await updateUser(formData.value)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    closeDialog()
    getTableData()
  }
}

getTableData()
</script>
