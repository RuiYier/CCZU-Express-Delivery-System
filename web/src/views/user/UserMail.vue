<template>
  <div class="user-mail">
    <h1>寄件服务</h1>

    <div class="mail-form-card">
      <h2>填写寄件信息</h2>

      <form @submit.prevent="handleSubmit" class="mail-form">
        <div class="form-row">
          <div class="form-group">
            <label for="shipper_phone">寄件人手机号</label>
            <input
              id="shipper_phone"
              v-model="formData.shipper_phone"
              type="tel"
              placeholder="请输入寄件人手机号"
              required
            />
            <small>{{ authStore.user?.phone }}</small>
          </div>

          <div class="form-group">
            <label for="shipping_address">寄件地址</label>
            <input
              id="shipping_address"
              v-model="formData.shipping_address"
              type="text"
              placeholder="请输入寄件地址"
              required
            />
          </div>
        </div>

        <div class="divider">
          <span>收件人信息</span>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="recipient">收件人姓名</label>
            <input
              id="recipient"
              v-model="formData.recipient"
              type="text"
              placeholder="请输入收件人姓名"
              required
            />
          </div>

          <div class="form-group">
            <label for="recipient_phone">收件人手机号</label>
            <input
              id="recipient_phone"
              v-model="formData.recipient_phone"
              type="tel"
              placeholder="必须是系统内已注册用户"
              required
            />
            <small>收件人必须是系统内已注册用户</small>
          </div>
        </div>

        <div class="form-group">
          <label for="reciving_address">收件地址</label>
          <input
            id="reciving_address"
            v-model="formData.reciving_address"
            type="text"
            placeholder="请输入收件地址"
            required
          />
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>

        <div class="form-actions">
          <button type="submit" class="btn-submit" :disabled="isLoading">
            {{ isLoading ? '提交中...' : '提交寄件' }}
          </button>
          <button type="button" @click="resetForm" class="btn-reset">重置</button>
        </div>
      </form>
    </div>

    <div v-if="successMessage" class="success-card">
      <div class="success-icon">✅</div>
      <h3>寄件成功！</h3>
      <p>{{ successMessage }}</p>
      <button @click="successMessage = ''" class="btn-close">知道了</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { usePackStore } from '@/stores/counter'

const authStore = useAuthStore()
const packStore = usePackStore()

const formData = reactive({
  shipping_address: '',
  recipient: '',
  reciving_address: '',
  shipper_phone: '',
  recipient_phone: ''
})

const isLoading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref('')

const handleSubmit = async () => {
  error.value = null
  isLoading.value = true

  try {
    const result = await packStore.createMailPack(formData)
    if (result) {
      successMessage.value = `寄件单号：${result.pack_id}`
      resetForm()
    }
  } catch (err: any) {
    error.value = err.response?.data?.message || '寄件失败，请检查信息'
  } finally {
    isLoading.value = false
  }
}

const resetForm = () => {
  formData.shipping_address = authStore.user?.address || ''
  formData.recipient = ''
  formData.reciving_address = ''
  formData.shipper_phone = authStore.user?.phone || ''
  formData.recipient_phone = ''
  error.value = null
}

onMounted(() => {
  // 预填充用户信息
  formData.shipper_phone = authStore.user?.phone || ''
  formData.shipping_address = authStore.user?.address || ''
})
</script>

<style scoped>
.user-mail {
  max-width: 800px;
}

h1 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 2rem;
}

.mail-form-card {
  background: white;
  padding: 2.5rem;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.mail-form-card h2 {
  font-size: 1.5rem;
  color: #333;
  margin: 0 0 2rem;
}

.mail-form {
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

.form-group small {
  color: #999;
  font-size: 0.85rem;
}

.divider {
  text-align: center;
  position: relative;
  margin: 1rem 0;
}

.divider::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  width: 100%;
  height: 1px;
  background: #e0e0e0;
}

.divider span {
  position: relative;
  background: white;
  padding: 0 1rem;
  color: #999;
  font-size: 0.9rem;
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
  animation: slideIn 0.3s ease-out;
}

.success-icon {
  font-size: 5rem;
  margin-bottom: 1rem;
}

.success-card h3 {
  font-size: 1.5rem;
  color: #333;
  margin: 0 0 1rem;
}

.success-card p {
  color: #666;
  margin: 0 0 2rem;
  font-size: 1.1rem;
}

.btn-close {
  padding: 0.875rem 2rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.3s ease;
}

.btn-close:hover {
  background: #5568d3;
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
  .form-row {
    grid-template-columns: 1fr;
  }

  .mail-form-card {
    padding: 1.5rem;
  }
}
</style>
