# 投票系统 API

这是一个基于 Gin 框架的投票系统 API，支持创建投票、编辑选项、进行投票以及查看投票结果。

## 功能特点

- 支持三种投票类型：二分选项（是/否）、单选、多选
- 用户可以创建、编辑和删除投票
- 用户可以添加、编辑和删除投票选项
- 用户可以进行投票，并根据投票类型进行相应的限制
- 支持查看投票结果和统计数据
- 支持设置投票截止时间
- 支持激活/停用投票
- **高级统计分析功能**：
  - 详细的投票统计信息，包括选项百分比、参与人数等
  - 投票时间分布分析
  - 热门投票排行榜
  - 用户投票行为分析
- **评论与讨论系统**：
  - 用户可以对投票进行评论
  - 支持评论回复功能，构建评论树
  - 评论作者可以编辑和删除自己的评论

## 技术栈

- Go 1.20
- Gin Web 框架
- GORM ORM 库
- SQLite 数据库

## 安装和运行

1. 确保已安装 Go 1.20 或更高版本
2. 克隆项目到本地
3. 安装依赖：
   ```
   go mod download
   ```
4. 运行项目：
   ```
   go run main.go
   ```
5. 服务器将在 http://localhost:8080 启动

## API 接口说明

### 用户相关接口

- `POST /api/users` - 创建用户
- `GET /api/users` - 获取用户列表
- `GET /api/users/:id` - 获取用户详情
- `GET /api/users/username/:username` - 通过用户名获取用户详情
- `GET /api/users/:id/stats` - 获取用户的投票统计信息

### 投票相关接口

- `POST /api/polls` - 创建投票
- `GET /api/polls` - 获取投票列表
- `GET /api/polls/:id` - 获取投票详情
- `PUT /api/polls/:id` - 更新投票信息
- `DELETE /api/polls/:id` - 删除投票
- `GET /api/polls/:id/results` - 获取投票结果
- `GET /api/polls/:id/stats` - 获取投票的详细统计信息

### 选项相关接口

- `POST /api/polls/:id/options` - 添加选项
- `PUT /api/polls/:id/options/:option_id` - 更新选项
- `DELETE /api/polls/:id/options/:option_id` - 删除选项

### 投票操作接口

- `POST /api/polls/:id/vote` - 进行投票
- `GET /api/polls/:id/user-votes` - 获取用户在特定投票中的投票记录

### 评论相关接口

- `POST /api/polls/:id/comments` - 添加评论
- `GET /api/polls/:id/comments` - 获取投票的评论
- `PUT /api/polls/:id/comments/:comment_id` - 更新评论
- `DELETE /api/polls/:id/comments/:comment_id` - 删除评论

### 统计和分析接口

- `GET /api/stats/trending` - 获取热门投票排行榜

## 示例请求

### 创建投票

```json
POST /api/polls
{
  "title": "最喜欢的编程语言",
  "description": "请选择你最喜欢的编程语言",
  "type": "single",
  "options": ["Go", "Python", "JavaScript", "Java"],
  "end_time": "2023-12-31T23:59:59Z"
}
```

### 进行投票

```json
POST /api/polls/:id/vote
{
  "option_ids": ["option_id_1"]
}
```

### 获取投票详细统计

```
GET /api/polls/:id/stats
```

响应示例：
```json
{
  "poll": {
    "id": "poll_id",
    "title": "最喜欢的编程语言",
    "description": "请选择你最喜欢的编程语言",
    "type": "single",
    "options": [...]
  },
  "total_votes": 150,
  "unique_voters": 120,
  "option_stats": [
    {
      "id": "option_id_1",
      "text": "Go",
      "count": 50,
      "percentage": 33.33
    },
    {
      "id": "option_id_2",
      "text": "Python",
      "count": 40,
      "percentage": 26.67
    },
    ...
  ],
  "time_distribution": [
    {"hour": 0, "count": 5},
    {"hour": 1, "count": 3},
    ...
    {"hour": 23, "count": 8}
  ]
}
```

### 添加评论

```json
POST /api/polls/:id/comments
{
  "content": "这是一个很有意思的投票！",
  "parent_id": "" // 可选，回复其他评论时提供
}
```

### 获取投票评论

```
GET /api/polls/:id/comments
```

响应示例：
```json
{
  "comments": [
    {
      "id": "comment_id_1",
      "poll_id": "poll_id",
      "user_id": "user_id_1",
      "content": "这是一个很有意思的投票！",
      "parent_id": "",
      "created_at": "2023-05-20T10:30:00Z",
      "updated_at": "2023-05-20T10:30:00Z",
      "user": {
        "id": "user_id_1",
        "username": "user1"
      },
      "replies": [
        {
          "id": "comment_id_2",
          "poll_id": "poll_id",
          "user_id": "user_id_2",
          "content": "我也这么认为！",
          "parent_id": "comment_id_1",
          "created_at": "2023-05-20T11:00:00Z",
          "updated_at": "2023-05-20T11:00:00Z",
          "user": {
            "id": "user_id_2",
            "username": "user2"
          }
        }
      ]
    }
  ]
}
```

## 注意事项

- 在实际生产环境中，应该添加适当的身份验证和授权机制
- 当前实现使用 SQLite 作为数据库，可以根据需要替换为其他数据库
- 为了简化演示，投票时如果没有提供用户 ID，系统会自动创建一个临时用户 