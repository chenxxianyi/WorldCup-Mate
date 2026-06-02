<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="admin-layout">
    <aside class="admin-sidebar">
      <div class="sidebar-brand">
        <div class="sidebar-mark">WM</div>
        <span>Admin</span>
      </div>
      <nav class="sidebar-nav">
        <button class="sidebar-link active" @click="router.push('/admin')">数据看板</button>
        <button class="sidebar-link" @click="router.push('/admin')">球队管理</button>
        <button class="sidebar-link" @click="router.push('/admin')">比赛管理</button>
        <button class="sidebar-link" @click="router.push('/admin')">积分榜管理</button>
      </nav>
      <button class="sidebar-link" @click="router.push('/')">返回前台</button>
    </aside>
    <main class="admin-main">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.admin-layout {
  display: grid;
  grid-template-columns: 220px 1fr;
  min-height: 100vh;
  background: var(--bg);
}

.admin-sidebar {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 16px 12px;
  border-right: 1px solid var(--line);
  background: var(--card);
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px;
  margin-bottom: 16px;
  font-weight: 800;
  font-size: 16px;
}

.sidebar-mark {
  width: 36px;
  height: 36px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  color: #fff;
  font-weight: 850;
  font-size: 14px;
  background: linear-gradient(145deg, var(--primary), var(--secondary));
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
}

.sidebar-link {
  width: 100%;
  padding: 10px 12px;
  border: 0;
  border-radius: 10px;
  text-align: left;
  color: var(--muted);
  background: transparent;
  font-size: 14px;
  font-weight: 600;
  transition: all 180ms ease;
}

.sidebar-link:hover {
  color: var(--text);
  background: var(--card-soft);
}

.sidebar-link.active {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 10%, transparent);
}

.admin-main {
  padding: 24px;
  overflow-x: auto;
}

@media (max-width: 768px) {
  .admin-layout {
    grid-template-columns: 1fr;
  }

  .admin-sidebar {
    flex-direction: row;
    overflow-x: auto;
    border-right: 0;
    border-bottom: 1px solid var(--line);
    padding: 8px;
  }

  .sidebar-nav {
    flex-direction: row;
    flex: 1;
  }

  .sidebar-link {
    white-space: nowrap;
  }
}
</style>
