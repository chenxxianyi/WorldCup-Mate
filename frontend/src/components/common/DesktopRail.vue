<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useSettingStore } from '@/stores/useSettingStore'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const settings = useSettingStore()

const baseNavItems = [
  { icon: 'home', route: '/', title: '首页' },
  { icon: 'calendar_month', route: '/schedule', title: '赛程' },
  { icon: 'flag', route: '/teams', title: '球队' },
  { icon: 'trophy', route: '/standings', title: '积分' },
  { icon: 'auto_awesome', route: '/ai', title: 'AI' },
  { icon: 'person', route: '/profile', title: '我的' },
]

const adminNavItems = [
  { icon: 'admin_panel_settings', route: '/admin', title: '后台' },
]

const navItems = computed(() => auth.isAdmin ? [...baseNavItems, ...adminNavItems] : baseNavItems)

function isActive(itemRoute: string) {
  if (itemRoute === '/') return route.path === '/'
  if (itemRoute === '/schedule') return route.path.startsWith('/schedule') || route.path.startsWith('/matches') || route.path.startsWith('/bracket')
  if (itemRoute === '/teams') return route.path.startsWith('/teams')
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
        :aria-label="item.title"
        :aria-current="isActive(item.route) ? 'page' : undefined"
        @click="router.push(item.route)"
      >
        <span class="rail-icon material-symbols-outlined" aria-hidden="true">{{ item.icon }}</span>
        <span class="rail-label">{{ item.title }}</span>
      </button>
    </nav>
    <button class="rail-btn" title="切换主题" aria-label="切换主题" @click="settings.toggleTheme">
      <span class="rail-icon material-symbols-outlined" aria-hidden="true">contrast</span>
      <span class="rail-label">主题</span>
    </button>
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
    overflow: hidden;
    transition:
      width 320ms cubic-bezier(0.2, 0.8, 0.2, 1),
      border-color 240ms ease,
      background 240ms ease,
      box-shadow 240ms ease;
  }

  .desktop-rail:hover,
  .desktop-rail:focus-within {
    width: 126px;
    border-color: color-mix(in srgb, var(--primary) 22%, var(--line));
    background: color-mix(in srgb, var(--card) 96%, transparent);
    box-shadow: var(--shadow-hover);
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
    flex-shrink: 0;
  }

  .rail-nav {
    display: grid;
    align-content: center;
    gap: 8px;
  }

  .rail-btn {
    width: 48px;
    height: 48px;
    display: inline-flex;
    align-items: center;
    justify-content: flex-start;
    gap: 10px;
    padding: 0 14px;
    border: 0;
    border-radius: 17px;
    color: var(--weak);
    background: transparent;
    overflow: hidden;
    white-space: nowrap;
    transition:
      width 320ms cubic-bezier(0.2, 0.8, 0.2, 1),
      color 240ms ease,
      background 240ms ease,
      border-radius 240ms ease;
  }

  .desktop-rail:hover .rail-btn,
  .desktop-rail:focus-within .rail-btn {
    width: 100%;
    border-radius: 16px;
  }

  .rail-icon {
    flex: 0 0 20px;
    font-size: 22px;
  }

  .rail-label {
    max-width: 0;
    overflow: hidden;
    color: var(--text);
    font-size: 14px;
    font-weight: 700;
    line-height: 1;
    opacity: 0;
    transform: translateX(-4px);
    transition:
      max-width 320ms cubic-bezier(0.2, 0.8, 0.2, 1),
      opacity 240ms ease,
      transform 260ms ease;
  }

  .desktop-rail:hover .rail-label,
  .desktop-rail:focus-within .rail-label {
    max-width: 44px;
    opacity: 1;
    transform: translateX(0);
  }

  .rail-btn:hover,
  .rail-btn:focus-visible {
    color: var(--text);
    background: var(--card-soft);
    outline: none;
  }

  .rail-btn.active {
    color: var(--primary);
    background: color-mix(in srgb, var(--primary) 12%, transparent);
  }

  .rail-btn.active .rail-label {
    color: var(--primary);
  }

  @media (prefers-reduced-motion: reduce) {
    .desktop-rail,
    .rail-btn,
    .rail-label {
      transition: none;
    }
  }
}
</style>
