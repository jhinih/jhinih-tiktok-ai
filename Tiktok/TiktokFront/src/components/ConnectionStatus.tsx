import { useState, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Wifi, WifiOff } from 'lucide-react'

export default function ConnectionStatus() {
  const [status, setStatus] = useState<'connecting' | 'connected' | 'error'>('connecting')
  const [message, setMessage] = useState('正在连接后端服务...')

  useEffect(() => {
    const checkConnection = async () => {
      try {
        const response = await fetch('/api/videos?page=1&page_size=1', {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        })
        
        if (response.ok) {
          setStatus('connected')
          setMessage('后端连接正常')
        } else {
          setStatus('error')
          setMessage(`连接失败: ${response.status}`)
        }
      } catch (error) {
        setStatus('error')
        setMessage('无法连接到后端服务')
        console.error('Connection test failed:', error)
      }
    }

    checkConnection()
    const interval = setInterval(checkConnection, 30000) // 每30秒检查一次

    return () => clearInterval(interval)
  }, [])

  if (status === 'connected') {
    return null // 连接正常时不显示
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: -50 }}
      animate={{ opacity: 1, y: 0 }}
      className="fixed top-4 right-4 z-50"
    >
      <div className={`flex items-center space-x-2 px-4 py-2 rounded-lg backdrop-blur-sm border ${
        status === 'connecting' 
          ? 'bg-blue-500/20 border-blue-500/30 text-blue-200' 
          : 'bg-red-500/20 border-red-500/30 text-red-200'
      }`}>
        {status === 'connecting' ? (
          <div className="w-4 h-4 border-2 border-blue-200 border-t-transparent rounded-full animate-spin" />
        ) : status === 'error' ? (
          <WifiOff size={16} />
        ) : (
          <Wifi size={16} />
        )}
        <span className="text-sm font-medium">{message}</span>
      </div>
    </motion.div>
  )
}