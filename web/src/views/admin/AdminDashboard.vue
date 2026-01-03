<template>
  <div class="admin-dashboard">
    <h1>管理控制台</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-info">
          <h3>总用户数</h3>
          <p class="stat-value">{{ totalUsers }}</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-info">
          <h3>总包裹数</h3>
          <p class="stat-value">{{ totalPacks }}</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📥</div>
        <div class="stat-info">
          <h3>待取包裹</h3>
          <p class="stat-value">{{ pendingPacks }}</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📤</div>
        <div class="stat-info">
          <h3>运输中</h3>
          <p class="stat-value">{{ inTransitPacks }}</p>
        </div>
      </div>
    </div>

    <div class="quick-actions">
      <h2>快捷操作</h2>
      <div class="action-buttons">
        <router-link to="/admin/checkin" class="action-btn">
          <span class="action-icon">📥</span>
          <span>包裹入库</span>
        </router-link>
        <router-link to="/admin/packs" class="action-btn">
          <span class="action-icon">📦</span>
          <span>包裹管理</span>
        </router-link>
        <router-link to="/admin/users" class="action-btn">
          <span class="action-icon">👥</span>
          <span>用户管理</span>
        </router-link>
      </div>
    </div>

    <div class="overview">
      <h2>系统概览</h2>
      <div class="overview-grid">
        <div class="overview-item">
          <h4>包裹状态分布</h4>
          <div class="status-list">
            <div class="status-item">
              <span class="status-label">待取</span>
              <span class="status-count">{{ pendingPacks }}</span>
            </div>
            <div class="status-item">
              <span class="status-label">已取件</span>
              <span class="status-count">{{ checkedOutPacks }}</span>
            </div>
            <div class="status-item">
              <span class="status-label">运输中</span>
              <span class="status-count">{{ inTransitPacks }}</span>
            </div>
          </div>
        </div>

        <div class="overview-item">
          <h4>系统状态</h4>
          <div class="system-status">
            <div class="status-item">
              <span class="status-label">系统运行</span>
              <span class="status-badge online">正常</span>
            </div>
            <div class="status-item">
              <span class="status-label">数据库</span>
              <span class="status-badge online">连接中</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { adminApi } from '@/api'
import type { User, Pack } from '@/types'

const users = ref<User[]>([])
const packs = ref<Pack[]>([])
const isLoading = ref(false)

const totalUsers = computed(() => users.value.length)
const totalPacks = computed(() => packs.value.length)
const pendingPacks = computed(
  () => packs.value.filter((p) => p.pack_status === 'pending' || p.pack_status === 'arrived').length
)
const checkedOutPacks = computed(() => packs.value.filter((p) => p.pack_status === 'checked_out').length)
const inTransitPacks = computed(
  () => packs.value.filter((p) => p.pack_status === 'in_transit' || p.pack_status === 'shipped').length
)

const fetchData = async () => {
  isLoading.value = true
  try {
    const [usersRes, packsRes] = await Promise.all([adminApi.getAllUsers(), adminApi.getAllPacks()])
    users.value = usersRes.data.users || []
    packs.value = packsRes.data.packs || []
  } catch (error) {
    console.error('获取数据失败:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-dashboard {
  max-width: 1400px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
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
  color: #f5576c;
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
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-btn {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
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
  box-shadow: 0 6px 20px rgba(245, 87, 108, 0.3);
}

.action-icon {
  font-size: 2rem;
}

.overview {
  margin-bottom: 3rem;
}

.overview h2 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1rem;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.overview-item {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.overview-item h4 {
  margin: 0 0 1.5rem;
  font-size: 1.1rem;
  color: #333;
}

.status-list,
.system-status {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 0;
  border-bottom: 1px solid #f0f0f0;
}

.status-item:last-child {
  border-bottom: none;
}

.status-label {
  color: #666;
  font-size: 0.95rem;
}

.status-count {
  font-size: 1.3rem;
  font-weight: 700;
  color: #f5576c;
}

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 2rem;
  font-size: 0.85rem;
  font-weight: 600;
}

.status-badge.online {
  background: #e8f5e9;
  color: #388e3c;
}

@media (max-width: 768px) {
  .action-buttons {
    grid-template-columns: 1fr;
  }
}
</style>
