<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useSettingStore } from '@/stores/useSettingStore'
import { useAuthStore } from '@/stores/useAuthStore'

const settings = useSettingStore()
const auth = useAuthStore()
const router = useRouter()
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
          <span>{{ auth.user?.username || '用户' }}</span>
          <template v-if="auth.user?.avatar && auth.user.avatar.startsWith('/')">
            <img class="user-avatar-img" :src="auth.user.avatar" alt="头像" />
          </template>
          <template v-else>
            <span class="user-avatar">{{ auth.user?.avatar || 'U' }}</span>
          </template>
        </button>
        <button class="icon-btn logout-btn" title="退出登录" @click="auth.logout()">
          <span>退出</span>
          <span class="theme-knob logout-knob">⏻</span>
        </button>
      </template>
      <template v-else>
        <button class="icon-btn login-btn" title="登录" @click="router.push('/login')">
          <span>登录</span>
          <span class="theme-knob" style="background: var(--green)">→</span>
        </button>
      </template>
      <button class="icon-btn" title="切换主题" @click="settings.toggleTheme">
        <span>主题切换</span>
        <span class="theme-knob">{{ settings.theme === 'dark' ? '☾' : '☀' }}</span>
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

.top-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.icon-btn {
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
  color: #111111;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 8px 18px rgba(30, 30, 40, 0.12);
  backdrop-filter: blur(16px);
  outline: none;
  transition: transform 180ms ease-out, box-shadow 180ms ease-out;
  font-size: 14px;
  font-weight: 850;
}

.icon-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 28px rgba(0, 0, 0, 0.09);
}

.theme-knob {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: #fff;
  background: var(--primary);
  box-shadow: inset 0 0 0 8px rgba(255, 255, 255, 0.08);
  font-size: 16px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  color: #fff;
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

.logout-knob {
  background: #dc3545;
}

.login-btn {
  min-width: 100px;
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

@media (max-width: 720px) {
  .topbar {
    align-items: center;
    flex-direction: row;
    padding: 10px 10px 10px 12px;
    border-radius: 22px;
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

  .icon-btn span:first-child {
    display: none;
  }

  .theme-knob {
    width: 30px;
    height: 30px;
  }
}
</style>
