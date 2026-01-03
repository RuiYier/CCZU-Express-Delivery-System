import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { Pack } from '@/types'
import { packApi } from '@/api'

export const usePackStore = defineStore('pack', () => {
  // 状态
  const packs = ref<Pack[]>([])
  const currentPack = ref<Pack | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // 获取用户所有包裹
  const fetchUserPacks = async (userId: number) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await packApi.getAllByUser(userId)
      packs.value = response.data.packs || []
      return packs.value
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取包裹列表失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 获取包裹详情
  const fetchPackDetails = async (packId: number) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await packApi.getDetails(packId)
      currentPack.value = response.data.pack || null
      return currentPack.value
    } catch (err: any) {
      error.value = err.response?.data?.message || '获取包裹详情失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 包裹出库
  const checkOutPack = async (packId: number, userId: number) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await packApi.checkOut({ pack_id: packId, user_id: userId })
      // 更新本地包裹状态
      const index = packs.value.findIndex(p => p.pack_id === packId)
      if (index !== -1 && response.data.pack) {
        packs.value[index] = response.data.pack
      }
      return response.data.pack
    } catch (err: any) {
      error.value = err.response?.data?.message || '取件失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 创建寄件
  const createMailPack = async (data: any) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await packApi.mailPack(data)
      const newPack = response.data.pack
      if (newPack) {
        packs.value.unshift(newPack)
      }
      return newPack
    } catch (err: any) {
      error.value = err.response?.data?.message || '寄件失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // 取消寄件
  const cancelMailPack = async (packId: number, userId: number) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await packApi.cancelMail({ pack_id: packId, user_id: userId })
      // 更新本地包裹状态
      const index = packs.value.findIndex(p => p.pack_id === packId)
      if (index !== -1 && packs.value[index]) {
        packs.value[index].pack_status = 'cancelled'
      }
      return response.data
    } catch (err: any) {
      error.value = err.response?.data?.message || '取消寄件失败'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  return {
    packs,
    currentPack,
    isLoading,
    error,
    fetchUserPacks,
    fetchPackDetails,
    checkOutPack,
    createMailPack,
    cancelMailPack
  }
})
