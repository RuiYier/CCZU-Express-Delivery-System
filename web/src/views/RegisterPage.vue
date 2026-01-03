<template>
  <div class="register-page">
    <div class="register-container">
      <div class="register-card">
        <div class="register-header">
          <div class="logo">
            <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
          </div>
          <h2>创建账号</h2>
          <p>注册您的校园快递系统账号</p>
        </div>

        <form @submit.prevent="handleRegister" class="register-form">
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
              v-model="formData.student_id"
              type="text"
              placeholder="请输入学号"
              required
            />
          </div>

          <div class="form-group">
            <label for="phone">手机号</label>
            <input
              id="phone"
              v-model="formData.phone"
              type="tel"
              placeholder="请输入手机号"
              required
            />
          </div>

          <div class="form-group">
            <label for="address">宿舍地址</label>
            <input
              id="address"
              v-model="formData.address"
              type="text"
              placeholder="例如：博雅居10幢548室"
              required
            />
          </div>

          <div class="form-group">
            <label for="password">密码</label>
            <input
              id="password"
              v-model="formData.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              required
              minlength="6"
            />
          </div>

          <div class="form-group">
            <label for="confirm_password">确认密码</label>
            <input
              id="confirm_password"
              v-model="confirmPassword"
              type="password"
              placeholder="请再次输入密码"
              required
            />
          </div>

          <div class="form-group role-group">
            <label>账号类型</label>
            <div class="radio-group">
              <label class="radio-label">
                <input type="radio" v-model="formData.role" value="user" />
                <span>普通用户</span>
              </label>
              <label class="radio-label">
                <input type="radio" v-model="formData.role" value="admin" />
                <span>管理员</span>
              </label>
            </div>
          </div>

          <div v-if="error" class="error-message">{{ error }}</div>

          <button type="submit" class="btn-submit" :disabled="isLoading">
            {{ isLoading ? '注册中...' : '注册' }}
          </button>

          <div class="form-footer">
            <p>已有账号？<router-link to="/login">立即登录</router-link></p>
            <router-link to="/" class="back-link">返回首页</router-link>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const formData = reactive({
  user_name: '',
  student_id: '',
  phone: '',
  address: '',
  password: '',
  role: 'user' as 'user' | 'admin'
})

const confirmPassword = ref('')
const isLoading = ref(false)
const error = ref<string | null>(null)

const handleRegister = async () => {
  error.value = null

  // 验证密码
  if (formData.password !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }

  if (formData.password.length < 6) {
    error.value = '密码至少需要6位'
    return
  }

  isLoading.value = true

  try {
    const user = await authStore.register(formData)

    // 根据用户角色跳转
    if (user.role === 'admin') {
      router.push({ name: 'admin-dashboard' })
    } else {
      router.push({ name: 'user-dashboard' })
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || '注册失败，请检查信息'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background-image: url("/resource/background.png");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.register-container {
  width: 100%;
  max-width: 500px;
}

.register-card {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(10px);
  border-radius: 1.5rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  padding: 3rem;
  animation: slideIn 0.5s ease-out;
  max-height: 90vh;
  overflow-y: auto;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.register-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo {
  margin-bottom: 1rem;
  color: #42b983;
  animation: bounce 2s infinite;
}

.register-header h2 {
  font-size: 2rem;
  color: #333;
  margin: 0 0 0.5rem;
}

.register-header p {
  color: #666;
  margin: 0;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
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
  border-color: #f5576c;
}

.role-group .radio-group {
  display: flex;
  gap: 1.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-weight: normal;
}

.radio-label input[type='radio'] {
  width: auto;
  cursor: pointer;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 0.875rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  border: 1px solid #fcc;
}

.btn-submit {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  padding: 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  margin-top: 0.5rem;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(245, 87, 108, 0.4);
}

.btn-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-footer {
  text-align: center;
  margin-top: 1rem;
}

.form-footer p {
  color: #666;
  margin: 0 0 1rem;
}

.form-footer a {
  color: #f5576c;
  text-decoration: none;
  font-weight: 600;
}

.form-footer a:hover {
  text-decoration: underline;
}

.back-link {
  display: inline-block;
  color: #999;
  font-size: 0.9rem;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes bounce {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

@media (max-width: 480px) {
  .register-card {
    padding: 2rem;
  }

  .register-header h2 {
    font-size: 1.5rem;
  }
}
</style>
