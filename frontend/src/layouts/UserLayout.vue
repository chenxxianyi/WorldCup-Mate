<script setup lang="ts">
import TopBar from '@/components/common/TopBar.vue'
import BottomTabBar from '@/components/common/BottomTabBar.vue'
import DesktopRail from '@/components/common/DesktopRail.vue'
</script>

<template>
  <DesktopRail />
  <main class="app-shell">
    <TopBar />
    <router-view v-slot="{ Component }">
      <transition name="page" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </main>
  <BottomTabBar />
</template>

<style scoped>
.app-shell {
  width: min(100%, 1280px);
  margin: 0 auto;
  padding: 16px 16px calc(var(--nav-h) + 26px);
  position: relative;
}

@media (min-width: 768px) {
  .app-shell {
    padding: 18px 22px 40px 104px;
  }
}
</style>
