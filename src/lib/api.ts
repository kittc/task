import axios, { AxiosInstance, AxiosResponse } from 'axios'
import { 
  ApiResponse, 
  AuthResponse, 
  LoginRequest, 
  RegisterRequest,
  User,
  Task,
  CreateTaskRequest,
  UpdateTaskRequest,
  TaskQueryOptions,
  PaginatedResponse,
  Organization,
  Department,
  CreateOrganizationRequest,
  CreateDepartmentRequest,
  CreateUserRequest
} from '@/types'

// 创建axios实例
const api: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加认证token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 清除本地存储的认证信息
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      // 重定向到登录页面
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// 认证相关API
export const authAPI = {
  // 登录
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post<ApiResponse<AuthResponse>>('/auth/login', data)
    return response.data.data!
  },

  // 注册
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<ApiResponse<AuthResponse>>('/auth/register', data)
    return response.data.data!
  },

  // 获取用户档案
  getProfile: async (): Promise<User> => {
    const response = await api.get<ApiResponse<User>>('/profile')
    return response.data.data!
  },

  // 修改密码
  changePassword: async (oldPassword: string, newPassword: string): Promise<void> => {
    await api.post('/change-password', {
      old_password: oldPassword,
      new_password: newPassword
    })
  },

  // 登出
  logout: async (): Promise<void> => {
    await api.post('/logout')
  }
}

// 任务相关API
export const taskAPI = {
  // 获取任务列表
  getTasks: async (options: TaskQueryOptions = {}): Promise<PaginatedResponse<Task>> => {
    const params = new URLSearchParams()
    
    Object.entries(options).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString())
      }
    })

    const response = await api.get<ApiResponse<PaginatedResponse<Task>>>(`/tasks?${params}`)
    return response.data.data!
  },

  // 创建任务
  createTask: async (data: CreateTaskRequest): Promise<Task> => {
    const response = await api.post<ApiResponse<Task>>('/tasks', data)
    return response.data.data!
  },

  // 获取任务详情
  getTask: async (id: number): Promise<Task> => {
    const response = await api.get<ApiResponse<Task>>(`/tasks/${id}`)
    return response.data.data!
  },

  // 更新任务
  updateTask: async (id: number, data: UpdateTaskRequest): Promise<Task> => {
    const response = await api.put<ApiResponse<Task>>(`/tasks/${id}`, data)
    return response.data.data!
  },

  // 删除任务
  deleteTask: async (id: number): Promise<void> => {
    await api.delete(`/tasks/${id}`)
  },

  // 添加任务成员
  addTaskMember: async (taskId: number, userId: number): Promise<void> => {
    await api.post(`/tasks/${taskId}/members`, { user_id: userId })
  },

  // 移除任务成员
  removeTaskMember: async (taskId: number, memberId: number): Promise<void> => {
    await api.delete(`/tasks/${taskId}/members/${memberId}`)
  },

  // 创建清单项
  createChecklist: async (taskId: number, title: string, position = 0): Promise<any> => {
    const response = await api.post(`/tasks/${taskId}/checklists`, { title, position })
    return response.data.data
  },

  // 更新清单项
  updateChecklist: async (checklistId: number, data: any): Promise<any> => {
    const response = await api.put(`/tasks/checklists/${checklistId}`, data)
    return response.data.data
  },

  // 删除清单项
  deleteChecklist: async (checklistId: number): Promise<void> => {
    await api.delete(`/tasks/checklists/${checklistId}`)
  }
}

// 组织架构相关API
export const organizationAPI = {
  // 获取组织列表
  getOrganizations: async (): Promise<Organization[]> => {
    const response = await api.get<ApiResponse<Organization[]>>('/organizations')
    return response.data.data!
  },

  // 创建组织
  createOrganization: async (data: CreateOrganizationRequest): Promise<Organization> => {
    const response = await api.post<ApiResponse<Organization>>('/organizations', data)
    return response.data.data!
  },

  // 获取组织详情
  getOrganization: async (id: number): Promise<Organization> => {
    const response = await api.get<ApiResponse<Organization>>(`/organizations/${id}`)
    return response.data.data!
  },

  // 更新组织
  updateOrganization: async (id: number, data: Partial<CreateOrganizationRequest>): Promise<Organization> => {
    const response = await api.put<ApiResponse<Organization>>(`/organizations/${id}`, data)
    return response.data.data!
  },

  // 删除组织
  deleteOrganization: async (id: number): Promise<void> => {
    await api.delete(`/organizations/${id}`)
  },

  // 获取部门列表
  getDepartments: async (organizationId?: number): Promise<Department[]> => {
    const params = organizationId ? `?organization_id=${organizationId}` : ''
    const response = await api.get<ApiResponse<Department[]>>(`/departments${params}`)
    return response.data.data!
  },

  // 创建部门
  createDepartment: async (data: CreateDepartmentRequest): Promise<Department> => {
    const response = await api.post<ApiResponse<Department>>('/departments', data)
    return response.data.data!
  },

  // 获取部门详情
  getDepartment: async (id: number): Promise<Department> => {
    const response = await api.get<ApiResponse<Department>>(`/departments/${id}`)
    return response.data.data!
  },

  // 更新部门
  updateDepartment: async (id: number, data: Partial<CreateDepartmentRequest>): Promise<Department> => {
    const response = await api.put<ApiResponse<Department>>(`/departments/${id}`, data)
    return response.data.data!
  },

  // 删除部门
  deleteDepartment: async (id: number): Promise<void> => {
    await api.delete(`/departments/${id}`)
  },

  // 获取用户列表
  getUsers: async (params: { 
    department_id?: number, 
    role?: string, 
    is_active?: boolean 
  } = {}): Promise<User[]> => {
    const queryParams = new URLSearchParams()
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        queryParams.append(key, value.toString())
      }
    })
    
    const response = await api.get<ApiResponse<User[]>>(`/users?${queryParams}`)
    return response.data.data!
  },

  // 创建用户
  createUser: async (data: CreateUserRequest): Promise<User> => {
    const response = await api.post<ApiResponse<User>>('/users', data)
    return response.data.data!
  },

  // 获取用户详情
  getUser: async (id: number): Promise<User> => {
    const response = await api.get<ApiResponse<User>>(`/users/${id}`)
    return response.data.data!
  },

  // 更新用户
  updateUser: async (id: number, data: Partial<CreateUserRequest>): Promise<User> => {
    const response = await api.put<ApiResponse<User>>(`/users/${id}`, data)
    return response.data.data!
  },

  // 删除用户
  deleteUser: async (id: number): Promise<void> => {
    await api.delete(`/users/${id}`)
  },

  // 转移用户
  transferUser: async (userId: number, departmentId: number): Promise<void> => {
    await api.post(`/users/${userId}/transfer`, { department_id: departmentId })
  },

  // 获取组织架构
  getOrganizationStructure: async (): Promise<Organization[]> => {
    const response = await api.get<ApiResponse<Organization[]>>('/organization-structure')
    return response.data.data!
  },

  // 获取部门统计
  getDepartmentStats: async (id: number): Promise<any> => {
    const response = await api.get<ApiResponse<any>>(`/departments/${id}/stats`)
    return response.data.data!
  }
}

export default api