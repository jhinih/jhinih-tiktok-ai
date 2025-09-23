import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Search, TrendingUp, Hash, Music, Users } from 'lucide-react'
import { motion } from 'framer-motion'
import { getVideosApi, allAIApi } from '@/services/api'
import { formatNumber } from '@/utils'
import type { Video } from '@/types'

interface TrendingItem {
  id: string
  title: string
  count: string
  type: 'hashtag' | 'music' | 'user'
}

const mockTrending: TrendingItem[] = [
  { id: '1', title: '#搞笑视频', count: '1.2M', type: 'hashtag' },
  { id: '2', title: '#舞蹈挑战', count: '890K', type: 'hashtag' },
  { id: '3', title: '热门音乐', count: '2.1M', type: 'music' },
  { id: '4', title: '#美食制作', count: '650K', type: 'hashtag' },
  { id: '5', title: '推荐用户', count: '340K', type: 'user' },
]

export default function DiscoverPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [aiQuery, setAiQuery] = useState('')
  const [aiResponse, setAiResponse] = useState('')
  const [isAiLoading, setIsAiLoading] = useState(false)
  const [activeTab, setActiveTab] = useState<'trending' | 'ai'>('trending')

  const { data: videosData, isLoading } = useQuery({
    queryKey: ['discover-videos'],
    queryFn: () => getVideosApi({ 
      page: '1', 
      page_size: '20',
      order_by: 'popular'
    }),
  })

  const videos = videosData?.data.data.data || []

  const handleAiSearch = async () => {
    if (!aiQuery.trim()) return

    setIsAiLoading(true)
    try {
      const response = await allAIApi({ ask: aiQuery })
      setAiResponse(response.data.data.anser)  // 修正字段名
    } catch (error) {
      setAiResponse('AI服务暂时不可用，请稍后再试。')
    } finally {
      setIsAiLoading(false)
    }
  }

  const getIcon = (type: string) => {
    switch (type) {
      case 'hashtag':
        return <Hash size={20} className="text-blue-400" />
      case 'music':
        return <Music size={20} className="text-purple-400" />
      case 'user':
        return <Users size={20} className="text-green-400" />
      default:
        return <TrendingUp size={20} className="text-gray-400" />
    }
  }

  return (
    <div className="min-h-screen bg-black text-white">
      {/* 搜索栏 */}
      <div className="sticky top-0 bg-black/80 backdrop-blur-lg border-b border-gray-800 p-4 z-10">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={20} />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="搜索视频、用户、话题..."
            className="w-full pl-10 pr-4 py-3 bg-gray-800 rounded-full text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      {/* 标签页 */}
      <div className="flex border-b border-gray-800">
        <button
          onClick={() => setActiveTab('trending')}
          className={`flex-1 py-4 text-center font-medium transition-colors ${
            activeTab === 'trending'
              ? 'text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-gray-300'
          }`}
        >
          热门发现
        </button>
        <button
          onClick={() => setActiveTab('ai')}
          className={`flex-1 py-4 text-center font-medium transition-colors ${
            activeTab === 'ai'
              ? 'text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-gray-300'
          }`}
        >
          AI助手
        </button>
      </div>

      {activeTab === 'trending' ? (
        <div className="p-4">
          {/* 热门趋势 */}
          <div className="mb-8">
            <h2 className="text-xl font-bold mb-4 flex items-center">
              <TrendingUp className="mr-2 text-red-500" size={24} />
              热门趋势
            </h2>
            <div className="space-y-3">
              {mockTrending.map((item, index) => (
                <motion.div
                  key={item.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.1 }}
                  className="flex items-center justify-between p-3 bg-gray-900 rounded-lg hover:bg-gray-800 transition-colors cursor-pointer"
                >
                  <div className="flex items-center space-x-3">
                    <span className="text-gray-400 font-mono text-sm w-6">
                      {(index + 1).toString().padStart(2, '0')}
                    </span>
                    {getIcon(item.type)}
                    <span className="font-medium">{item.title}</span>
                  </div>
                  <span className="text-gray-400 text-sm">{item.count}</span>
                </motion.div>
              ))}
            </div>
          </div>

          {/* 热门视频网格 */}
          <div className="mb-8">
            <h2 className="text-xl font-bold mb-4">热门视频</h2>
            {isLoading ? (
              <div className="grid grid-cols-2 gap-2">
                {Array.from({ length: 6 }).map((_, i) => (
                  <div key={i} className="aspect-[9/16] bg-gray-800 rounded-lg animate-pulse" />
                ))}
              </div>
            ) : (
              <div className="grid grid-cols-2 gap-2">
                {videos.slice(0, 6).map((video: Video) => (
                  <motion.div
                    key={video.id}
                    whileHover={{ scale: 1.02 }}
                    whileTap={{ scale: 0.98 }}
                    className="relative aspect-[9/16] bg-gray-800 rounded-lg overflow-hidden cursor-pointer group"
                  >
                    <img
                      src={video.cover || '/placeholder-video.jpg'}
                      alt={video.title}
                      className="w-full h-full object-cover"
                    />
                    <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
                    <div className="absolute bottom-2 left-2 right-2">
                      <p className="text-white text-sm font-medium line-clamp-2 mb-1">
                        {video.title}
                      </p>
                      <div className="flex items-center space-x-3 text-xs text-gray-300">
                        <span className="flex items-center">
                          ❤️ {formatNumber(video.likes)}
                        </span>
                        <span className="flex items-center">
                          💬 {formatNumber(video.comments)}
                        </span>
                      </div>
                    </div>
                    <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                      <div className="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                        <div className="w-0 h-0 border-l-[8px] border-l-white border-y-[6px] border-y-transparent ml-1" />
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            )}
          </div>
        </div>
      ) : (
        <div className="p-4">
          {/* AI助手 */}
          <div className="mb-6">
            <h2 className="text-xl font-bold mb-4 flex items-center">
              🤖 AI智能助手
            </h2>
            <p className="text-gray-400 mb-4">
              我可以帮您推荐视频、回答问题、提供创作建议等。试试问我一些问题吧！
            </p>
            
            <div className="space-y-4">
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={aiQuery}
                  onChange={(e) => setAiQuery(e.target.value)}
                  placeholder="问我任何问题..."
                  className="flex-1 px-4 py-3 bg-gray-800 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  onKeyPress={(e) => e.key === 'Enter' && handleAiSearch()}
                />
                <motion.button
                  whileTap={{ scale: 0.95 }}
                  onClick={handleAiSearch}
                  disabled={isAiLoading || !aiQuery.trim()}
                  className="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {isAiLoading ? '思考中...' : '提问'}
                </motion.button>
              </div>

              {/* 快捷问题 */}
              <div className="flex flex-wrap gap-2">
                {[
                  '推荐一些热门视频',
                  '如何制作有趣的短视频',
                  '今天有什么热门话题',
                  '帮我想个创意点子'
                ].map((question) => (
                  <button
                    key={question}
                    onClick={() => setAiQuery(question)}
                    className="px-3 py-1 bg-gray-700 text-gray-300 rounded-full text-sm hover:bg-gray-600 transition-colors"
                  >
                    {question}
                  </button>
                ))}
              </div>

              {/* AI回复 */}
              {aiResponse && (
                <motion.div
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  className="p-4 bg-gray-800 rounded-lg border-l-4 border-blue-500"
                >
                  <div className="flex items-start space-x-3">
                    <div className="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center flex-shrink-0">
                      🤖
                    </div>
                    <div className="flex-1">
                      <p className="text-white whitespace-pre-wrap">{aiResponse}</p>
                    </div>
                  </div>
                </motion.div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}