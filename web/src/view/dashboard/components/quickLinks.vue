<template>
  <div class="mt-4 w-full">
    <div class="text-xs tracking-wide text-black/60 dark:text-white/60">快捷入口</div>
    <div class="mt-3 grid grid-cols-3 gap-3 sm:grid-cols-4">
      <div
        v-for="(item, index) in shortcuts"
        :key="index"
        class="flex flex-col items-center group cursor-pointer"
        @click="toPath(item)"
      >
        <div
          class="w-10 h-10 rounded-lg border border-black/10 dark:border-white/10 flex items-center justify-center text-black/70 dark:text-white/70 group-hover:bg-[var(--el-color-primary)] group-hover:text-white transition-colors"
        >
          <el-icon><component :is="item.icon" /></el-icon>
        </div>
        <div class="mt-2 text-[11px] text-black/70 dark:text-white/70">
          {{ item.title }}
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
  import {
    Menu,
    Link,
    User,
    Service,
    Document,
    Reading,
    Files,
    Memo
  } from '@element-plus/icons-vue'
  import { useRouter } from 'vue-router'
  const router = useRouter()

  const toPath = (item) => {
    router.push({ name: item.path })
  }

  const openLink = (item) => {
    window.open(item.path, '_blank')
  }
  const shortcuts = [
    {
      icon: Menu,
      title: '菜单管理',
      path: 'menu'
    },
    {
      icon: Service,
      title: '角色管理',
      path: 'authority'
    },
    {
      icon: User,
      title: '用户管理',
      path: 'user'
    },
  ]

  
</script>

<style scoped lang="scss"></style>
