<script setup lang="ts">
import { watchEffect, onMounted } from 'vue'
import { useSettingStore } from '@/stores/useSettingStore'
import { useAuthStore } from '@/stores/useAuthStore'
import PullToRefresh from '@/components/common/PullToRefresh.vue'

const settings = useSettingStore()
const auth = useAuthStore()

watchEffect(() => {
  document.documentElement.dataset.theme = settings.theme
})

onMounted(() => {
  if (auth.token) {
    auth.fetchProfile()
  }
})
</script>

<template>
  <PullToRefresh />
  <router-view />
</template>
