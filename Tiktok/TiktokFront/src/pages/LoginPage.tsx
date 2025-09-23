import { useState } from 'react'
import { motion } from 'framer-motion'
import { Eye, EyeOff, Mail, Lock, User, ArrowRight, Loader2 } from 'lucide-react'
import { useAuthStore } from '@/stores/authStore'
import { loginApi, registerApi, sendCodeApi } from '@/services/api'
import { toast } from '@/utils/toast'
import type { User as UserType } from '@/types'

export default function LoginPage() {
  const [isLogin, setIsLogin] = useState(true)
  const [showPassword, setShowPassword] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [isSendingCode, setIsSendingCode] = useState(false)
  const [codeCountdown, setCodeCountdown] = useState(0)
  
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    username: '',
    code: ''
  })

  const { login } = useAuthStore()

  // 发送验证码
  const handleSendCode = async () => {
    if (!formData.email) {
      toast.error('请输入邮箱地址')
      return
    }

    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      toast.error('请输入有效的邮箱地址')
      return
    }

    setIsSendingCode(true)
    try {
      const response = await sendCodeApi(formData.email)
      
      // 检查响应结构
      if (response.data && response.data.code === 20000) {
        toast.success('验证码已发送到您的邮箱')
        
        // 开始倒计时
        setCodeCountdown(60)
        const timer = setInterval(() => {
          setCodeCountdown((prev) => {
            if (prev <= 1) {
              clearInterval(timer)
              return 0
            }
            return prev - 1
          })
        }, 1000)
      } else {
        // 处理发送失败的情况
        const errorMessage = response.data?.message || '发送验证码失败'
        toast.error(errorMessage)
      }
    } catch (error: any) {
      console.error('发送验证码失败:', error)
      const errorMessage = error.response?.data?.message || error.message || '发送验证码失败'
      toast.error(errorMessage)
    } finally {
      setIsSendingCode(false)
    }
  }

  // 处理登录
  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.email || !formData.password) {
      toast.error('请填写完整的登录信息')
      return
    }

    setIsLoading(true)
    try {
      const response = await loginApi({
        email: formData.email,
        password: formData.password
      })

      console.log('登录响应:', response.data)
      console.log('完整响应数据:', JSON.stringify(response.data, null, 2))

      // 检查响应结构
      if (response.data && response.data.code === 20000) {
        const responseData = response.data.data
        console.log('响应数据:', responseData)
        
        if (responseData && responseData.atoken) {
          const { atoken, rtoken } = responseData
          console.log('获取到的tokens:', { atoken, rtoken })
          
          // 从token中解析用户信息
          try {
            console.log('开始解析token:', atoken)
            console.log('Token类型:', typeof atoken)
            console.log('Token长度:', atoken.length)
            
            // 检查token格式
            if (!atoken || typeof atoken !== 'string') {
              throw new Error('Token不是有效的字符串')
            }
            
            if (!atoken.includes('.')) {
              throw new Error('Token格式不正确，不包含点分隔符')
            }
            
            const tokenParts = atoken.split('.')
            console.log('Token分割结果:', tokenParts)
            console.log('Token部分数量:', tokenParts.length)
            
            if (tokenParts.length < 2) {
              throw new Error('Token格式不正确，部分数量不足')
            }
            
            console.log('准备解码的部分:', tokenParts[1])
            
            // 安全的JWT payload解码函数
            function decodeJWTPayload(base64Payload: string) {
              try {
                // 添加必要的padding
                let base64 = base64Payload.replace(/-/g, '+').replace(/_/g, '/')
                const padding = base64.length % 4
                if (padding) {
                  base64 += '='.repeat(4 - padding)
                }
                
                // 解码base64
                const decoded = atob(base64)
                
                // 处理UTF-8字符（包括中文）
                const utf8Decoded = decodeURIComponent(
                  decoded.split('').map(c => 
                    '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
                  ).join('')
                )
                
                return JSON.parse(utf8Decoded)
              } catch (error) {
                console.error('JWT解码失败:', error)
                throw new Error('Token解码失败')
              }
            }
            
            const payload = decodeJWTPayload(tokenParts[1])
            console.log('解析后的payload:', payload)
            const user: UserType = {
              id: payload.userid,
              username: payload.username,
              email: formData.email,
              role: payload.role.toString(),
              avatar: '',
              bio: '',
              CreatedTime: Date.now().toString(),
              UpdatedTime: Date.now().toString(),
              Phone: '',
              ClientIp: '',
              ClientPort: '',
              LoginTime: Date.now().toString(),
              HeartbeatTime: Date.now().toString(),
              LoginOutTime: '',
              IsLogout: false,
              DeviceInfo: ''
            }

            // 登录成功，保存用户信息和双token
            login(user, atoken, rtoken || '')
            toast.success('登录成功！')
            
            // 跳转到首页
            window.location.href = '/'
          } catch (tokenError) {
            console.error('解析token失败:', tokenError)
            toast.error('登录数据解析失败')
          }
        } else {
          toast.error('登录响应数据异常')
        }
      } else {
        // 处理登录失败的情况
        const errorMessage = response.data?.message || '登录失败'
        toast.error(errorMessage)
      }
    } catch (error: any) {
      console.error('登录失败:', error)
      toast.error(error.response?.data?.message || '登录失败，请检查邮箱和密码')
    } finally {
      setIsLoading(false)
    }
  }

  // 处理注册
  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.email || !formData.password || !formData.username || !formData.code) {
      toast.error('请填写完整的注册信息')
      return
    }

    if (formData.password.length < 6) {
      toast.error('密码长度至少6位')
      return
    }

    setIsLoading(true)
    try {
      const response = await registerApi({
        email: formData.email,
        password: formData.password,
        username: formData.username,
        code: formData.code,
        avatar: ''
      })

      console.log('注册响应:', response.data)

      // 检查响应结构
      if (response.data && response.data.code === 20000) {
        const responseData = response.data.data
        if (responseData && responseData.atoken) {
          const { atoken } = responseData
          
          // 从token中解析用户信息
          try {
            const payload = JSON.parse(atob(atoken.split('.')[1]))
            const user: UserType = {
              id: payload.userid,
              username: payload.username,
              email: formData.email,
              role: payload.role.toString(),
              avatar: '',
              bio: formData.username + '很懒，什么都没留下',
              CreatedTime: Date.now().toString(),
              UpdatedTime: Date.now().toString(),
              Phone: '',
              ClientIp: '',
              ClientPort: '',
              LoginTime: Date.now().toString(),
              HeartbeatTime: Date.now().toString(),
              LoginOutTime: '',
              IsLogout: false,
              DeviceInfo: ''
            }

            // 注册成功，保存用户信息（注册只返回atoken）
            login(user, atoken, '')
            toast.success('注册成功！')
            
            // 跳转到首页
            window.location.href = '/'
          } catch (tokenError) {
            console.error('解析token失败:', tokenError)
            toast.error('注册数据解析失败')
          }
        } else {
          toast.error('注册响应数据异常')
        }
      } else {
        // 处理注册失败的情况
        const errorMessage = response.data?.message || '注册失败'
        toast.error(errorMessage)
      }
    } catch (error: any) {
      console.error('注册失败:', error)
      toast.error(error.response?.data?.message || '注册失败，请检查信息是否正确')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-600 via-pink-600 to-red-600 flex items-center justify-center p-4">
      {/* 背景装饰 */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-white/10 rounded-full blur-3xl" />
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-white/10 rounded-full blur-3xl" />
      </div>

      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="relative w-full max-w-md"
      >
        {/* Logo和标题 */}
        <div className="text-center mb-8">
          <motion.div
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.2, type: "spring", stiffness: 200 }}
            className="w-20 h-20 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center mx-auto mb-4"
          >
            <span className="text-3xl font-bold text-white">T</span>
          </motion.div>
          <h1 className="text-3xl font-bold text-white mb-2">TikTok</h1>
          <p className="text-white/80">
            {isLogin ? '欢迎回来！' : '加入我们的社区'}
          </p>
        </div>

        {/* 登录/注册表单 */}
        <motion.div
          layout
          className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20"
        >
          {/* 切换按钮 */}
          <div className="flex bg-white/10 rounded-xl p-1 mb-6">
            <button
              onClick={() => setIsLogin(true)}
              className={`flex-1 py-2 px-4 rounded-lg font-medium transition-all ${
                isLogin 
                  ? 'bg-white text-purple-600 shadow-lg' 
                  : 'text-white/80 hover:text-white'
              }`}
            >
              登录
            </button>
            <button
              onClick={() => setIsLogin(false)}
              className={`flex-1 py-2 px-4 rounded-lg font-medium transition-all ${
                !isLogin 
                  ? 'bg-white text-purple-600 shadow-lg' 
                  : 'text-white/80 hover:text-white'
              }`}
            >
              注册
            </button>
          </div>

          <form onSubmit={isLogin ? handleLogin : handleRegister} className="space-y-4">
            {/* 用户名 - 仅注册时显示 */}
            {!isLogin && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: 'auto' }}
                exit={{ opacity: 0, height: 0 }}
                className="relative"
              >
                <User className="absolute left-3 top-1/2 transform -translate-y-1/2 text-white/60" size={20} />
                <input
                  type="text"
                  placeholder="用户名"
                  value={formData.username}
                  onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                  className="w-full pl-12 pr-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-white/60 focus:outline-none focus:border-white/40 focus:bg-white/20 transition-all"
                />
              </motion.div>
            )}

            {/* 邮箱 */}
            <div className="relative">
              <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 text-white/60" size={20} />
              <input
                type="email"
                placeholder="邮箱地址"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                className="w-full pl-12 pr-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-white/60 focus:outline-none focus:border-white/40 focus:bg-white/20 transition-all"
              />
            </div>

            {/* 验证码 - 仅注册时显示 */}
            {!isLogin && (
              <motion.div
                initial={{ opacity: 0, height: 0 }}
                animate={{ opacity: 1, height: 'auto' }}
                exit={{ opacity: 0, height: 0 }}
                className="flex space-x-2"
              >
                <input
                  type="text"
                  placeholder="验证码"
                  value={formData.code}
                  onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                  className="flex-1 px-4 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-white/60 focus:outline-none focus:border-white/40 focus:bg-white/20 transition-all"
                />
                <button
                  type="button"
                  onClick={handleSendCode}
                  disabled={isSendingCode || codeCountdown > 0}
                  className="px-4 py-3 bg-white/20 border border-white/20 rounded-xl text-white font-medium hover:bg-white/30 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
                >
                  {isSendingCode ? (
                    <Loader2 size={16} className="animate-spin" />
                  ) : codeCountdown > 0 ? (
                    `${codeCountdown}s`
                  ) : (
                    '发送'
                  )}
                </button>
              </motion.div>
            )}

            {/* 密码 */}
            <div className="relative">
              <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 text-white/60" size={20} />
              <input
                type={showPassword ? 'text' : 'password'}
                placeholder="密码"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                className="w-full pl-12 pr-12 py-3 bg-white/10 border border-white/20 rounded-xl text-white placeholder-white/60 focus:outline-none focus:border-white/40 focus:bg-white/20 transition-all"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 transform -translate-y-1/2 text-white/60 hover:text-white transition-colors"
              >
                {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
              </button>
            </div>

            {/* 提交按钮 */}
            <motion.button
              whileTap={{ scale: 0.98 }}
              type="submit"
              disabled={isLoading}
              className="w-full py-3 bg-white text-purple-600 rounded-xl font-bold hover:bg-white/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center space-x-2"
            >
              {isLoading ? (
                <Loader2 size={20} className="animate-spin" />
              ) : (
                <>
                  <span>{isLogin ? '登录' : '注册'}</span>
                  <ArrowRight size={20} />
                </>
              )}
            </motion.button>
          </form>

          {/* 忘记密码 - 仅登录时显示 */}
          {isLogin && (
            <div className="text-center mt-4">
              <button className="text-white/80 hover:text-white text-sm transition-colors">
                忘记密码？
              </button>
            </div>
          )}
        </motion.div>

        {/* 底部提示 */}
        <div className="text-center mt-6 text-white/60 text-sm">
          <p>登录即表示您同意我们的</p>
          <div className="flex justify-center space-x-4 mt-1">
            <button className="hover:text-white transition-colors">服务条款</button>
            <span>·</span>
            <button className="hover:text-white transition-colors">隐私政策</button>
          </div>
        </div>
      </motion.div>
    </div>
  )
}