import { create } from 'zustand'
import type { Message, User, Community, ChatRoom } from '@/types'

interface ChatState {
  messages: Message[]
  friends: User[]
  groups: Community[]
  currentChat: {
    type: 'friend' | 'group' | null
    id: string | null
    name: string | null
  }
  activeChat: ChatRoom | null
  onlineUsers: User[]
  isConnected: boolean
  
  // Actions
  setMessages: (messages: Message[]) => void
  addMessage: (message: Message) => void
  setFriends: (friends: User[]) => void
  setGroups: (groups: Community[]) => void
  setCurrentChat: (type: 'friend' | 'group' | null, id: string | null, name: string | null) => void
  setActiveChat: (chat: ChatRoom | null) => void
  setOnlineUsers: (users: User[]) => void
  setIsConnected: (connected: boolean) => void
  addFriend: (friend: User) => void
  addGroup: (group: Community) => void
  updateMessageStatus: (messageId: string, status: 'sent' | 'delivered' | 'read') => void
}

export const useChatStore = create<ChatState>((set, get) => ({
  messages: [],
  friends: [],
  groups: [],
  currentChat: {
    type: null,
    id: null,
    name: null,
  },
  activeChat: null,
  onlineUsers: [],
  isConnected: false,

  setMessages: (messages) => {
    set({ messages })
  },

  addMessage: (message) => {
    const { messages } = get()
    set({ messages: [...messages, message] })
  },

  setFriends: (friends) => {
    set({ friends })
  },

  setGroups: (groups) => {
    set({ groups })
  },

  setCurrentChat: (type, id, name) => {
    set({
      currentChat: { type, id, name }
    })
  },

  setActiveChat: (chat) => {
    set({ activeChat: chat })
  },

  setOnlineUsers: (users) => {
    set({ onlineUsers: users })
  },

  setIsConnected: (connected) => {
    set({ isConnected: connected })
  },

  addFriend: (friend) => {
    const { friends } = get()
    if (!friends.find(f => f.id === friend.id)) {
      set({ friends: [...friends, friend] })
    }
  },

  addGroup: (group) => {
    const { groups } = get()
    if (!groups.find(g => g.id === group.id)) {
      set({ groups: [...groups, group] })
    }
  },

  updateMessageStatus: (messageId, status) => {
    const { messages } = get()
    const updatedMessages = messages.map(msg => {
      if (msg.id === messageId) {
        return { ...msg, status }
      }
      return msg
    })
    set({ messages: updatedMessages })
  },
}))