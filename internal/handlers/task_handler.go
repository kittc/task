package handlers

import (
	"net/http"
	"strconv"
	"time"

	"task-management-platform/internal/middleware"
	"task-management-platform/internal/models"
	"task-management-platform/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	var req services.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	task, err := h.taskService.CreateTask(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": task,
	})
}

// GetTasks 获取任务列表
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	// 解析查询参数
	options := &services.TaskQueryOptions{}
	
	if status := c.Query("status"); status != "" {
		taskStatus := models.TaskStatus(status)
		options.Status = &taskStatus
	}
	
	if priority := c.Query("priority"); priority != "" {
		taskPriority := models.Priority(priority)
		options.Priority = &taskPriority
	}
	
	if creatorID := c.Query("creator_id"); creatorID != "" {
		if id, err := strconv.ParseUint(creatorID, 10, 32); err == nil {
			creatorIDUint := uint(id)
			options.CreatorID = &creatorIDUint
		}
	}
	
	if assigneeID := c.Query("assignee_id"); assigneeID != "" {
		if id, err := strconv.ParseUint(assigneeID, 10, 32); err == nil {
			assigneeIDUint := uint(id)
			options.AssigneeID = &assigneeIDUint
		}
	}
	
	if memberID := c.Query("member_id"); memberID != "" {
		if id, err := strconv.ParseUint(memberID, 10, 32); err == nil {
			memberIDUint := uint(id)
			options.MemberID = &memberIDUint
		}
	}
	
	if dueFrom := c.Query("due_from"); dueFrom != "" {
		if t, err := time.Parse("2006-01-02", dueFrom); err == nil {
			options.DueFrom = &t
		}
	}
	
	if dueTo := c.Query("due_to"); dueTo != "" {
		if t, err := time.Parse("2006-01-02", dueTo); err == nil {
			options.DueTo = &t
		}
	}
	
	options.Search = c.Query("search")
	
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			options.Page = p
		}
	}
	
	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil {
			options.PageSize = ps
		}
	}
	
	options.SortBy = c.Query("sort_by")
	options.SortOrder = c.Query("sort_order")

	response, err := h.taskService.GetTasks(userID, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": response,
	})
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的任务ID",
		})
		return
	}

	task, err := h.taskService.GetTaskByID(uint(taskID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": task,
	})
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的任务ID",
		})
		return
	}

	var req services.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	task, err := h.taskService.UpdateTask(uint(taskID), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": task,
	})
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的任务ID",
		})
		return
	}

	err = h.taskService.DeleteTask(uint(taskID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "任务删除成功",
	})
}

// AddTaskMember 添加任务成员
func (h *TaskHandler) AddTaskMember(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的任务ID",
		})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	err = h.taskService.AddTaskMember(uint(taskID), userID, req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成员添加成功",
	})
}

// RemoveTaskMember 移除任务成员
func (h *TaskHandler) RemoveTaskMember(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的任务ID",
		})
		return
	}

	memberIDStr := c.Param("member_id")
	memberID, err := strconv.ParseUint(memberIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的成员ID",
		})
		return
	}

	err = h.taskService.RemoveTaskMember(uint(taskID), userID, uint(memberID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成员移除成功",
	})
}

// CreateChecklist 创建任务清单
func (h *TaskHandler) CreateChecklist(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的任务ID",
		})
		return
	}

	var req services.CreateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	checklist, err := h.taskService.CreateChecklist(uint(taskID), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": checklist,
	})
}

// UpdateChecklist 更新任务清单
func (h *TaskHandler) UpdateChecklist(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	checklistIDStr := c.Param("checklist_id")
	checklistID, err := strconv.ParseUint(checklistIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的清单ID",
		})
		return
	}

	var req services.UpdateChecklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	checklist, err := h.taskService.UpdateChecklist(uint(checklistID), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": checklist,
	})
}

// DeleteChecklist 删除任务清单
func (h *TaskHandler) DeleteChecklist(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未认证",
		})
		return
	}

	checklistIDStr := c.Param("checklist_id")
	checklistID, err := strconv.ParseUint(checklistIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的清单ID",
		})
		return
	}

	err = h.taskService.DeleteChecklist(uint(checklistID), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "清单删除成功",
	})
}