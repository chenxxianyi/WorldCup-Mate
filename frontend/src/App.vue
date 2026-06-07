<script setup lang="ts">
import { watchEffect, onMounted } from 'vue'
import { useSettingStore } from '@/stores/useSettingStore'
import { useAuthStore } from '@/stores/useAuthStore'
import PullToRefresh from '@/components/common/PullToRefresh.vue'
import { LOGOUT_QUERY_KEY, LOGOUT_QUERY_VALUE, clearAuthStorage } from '@/utils/logout'

const settings = useSettingStore()
const auth = useAuthStore()

watchEffect(() => {
  document.documentElement.dataset.theme = settings.theme
})

onMounted(async () => {
  const query = new URLSearchParams(window.location.search)
  if (query.get(LOGOUT_QUERY_KEY) === LOGOUT_QUERY_VALUE) {
    clearAuthStorage()
    auth.logout()
    return
  }

  if (auth.token) {
    await auth.fetchProfile()
  }
})
</script>

<template>
  <PullToRefresh />
  <router-view />
</template>
