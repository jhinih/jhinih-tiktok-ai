# TikTok Frontend

一个现代化的TikTok前端应用，使用React + TypeScript + Vite构建，具有完整的短视频社交功能。

## 🚀 功能特性

### 核心功能
- **AI推荐算法集成** - 智能视频推荐和内容发现
- **短视频流展示** - 沉浸式全屏视频播放体验
- **好友关系管理** - 关注/取消关注，好友列表管理
- **群组聊天室** - 实时群组聊天功能
- **私信功能** - 一对一私人消息
- **视频上传** - 支持多格式视频上传和编辑
- **个人资料** - 完整的用户资料管理

### 技术特性
- **移动端优先** - 响应式设计，完美适配移动设备
- **实时通信** - WebSocket实时消息推送
- **状态管理** - Zustand轻量级状态管理
- **数据获取** - React Query智能数据缓存
- **动画效果** - Framer Motion流畅动画
- **类型安全** - 完整的TypeScript类型定义

## 🛠️ 技术栈

- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **样式**: Tailwind CSS
- **状态管理**: Zustand
- **数据获取**: TanStack React Query
- **路由**: React Router DOM
- **表单**: React Hook Form
- **动画**: Framer Motion
- **图标**: Lucide React
- **视频播放**: React Player
- **实时通信**: Socket.IO Client
- **通知**: React Hot Toast

## 📦 安装和运行

### 环境要求
- Node.js >= 16
- npm >= 8

### 安装依赖
```bash
cd Tiktok/TiktokFront
npm install
```

### 开发模式
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

### 预览生产版本
```bash
npm run preview
```

## 🏗️ 项目结构

```
src/
├── components/          # 可复用组件
│   ├── Layout.tsx      # 布局组件
│   ├── BottomNavigation.tsx  # 底部导航
│   └── ProtectedRoute.tsx    # 路由保护
├── pages/              # 页面组件
│   ├── HomePage.tsx    # 首页视频流
│   ├── DiscoverPage.tsx # 发现页面
│   ├── ChatPage.tsx    # 聊天页面
│   ├── UploadPage.tsx  # 视频上传
│   ├── ProfilePage.tsx # 个人资料
│   ├── LoginPage.tsx   # 登录页面
│   └── RegisterPage.tsx # 注册页面
├── services/           # API服务
│   ├── api.ts         # HTTP API
│   └── websocket.ts   # WebSocket服务
├── stores/            # 状态管理
│   ├── authStore.ts   # 认证状态
│   ├── videoStore.ts  # 视频状态
│   └── chatStore.ts   # 聊天状态
├── types/             # 类型定义
│   └── index.ts       # 全局类型
├── utils/             # 工具函数
│   └── index.ts       # 通用工具
├── App.tsx            # 应用根组件
├── main.tsx           # 应用入口
└── index.css          # 全局样式
```

## 🔧 配置

### 环境变量
创建 `.env` 文件：
```env
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_AI_API_URL=http://localhost:8081
```

### API集成
项目已配置与后端API完全适配：
- 用户认证 (`/api/user/login`, `/api/user/register`)
- 视频管理 (`/api/video/*`)
- 聊天功能 (`/api/chat/*`)
- AI推荐 (`/api/ai/*`)

## 📱 页面功能

### 首页 (/)
- 全屏短视频播放
- 上下滑动切换视频
- 点赞、评论、分享功能
- AI推荐算法驱动

### 发现页 (/discover)
- 热门趋势展示
- AI智能助手
- 视频搜索功能
- 话题和用户推荐

### 聊天页 (/chat)
- 聊天列表管理
- 实时消息收发
- 群组聊天创建
- 文件和图片分享

### 上传页 (/upload)
- 视频文件上传
- 视频预览播放
- 标题和描述编辑
- 隐私设置配置

### 个人资料 (/profile)
- 用户信息展示
- 资料编辑功能
- 作品和喜欢列表
- 关注关系管理

## 🎨 设计特色

### 移动端优先
- 响应式布局设计
- 触摸友好的交互
- 优化的移动端性能

### 现代化UI
- 深色主题设计
- 流畅的动画效果
- 直观的用户界面

### 性能优化
- 懒加载和虚拟滚动
- 图片和视频优化
- 智能缓存策略

## 🔐 安全特性

- JWT令牌认证
- 路由权限保护
- XSS防护
- 安全的文件上传

## 🌐 浏览器支持

- Chrome >= 88
- Firefox >= 85
- Safari >= 14
- Edge >= 88

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request来改进项目。

## 📞 联系

如有问题或建议，请通过以下方式联系：
- 创建Issue
- 发送邮件

---

**注意**: 确保后端服务正在运行，并且API端点配置正确。