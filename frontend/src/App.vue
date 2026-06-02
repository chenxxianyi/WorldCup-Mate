<script setup lang="ts">
import { watchEffect, onMounted } from 'vue'
import { useSettingStore } from '@/stores/useSettingStore'
import { useAuthStore } from '@/stores/useAuthStore'

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
  <router-view />
</template>
