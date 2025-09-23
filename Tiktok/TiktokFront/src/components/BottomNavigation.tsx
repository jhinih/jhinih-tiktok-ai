import { NavLink, useLocation } from 'react-router-dom'
import { 
  Home, 
  Search, 
  PlusCircle, 
  MessageCircle, 
  User 
} from 'lucide-react'
import { cn } from '@/utils'

const navItems = [
  {
    path: '/',
    icon: Home,
    label: '首页',
  },
  {
    path: '/discover',
    icon: Search,
    label: '发现',
  },
  {
    path: '/upload',
    icon: PlusCircle,
    label: '发布',
  },
  {
    path: '/chat',
    icon: MessageCircle,
    label: '消息',
  },
  {
    path: '/profile',
    icon: User,
    label: '我的',
  },
]

export default function BottomNavigation() {
  const location = useLocation()

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-black border-t border-gray-800 z-50">
      <div className="flex items-center justify-around h-16 px-2">
        {navItems.map((item) => {
          const Icon = item.icon
          const isActive = location.pathname === item.path
          
          return (
            <NavLink
              key={item.path}
              to={item.path}
              className={cn(
                'flex flex-col items-center justify-center flex-1 py-2 px-1 transition-colors',
                isActive 
                  ? 'text-white' 
                  : 'text-gray-400 hover:text-gray-300'
              )}
            >
              <Icon 
                size={24} 
                className={cn(
                  'mb-1',
                  isActive && item.path === '/upload' && 'text-red-500'
                )}
              />
              <span className="text-xs font-medium">{item.label}</span>
            </NavLink>
          )
        })}
      </div>
    </nav>
  )
}