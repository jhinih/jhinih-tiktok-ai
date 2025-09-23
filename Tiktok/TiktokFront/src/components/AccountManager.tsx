import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { 
  User, 
  LogOut, 
  Settings, 
  UserPlus, 
  ChevronDown,
  Shield,
  Clock,
  Wifi,
  WifiOff
} from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { tokenService } from '@/services/tokenService'
import { toast } from '@/utils/toast'

interface AccountManagerProps {
  className?: string
}

export default function AccountManager({ className = '' }: AccountManagerProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [showLogoutConfirm, setShowLogoutConfirm] = useState(false)
  const { user, isAuthenticated, isRefreshing, logout } = useAuthStore()

  const handleLogout = () => {
    setShowLogoutConfirm(true)
  }

  const confirmLogout = () => {
    logout()
    tokenService.clearTokens()
    toast.success('已安全退出登录')
    setShowLogoutConfirm(false)
    setIsOpen(false)
    // 跳转到登录页
    window.location.href = '/login'
  }

  const handleSwitchAccount = () => {
    // 清除当前认证状态，但不显示退出提示
    logout()
    tokenService.clearTokens()
    toast.info('请登录其他账号')
    setIsOpen(false)
    // 跳转到登录页
    window.location.href = '/login'
  }

  const handleRefreshToken = async () => {
    try {
      toast.info('正在刷新登录状态...')
      const newToken = await tokenService.refreshAccessToken()
      if (newToken) {
        toast.success('登录状态已刷新')
      } else {
        toast.error('刷新失败，请重新登录')
      }
    } catch (error) {
      toast.error('刷新失败，请重新登录')
    }
  }

  if (!isAuthenticated || !user) {
    return null
  }

  return (
    <div className={`relative ${className}`}>
      {/* 用户头像按钮 */}
      <motion.button
        whileTap={{ scale: 0.95 }}
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center space-x-2 p-2 rounded-full bg-white/10 backdrop-blur-sm border border-white/20 hover:bg-white/20 transition-colors"
      >
        <div className="w-8 h-8 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
          {user.avatar ? (
            <img 
              src={user.avatar} 
              alt={user.username}
              className="w-full h-full rounded-full object-cover"
            />
          ) : (
            <span className="text-white font-bold text-sm">
              {user.username.charAt(0).toUpperCase()}
            </span>
          )}
        </div>
        
        {/* 刷新状态指示器 */}
        {isRefreshing && (
          <div className="w-3 h-3 border border-white border-t-transparent rounded-full animate-spin" />
        )}
        
        <ChevronDown 
          size={16} 
          className={`text-white transition-transform ${isOpen ? 'rotate-180' : ''}`} 
        />
      </motion.button>

      {/* 下拉菜单 */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -10, scale: 0.95 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: -10, scale: 0.95 }}
            transition={{ duration: 0.2 }}
            className="absolute top-full right-0 mt-2 w-72 bg-white/95 backdrop-blur-md rounded-xl shadow-xl border border-white/20 overflow-hidden z-50"
          >
            {/* 用户信息头部 */}
            <div className="p-4 bg-gradient-to-r from-purple-500/20 to-pink-500/20 border-b border-white/10">
              <div className="flex items-center space-x-3">
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
                  {user.avatar ? (
                    <img 
                      src={user.avatar} 
                      alt={user.username}
                      className="w-full h-full rounded-full object-cover"
                    />
                  ) : (
                    <span className="text-white font-bold">
                      {user.username.charAt(0).toUpperCase()}
                    </span>
                  )}
                </div>
                <div className="flex-1">
                  <h3 className="font-bold text-gray-800">{user.username}</h3>
                  <p className="text-sm text-gray-600">{user.email}</p>
                  {user.bio && (
                    <p className="text-xs text-gray-500 mt-1 line-clamp-1">{user.bio}</p>
                  )}
                </div>
                <div className="flex items-center space-x-1">
                  {isRefreshing ? (
                    <WifiOff size={16} className="text-orange-500" />
                  ) : (
                    <Wifi size={16} className="text-green-500" />
                  )}
                  <Shield size={16} className="text-blue-500" />
                </div>
              </div>
            </div>

            {/* 菜单选项 */}
            <div className="py-2">
              {/* 个人资料 */}
              <button
                onClick={() => {
                  setIsOpen(false)
                  window.location.href = '/profile'
                }}
                className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-100/50 transition-colors"
              >
                <User size={18} className="text-gray-600" />
                <span className="text-gray-800">个人资料</span>
              </button>

              {/* 设置 */}
              <button
                onClick={() => {
                  setIsOpen(false)
                  // 这里可以打开设置页面
                  toast.info('设置功能开发中...')
                }}
                className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-100/50 transition-colors"
              >
                <Settings size={18} className="text-gray-600" />
                <span className="text-gray-800">设置</span>
              </button>

              {/* 刷新登录状态 */}
              <button
                onClick={handleRefreshToken}
                disabled={isRefreshing}
                className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-100/50 transition-colors disabled:opacity-50"
              >
                <Clock size={18} className="text-gray-600" />
                <span className="text-gray-800">
                  {isRefreshing ? '刷新中...' : '刷新登录状态'}
                </span>
              </button>

              <div className="border-t border-gray-200/50 my-2" />

              {/* 切换账号 */}
              <button
                onClick={handleSwitchAccount}
                className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-gray-100/50 transition-colors"
              >
                <UserPlus size={18} className="text-blue-600" />
                <span className="text-gray-800">切换账号</span>
              </button>

              {/* 退出登录 */}
              <button
                onClick={handleLogout}
                className="w-full flex items-center space-x-3 px-4 py-3 hover:bg-red-50 transition-colors"
              >
                <LogOut size={18} className="text-red-600" />
                <span className="text-red-600">退出登录</span>
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* 退出确认对话框 */}
      <AnimatePresence>
        {showLogoutConfirm && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-[100]"
            onClick={() => setShowLogoutConfirm(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-white rounded-2xl p-6 max-w-sm mx-4 shadow-2xl"
            >
              <div className="text-center">
                <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <LogOut size={32} className="text-red-600" />
                </div>
                <h3 className="text-xl font-bold text-gray-800 mb-2">确认退出登录</h3>
                <p className="text-gray-600 mb-6">
                  您确定要退出当前账号吗？退出后需要重新登录才能使用。
                </p>
                <div className="flex space-x-3">
                  <button
                    onClick={() => setShowLogoutConfirm(false)}
                    className="flex-1 py-3 px-4 bg-gray-100 text-gray-800 rounded-xl font-medium hover:bg-gray-200 transition-colors"
                  >
                    取消
                  </button>
                  <button
                    onClick={confirmLogout}
                    className="flex-1 py-3 px-4 bg-red-500 text-white rounded-xl font-medium hover:bg-red-600 transition-colors"
                  >
                    确认退出
                  </button>
                </div>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* 点击外部关闭菜单 */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40"
          onClick={() => setIsOpen(false)}
        />
      )}
    </div>
  )
}