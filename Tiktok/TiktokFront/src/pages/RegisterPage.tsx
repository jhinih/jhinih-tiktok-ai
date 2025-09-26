import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { Eye, EyeOff, Mail, Lock, User, Image } from 'lucide-react'
import { motion } from 'framer-motion'
import toast from 'react-hot-toast'
import { registerApi, sendCodeApi } from '@/services/api'
import { useAuthStore } from '@/stores/authStore'
import { cn } from '@/utils'
import type { RegisterRequest } from '@/types'

export default function RegisterPage() {
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [isSendingCode, setIsSendingCode] = useState(false)
  const [countdown, setCountdown] = useState(0)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<RegisterRequest>()

  const email = watch('email')

  const sendVerificationCode = async () => {
    if (!email) {
      toast.error('请先输入邮箱地址')
      return
    }

    setIsSendingCode(true)
    try {
      await sendCodeApi(email)
      toast.success('验证码已发送')
      
      // 开始倒计时
      setCountdown(60)
      const timer = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            clearInterval(timer)
            return 0
          }
          return prev - 1
        })
      }, 1000)
    } catch (error) {
      toast.error('发送验证码失败')
    } finally {
      setIsSendingCode(false)
    }
  }

  const onSubmit = async (data: RegisterRequest) => {
    setIsLoading(true)
    try {
      const response = await registerApi(data)
      const { atoken } = response.data.data
      
      // 创建用户对象
      const user = {
        id: '1', // 实际应该从响应中获取
        username: data.username,
        email: data.email,
        avatar: data.avatar || '',
        role: '0',  // 修改为字符串类型
        bio: '',
        CreatedTime: '',
        UpdatedTime: '',
        Phone: '',
        ClientIp: '',
        ClientPort: '',
        LoginTime: '',
        HeartbeatTime: '',
        LoginOutTime: '',  // 修正字段名
        IsLogout: false,
        DeviceInfo: '',
      }
      
      login(user, atoken)
      toast.success('注册成功')
      navigate('/')
    } catch (error) {
      toast.error('注册失败，请检查信息')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900 flex items-center justify-center p-4">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md"
      >
        <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-8 shadow-2xl border border-white/20">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-white mb-2">加入TikTok</h1>
            <p className="text-gray-300">创建您的账户开始分享</p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                用户名
              </label>
              <div className="relative">
                <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                <input
                  {...register('username', {
                    required: '请输入用户名',
                    minLength: {
                      value: 2,
                      message: '用户名至少2位',
                    },
                  })}
                  type="text"
                  className={cn(
                    'w-full pl-10 pr-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
                    errors.username && 'border-red-500'
                  )}
                  placeholder="输入用户名"
                />
              </div>
              {errors.username && (
                <p className="mt-1 text-sm text-red-400">{errors.username.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                邮箱地址
              </label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                <input
                  {...register('email', {
                    required: '请输入邮箱地址',
                    pattern: {
                      value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                      message: '请输入有效的邮箱地址',
                    },
                  })}
                  type="email"
                  className={cn(
                    'w-full pl-10 pr-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
                    errors.email && 'border-red-500'
                  )}
                  placeholder="输入邮箱地址"
                />
              </div>
              {errors.email && (
                <p className="mt-1 text-sm text-red-400">{errors.email.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                验证码
              </label>
              <div className="flex space-x-2">
                <input
                  {...register('code', {
                    required: '请输入验证码',
                    pattern: {
                      value: /^\d{6}$/,
                      message: '验证码为6位数字',
                    },
                  })}
                  type="text"
                  className={cn(
                    'flex-1 px-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
                    errors.code && 'border-red-500'
                  )}
                  placeholder="输入验证码"
                />
                <button
                  type="button"
                  onClick={sendVerificationCode}
                  disabled={isSendingCode || countdown > 0}
                  className="px-4 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors whitespace-nowrap"
                >
                  {countdown > 0 ? `${countdown}s` : isSendingCode ? '发送中' : '发送验证码'}
                </button>
              </div>
              {errors.code && (
                <p className="mt-1 text-sm text-red-400">{errors.code.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                密码
              </label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                <input
                  {...register('password', {
                    required: '请输入密码',
                    minLength: {
                      value: 6,
                      message: '密码至少6位',
                    },
                  })}
                  type={showPassword ? 'text' : 'password'}
                  className={cn(
                    'w-full pl-10 pr-12 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
                    errors.password && 'border-red-500'
                  )}
                  placeholder="输入密码"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-300"
                >
                  {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                </button>
              </div>
              {errors.password && (
                <p className="mt-1 text-sm text-red-400">{errors.password.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                头像链接（可选）
              </label>
              <div className="relative">
                <Image className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
                <input
                  {...register('avatar')}
                  type="url"
                  className="w-full pl-10 pr-4 py-3 bg-white/10 border border-white/20 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="输入头像链接"
                />
              </div>
            </div>

            <motion.button
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              type="submit"
              disabled={isLoading}
              className="w-full bg-gradient-to-r from-blue-500 to-purple-600 text-white py-3 rounded-lg font-semibold hover:from-blue-600 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-transparent disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
            >
              {isLoading ? '注册中...' : '注册'}
            </motion.button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-gray-300">
              已有账户？{' '}
              <Link
                to="/login"
                className="text-blue-400 hover:text-blue-300 font-semibold transition-colors"
              >
                立即登录
              </Link>
            </p>
          </div>
        </div>
      </motion.div>
    </div>
  )
}