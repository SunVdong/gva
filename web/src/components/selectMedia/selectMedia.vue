<template>
  <div>
    <el-upload
      v-model:file-list="fileList"
      multiple
      :http-request="customUpload"
      :on-remove="handleRemove"
      :before-upload="beforeUpload"
      list-type="picture-card"
      :limit="limit"
      :accept="effectiveAccept"
      class="select-media-upload"
    >
      <el-icon><Plus /></el-icon>

      <template #file="{ file }">
        <div class="file-card">
          <template v-if="isVideo(file)">
            <video
              :src="file.url"
              class="thumb"
              muted
              preload="metadata"
            />
            <div class="video-badge">
              <el-icon :size="16"><VideoPlay /></el-icon>
            </div>
          </template>
          <el-image
            v-else
            :src="file.url"
            fit="cover"
            class="thumb"
          />
          <span class="actions">
            <span class="action-btn" @click.stop="handlePreview(file)">
              <el-icon><ZoomIn /></el-icon>
            </span>
            <span class="action-btn" @click.stop="handleCardRemove(file)">
              <el-icon><Delete /></el-icon>
            </span>
          </span>
        </div>
      </template>
    </el-upload>

    <el-image-viewer
      v-if="previewVisible && !previewIsVideo"
      :url-list="[previewUrl]"
      @close="previewVisible = false"
      teleported
    />

    <el-dialog
      v-model="videoDialogVisible"
      title="视频预览"
      width="720px"
      destroy-on-close
      append-to-body
    >
      <video
        :src="previewUrl"
        controls
        autoplay
        class="preview-video"
      />
    </el-dialog>
  </div>
</template>

<script setup>
  import { ref, watch, computed } from 'vue'
  import { ElMessage } from 'element-plus'
  import { Plus, ZoomIn, Delete, VideoPlay } from '@element-plus/icons-vue'
  import { getUrl, isImageMime, isVideoExt } from '@/utils/image'
  import service from '@/utils/request'

  defineOptions({
    name: 'SelectMedia'
  })

  const props = defineProps({
    limit: {
      type: Number,
      default: 20
    },
    accept: {
      type: String,
      default: 'image/*,video/*'
    },
    /** 为 true 时允许图片；视频仅允许 MP4（活动指南等场景） */
    restrictVideoToMp4: {
      type: Boolean,
      default: false
    },
    /** 与 restrictVideoToMp4 配合：单个视频最大体积（MB），0 表示不限制 */
    maxVideoSizeMB: {
      type: Number,
      default: 50
    }
  })

  const imageFilenameExt = /\.(jpe?g|png|gif|webp|svg|bmp)$/i

  const effectiveAccept = computed(() =>
    props.restrictVideoToMp4 ? 'image/*,video/mp4,.mp4' : props.accept
  )

  const beforeUpload = (rawFile) => {
    if (!props.restrictVideoToMp4) {
      return true
    }
    const name = (rawFile.name || '').toLowerCase()
    const type = (rawFile.type || '').toLowerCase()

    if (type.startsWith('image/') || isImageMime(type)) {
      return true
    }
    if (imageFilenameExt.test(name)) {
      return true
    }

    let asVideo = false
    if (name.endsWith('.mp4')) {
      if (
        !type ||
        type === 'video/mp4' ||
        type === 'application/octet-stream'
      ) {
        asVideo = true
      } else if (type.startsWith('video/') && type !== 'video/mp4') {
        ElMessage.error('视频仅支持 MP4 格式')
        return false
      } else {
        asVideo = true
      }
    } else if (type === 'video/mp4') {
      asVideo = true
    } else if (type.startsWith('video/')) {
      ElMessage.error('视频仅支持 MP4 格式')
      return false
    } else {
      ElMessage.error('请上传图片或 MP4 视频')
      return false
    }

    if (
      asVideo &&
      props.maxVideoSizeMB > 0 &&
      rawFile.size > props.maxVideoSizeMB * 1024 * 1024
    ) {
      ElMessage.error(`视频大小不能超过 ${props.maxVideoSizeMB}MB`)
      return false
    }
    return true
  }

  const model = defineModel({ type: Array })
  const emits = defineEmits(['on-success', 'on-error'])

  const buildFileList = (items) =>
    (items || []).map(item => ({
      name: item.name,
      url: getUrl(item.url),
      _rawUrl: item.url,
    }))

  const fileList = ref(buildFileList(model.value))

  watch(() => model.value, (val) => {
    if (!val || !val.length) {
      fileList.value = []
    }
  })

  const isVideo = (file) => isVideoExt(file.url) || isVideoExt(file.name)

  const previewVisible = ref(false)
  const videoDialogVisible = ref(false)
  const previewUrl = ref('')
  const previewIsVideo = ref(false)

  const handlePreview = (file) => {
    const url = file.url || getUrl(file._rawUrl)
    previewUrl.value = url
    previewIsVideo.value = isVideo(file)
    if (previewIsVideo.value) {
      videoDialogVisible.value = true
    } else {
      previewVisible.value = true
    }
  }

  const customUpload = async (options) => {
    const { file, onSuccess, onError, onProgress } = options
    const formData = new FormData()
    formData.append('file', file)

    try {
      const res = await service({
        url: '/fileUploadAndDownload/upload?noSave=1',
        method: 'post',
        data: formData,
        headers: { 'Content-Type': 'multipart/form-data' },
        donNotShowLoading: true,
        onUploadProgress: (e) => {
          if (e.total > 0) {
            onProgress({ percent: Math.round((e.loaded / e.total) * 100) })
          }
        }
      })

      if (res.code === 0) {
        const serverUrl = getUrl(res.data.file.url)
        onSuccess(res)
        const uploaded = fileList.value.find(f => f.uid === options.file.uid)
        if (uploaded) {
          uploaded.url = serverUrl
          uploaded._rawUrl = res.data.file.url
        }
        model.value.push({
          name: res.data.file.name,
          url: res.data.file.url
        })
        emits('on-success', res)
      } else {
        onError(new Error(res.msg || '上传失败'))
        ElMessage({ type: 'error', message: '上传失败: ' + (res.msg || '') })
        fileList.value.pop()
      }
    } catch (err) {
      onError(err)
      ElMessage({ type: 'error', message: '上传失败' })
      fileList.value.pop()
      emits('on-error', err)
    }
  }

  const handleRemove = (file) => {
    const rawUrl = file._rawUrl || file.url
    const idx = model.value.findIndex(
      item => item.url === rawUrl || getUrl(item.url) === file.url || item.name === file.name
    )
    if (idx > -1) {
      model.value.splice(idx, 1)
    }
  }

  const handleCardRemove = (file) => {
    handleRemove(file)
    const idx = fileList.value.findIndex(f => f.uid === file.uid || f.url === file.url)
    if (idx > -1) {
      fileList.value.splice(idx, 1)
    }
  }
</script>

<style scoped>
.select-media-upload :deep(.el-upload-list__item) {
  overflow: hidden;
  border-radius: 6px;
}

.file-card {
  position: relative;
  width: 100%;
  height: 100%;
}

.file-card .thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.file-card .video-badge {
  position: absolute;
  left: 6px;
  bottom: 6px;
  background: rgba(0, 0, 0, 0.55);
  color: #fff;
  border-radius: 4px;
  padding: 2px 6px;
  display: flex;
  align-items: center;
}

.file-card .actions {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: rgba(0, 0, 0, 0.45);
  opacity: 0;
  transition: opacity 0.2s;
}

.file-card:hover .actions {
  opacity: 1;
}

.action-btn {
  color: #fff;
  font-size: 18px;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: transform 0.15s;
}

.action-btn:hover {
  transform: scale(1.2);
}

.preview-video {
  width: 100%;
  max-height: 480px;
  outline: none;
}
</style>
