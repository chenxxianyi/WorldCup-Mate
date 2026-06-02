<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useSettingStore } from '@/stores/useSettingStore'

const route = useRoute()
const router = useRouter()
const settings = useSettingStore()

const navItems = [
  { icon: '⌂', route: '/', title: '首页' },
  { icon: '□', route: '/schedule', title: '赛程' },
  { icon: '⚑', route: '/teams', title: '球队' },
  { icon: '▥', route: '/standings', title: '积分榜' },
  { icon: '○', route: '/profile', title: '我的' },
  { icon: '▤', route: '/admin', title: '后台' },
]

function isActive(itemRoute: string) {
  if (itemRoute === '/') return route.path === '/'
  return route.path.startsWith(itemRoute)
}
</script>

<template>
  <aside class="desktop-rail" aria-label="桌面导航">
    <div class="rail-mark">WM</div>
    <nav class="rail-nav">
      <button
        v-for="item in navItems"
        :key="item.route"
        class="rail-btn"
        :class="{ active: isActive(item.route) }"
        :title="item.title"
        @click="router.push(item.route)"
      >
        {{ item.icon }}
      </button>
    </nav>
    <button class="rail-btn" title="切换主题" @click="settings.toggleTheme">◐</button>
  </aside>
</template>

<style scoped>
.desktop-rail {
  display: none;
}

@media (min-width: 768px) {
  .desktop-rail {
    position: fixed;
    inset: 18px auto 18px 20px;
    width: 70px;
    z-index: 20;
    display: grid;
    grid-template-rows: auto 1fr auto;
    gap: 12px;
    padding: 10px;
    border: 1px solid var(--line);
    border-radius: 24px;
    background: color-mix(in srgb, var(--card) 88%, transparent);
    box-shadow: var(--shadow);
    backdrop-filter: blur(18px);
  }

  .rail-mark {
    width: 48px;
    height: 48px;
    display: grid;
    place-items: center;
    border-radius: 17px;
    color: #fff;
    font-weight: 850;
    background: linear-gradient(145deg, var(--primary), var(--secondary));
  }

  .rail-nav {
    display: grid;
    align-content: center;
    gap: 8px;
  }

  .rail-btn {
    width: 48px;
    height: 48px;
    border: 0;
    border-radius: 17px;
    color: var(--weak);
    background: transparent;
    font-size: 18px;
  }

  .rail-btn.active {
    color: var(--primary);
    background: color-mix(in srgb, var(--primary) 12%, transparent);
  }
}
</style>
