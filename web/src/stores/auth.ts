import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { authApi } from '@/api'
import type { LoginRequest, RegisterRequest } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!accessToken.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isUser = computed(() => user.value?.role === 'user')

  // 从localStorage加载用户信息
  const loadFromStorage = () => {
    const storedUser = localStorage.getItem('user')
    const storedAccessToken = localStorage.getItem('access_token')
    const storedRefreshToken = localStorage.getItem('refresh_token')

    if (storedUser && storedAccessToken) {
      user.value = JSON.parse(storedUser)
      accessToken.value = storedAccessToken
      refreshToken.value = storedRefreshToken
    }
  }

  // 保存到localStorage
  const saveToStorage = () => {
    if (user.value && accessToken.value) {
      localStorage.setItem('user', JSON.stringify(user.value))
      localStorage.setItem('access_token', accessToken.value)
      if (refreshToken.value) {
        localStorage.setItem('refresh_token', refreshToken.value)
      }
    }
  }

  // 清除存储
  const clearStorage = () => {
    localStorage.removeItem('user')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  // 登录
  const login = async (credentials: LoginRequest) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.login(credentials)
      const { user: userData, access_token, refresh_token } = response.data

      user.value = userData
      accessToken.value = access_token
      refreshToken.value = refresh_token

      saveToStorage()
      return userData
    } catch (err: any) {
      error.value = err.response?.data?.message || '登录失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 注册
  const register = async (data: RegisterRequest) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await authApi.register(data)
      const { user: userData, access_token, refresh_token } = response.data

      user.value = userData
      accessToken.value = access_token
      refreshToken.value = refresh_token

      saveToStorage()
      return userData
    } catch (err: any) {
      error.value = err.response?.data?.message || '注册失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 登出
  const logout = () => {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    clearStorage()
  }

  // 更新用户信息
  const updateUser = (updatedUser: User) => {
    user.value = updatedUser
    saveToStorage()
  }

  // 初始化时加载
  loadFromStorage()

  return {
    user,
    accessToken,
    refreshToken,
    isLoading,
    error,
    isAuthenticated,
    isAdmin,
    isUser,
    login,
    register,
    logout,
    updateUser,
    loadFromStorage
  }
})
