// 用户相关类型 - 完全匹配后端User结构
export interface User {
  id: string
  username: string
  password?: string
  email: string
  avatar: string
  role: string  // 后端使用字符串类型
  bio: string
  CreatedTime: string  // 毫秒时间戳字符串
  UpdatedTime: string  // 毫秒时间戳字符串
  Phone: string
  ClientIp: string
  ClientPort: string
  LoginTime: string
  HeartbeatTime: string
  LoginOutTime: string  // 匹配后端字段名
  IsLogout: boolean
  DeviceInfo: string
  nickname?: string  // 兼容字段
  following_count?: number  // 关注数
  followers_count?: number  // 粉丝数
  likes_count?: number  // 获赞数
  location?: string  // 位置
  website?: string  // 网站
}

// 视频相关类型 - 完全匹配后端Video结构
export interface Video {
  id: string
  CreatedTime: string  // 毫秒时间戳字符串
  UpdatedTime: string  // 毫秒时间戳字符串
  title: string
  description: string
  url: string
  cover: string
  likes: string  // 后端使用字符串类型
  comments: string  // 后端使用字符串类型
  shares: string  // 后端使用字符串类型
  user_id: string
  PublishTime: string
  type: string
  is_private: boolean
}

// 评论相关类型
export interface Comment {
  id: string
  CreatedTime: string
  UpdatedTime: string
  content: string
  video_id: string
  father_id: string
  user_id: string
  owner_id: string
  likes: string
  comments: string
}

// 社区/群组相关类型
export interface Community {
  id: string
  CreatedTime: string
  UpdatedTime: string
  name: string
  owner_id: string
  owner_name: string
  img: string
  desc: string
  description?: string  // 兼容字段
  member_count?: number  // 兼容字段
}

// 消息相关类型 - 兼容前后端字段
export interface Message {
  id: string
  UserId: string
  UserName: string
  TargetId: string
  TargetName: string
  Type: number // 1私聊 2群聊 3心跳
  Media: number // 1文字 2表情包 3语音 4图片
  Content: string
  CreateTime: number
  ReadTime: number
  Pic: string
  Url: string
  Desc: string
  Amount: number
  // 前端兼容字段
  sender_id?: string
  sender_name?: string
  type?: string
  content?: string
  created_at?: string
}

// API响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 登录相关类型
export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  code: string
  password: string
  username: string
  avatar: string
}

export interface LoginResponse {
  atoken: string
  rtoken?: string
}

// 视频列表响应类型
export interface VideoListResponse {
  data: Video[]
  page: string
  page_size: string
  total: string
  has_more: boolean
}

// 好友列表响应类型
export interface FriendListResponse {
  users: User[]
}

// 群组列表响应类型
export interface GroupListResponse {
  groups: Community[]
}

// 文件上传响应类型
export interface UploadResponse {
  url: string
}

// AI聊天请求类型 - 匹配后端AIRequest/AIResponse
export interface AIRequest {
  ask: string
}

export interface AIResponse {
  anser: string  // 注意：后端拼写是anser不是answer
}

// WebSocket票据响应类型
export interface WSTicketResponse {
  ticket: string
}

// 视频创建请求类型
export interface CreateVideoRequest {
  video_path: string
  cover_path: string
  title: string
  description: string
  is_private: boolean
  type: string
  user_id: string
}

// 点赞请求类型
export interface LikeVideoRequest {
  video_id: string
  owner_id: string
  user_id: string
}

// 评论请求类型
export interface CommentVideoRequest {
  video_id: string
  content: string
  owner_id: string
  user_id: string
}

// 添加好友请求类型 - 匹配后端AddFriendRequest
export interface AddFriendRequest {
  user_id: string
  user_name: string  // 后端字段名是user_name不是username
}

// 创建群组请求类型
export interface CreateCommunityRequest {
  owner_name: string
  name: string
  icon: string
  desc: string
}

// 加入群组请求类型
export interface JoinCommunityRequest {
  community_id: string
}

// 聊天室类型
export interface ChatRoom {
  id: string
  name: string
  type: 'private' | 'group'
  online: boolean
  lastMessage?: Message | null
  unreadCount: number
}

// 更新资料请求类型
export interface UpdateProfileRequest {
  username?: string
  bio?: string
  avatar?: string
}

// 上传视频请求类型
export interface UploadVideoRequest {
  video_path: string
  cover_path: string
  title: string
  description: string
  is_private: boolean
  type: string
}