# 后端API文档
FatherID, _ := strconv.ParseInt(req.FatherID, 10, 64)//string---int64
## 登录相关接口

### 发送验证码

- 路径: api/login/send-code
- 方法: POST
- 限流: 每分钟4次
- 请求参数:

```json
{
  "email": "string"
}
```

- 成功响应:空响应

- postman测试示范

- 请求：

```json
{
    "email":"jhinih@163.com"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 用户注册

- 路径: api/login/register
- 方法: POST
- 限流: 每分钟4次

- 请求参数:

```json
{
  "email": "string",
  "code": "string",
  "password": "string",
  "username": "string",
  "avatar": "string"
}
```

- 成功响应:

```json
{
  "atoken": "string"
}
```

- postman测试示范

- 请求：

```json
{
    "email":"jhinih@163.com",
    "code": "786739",
    "password": "123456",
    "username": "wxy",
    "avatar": "https://xsg-bucket.oss-cn-shenzhen.aliyuncs.com/ca3a10f8e7fbe4291206386d8f17e6036168c18d454ca595f668cb6dc6038b84.png"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "atoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUzMjcyMzgsInJvbGUiOjAsInVzZXJpZCI6IjE5NTY1ODAxNzA2OTc1NDM2ODAiLCJ1c2VybmFtZSI6Ind4eSJ9.Ud-9jyAdKIdscuhXrNgJ4pqUKvCjE5L34jKrhTrWBUI"
    }
}
```

### 用户登录

- 路径: api/login/login
- 方法: POST
- 限流: 每分钟4次

- 请求参数:

```json
{
  "email": "string",
  "password": "string"
}
```

- 成功响应:

```json
{
  "atoken": "string",
  "rtoken": "string"
}
```

- postman测试示范

- 请求：

```json
{
    "email":"jhinih@163.com",
    "password": "123456"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "atoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUzMjc0MzIsInJvbGUiOjAsInVzZXJpZCI6IjE5NTY1ODAxNzA2OTc1NDM2ODAiLCJ1c2VybmFtZSI6Ind4eSJ9.YhGYtaU3J0YzDLTDXVO4M9t_3ZTBFhScxKZ4HLvNjv4",
        "rtoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU5MjUwMzIsInJvbGUiOjAsInVzZXJpZCI6IjE5NTY1ODAxNzA2OTc1NDM2ODAiLCJ1c2VybmFtZSI6Ind4eSJ9.8bMtXpqAvgnqQV5YeZ8HgeA9IVPc5ns3DWxyMt_T3rA"
    }
}
```

### 刷新Token

- 路径: api/login/refresh-token
- 方法: POST
- 限流: 每4秒8次
- 请求参数:

```json
{
  "rtoken": "string"
}
```

- 成功响应:

```json
{
  "atoken": "string"
}
```

postman测试示范

- 请求：

```json
{
    "rtoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU5MjUwMzIsInJvbGUiOjAsInVzZXJpZCI6IjE5NTY1ODAxNzA2OTc1NDM2ODAiLCJ1c2VybmFtZSI6Ind4eSJ9.8bMtXpqAvgnqQV5YeZ8HgeA9IVPc5ns3DWxyMt_T3rA"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "atoken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUzMjc0OTYsInJvbGUiOjAsInVzZXJpZCI6IjE5NTY1ODAxNzA2OTc1NDM2ODAiLCJ1c2VybmFtZSI6Ind4eSJ9.aRzauqn-XZ8Z8oFX0JJ5VsJwA7uTY_WrFNpHOlsp7IQ"
    }
}
```

## 用户相关接口

### 获取用户信息

- 路径: api/user/get-user-info
- 方法: GET
- 限流: 每秒20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string"
}
```

- 成功响应:

```json
{
  "id": "string",
  "username": "string",
  "avatar": "string",
  "role": "string",
  "bio": "string",
  "create_at": "string",
  "email": "string"
}
```

- postman测试示范

- 请求：

```json
{
    "id":"1954574560636440576"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "user": {
            "id": "1956682786274283520",
            "CreatedTime": "1755344504314",
            "UpdatedTime": "1755344504314",
            "username": "jhinih",
            "password": "$2a$10$cQBULo6pIuBd/TWz0DTrSew.B5fJQe6OI4wa8K1l4Jlv6Lcl7TIUa",
            "email": "3090497217@qq.com",
            "avatar": "",
            "role": "0",
            "Phone": "req.Phone",
            "ClientIp": "req.ClientIp",
            "ClientPort": "req.ClientPort",
            "LoginTime": "2025-08-16T19:41:44.314+08:00",
            "HeartbeatTime": "2025-08-16T19:41:44.314+08:00",
            "login_out_time": "2025-08-16T19:41:44.314+08:00",
            "IsLogout": false,
            "DeviceInfo": "req.DeviceInfo",
            "Bio": "jhinih很懒，什么都没留下"
        }
    }
}
```

### 获取我的信息

- 路径: api/user/get-my-info
- 方法: GET
- 限流: 每秒5次，10次突发
- 认证: 需要
- 请求参数: 无
- 成功响应:

```json
{
  "id": "string",
  "username": "string",
  "avatar": "string",
  "role": "string",
  "bio": "string",
  "create_at": "string",
  "email": "string"
}
```

- postman测试示范

- 请求：

```json

```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "user": {
            "id": "1956682786274283520",
            "CreatedTime": "1755344504314",
            "UpdatedTime": "1755344504314",
            "username": "jhinih",
            "password": "$2a$10$cQBULo6pIuBd/TWz0DTrSew.B5fJQe6OI4wa8K1l4Jlv6Lcl7TIUa",
            "email": "3090497217@qq.com",
            "avatar": "",
            "role": "0",
            "Phone": "req.Phone",
            "ClientIp": "req.ClientIp",
            "ClientPort": "req.ClientPort",
            "LoginTime": "2025-08-16T19:41:44.314+08:00",
            "HeartbeatTime": "2025-08-16T19:41:44.314+08:00",
            "login_out_time": "2025-08-16T19:41:44.314+08:00",
            "IsLogout": false,
            "DeviceInfo": "req.DeviceInfo",
            "Bio": "jhinih很懒，什么都没留下"
        }
    }
}
```

### 获取用户资料

- 路径: api/user/get-profile
- 方法: GET
- 限流: 每秒10次，20次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string"
}
```

- 成功响应:

```json
{
  "id": "string",
  "username": "string",
  "avatar": "string",
  "role": "string",
  "bio": "string",
  "email": "string",
  "create_at": "string"
}
```

- postman测试示范

- 请求：

```json
{
    "id":"1956580170697543680"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "user": {
            "id": "1956682786274283520",
            "CreatedTime": "1755344504314",
            "UpdatedTime": "1755344504314",
            "username": "jhinih",
            "password": "$2a$10$cQBULo6pIuBd/TWz0DTrSew.B5fJQe6OI4wa8K1l4Jlv6Lcl7TIUa",
            "email": "3090497217@qq.com",
            "avatar": "",
            "role": "0",
            "Phone": "req.Phone",
            "ClientIp": "req.ClientIp",
            "ClientPort": "req.ClientPort",
            "LoginTime": "2025-08-16T19:41:44.314+08:00",
            "HeartbeatTime": "2025-08-16T19:41:44.314+08:00",
            "login_out_time": "2025-08-16T19:41:44.314+08:00",
            "IsLogout": false,
            "DeviceInfo": "req.DeviceInfo",
            "Bio": "jhinih很懒，什么都没留下"
        }
    }
}
```

### 设置用户资料

- 路径: api/user/set-profile
- 方法: POST
- 限流: 每秒4次，8次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string",
  "username": "string",
  "avatar": "string",
  "bio": "string",
  "email": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
  "id": "1956580170697543680",
  "user": {
            "id": "1956682786274283520",
            "CreatedTime": "1755344504314",
            "UpdatedTime": "1755344504314",
            "username": "jhinih",
            "password": "$2a$10$cQBULo6pIuBd/TWz0DTrSew.B5fJQe6OI4wa8K1l4Jlv6Lcl7TIUa",
            "email": "3090497217@qq.com",
            "avatar": "",
            "role": "0",
            "Phone": "req.Phone",
            "ClientIp": "req.ClientIp",
            "ClientPort": "req.ClientPort",
            "LoginTime": "2025-08-16T19:41:44.314+08:00",
            "HeartbeatTime": "2025-08-16T19:41:44.314+08:00",
            "login_out_time": "2025-08-16T19:41:44.314+08:00",
            "IsLogout": false,
            "DeviceInfo": "req.DeviceInfo",
            "Bio": "jhinih很懒，什么都没留下"
        }
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 设置用户角色

- 路径: api/user/set-role
- 方法: POST
- 限流: 每秒4次，8次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string",
  "role": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
  "id": "1956682786274283520",
  "role": "1"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

## 文件上传接口

### 上传文件

- 路径: api/file/upload
- 方法: POST
- 限流: 每分钟3次，5次突发
- 认证: 需要
- 请求参数: multipart/form-data格式的文件上传
    - 参数名: file
    - 类型: 文件

- 成功响应:

```json
{
  "url": "string"
}
```

- postman测试示范

- 请求：

form-data格式上传文件
key: file
类型：File
value: 选择文件
格式：任意格式文件

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "url": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/59dbd7311ac3ce04d757421fef15964c01a2841149f388cac87f528d0d151ae9.mp4"
    }
}
```

## 视频相关接口

### 获取视频

- 路径: api/videos
- 方法: GET
- 认证: 需要
- 请求参数:

```json
{
  "page": "string",
  "page_size": "string",
  "order_by": "string"
}
```

- 成功响应:

```json
{
  "data": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "video_url": "string",
      "cover_url": "string",
      "like_count": "string",
      "comment_count": "string",
      "share_count": "string",
      "view_count": "string",
      "user_id": "string",
      "created_at": "string",
      "updated_at": "string"
    }
  ],
  "page": "string",
  "page_size": "string",
  "total": "string",
  "has_more": "boolean"
}
```

- postman测试示范

- 请求：

```json

```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "data": [
            {
                "id": "1956251539475533824",
                "title": "社团学姐你最喜欢哪一位",
                "description": "美丽的女士",
                "url": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/59dbd7311ac3ce04d757421fef15964c01a2841149f388cac87f528d0d151ae9.mp4",
                "cover": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/3c65c9a4a3d21cadf6df861e13e24743877e99af6532537c5fc8dc5ea6760e8e.jpg",
                "likes": "1",
                "comments": "4",
                "shares": "0",
                "user_id": "0",
                "created_at": "1755241687",
                "PublishTime": "2025-08-15T15:08:07.064+08:00",
                "type": "video",
                "is_private": false
            },
            {
                "id": "1956590790532665344",
                "title": "社团学姐你最喜欢哪一位",
                "description": "美丽的女士",
                "url": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/59dbd7311ac3ce04d757421fef15964c01a2841149f388cac87f528d0d151ae9.mp4",
                "cover": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/3c65c9a4a3d21cadf6df861e13e24743877e99af6532537c5fc8dc5ea6760e8e.jpg",
                "likes": "1",
                "comments": "0",
                "shares": "0",
                "user_id": "1956580170697543680",
                "created_at": "1755322570",
                "PublishTime": "2025-08-16T13:36:10.819+08:00",
                "type": "video",
                "is_private": false
            },
            {
                "id": "1956596075691249664",
                "title": "社团学姐你最喜欢哪一位",
                "description": "美丽的女士",
                "url": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/59dbd7311ac3ce04d757421fef15964c01a2841149f388cac87f528d0d151ae9.mp4",
                "cover": "https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/3c65c9a4a3d21cadf6df861e13e24743877e99af6532537c5fc8dc5ea6760e8e.jpg",
                "likes": "0",
                "comments": "0",
                "shares": "0",
                "user_id": "1956580170697543680",
                "created_at": "1755323830",
                "PublishTime": "2025-08-16T13:57:10.9+08:00",
                "type": "video",
                "is_private": false
            }
        ],
        "page": "0",
        "page_size": "0",
        "total": "3",
        "has_more": true
    }
}
```

### 创建视频

- 路径: api/videos/create-video
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "video_path":"string",
 "cover_path":"string",
 "title":  "string",
 "description": "string",
 "is_private":  "boolean",
 "type": "string",
  "user_id": "string"
}
```

- 成功响应:空对象

- postman测试示范

- 请求：

```json
{
 "video_path":"https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/59dbd7311ac3ce04d757421fef15964c01a2841149f388cac87f528d0d151ae9.mp4",
 "cover_path":"https://jhinih-tiktok.oss-cn-hongkong.aliyuncs.com/3c65c9a4a3d21cadf6df861e13e24743877e99af6532537c5fc8dc5ea6760e8e.jpg",
 "title":  "社团学姐你最喜欢哪一位",
 "description": "美丽的女士",
 "is_private":  false,
 "type": "video",
 "user_id": "1956580170697543680"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 点赞视频

- 路径: api/videos/like-video
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "video_id": "string",
  "owner_id": "string",
  "user_id": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
  "video_id": "1956590790532665344",
  "owner_id": "1956580170697543680",
  "user_id": "1956580170697543680"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 获取视频点赞数

- 路径: api/videos/get-video-likes
- 方法: GET
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "video_id": "string"
}
```

- 成功响应:

```json
{
  "video_likes": "string"
}
```

- postman测试示范

- 请求：

- <http://localhost:8080/api/videos/get-video-likes?video_id=1956590790532665344>

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "video_likes": "1"
    }
}
```

### 评论视频

- 路径: api/videos/comment-video
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "video_id": "string",
  "content": "any",
  "owner_id": "string",
  "user_id": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
  "video_id": "1956590792665344",
  "content": "hongkong.aliyuncs.com/3c65c9a4a3d21cadf6df861e13e24743877e99af6532537c5fc8dc5ea6760e8e.jpg",
  "user_id": "1956580170697543680"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 评论评论

- 路径: api/videos/comment-comments
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "content": "any",
  "owner_id": "string",
  "user_id": "string",
  "father_id": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
  "content": "hongkong.aliyuncs.com/3c65c9a4a3d21cadf6df861e13e24743877e99af6532537c5fc8dc5ea6760e8e.jpg",
  "user_id": "1956580170697543680",
  "father_id": "1956603758158811136"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 点赞评论

- 路径: api/videos/like-comments
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string",
  "owner_id": "string",
  "user_id": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
  "user_id": "1956580170697543680",
  "comment_id": "1956603758158811136"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 获取评论

- 路径: api/videos/get-comments
- 方法: GET
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string",
  "is_video": "boolean",
  "before_id":"string",
  "page":"string",
  "page_size":"string",
  "order_by":"string"
}
```

- 成功响应:

```json
{
  "comments": [
    {
      "id": "string",
      "content": "string",
      "user_id": "string",
      "video_id": "string",
      "like_count": "string",
      "created_at": "string",
      "updated_at": "string"
    }
  ],
  "length": "string"
}
```

- postman测试示范

- 请求：

- <http://localhost:8080/api/videos/get-comments?id=1956725200238153728&is_video=false>

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "comments": [
            {
                "id": "1956725555181129728",
                "CreatedTime": "1755354701195",
                "UpdatedTime": "1755354701195",
                "content": "aG9uZ2tvbmcuYWxpeXVuY3MuY29tLzNjNjVjOWE0YTNkMjFjYWRmNmRmODYxZTEzZTI0NzQzODc3ZTk5YWY2NTMyNTM3YzVmYzhkYzVlYTY3NjBlOGUuanBn",
                "video_id": "0",
                "father_id": "1956725200238153728",
                "user_id": "1956682786274283520",
                "owner_id": "1956682786274283520",
                "likes": "0",
                "comments": "0"
            },
            {
                "id": "1956725200238153728",
                "CreatedTime": "1755354616574",
                "UpdatedTime": "1755354616574",
                "content": "aG9uZ2tvbmcuYWxpeXVuY3MuY29tLzNjNjVjOWE0YTNkMjFjYWRmNmRmODYxZTEzZTI0NzQzODc3ZTk5YWY2NTMyNTM3YzVmYzhkYzVlYTY3NjBlOGUuanBn",
                "video_id": "1956722990410371072",
                "father_id": "1956725200238153728",
                "user_id": "1956682786274283520",
                "owner_id": "1956682786274283520",
                "likes": "1",
                "comments": "1"
            }
        ],
        "length": "2"
    }
}
```

### 获取评论详情

- 路径: api/videos/get-comment-all
- 方法: GET
- 限流: 每分钟20次，40次突发
- 认证: 需要

- 请求参数:

```json
{
  "id": "string"
}
```

- 成功响应:

```json
{
  "comment": {
    "id": "string",
    "content": "string",
    "user_id": "string",
    "video_id": "string",
    "like_count": "string",
    "created_at": "string",
    "updated_at": "string"
  }
}
```

- postman测试示范

- 请求：

- <http://localhost:8080/api/videos/get-comment-all?comment_id=1956725200238153728>

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "comment": {
            "id": "1956725200238153728",
            "CreatedTime": "1755354616574",
            "UpdatedTime": "1755354616574",
            "content": "aG9uZ2tvbmcuYWxpeXVuY3MuY29tLzNjNjVjOWE0YTNkMjFjYWRmNmRmODYxZTEzZTI0NzQzODc3ZTk5YWY2NTMyNTM3YzVmYzhkYzVlYTY3NjBlOGUuanBn",
            "video_id": "1956722990410371072",
            "father_id": "1956725200238153728",
            "user_id": "1956682786274283520",
            "owner_id": "1956682786274283520",
            "likes": "1",
            "comments": "1"
        }
    }
}
```

### 获取评论数

- 路径: api/videos/get-comment-member
- 方法: GET
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "id": "string",
  "is_video": "boolean"
}
```

- 成功响应:

```json
{
  "member":"string"
}
```

- postman测试示范

- 请求：

- <http://localhost:8080/api/videos/get-comment-member?id=1956725200238153728&is_video=false>

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "member": "1"
    }
}
```

## 关系相关接口

### 添加好友

- 路径: api/contact/add-friend
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "username": "string"
}
```

- 成功响应: 空对象

- postman测试示范

- 请求：

```json
{
    "user_name":"jhinih"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 获取好友列表

- 路径: api/contact/get-friend-list
- 方法: POST
- 限流: 每分钟20次，40次突发
- 请求参数:

```json
{
  "user_id": "string",
  "page": "string",
  "page_size": "string",
  "order_by": "string"
}
```

- 成功响应:

```json
{
    "users": [
    {
      "id": "string",
      "CreatedTime": "string",
      "UpdatedTime": "string",
      "username": "string",
      "password": "string",
      "email": "string",
      "avatar": "string",
      "role": "string",
      "Phone": "string",
      "ClientIp": "string",
      "ClientPort": "string",
      "LoginTime": "string",
      "HeartbeatTime": "string",
      "login_out_time": "string",
      "IsLogout": "boolean",
      "DeviceInfo": "string",
      "Bio": "string"
    }
  ]
}
```

- postman测试示范

- 请求：

```json
{

}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "users": [
            {
                "id": "1956744637980872704",
                "CreatedTime": "1755359250905",
                "UpdatedTime": "1755359250905",
                "username": "王鑫宇",
                "password": "$2a$10$EV6OWJVREQEmj3gbkgdtnOe2sVQhkTqvpqsKQ1taORsHEUdcMb5aq",
                "email": "jhinih123@gmail.com",
                "avatar": "",
                "role": "0",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T23:47:30.906+08:00",
                "HeartbeatTime": "2025-08-16T23:47:30.906+08:00",
                "login_out_time": "2025-08-16T23:47:30.906+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "王鑫宇很懒，什么都没留下"
            },
            {
                "id": "1956733730454245376",
                "CreatedTime": "1755356650350",
                "UpdatedTime": "1755356650350",
                "username": "wxy",
                "password": "$2a$10$VBVruJABtOej3YXbBvp01.1UAiqfbCdr0qUjgx/30J3imlCcqLKce",
                "email": "jhinih@163.com",
                "avatar": "",
                "role": "0",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T23:04:10.349+08:00",
                "HeartbeatTime": "2025-08-16T23:04:10.349+08:00",
                "login_out_time": "2025-08-16T23:04:10.349+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "wxy很懒，什么都没留下"
            }
        ]
    }
}
```

### 获取在线用户列表

- 路径: api/contact/get-user-list-online
- 方法: GET
- 限流: 每分钟20次，40次突发
- 请求参数: 无
- 成功响应:

```json
{
  "users": [
    {
      "id": "string",
      "username": "string",
      "avatar": "string",
      "role": "string",
      "bio": "string",
      "email": "string",
      "create_at": "string"
    }
  ]
}
```

- postman测试示范

- 请求：

```json
{

}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "users": [
            {
                "id": "1956682786274283520",
                "CreatedTime": "1755344504314",
                "UpdatedTime": "1755344504314",
                "username": "jhinih",
                "password": "$2a$10$cQBULo6pIuBd/TWz0DTrSew.B5fJQe6OI4wa8K1l4Jlv6Lcl7TIUa",
                "email": "3090497217@qq.com",
                "avatar": "",
                "role": "1",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T19:41:44.314+08:00",
                "HeartbeatTime": "2025-08-16T19:41:44.314+08:00",
                "login_out_time": "2025-08-16T19:41:44.314+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "jhinih很懒，什么都没留下"
            },
            {
                "id": "1956733730454245376",
                "CreatedTime": "1755356650350",
                "UpdatedTime": "1755356650350",
                "username": "wxy",
                "password": "$2a$10$VBVruJABtOej3YXbBvp01.1UAiqfbCdr0qUjgx/30J3imlCcqLKce",
                "email": "jhinih@163.com",
                "avatar": "",
                "role": "0",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T23:04:10.349+08:00",
                "HeartbeatTime": "2025-08-16T23:04:10.349+08:00",
                "login_out_time": "2025-08-16T23:04:10.349+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "wxy很懒，什么都没留下"
            },
            {
                "id": "1956744637980872704",
                "CreatedTime": "1755359250905",
                "UpdatedTime": "1755359250905",
                "username": "王鑫宇",
                "password": "$2a$10$EV6OWJVREQEmj3gbkgdtnOe2sVQhkTqvpqsKQ1taORsHEUdcMb5aq",
                "email": "jhinih123@gmail.com",
                "avatar": "",
                "role": "0",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T23:47:30.906+08:00",
                "HeartbeatTime": "2025-08-16T23:47:30.906+08:00",
                "login_out_time": "2025-08-16T23:47:30.906+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "王鑫宇很懒，什么都没留下"
            }
        ]
    }
}
```

### 获取群组成员

- 路径: api/contact/get-group-users
- 方法: GET
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "community_id": "string"
}
```

- 成功响应:

```json
{
  "users": [
    {
      "id": "string",
      "username": "string",
      "avatar": "string",
      "role": "string",
      "bio": "string",
      "email": "string",
      "create_at": "string"
    }
  ]
}
```

- postman测试示范

- 请求：

```json
{

}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "users": [
            {
                "id": "1956682786274283520",
                "CreatedTime": "1755344504314",
                "UpdatedTime": "1755344504314",
                "username": "jhinih",
                "password": "$2a$10$cQBULo6pIuBd/TWz0DTrSew.B5fJQe6OI4wa8K1l4Jlv6Lcl7TIUa",
                "email": "3090497217@qq.com",
                "avatar": "",
                "role": "1",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T19:41:44.314+08:00",
                "HeartbeatTime": "2025-08-16T19:41:44.314+08:00",
                "login_out_time": "2025-08-16T19:41:44.314+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "jhinih很懒，什么都没留下"
            },
            {
                "id": "1956733730454245376",
                "CreatedTime": "1755356650350",
                "UpdatedTime": "1755356650350",
                "username": "wxy",
                "password": "$2a$10$VBVruJABtOej3YXbBvp01.1UAiqfbCdr0qUjgx/30J3imlCcqLKce",
                "email": "jhinih@163.com",
                "avatar": "",
                "role": "0",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T23:04:10.349+08:00",
                "HeartbeatTime": "2025-08-16T23:04:10.349+08:00",
                "login_out_time": "2025-08-16T23:04:10.349+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "wxy很懒，什么都没留下"
            },
            {
                "id": "1956744637980872704",
                "CreatedTime": "1755359250905",
                "UpdatedTime": "1755359250905",
                "username": "王鑫宇",
                "password": "$2a$10$EV6OWJVREQEmj3gbkgdtnOe2sVQhkTqvpqsKQ1taORsHEUdcMb5aq",
                "email": "jhinih123@gmail.com",
                "avatar": "",
                "role": "0",
                "Phone": "req.Phone",
                "ClientIp": "req.ClientIp",
                "ClientPort": "req.ClientPort",
                "LoginTime": "2025-08-16T23:47:30.906+08:00",
                "HeartbeatTime": "2025-08-16T23:47:30.906+08:00",
                "login_out_time": "2025-08-16T23:47:30.906+08:00",
                "IsLogout": false,
                "DeviceInfo": "req.DeviceInfo",
                "Bio": "王鑫宇很懒，什么都没留下"
            }
        ]
    }
}
```

### 获取群组列表

- 路径: api/contact/get-group-list
- 方法: GET
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
  "user_id": "string"
}
```

- 成功响应:

```json
{
  "groups": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "creator_id": "string",
      "created_at": "string",
      "updated_at": "string"
    }
  ]
}
```

- postman测试示范

- 请求：

```json
{
 "owner_name": "jhinih",
 "name": "一",
 "icon": "",
 "desc": ""
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 创建群聊

- 路径: api/contact/create-community
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
 "owner_id": "string",
 "owner_name": "string",
 "name": "string",
 "icon": "string",
 "desc": "string"
}
```

- 成功响应:空响应

- postman测试示范

- 请求：

```json
{
 "owner_name": "jhinih",
 "name": "一",
 "icon": "",
 "desc": ""
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 加入群聊

- 路径: api/contact/join-community
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
 "owner_id": "string",
 "community_id": "string"
}
```

- 成功响应:空响应

- postman测试示范

- 请求：

```json
{
 "community_id": "1"
}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {}
}
```

### 加载群聊

- 路径: api/contact/load-community
- 方法: POST
- 限流: 每分钟20次，40次突发
- 认证: 需要
- 请求参数:

```json
{
 "owner_id": "string"
}
```

- 成功响应:

```json
{
    "groups": [
    {
      "id": "string",
      "name": "string",
      "description": "string",
      "creator_id": "string",
      "created_at": "string",
      "updated_at": "string"
    }
  ]
}
```

- postman测试示范

- 请求：

```json
{

}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "groups": [
            {
                "id": "1",
                "CreatedTime": "1755364195259",
                "UpdatedTime": "1755364195259",
                "name": "一",
                "owner_id": "1956682786274283520",
                "owner_name": "jhinih",
                "img": "",
                "desc": ""
            },
            {
                "id": "2",
                "CreatedTime": "1755364336138",
                "UpdatedTime": "1755364336138",
                "name": "二",
                "owner_id": "1956682786274283520",
                "owner_name": "jhinih",
                "img": "",
                "desc": ""
            },
            {
                "id": "1956774945174327296",
                "CreatedTime": "1755366476712",
                "UpdatedTime": "1755366476712",
                "name": "三",
                "owner_id": "1956682786274283520",
                "owner_name": "jhinih",
                "img": "",
                "desc": ""
            },
            {
                "id": "1956783355366215680",
                "CreatedTime": "1755368481855",
                "UpdatedTime": "1755368481855",
                "name": "四",
                "owner_id": "1956682786274283520",
                "owner_name": "jhinih",
                "img": "",
                "desc": ""
            }
        ]
    }
}
```

## 聊天相关接口

### 生成一次性 30 秒 WebSocket ticket

- 路径: api/chat/ws-ticket
- 方法: GET
- 协议: WebSocket
- 认证: 需要
- 请求参数:

```json
{
   "user_id":"string",
   "user_name":"string", 
   "role": "string"
}
```

- 成功响应:

```json
{
  "ticket": "string"
}
```

- postman测试示范

- 请求：

```json
{

}
```

- 响应：

```json
{
    "code": 20000,
    "message": "成功",
    "data": {
        "ticket": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTUzNzA3NjIsInJvbGUiOjAsInVzZXJpZCI6IjE5NTY2ODI3ODYyNzQyODM1MjAiLCJ1c2VybmFtZSI6IiJ9.esK82Jw7DBrdm-oyuazmIyXH1XHWzYryLK5NB4Dys5o"
    }
}
```

### WebSocket连接

- 路径: api/chat/ws
- 方法: GET
- 协议: WebSocket
- 认证: 需要
- 请求参数: 无
- 成功响应: WebSocket连接建立
