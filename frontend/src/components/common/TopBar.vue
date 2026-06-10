<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/useAuthStore'
import { useNotificationStore } from '@/stores/useNotificationStore'
import { useSettingStore } from '@/stores/useSettingStore'
import { LOGOUT_LOGIN_PATH, clearAuthStorage } from '@/utils/logout'

const settings = useSettingStore()
const auth = useAuthStore()
const notification = useNotificationStore()

function refreshUnread(loggedIn = auth.isLoggedIn) {
  if (loggedIn) notification.fetchUnreadCount()
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
        <h1 class="brand-title" aria-label="WorldCup Mate">
          <span class="title-world">WorldCup</span>
          <span class="title-mate">Mate</span>
        </h1>
      </div>
    </div>

    <div class="top-actions">
      <template v-if="auth.isLoggedIn">
        <button class="icon-btn user-btn" title="个人中心" @click="$router.push('/profile')">
          <span class="btn-label">{{ auth.user?.username || '用户' }}</span>
          <span v-if="notification.unreadCount > 0" class="badge-dot">{{ notification.unreadCount > 9 ? '9+' : notification.unreadCount }}</span>
          <template v-if="auth.user?.avatar && auth.user.avatar.startsWith('/')">
            <img class="user-avatar-img" :src="auth.user.avatar" alt="头像" />
          </template>
          <template v-else>
            <span class="user-avatar">{{ auth.user?.avatar || 'U' }}</span>
          </template>
        </button>
        <a class="icon-btn logout-btn" :href="LOGOUT_LOGIN_PATH" title="退出登录" @click.capture="clearAuthStorage">
          <span class="btn-label">退出</span>
          <span class="action-orb logout-orb" aria-hidden="true">
            <svg class="action-icon" viewBox="0 0 24 24">
              <path d="M10 5H6.8A1.8 1.8 0 0 0 5 6.8v10.4A1.8 1.8 0 0 0 6.8 19H10" />
              <path d="M14 8l4 4-4 4" />
              <path d="M18 12H9" />
            </svg>
          </span>
        </a>
      </template>
      <template v-else>
        <button class="icon-btn login-btn primary-btn" title="登录" @click="$router.push('/login')">
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
  --topbar-glass: rgba(255, 255, 255, 0.82);
  --topbar-glass-soft: rgba(255, 255, 255, 0.72);
  --topbar-stroke: rgba(255, 255, 255, 0.9);
  position: relative;
  z-index: 10;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin: 0 0 18px;
  padding: 10px 12px;
  border: 1px solid var(--topbar-stroke);
  border-radius: 26px;
  background:
    radial-gradient(86% 150% at 10% 4%, color-mix(in srgb, var(--primary) 12%, transparent), transparent 44%),
    radial-gradient(96% 150% at 60% -28%, color-mix(in srgb, var(--accent) 22%, transparent), transparent 48%),
    radial-gradient(112% 160% at 96% 12%, color-mix(in srgb, var(--green-2) 13%, transparent), transparent 50%),
    linear-gradient(135deg, rgba(255, 255, 255, 0.96) 0%, #fff8ed 36%, #eefbf4 70%, #f3f7ff 100%),
    var(--topbar-glass);
  box-shadow:
    0 20px 54px rgba(15, 23, 42, 0.09),
    0 10px 24px color-mix(in srgb, var(--green) 8%, transparent),
    0 2px 10px rgba(255, 255, 255, 0.52) inset,
    0 -1px 0 rgba(255, 255, 255, 0.56) inset,
    0 0 0 1px rgba(255, 255, 255, 0.5) inset;
  backdrop-filter: blur(28px) saturate(1.45) contrast(1.03);
  -webkit-backdrop-filter: blur(28px) saturate(1.45) contrast(1.03);
  isolation: isolate;
  overflow: hidden;
}

.topbar::before {
  content: "";
  position: absolute;
  inset: 1px;
  z-index: 0;
  border-radius: 25px;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.78), transparent 31%),
    linear-gradient(315deg, rgba(255, 255, 255, 0.36), transparent 46%),
    radial-gradient(88% 150% at 50% -54%, rgba(255, 255, 255, 0.78), transparent 50%);
  mix-blend-mode: screen;
  pointer-events: none;
}

.topbar::after {
  content: "";
  position: absolute;
  left: 18px;
  right: 18px;
  bottom: 6px;
  z-index: 0;
  height: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.54);
  border-radius: 0 0 999px 999px;
  background:
    linear-gradient(
      90deg,
      transparent,
      color-mix(in srgb, var(--primary) 28%, transparent) 14%,
      color-mix(in srgb, var(--accent) 34%, transparent) 48%,
      color-mix(in srgb, var(--green-2) 28%, transparent) 74%,
      color-mix(in srgb, var(--sky) 20%, transparent) 88%,
      transparent
    );
  filter: blur(0.2px);
  opacity: 0.86;
  pointer-events: none;
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
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border: 1px solid color-mix(in srgb, var(--primary) 14%, rgba(255, 255, 255, 0.68));
  border-radius: 13px;
  color: var(--primary);
  background:
    radial-gradient(circle at 72% 20%, color-mix(in srgb, var(--accent) 15%, transparent), transparent 34%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.9), rgba(255, 247, 250, 0.58)),
    color-mix(in srgb, var(--primary) 5%, var(--card));
  box-shadow:
    0 12px 28px color-mix(in srgb, var(--primary) 14%, transparent),
    inset 0 1px 0 rgba(255, 255, 255, 0.72),
    inset 0 -1px 0 rgba(255, 255, 255, 0.18);
  backdrop-filter: blur(20px) saturate(1.45);
  -webkit-backdrop-filter: blur(20px) saturate(1.45);
  flex-shrink: 0;
}

.brand-mark svg {
  width: 24px;
  height: 24px;
  display: block;
  fill: currentColor;
}

.brand-title {
  position: relative;
  display: inline-flex;
  align-items: baseline;
  gap: 5px;
  margin: 0;
  font-size: clamp(20px, 3vw, 28px);
  line-height: 1;
  font-family: "Arial Black", "Impact", system-ui, sans-serif;
  font-weight: 900;
  letter-spacing: 0;
  transform: skewX(-4deg);
  text-transform: none;
  filter: drop-shadow(0 2px 0 rgba(15, 23, 42, 0.08));
}

.brand-title::before {
  content: "";
  position: absolute;
  left: 4px;
  right: 2px;
  bottom: -5px;
  height: 8px;
  border-bottom: 2px solid color-mix(in srgb, var(--primary) 58%, #19a974);
  border-radius: 0 0 999px 999px;
  transform: skewX(4deg) rotate(-1deg);
  opacity: 0.88;
}

.brand-title::after {
  content: "";
  position: absolute;
  left: 16%;
  right: 28%;
  bottom: -2px;
  height: 2px;
  border-radius: 999px;
  background: linear-gradient(90deg, transparent, #ffffff 24%, color-mix(in srgb, var(--primary) 66%, #ffffff) 50%, transparent);
  transform: skewX(4deg);
  opacity: 0.78;
}

.title-world {
  position: relative;
  color: transparent;
  background:
    linear-gradient(180deg, #10223f 0%, var(--text) 44%, #020617 100%);
  -webkit-background-clip: text;
  background-clip: text;
  text-shadow:
    0 1px 0 rgba(255, 255, 255, 0.55),
    0 5px 14px rgba(15, 23, 42, 0.12);
}

.title-world::after {
  content: "";
  position: absolute;
  left: 8%;
  right: 10%;
  top: 18%;
  height: 2px;
  border-radius: 999px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.68), transparent);
  opacity: 0.7;
}

.title-mate {
  position: relative;
  color: transparent;
  background:
    linear-gradient(180deg, #fff7cc 0%, #f4b940 42%, #b7791f 100%);
  -webkit-background-clip: text;
  background-clip: text;
  text-shadow:
    0 1px 0 rgba(255, 255, 255, 0.45),
    0 5px 14px rgba(180, 83, 9, 0.14);
}

.top-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.icon-btn {
  --icon-glass-bg: rgba(255, 255, 255, 0.28);
  --icon-glass-highlight: rgba(255, 255, 255, 0.5);
  --icon-glass-stroke: rgba(255, 255, 255, 0.62);
  position: relative;
  width: auto;
  min-width: 108px;
  height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 4px 6px 4px 12px;
  border: 1px solid var(--icon-glass-stroke);
  border-radius: 999px;
  color: var(--text);
  background:
    radial-gradient(circle at 30% 18%, var(--icon-glass-highlight), transparent 38%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.34), rgba(255, 255, 255, 0.08)),
    var(--icon-glass-bg);
  box-shadow:
    0 12px 26px rgba(15, 23, 42, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.62),
    inset 0 -1px 0 rgba(255, 255, 255, 0.16),
    inset 0 0 0 1px rgba(255, 255, 255, 0.12);
  backdrop-filter: blur(24px) saturate(1.5);
  -webkit-backdrop-filter: blur(24px) saturate(1.5);
  outline: none;
  transition:
    transform 180ms cubic-bezier(0.34, 1.56, 0.64, 1),
    border-color 180ms ease-out,
    background 180ms ease-out,
    box-shadow 180ms ease-out;
  font-size: 13px;
  font-weight: 750;
  text-decoration: none;
  cursor: pointer;
}

.icon-btn:hover {
  transform: translateY(-1px);
  border-color: color-mix(in srgb, var(--primary) 18%, rgba(255, 255, 255, 0.7));
  background:
    radial-gradient(circle at 30% 18%, rgba(255, 255, 255, 0.62), transparent 40%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.42), rgba(255, 255, 255, 0.12)),
    color-mix(in srgb, var(--accent) 8%, rgba(255, 255, 255, 0.24));
  box-shadow:
    0 16px 34px rgba(15, 23, 42, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.72),
    inset 0 0 0 1px rgba(255, 255, 255, 0.16);
}

.icon-btn:active {
  transform: scale(0.97);
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
  color: var(--muted);
  background:
    radial-gradient(circle at 34% 22%, rgba(255, 255, 255, 0.5), transparent 42%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.22), rgba(255, 255, 255, 0.06)),
    rgba(255, 255, 255, 0.14);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.58),
    inset 0 0 0 1px rgba(255, 255, 255, 0.18),
    0 8px 18px rgba(15, 23, 42, 0.06);
  backdrop-filter: blur(18px) saturate(1.35);
  -webkit-backdrop-filter: blur(18px) saturate(1.35);
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
  background:
    radial-gradient(circle at 32% 28%, color-mix(in srgb, var(--accent) 24%, rgba(255, 255, 255, 0.2)), transparent 46%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.26), rgba(255, 255, 255, 0.06)),
    rgba(255, 255, 255, 0.12);
}

.theme-orb.dark {
  color: #fff;
  background:
    linear-gradient(145deg, color-mix(in srgb, var(--blue) 74%, #000), color-mix(in srgb, var(--secondary) 58%, #000)),
    var(--blue);
}

.login-orb {
  background: rgba(255, 255, 255, 0.22);
}

.logout-orb {
  color: var(--muted);
  background:
    radial-gradient(circle at 34% 22%, rgba(255, 255, 255, 0.46), transparent 42%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.06)),
    color-mix(in srgb, var(--sky) 8%, rgba(255, 255, 255, 0.1));
}

.user-avatar {
  color: var(--primary);
  background:
    radial-gradient(circle at 34% 22%, rgba(255, 255, 255, 0.46), transparent 42%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.22), rgba(255, 255, 255, 0.06)),
    color-mix(in srgb, var(--primary) 9%, rgba(255, 255, 255, 0.08));
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

.primary-btn {
  color: #fff;
  background:
    linear-gradient(135deg, color-mix(in srgb, var(--primary) 86%, #fff), var(--primary));
  border-color: color-mix(in srgb, var(--primary) 52%, rgba(255, 255, 255, 0.56));
  box-shadow:
    0 12px 28px color-mix(in srgb, var(--primary) 24%, transparent),
    inset 0 1px 0 rgba(255, 255, 255, 0.34);
}

.primary-btn:hover {
  background:
    linear-gradient(135deg, color-mix(in srgb, var(--primary) 92%, #fff), color-mix(in srgb, var(--primary) 88%, #000));
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
  --topbar-glass: rgba(11, 18, 24, 0.78);
  --topbar-glass-soft: rgba(18, 28, 34, 0.72);
  --topbar-stroke: rgba(255, 255, 255, 0.18);
  background:
    radial-gradient(90% 150% at 9% 0%, color-mix(in srgb, var(--primary) 18%, transparent), transparent 42%),
    radial-gradient(94% 150% at 58% -28%, color-mix(in srgb, var(--accent) 14%, transparent), transparent 46%),
    radial-gradient(116% 160% at 98% 12%, color-mix(in srgb, var(--green-2) 16%, transparent), transparent 50%),
    linear-gradient(135deg, rgba(18, 28, 34, 0.92), rgba(10, 23, 27, 0.82) 52%, rgba(14, 22, 40, 0.86)),
    var(--topbar-glass);
  box-shadow:
    0 22px 60px rgba(0, 0, 0, 0.28),
    inset 0 1px 0 rgba(255, 255, 255, 0.18),
    inset 0 -1px 0 rgba(255, 255, 255, 0.08),
    inset 0 0 0 1px rgba(255, 255, 255, 0.04);
}

[data-theme='dark'] .topbar::before {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.22), transparent 34%),
    linear-gradient(315deg, rgba(255, 255, 255, 0.1), transparent 42%),
    radial-gradient(95% 180% at 50% -52%, rgba(255, 255, 255, 0.16), transparent 48%);
}

[data-theme='dark'] .title-world {
  background:
    linear-gradient(180deg, #ffffff 0%, #d9f99d 28%, #7dd3fc 100%);
  -webkit-background-clip: text;
  background-clip: text;
  text-shadow:
    0 1px 0 rgba(255, 255, 255, 0.12),
    0 8px 18px rgba(56, 189, 248, 0.2);
}

[data-theme='dark'] .brand-title::before {
  border-bottom-color: rgba(52, 211, 153, 0.82);
}

[data-theme='dark'] .brand-title::after {
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.82), rgba(52, 211, 153, 0.8), transparent);
}

[data-theme='dark'] .icon-btn {
  --icon-glass-bg: rgba(255, 255, 255, 0.08);
  --icon-glass-highlight: rgba(255, 255, 255, 0.2);
  --icon-glass-stroke: rgba(255, 255, 255, 0.14);
  color: var(--text);
  background:
    radial-gradient(circle at 30% 18%, var(--icon-glass-highlight), transparent 40%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.03)),
    var(--icon-glass-bg);
  box-shadow:
    0 12px 28px rgba(0, 0, 0, 0.24),
    inset 0 1px 0 rgba(255, 255, 255, 0.16),
    inset 0 0 0 1px rgba(255, 255, 255, 0.04);
}

[data-theme='dark'] .brand-mark,
[data-theme='dark'] .action-orb,
[data-theme='dark'] .user-avatar {
  border-color: rgba(255, 255, 255, 0.1);
  background:
    radial-gradient(circle at 34% 22%, rgba(255, 255, 255, 0.16), transparent 42%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.03)),
    color-mix(in srgb, var(--primary) 6%, transparent);
}

@media (max-width: 720px) {
  .topbar {
    align-items: center;
    flex-direction: row;
    padding: 8px;
    border-radius: 24px;
  }

  .brand {
    gap: 9px;
    flex: 1;
  }

  .brand-mark {
    width: 34px;
    height: 34px;
  }

  .brand-mark svg {
    width: 24px;
    height: 24px;
  }

  .brand-title {
    font-size: 20px;
  }

  .icon-btn {
    width: 44px;
    min-width: 0;
    height: 44px;
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
    width: 32px;
    height: 32px;
  }

  .logout-btn,
  .theme-btn {
    border-color: var(--icon-glass-stroke);
    background:
      radial-gradient(circle at 30% 18%, var(--icon-glass-highlight), transparent 38%),
      linear-gradient(145deg, rgba(255, 255, 255, 0.34), rgba(255, 255, 255, 0.08)),
      var(--icon-glass-bg);
  }

  [data-theme='dark'] .logout-btn,
  [data-theme='dark'] .theme-btn {
    border-color: var(--icon-glass-stroke);
    background:
      radial-gradient(circle at 30% 18%, var(--icon-glass-highlight), transparent 40%),
      linear-gradient(145deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.03)),
      var(--icon-glass-bg);
  }
}

@media (max-width: 380px) {
  .topbar {
    gap: 8px;
    padding: 7px;
  }

  .brand {
    gap: 7px;
  }

  .brand-title {
    font-size: 18px;
  }

  .top-actions {
    gap: 6px;
  }

  .icon-btn {
    width: 44px;
    height: 44px;
  }
}
</style>
