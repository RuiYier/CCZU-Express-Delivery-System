<template>
  <div class="user-profile">
    <h1>个人信息</h1>

    <div class="profile-card">
      <div class="profile-header">
        <div class="avatar">{{ authStore.user?.user_name?.charAt(0) }}</div>
        <div class="user-basic">
          <h2>{{ authStore.user?.user_name }}</h2>
          <p class="role-badge">{{ authStore.user?.role === 'admin' ? '管理员' : '普通用户' }}</p>
        </div>
      </div>

      <form @submit.prevent="handleUpdate" class="profile-form">
        <div class="form-group">
          <label for="user_name">姓名</label>
          <input
            id="user_name"
            v-model="formData.user_name"
            type="text"
            placeholder="请输入姓名"
            required
          />
        </div>

        <div class="form-group">
          <label for="student_id">学号</label>
          <input
            id="student_id"
            :value="authStore.user?.student_id"
            type="text"
            disabled
            readonly
          />
          <small>学号不可修改</small>
        </div>

        <div class="form-group">
          <label for="phone">手机号</label>
          <input id="phone" v-model="formData.phone" type="tel" placeholder="请输入手机号" required />
        </div>

        <div class="form-group">
          <label for="address">宿舍地址</label>
          <input
            id="address"
            v-model="formData.address"
            type="text"
            placeholder="请输入宿舍地址"
            required
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="successMessage" class="success-message">{{ successMessage }}</div>

        <div class="form-actions">
          <button type="submit" class="btn-submit" :disabled="isLoading">
            {{ isLoading ? '保存中...' : '保存修改' }}
          </button>
          <button type="button" @click="resetForm" class="btn-reset">重置</button>
        </div>
      </form>
    </div>

    <div class="info-section">
      <h3>账号信息</h3>
      <div class="info-grid">
        <div class="info-item">
          <label>用户ID</label>
          <span>{{ authStore.user?.user_id }}</span>
        </div>
        <div class="info-item">
          <label>注册时间</label>
          <span>{{ formatTime(authStore.user?.created_at) }}</span>
        </div>
        <div class="info-item">
          <label>更新时间</label>
          <span>{{ formatTime(authStore.user?.updated_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { userApi } from '@/api'

const authStore = useAuthStore()

const formData = reactive({
  user_name: '',
  phone: '',
  address: ''
})

const isLoading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

const handleUpdate = async () => {
  if (!authStore.user?.user_id) return

  error.value = null
  successMessage.value = null
  isLoading.value = true

  try {
    await userApi.updateInfo({
      user_id: authStore.user.user_id,
      ...formData
    })

    // 更新本地用户信息
    authStore.updateUser({
      ...authStore.user,
      ...formData
    })

    successMessage.value = '信息更新成功！'
  } catch (err: any) {
    error.value = err.response?.data?.message || '更新失败'
  } finally {
    isLoading.value = false
  }
}

const resetForm = () => {
  if (authStore.user) {
    formData.user_name = authStore.user.user_name
    formData.phone = authStore.user.phone
    formData.address = authStore.user.address
  }
  error.value = null
  successMessage.value = null
}

const formatTime = (time?: string) => {
  if (!time) return '未知'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  resetForm()
})
</script>

<style scoped>
.user-profile {
  max-width: 800px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.profile-card {
  background: white;
  padding: 2.5rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 2rem;
  margin-bottom: 2.5rem;
  padding-bottom: 2rem;
  border-bottom: 2px solid #f0f0f0;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 3rem;
  font-weight: 700;
}

.user-basic h2 {
  margin: 0 0 0.5rem;
  font-size: 1.8rem;
  color: #333;
}

.role-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  background: #e3f2fd;
  color: #1976d2;
  border-radius: 2rem;
  font-size: 0.9rem;
  font-weight: 600;
  margin: 0;
}

.profile-form {
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

.form-group input {
  padding: 0.875rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-group input:disabled {
  background: #f5f5f5;
  cursor: not-allowed;
}

.form-group small {
  color: #999;
  font-size: 0.85rem;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 0.875rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  border: 1px solid #fcc;
}

.success-message {
  background: #e8f5e9;
  color: #2e7d32;
  padding: 0.875rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  border: 1px solid #a5d6a7;
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
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

.info-section {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.info-section h3 {
  margin: 0 0 1.5rem;
  font-size: 1.3rem;
  color: #333;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item label {
  color: #999;
  font-size: 0.9rem;
}

.info-item span {
  color: #333;
  font-weight: 500;
}

@media (max-width: 768px) {
  .profile-header {
    flex-direction: column;
    text-align: center;
  }

  .profile-card {
    padding: 1.5rem;
  }
}
</style>
