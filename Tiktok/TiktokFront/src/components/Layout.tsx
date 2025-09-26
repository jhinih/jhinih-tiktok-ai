import { Outlet } from 'react-router-dom'
import { useEffect } from 'react'
import BottomNavigation from './BottomNavigation'
import ConnectionStatus from './ConnectionStatus'
import { wsService } from '@/services/websocket'
import { useAuthStore } from '@/stores/authStore'

export default function Layout() {
  const { isAuthenticated } = useAuthStore()

  useEffect(() => {
    // 建立WebSocket连接
    if (isAuthenticated) {
      wsService.connect()
    }

    // 清理函数
    return () => {
      wsService.disconnect()
    }
  }, [isAuthenticated])

  return (
    <div className="min-h-screen bg-black text-white">
      <ConnectionStatus />
      <main className="pb-16">
        <Outlet />
      </main>
      <BottomNavigation />
    </div>
  )
}