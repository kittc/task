import { Routes, Route, Navigate } from 'react-router-dom'
import { useEffect, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { authAPI } from '@/lib/api'
import { User } from '@/types'
import Layout from '@/components/Layout'
import Login from '@/pages/Login'
import Dashboard from '@/pages/Dashboard'
import Tasks from '@/pages/Tasks'
import TaskDetail from '@/pages/TaskDetail'
import Organizations from '@/pages/Organizations'
import Departments from '@/pages/Departments'
import Users from '@/pages/Users'
import Profile from '@/pages/Profile'
import { AuthProvider } from '@/contexts/AuthContext'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null)

  // 检查本地存储的token
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      setIsAuthenticated(true)
    } else {
      setIsAuthenticated(false)
    }
  }, [])

  // 如果还在检查认证状态，显示加载页面
  if (isAuthenticated === null) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">加载中...</p>
        </div>
      </div>
    )
  }

  return (
    <AuthProvider>
      <Routes>
        {/* 公开路由 */}
        <Route 
          path="/login" 
          element={!isAuthenticated ? <Login /> : <Navigate to="/dashboard" replace />} 
        />
        
        {/* 受保护的路由 */}
        <Route
          path="/*"
          element={
            isAuthenticated ? (
              <Layout>
                <Routes>
                  <Route path="/" element={<Navigate to="/dashboard" replace />} />
                  <Route path="/dashboard" element={<Dashboard />} />
                  <Route path="/tasks" element={<Tasks />} />
                  <Route path="/tasks/:id" element={<TaskDetail />} />
                  <Route path="/organizations" element={<Organizations />} />
                  <Route path="/departments" element={<Departments />} />
                  <Route path="/users" element={<Users />} />
                  <Route path="/profile" element={<Profile />} />
                </Routes>
              </Layout>
            ) : (
              <Navigate to="/login" replace />
            )
          }
        />
      </Routes>
    </AuthProvider>
  )
}

export default App