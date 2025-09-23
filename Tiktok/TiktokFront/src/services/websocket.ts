import { useAuthStore } from '@/stores/authStore'
import { toast } from '@/utils/toast'

class WebSocketService {
  private ws: WebSocket | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectInterval = 3000
  private messageHandlers: ((data: any) => void)[] = []

  async connect() {
    try {
      // 获取当前用户的token
      const { accessToken } = useAuthStore.getState()
      if (!accessToken) {
        throw new Error('未找到访问令牌')
      }
      
      // 建立WebSocket连接，通过查询参数传递token
      const wsUrl = `ws://localhost:8080/ws?token=Bearer ${encodeURIComponent(accessToken)}`
      this.ws = new WebSocket(wsUrl)
      
      this.ws.onopen = () => {
        console.log('WebSocket连接已建立')
        this.reconnectAttempts = 0
        toast.success('WebSocket连接成功')
      }
      
      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          console.log('收到WebSocket消息:', data)
          this.messageHandlers.forEach(handler => handler(data))
        } catch (error) {
          console.error('解析WebSocket消息失败:', error)
        }
      }
      
      this.ws.onclose = (event) => {
        console.log('WebSocket连接已关闭', event.code, event.reason)
        this.reconnect()
      }
      
      this.ws.onerror = (error) => {
        console.error('WebSocket错误:', error)
      }
      
    } catch (error) {
      console.error('WebSocket连接失败:', error)
      throw error
    }
  }



  private reconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
      
      setTimeout(() => {
        this.connect()
      }, this.reconnectInterval)
    } else {
      console.error('WebSocket重连失败，已达到最大重试次数')
      toast.error('连接失败，请刷新页面重试')
    }
  }

  // 发送消息
  sendMessage(data: {
    type: 'private' | 'group' | 'chat'
    receiverId?: string
    groupId?: string
    content: string
    messageType?: 'text' | 'image' | 'video' | 'file'
  }) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message = {
        type: data.type,
        content: data.content,
        to_type: data.type === 'private' ? 'user' : 'group',
        to: data.receiverId || data.groupId || 0,
        message_type: data.messageType || 'text'
      }
      this.ws.send(JSON.stringify(message))
      console.log('发送消息:', message)
    } else {
      toast.error('连接已断开，请重新连接')
    }
  }

  // 添加消息处理器
  onMessage(handler: (data: any) => void) {
    this.messageHandlers.push(handler)
  }

  // 移除消息处理器
  offMessage(handler: (data: any) => void) {
    const index = this.messageHandlers.indexOf(handler)
    if (index > -1) {
      this.messageHandlers.splice(index, 1)
    }
  }

  // 断开连接
  disconnect() {
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.messageHandlers = []
  }

  // 获取连接状态
  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN
  }

  // 获取连接状态文本
  getConnectionStatus(): string {
    if (!this.ws) return 'disconnected'
    
    switch (this.ws.readyState) {
      case WebSocket.CONNECTING:
        return 'connecting'
      case WebSocket.OPEN:
        return 'connected'
      case WebSocket.CLOSING:
        return 'closing'
      case WebSocket.CLOSED:
        return 'closed'
      default:
        return 'unknown'
    }
  }
}

export const wsService = new WebSocketService()

// React Hook for WebSocket
export function useWebSocket() {
  const connect = async () => {
    try {
      await wsService.connect()
    } catch (error) {
      console.error('WebSocket连接失败:', error)
    }
  }

  const sendMessage = (data: {
    type: 'private' | 'group' | 'chat'
    receiverId?: string
    groupId?: string
    content: string
    messageType?: 'text' | 'image' | 'video' | 'file'
  }) => {
    wsService.sendMessage(data)
  }

  const disconnect = () => {
    wsService.disconnect()
  }

  return {
    connect,
    sendMessage,
    disconnect,
    isConnected: wsService.isConnected(),
    getConnectionStatus: wsService.getConnectionStatus(),
    onMessage: wsService.onMessage.bind(wsService),
    offMessage: wsService.offMessage.bind(wsService)
  }
}

export default wsService