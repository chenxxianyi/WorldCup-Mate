<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import { useSettingStore } from '@/stores/useSettingStore'
import { useCompetitionStore } from '@/stores/useCompetitionStore'
import { seasonLabel } from '@/types/competition'

const settings = useSettingStore()
const auth = useAuthStore()
const notification = useNotificationStore()
const comp = useCompetitionStore()
const router = useRouter()

const brandSub = computed(() => {
  if (comp.isWorldCup) return '2026 世界杯赛程助手'
  return comp.current ? `${comp.current.name} 赛程助手` : '足球赛事助手'
})

const versionBadge = computed(() => {
  if (comp.isWorldCup) return '美加墨'
  return comp.current ? seasonLabel(comp.current.season) : ''
})

function refreshUnread() {
  if (auth.isLoggedIn) notification.fetchUnreadCount()
}

onMounted(() => {
  refreshUnread()
  comp.fetchCompetitions()
})
watch(() => auth.isLoggedIn, refreshUnread)
</script>

<template>
  <header class="topbar">
    <div class="brand">
      <div class="brand-mark">
        <svg viewBox="0 0 24 24" aria-label="WorldCup Mate">
          <path
            d="M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm0 3.2 2.8 2-1.1 3.3h-3.4L9.2 7.2 12 5.2Zm-6.5 4 2.9.8 1 3.2-2.2 1.9a7 7 0 0 1-2.3-1.3A7.9 7.9 0 0 1 5.5 9.2Zm2.4 8.4a9 9 0 0 1 1.8-2l2.3 1.3 2.3-1.3a9 9 0 0 1 1.8 2 7.8 7.8 0 0 1-8.2 0Zm9.9-2.5-2.2-1.9 1-3.2 2.9-.8a7.9 7.9 0 0 1 .6 3.6 7 7 0 0 1-2.3 1.3ZM14.5 4a7.8 7.8 0 0 1 4 2.5l-2.8.8-2.5-1.8V3.9c.4 0 .8.1 1.3.1ZM10.8 3.9v1.6L8.3 7.3l-2.8-.8a7.8 7.8 0 0 1 4-2.5c.5 0 .9-.1 1.3-.1Z"
          />
        </svg>
      </div>
      <div>
        <h1 class="brand-title">WorldCup Mate</h1>
        <p class="brand-sub">{{ brandSub }} <span v-if="versionBadge" class="version-badge">{{ versionBadge }}</span></p>
      </div>
    </div>

    <div class="comp-switch">
      <select
        v-model="comp.currentCode"
        class="comp-select"
        aria-label="切换赛事"
        title="切换赛事"
        @change="comp.setCurrent(comp.currentCode)"
      >
        <option v-for="c in comp.competitions" :key="c.code" :value="c.code">{{ c.name }}</option>
      </select>
      <span class="material-symbols-outlined comp-caret" aria-hidden="true">expand_more</span>
    </div>

    <div class="top-actions">
      <template v-if="auth.isLoggedIn">
        <button class="icon-btn user-btn" title="个人中心" @click="router.push('/profile')">
          <span class="btn-label">{{ auth.user?.username || '用户' }}</span>
          <span v-if="notification.unreadCount > 0" class="badge-dot">{{ notification.unreadCount > 9 ? '9+' : notification.unreadCount }}</span>
          <template v-if="auth.user?.avatar && auth.user.avatar.startsWith('/')">
            <img class="user-avatar-img" :src="auth.user.avatar" alt="头像" />
          </template>
          <template v-else>
            <span class="user-avatar">{{ auth.user?.avatar || 'U' }}</span>
          </template>
        </button>
        <button class="icon-btn logout-btn" title="退出登录" @click="auth.logout()">
          <span class="btn-label">退出</span>
          <span class="action-orb logout-orb" aria-hidden="true">
            <svg class="action-icon" viewBox="0 0 24 24">
              <path d="M10 5H6.8A1.8 1.8 0 0 0 5 6.8v10.4A1.8 1.8 0 0 0 6.8 19H10" />
              <path d="M14 8l4 4-4 4" />
              <path d="M18 12H9" />
            </svg>
          </span>
        </button>
      </template>
      <template v-else>
        <button class="icon-btn login-btn" title="登录" @click="router.push('/login')">
          <span class="btn-label">登录</span>
          <span class="action-orb login-orb" aria-hidden="true">
            <svg class="action-icon" viewBox="0 0 24 24">
              <path d="M14 8l4 4-4 4" />
              <path d="M18 12H9" />
              <path d="M10 19h7.2a1.8 1.8 0 0 0 1.8-1.8V6.8A1.8 1.8 0 0 0 17.2 5H10" />
            </svg>
          </span>
        </button>
      </template>
      <button class="icon-btn theme-btn" :title="settings.theme === 'dark' ? '切换浅色模式' : '切换深色模式'" @click="settings.toggleTheme">
        <span class="btn-label">{{ settings.theme === 'dark' ? '深色' : '浅色' }}</span>
        <span class="action-orb theme-orb" :class="{ dark: settings.theme === 'dark' }" aria-hidden="true">
          <svg v-if="settings.theme === 'dark'" class="action-icon" viewBox="0 0 24 24">
            <path d="M19 14.2A7 7 0 0 1 9.8 5a7.5 7.5 0 1 0 9.2 9.2Z" />
          </svg>
          <svg v-else class="action-icon" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="4" />
            <path d="M12 2v2.2M12 19.8V22M4.9 4.9l1.6 1.6M17.5 17.5l1.6 1.6M2 12h2.2M19.8 12H22M4.9 19.1l1.6-1.6M17.5 6.5l1.6-1.6" />
          </svg>
        </span>
      </button>
    </div>
  </header>
</template>

<style scoped>
.topbar {
  position: relative;
  z-index: 10;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin: 2px 0 16px;
  padding: 12px 14px;
  border: 1px solid rgba(255, 255, 255, 0.62);
  border-radius: 24px;
  background:
    linear-gradient(102deg, rgba(246, 218, 229, 0.96), rgba(224, 226, 239, 0.96)),
    var(--card);
  box-shadow: 0 18px 48px rgba(32, 34, 45, 0.08);
  backdrop-filter: blur(18px);
  overflow: hidden;
}

.topbar::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 16% 14%, rgba(220, 20, 60, 0.1), transparent 30%),
    radial-gradient(circle at 86% 30%, rgba(26, 35, 126, 0.08), transparent 34%);
  pointer-events: none;
}

.topbar > * {
  position: relative;
  z-index: 1;
}

.brand {
  display: flex;
  align-items: center;
  gap: 18px;
  min-width: 0;
}

.brand-mark {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 14px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 10%, transparent);
  flex-shrink: 0;
}

.brand-mark svg {
  width: 32px;
  height: 32px;
  display: block;
  fill: currentColor;
}

.brand-title {
  margin: 0;
  font-size: clamp(23px, 3vw, 34px);
  line-height: 1.05;
  font-weight: 900;
}

.brand-sub {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: clamp(12px, 1.35vw, 16px);
  line-height: 1.25;
}

.version-badge {
  display: inline-flex;
  align-items: center;
  min-height: 20px;
  margin-left: 8px;
  padding: 0 8px;
  border-radius: 999px;
  color: #0b0b0b;
  background: var(--accent);
  font-size: 12px;
  font-weight: 900;
}

.comp-switch {
  position: relative;
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
}

.comp-select {
  appearance: none;
  -webkit-appearance: none;
  min-width: 132px;
  height: 40px;
  padding: 0 34px 0 14px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 14px;
  color: #111;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 18px rgba(30, 30, 40, 0.12);
  font-size: 14px;
  font-weight: 800;
  outline: none;
  cursor: pointer;
  transition: transform 180ms ease-out, box-shadow 180ms ease-out;
}

.comp-select:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.09);
}

.comp-caret {
  position: absolute;
  right: 10px;
  pointer-events: none;
  color: var(--muted);
  font-size: 18px;
}

.top-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.icon-btn {
  position: relative;
  width: auto;
  min-width: 146px;
  height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 6px 8px 6px 14px;
  border: 1px solid rgba(255, 255, 255, 0.72);
  border-radius: 18px;
  color: #111;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 18px rgba(30, 30, 40, 0.12);
  backdrop-filter: blur(16px);
  outline: none;
  transition: transform 180ms ease-out, box-shadow 180ms ease-out;
  font-size: 14px;
  font-weight: 800;
  cursor: pointer;
}

.icon-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.09);
}

.badge-dot {
  position: absolute;
  top: -5px;
  right: -5px;
  min-width: 18px;
  height: 18px;
  display: grid;
  place-items: center;
  padding: 0 5px;
  border: 2px solid #fff;
  border-radius: 999px;
  color: #fff;
  background: var(--primary);
  font-size: 10px;
  line-height: 1;
}

.action-orb,
.user-avatar {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: #fff;
  background: var(--primary);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.22);
  font-size: 16px;
}

.action-icon {
  width: 18px;
  height: 18px;
  display: block;
  fill: none;
  stroke: currentColor;
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.theme-orb {
  color: #fff;
  background:
    radial-gradient(circle at 32% 30%, rgba(255, 255, 255, 0.28) 0 18%, transparent 19%),
    linear-gradient(145deg, var(--green-2), var(--green));
}

.theme-orb.dark {
  color: #fff;
  background:
    radial-gradient(circle at 68% 34%, rgba(255, 255, 255, 0.22) 0 9%, transparent 10%),
    linear-gradient(145deg, #00a66b, #004d35);
}

.login-orb {
  background: linear-gradient(145deg, var(--green), var(--green-2));
}

.logout-orb {
  color: #fff;
  background: linear-gradient(145deg, #e5484d, #b4232c);
}

.user-avatar {
  background: linear-gradient(145deg, var(--primary), var(--secondary));
  font-size: 14px;
  font-weight: 800;
}

.user-avatar-img {
  width: 32px;
  height: 32px;
  border-radius: 999px;
  object-fit: cover;
}

.logout-btn {
  min-width: auto;
  padding: 6px 14px;
  color: var(--muted);
}

.login-btn {
  min-width: 100px;
}

.theme-btn {
  min-width: 116px;
}

[data-theme='dark'] .topbar {
  border-color: rgba(255, 255, 255, 0.08);
  background:
    linear-gradient(102deg, rgba(44, 16, 28, 0.96), rgba(25, 28, 55, 0.96)),
    var(--card);
}

[data-theme='dark'] .icon-btn {
  color: var(--text);
  background: rgba(17, 17, 17, 0.82);
  border-color: rgba(255, 255, 255, 0.08);
}

[data-theme='dark'] .comp-select {
  color: var(--text);
  background: rgba(17, 17, 17, 0.82);
  border-color: rgba(255, 255, 255, 0.08);
}

@media (max-width: 720px) {
  .topbar {
    align-items: center;
    flex-direction: row;
    flex-wrap: wrap;
    padding: 10px 10px 10px 12px;
    border-radius: 22px;
  }

  .comp-switch {
    order: 3;
    width: 100%;
  }

  .comp-select {
    width: 100%;
  }

  .brand {
    gap: 9px;
    flex: 1;
  }

  .brand-mark {
    width: 36px;
    height: 36px;
  }

  .brand-mark svg {
    width: 28px;
    height: 28px;
  }

  .brand-title {
    font-size: 21px;
  }

  .brand-sub {
    font-size: 11px;
  }

  .version-badge {
    min-height: 18px;
    padding: 0 6px;
    font-size: 10px;
  }

  .icon-btn {
    width: 42px;
    min-width: 0;
    height: 42px;
    justify-content: center;
    padding: 0;
    border-radius: 16px;
    font-size: 16px;
  }

  .btn-label {
    display: none;
  }

  .action-orb,
  .user-avatar,
  .user-avatar-img {
    width: 30px;
    height: 30px;
  }

  .logout-btn,
  .theme-btn {
    border-color: rgba(255, 255, 255, 0.82);
    background: rgba(255, 255, 255, 0.86);
  }
}
</style>
