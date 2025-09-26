import { useState, useEffect, useRef } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Search, Send, Phone, Video, MoreHorizontal, ArrowLeft, Users } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { useNavigate, useParams } from 'react-router-dom'
import { getFriendListApi, createCommunityApi } from '@/services/api'
import { useAuthStore } from '@/stores/authStore'
import { useChatStore } from '@/stores/chatStore'
import { useWebSocket } from '@/services/websocket'
import { formatTime } from '@/utils'
import toast from 'react-hot-toast'
import type { Message, ChatRoom } from '@/types'

interface ChatListItemProps {
  chat: ChatRoom
  onClick: () => void
  isActive: boolean
}

function ChatListItem({ chat, onClick, isActive }: ChatListItemProps) {
  return (
    <motion.div
      whileHover={{ backgroundColor: 'rgba(255, 255, 255, 0.05)' }}
      onClick={onClick}
      className={`p-4 cursor-pointer border-b border-gray-800 ${
        isActive ? 'bg-blue-500/20' : ''
      }`}
    >
      <div className="flex items-center space-x-3">
        <div className="relative">
          <div className="w-12 h-12 bg-gradient-to-br from-purple-400 to-pink-400 rounded-full flex items-center justify-center">
            {chat.type === 'group' ? (
              <Users size={20} className="text-white" />
            ) : (
              <span className="text-white font-semibold">
                {chat.name.charAt(0).toUpperCase()}
              </span>
            )}
          </div>
          {chat.online && (
            <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-gray-900" />
          )}
        </div>
        
        <div className="flex-1 min-w-0">
          <div className="flex items-center justify-between">
            <h3 className="text-white font-medium truncate">{chat.name}</h3>
            <span className="text-gray-400 text-xs">
              {formatTime(chat.lastMessage?.created_at || '')}
            </span>
          </div>
          <p className="text-gray-400 text-sm truncate mt-1">
            {chat.lastMessage?.content || '暂无消息'}
          </p>
        </div>
        
        {chat.unreadCount > 0 && (
          <div className="w-5 h-5 bg-red-500 rounded-full flex items-center justify-center">
            <span className="text-white text-xs font-medium">
              {chat.unreadCount > 99 ? '99+' : chat.unreadCount}
            </span>
          </div>
        )}
      </div>
    </motion.div>
  )
}

interface MessageBubbleProps {
  message: Message
  isOwn: boolean
}

function MessageBubble({ message, isOwn }: MessageBubbleProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className={`flex ${isOwn ? 'justify-end' : 'justify-start'} mb-4`}
    >
      <div className={`max-w-xs lg:max-w-md ${isOwn ? 'order-2' : 'order-1'}`}>
        {!isOwn && (
          <p className="text-gray-400 text-xs mb-1 px-3">{message.sender_name}</p>
        )}
        <div
          className={`px-4 py-2 rounded-2xl ${
            isOwn
              ? 'bg-blue-500 text-white'
              : 'bg-gray-700 text-white'
          }`}
        >
          {message.type === 'text' ? (
            <p className="text-sm">{message.content}</p>
          ) : message.type === 'image' ? (
            <img
              src={message.content}
              alt="图片消息"
              className="max-w-full h-auto rounded-lg"
            />
          ) : (
            <div className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-white/20 rounded-full flex items-center justify-center">
                📁
              </div>
              <span className="text-sm">文件消息</span>
            </div>
          )}
        </div>
        <p className="text-gray-500 text-xs mt-1 px-3">
          {formatTime(message.created_at || '')}
        </p>
      </div>
    </motion.div>
  )
}

export default function ChatPage() {
  const { chatId } = useParams()
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const { setActiveChat } = useChatStore()
  const [message, setMessage] = useState('')
  const [searchQuery, setSearchQuery] = useState('')
  const [showCreateGroup, setShowCreateGroup] = useState(false)
  const [groupName, setGroupName] = useState('')
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const queryClient = useQueryClient()

  // WebSocket连接
  const { sendMessage: sendWsMessage } = useWebSocket()

  // 获取聊天列表 - 修正API调用
  const { data: chatListData } = useQuery({
    queryKey: ['chatList'],
    queryFn: () => getFriendListApi({ user_id: user?.id || '' }),
  })

  const chatList: ChatRoom[] = (chatListData?.data.data?.users || []).map((friend: any) => ({
    id: friend.id,
    name: friend.username,
    type: 'private' as const,
    online: !friend.IsLogout,
    lastMessage: null,
    unreadCount: 0
  }))

  // 获取当前聊天消息 - 暂时使用模拟数据，实际需要实现消息历史API
  const messages: Message[] = []
  const messagesLoading = false

  // 发送消息 - 使用WebSocket发送
  const handleSendMessage = () => {
    if (!message.trim() || !chatId || !user) return

    const messageData = {
      targetID: chatId,
      targetName: currentChat?.name || '',
      type: 'text',
      media: '1', // 1表示文字消息
      content: message.trim(),
      createTime: Date.now().toString(),
    }

    // 通过WebSocket发送实时消息
    sendWsMessage('message', messageData)
    setMessage('')
    toast.success('消息已发送')
  }

  // 创建群组
  const createGroupMutation = useMutation({
    mutationFn: (data: { name: string; owner_name: string; icon: string; desc: string }) => 
      createCommunityApi(data),
    onSuccess: () => {
      setShowCreateGroup(false)
      setGroupName('')
      queryClient.invalidateQueries({ queryKey: ['chatList'] })
      toast.success('群组创建成功')
    },
    onError: () => {
      toast.error('创建群组失败')
    }
  })

  const handleCreateGroup = () => {
    if (!groupName.trim() || !user) return

    createGroupMutation.mutate({
      name: groupName.trim(),
      owner_name: user.username,
      icon: '',
      desc: `由 ${user.username} 创建的群组`,
    })
  }

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const filteredChats = chatList.filter(chat =>
    chat.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const currentChat = chatList.find(chat => chat.id === chatId)

  return (
    <div className="h-screen bg-black text-white flex">
      {/* 左侧聊天列表 */}
      <div className="w-80 border-r border-gray-800 flex flex-col">
        {/* 头部 */}
        <div className="p-4 border-b border-gray-800">
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-xl font-bold">消息</h1>
            <button
              onClick={() => setShowCreateGroup(true)}
              className="p-2 bg-blue-500 rounded-full hover:bg-blue-600 transition-colors"
            >
              <Users size={20} />
            </button>
          </div>
          
          {/* 搜索框 */}
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={16} />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="搜索聊天..."
              className="w-full pl-9 pr-4 py-2 bg-gray-800 rounded-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        {/* 聊天列表 */}
        <div className="flex-1 overflow-y-auto">
          {filteredChats.map(chat => (
            <ChatListItem
              key={chat.id}
              chat={chat}
              onClick={() => {
                setActiveChat(chat)
                navigate(`/chat/${chat.id}`)
              }}
              isActive={chat.id === chatId}
            />
          ))}
        </div>
      </div>

      {/* 右侧聊天区域 */}
      <div className="flex-1 flex flex-col">
        {chatId && currentChat ? (
          <>
            {/* 聊天头部 */}
            <div className="p-4 border-b border-gray-800 flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <button
                  onClick={() => navigate('/chat')}
                  className="lg:hidden p-2 hover:bg-gray-800 rounded-full"
                >
                  <ArrowLeft size={20} />
                </button>
                <div className="w-10 h-10 bg-gradient-to-br from-purple-400 to-pink-400 rounded-full flex items-center justify-center">
                  {currentChat.type === 'group' ? (
                    <Users size={20} className="text-white" />
                  ) : (
                    <span className="text-white font-semibold">
                      {currentChat.name.charAt(0).toUpperCase()}
                    </span>
                  )}
                </div>
                <div>
                  <h2 className="font-semibold">{currentChat.name}</h2>
                  <p className="text-gray-400 text-sm">
                    {currentChat.online ? '在线' : '离线'}
                  </p>
                </div>
              </div>
              
              <div className="flex items-center space-x-2">
                <button className="p-2 hover:bg-gray-800 rounded-full">
                  <Phone size={20} />
                </button>
                <button className="p-2 hover:bg-gray-800 rounded-full">
                  <Video size={20} />
                </button>
                <button className="p-2 hover:bg-gray-800 rounded-full">
                  <MoreHorizontal size={20} />
                </button>
              </div>
            </div>

            {/* 消息区域 */}
            <div className="flex-1 overflow-y-auto p-4">
              {messagesLoading ? (
                <div className="flex items-center justify-center h-full">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
                </div>
              ) : (
                <>
                  {messages.map((msg: Message) => (
                    <MessageBubble
                      key={msg.id}
                      message={msg}
                      isOwn={msg.sender_id === user?.id}
                    />
                  ))}
                  <div ref={messagesEndRef} />
                </>
              )}
            </div>

            {/* 输入区域 */}
            <div className="p-4 border-t border-gray-800">
              <div className="flex items-center space-x-2">
                <input
                  type="text"
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  placeholder="输入消息..."
                  className="flex-1 px-4 py-2 bg-gray-800 rounded-full focus:outline-none focus:ring-2 focus:ring-blue-500"
                  onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                />
                <motion.button
                  whileTap={{ scale: 0.95 }}
                  onClick={handleSendMessage}
                  disabled={!message.trim()}
                  className="p-2 bg-blue-500 rounded-full hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  <Send size={20} />
                </motion.button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center text-gray-400">
              <div className="w-20 h-20 bg-gray-800 rounded-full flex items-center justify-center mx-auto mb-4">
                💬
              </div>
              <h3 className="text-xl font-semibold mb-2">选择一个聊天</h3>
              <p>从左侧选择一个聊天开始对话</p>
            </div>
          </div>
        )}
      </div>

      {/* 创建群组弹窗 */}
      <AnimatePresence>
        {showCreateGroup && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
            onClick={() => setShowCreateGroup(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-gray-900 rounded-lg p-6 w-96 max-w-[90vw]"
            >
              <h3 className="text-xl font-bold mb-4">创建群组</h3>
              <input
                type="text"
                value={groupName}
                onChange={(e) => setGroupName(e.target.value)}
                placeholder="输入群组名称"
                className="w-full px-4 py-2 bg-gray-800 rounded-lg mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500"
                autoFocus
              />
              <div className="flex space-x-3">
                <button
                  onClick={() => setShowCreateGroup(false)}
                  className="flex-1 py-2 bg-gray-700 rounded-lg hover:bg-gray-600 transition-colors"
                >
                  取消
                </button>
                <button
                  onClick={handleCreateGroup}
                  disabled={!groupName.trim() || createGroupMutation.isLoading}
                  className="flex-1 py-2 bg-blue-500 rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {createGroupMutation.isLoading ? '创建中...' : '创建'}
                </button>
              </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}