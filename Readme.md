# 常大校园快递收发系统 (CCZU Express Delivery System) - 学校课程实训项目

## 📖 项目概述

**常大校园快递收发系统** 是一个现代化校园快递管理平台，旨在解决校内快递收发效率低、信息不透明等问题。本项目为**学校实训项目**，采用前后端分离架构，模拟真实开发流程，涵盖了从需求分析、系统设计到编码实现的全过程。

- **后端**: 基于 Go (Gin) + PostgreSQL 开发，提供高性能的 API 服务。
- **前端**: 基于 Vue 3 + TypeScript + Vite 开发，提供现代化的用户交互体验。

## 🏗 系统架构

```
cczu_expressdelivery_system/
├── server/              # 后端服务 (Go)
│   ├── controllers/     # 业务逻辑控制器
│   ├── models/         # GORM 数据模型
│   ├── routes/         # 路由定义
│   ├── middlewares/    # 认证与权限中间件
│   ├── database/       # 数据库连接配置
│   ├── utils/          # 工具函数 (Snowflake, JWT等)
│   └── config/         # 配置文件
│
└── web/                # 前端应用 (Vue 3)
    ├── src/
    │   ├── api/        # Axios API 封装
    │   ├── stores/     # Pinia 状态管理
    │   ├── router/     # Vue Router 路由配置
    │   ├── views/      # 页面视图组件
    │   └── types/      # TypeScript 类型定义
    └── public/         # 静态资源 (背景图, SVG等)
```

## 🚀 快速启动

### 1. 启动后端服务

```bash
cd server
go mod tidy
go run .
```

后端服务默认运行在 `http://localhost:8088`

### 2. 启动前端应用

```bash
cd web
npm install
npm run dev
```

前端应用默认运行在 `http://localhost:7999`

## 📡 API对接说明

前端通过 Axios 与后端 API 通信，配置位于 `web/src/api/client.ts`。
开发环境下，Vite 配置了代理 (`vite.config.ts`) 转发 `/api` 请求至后端服务器。

##  认证机制

1. **JWT 认证**: 用户登录/注册后获取 Access Token 和 Refresh Token。
2. **持久化**: Token 存储在 localStorage 中。
3. **请求拦截**: 每次 API 请求自动在 Header 中携带 Token。
4. **安全验证**: 后端中间件验证 Token 合法性及过期时间。

##  用户角色

系统支持两种用户角色，拥有不同的权限视图：

### 普通用户 (user)
- **我的包裹**: 查看待取、已取包裹，支持状态筛选。
- **在线寄件**: 填写寄件信息，发起寄件请求。
- **个人中心**: 管理个人资料，查看寄件历史。

### 管理员 (admin)
- **控制台**: 查看系统数据概览（用户数、包裹数等）。
- **包裹管理**: 查看所有包裹，支持编辑、状态更新。
- **用户管理**: 查看所有用户，支持搜索、删除用户。
- **包裹入库**: 扫码/手动录入包裹，自动生成取件码。

##  界面设计

项目采用了现代化的 UI 设计风格：

- **首页**: 
  - 动态流光标题。
  - 玻璃拟态 (Glassmorphism) 风格的功能卡片。
  - 商务风格 SVG 图标，视觉统一。
  - 响应式布局，适配移动端。
- **登录/注册**: 
  - 与首页统一的背景与毛玻璃卡片风格。
  - 简洁的表单设计与交互反馈。
- **后台管理**:
  - 侧边栏导航，清晰的功能分区。
  - 统一的 SVG 图标系统。
  - 数据可视化展示。

##  技术栈

### 前端技术
- **Vue 3**
- **TypeScript**
- **Vite**
- **Pinia**
- **Vue Router 4**
- **CSS3**

### 后端技术
- **Go (Golang)**
- **Gin**
- **GORM**
- **PostgreSQL**
- **JWT**
- **Snowflake**

##  功能清单

### ✅ 已实现功能

**用户功能**
- 用户注册与登录
- 首页功能展示
- 包裹列表查看与筛选
- 包裹取件操作
- 在线寄件与取消
- 个人信息管理

**管理员功能**
- 数据看板
- 包裹入库（生成取件码）
- 包裹全量管理（编辑、状态流转）
- 用户列表管理
- 删除用户（新增，含权限控制）
- 用户搜索（支持姓名、学号、手机号）

###  后续规划

- [ ] 接入微信/短信通知服务
- [ ] 引入 ECharts 图表统计
- [ ] 支持 Excel 数据导入导出
- [ ] 移动端适配优化 (H5/小程序)


##  开发团队
  后端开发 [yurin He](https://github.com/yurin-kami)

  前端开发 [RuiYier](https://github.com/RuiYier)


---
