// 用户相关类型
export interface User {
  user_id: number
  user_name: string
  student_id: string
  phone: string
  address: string
  role: 'user' | 'admin'
  created_at?: string
  updated_at?: string
}

// 登录请求
export interface LoginRequest {
  student_id: string
  password: string
}

// 注册请求
export interface RegisterRequest {
  user_name: string
  password: string
  student_id: string
  phone: string
  address: string
  role?: 'user' | 'admin'
}

// 认证响应
export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
}

// 包裹状态类型
export type PackStatus = 'pending' | 'checked_out' | 'in_transit' | 'cancelled' | 'arrived' | 'shipped'

// 包裹信息
export interface Pack {
  pack_id: number
  user_id: number
  pack_status: PackStatus
  pickup_code?: string
  shelf_code?: number
  check_in_time?: string
  check_out_time?: string
  shipping_address?: string
  recipient?: string
  reciving_address?: string
  shipper_phone?: string
  recipient_phone?: string
  created_at?: string
  updated_at?: string
}

// 包裹入库请求
export interface PackCheckInRequest {
  pack_id: number
  user_id: number
  shelf_code: number
}

// 包裹出库请求
export interface PackCheckOutRequest {
  pack_id: number
  user_id: number
}

// 寄件请求
export interface MailPackRequest {
  shipping_address: string
  recipient: string
  reciving_address: string
  shipper_phone: string
  recipient_phone: string
}

// 取消寄件请求
export interface CancelMailRequest {
  pack_id: number
  user_id: number
}

// 更新包裹状态请求
export interface UpdatePackStatusRequest {
  pack_id: number
  pack_status: PackStatus
}

// 更新用户信息请求
export interface UpdateUserInfoRequest {
  user_id: number
  user_name?: string
  address?: string
  phone?: string
}

// API响应通用格式
export interface ApiResponse<T = any> {
  message?: string
  data?: T
  user?: User
  pack?: Pack
  packs?: Pack[]
  users?: User[]
  access_token?: string
  refresh_token?: string
}
