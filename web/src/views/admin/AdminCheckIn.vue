<template>
  <div class="admin-checkin">
    <h1>包裹入库</h1>

    <div class="checkin-card">
      <h2>扫描或输入包裹信息</h2>

      <form @submit.prevent="handleCheckIn" class="checkin-form">
        <div class="form-row">
          <div class="form-group">
            <label for="pack_id">快递单号</label>
            <input
              id="pack_id"
              v-model.number="formData.pack_id"
              type="number"
              placeholder="请输入或扫描快递单号"
              required
              autofocus
            />
          </div>

          <div class="form-group">
            <label for="shelf_code">货架号</label>
            <input
              id="shelf_code"
              v-model.number="formData.shelf_code"
              type="number"
              placeholder="请输入货架号"
              required
            />
          </div>
        </div>

        <div class="form-group">
          <label for="user_search">收件人</label>
          <input
            id="user_search"
            v-model="userSearch"
            type="text"
            placeholder="搜索用户（姓名、学号、手机号）"
            @input="searchUsers"
          />

          <div v-if="searchResults.length > 0" class="search-results">
            <div
              v-for="user in searchResults"
              :key="user.user_id"
              @click="selectUser(user)"
              class="search-result-item"
            >
              <div class="user-info">
                <strong>{{ user.user_name }}</strong>
                <span>{{ user.student_id }}</span>
                <span>{{ user.phone }}</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="selectedUser" class="selected-user">
          <h4>已选择收件人:</h4>
          <div class="user-details">
            <p><strong>姓名:</strong> {{ selectedUser.user_name }}</p>
            <p><strong>学号:</strong> {{ selectedUser.student_id }}</p>
            <p><strong>手机:</strong> {{ selectedUser.phone }}</p>
            <p><strong>地址:</strong> {{ selectedUser.address }}</p>
          </div>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div class="form-actions">
          <button type="submit" class="btn-submit" :disabled="isLoading || !selectedUser">
            {{ isLoading ? '入库中...' : '确认入库' }}
          </button>
          <button type="button" @click="resetForm" class="btn-reset">重置</button>
        </div>
      </form>
    </div>

    <div v-if="successMessage" class="success-card">
      <div class="success-icon">✅</div>
      <h3>入库成功！</h3>
      <div class="pack-info">
        <p><strong>包裹号:</strong> {{ lastPack?.pack_id }}</p>
        <p class="pickup-code">
          <strong>取件码:</strong> <span>{{ lastPack?.pickup_code }}</span>
        </p>
        <p><strong>货架号:</strong> {{ lastPack?.shelf_code }}</p>
      </div>
      <button @click="closeSuccess" class="btn-close">继续入库</button>
    </div>

    <div class="recent-checkins">
      <h2>最近入库记录</h2>
      <div v-if="recentPacks.length === 0" class="empty-recent">暂无记录</div>
      <div v-else class="recent-list">
        <div v-for="pack in recentPacks" :key="pack.pack_id" class="recent-item">
          <div class="pack-id">{{ pack.pack_id }}</div>
          <div class="pack-details">
            <span>取件码: {{ pack.pickup_code }}</span>
            <span>货架: {{ pack.shelf_code }}</span>
            <span>{{ formatTime(pack.check_in_time) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { adminApi, packApi } from '@/api'
import type { User, Pack } from '@/types'

const formData = reactive({
  pack_id: 0,
  user_id: 0,
  shelf_code: 0
})

const userSearch = ref('')
const searchResults = ref<User[]>([])
const selectedUser = ref<User | null>(null)
const allUsers = ref<User[]>([])

const isLoading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref(false)
const lastPack = ref<Pack | null>(null)
const recentPacks = ref<Pack[]>([])

const searchUsers = () => {
  if (!userSearch.value) {
    searchResults.value = []
    return
  }

  const query = userSearch.value.toLowerCase()
  searchResults.value = allUsers.value
    .filter(
      (user) =>
        user.user_name.toLowerCase().includes(query) ||
        user.student_id.toLowerCase().includes(query) ||
        user.phone.includes(query)
    )
    .slice(0, 5)
}

const selectUser = (user: User) => {
  selectedUser.value = user
  formData.user_id = user.user_id
  userSearch.value = ''
  searchResults.value = []
}

const handleCheckIn = async () => {
  if (!selectedUser.value) {
    error.value = '请选择收件人'
    return
  }

  error.value = null
  isLoading.value = true

  try {
    const response = await packApi.checkIn(formData)
    lastPack.value = response.data.pack || null
    successMessage.value = true

    // 添加到最近记录
    if (lastPack.value) {
      recentPacks.value.unshift(lastPack.value)
      if (recentPacks.value.length > 10) {
        recentPacks.value.pop()
      }
    }

    resetForm()
  } catch (err: any) {
    error.value = err.response?.data?.message || '入库失败，请检查信息'
  } finally {
    isLoading.value = false
  }
}

const resetForm = () => {
  formData.pack_id = 0
  formData.user_id = 0
  formData.shelf_code = 0
  userSearch.value = ''
  searchResults.value = []
  selectedUser.value = null
  error.value = null
}

const closeSuccess = () => {
  successMessage.value = false
  lastPack.value = null
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

onMounted(async () => {
  try {
    const response = await adminApi.getAllUsers()
    allUsers.value = response.data.users || []
  } catch (error) {
    console.error('获取用户列表失败:', error)
  }
})
</script>

<style scoped>
.admin-checkin {
  max-width: 900px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.checkin-card {
  background: white;
  padding: 2.5rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.checkin-card h2 {
  font-size: 1.5rem;
  color: #333;
  margin: 0 0 2rem;
}

.checkin-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  position: relative;
}

.form-group label {
  font-weight: 600;
  color: #333;
  font-size: 0.95rem;
}

.form-group input {
  padding: 0.875rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #f5576c;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 2px solid #f5576c;
  border-radius: 0.5rem;
  max-height: 200px;
  overflow-y: auto;
  z-index: 10;
  margin-top: 0.5rem;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.search-result-item {
  padding: 1rem;
  cursor: pointer;
  transition: background 0.2s ease;
  border-bottom: 1px solid #f0f0f0;
}

.search-result-item:hover {
  background: #f8f9fa;
}

.search-result-item:last-child {
  border-bottom: none;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.user-info strong {
  color: #333;
}

.user-info span {
  color: #666;
  font-size: 0.9rem;
}

.selected-user {
  background: #e8f5e9;
  padding: 1.5rem;
  border-radius: 0.5rem;
  border: 2px solid #4caf50;
}

.selected-user h4 {
  margin: 0 0 1rem;
  color: #2e7d32;
}

.user-details {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem;
}

.user-details p {
  margin: 0;
  color: #333;
  font-size: 0.95rem;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 0.875rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  border: 1px solid #fcc;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
}

.btn-submit,
.btn-reset {
  flex: 1;
  padding: 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-submit {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(245, 87, 108, 0.4);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-reset {
  background: #f5f5f5;
  color: #666;
}

.btn-reset:hover {
  background: #e0e0e0;
}

.success-card {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: white;
  padding: 3rem;
  border-radius: 1rem;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.3);
  text-align: center;
  z-index: 1000;
  min-width: 400px;
  animation: slideIn 0.3s ease-out;
}

.success-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
}

.success-card h3 {
  font-size: 1.5rem;
  color: #333;
  margin: 0 0 1.5rem;
}

.pack-info {
  text-align: left;
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 0.5rem;
  margin-bottom: 2rem;
}

.pack-info p {
  margin: 0.75rem 0;
  color: #666;
}

.pickup-code {
  font-size: 1.2rem;
}

.pickup-code span {
  color: #f5576c;
  font-weight: 700;
  font-size: 1.5rem;
}

.btn-close {
  padding: 0.875rem 2rem;
  background: #f5576c;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s ease;
}

.btn-close:hover {
  background: #d94558;
}

.recent-checkins {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.recent-checkins h2 {
  font-size: 1.3rem;
  color: #333;
  margin: 0 0 1.5rem;
}

.empty-recent {
  text-align: center;
  padding: 2rem;
  color: #999;
}

.recent-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.recent-item {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 0.5rem;
}

.pack-id {
  font-weight: 700;
  color: #f5576c;
  font-family: monospace;
  font-size: 1.1rem;
}

.pack-details {
  display: flex;
  gap: 1.5rem;
  flex: 1;
  color: #666;
  font-size: 0.9rem;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translate(-50%, -60%);
  }
  to {
    opacity: 1;
    transform: translate(-50%, -50%);
  }
}

@media (max-width: 768px) {
  .form-row,
  .user-details {
    grid-template-columns: 1fr;
  }

  .success-card {
    min-width: auto;
    width: 90%;
  }

  .pack-details {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>
