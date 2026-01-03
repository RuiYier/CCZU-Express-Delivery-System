<template>
  <div class="admin-users">
    <h1>用户管理</h1>

    <div class="toolbar">
      <div class="search-box">
        <input v-model="searchQuery" type="text" placeholder="搜索用户（姓名、学号、手机号）" />
      </div>
      <button @click="fetchData" class="btn-refresh">🔄 刷新</button>
    </div>

    <div v-if="isLoading" class="loading">加载中...</div>

    <div v-else-if="filteredUsers.length === 0" class="empty">
      <div class="empty-icon">👥</div>
      <p>暂无用户数据</p>
    </div>

    <div v-else class="table-container">
      <table class="users-table">
        <thead>
          <tr>
            <th>用户ID</th>
            <th>姓名</th>
            <th>学号</th>
            <th>手机号</th>
            <th>地址</th>
            <th>角色</th>
            <th>注册时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.user_id">
            <td>{{ user.user_id }}</td>
            <td class="user-name">{{ user.user_name }}</td>
            <td class="student-id">{{ user.student_id }}</td>
            <td>{{ user.phone }}</td>
            <td class="address">{{ user.address }}</td>
            <td>
              <span class="role-badge" :class="user.role">
                {{ user.role === 'admin' ? '管理员' : '用户' }}
              </span>
            </td>
            <td>{{ formatTime(user.created_at) }}</td>
            <td>
              <button 
                v-if="user.user_id !== authStore.user?.user_id" 
                @click="handleDelete(user)" 
                class="btn-delete"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="stats">
      <p>总用户数: <strong>{{ users.length }}</strong></p>
      <p>管理员: <strong>{{ adminCount }}</strong></p>
      <p>普通用户: <strong>{{ userCount }}</strong></p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { adminApi } from '@/api'
import type { User } from '@/types'

const authStore = useAuthStore()
const users = ref<User[]>([])
const isLoading = ref(false)
const searchQuery = ref('')

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value

  const query = searchQuery.value.toLowerCase()
  return users.value.filter(
    (user) =>
      user.user_name.toLowerCase().includes(query) ||
      user.student_id.toLowerCase().includes(query) ||
      user.phone.includes(query)
  )
})

const adminCount = computed(() => users.value.filter((u) => u.role === 'admin').length)
const userCount = computed(() => users.value.filter((u) => u.role === 'user').length)

const formatTime = (time?: string) => {
  if (!time) return '-'
  const date = new Date(time)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const handleDelete = async (user: User) => {
  if (!confirm(`确定要删除用户 "${user.user_name}" 吗？此操作不可恢复。`)) return

  try {
    await adminApi.deleteUser(user.user_id)
    await fetchData()
    alert('用户删除成功')
  } catch (error: any) {
    console.error('删除失败:', error)
    alert(error.response?.data?.message || '删除失败，请重试')
  }
}

const fetchData = async () => {
  isLoading.value = true
  try {
    const response = await adminApi.getAllUsers()
    users.value = response.data.users || []
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.admin-users {
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

.search-box {
  flex: 1;
  max-width: 400px;
}

.search-box input {
  width: 100%;
  padding: 0.875rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 0.5rem;
  font-size: 1rem;
  transition: border-color 0.3s ease;
}

.search-box input:focus {
  outline: none;
  border-color: #f5576c;
}

.btn-refresh {
  padding: 0.875rem 1.5rem;
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
  margin-bottom: 2rem;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}

.users-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #666;
  font-size: 0.9rem;
  text-transform: uppercase;
}

.users-table tbody tr:hover {
  background: #f8f9fa;
}

.user-name {
  font-weight: 600;
  color: #333;
}

.student-id {
  font-family: monospace;
  color: #666;
}

.address {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.role-badge {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 2rem;
  font-size: 0.85rem;
  font-weight: 600;
}

.role-badge.admin {
  background: #ffebee;
  color: #d32f2f;
}

.role-badge.user {
  background: #e3f2fd;
  color: #1976d2;
}

.btn-delete {
  padding: 0.4rem 0.8rem;
  background: #ff4d4f;
  color: white;
  border: none;
  border-radius: 0.3rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.3s;
}

.btn-delete:hover {
  background: #d9363e;
}

.stats {
  background: white;
  padding: 1.5rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}

.stats p {
  margin: 0;
  color: #666;
}

.stats strong {
  color: #f5576c;
  font-size: 1.2rem;
  margin-left: 0.5rem;
}

@media (max-width: 768px) {
  .table-container {
    overflow-x: auto;
  }

  .users-table {
    min-width: 900px;
  }

  .stats {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
