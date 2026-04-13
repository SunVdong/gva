<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="内容关键词">
          <el-input v-model="searchInfo.content" placeholder="模糊搜索" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          icon="delete"
          style="margin-left: 10px;"
          :disabled="!multipleSelection.length"
          @click="onDeleteBatch"
        >删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="72" />
        <el-table-column align="left" label="用户" width="140">
          <template #default="scope">
            {{ scope.row.user?.nickname || ('#' + scope.row.userId) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="内容" prop="content" min-width="240" show-overflow-tooltip />
        <el-table-column align="left" label="提交时间" width="168">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="操作" width="160" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="openDetail(scope.row)">查看</el-button>
            <el-popover v-model:visible="scope.row.visible" placement="top" width="160">
              <p>确定要删除吗？</p>
              <div style="text-align: right; margin: 0">
                <el-button size="small" type="primary" link @click="scope.row.visible = false">取消</el-button>
                <el-button size="small" type="primary" @click="deleteRow(scope.row)">确定</el-button>
              </div>
              <template #reference>
                <el-button type="primary" link @click="scope.row.visible = true">删除</el-button>
              </template>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" title="反馈详情" width="560px" destroy-on-close>
      <template v-if="detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="反馈内容">
            <div class="content-pre">{{ detail.content }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="用户">
            {{ detail.user?.nickname || ('#' + detail.userId) }}
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">{{ formatDate(detail.CreatedAt) }}</el-descriptions-item>
        </el-descriptions>
      </template>
      <template #footer>
        <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  getFeedbackList,
  deleteFeedback,
  deleteFeedbackByIds,
  findFeedback
} from '@/api/feedback'
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const multipleSelection = ref([])
const dialogVisible = ref(false)
const detail = ref(null)

const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const getTableData = async () => {
  const res = await getFeedbackList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (res.code === 0) {
    tableData.value = res.data.list
    total.value = res.data.total
    page.value = res.data.page
    pageSize.value = res.data.pageSize
  }
}

const deleteRow = async (row) => {
  row.visible = false
  const res = await deleteFeedback(row)
  if (res.code === 0) {
    ElMessage.success('删除成功')
    if (tableData.value.length === 1 && page.value > 1) {
      page.value--
    }
    getTableData()
  }
}

const onDeleteBatch = () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = multipleSelection.value.map((item) => item.ID)
    const res = await deleteFeedbackByIds({ ids })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === ids.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  }).catch(() => {})
}

const openDetail = async (row) => {
  const res = await findFeedback({ ID: row.ID })
  if (res.code !== 0) {
    return
  }
  detail.value = res.data
  dialogVisible.value = true
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  page.value = 1
  pageSize.value = 10
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

getTableData()
</script>

<style scoped>
.content-pre {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
