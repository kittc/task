# 任务管理平台

一个基于Go+React+SQLite的现代化任务管理平台，类似于Trello，提供完整的组织架构管理和任务协作功能。

## 功能特性

### 🏢 组织架构管理
- **多级组织结构**：公司 → 部门 → 用户
- **角色权限管理**：管理员、经理、主管、职员四级权限
- **用户管理**：用户创建、编辑、转移、禁用等操作
- **部门统计**：实时统计部门人员和任务情况

### 📋 任务管理
- **任务CRUD**：创建、查看、编辑、删除任务
- **任务状态**：待办 → 进行中 → 异常（延期/暂停/取消）→ 完成
- **优先级设置**：低、中、高、紧急四个级别
- **任务分配**：支持指定负责人和多个成员
- **子任务清单**：可添加多个检查项
- **截止时间**：设置任务截止时间，自动超时提醒

### 👥 协作功能
- **任务成员**：添加/移除任务参与者
- **评论系统**：任务讨论和沟通
- **文件附件**：支持文件上传和管理
- **实时通知**：任务分配、状态变更、到期提醒

### 🔐 权限控制
- **基于角色的访问控制**：不同角色拥有不同权限
- **任务访问控制**：只能访问相关任务
- **数据隔离**：部门间数据隔离

## 技术架构

### 后端技术栈
- **Go 1.21+**：高性能后端服务
- **Gin**：Web框架
- **GORM**：ORM框架
- **SQLite**：轻量级数据库
- **JWT**：认证和授权
- **bcrypt**：密码加密

### 前端技术栈
- **React 18**：现代化前端框架
- **TypeScript**：类型安全
- **Vite**：快速构建工具
- **Tailwind CSS**：原子化CSS框架
- **shadcn/ui**：现代化UI组件库
- **React Query**：数据获取和缓存
- **React Router**：路由管理
- **React Hook Form**：表单处理

### 项目结构

```
task-management-platform/
├── cmd/
│   └── server/          # 服务器入口
├── internal/
│   ├── models/          # 数据模型
│   ├── services/        # 业务逻辑
│   ├── handlers/        # HTTP处理器
│   ├── middleware/      # 中间件
│   └── database/        # 数据库配置
├── src/                 # 前端源码
│   ├── components/      # React组件
│   ├── pages/          # 页面组件
│   ├── hooks/          # 自定义Hooks
│   ├── contexts/       # React Context
│   ├── lib/            # 工具库
│   └── types/          # TypeScript类型
├── go.mod              # Go依赖管理
├── package.json        # Node.js依赖管理
├── Makefile           # 构建脚本
└── README.md          # 项目文档
```

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 安装依赖

```bash
# 安装所有依赖
make install

# 或者分别安装
go mod tidy    # Go依赖
npm install    # Node.js依赖
```

### 开发环境

```bash
# 启动开发环境（并行启动后端和前端）
make dev

# 或者分别启动
make start-backend   # 后端服务 (http://localhost:8080)
make start-frontend  # 前端服务 (http://localhost:5173)
```

### 构建项目

```bash
# 构建项目
make build

# 生产环境构建
make prod-build
```

### 默认账户

系统初始化时会创建以下默认账户：

| 用户名 | 密码 | 角色 | 说明 |
|--------|------|------|------|
| admin | admin123 | 管理员 | 系统管理员 |
| manager1 | admin123 | 经理 | 部门经理 |
| supervisor1 | admin123 | 主管 | 部门主管 |
| employee1 | admin123 | 职员 | 普通职员 |

## API文档

### 认证API
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/refresh` - 刷新Token
- `GET /api/v1/profile` - 获取用户信息
- `POST /api/v1/change-password` - 修改密码

### 任务API
- `GET /api/v1/tasks` - 获取任务列表
- `POST /api/v1/tasks` - 创建任务
- `GET /api/v1/tasks/:id` - 获取任务详情
- `PUT /api/v1/tasks/:id` - 更新任务
- `DELETE /api/v1/tasks/:id` - 删除任务
- `POST /api/v1/tasks/:id/members` - 添加任务成员
- `DELETE /api/v1/tasks/:id/members/:member_id` - 移除任务成员

### 组织架构API
- `GET /api/v1/organizations` - 获取组织列表
- `POST /api/v1/organizations` - 创建组织
- `GET /api/v1/departments` - 获取部门列表
- `POST /api/v1/departments` - 创建部门
- `GET /api/v1/users` - 获取用户列表（需要管理权限）
- `POST /api/v1/users` - 创建用户（需要管理权限）

## 配置选项

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `PORT` | 8080 | 服务器端口 |
| `DB_PATH` | ./task_management.db | 数据库文件路径 |
| `JWT_SECRET` | your-secret-key-here | JWT签名密钥 |
| `GIN_MODE` | debug | Gin运行模式 |

### 数据库配置

系统使用SQLite作为数据库，支持自动迁移。首次启动时会自动创建数据库表并插入初始数据。

## 开发指南

### 添加新功能

1. **后端**：
   - 在 `internal/models/` 中定义数据模型
   - 在 `internal/services/` 中实现业务逻辑
   - 在 `internal/handlers/` 中创建HTTP处理器
   - 在 `cmd/server/main.go` 中注册路由

2. **前端**：
   - 在 `src/types/` 中定义TypeScript类型
   - 在 `src/lib/api.ts` 中添加API调用
   - 在 `src/components/` 或 `src/pages/` 中创建组件
   - 在 `src/App.tsx` 中添加路由

### 代码规范

- **Go**：遵循官方Go代码规范
- **TypeScript**：使用严格模式，所有变量需要类型声明
- **CSS**：使用Tailwind CSS原子类，避免自定义CSS

### 测试

```bash
# 运行所有测试
make test

# 仅运行后端测试
go test ./...

# 仅运行前端测试
npm test
```

## 部署

### Docker部署

```dockerfile
# Dockerfile示例
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o server cmd/server/main.go

FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/
COPY --from=backend-builder /app/server .
COPY --from=frontend-builder /app/dist ./static
CMD ["./server"]
```

### 生产环境配置

```bash
# 设置环境变量
export GIN_MODE=release
export JWT_SECRET=your-production-secret-key
export DB_PATH=/data/task_management.db

# 启动服务
make prod-start
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 实现基础的任务管理功能
- 完整的组织架构管理
- 用户认证和权限控制
- 现代化的Web界面

## 联系方式

如有问题或建议，请通过以下方式联系：

- GitHub Issues: [项目Issues页面]
- Email: [您的邮箱]

---

**注意**：这是一个演示项目，请根据实际需求进行安全加固和性能优化后再用于生产环境。