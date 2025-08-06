// 用户角色类型
export type Role = 'employee' | 'supervisor' | 'manager' | 'admin'

// 任务状态类型
export type TaskStatus = 'todo' | 'in_progress' | 'overdue' | 'paused' | 'cancelled' | 'completed'

// 优先级类型
export type Priority = 'low' | 'medium' | 'high' | 'urgent'

// 用户类型
export interface User {
  id: number
  username: string
  email: string
  full_name: string
  avatar?: string
  role: Role
  is_active: boolean
  department_id?: number
  department?: Department
  last_login_at?: string
  created_at: string
  updated_at: string
}

// 组织类型
export interface Organization {
  id: number
  name: string
  description?: string
  created_at: string
  updated_at: string
  departments?: Department[]
}

// 部门类型
export interface Department {
  id: number
  organization_id: number
  name: string
  description?: string
  created_at: string
  updated_at: string
  organization?: Organization
  users?: User[]
}

// 任务类型
export interface Task {
  id: number
  title: string
  description?: string
  status: TaskStatus
  priority: Priority
  creator_id: number
  assignee_id?: number
  due_date?: string
  completed_at?: string
  estimated_hours?: number
  actual_hours?: number
  created_at: string
  updated_at: string
  creator: User
  assignee?: User
  members?: TaskMember[]
  checklists?: Checklist[]
  comments?: Comment[]
  attachments?: Attachment[]
}

// 任务成员类型
export interface TaskMember {
  id: number
  task_id: number
  user_id: number
  joined_at: string
  user: User
}

// 清单项类型
export interface Checklist {
  id: number
  task_id: number
  title: string
  is_completed: boolean
  position: number
  created_at: string
  updated_at: string
}

// 评论类型
export interface Comment {
  id: number
  task_id: number
  user_id: number
  content: string
  created_at: string
  updated_at: string
  user: User
}

// 附件类型
export interface Attachment {
  id: number
  task_id: number
  user_id: number
  file_name: string
  file_path: string
  file_size: number
  mime_type: string
  created_at: string
  user: User
}

// 通知类型
export interface Notification {
  id: number
  user_id: number
  task_id?: number
  type: string
  title: string
  content: string
  is_read: boolean
  created_at: string
  updated_at: string
  task?: Task
}

// API响应类型
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

// 分页响应类型
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 登录请求类型
export interface LoginRequest {
  username: string
  password: string
}

// 注册请求类型
export interface RegisterRequest {
  username: string
  email: string
  password: string
  full_name: string
  department_id?: number
  role?: Role
}

// 认证响应类型
export interface AuthResponse {
  token: string
  user: User
  expires_at: string
}

// 创建任务请求类型
export interface CreateTaskRequest {
  title: string
  description?: string
  priority?: Priority
  assignee_id?: number
  due_date?: string
  estimated_hours?: number
  member_ids?: number[]
  checklists?: CreateChecklistRequest[]
}

// 更新任务请求类型
export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: TaskStatus
  priority?: Priority
  assignee_id?: number
  due_date?: string
  estimated_hours?: number
  actual_hours?: number
}

// 创建清单项请求类型
export interface CreateChecklistRequest {
  title: string
  position?: number
}

// 任务查询选项类型
export interface TaskQueryOptions {
  status?: TaskStatus
  priority?: Priority
  creator_id?: number
  assignee_id?: number
  member_id?: number
  due_from?: string
  due_to?: string
  search?: string
  page?: number
  page_size?: number
  sort_by?: string
  sort_order?: string
}

// 创建组织请求类型
export interface CreateOrganizationRequest {
  name: string
  description?: string
}

// 创建部门请求类型
export interface CreateDepartmentRequest {
  organization_id: number
  name: string
  description?: string
}

// 创建用户请求类型
export interface CreateUserRequest {
  department_id?: number
  username: string
  email: string
  password: string
  full_name: string
  role: Role
}

// 状态标签映射
export const statusLabels: Record<TaskStatus, string> = {
  todo: '待办',
  in_progress: '进行中',
  overdue: '延期',
  paused: '暂停',
  cancelled: '取消',
  completed: '完成'
}

// 优先级标签映射
export const priorityLabels: Record<Priority, string> = {
  low: '低',
  medium: '中',
  high: '高',
  urgent: '紧急'
}

// 角色标签映射
export const roleLabels: Record<Role, string> = {
  employee: '职员',
  supervisor: '主管',
  manager: '经理',
  admin: '管理员'
}