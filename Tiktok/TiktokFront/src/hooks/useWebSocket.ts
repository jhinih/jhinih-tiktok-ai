import { useEffect, useRef, useState } from 'react'
import { wsService } from '@/services/websocket'
import { useAuthStore } from '@/stores/authStore'
import toast from 'react-hot-toast'

interface UseWebSocketOptions {
  onMessage?: (data: any) => void
  onUserStatus?: (data: any) => void
  autoConnect?: boolean
}

export const useWebSocket = (options: UseWebSocketOptions = {}) => {
  const { onMessage, onUserStatus, autoConnect = true } = options
  const { isAuthenticated } = useAuthStore()
  const [isConnected, setIsConnected] = useState(false)
  const [isConnecting, setIsConnecting] = useState(false)
  const reconnectTimeoutRef = useRef<number | null>(null)

  const connect = async () => {
    if (!isAuthenticated || isConnecting || isConnected) {
      return
    }

    try {
      setIsConnecting(true)
      const socket = await wsService.connect()
      
      socket.on('connect', () => {
        setIsConnected(true)
        setIsConnecting(false)
        console.log('WebSocket连接成功')
      })

      socket.on('disconnect', () => {
        setIsConnected(false)
        setIsConnecting(false)
        console.log('WebSocket连接断开')
        
        // 自动重连
        if (autoConnect && isAuthenticated) {
          reconnectTimeoutRef.current = window.setTimeout(() => {
            connect()
          }, 3000)
        }
      })

      socket.on('connect_error', (error) => {
        setIsConnected(false)
        setIsConnecting(false)
        console.error('WebSocket连接错误:', error)
        toast.error('连接失败，正在重试...')
      })

      // 设置消息监听器
      if (onMessage) {
        wsService.onMessage(onMessage)
      }

      if (onUserStatus) {
        wsService.onUserStatus(onUserStatus)
      }

    } catch (error) {
      setIsConnecting(false)
      console.error('WebSocket连接失败:', error)
      toast.error('连接失败')
    }
  }

  const disconnect = () => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }
    wsService.disconnect()
    setIsConnected(false)
    setIsConnecting(false)
  }

  const sendMessage = (data: {
    type: 'private' | 'group'
    receiverId?: string
    groupId?: string
    content: string
    messageType: 'text' | 'image' | 'video' | 'file'
  }) => {
    if (isConnected) {
      wsService.sendMessage(data)
    } else {
      toast.error('连接已断开，请重新连接')
    }
  }

  const joinGroup = (groupId: string) => {
    if (isConnected) {
      wsService.joinGroup(groupId)
    }
  }

  const leaveGroup = (groupId: string) => {
    if (isConnected) {
      wsService.leaveGroup(groupId)
    }
  }

  useEffect(() => {
    if (autoConnect && isAuthenticated && !isConnected && !isConnecting) {
      connect()
    }

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
    }
  }, [isAuthenticated, autoConnect])

  useEffect(() => {
    return () => {
      disconnect()
    }
  }, [])

  return {
    isConnected,
    isConnecting,
    connect,
    disconnect,
    sendMessage,
    joinGroup,
    leaveGroup
  }
}