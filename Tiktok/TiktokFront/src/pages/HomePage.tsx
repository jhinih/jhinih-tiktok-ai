import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { Heart, MessageCircle, Share, MoreHorizontal } from 'lucide-react'
import { getVideosApi, likeVideoApi } from '@/services/api'
import { useVideoStore } from '@/stores/videoStore'
import { useAuthStore } from '@/stores/authStore'
import { formatNumber, formatTime } from '@/utils'
import VideoPlayer from '@/components/VideoPlayer'
import AccountManager from '@/components/AccountManager'
import { toast } from '@/utils/toast'
import type { Video } from '@/types'



interface VideoCardProps {
  video: Video
  isActive: boolean
  isMuted: boolean
  onToggleMute: () => void
}

function VideoCard({ video, isActive, isMuted, onToggleMute }: VideoCardProps) {
  const [liked, setLiked] = useState(false)
  const [likeCount, setLikeCount] = useState(parseInt(video.likes) || 0)
  const { user } = useAuthStore()

  const handleLike = async () => {
    if (!user) {
      toast.error('请先登录')
      return
    }

    try {
      await likeVideoApi({
        video_id: video.id,
        owner_id: video.user_id,
        user_id: user.id,
      })
      
      setLiked(!liked)
      setLikeCount(prev => liked ? prev - 1 : prev + 1)
      
      if (!liked) {
        toast.success('点赞成功', { icon: '❤️' })
      }
    } catch (error) {
      toast.error('操作失败')
    }
  }

  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: video.title,
        text: video.description,
        url: window.location.href,
      })
    } else {
      navigator.clipboard.writeText(window.location.href)
      toast.success('链接已复制到剪贴板')
    }
  }

  return (
    <div className="relative w-full h-screen bg-black overflow-hidden">
      {/* 视频播放器 */}
      <VideoPlayer
        video={video}
        isActive={isActive}
        isMuted={isMuted}
        onToggleMute={onToggleMute}
        className="absolute inset-0"
      />
      
      {/* 视频遮罩 */}
      <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent pointer-events-none" />
      
      {/* 右侧操作栏 - 改善样式 */}
      <div className="absolute right-4 bottom-32 flex flex-col items-center space-y-4 z-30">
        <motion.button
          whileTap={{ scale: 0.9 }}
          onClick={handleLike}
          className="flex flex-col items-center space-y-2"
        >
          <div className={`p-4 rounded-full ${liked ? 'bg-red-500 shadow-lg shadow-red-500/30' : 'bg-black/30'} backdrop-blur-sm border border-white/10`}>
            <Heart 
              size={28} 
              className={liked ? 'text-white fill-current' : 'text-white'} 
            />
          </div>
          <span className="text-white text-sm font-bold drop-shadow-lg">
            {formatNumber(likeCount)}
          </span>
        </motion.button>

        <motion.button
          whileTap={{ scale: 0.9 }}
          className="flex flex-col items-center space-y-2"
        >
          <div className="p-4 rounded-full bg-black/30 backdrop-blur-sm border border-white/10">
            <MessageCircle size={28} className="text-white" />
          </div>
          <span className="text-white text-sm font-bold drop-shadow-lg">
            {formatNumber(video.comments)}
          </span>
        </motion.button>

        <motion.button
          whileTap={{ scale: 0.9 }}
          onClick={handleShare}
          className="flex flex-col items-center space-y-2"
        >
          <div className="p-4 rounded-full bg-black/30 backdrop-blur-sm border border-white/10">
            <Share size={28} className="text-white" />
          </div>
          <span className="text-white text-sm font-bold drop-shadow-lg">
            {formatNumber(video.shares)}
          </span>
        </motion.button>

        <motion.button
          whileTap={{ scale: 0.9 }}
          className="flex flex-col items-center space-y-2"
        >
          <div className="p-4 rounded-full bg-black/30 backdrop-blur-sm border border-white/10">
            <MoreHorizontal size={28} className="text-white" />
          </div>
        </motion.button>

        {/* 音乐图标 */}
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 3, repeat: Infinity, ease: "linear" }}
          className="mt-4"
        >
          <div className="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 border-2 border-white/20 flex items-center justify-center">
            <span className="text-white text-xs">♪</span>
          </div>
        </motion.div>
      </div>

      {/* 底部信息 - 改善排版 */}
      <div className="absolute bottom-0 left-0 right-0 z-30">
        <div className="bg-gradient-to-t from-black/80 via-black/40 to-transparent p-6 pb-24">
          <div className="space-y-4">
            {/* 用户信息 */}
            <div className="flex items-center space-x-3">
              <div className="w-12 h-12 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 overflow-hidden border-2 border-white/20">
                <div className="w-full h-full flex items-center justify-center">
                  <span className="text-white font-bold text-lg">
                    {video.user_id.charAt(0).toUpperCase()}
                  </span>
                </div>
              </div>
              <div className="flex-1">
                <p className="text-white font-bold text-base">@{video.user_id}</p>
                <p className="text-gray-300 text-sm">{formatTime(video.CreatedTime || '')}</p>
              </div>
              <button className="px-4 py-1.5 bg-red-500 text-white text-sm font-medium rounded-full hover:bg-red-600 transition-colors">
                关注
              </button>
            </div>
            
            {/* 视频信息 */}
            <div className="space-y-2">
              <h3 className="text-white font-bold text-xl leading-tight">{video.title}</h3>
              {video.description && (
                <p className="text-gray-200 text-base leading-relaxed line-clamp-3">
                  {video.description}
                </p>
              )}
            </div>

            {/* 标签和话题 */}
            <div className="flex flex-wrap gap-2">
              <span className="px-2 py-1 bg-white/10 text-white text-sm rounded-full">
                #{video.type || 'video'}
              </span>
              <span className="px-2 py-1 bg-white/10 text-white text-sm rounded-full">
                #推荐
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default function HomePage() {
  const { videos, currentVideoIndex, isMuted, setVideos, setCurrentVideoIndex, setIsMuted } = useVideoStore()
  const [page, setPage] = useState(1)
  const [isRefreshing, setIsRefreshing] = useState(false)
  const [pullDistance, setPullDistance] = useState(0)

  // 获取视频数据
  const { isLoading, error } = useQuery({
    queryKey: ['videos', page],
    queryFn: () => getVideosApi({ 
      page: page.toString(), 
      page_size: '10',
      order_by: 'CreatedTime'
    }),
    onSuccess: (response) => {
      if (response.data.code === 200 && response.data.data?.data) {
        const newVideos = response.data.data.data
        if (page === 1) {
          setVideos(newVideos)
        } else {
          setVideos([...videos, ...newVideos])
        }
      }
    },
    onError: (error) => {
      console.error('获取视频失败:', error)
      toast.error('获取视频失败，请检查网络连接')
    },
    retry: 2,
    retryDelay: 1000
  })

  const handleScroll = (direction: 'up' | 'down') => {
    if (videos.length === 0) return
    
    if (direction === 'down' && currentVideoIndex < videos.length - 1) {
      setCurrentVideoIndex(currentVideoIndex + 1)
    } else if (direction === 'up' && currentVideoIndex > 0) {
      setCurrentVideoIndex(currentVideoIndex - 1)
    }
    
    // 预加载更多视频
    if (currentVideoIndex >= videos.length - 3 && !isLoading) {
      setPage(prev => prev + 1)
    }
  }

  // 刷新视频数据
  const refreshVideos = async () => {
    setIsRefreshing(true)
    try {
      setPage(1)
      setCurrentVideoIndex(0)
      // 触发重新获取数据
      window.location.reload()
    } catch (error) {
      console.error('刷新失败:', error)
      toast.error('刷新失败，请重试')
    } finally {
      setIsRefreshing(false)
      setPullDistance(0)
    }
  }

  // 键盘控制
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'ArrowUp':
        case 'w':
        case 'W':
          e.preventDefault()
          handleScroll('up')
          break
        case 'ArrowDown':
        case 's':
        case 'S':
          e.preventDefault()
          handleScroll('down')
          break
        case ' ':
        case 'Spacebar':
          e.preventDefault()
          setIsMuted(!isMuted)
          break
        case 'r':
        case 'R':
          e.preventDefault()
          refreshVideos()
          break
        case 'f':
        case 'F':
          e.preventDefault()
          // 全屏切换
          if (document.fullscreenElement) {
            document.exitFullscreen()
          } else {
            document.documentElement.requestFullscreen()
          }
          break
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [currentVideoIndex, videos.length, isMuted, setIsMuted, handleScroll, refreshVideos])

  // 改进的触摸控制
  useEffect(() => {
    let startY = 0
    let startX = 0
    let currentY = 0
    let isScrolling = false
    let isPulling = false

    const handleTouchStart = (e: TouchEvent) => {
      startY = e.touches[0].clientY
      startX = e.touches[0].clientX
      currentY = startY
      isScrolling = false
      isPulling = false
    }

    const handleTouchMove = (e: TouchEvent) => {
      if (!startY) return
      
      currentY = e.touches[0].clientY
      const currentX = e.touches[0].clientX
      const diffY = currentY - startY
      const diffX = Math.abs(currentX - startX)
      
      // 如果水平滑动距离大于垂直滑动，忽略
      if (diffX > Math.abs(diffY)) return
      
      // 检测下拉刷新
      if (diffY > 0 && currentVideoIndex === 0 && !isScrolling) {
        isPulling = true
        e.preventDefault()
        const distance = Math.min(diffY * 0.5, 100) // 限制最大下拉距离
        setPullDistance(distance)
      }
    }

    const handleTouchEnd = (e: TouchEvent) => {
      if (!startY) return
      
      const endY = e.changedTouches[0].clientY
      const diff = startY - endY
      const absDiff = Math.abs(diff)
      
      // 处理下拉刷新
      if (isPulling && pullDistance > 60) {
        refreshVideos()
        return
      }
      
      setPullDistance(0)
      
      // 处理视频切换 - 增加灵敏度
      if (absDiff > 30 && !isPulling) { // 降低阈值从50到30
        if (diff > 0) {
          // 向上滑动，下一个视频
          handleScroll('down')
        } else {
          // 向下滑动，上一个视频
          handleScroll('up')
        }
      }
      
      // 重置状态
      startY = 0
      currentY = 0
      isScrolling = false
      isPulling = false
    }

    const container = document.getElementById('video-container')
    if (container) {
      container.addEventListener('touchstart', handleTouchStart, { passive: false })
      container.addEventListener('touchmove', handleTouchMove, { passive: false })
      container.addEventListener('touchend', handleTouchEnd, { passive: false })
      
      return () => {
        container.removeEventListener('touchstart', handleTouchStart)
        container.removeEventListener('touchmove', handleTouchMove)
        container.removeEventListener('touchend', handleTouchEnd)
      }
    }
  }, [currentVideoIndex, videos.length, pullDistance, handleScroll, refreshVideos])

  if (isLoading && videos.length === 0) {
    return (
      <div className="h-screen bg-black flex items-center justify-center">
        <div className="text-white text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p>正在加载视频...</p>
        </div>
      </div>
    )
  }

  if (error && videos.length === 0) {
    return (
      <div className="h-screen bg-black flex items-center justify-center">
        <div className="text-white text-center">
          <div className="text-red-400 text-6xl mb-4">⚠️</div>
          <h2 className="text-xl font-bold mb-2">无法加载视频</h2>
          <p className="text-gray-400 mb-4">请检查网络连接或稍后重试</p>
          <button
            onClick={() => window.location.reload()}
            className="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            重新加载
          </button>
        </div>
      </div>
    )
  }

  if (videos.length === 0) {
    return (
      <div className="h-screen bg-black flex items-center justify-center">
        <div className="text-white text-center">
          <div className="text-gray-400 text-6xl mb-4">📹</div>
          <h2 className="text-xl font-bold mb-2">暂无视频</h2>
          <p className="text-gray-400">快去上传第一个视频吧！</p>
        </div>
      </div>
    )
  }

  return (
    <div id="video-container" className="relative h-screen overflow-hidden bg-black">
      {/* 下拉刷新指示器 */}
      {pullDistance > 0 && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="absolute top-0 left-0 right-0 z-50 flex items-center justify-center"
          style={{ transform: `translateY(${pullDistance - 60}px)` }}
        >
          <div className="bg-black/50 backdrop-blur-sm rounded-full px-4 py-2 flex items-center space-x-2">
            {isRefreshing ? (
              <>
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                <span className="text-white text-sm">刷新中...</span>
              </>
            ) : pullDistance > 60 ? (
              <>
                <span className="text-white text-lg">↻</span>
                <span className="text-white text-sm">松开刷新</span>
              </>
            ) : (
              <>
                <span className="text-white text-lg">↓</span>
                <span className="text-white text-sm">下拉刷新</span>
              </>
            )}
          </div>
        </motion.div>
      )}

      <AnimatePresence>
        {videos.map((video: Video, index: number) => (
          <motion.div
            key={video.id}
            initial={{ y: index > currentVideoIndex ? '100%' : index < currentVideoIndex ? '-100%' : 0 }}
            animate={{ y: index === currentVideoIndex ? 0 : index > currentVideoIndex ? '100%' : '-100%' }}
            exit={{ y: index < currentVideoIndex ? '-100%' : '100%' }}
            transition={{ type: 'tween', duration: 0.3 }}
            className="absolute inset-0"
            style={{ display: Math.abs(index - currentVideoIndex) <= 1 ? 'block' : 'none' }}
          >
            <VideoCard
              video={video}
              isActive={index === currentVideoIndex}
              isMuted={isMuted}
              onToggleMute={() => setIsMuted(!isMuted)}
            />
          </motion.div>
        ))}
      </AnimatePresence>



      {/* 使用说明 */}
      <div className="absolute top-4 left-4 bg-black/30 backdrop-blur-sm rounded-lg p-3 z-40">
        <p className="text-white text-sm">
          📱 滑动切换 | ⌨️ W/S或↑↓切换 | 空格静音 | R刷新 | F全屏
        </p>
      </div>

      {/* 顶部右侧控制区 */}
      <div className="absolute top-4 right-4 flex items-center space-x-3 z-40">
        {/* 视频计数器 */}
        {videos.length > 0 && (
          <div className="bg-black/30 backdrop-blur-sm rounded-lg px-3 py-2">
            <p className="text-white text-sm font-medium">
              {currentVideoIndex + 1} / {videos.length}
            </p>
          </div>
        )}
        
        {/* 账号管理 */}
        <AccountManager />
      </div>
    </div>
  )
}