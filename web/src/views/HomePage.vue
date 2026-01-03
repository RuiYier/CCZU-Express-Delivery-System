<template>
  <div class="home-page">
    <!-- 顶部导航 -->
    <header class="header">
      <div class="auth-nav">
        <router-link to="/login" class="nav-link">登录</router-link>
        <span class="divider">/</span>
        <router-link to="/register" class="nav-link">注册</router-link>
      </div>
    </header>

    <!-- 主要内容区 -->
    <main class="main-content">
      <div class="hero-text">
        <h1 class="title">常大校园快递收发系统</h1>
        <p class="subtitle">CCZU Express Delivery System</p>
      </div>

      <!-- 功能展示 - 单行排列，无边框 -->
      <div class="features-row">
        <div class="feature-item">
          <div class="icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path><polyline points="3.27 6.96 12 12.01 20.73 6.96"></polyline><line x1="12" y1="22.08" x2="12" y2="12"></line></svg>
          </div>
          <h3>快递入库</h3>
          <p>智能扫码，快速入库</p>
        </div>

        <div class="feature-item">
          <div class="icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
          </div>
          <h3>在线寄件</h3>
          <p>便捷下单，上门取件</p>
        </div>

        <div class="feature-item">
          <div class="icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          </div>
          <h3>实时追踪</h3>
          <p>物流状态，随时掌握</p>
        </div>

        <div class="feature-item">
          <div class="icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect><line x1="8" y1="21" x2="16" y2="21"></line><line x1="12" y1="17" x2="12" y2="21"></line></svg>
          </div>
          <h3>数据管理</h3>
          <p>高效统计，智慧运营</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

onMounted(() => {
  // 如果已登录，重定向到对应界面
  if (authStore.isAuthenticated) {
    if (authStore.isAdmin) {
      router.push({ name: 'admin-dashboard' })
    } else {
      router.push({ name: 'user-dashboard' })
    }
  }
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  width: 100%;
  position: relative;
  background-image: url("/resource/background.png");
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  color: #2c3e50;
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  overflow: hidden;
}

/* 顶部导航 */
.header {
  position: absolute;
  top: 2rem;
  right: 3rem;
  z-index: 10;
}

.auth-nav {
  display: flex;
  align-items: center;
  gap: 1rem;
  font-size: 1.1rem;
  background: rgba(255, 255, 255, 0.6);
  padding: 0.5rem 1.5rem;
  border-radius: 2rem;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.nav-link {
  color: #2c3e50;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.3s;
}

.nav-link:hover {
  color: #42b983;
}

.divider {
  color: #cbd5e0;
}

/* 主要内容 */
.main-content {
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding-bottom: 5vh;
}

.hero-text {
  text-align: center;
  margin-bottom: 8rem;
  animation: fadeInDown 1s ease-out;
}

.title {
  font-size: 4rem;
  font-weight: 800;
  margin: 0;
  letter-spacing: 2px;
  background: linear-gradient(to right, #000000 0%,#242424 30%,#707070 45%,#494949 60%,#272727 75%,#000000 100%);
  background-size: 200% 100%;
  color: transparent;
  -webkit-background-clip: text;
  background-clip: text;
  animation: shine 10s linear infinite;
  text-shadow: 
  0 2px 2px rgba(68, 68, 68, 0.2);
}


.subtitle {
  font-size: 1.2rem;
  margin-top: 1rem;
  color: #ffffffa9;
  font-weight: 400;
  letter-spacing: 4px;
  text-transform: uppercase;
}

/* 功能展示行 */
.features-row {
  display: flex;
  justify-content: center;
  gap: 2rem;
  width: 100%;
  max-width: 1200px;
  padding: 0 2rem;
  animation: fadeInUp 1s ease-out 0.3s backwards;
}

.feature-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  color: #2d3748;
  
  /* 卡片样式 */
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  padding: 1.5rem 1.5rem;
  border-radius: 1rem;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.feature-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 15px rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.8);
}

.icon-wrapper {
  margin-bottom: 1.5rem;
  color: #2c3e50;
  transition: transform 0.3s ease;
}

.feature-item:hover .icon-wrapper {
  transform: scale(1.1);
  color: #42b983;
}

.feature-item h3 {
  font-size: 1.2rem;
  margin-bottom: 0.8rem;
  font-weight: 700;
  color: #1a202c;
}

.feature-item p {
  font-size: 0.95rem;
  color: #4a5568;
  line-height: 1.6;
}

/* 动画 */
@keyframes fadeInDown {
  from { opacity: 0; transform: translateY(-30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes shine {
  to {
    background-position: 200% center;
  }
}

/* 响应式 */
@media (max-width: 768px) {
  .features-row {
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .title {
    font-size: 2.5rem;
  }
  
  .header {
    top: 1rem;
    right: 1rem;
  }
}
</style>
