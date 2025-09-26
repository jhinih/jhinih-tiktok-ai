import { useState, useRef, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { Upload, X, Play, Pause, Volume2, VolumeX, ArrowLeft } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import { toast } from '@/utils/toast'
import { uploadVideoApi } from '@/services/api'
import { useAuthStore } from '@/stores/authStore'
import { cn } from '@/utils'


interface VideoPreviewProps {
  file: File
  onRemove: () => void
}

function VideoPreview({ file, onRemove }: VideoPreviewProps) {
  const [isPlaying, setIsPlaying] = useState(false)
  const [isMuted, setIsMuted] = useState(true)
  const videoUrl = URL.createObjectURL(file)
  const videoRef = useRef<HTMLVideoElement>(null)

  const togglePlay = () => {
    const video = videoRef.current
    if (!video) return

    if (isPlaying) {
      video.pause()
    } else {
      video.play()
    }
    setIsPlaying(!isPlaying)
  }

  return (
    <div className="relative w-full aspect-[9/16] bg-black rounded-lg overflow-hidden">
      <video
        ref={videoRef}
        src={videoUrl}
        muted={isMuted}
        loop
        playsInline
        className="w-full h-full object-cover"
        onPlay={() => setIsPlaying(true)}
        onPause={() => setIsPlaying(false)}
      />
      
      {/* 控制按钮 */}
      <div className="absolute inset-0 flex items-center justify-center">
        <motion.button
          whileTap={{ scale: 0.9 }}
          onClick={togglePlay}
          className="w-16 h-16 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center"
        >
          {isPlaying ? (
            <Pause size={24} className="text-white" />
          ) : (
            <Play size={24} className="text-white ml-1" />
          )}
        </motion.button>
      </div>

      {/* 音量控制 */}
      <button
        onClick={() => {
          setIsMuted(!isMuted)
          if (videoRef.current) {
            videoRef.current.muted = !isMuted
          }
        }}
        className="absolute top-4 right-4 p-2 bg-white/20 backdrop-blur-sm rounded-full"
      >
        {isMuted ? (
          <VolumeX size={20} className="text-white" />
        ) : (
          <Volume2 size={20} className="text-white" />
        )}
      </button>

      {/* 删除按钮 */}
      <button
        onClick={onRemove}
        className="absolute top-4 left-4 p-2 bg-red-500/80 backdrop-blur-sm rounded-full hover:bg-red-500 transition-colors"
      >
        <X size={20} className="text-white" />
      </button>

      {/* 文件信息 */}
      <div className="absolute bottom-4 left-4 right-4">
        <div className="bg-black/50 backdrop-blur-sm rounded-lg p-3">
          <p className="text-white text-sm font-medium truncate">{file.name}</p>
          <p className="text-gray-300 text-xs">
            {(file.size / (1024 * 1024)).toFixed(2)} MB
          </p>
        </div>
      </div>
    </div>
  )
}

export default function UploadPage() {
  const navigate = useNavigate()
  const { user } = useAuthStore()
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [dragActive, setDragActive] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    is_private: 'public'
  })
  const [errors, setErrors] = useState<Record<string, string>>({})
  const fileInputRef = useRef<HTMLInputElement>(null)

  const uploadMutation = useMutation({
    mutationFn: uploadVideoApi,
    onSuccess: () => {
      toast.success('视频上传成功！')
      setFormData({ title: '', description: '', is_private: 'public' })
      setSelectedFile(null)
      navigate('/')
    },
    onError: (error: any) => {
      console.error('Upload error:', error)
      toast.error('上传失败，请重试')
    }
  })

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true)
    } else if (e.type === 'dragleave') {
      setDragActive(false)
    }
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setDragActive(false)

    const files = e.dataTransfer.files
    if (files && files[0]) {
      handleFileSelect(files[0])
    }
  }, [])

  const handleFileSelect = (file: File) => {
    // 验证文件类型
    if (!file.type.startsWith('video/')) {
      toast.error('请选择视频文件')
      return
    }

    // 验证文件大小 (100MB)
    if (file.size > 100 * 1024 * 1024) {
      toast.error('文件大小不能超过100MB')
      return
    }

    setSelectedFile(file)
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (files && files[0]) {
      handleFileSelect(files[0])
    }
  }

  const validateForm = () => {
    const newErrors: Record<string, string> = {}

    if (!formData.title.trim()) {
      newErrors.title = '请输入视频标题'
    } else if (formData.title.length > 100) {
      newErrors.title = '标题不能超过100个字符'
    }

    if (formData.description.length > 500) {
      newErrors.description = '描述不能超过500个字符'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!selectedFile || !user) {
      toast.error('请选择视频文件')
      return
    }

    if (!validateForm()) {
      return
    }

    const uploadFormData = new FormData()
    uploadFormData.append('video', selectedFile)
    uploadFormData.append('title', formData.title)
    uploadFormData.append('description', formData.description)
    uploadFormData.append('user_id', user.id)
    uploadFormData.append('is_private', formData.is_private === 'private' ? 'true' : 'false')

    uploadMutation.mutate(uploadFormData as any)
  }

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))
    // 清除对应字段的错误
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }))
    }
  }

  return (
    <div className="min-h-screen bg-black text-white">
      {/* 头部 */}
      <div className="sticky top-0 bg-black/80 backdrop-blur-lg border-b border-gray-800 p-4 z-10">
        <div className="flex items-center space-x-4">
          <button
            onClick={() => navigate(-1)}
            className="p-2 hover:bg-gray-800 rounded-full transition-colors"
          >
            <ArrowLeft size={24} />
          </button>
          <h1 className="text-xl font-bold">上传视频</h1>
        </div>
      </div>

      <div className="p-4 max-w-4xl mx-auto">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* 左侧：文件上传区域 */}
          <div className="space-y-6">
            {!selectedFile ? (
              <div
                onDragEnter={handleDrag}
                onDragLeave={handleDrag}
                onDragOver={handleDrag}
                onDrop={handleDrop}
                className={cn(
                  'border-2 border-dashed rounded-lg p-8 text-center transition-colors cursor-pointer',
                  dragActive
                    ? 'border-blue-500 bg-blue-500/10'
                    : 'border-gray-600 hover:border-gray-500'
                )}
                onClick={() => fileInputRef.current?.click()}
              >
                <div className="space-y-4">
                  <div className="w-16 h-16 bg-gray-800 rounded-full flex items-center justify-center mx-auto">
                    <Upload size={32} className="text-gray-400" />
                  </div>
                  <div>
                    <h3 className="text-lg font-semibold mb-2">选择视频上传</h3>
                    <p className="text-gray-400 mb-4">
                      拖拽视频文件到这里，或点击选择文件
                    </p>
                    <div className="text-sm text-gray-500 space-y-1">
                      <p>支持格式：MP4, MOV, AVI, MKV</p>
                      <p>最大文件大小：100MB</p>
                      <p>推荐分辨率：1080x1920 (9:16)</p>
                    </div>
                  </div>
                  <motion.button
                    whileHover={{ scale: 1.05 }}
                    whileTap={{ scale: 0.95 }}
                    className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                  >
                    选择文件
                  </motion.button>
                </div>
              </div>
            ) : (
              <VideoPreview
                file={selectedFile}
                onRemove={() => setSelectedFile(null)}
              />
            )}

            <input
              ref={fileInputRef}
              type="file"
              accept="video/*"
              onChange={handleFileInputChange}
              className="hidden"
            />
          </div>

          {/* 右侧：视频信息表单 */}
          <div className="space-y-6">
            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-200 mb-2">
                  视频标题 *
                </label>
                <input
                  type="text"
                  value={formData.title}
                  onChange={(e) => handleInputChange('title', e.target.value)}
                  className={cn(
                    'w-full px-4 py-3 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent',
                    errors.title && 'border-red-500'
                  )}
                  placeholder="给你的视频起个标题..."
                />
                {errors.title && (
                  <p className="mt-1 text-sm text-red-400">{errors.title}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-200 mb-2">
                  视频描述
                </label>
                <textarea
                  rows={4}
                  value={formData.description}
                  onChange={(e) => handleInputChange('description', e.target.value)}
                  className={cn(
                    'w-full px-4 py-3 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none',
                    errors.description && 'border-red-500'
                  )}
                  placeholder="描述一下你的视频内容..."
                />
                {errors.description && (
                  <p className="mt-1 text-sm text-red-400">{errors.description}</p>
                )}
              </div>

              {/* 隐私设置 */}
              <div>
                <label className="block text-sm font-medium text-gray-200 mb-3">
                  隐私设置
                </label>
                <div className="space-y-3">
                  <label className="flex items-center space-x-3 cursor-pointer">
                    <input
                      type="radio"
                      name="privacy"
                      value="public"
                      checked={formData.is_private === 'public'}
                      onChange={(e) => handleInputChange('is_private', e.target.value)}
                      className="w-4 h-4 text-blue-500 bg-gray-800 border-gray-600 focus:ring-blue-500"
                    />
                    <div>
                      <p className="text-white font-medium">公开</p>
                      <p className="text-gray-400 text-sm">所有人都可以看到你的视频</p>
                    </div>
                  </label>
                  <label className="flex items-center space-x-3 cursor-pointer">
                    <input
                      type="radio"
                      name="privacy"
                      value="friends"
                      checked={formData.is_private === 'friends'}
                      onChange={(e) => handleInputChange('is_private', e.target.value)}
                      className="w-4 h-4 text-blue-500 bg-gray-800 border-gray-600 focus:ring-blue-500"
                    />
                    <div>
                      <p className="text-white font-medium">好友可见</p>
                      <p className="text-gray-400 text-sm">只有你的好友可以看到</p>
                    </div>
                  </label>
                  <label className="flex items-center space-x-3 cursor-pointer">
                    <input
                      type="radio"
                      name="privacy"
                      value="private"
                      checked={formData.is_private === 'private'}
                      onChange={(e) => handleInputChange('is_private', e.target.value)}
                      className="w-4 h-4 text-blue-500 bg-gray-800 border-gray-600 focus:ring-blue-500"
                    />
                    <div>
                      <p className="text-white font-medium">私密</p>
                      <p className="text-gray-400 text-sm">只有你可以看到</p>
                    </div>
                  </label>
                </div>
              </div>

              {/* 上传按钮 */}
              <div className="flex space-x-4">
                <button
                  type="button"
                  onClick={() => navigate(-1)}
                  className="flex-1 py-3 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition-colors"
                >
                  取消
                </button>
                <motion.button
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                  type="submit"
                  disabled={!selectedFile || uploadMutation.isLoading}
                  className="flex-1 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {uploadMutation.isLoading ? '上传中...' : '发布视频'}
                </motion.button>
              </div>
            </form>

            {/* 上传进度 */}
            <AnimatePresence>
              {uploadMutation.isLoading && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -20 }}
                  className="bg-gray-800 rounded-lg p-4"
                >
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium">上传进度</span>
                    <span className="text-sm text-gray-400">处理中...</span>
                  </div>
                  <div className="w-full bg-gray-700 rounded-full h-2">
                    <div className="bg-blue-500 h-2 rounded-full animate-pulse" style={{ width: '60%' }} />
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </div>
        </div>
      </div>
    </div>
  )
}