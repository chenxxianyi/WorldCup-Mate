<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const showBackTop = ref(false)

function updateBackTopVisibility() {
  showBackTop.value = window.scrollY > 420
}

function backToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', updateBackTopVisibility, { passive: true })
  updateBackTopVisibility()
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', updateBackTopVisibility)
})

watch(
  () => route.fullPath,
  async () => {
    await nextTick()
    updateBackTopVisibility()
  },
)
</script>

<template>
  <Transition name="back-top">
    <button
      v-if="showBackTop"
      class="back-top-btn"
      type="button"
      aria-label="返回顶部"
      title="返回顶部"
      @click="backToTop"
    >
      <span class="material-symbols-outlined">keyboard_arrow_up</span>
    </button>
  </Transition>
</template>
