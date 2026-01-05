<template>
  <div class="user-dashboard">
    <h1>欢迎回来，{{ authStore.user?.user_name }}！</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 16 12 14 15 10 15 8 12 2 12"></polyline><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"></path></svg>
        </div>
        <div class="stat-info">
          <h3>待取包裹</h3>
          <p class="stat-value">{{ pendingCount }}</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
        </div>
        <div class="stat-info">
          <h3>寄件中</h3>
          <p class="stat-value">{{ inTransitCount }}</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
        </div>
        <div class="stat-info">
          <h3>已取包裹</h3>
          <p class="stat-value">{{ checkedOutCount }}</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="20" x2="18" y2="10"></line><line x1="12" y1="20" x2="12" y2="4"></line><line x1="6" y1="20" x2="6" y2="14"></line></svg>
        </div>
        <div class="stat-info">
          <h3>总计</h3>
          <p class="stat-value">{{ totalCount }}</p>
        </div>
      </div>
    </div>

    <div class="quick-actions">
      <h2>快捷操作</h2>
      <div class="action-buttons">
        <router-link to="/user/packs" class="action-btn">
          <span class="action-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
          </span>
          <span>查看包裹</span>
        </router-link>
        <router-link to="/user/mail" class="action-btn">
          <span class="action-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
          </span>
          <span>我要寄件</span>
        </router-link>
      </div>
    </div>

    <div v-if="recentPacks.length > 0" class="recent-packs">
      <h2>最近包裹</h2>
      <div class="pack-list">
        <div v-for="pack in recentPacks" :key="pack.pack_id" class="pack-item">
          <div class="pack-status" :class="pack.pack_status">
            {{ getStatusText(pack.pack_status) }}
          </div>
          <div class="pack-details">
            <p class="pack-id">包裹号: {{ pack.pack_id }}</p>
            <p v-if="pack.pickup_code" class="pickup-code">取件码: {{ pack.pickup_code }}</p>
            <p class="pack-time">{{ formatTime(pack.created_at) }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usePackStore } from '@/stores/counter'

const authStore = useAuthStore()
const packStore = usePackStore()

const packs = computed(() => packStore.packs)

const pendingCount = computed(
  () => packs.value.filter((p) => p.pack_status === 'pending' || p.pack_status === 'arrived').length
)
const inTransitCount = computed(
  () => packs.value.filter((p) => p.pack_status === 'in_transit' || p.pack_status === 'shipped').length
)
const checkedOutCount = computed(() => packs.value.filter((p) => p.pack_status === 'checked_out').length)
const totalCount = computed(() => packs.value.length)

const recentPacks = computed(() => packs.value.slice(0, 5))

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    pending: '待取',
    arrived: '已到达',
    checked_out: '已取件',
    in_transit: '运输中',
    shipped: '已发货',
    cancelled: '已取消'
  }
  return statusMap[status] || status
}

const formatTime = (time?: string) => {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(async () => {
  if (authStore.user?.user_id) {
    try {
      await packStore.fetchUserPacks(authStore.user.user_id)
    } catch (error) {
      console.error('获取包裹列表失败:', error)
    }
  }
})
</script>

<style scoped>
.user-dashboard {
  max-width: 1200px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 1.5rem;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 3rem;
}

.stat-info h3 {
  margin: 0;
  color: #666;
  font-size: 0.95rem;
  font-weight: 500;
}

.stat-value {
  margin: 0.5rem 0 0;
  font-size: 2rem;
  font-weight: 700;
  color: #667eea;
}

.quick-actions {
  margin-bottom: 3rem;
}

.quick-actions h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
}

.action-btn {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.5rem;
  border-radius: 1rem;
  text-decoration: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  font-size: 1.1rem;
  font-weight: 600;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.3);
}

.action-icon {
  font-size: 2rem;
}

.recent-packs h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1rem;
}

.pack-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.pack-item {
  background: white;
  padding: 1.5rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.pack-status {
  padding: 0.5rem 1rem;
  border-radius: 2rem;
  font-size: 0.9rem;
  font-weight: 600;
  white-space: nowrap;
}

.pack-status.pending,
.pack-status.arrived {
  background: #e3f2fd;
  color: #1976d2;
}

.pack-status.checked_out {
  background: #e8f5e9;
  color: #388e3c;
}

.pack-status.in_transit,
.pack-status.shipped {
  background: #fff3e0;
  color: #f57c00;
}

.pack-status.cancelled {
  background: #ffebee;
  color: #d32f2f;
}

.pack-details {
  flex: 1;
}

.pack-details p {
  margin: 0.25rem 0;
  color: #666;
}

.pack-id {
  font-weight: 600;
  color: #333;
}

.pickup-code {
  font-size: 1.1rem;
  color: #667eea;
  font-weight: 600;
}

@media (max-width: 768px) {
  .action-buttons {
    flex-direction: column;
  }

  .pack-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
