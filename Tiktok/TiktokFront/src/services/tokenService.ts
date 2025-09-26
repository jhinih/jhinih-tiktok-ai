import axios, { AxiosError } from 'axios'
import { useAuthStore } from '@/stores/authStore'
import { toast } from '@/utils/toast'

// Token刷新相关的API
export const refreshTokenApi = async (refreshToken: string) => {
  const response = await axios.post('/api/refresh-token', {
    rtoken: refreshToken
  })
  return response.data
}

// Token管理服务
class TokenService {
  private refreshPromise: Promise<string> | null = null

  // 检查token是否即将过期
  private isTokenExpiringSoon(token: string): boolean {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      const exp = payload.exp * 1000 // 转换为毫秒
      const now = Date.now()
      // 如果token在5分钟内过期，则认为即将过期
      return exp - now < 5 * 60 * 1000
    } catch {
      return true // 如果解析失败，认为token无效
    }
  }

  // 检查token是否已过期
  private isTokenExpired(token: string): boolean {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      const exp = payload.exp * 1000
      return Date.now() >= exp
    } catch {
      return true
    }
  }

  // 刷新访问token
  async refreshAccessToken(): Promise<string | null> {
    const { refreshToken, setRefreshing, setToken, clearAuth } = useAuthStore.getState()

    if (!refreshToken) {
      console.warn('没有refresh token，无法刷新')
      return null
    }

    // 检查refresh token是否过期
    if (this.isTokenExpired(refreshToken)) {
      console.warn('Refresh token已过期，需要重新登录')
      clearAuth()
      toast.error('登录已过期，请重新登录')
      // 跳转到登录页
      window.location.href = '/login'
      return null
    }

    // 如果已经在刷新中，返回现有的Promise
    if (this.refreshPromise) {
      return this.refreshPromise
    }

    setRefreshing(true)

    this.refreshPromise = (async () => {
      try {
        console.log('正在刷新access token...')
        const response = await refreshTokenApi(refreshToken)
        
        if (response.code === 200 && response.data?.atoken) {
          const newAccessToken = response.data.atoken
          setToken(newAccessToken)
          console.log('Access token刷新成功')
          toast.success('登录状态已更新')
          return newAccessToken
        } else {
          throw new Error('刷新token失败: ' + (response.message || '未知错误'))
        }
      } catch (error) {
        console.error('刷新token失败:', error)
        
        if (error instanceof AxiosError) {
          if (error.response?.status === 401 || error.response?.status === 403) {
            // refresh token无效，清除认证状态
            clearAuth()
            toast.error('登录已过期，请重新登录')
            window.location.href = '/login'
          } else {
            toast.error('网络错误，请检查连接')
          }
        } else {
          toast.error('刷新登录状态失败')
        }
        
        return null
      } finally {
        setRefreshing(false)
        this.refreshPromise = null
      }
    })()

    return this.refreshPromise
  }

  // 获取有效的访问token
  async getValidAccessToken(): Promise<string | null> {
    const { token, isRefreshing } = useAuthStore.getState()

    if (!token) {
      return null
    }

    // 如果正在刷新，等待刷新完成
    if (isRefreshing && this.refreshPromise) {
      return await this.refreshPromise
    }

    // 如果token已过期，尝试刷新
    if (this.isTokenExpired(token)) {
      console.log('Access token已过期，尝试刷新...')
      return await this.refreshAccessToken()
    }

    // 如果token即将过期，预先刷新
    if (this.isTokenExpiringSoon(token)) {
      console.log('Access token即将过期，预先刷新...')
      // 异步刷新，不阻塞当前请求
      this.refreshAccessToken().catch(console.error)
    }

    return token
  }

  // 清除所有token相关状态
  clearTokens() {
    const { clearAuth } = useAuthStore.getState()
    clearAuth()
    this.refreshPromise = null
  }

  // 启动定时检查token状态
  startTokenMonitoring() {
    // 每分钟检查一次token状态
    setInterval(() => {
      const { token, isAuthenticated } = useAuthStore.getState()
      
      if (!isAuthenticated || !token) {
        return
      }

      // 如果token即将过期，预先刷新
      if (this.isTokenExpiringSoon(token)) {
        console.log('定时检查: token即将过期，开始刷新...')
        this.refreshAccessToken().catch(console.error)
      }
    }, 60 * 1000) // 每分钟检查一次
  }
}

export const tokenService = new TokenService()

// 在应用启动时开始监控
if (typeof window !== 'undefined') {
  tokenService.startTokenMonitoring()
}