import { useState, useRef, useEffect } from 'react'
import { motion } from 'framer-motion'
import { Play, Pause, Volume2, VolumeX, AlertCircle } from 'lucide-react'
import type { Video } from '@/types'

interface VideoPlayerProps {
  video: Video
  isActive: boolean
  isMuted: boolean
  onToggleMute: () => void
  className?: string
}

export default function VideoPlayer({ 
  video, 
  isActive, 
  isMuted, 
  onToggleMute,
  className = ""
}: VideoPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(false)
  const [hasError, setHasError] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [showControls, setShowControls] = useState(false)
  const videoRef = useRef<HTMLVideoElement>(null)

  // 检查视频URL是否有效
  const isValidVideoUrl = (url: string) => {
    if (!url) return false
    
    // 检查是否是有效的URL格式
    try {
      new URL(url)
      return true
    } catch {
      // 如果不是完整URL，检查是否是相对路径
      return url.startsWith('/') || url.includes('.')
    }
  }

  useEffect(() => {
    const video = videoRef.current
    if (!video) return

    if (isActive) {
      video.currentTime = 0
      // 添加延迟以确保DOM已经准备好
      const timer = setTimeout(() => {
        if (video && document.contains(video)) {
          video.play().then(() => {
            setIsPlaying(true)
            setIsLoading(false)
          }).catch((error) => {
            console.error('视频播放失败:', error)
            setHasError(true)
            setIsLoading(false)
          })
        }
      }, 100)
      
      return () => clearTimeout(timer)
    } else {
      video.pause()
      setIsPlaying(false)
    }
  }, [isActive])

  useEffect(() => {
    const video = videoRef.current
    if (!video) return

    video.muted = isMuted
  }, [isMuted])

  const handleVideoClick = () => {
    const video = videoRef.current
    if (!video || hasError) return
    
    if (isPlaying) {
      video.pause()
      setIsPlaying(false)
    } else {
      video.play().then(() => {
        setIsPlaying(true)
      }).catch((error) => {
        console.error('视频播放失败:', error)
        setHasError(true)
      })
    }
    
    setShowControls(true)
    
    // 3秒后隐藏控制按钮
    setTimeout(() => {
      setShowControls(false)
    }, 3000)
  }

  const handleLoadedData = () => {
    setIsLoading(false)
    setHasError(false)
  }

  const handleError = () => {
    console.error('视频加载失败:', video.url)
    setHasError(true)
    setIsLoading(false)
  }

  const handlePlay = () => {
    setIsPlaying(true)
  }

  const handlePause = () => {
    setIsPlaying(false)
  }

  // 如果视频URL无效，显示错误状态
  if (!isValidVideoUrl(video.url)) {
    return (
      <div className={`relative w-full h-full bg-gray-900 flex items-center justify-center ${className}`}>
        {/* 封面图片作为fallback */}
        {video.cover && (
          <img
            src={video.cover}
            alt={video.title}
            className="absolute inset-0 w-full h-full object-cover"
          />
        )}
        
        <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
          <div className="text-center text-white">
            <AlertCircle size={48} className="mx-auto mb-4 text-gray-400" />
            <p className="text-lg font-medium mb-2">视频暂时无法播放</p>
            <p className="text-sm text-gray-300">URL: {video.url}</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className={`relative w-full h-full bg-black ${className}`}>
      {/* 加载状态 */}
      {isLoading && (
        <div className="absolute inset-0 bg-black flex items-center justify-center z-10">
          <div className="text-white text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
            <p>加载中...</p>
          </div>
        </div>
      )}

      {/* 错误状态 */}
      {hasError && (
        <div className="absolute inset-0 bg-gray-900 flex items-center justify-center z-10">
          {/* 显示封面作为fallback */}
          {video.cover && (
            <img
              src={video.cover}
              alt={video.title}
              className="absolute inset-0 w-full h-full object-cover"
            />
          )}
          
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
            <div className="text-center text-white">
              <AlertCircle size={48} className="mx-auto mb-4 text-gray-400" />
              <p className="text-lg font-medium mb-2">视频加载失败</p>
              <p className="text-sm text-gray-300 mb-4">URL: {video.url}</p>
              <button
                onClick={() => {
                  setHasError(false)
                  setIsLoading(true)
                  const videoEl = videoRef.current
                  if (videoEl) {
                    videoEl.load()
                  }
                }}
                className="px-4 py-2 bg-blue-500 rounded-lg hover:bg-blue-600 transition-colors"
              >
                重试
              </button>
            </div>
          </div>
        </div>
      )}

      {/* HTML5 视频播放器 */}
      <video
        ref={videoRef}
        src={video.url}
        poster={video.cover}
        muted={isMuted}
        loop
        playsInline
        webkit-playsinline="true"
        preload="metadata"
        className="absolute inset-0 w-full h-full object-cover"
        onLoadedData={handleLoadedData}
        onError={handleError}
        onPlay={handlePlay}
        onPause={handlePause}
        onCanPlay={() => setIsLoading(false)}
        style={{ display: hasError ? 'none' : 'block' }}
      />

      {/* 点击区域 */}
      <div
        className="absolute inset-0 z-20"
        onClick={handleVideoClick}
      />

      {/* 播放/暂停控制 */}
      {showControls && !hasError && (
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          exit={{ opacity: 0, scale: 0.8 }}
          className="absolute inset-0 flex items-center justify-center z-30 pointer-events-none"
        >
          <div className="p-4 rounded-full bg-black/50 backdrop-blur-sm">
            {isPlaying ? (
              <Pause size={32} className="text-white" />
            ) : (
              <Play size={32} className="text-white ml-1" />
            )}
          </div>
        </motion.div>
      )}

      {/* 音量控制 */}
      <motion.button
        whileTap={{ scale: 0.9 }}
        onClick={(e) => {
          e.stopPropagation()
          onToggleMute()
        }}
        className="absolute top-4 right-4 p-2 rounded-full bg-black/30 backdrop-blur-sm z-30"
      >
        {isMuted ? (
          <VolumeX size={20} className="text-white" />
        ) : (
          <Volume2 size={20} className="text-white" />
        )}
      </motion.button>

      {/* 调试信息 (开发环境) */}
      {import.meta.env.DEV && (
        <div className="absolute bottom-4 left-4 bg-black/50 text-white text-xs p-2 rounded z-30">
          <div>URL: {video.url}</div>
          <div>Playing: {isPlaying ? 'Yes' : 'No'}</div>
          <div>Active: {isActive ? 'Yes' : 'No'}</div>
          <div>Muted: {isMuted ? 'Yes' : 'No'}</div>
          <div>Error: {hasError ? 'Yes' : 'No'}</div>
        </div>
      )}
    </div>
  )
}