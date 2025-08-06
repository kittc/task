package handlers

import (
	"net/http"
	"strconv"

	"task-management-platform/internal/middleware"
	"task-management-platform/internal/models"
	"task-management-platform/internal/services"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgService *services.OrganizationService
}

func NewOrganizationHandler(orgService *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// Organization handlers

// CreateOrganization 创建组织
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req services.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	org, err := h.orgService.CreateOrganization(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": org,
	})
}

// GetOrganizations 获取组织列表
func (h *OrganizationHandler) GetOrganizations(c *gin.Context) {
	orgs, err := h.orgService.GetOrganizations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": orgs,
	})
}

// GetOrganization 获取组织详情
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的组织ID",
		})
		return
	}

	org, err := h.orgService.GetOrganizationByID(uint(orgID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": org,
	})
}

// UpdateOrganization 更新组织
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的组织ID",
		})
		return
	}

	var req services.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	org, err := h.orgService.UpdateOrganization(uint(orgID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": org,
	})
}

// DeleteOrganization 删除组织
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的组织ID",
		})
		return
	}

	err = h.orgService.DeleteOrganization(uint(orgID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "组织删除成功",
	})
}

// GetOrganizationStats 获取组织统计信息
func (h *OrganizationHandler) GetOrganizationStats(c *gin.Context) {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的组织ID",
		})
		return
	}

	// 这里可以实现组织统计逻辑
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id": orgID,
			"message": "组织统计功能待实现",
		},
	})
}

// Department handlers

// CreateDepartment 创建部门
func (h *OrganizationHandler) CreateDepartment(c *gin.Context) {
	var req services.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	dept, err := h.orgService.CreateDepartment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": dept,
	})
}

// GetDepartments 获取部门列表
func (h *OrganizationHandler) GetDepartments(c *gin.Context) {
	var orgID *uint
	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		if id, err := strconv.ParseUint(orgIDStr, 10, 32); err == nil {
			orgIDUint := uint(id)
			orgID = &orgIDUint
		}
	}

	depts, err := h.orgService.GetDepartments(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": depts,
	})
}

// GetDepartment 获取部门详情
func (h *OrganizationHandler) GetDepartment(c *gin.Context) {
	deptIDStr := c.Param("id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的部门ID",
		})
		return
	}

	dept, err := h.orgService.GetDepartmentByID(uint(deptID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dept,
	})
}

// UpdateDepartment 更新部门
func (h *OrganizationHandler) UpdateDepartment(c *gin.Context) {
	deptIDStr := c.Param("id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的部门ID",
		})
		return
	}

	var req services.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	dept, err := h.orgService.UpdateDepartment(uint(deptID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": dept,
	})
}

// DeleteDepartment 删除部门
func (h *OrganizationHandler) DeleteDepartment(c *gin.Context) {
	deptIDStr := c.Param("id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的部门ID",
		})
		return
	}

	err = h.orgService.DeleteDepartment(uint(deptID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "部门删除成功",
	})
}

// GetDepartmentStats 获取部门统计信息
func (h *OrganizationHandler) GetDepartmentStats(c *gin.Context) {
	deptIDStr := c.Param("id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的部门ID",
		})
		return
	}

	stats, err := h.orgService.GetDepartmentStats(uint(deptID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": stats,
	})
}

// User handlers

// CreateUser 创建用户
func (h *OrganizationHandler) CreateUser(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	user, err := h.orgService.CreateUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": user,
	})
}

// GetUsers 获取用户列表
func (h *OrganizationHandler) GetUsers(c *gin.Context) {
	var deptID *uint
	var role *models.Role
	var isActive *bool

	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		if id, err := strconv.ParseUint(deptIDStr, 10, 32); err == nil {
			deptIDUint := uint(id)
			deptID = &deptIDUint
		}
	}

	if roleStr := c.Query("role"); roleStr != "" {
		roleValue := models.Role(roleStr)
		role = &roleValue
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	users, err := h.orgService.GetUsers(deptID, role, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": users,
	})
}

// GetUser 获取用户详情
func (h *OrganizationHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	user, err := h.orgService.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": user,
	})
}

// UpdateUser 更新用户
func (h *OrganizationHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	user, err := h.orgService.UpdateUser(uint(userID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": user,
	})
}

// DeleteUser 删除用户
func (h *OrganizationHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	err = h.orgService.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户删除成功",
	})
}

// TransferUser 转移用户到其他部门
func (h *OrganizationHandler) TransferUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	var req struct {
		DepartmentID uint `json:"department_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	err = h.orgService.TransferUser(uint(userID), req.DepartmentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "用户转移成功",
	})
}

// GetOrganizationStructure 获取组织架构
func (h *OrganizationHandler) GetOrganizationStructure(c *gin.Context) {
	structure, err := h.orgService.GetOrganizationStructure()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": structure,
	})
}