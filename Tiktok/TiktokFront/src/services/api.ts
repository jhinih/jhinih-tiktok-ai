import axios, { AxiosResponse, AxiosError } from 'axios'
import { tokenService } from './tokenService'
import { toast } from '@/utils/toast'
import type {
  ApiResponse,
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  VideoListResponse,
  FriendListResponse,
  GroupListResponse,
  UploadResponse,
  AIRequest,
  AIResponse,
  WSTicketResponse,
  CreateVideoRequest,
  LikeVideoRequest,
  CommentVideoRequest,
  AddFriendRequest,
  CreateCommunityRequest,
  JoinCommunityRequest,
  User,
  Comment
} from '@/types'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 自动添加有效的访问token
api.interceptors.request.use(
  async (config) => {
    // 对于登录、注册、发送验证码等不需要token的接口，跳过token检查
    const noAuthEndpoints = ['/login', '/register', '/send-code', '/refresh-token']
    const needsAuth = !noAuthEndpoints.some(endpoint => config.url?.includes(endpoint))
    
    if (needsAuth) {
      const validToken = await tokenService.getValidAccessToken()
      if (validToken) {
        config.headers.Authorization = `Bearer ${validToken}`
      }
    }
    
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误和token刷新
api.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    // 直接返回响应，让具体的API调用处理业务逻辑
    return response
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as any
    
    // 如果是401错误且不是刷新token的请求，尝试刷新token
    if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes('/refresh-token')) {
      originalRequest._retry = true
      
      try {
        const newToken = await tokenService.refreshAccessToken()
        if (newToken) {
          // 重新发送原始请求
          originalRequest.headers.Authorization = `Bearer ${newToken}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        console.error('Token刷新失败:', refreshError)
        // 刷新失败，清除认证状态并跳转到登录页
        tokenService.clearTokens()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }
    
    // 处理其他错误
    const response = error.response
    if (response?.data && typeof response.data === 'object' && 'message' in response.data) {
      toast.error(response.data.message as string)
    } else {
      toast.error('网络错误，请检查连接')
    }
    
    return Promise.reject(error)
  }
)

// 登录相关API
export const loginApi = (data: LoginRequest) => 
  api.post<ApiResponse<LoginResponse>>('/login/login', data)

export const registerApi = (data: RegisterRequest) => 
  api.post<ApiResponse<LoginResponse>>('/login/register', data)

export const sendCodeApi = (email: string) => 
  api.post<ApiResponse>('/login/send-code', { email })

export const refreshTokenApi = (rtoken: string) => 
  api.post<ApiResponse<{ atoken: string }>>('/login/refresh-token', { rtoken })

// 用户相关API
export const getUserInfoApi = (id: string) => 
  api.get<ApiResponse<{ user: User }>>('/user/get-user-info', { params: { id } })

export const getMyInfoApi = () => 
  api.get<ApiResponse<{ user: User }>>('/user/get-my-info')

export const getProfileApi = (id: string) => 
  api.get<ApiResponse<{ user: User }>>('/user/get-profile', { params: { id } })

export const setProfileApi = (data: Partial<User>) => 
  api.post<ApiResponse>('/user/set-profile', data)

// 视频相关API
export const getVideosApi = (params?: { page?: string; page_size?: string; order_by?: string }) => 
  api.get<ApiResponse<VideoListResponse>>('/videos', { params })

export const createVideoApi = (data: CreateVideoRequest) => 
  api.post<ApiResponse>('/videos/create-video', data)

export const likeVideoApi = (data: LikeVideoRequest) => 
  api.post<ApiResponse>('/videos/like-video', data)

export const getVideoLikesApi = (videoId: string) => 
  api.get<ApiResponse<{ video_likes: string }>>('/videos/get-video-likes', { params: { video_id: videoId } })

export const commentVideoApi = (data: CommentVideoRequest) => 
  api.post<ApiResponse>('/videos/comment-video', data)

export const getCommentsApi = (params: { id: string; is_video: boolean; before_id?: string; page?: string; page_size?: string; order_by?: string }) => 
  api.get<ApiResponse<{ comments: Comment[]; length: string }>>('/videos/get-comments', { params })

// 文件上传API
export const uploadFileApi = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return api.post<ApiResponse<UploadResponse>>('/file/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

// 好友关系API
export const addFriendApi = (data: AddFriendRequest) => 
  api.post<ApiResponse>('/contact/add-friend', data)

export const getFriendListApi = (data?: { user_id?: string; page?: string; page_size?: string; order_by?: string }) => 
  api.post<ApiResponse<FriendListResponse>>('/contact/get-friend-list', data)

export const getOnlineUsersApi = () => 
  api.get<ApiResponse<FriendListResponse>>('/contact/get-user-list-online')

// 群组相关API
export const createCommunityApi = (data: CreateCommunityRequest) => 
  api.post<ApiResponse>('/contact/create-community', data)

export const joinCommunityApi = (data: JoinCommunityRequest) => 
  api.post<ApiResponse>('/contact/join-community', data)

export const loadCommunityApi = (data?: { owner_id?: string }) => 
  api.post<ApiResponse<GroupListResponse>>('/contact/load-community', data)

export const getGroupUsersApi = (communityId: string) => 
  api.get<ApiResponse<FriendListResponse>>('/contact/get-group-users', { params: { community_id: communityId } })

// WebSocket相关API已移除 - 后端使用直接连接方式

// AI相关API
export const commonAIApi = (data: AIRequest) => 
  api.post<ApiResponse<AIResponse>>('/ai/common_ai', data)

export const videoAIApi = (data: AIRequest) => 
  api.post<ApiResponse<AIResponse>>('/ai/video_ai', data)

export const sendCodeAIApi = (data: AIRequest) => 
  api.post<ApiResponse<AIResponse>>('/ai/send_code_ai', data)

export const getUserInfoAIApi = (data: AIRequest) => 
  api.post<ApiResponse<AIResponse>>('/ai/get_user_info_ai', data)

export const allAIApi = (data: AIRequest) => 
  api.post<ApiResponse<AIResponse>>('/ai/ai', data)

// 发送消息的别名函数 - 先声明
const sendMessage = (data: any) => 
  api.post<ApiResponse>('/message/send', data)

// 导出缺失的API函数别名
export const getUserProfileApi = getUserInfoApi
export const updateUserProfileApi = setProfileApi
export const followUserApi = (userId: string) => 
  api.post<ApiResponse>('/user/follow', { user_id: userId })
export const unfollowUserApi = (userId: string) => 
  api.delete<ApiResponse>(`/user/follow/${userId}`)
export const sendMessageApi = sendMessage
export const getChatListApi = () => 
  api.get<ApiResponse<{ chats: any[] }>>('/chat/list')
export const getChatMessagesApi = (chatId: string) => 
  api.get<ApiResponse<{ messages: any[] }>>(`/chat/${chatId}/messages`)
export const createGroupApi = createCommunityApi
// 视频上传API - 分两步：1.上传文件 2.创建视频记录
export const uploadVideoApi = async (formData: FormData) => {
  // 第一步：上传视频文件
  const videoFile = formData.get('video') as File
  if (!videoFile) {
    throw new Error('没有选择视频文件')
  }

  const fileFormData = new FormData()
  fileFormData.append('file', videoFile)
  
  const uploadResponse = await api.post<ApiResponse<UploadResponse>>('/file/upload', fileFormData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })

  if (uploadResponse.data.code !== 200) {
    throw new Error(uploadResponse.data.message || '文件上传失败')
  }

  // 第二步：创建视频记录
  const videoData = {
    video_path: uploadResponse.data.data.url,
    cover_path: uploadResponse.data.data.url, // 暂时使用视频URL作为封面
    title: formData.get('title') as string,
    description: formData.get('description') as string || '',
    is_private: formData.get('is_private') === 'true',
    type: 'video'
  }

  return api.post<ApiResponse>('/videos/create-video', videoData)
}

export default api