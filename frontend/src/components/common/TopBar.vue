<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import { useSettingStore } from '@/stores/useSettingStore'

const settings = useSettingStore()
const auth = useAuthStore()
const notification = useNotificationStore()
const router = useRouter()

function refreshUnread() {
  if (auth.isLoggedIn) notification.fetchUnreadCount()
}

onMounted(refreshUnread)
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
        <p class="brand-sub">2026 世界杯赛程助手 <span class="version-badge">美加墨</span></p>
      </div>
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
        <button class="icon-btn login-btn primary-btn" title="登录" @click="router.push('/login')">
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
  gap: 12px;
  margin: 0 0 18px;
  padding: 10px 2px;
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  backdrop-filter: none;
  overflow: hidden;
}

.topbar::before {
  content: none;
}

.topbar > * {
  position: relative;
  z-index: 1;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.brand-mark {
  width: 34px;
  height: 34px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 10px;
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 8%, var(--card));
  flex-shrink: 0;
}

.brand-mark svg {
  width: 24px;
  height: 24px;
  display: block;
  fill: currentColor;
}

.brand-title {
  margin: 0;
  font-size: clamp(20px, 3vw, 28px);
  line-height: 1.08;
  font-weight: 850;
}

.brand-sub {
  margin: 4px 0 0;
  color: var(--muted);
  font-size: clamp(11px, 1.35vw, 14px);
  line-height: 1.25;
}

.version-badge {
  display: inline-flex;
  align-items: center;
  min-height: 18px;
  margin-left: 6px;
  padding: 0 7px;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  font-size: 11px;
  font-weight: 750;
}

.top-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.icon-btn {
  position: relative;
  width: auto;
  min-width: 108px;
  height: 38px;
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 4px 6px 4px 12px;
  border: 1px solid var(--line);
  border-radius: 999px;
  color: var(--text);
  background: var(--card);
  box-shadow: none;
  backdrop-filter: none;
  outline: none;
  transition: border-color 180ms ease-out, background 180ms ease-out;
  font-size: 13px;
  font-weight: 750;
  cursor: pointer;
}

.icon-btn:hover {
  border-color: color-mix(in srgb, var(--primary) 24%, var(--line));
  background: color-mix(in srgb, var(--primary) 4%, var(--card));
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
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: var(--muted);
  background: var(--card-soft);
  box-shadow: none;
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
  color: var(--muted);
  background: var(--card-soft);
}

.theme-orb.dark {
  color: var(--muted);
  background: var(--card-soft);
}

.login-orb {
  background: var(--card-soft);
}

.logout-orb {
  color: var(--muted);
  background: var(--card-soft);
}

.user-avatar {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 9%, var(--card));
  font-size: 14px;
  font-weight: 800;
}

.user-avatar-img {
  width: 28px;
  height: 28px;
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

.primary-btn {
  color: #fff;
  background: var(--primary);
  border-color: var(--primary);
}

.primary-btn:hover {
  background: color-mix(in srgb, var(--primary) 85%, #000);
  border-color: color-mix(in srgb, var(--primary) 85%, #000);
}

.primary-btn .action-orb {
  color: #fff;
  background: rgba(255, 255, 255, 0.2);
}

.theme-btn {
  min-width: 116px;
}

[data-theme='dark'] .topbar {
  background: transparent;
}

[data-theme='dark'] .icon-btn {
  color: var(--text);
  background: var(--card);
  border-color: var(--line);
}

@media (max-width: 720px) {
  .topbar {
    align-items: center;
    flex-direction: row;
    padding: 8px 0 10px;
    border-radius: 0;
  }

  .brand {
    gap: 9px;
    flex: 1;
  }

  .brand-mark {
    width: 32px;
    height: 32px;
  }

  .brand-mark svg {
    width: 23px;
    height: 23px;
  }

  .brand-title {
    font-size: 20px;
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
    width: 38px;
    min-width: 0;
    height: 38px;
    justify-content: center;
    padding: 0;
    border-radius: 999px;
    font-size: 16px;
  }

  .btn-label {
    display: none;
  }

  .action-orb,
  .user-avatar,
  .user-avatar-img {
    width: 28px;
    height: 28px;
  }

  .logout-btn,
  .theme-btn {
    border-color: var(--line);
    background: var(--card);
  }
}
</style>
