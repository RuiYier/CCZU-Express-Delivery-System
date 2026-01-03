<template>
  <div class="admin-packs">
    <h1>包裹管理</h1>

    <div class="toolbar">
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

      <button @click="fetchData" class="btn-refresh">🔄 刷新</button>
    </div>

    <div v-if="isLoading" class="loading">加载中...</div>

    <div v-else-if="filteredPacks.length === 0" class="empty">
      <div class="empty-icon">📭</div>
      <p>暂无包裹数据</p>
    </div>

    <div v-else class="table-container">
      <table class="packs-table">
        <thead>
          <tr>
            <th>包裹ID</th>
            <th>用户ID</th>
            <th>状态</th>
            <th>取件码</th>
            <th>货架号</th>
            <th>入库时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pack in filteredPacks" :key="pack.pack_id">
            <td class="pack-id">{{ pack.pack_id }}</td>
            <td>{{ pack.user_id }}</td>
            <td>
              <span class="pack-status" :class="pack.pack_status">
                {{ getStatusText(pack.pack_status) }}
              </span>
            </td>
            <td class="pickup-code">{{ pack.pickup_code || '-' }}</td>
            <td>{{ pack.shelf_code || '-' }}</td>
            <td>{{ formatTime(pack.check_in_time) }}</td>
            <td class="actions">
              <button @click="editPack(pack)" class="btn-edit">编辑</button>
              <button @click="updateStatus(pack)" class="btn-update">更新状态</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 编辑模态框 -->
    <div v-if="showEditModal" class="modal" @click.self="closeEditModal">
      <div class="modal-content">
        <h2>编辑包裹</h2>
        <form @submit.prevent="handleUpdate" class="edit-form">
          <div class="form-group">
            <label>包裹ID</label>
            <input v-model="editForm.pack_id" type="number" disabled />
          </div>

          <div class="form-group">
            <label>用户ID</label>
            <input v-model="editForm.user_id" type="number" required />
          </div>

          <div class="form-group">
            <label>状态</label>
            <select v-model="editForm.pack_status" required>
              <option value="pending">待取</option>
              <option value="arrived">已到达</option>
              <option value="checked_out">已取件</option>
              <option value="in_transit">运输中</option>
              <option value="shipped">已发货</option>
              <option value="cancelled">已取消</option>
            </select>
          </div>

          <div class="form-group">
            <label>取件码</label>
            <input v-model="editForm.pickup_code" type="text" />
          </div>

          <div class="form-group">
            <label>货架号</label>
            <input v-model="editForm.shelf_code" type="number" />
          </div>

          <div v-if="error" class="error-message">{{ error }}</div>

          <div class="modal-actions">
            <button type="submit" class="btn-submit" :disabled="isUpdating">
              {{ isUpdating ? '保存中...' : '保存' }}
            </button>
            <button type="button" @click="closeEditModal" class="btn-cancel">取消</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { adminApi } from '@/api'
import type { Pack } from '@/types'

const packs = ref<Pack[]>([])
const isLoading = ref(false)
const currentFilter = ref('all')
const showEditModal = ref(false)
const isUpdating = ref(false)
const error = ref<string | null>(null)

const editForm = ref({
  pack_id: 0,
  user_id: 0,
  pack_status: 'pending' as any,
  pickup_code: '',
  shelf_code: undefined as number | undefined
})

const statusFilters = [
  { label: '全部', value: 'all' },
  { label: '待取', value: 'pending' },
  { label: '已到达', value: 'arrived' },
  { label: '运输中', value: 'in_transit' },
  { label: '已取件', value: 'checked_out' }
]

const filteredPacks = computed(() => {
  if (currentFilter.value === 'all') {
    return packs.value
  }
  return packs.value.filter((p) => p.pack_status === currentFilter.value)
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
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getAllPacks()
    packs.value = response.data.packs || []
  } catch (error) {
    console.error('获取包裹列表失败:', error)
  } finally {
    isLoading.value = false
  }
}

const editPack = (pack: Pack) => {
  editForm.value = {
    pack_id: pack.pack_id,
    user_id: pack.user_id,
    pack_status: pack.pack_status,
    pickup_code: pack.pickup_code || '',
    shelf_code: pack.shelf_code
  }
  showEditModal.value = true
  error.value = null
}

const closeEditModal = () => {
  showEditModal.value = false
  error.value = null
}

const handleUpdate = async () => {
  isUpdating.value = true
  error.value = null

  try {
    await adminApi.updatePack(editForm.value)
    alert('更新成功！')
    closeEditModal()
    await fetchData()
  } catch (err: any) {
    error.value = err.response?.data?.message || '更新失败'
  } finally {
    isUpdating.value = false
  }
}

const updateStatus = async (pack: Pack) => {
  const newStatus = prompt(
    '请输入新状态 (pending/arrived/checked_out/in_transit/shipped/cancelled):',
    pack.pack_status
  )

  if (!newStatus) return

  try {
    await adminApi.updatePack({
      pack_id: pack.pack_id,
      pack_status: newStatus as any
    })
    alert('状态更新成功！')
    await fetchData()
  } catch (error: any) {
    alert(error.response?.data?.message || '更新失败')
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-packs {
  max-width: 1400px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1rem;
  flex-wrap: wrap;
}

.filters {
  display: flex;
  gap: 1rem;
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
  border-color: #f5576c;
  color: #f5576c;
}

.filter-btn.active {
  background: #f5576c;
  color: white;
  border-color: #f5576c;
}

.btn-refresh {
  padding: 0.75rem 1.5rem;
  background: #f5576c;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 0.95rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s ease;
}

.btn-refresh:hover {
  background: #d94558;
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

.table-container {
  background: white;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.packs-table {
  width: 100%;
  border-collapse: collapse;
}

.packs-table th,
.packs-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.packs-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #666;
  font-size: 0.9rem;
  text-transform: uppercase;
}

.packs-table tbody tr:hover {
  background: #f8f9fa;
}

.pack-id {
  font-family: monospace;
  color: #666;
}

.pickup-code {
  font-weight: 600;
  color: #f5576c;
}

.pack-status {
  display: inline-block;
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

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn-edit,
.btn-update {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-edit {
  background: #f5576c;
  color: white;
}

.btn-edit:hover {
  background: #d94558;
}

.btn-update {
  background: #667eea;
  color: white;
}

.btn-update:hover {
  background: #5568d3;
}

/* Modal */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 2.5rem;
  border-radius: 1rem;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content h2 {
  margin: 0 0 2rem;
  font-size: 1.5rem;
  color: #333;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: #333;
  font-size: 0.95rem;
}

.form-group input,
.form-group select {
  padding: 0.875rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #f5576c;
}

.form-group input:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 0.875rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  border: 1px solid #fcc;
}

.modal-actions {
  display: flex;
  gap: 1rem;
}

.btn-submit,
.btn-cancel {
  flex: 1;
  padding: 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-submit {
  background: #f5576c;
  color: white;
}

.btn-submit:hover:not(:disabled) {
  background: #d94558;
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-cancel {
  background: #f5f5f5;
  color: #666;
}

.btn-cancel:hover {
  background: #e0e0e0;
}

@media (max-width: 768px) {
  .table-container {
    overflow-x: auto;
  }

  .packs-table {
    min-width: 800px;
  }
}
</style>
