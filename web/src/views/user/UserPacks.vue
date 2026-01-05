<template>
  <div class="user-packs">
    <h1>我的包裹</h1>

    <div class="filters">
      <button
        v-for="status in statusFilters"
        :key="status.value"
        @click="currentFilter = status.value"
        :class="{ active: currentFilter === status.value }"
        class="filter-btn"
      >
        {{ status.label }}
      </button>
    </div>

    <div v-if="packStore.isLoading" class="loading">加载中...</div>

    <div v-else-if="filteredPacks.length === 0" class="empty">
      <div class="empty-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-6l-2 3h-4l-2-3H2"></path><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"></path></svg>
      </div>
      <p>暂无包裹</p>
    </div>

    <div v-else class="pack-grid">
      <div v-for="pack in filteredPacks" :key="pack.pack_id" class="pack-card">
        <div class="pack-header">
          <span class="pack-status" :class="pack.pack_status">
            {{ getStatusText(pack.pack_status) }}
          </span>
          <span class="pack-id">{{ pack.pack_id }}</span>
        </div>

        <div class="pack-body">
          <div v-if="pack.pickup_code" class="pickup-code-section">
            <label>取件码</label>
            <div class="pickup-code">{{ pack.pickup_code }}</div>
          </div>

          <div v-if="pack.shelf_code" class="info-item">
            <label>货架号</label>
            <span>{{ pack.shelf_code }}</span>
          </div>

          <div v-if="pack.recipient" class="info-item">
            <label>收件人</label>
            <span>{{ pack.recipient }}</span>
          </div>

          <div v-if="pack.reciving_address" class="info-item">
            <label>收件地址</label>
            <span>{{ pack.reciving_address }}</span>
          </div>

          <div v-if="pack.check_in_time" class="info-item">
            <label>入库时间</label>
            <span>{{ formatTime(pack.check_in_time) }}</span>
          </div>

          <div v-if="pack.check_out_time" class="info-item">
            <label>取件时间</label>
            <span>{{ formatTime(pack.check_out_time) }}</span>
          </div>
        </div>

        <div class="pack-actions">
          <button
            v-if="pack.pack_status === 'pending' || pack.pack_status === 'arrived'"
            @click="handleCheckout(pack)"
            class="btn-checkout"
          >
            确认取件
          </button>
          <button
            v-if="pack.pack_status === 'in_transit'"
            @click="handleCancelMail(pack)"
            class="btn-cancel"
          >
            取消寄件
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usePackStore } from '@/stores/counter'
import type { Pack } from '@/types'

const authStore = useAuthStore()
const packStore = usePackStore()

const currentFilter = ref('all')

const statusFilters = [
  { label: '全部', value: 'all' },
  { label: '待取', value: 'pending' },
  { label: '已到达', value: 'arrived' },
  { label: '运输中', value: 'in_transit' },
  { label: '已取件', value: 'checked_out' }
]

const filteredPacks = computed(() => {
  if (currentFilter.value === 'all') {
    return packStore.packs
  }
  return packStore.packs.filter((p) => p.pack_status === currentFilter.value)
})

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
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleCheckout = async (pack: Pack) => {
  if (!authStore.user?.user_id) return

  if (!confirm(`确认取件 ${pack.pack_id} 吗？`)) return

  try {
    await packStore.checkOutPack(pack.pack_id, authStore.user.user_id)
    alert('取件成功！')
  } catch (error: any) {
    alert(error.response?.data?.message || '取件失败')
  }
}

const handleCancelMail = async (pack: Pack) => {
  if (!authStore.user?.user_id) return

  if (!confirm(`确认取消寄件 ${pack.pack_id} 吗？`)) return

  try {
    await packStore.cancelMailPack(pack.pack_id, authStore.user.user_id)
    alert('取消成功！')
  } catch (error: any) {
    alert(error.response?.data?.message || '取消失败')
  }
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
.user-packs {
  max-width: 1200px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.filters {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 0.75rem 1.5rem;
  border: 2px solid #e0e0e0;
  background: white;
  border-radius: 2rem;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-btn:hover {
  border-color: #667eea;
  color: #667eea;
}

.filter-btn.active {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.loading,
.empty {
  text-align: center;
  padding: 4rem 2rem;
  color: #999;
}

.empty-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
}

.pack-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.pack-card {
  background: white;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.pack-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.pack-header {
  padding: 1.5rem;
  background: #f8f9fa;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e0e0e0;
}

.pack-status {
  padding: 0.5rem 1rem;
  border-radius: 2rem;
  font-size: 0.85rem;
  font-weight: 600;
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

.pack-id {
  color: #999;
  font-size: 0.9rem;
}

.pack-body {
  padding: 1.5rem;
}

.pickup-code-section {
  text-align: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1.5rem;
  border-bottom: 2px dashed #e0e0e0;
}

.pickup-code-section label {
  display: block;
  color: #999;
  font-size: 0.85rem;
  margin-bottom: 0.5rem;
}

.pickup-code {
  font-size: 2rem;
  font-weight: 700;
  color: #667eea;
  letter-spacing: 0.1em;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
  font-size: 0.9rem;
}

.info-item label {
  color: #999;
  font-weight: 500;
}

.info-item span {
  color: #333;
}

.pack-actions {
  padding: 1rem 1.5rem;
  background: #f8f9fa;
  border-top: 1px solid #e0e0e0;
}

.pack-actions button {
  width: 100%;
  padding: 0.75rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-checkout {
  background: #667eea;
  color: white;
}

.btn-checkout:hover {
  background: #5568d3;
}

.btn-cancel {
  background: #f5576c;
  color: white;
}

.btn-cancel:hover {
  background: #d94558;
}

@media (max-width: 768px) {
  .pack-grid {
    grid-template-columns: 1fr;
  }
}
</style>
