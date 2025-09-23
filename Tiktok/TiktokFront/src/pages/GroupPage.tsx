import React, { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion } from 'framer-motion'
import { 
  ArrowLeftIcon, 
  UserGroupIcon, 
  PlusIcon,
  UserPlusIcon
} from '@heroicons/react/24/outline'
import { 
  loadCommunityApi, 
  getGroupUsersApi, 
  joinCommunityApi,
  createCommunityApi 
} from '@/services/api'
import { useAuthStore } from '@/stores/authStore'
import toast from 'react-hot-toast'
import type { Community, User } from '@/types'

const GroupPage: React.FC = () => {
  const { groupId } = useParams<{ groupId: string }>()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const { user } = useAuthStore()
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showJoinModal, setShowJoinModal] = useState(false)
  const [newGroupName, setNewGroupName] = useState('')
  const [newGroupDesc, setNewGroupDesc] = useState('')
  const [joinGroupId, setJoinGroupId] = useState('')

  // 获取群组列表
  const { data: groupsData, isLoading: groupsLoading } = useQuery({
    queryKey: ['groups'],
    queryFn: () => loadCommunityApi()
  })

  // 获取群组成员
  const { data: membersData, isLoading: membersLoading } = useQuery({
    queryKey: ['group-members', groupId],
    queryFn: () => groupId ? getGroupUsersApi(groupId) : null,
    enabled: !!groupId
  })

  // 创建群组
  const createGroupMutation = useMutation({
    mutationFn: createCommunityApi,
    onSuccess: () => {
      toast.success('群组创建成功')
      setShowCreateModal(false)
      setNewGroupName('')
      setNewGroupDesc('')
      queryClient.invalidateQueries({ queryKey: ['groups'] })
    },
    onError: (error: any) => {
      toast.error(error.message || '创建群组失败')
    }
  })

  // 加入群组
  const joinGroupMutation = useMutation({
    mutationFn: joinCommunityApi,
    onSuccess: () => {
      toast.success('加入群组成功')
      setShowJoinModal(false)
      setJoinGroupId('')
      queryClient.invalidateQueries({ queryKey: ['groups'] })
    },
    onError: (error: any) => {
      toast.error(error.message || '加入群组失败')
    }
  })

  const groups = Array.isArray(groupsData?.data?.data) ? groupsData.data.data : []
  const members = membersData?.data?.data?.users || []

  const handleCreateGroup = () => {
    if (!newGroupName.trim()) {
      toast.error('请输入群组名称')
      return
    }

    createGroupMutation.mutate({
      name: newGroupName,
      desc: newGroupDesc,
      owner_name: user?.username || '',
      icon: ''
    })
  }

  const handleJoinGroup = () => {
    if (!joinGroupId.trim()) {
      toast.error('请输入群组ID')
      return
    }

    joinGroupMutation.mutate({
      community_id: joinGroupId
    })
  }

  const handleGroupClick = (group: Community) => {
    navigate(`/group/${group.id}`)
  }

  const handleChatClick = (group: Community) => {
    navigate(`/chat?type=group&id=${group.id}`)
  }

  if (groupsLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 头部 */}
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-md mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <button
              onClick={() => navigate(-1)}
              className="p-2 hover:bg-gray-100 rounded-full transition-colors"
            >
              <ArrowLeftIcon className="w-5 h-5" />
            </button>
            <h1 className="text-lg font-semibold">
              {groupId ? '群组详情' : '我的群组'}
            </h1>
          </div>
          
          {!groupId && (
            <div className="flex items-center space-x-2">
              <button
                onClick={() => setShowJoinModal(true)}
                className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                title="加入群组"
              >
                <UserPlusIcon className="w-5 h-5" />
              </button>
              <button
                onClick={() => setShowCreateModal(true)}
                className="p-2 hover:bg-gray-100 rounded-full transition-colors"
                title="创建群组"
              >
                <PlusIcon className="w-5 h-5" />
              </button>
            </div>
          )}
        </div>
      </div>

      <div className="max-w-md mx-auto">
        {groupId ? (
          // 群组详情页面
          <div className="p-4">
            {membersLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-500"></div>
              </div>
            ) : (
              <div className="space-y-4">
                <div className="bg-white rounded-lg p-4 shadow-sm">
                  <h2 className="text-lg font-semibold mb-2">群组成员</h2>
                  <p className="text-gray-500 text-sm mb-4">共 {members.length} 人</p>
                  
                  <div className="space-y-3">
                    {members.map((member: User) => (
                      <div key={member.id} className="flex items-center space-x-3">
                        <img
                          src={member.avatar || '/default-avatar.png'}
                          alt={member.nickname || member.username}
                          className="w-10 h-10 rounded-full object-cover"
                        />
                        <div className="flex-1">
                          <p className="font-medium">{member.nickname || member.username}</p>
                          <p className="text-sm text-gray-500">@{member.username}</p>
                        </div>
                        {member.id === user?.id && (
                          <span className="text-xs bg-blue-100 text-blue-600 px-2 py-1 rounded">
                            我
                          </span>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        ) : (
          // 群组列表页面
          <div className="p-4">
            {groups.length === 0 ? (
              <div className="text-center py-12">
                <UserGroupIcon className="w-16 h-16 text-gray-300 mx-auto mb-4" />
                <p className="text-gray-500 mb-4">还没有加入任何群组</p>
                <button
                  onClick={() => setShowCreateModal(true)}
                  className="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600 transition-colors"
                >
                  创建群组
                </button>
              </div>
            ) : (
              <div className="space-y-3">
                {groups.map((group: Community) => (
                  <motion.div
                    key={group.id}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    className="bg-white rounded-lg p-4 shadow-sm"
                  >
                    <div className="flex items-center space-x-3">
                      <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-purple-500 rounded-lg flex items-center justify-center">
                        <UserGroupIcon className="w-6 h-6 text-white" />
                      </div>
                      <div className="flex-1">
                        <h3 className="font-semibold">{group.name}</h3>
                        <p className="text-sm text-gray-500">{group.description || group.desc}</p>
                        <p className="text-xs text-gray-400 mt-1">
                          {group.member_count || 0} 人
                        </p>
                      </div>
                      <div className="flex flex-col space-y-2">
                        <button
                          onClick={() => handleChatClick(group)}
                          className="text-blue-500 text-sm hover:bg-blue-50 px-3 py-1 rounded transition-colors"
                        >
                          聊天
                        </button>
                        <button
                          onClick={() => handleGroupClick(group)}
                          className="text-gray-500 text-sm hover:bg-gray-50 px-3 py-1 rounded transition-colors"
                        >
                          详情
                        </button>
                      </div>
                    </div>
                  </motion.div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>

      {/* 创建群组模态框 */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white rounded-lg p-6 w-full max-w-sm"
          >
            <h2 className="text-lg font-semibold mb-4">创建群组</h2>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  群组名称
                </label>
                <input
                  type="text"
                  value={newGroupName}
                  onChange={(e) => setNewGroupName(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="输入群组名称"
                />
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  群组描述
                </label>
                <textarea
                  value={newGroupDesc}
                  onChange={(e) => setNewGroupDesc(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="输入群组描述（可选）"
                  rows={3}
                />
              </div>
            </div>
            
            <div className="flex space-x-3 mt-6">
              <button
                onClick={() => setShowCreateModal(false)}
                className="flex-1 px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleCreateGroup}
                disabled={createGroupMutation.isPending}
                className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
              >
                {createGroupMutation.isPending ? '创建中...' : '创建'}
              </button>
            </div>
          </motion.div>
        </div>
      )}

      {/* 加入群组模态框 */}
      {showJoinModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <motion.div
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            className="bg-white rounded-lg p-6 w-full max-w-sm"
          >
            <h2 className="text-lg font-semibold mb-4">加入群组</h2>
            
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  群组ID
                </label>
                <input
                  type="text"
                  value={joinGroupId}
                  onChange={(e) => setJoinGroupId(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="输入群组ID"
                />
              </div>
            </div>
            
            <div className="flex space-x-3 mt-6">
              <button
                onClick={() => setShowJoinModal(false)}
                className="flex-1 px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleJoinGroup}
                disabled={joinGroupMutation.isPending}
                className="flex-1 px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50"
              >
                {joinGroupMutation.isPending ? '加入中...' : '加入'}
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </div>
  )
}

export default GroupPage