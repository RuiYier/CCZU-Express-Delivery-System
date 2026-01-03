<template>
  <div class="not-found">
    <div class="container">
      <div class="error-code">404</div>
      <h1>页面未找到</h1>
      <p>抱歉，您访问的页面不存在</p>
      <router-link to="/" class="btn-home">返回首页</router-link>
    </div>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.not-found {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.container {
  text-align: center;
  color: white;
}

.error-code {
  font-size: 10rem;
  font-weight: 700;
  margin: 0;
  line-height: 1;
  opacity: 0.5;
}

h1 {
  font-size: 2.5rem;
  margin: 1rem 0;
}

p {
  font-size: 1.2rem;
  margin: 1rem 0 2rem;
  opacity: 0.9;
}

.btn-home {
  display: inline-block;
  padding: 1rem 2rem;
  background: white;
  color: #667eea;
  text-decoration: none;
  border-radius: 2rem;
  font-weight: 600;
  transition: transform 0.2s ease;
}

.btn-home:hover {
  transform: translateY(-2px);
}
</style>
