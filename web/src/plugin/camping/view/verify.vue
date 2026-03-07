<template>
  <div class="p-4">
    <el-card header="核销" class="max-w-md">
      <el-form label-position="top">
        <el-form-item label="核销码">
          <el-input
            v-model="verifyCode"
            placeholder="请输入核销码或扫描二维码"
            clearable
            size="large"
            @keyup.enter="doVerify"
          >
            <template #append>
              <el-button type="primary" :loading="loading" @click="doVerify">核销</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>
      <el-alert v-if="lastResult" :title="lastResult" :type="lastSuccess ? 'success' : 'error'" show-icon class="mt-2" />
    </el-card>
  </div>
</template>

<script setup>
import { verifyReservationByCode } from '@/plugin/camping/api/reservation'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'

defineOptions({ name: 'CampingVerify' })

const verifyCode = ref('')
const loading = ref(false)
const lastResult = ref('')
const lastSuccess = ref(false)

const doVerify = async () => {
  const code = verifyCode.value?.trim()
  if (!code) {
    ElMessage.warning('请输入核销码')
    return
  }
  loading.value = true
  lastResult.value = ''
  try {
    const res = await verifyReservationByCode({ code })
    if (res.code === 0) {
      lastSuccess.value = true
      lastResult.value = '核销成功'
      verifyCode.value = ''
    } else {
      lastSuccess.value = false
      lastResult.value = res.msg || '核销失败'
    }
  } catch (e) {
    lastSuccess.value = false
    lastResult.value = '请求失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>
