import { create } from 'zustand'
import type { Video } from '@/types'

interface VideoState {
  videos: Video[]
  currentVideoIndex: number
  isPlaying: boolean
  isMuted: boolean
  currentVideo: Video | null
  setVideos: (videos: Video[]) => void
  addVideos: (videos: Video[]) => void
  setCurrentVideoIndex: (index: number) => void
  setIsPlaying: (playing: boolean) => void
  setIsMuted: (muted: boolean) => void
  nextVideo: () => void
  previousVideo: () => void
  getCurrentVideo: () => Video | null
  likeVideo: (videoId: string) => void
  updateVideoStats: (videoId: string, stats: Partial<Pick<Video, 'likes' | 'comments' | 'shares'>>) => void
}

export const useVideoStore = create<VideoState>((set, get) => ({
  videos: [],
  currentVideoIndex: 0,
  isPlaying: false,
  isMuted: false,
  currentVideo: null,

  setVideos: (videos) => {
    set({ 
      videos, 
      currentVideoIndex: 0,
      currentVideo: videos[0] || null 
    })
  },

  addVideos: (newVideos) => {
    const { videos } = get()
    set({ videos: [...videos, ...newVideos] })
  },

  setCurrentVideoIndex: (index) => {
    const { videos } = get()
    if (index >= 0 && index < videos.length) {
      set({ 
        currentVideoIndex: index,
        currentVideo: videos[index]
      })
    }
  },

  setIsPlaying: (playing) => {
    set({ isPlaying: playing })
  },

  setIsMuted: (muted) => {
    set({ isMuted: muted })
  },

  nextVideo: () => {
    const { videos, currentVideoIndex } = get()
    const nextIndex = (currentVideoIndex + 1) % videos.length
    set({ 
      currentVideoIndex: nextIndex,
      currentVideo: videos[nextIndex]
    })
  },

  previousVideo: () => {
    const { videos, currentVideoIndex } = get()
    const prevIndex = currentVideoIndex === 0 ? videos.length - 1 : currentVideoIndex - 1
    set({ 
      currentVideoIndex: prevIndex,
      currentVideo: videos[prevIndex]
    })
  },

  getCurrentVideo: () => {
    const { videos, currentVideoIndex } = get()
    return videos[currentVideoIndex] || null
  },

  likeVideo: (videoId) => {
    const { videos } = get()
    const updatedVideos = videos.map(video => {
      if (video.id === videoId) {
        const currentLikes = parseInt(video.likes) || 0
        return {
          ...video,
          likes: (currentLikes + 1).toString()
        }
      }
      return video
    })
    set({ videos: updatedVideos })
  },

  updateVideoStats: (videoId, stats) => {
    const { videos } = get()
    const updatedVideos = videos.map(video => {
      if (video.id === videoId) {
        return { ...video, ...stats }
      }
      return video
    })
    set({ videos: updatedVideos })
  },
}))