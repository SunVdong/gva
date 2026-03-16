<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
        <el-form-item label="景区">
          <el-select v-model="searchInfo.scenicId" placeholder="全部" clearable filterable style="width: 200px">
            <el-option v-for="s in scenicOptions" :key="s.ID" :label="s.name" :value="s.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名称">
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
        <el-button type="primary" icon="Plus" @click="openDialog">新增门票商品</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="ID" width="80" />
        <el-table-column align="left" label="商品名称" prop="name" min-width="140" />
        <el-table-column align="left" label="所属景区" min-width="120">
          <template #default="{ row }">
            {{ scenicName(row.scenicId) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">{{ row.status === 1 ? '启用' : '禁用' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="200">
          <template #default="{ row }">
            <el-button type="primary" link icon="Edit" @click="updateFunc(row)">编辑</el-button>
            <el-button type="primary" link @click="openSkuDrawer(row)">SKU/规则</el-button>
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

    <el-drawer v-model="dialogVisible" destroy-on-close size="700" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增门票商品' : '编辑门票商品' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确定</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>
      <el-form ref="formRef" :model="formData" label-position="top" :rules="rules" label-width="100px">
        <el-form-item label="所属景区" prop="scenicId">
          <el-select v-model="formData.scenicId" placeholder="请选择景区" filterable style="width:100%">
            <el-option v-for="s in scenicOptions" :key="s.ID" :label="s.name" :value="s.ID" />
          </el-select>
        </el-form-item>
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="formData.name" placeholder="如：成人门票" clearable />
        </el-form-item>
        <el-form-item label="门票说明" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="选填" />
        </el-form-item>
        <el-form-item label="适用人群" prop="audience">
          <el-input
            v-model="formData.audience"
            type="textarea"
            :rows="3"
            placeholder="如：成人、儿童、老人等，选填"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer v-model="skuDrawerVisible" destroy-on-close size="960" title="SKU 与规则">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">「{{ currentProductName }}」SKU 与规则</span>
          <div>
            <el-button type="primary" size="small" @click="addSkuRow">添加 SKU</el-button>
            <el-button type="primary" @click="saveSkuAndRule">保存</el-button>
            <el-button @click="skuDrawerVisible = false">关闭</el-button>
          </div>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="门票 SKU" name="sku">
          <el-table :data="skuList" border size="small">
            <el-table-column label="SKU 名称" min-width="130">
              <template #default="{ row }">
                <el-input v-model="row.name" placeholder="如成人票" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="销售价" width="118">
              <template #default="{ row }">
                <el-input-number v-model="row.price" :min="0" :precision="2" size="small" style="width:100%" />
              </template>
            </el-table-column>
            <el-table-column label="市场价" width="118">
              <template #default="{ row }">
                <el-input-number v-model="row.marketPrice" :min="0" :precision="2" size="small" style="width:100%" />
              </template>
            </el-table-column>
            <el-table-column label="库存" width="118">
              <template #default="{ row }">
                <el-input-number v-model="row.stock" :min="0" size="small" style="width:100%" />
              </template>
            </el-table-column>
            <el-table-column label="限购" width="96">
              <template #default="{ row }">
                <el-input-number v-model="row.limitBuy" :min="0" size="small" style="width:100%" />
              </template>
            </el-table-column>
            <el-table-column label="状态" width="96">
              <template #default="{ row }">
                <el-select v-model="row.status" size="small" style="width:100%">
                  <el-option label="启用" :value="1" />
                  <el-option label="禁用" :value="0" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="90" fixed="right">
              <template #default="{ row, $index }">
                <el-button v-if="!row.ID" type="danger" link size="small" @click="skuList.splice($index, 1)">删除</el-button>
                <el-button v-else type="danger" link size="small" @click="deleteSkuRow(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="门票规则" name="rule">
          <div class="flex justify-between items-center mb-2">
            <span class="text-gray-500 text-sm">购买须知、退改规则等</span>
            <el-button type="primary" size="small" @click="addRuleRow">添加规则</el-button>
          </div>
          <el-table :data="ruleList" border size="small">
            <el-table-column label="标题" min-width="100">
              <template #default="{ row }">
                <el-input v-model="row.title" placeholder="规则标题" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="内容" min-width="180">
              <template #default="{ row }">
                <el-input v-model="row.content" type="textarea" :rows="2" placeholder="规则内容" size="small" />
              </template>
            </el-table-column>
            <el-table-column label="排序" width="100">
              <template #default="{ row }">
                <el-input-number v-model="row.sort" :min="0" size="small" style="width:100%" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="80">
              <template #default="{ $index }">
                <el-button type="danger" link size="small" @click="ruleList.splice($index, 1)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-drawer>

    <el-dialog
      v-model="audienceDialogVisible"
      title="适用人群"
      width="500"
      destroy-on-close
      @open="onAudienceDialogOpen"
    >
      <template #header>
        <span>适用人群 — {{ currentAudienceSku.name || 'SKU' }}</span>
      </template>
      <div class="space-y-2">
        <div
          v-for="(item, idx) in audienceList"
          :key="idx"
          class="flex items-center gap-2"
        >
          <el-input v-model="item.audienceType" placeholder="如：成人、儿童" size="small" style="width:140px" />
          <el-input v-model="item.description" placeholder="说明（选填）" size="small" class="flex-1" />
          <el-button type="danger" link size="small" icon="Delete" @click="audienceList.splice(idx, 1)" />
        </div>
        <el-button type="primary" link size="small" icon="Plus" @click="addAudienceItem">添加一项</el-button>
      </div>
      <template #footer>
        <el-button @click="audienceDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="audienceSaving" @click="saveAudienceDialog">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { getScenicList } from '@/plugin/ticket/api/scenic'
import {
  createProduct,
  deleteProduct,
  deleteProductByIds,
  updateProduct,
  findProduct,
  getProductList,
  getSkuList,
  createSku,
  updateSku,
  deleteSku,
  getRuleByProduct,
  saveRule,
  getAudienceBySku,
  saveAudience
} from '@/plugin/ticket/api/product'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'

defineOptions({ name: 'TicketProduct' })

const formRef = ref()
const dialogVisible = ref(false)
const skuDrawerVisible = ref(false)
const type = ref('create')
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const multipleSelection = ref([])
const searchInfo = ref({})
const scenicOptions = ref([])
const currentProductId = ref(0)
const currentProductName = ref('')
const skuList = ref([])
const ruleList = ref([])
const activeTab = ref('sku')
const audienceDialogVisible = ref(false)
const currentAudienceSku = ref({ id: 0, name: '' })
const audienceList = ref([])
const audienceSaving = ref(false)

const formData = ref({
  scenicId: undefined,
  name: '',
  description: '',
  audience: '',
  status: 1
})

const rules = reactive({
  scenicId: [{ required: true, message: '请选择景区', trigger: 'change' }],
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }]
})

function scenicName(id) {
  if (!id) return '-'
  const s = scenicOptions.value.find((x) => x.ID === id)
  return s ? s.name : '-'
}

const getScenicOptions = async () => {
  const res = await getScenicList({ page: 1, pageSize: 500 })
  if (res.code === 0 && res.data.list) scenicOptions.value = res.data.list
}

const getTableData = async () => {
  const res = await getProductList({
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
  formData.value = { scenicId: undefined, name: '', description: '', audience: '', status: 1 }
  dialogVisible.value = true
}

const updateFunc = async (row) => {
  const res = await findProduct({ id: row.ID })
  if (res.code === 0) {
    const d = res.data
    formData.value = {
      ID: d.ID,
      scenicId: d.scenicId,
      name: d.name || '',
      description: d.description || '',
      audience: d.audience || '',
      status: d.status ?? 1
    }
    type.value = 'update'
    dialogVisible.value = true
  }
}

const enterDialog = async () => {
  await formRef.value?.validate(async (valid) => {
    if (!valid) return
    const res = type.value === 'create' ? await createProduct(formData.value) : await updateProduct(formData.value)
    if (res.code === 0) {
      ElMessage.success('操作成功')
      closeDialog()
      getTableData()
    }
  })
}

const closeDialog = () => { dialogVisible.value = false }

const openSkuDrawer = async (row) => {
  currentProductId.value = row.ID
  currentProductName.value = row.name
  skuDrawerVisible.value = true
  const [skuRes, ruleRes] = await Promise.all([
    getSkuList({ productId: row.ID, page: 1, pageSize: 100 }),
    getRuleByProduct({ productId: row.ID })
  ])
  skuList.value = (skuRes.code === 0 && skuRes.data.list) ? skuRes.data.list.map((s) => ({ ...s, marketPrice: s.marketPrice ?? undefined })) : []
  ruleList.value = (ruleRes.code === 0 && ruleRes.data) ? ruleRes.data.map((r) => ({ ...r })) : []
  activeTab.value = 'sku'
}

function addSkuRow() {
  skuList.value.push({
    productId: currentProductId.value,
    name: '',
    price: 0,
    marketPrice: undefined,
    stock: 0,
    limitBuy: 0,
    status: 1
  })
}

async function deleteSkuRow(row) {
  if (!row.ID) { skuList.value = skuList.value.filter((s) => s !== row); return }
  await ElMessageBox.confirm('确定删除该 SKU？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' })
  const res = await deleteSku({ id: row.ID })
  if (res.code === 0) {
    ElMessage.success('已删除')
    skuList.value = skuList.value.filter((s) => s.ID !== row.ID)
  }
}

function addAudienceItem() {
  audienceList.value.push({ audienceType: '', description: '' })
}

async function openAudienceDialog(row) {
  currentAudienceSku.value = { id: row.ID, name: row.name }
  audienceDialogVisible.value = true
}

async function onAudienceDialogOpen() {
  const res = await getAudienceBySku({ skuId: currentAudienceSku.value.id })
  if (res.code === 0 && Array.isArray(res.data)) {
    audienceList.value = res.data.map((x) => ({
      audienceType: x.audienceType || '',
      description: x.description || ''
    }))
  } else {
    audienceList.value = []
  }
  if (audienceList.value.length === 0) audienceList.value.push({ audienceType: '', description: '' })
}

async function saveAudienceDialog() {
  const list = audienceList.value
    .filter((x) => x.audienceType && x.audienceType.trim())
    .map((x) => ({ audienceType: x.audienceType.trim(), description: (x.description || '').trim() }))
  if (list.length === 0) {
    ElMessage.warning('请至少填写一项适用人群')
    return
  }
  audienceSaving.value = true
  try {
    const res = await saveAudience({ skuId: currentAudienceSku.value.id, list })
    if (res.code === 0) {
      ElMessage.success('保存成功')
      audienceDialogVisible.value = false
    } else {
      ElMessage.error(res.msg || '保存失败')
    }
  } finally {
    audienceSaving.value = false
  }
}

function addRuleRow() {
  ruleList.value.push({ title: '', content: '', sort: ruleList.value.length })
}

const saveSkuAndRule = async () => {
  for (const s of skuList.value) {
    if (!s.name || s.price == null) {
      ElMessage.warning('请填写 SKU 名称和销售价')
      return
    }
    const payload = {
      productId: currentProductId.value,
      name: s.name,
      price: Number(s.price),
      marketPrice: s.marketPrice != null ? Number(s.marketPrice) : null,
      stock: Number(s.stock) || 0,
      limitBuy: Number(s.limitBuy) || 0,
      status: s.status ?? 1
    }
    if (s.ID) {
      await updateSku({ ...payload, ID: s.ID })
    } else {
      await createSku(payload)
    }
  }
  const rulePayload = {
    productId: currentProductId.value,
    list: ruleList.value.map((r) => ({ title: r.title, content: r.content, sort: r.sort }))
  }
  const res = await saveRule(rulePayload)
  if (res.code === 0) {
    ElMessage.success('保存成功')
    openSkuDrawer({ ID: currentProductId.value, name: currentProductName.value })
  }
}

const deleteRow = (row) => {
  ElMessageBox.confirm('确定删除该门票商品？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }).then(async () => {
    const res = await deleteProduct({ id: row.ID })
    if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
  })
}

const onDelete = () => {
  if (!multipleSelection.value.length) { ElMessage.warning('请选择要删除的数据'); return }
  ElMessageBox.confirm('确定删除所选门票商品？', '提示', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }).then(async () => {
    const ids = multipleSelection.value.map((r) => r.ID)
    const res = await deleteProductByIds(ids)
    if (res.code === 0) { ElMessage.success('删除成功'); getTableData() }
  })
}

onMounted(() => {
  getScenicOptions()
  getTableData()
})
</script>
