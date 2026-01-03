import apiClient from './client'
import type {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  User,
  Pack,
  PackCheckInRequest,
  PackCheckOutRequest,
  MailPackRequest,
  CancelMailRequest,
  UpdatePackStatusRequest,
  UpdateUserInfoRequest,
  ApiResponse
} from '@/types'

// ============ 认证相关 ============
export const authApi = {
  // 用户登录
  login: (data: LoginRequest) => 
    apiClient.post<AuthResponse>('/login', data),

  // 用户注册
  register: (data: RegisterRequest) => 
    apiClient.post<AuthResponse>('/register', data),

  // 健康检查
  ping: () => 
    apiClient.get('/ping')
}

// ============ 包裹相关 ============
export const packApi = {
  // 包裹入库
  checkIn: (data: PackCheckInRequest) => 
    apiClient.post<ApiResponse>('/packCheckIn', data),

  // 包裹出库
  checkOut: (data: PackCheckOutRequest) => 
    apiClient.post<ApiResponse>('/packCheckout', data),

  // 创建寄件
  mailPack: (data: MailPackRequest) => 
    apiClient.post<ApiResponse>('/mailPack', data),

  // 取消寄件
  cancelMail: (data: CancelMailRequest) => 
    apiClient.post<ApiResponse>('/cancelMail', data),

  // 更新包裹状态
  updateStatus: (data: UpdatePackStatusRequest) => 
    apiClient.post<ApiResponse>('/updatePackStatus', data),

  // 获取包裹详情
  getDetails: (packId: number) => 
    apiClient.get<ApiResponse>(`/getPackDetails/${packId}`),

  // 获取用户所有包裹
  getAllByUser: (userId: number) => 
    apiClient.get<ApiResponse>(`/allPacks/${userId}`)
}

// ============ 用户相关 ============
export const userApi = {
  // 更新用户信息
  updateInfo: (data: UpdateUserInfoRequest) => 
    apiClient.post<ApiResponse>('/updateUserInfo', data)
}

// ============ 管理员相关 ============
export const adminApi = {
  // 获取所有用户
  getAllUsers: () => 
    apiClient.get<ApiResponse>('/admin/users'),

  // 获取所有包裹
  getAllPacks: (status?: string) => 
    apiClient.get<ApiResponse>('/admin/packs', { params: { status } }),

  // 更新包裹信息
  updatePack: (data: Partial<Pack>) => 
    apiClient.put<ApiResponse>('/admin/pack', data),

  // 删除用户
  deleteUser: (userId: number) => 
    apiClient.delete<ApiResponse>('/admin/deleteUser', { params: { user_id: userId } })
}
