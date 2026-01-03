import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomePage.vue'),
      meta: { public: true }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginPage.vue'),
      meta: { public: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterPage.vue'),
      meta: { public: true }
    },
    {
      path: '/user',
      name: 'user',
      component: () => import('@/views/user/UserLayout.vue'),
      meta: { requiresAuth: true, role: 'user' },
      children: [
        {
          path: '',
          name: 'user-dashboard',
          component: () => import('@/views/user/UserDashboard.vue')
        },
        {
          path: 'packs',
          name: 'user-packs',
          component: () => import('@/views/user/UserPacks.vue')
        },
        {
          path: 'mail',
          name: 'user-mail',
          component: () => import('@/views/user/UserMail.vue')
        },
        {
          path: 'profile',
          name: 'user-profile',
          component: () => import('@/views/user/UserProfile.vue')
        }
      ]
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true, role: 'admin' },
      children: [
        {
          path: '',
          name: 'admin-dashboard',
          component: () => import('@/views/admin/AdminDashboard.vue')
        },
        {
          path: 'packs',
          name: 'admin-packs',
          component: () => import('@/views/admin/AdminPacks.vue')
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('@/views/admin/AdminUsers.vue')
        },
        {
          path: 'checkin',
          name: 'admin-checkin',
          component: () => import('@/views/admin/AdminCheckIn.vue')
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'notfound',
      component: () => import('@/views/NotFound.vue')
    }
  ],
})

// 导航守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  // 公开页面直接放行
  if (to.meta.public) {
    // 如果已登录，访问登录/注册页时重定向到对应的界面
    if (authStore.isAuthenticated && (to.name === 'login' || to.name === 'register')) {
      if (authStore.isAdmin) {
        next({ name: 'admin-dashboard' })
      } else {
        next({ name: 'user-dashboard' })
      }
      return
    }
    next()
    return
  }

  // 需要认证的页面
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      next({ name: 'login', query: { redirect: to.fullPath } })
      return
    }

    // 检查角色权限
    if (to.meta.role) {
      if (authStore.user?.role !== to.meta.role) {
        // 角色不匹配，重定向到对应的主页
        if (authStore.isAdmin) {
          next({ name: 'admin-dashboard' })
        } else {
          next({ name: 'user-dashboard' })
        }
        return
      }
    }
  }

  next()
})

export default router

