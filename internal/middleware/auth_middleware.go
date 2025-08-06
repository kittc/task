package middleware

import (
	"net/http"
	"strings"

	"task-management-platform/internal/models"
	"task-management-platform/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth 需要认证的中间件
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证令牌",
			})
			c.Abort()
			return
		}

		user, err := m.authService.GetUserByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		
		c.Next()
	}
}

// RequireRole 需要特定角色的中间件
func (m *AuthMiddleware) RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未认证",
			})
			c.Abort()
			return
		}

		role := userRole.(models.Role)
		
		// 管理员拥有所有权限
		if role == models.RoleAdmin {
			c.Next()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		for _, allowedRole := range roles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "权限不足",
		})
		c.Abort()
	}
}

// RequireAdminOrManager 需要管理员或经理权限
func (m *AuthMiddleware) RequireAdminOrManager() gin.HandlerFunc {
	return m.RequireRole(models.RoleAdmin, models.RoleManager)
}

// RequireAdminOrManagerOrSupervisor 需要管理员、经理或主管权限
func (m *AuthMiddleware) RequireAdminOrManagerOrSupervisor() gin.HandlerFunc {
	return m.RequireRole(models.RoleAdmin, models.RoleManager, models.RoleSupervisor)
}

// OptionalAuth 可选认证中间件（用户可以是游客）
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			user, err := m.authService.GetUserByToken(token)
			if err == nil {
				c.Set("user", user)
				c.Set("user_id", user.ID)
				c.Set("user_role", user.Role)
			}
		}
		c.Next()
	}
}

// extractToken 从请求中提取token
func extractToken(c *gin.Context) string {
	// 首先从 Authorization header 中获取
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	// 然后从查询参数中获取
	token := c.Query("token")
	if token != "" {
		return token
	}

	// 最后从 cookie 中获取
	cookie, err := c.Cookie("token")
	if err == nil {
		return cookie
	}

	return ""
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// GetCurrentUserID 获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetCurrentUserRole 获取当前用户角色
func GetCurrentUserRole(c *gin.Context) (models.Role, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	return role.(models.Role), true
}