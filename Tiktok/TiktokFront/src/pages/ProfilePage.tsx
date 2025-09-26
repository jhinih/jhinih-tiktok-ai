import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useParams, useNavigate } from 'react-router-dom'
import { 
  Settings, 
  Edit, 
  Share, 
  MoreHorizontal, 
  Grid, 
  Heart, 
  MessageCircle,
  ArrowLeft,
  Camera,
  MapPin,
  Calendar,
  Link as LinkIcon
} from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { getUserProfileApi, updateUserProfileApi, followUserApi, unfollowUserApi } from '@/services/api'
import { useAuthStore } from '@/stores/authStore'
import { formatNumber, formatDate } from '@/utils'
import { cn } from '@/utils'
import type { UpdateProfileRequest } from '@/types'

interface ProfileStatsProps {
  label: string
  count: number
  onClick?: () => void
}

function ProfileStats({ label, count, onClick }: ProfileStatsProps) {
  return (
    <motion.button
      whileTap={{ scale: 0.95 }}
      onClick={onClick}
      className="text-center hover:bg-gray-800/50 rounded-lg p-2 transition-colors"
    >
      <div className="text-xl font-bold text-white">{formatNumber(count)}</div>
      <div className="text-gray-400 text-sm">{label}</div>
    </motion.button>
  )
}

interface VideoGridProps {
  videos: any[]
  isLoading: boolean
}

function VideoGrid({ videos, isLoading }: VideoGridProps) {
  if (isLoading) {
    return (
      <div className="grid grid-cols-3 gap-1">
        {Array.from({ length: 9 }).map((_, i) => (
          <div key={i} className="aspect-[9/16] bg-gray-800 animate-pulse" />
        ))}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-3 gap-1">
      {videos.map((video, index) => (
        <motion.div
          key={video.id}
          initial={{ opacity: 0, scale: 0.9 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ delay: index * 0.1 }}
          whileHover={{ scale: 1.05 }}
          className="relative aspect-[9/16] bg-gray-800 rounded-lg overflow-hidden cursor-pointer group"
        >
          <img
            src={video.cover || '/placeholder-video.jpg'}
            alt={video.title}
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity" />
          <div className="absolute bottom-2 left-2 right-2">
            <div className="flex items-center justify-between text-white text-xs">
              <span className="flex items-center">
                <Heart size={12} className="mr-1" />
                {formatNumber(video.likes)}
              </span>
              <span className="flex items-center">
                <MessageCircle size={12} className="mr-1" />
                {formatNumber(video.comments)}
              </span>
            </div>
          </div>
        </motion.div>
      ))}
    </div>
  )
}

export default function ProfilePage() {
  const { userId } = useParams()
  const navigate = useNavigate()
  const { user: currentUser, updateUser } = useAuthStore()
  const [activeTab, setActiveTab] = useState<'videos' | 'liked'>('videos')
  const [showEditModal, setShowEditModal] = useState(false)
  const [isFollowing, setIsFollowing] = useState(false)
  const queryClient = useQueryClient()

  const isOwnProfile = !userId || userId === currentUser?.id

  // 获取用户资料
  const { data: profileData, isLoading: profileLoading } = useQuery({
    queryKey: ['userProfile', userId || currentUser?.id],
    queryFn: () => getUserProfileApi(userId || currentUser?.id || ''),
    enabled: !!(userId || currentUser?.id),
  })

  const profile = profileData?.data.data || currentUser
  
  // 辅助函数：获取用户数据
  const getUser = (data: any) => {
    if (!data) return null
    return 'user' in data ? data.user : data
  }
  
  const user = getUser(profile)

  // 更新资料
  const updateProfileMutation = useMutation({
    mutationFn: updateUserProfileApi,
    onSuccess: (response) => {
      const updatedUser = response.data.data
      if (isOwnProfile) {
        updateUser(updatedUser)
      }
      queryClient.invalidateQueries({ queryKey: ['userProfile'] })
      setShowEditModal(false)
      toast.success('资料更新成功')
    },
    onError: () => {
      toast.error('更新失败')
    }
  })

  // 关注/取消关注
  const followMutation = useMutation({
    mutationFn: isFollowing ? unfollowUserApi : followUserApi,
    onSuccess: () => {
      setIsFollowing(!isFollowing)
      queryClient.invalidateQueries({ queryKey: ['userProfile'] })
      toast.success(isFollowing ? '已取消关注' : '关注成功')
    },
    onError: () => {
      toast.error('操作失败')
    }
  })

  const {
    register,
    handleSubmit,
    formState: { errors },

  } = useForm<UpdateProfileRequest>()

  const onSubmit = (data: UpdateProfileRequest) => {
    if (!currentUser) return
    
    updateProfileMutation.mutate(data)
  }

  const handleFollow = () => {
    if (!currentUser || !user) return
    
    followMutation.mutate(user.id)
  }

  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: `${user?.username}的主页`,
        url: window.location.href,
      })
    } else {
      navigator.clipboard.writeText(window.location.href)
      toast.success('链接已复制到剪贴板')
    }
  }

  if (profileLoading) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-black text-white">
      {/* 头部 */}
      <div className="sticky top-0 bg-black/80 backdrop-blur-lg border-b border-gray-800 p-4 z-10">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <button
              onClick={() => navigate(-1)}
              className="p-2 hover:bg-gray-800 rounded-full transition-colors"
            >
              <ArrowLeft size={24} />
            </button>
            <h1 className="text-xl font-bold">{user?.username}</h1>
          </div>
          
          <div className="flex items-center space-x-2">
            <button
              onClick={handleShare}
              className="p-2 hover:bg-gray-800 rounded-full transition-colors"
            >
              <Share size={20} />
            </button>
            <button className="p-2 hover:bg-gray-800 rounded-full transition-colors">
              <MoreHorizontal size={20} />
            </button>
          </div>
        </div>
      </div>

      <div className="p-4">
        {/* 用户信息 */}
        <div className="flex items-start space-x-4 mb-6">
          <div className="relative">
            <div className="w-24 h-24 bg-gradient-to-br from-purple-400 to-pink-400 rounded-full flex items-center justify-center overflow-hidden">
              {user?.avatar ? (
                <img
                  src={user.avatar}
                  alt={user.username}
                  className="w-full h-full object-cover"
                />
              ) : (
                <span className="text-white text-2xl font-bold">
                  {user?.username?.charAt(0).toUpperCase()}
                </span>
              )}
            </div>
            {isOwnProfile && (
              <button
                onClick={() => setShowEditModal(true)}
                className="absolute -bottom-1 -right-1 w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center hover:bg-blue-600 transition-colors"
              >
                <Camera size={16} className="text-white" />
              </button>
            )}
          </div>
          
          <div className="flex-1">
            <h2 className="text-2xl font-bold mb-1">{user?.username}</h2>
            <p className="text-gray-400 mb-3">@{user?.username}</p>
            
            {/* 统计数据 */}
            <div className="flex space-x-6 mb-4">
              <ProfileStats label="关注" count={user?.following_count || 0} />
              <ProfileStats label="粉丝" count={user?.followers_count || 0} />
              <ProfileStats label="获赞" count={user?.likes_count || 0} />
            </div>
            
            {/* 操作按钮 */}
            <div className="flex space-x-3">
              {isOwnProfile ? (
                <>
                  <motion.button
                    whileTap={{ scale: 0.95 }}
                    onClick={() => setShowEditModal(true)}
                    className="flex-1 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors flex items-center justify-center space-x-2"
                  >
                    <Edit size={16} />
                    <span>编辑资料</span>
                  </motion.button>
                  <motion.button
                    whileTap={{ scale: 0.95 }}
                    onClick={() => navigate('/settings')}
                    className="p-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors"
                  >
                    <Settings size={16} />
                  </motion.button>
                </>
              ) : (
                <>
                  <motion.button
                    whileTap={{ scale: 0.95 }}
                    onClick={handleFollow}
                    disabled={followMutation.isLoading}
                    className={cn(
                      'flex-1 py-2 rounded-lg transition-colors font-medium',
                      isFollowing
                        ? 'bg-gray-800 text-white hover:bg-gray-700'
                        : 'bg-blue-500 text-white hover:bg-blue-600'
                    )}
                  >
                    {followMutation.isLoading
                      ? '处理中...'
                      : isFollowing
                      ? '已关注'
                      : '关注'
                    }
                  </motion.button>
                  <motion.button
                    whileTap={{ scale: 0.95 }}
                    onClick={() => navigate(`/chat/new?userId=${user?.id}`)}
                    className="px-4 py-2 bg-gray-800 text-white rounded-lg hover:bg-gray-700 transition-colors"
                  >
                    消息
                  </motion.button>
                </>
              )}
            </div>
          </div>
        </div>

        {/* 个人简介 */}
        {user?.bio && (
          <div className="mb-6">
            <p className="text-white leading-relaxed">{user.bio}</p>
          </div>
        )}

        {/* 其他信息 */}
        <div className="space-y-2 mb-6 text-gray-400 text-sm">
          {user?.location && (
            <div className="flex items-center space-x-2">
              <MapPin size={16} />
              <span>{user.location}</span>
            </div>
          )}
          {user?.website && (
            <div className="flex items-center space-x-2">
              <LinkIcon size={16} />
              <a
                href={user.website}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-400 hover:underline"
              >
                {user.website}
              </a>
            </div>
          )}
          <div className="flex items-center space-x-2">
            <Calendar size={16} />
            <span>加入于 {formatDate(user?.CreatedTime || '')}</span>
          </div>
        </div>

        {/* 标签页 */}
        <div className="flex border-b border-gray-800 mb-4">
          <button
            onClick={() => setActiveTab('videos')}
            className={cn(
              'flex-1 py-3 text-center font-medium transition-colors relative',
              activeTab === 'videos'
                ? 'text-white'
                : 'text-gray-400 hover:text-gray-300'
            )}
          >
            <Grid size={20} className="mx-auto mb-1" />
            作品
            {activeTab === 'videos' && (
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-white" />
            )}
          </button>
          <button
            onClick={() => setActiveTab('liked')}
            className={cn(
              'flex-1 py-3 text-center font-medium transition-colors relative',
              activeTab === 'liked'
                ? 'text-white'
                : 'text-gray-400 hover:text-gray-300'
            )}
          >
            <Heart size={20} className="mx-auto mb-1" />
            喜欢
            {activeTab === 'liked' && (
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-white" />
            )}
          </button>
        </div>

        {/* 视频网格 */}
        <VideoGrid
          videos={[]} // 这里应该根据activeTab加载对应的视频数据
          isLoading={false}
        />
      </div>

      {/* 编辑资料弹窗 */}
      <AnimatePresence>
        {showEditModal && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
            onClick={() => setShowEditModal(false)}
          >
            <motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              onClick={(e) => e.stopPropagation()}
              className="bg-gray-900 rounded-lg p-6 w-full max-w-md max-h-[90vh] overflow-y-auto"
            >
              <h3 className="text-xl font-bold mb-4">编辑资料</h3>
              
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-200 mb-2">
                    用户名
                  </label>
                  <input
                    {...register('username', {
                      required: '请输入用户名',
                      minLength: { value: 2, message: '用户名至少2位' },
                    })}
                    type="text"
                    defaultValue={user?.username}
                    className="w-full px-4 py-2 bg-gray-800 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                  {errors.username && (
                    <p className="mt-1 text-sm text-red-400">{errors.username.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-200 mb-2">
                    个人简介
                  </label>
                  <textarea
                    {...register('bio')}
                    rows={3}
                    defaultValue={user?.bio}
                    className="w-full px-4 py-2 bg-gray-800 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                    placeholder="介绍一下自己..."
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-200 mb-2">
                    头像链接
                  </label>
                  <input
                    {...register('avatar')}
                    type="url"
                    defaultValue={user?.avatar}
                    className="w-full px-4 py-2 bg-gray-800 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="输入头像链接..."
                  />
                </div>

                <div className="flex space-x-3 pt-4">
                  <button
                    type="button"
                    onClick={() => setShowEditModal(false)}
                    className="flex-1 py-2 bg-gray-700 rounded-lg hover:bg-gray-600 transition-colors"
                  >
                    取消
                  </button>
                  <button
                    type="submit"
                    disabled={updateProfileMutation.isLoading}
                    className="flex-1 py-2 bg-blue-500 rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                  >
                    {updateProfileMutation.isLoading ? '保存中...' : '保存'}
                  </button>
                </div>
              </form>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  )
}