<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ThemeIcon from './ThemeIcon.vue'

const route = useRoute()
const router = useRouter()
const navItems = [
  { path: '/', label: '首页', icon: 'home' },
  { path: '/schedule', label: '赛程', icon: 'calendar' },
  { path: '/teams', label: '球队', icon: 'shield' },
  { path: '/standings', label: '积分榜', icon: 'table' },
  { path: '/profile', label: '我的', icon: 'user' },
]

const activePath = computed(() => {
  if (route.path.startsWith('/matches/')) return '/schedule'
  if (route.path.startsWith('/teams/')) return '/teams'
  if (route.path === '/login') return '/profile'
  return route.path
})
</script>

<template>
  <aside
    class="desktop-rail"
    aria-label="主要导航"
  >
    <button
      v-for="item in navItems"
      :key="item.path"
      class="rail-link"
      :class="{ active: activePath === item.path }"
      type="button"
      :aria-current="activePath === item.path ? 'page' : undefined"
      @click="router.push(item.path)"
    >
      <ThemeIcon :name="item.icon" /><span>{{ item.label }}</span>
    </button>
  </aside>

  <nav
    class="bottom-nav"
    aria-label="主要导航"
  >
    <button
      v-for="item in navItems"
      :key="item.path"
      class="bottom-link"
      :class="{ active: activePath === item.path }"
      type="button"
      :aria-current="activePath === item.path ? 'page' : undefined"
      @click="router.push(item.path)"
    >
      <ThemeIcon :name="item.icon" /><span>{{ item.label }}</span>
    </button>
  </nav>
</template>
