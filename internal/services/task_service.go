package services

import (
	"errors"
	"fmt"
	"time"

	"task-management-platform/internal/database"
	"task-management-platform/internal/models"

	"gorm.io/gorm"
)

type TaskService struct{}

type CreateTaskRequest struct {
	Title          string             `json:"title" binding:"required"`
	Description    string             `json:"description"`
	Priority       models.Priority    `json:"priority"`
	AssigneeID     *uint              `json:"assignee_id"`
	DueDate        *time.Time         `json:"due_date"`
	EstimatedHours *float64           `json:"estimated_hours"`
	MemberIDs      []uint             `json:"member_ids"`
	Checklists     []CreateChecklistRequest `json:"checklists"`
}

type UpdateTaskRequest struct {
	Title          *string            `json:"title"`
	Description    *string            `json:"description"`
	Status         *models.TaskStatus `json:"status"`
	Priority       *models.Priority   `json:"priority"`
	AssigneeID     *uint              `json:"assignee_id"`
	DueDate        *time.Time         `json:"due_date"`
	EstimatedHours *float64           `json:"estimated_hours"`
	ActualHours    *float64           `json:"actual_hours"`
}

type CreateChecklistRequest struct {
	Title    string `json:"title" binding:"required"`
	Position int    `json:"position"`
}

type UpdateChecklistRequest struct {
	Title       *string `json:"title"`
	IsCompleted *bool   `json:"is_completed"`
	Position    *int    `json:"position"`
}

type TaskQueryOptions struct {
	Status     *models.TaskStatus `json:"status"`
	Priority   *models.Priority   `json:"priority"`
	CreatorID  *uint              `json:"creator_id"`
	AssigneeID *uint              `json:"assignee_id"`
	MemberID   *uint              `json:"member_id"`
	DueFrom    *time.Time         `json:"due_from"`
	DueTo      *time.Time         `json:"due_to"`
	Search     string             `json:"search"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	SortBy     string             `json:"sort_by"`
	SortOrder  string             `json:"sort_order"`
}

type TaskListResponse struct {
	Tasks      []models.Task `json:"tasks"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(userID uint, req *CreateTaskRequest) (*models.Task, error) {
	db := database.GetDB()
	
	// 创建任务
	task := &models.Task{
		Title:          req.Title,
		Description:    req.Description,
		Priority:       req.Priority,
		CreatorID:      userID,
		AssigneeID:     req.AssigneeID,
		DueDate:        req.DueDate,
		EstimatedHours: req.EstimatedHours,
		Status:         models.TaskStatusTodo,
	}
	
	// 如果没有指定优先级，默认为中等
	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}
	
	// 开始事务
	tx := db.Begin()
	
	// 创建任务
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	
	// 添加任务成员
	if len(req.MemberIDs) > 0 {
		for _, memberID := range req.MemberIDs {
			taskMember := &models.TaskMember{
				TaskID:   task.ID,
				UserID:   memberID,
				JoinedAt: time.Now(),
			}
			if err := tx.Create(taskMember).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	
	// 创建任务清单
	if len(req.Checklists) > 0 {
		for _, checklistReq := range req.Checklists {
			checklist := &models.Checklist{
				TaskID:   task.ID,
				Title:    checklistReq.Title,
				Position: checklistReq.Position,
			}
			if err := tx.Create(checklist).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	
	// 重新加载任务数据
	if err := s.loadTaskRelations(db, task); err != nil {
		return nil, err
	}
	
	// 发送通知
	go s.sendTaskNotification(task, "task_created")
	
	return task, nil
}

// GetTaskByID 根据ID获取任务
func (s *TaskService) GetTaskByID(taskID uint, userID uint) (*models.Task, error) {
	db := database.GetDB()
	
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, errors.New("任务不存在")
	}
	
	// 检查权限
	if !s.canAccessTask(userID, &task) {
		return nil, errors.New("无权访问此任务")
	}
	
	// 加载关联数据
	if err := s.loadTaskRelations(db, &task); err != nil {
		return nil, err
	}
	
	return &task, nil
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(taskID uint, userID uint, req *UpdateTaskRequest) (*models.Task, error) {
	db := database.GetDB()
	
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, errors.New("任务不存在")
	}
	
	// 检查权限
	if !s.canModifyTask(userID, &task) {
		return nil, errors.New("无权修改此任务")
	}
	
	// 记录状态变更
	oldStatus := task.Status
	
	// 更新任务字段
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
		// 如果任务完成，记录完成时间
		if *req.Status == models.TaskStatusCompleted && task.CompletedAt == nil {
			now := time.Now()
			task.CompletedAt = &now
		}
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.EstimatedHours != nil {
		task.EstimatedHours = req.EstimatedHours
	}
	if req.ActualHours != nil {
		task.ActualHours = req.ActualHours
	}
	
	// 保存更新
	if err := db.Save(&task).Error; err != nil {
		return nil, err
	}
	
	// 重新加载任务数据
	if err := s.loadTaskRelations(db, &task); err != nil {
		return nil, err
	}
	
	// 如果状态发生变化，发送通知
	if req.Status != nil && oldStatus != *req.Status {
		go s.sendTaskNotification(&task, "task_status_changed")
	}
	
	return &task, nil
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(taskID uint, userID uint) error {
	db := database.GetDB()
	
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}
	
	// 检查权限
	if !s.canModifyTask(userID, &task) {
		return errors.New("无权删除此任务")
	}
	
	// 软删除任务
	return db.Delete(&task).Error
}

// GetTasks 获取任务列表
func (s *TaskService) GetTasks(userID uint, options *TaskQueryOptions) (*TaskListResponse, error) {
	db := database.GetDB()
	
	// 构建查询
	query := db.Model(&models.Task{})
	
	// 根据用户权限过滤任务
	user, err := s.getUserByID(userID)
	if err != nil {
		return nil, err
	}
	
	// 如果不是管理员，只能看到相关的任务
	if user.Role != models.RoleAdmin {
		query = query.Where("creator_id = ? OR assignee_id = ? OR id IN (SELECT task_id FROM task_members WHERE user_id = ?)", 
			userID, userID, userID)
	}
	
	// 应用过滤条件
	if options.Status != nil {
		query = query.Where("status = ?", *options.Status)
	}
	if options.Priority != nil {
		query = query.Where("priority = ?", *options.Priority)
	}
	if options.CreatorID != nil {
		query = query.Where("creator_id = ?", *options.CreatorID)
	}
	if options.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *options.AssigneeID)
	}
	if options.MemberID != nil {
		query = query.Where("id IN (SELECT task_id FROM task_members WHERE user_id = ?)", *options.MemberID)
	}
	if options.DueFrom != nil {
		query = query.Where("due_date >= ?", *options.DueFrom)
	}
	if options.DueTo != nil {
		query = query.Where("due_date <= ?", *options.DueTo)
	}
	if options.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+options.Search+"%", "%"+options.Search+"%")
	}
	
	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 设置默认分页参数
	if options.Page <= 0 {
		options.Page = 1
	}
	if options.PageSize <= 0 {
		options.PageSize = 20
	}
	
	// 设置排序
	orderBy := "created_at DESC"
	if options.SortBy != "" {
		order := "ASC"
		if options.SortOrder == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("%s %s", options.SortBy, order)
	}
	
	// 分页查询
	offset := (options.Page - 1) * options.PageSize
	var tasks []models.Task
	if err := query.Order(orderBy).Offset(offset).Limit(options.PageSize).
		Preload("Creator").
		Preload("Assignee").
		Preload("Members.User").
		Preload("Checklists").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	
	totalPages := int((total + int64(options.PageSize) - 1) / int64(options.PageSize))
	
	return &TaskListResponse{
		Tasks:      tasks,
		Total:      total,
		Page:       options.Page,
		PageSize:   options.PageSize,
		TotalPages: totalPages,
	}, nil
}

// AddTaskMember 添加任务成员
func (s *TaskService) AddTaskMember(taskID uint, userID uint, memberID uint) error {
	db := database.GetDB()
	
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}
	
	// 检查权限
	if !s.canModifyTask(userID, &task) {
		return errors.New("无权修改此任务")
	}
	
	// 检查成员是否已存在
	var existingMember models.TaskMember
	if err := db.Where("task_id = ? AND user_id = ?", taskID, memberID).First(&existingMember).Error; err == nil {
		return errors.New("用户已是任务成员")
	}
	
	// 添加成员
	taskMember := &models.TaskMember{
		TaskID:   taskID,
		UserID:   memberID,
		JoinedAt: time.Now(),
	}
	
	return db.Create(taskMember).Error
}

// RemoveTaskMember 移除任务成员
func (s *TaskService) RemoveTaskMember(taskID uint, userID uint, memberID uint) error {
	db := database.GetDB()
	
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}
	
	// 检查权限
	if !s.canModifyTask(userID, &task) {
		return errors.New("无权修改此任务")
	}
	
	// 删除成员
	return db.Where("task_id = ? AND user_id = ?", taskID, memberID).Delete(&models.TaskMember{}).Error
}

// CreateChecklist 创建任务清单
func (s *TaskService) CreateChecklist(taskID uint, userID uint, req *CreateChecklistRequest) (*models.Checklist, error) {
	db := database.GetDB()
	
	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, errors.New("任务不存在")
	}
	
	// 检查权限
	if !s.canAccessTask(userID, &task) {
		return nil, errors.New("无权访问此任务")
	}
	
	checklist := &models.Checklist{
		TaskID:   taskID,
		Title:    req.Title,
		Position: req.Position,
	}
	
	if err := db.Create(checklist).Error; err != nil {
		return nil, err
	}
	
	return checklist, nil
}

// UpdateChecklist 更新任务清单
func (s *TaskService) UpdateChecklist(checklistID uint, userID uint, req *UpdateChecklistRequest) (*models.Checklist, error) {
	db := database.GetDB()
	
	var checklist models.Checklist
	if err := db.Preload("Task").Where("id = ?", checklistID).First(&checklist).Error; err != nil {
		return nil, errors.New("清单项不存在")
	}
	
	// 检查权限
	if !s.canAccessTask(userID, &checklist.Task) {
		return nil, errors.New("无权访问此任务")
	}
	
	// 更新字段
	if req.Title != nil {
		checklist.Title = *req.Title
	}
	if req.IsCompleted != nil {
		checklist.IsCompleted = *req.IsCompleted
	}
	if req.Position != nil {
		checklist.Position = *req.Position
	}
	
	if err := db.Save(&checklist).Error; err != nil {
		return nil, err
	}
	
	return &checklist, nil
}

// DeleteChecklist 删除任务清单
func (s *TaskService) DeleteChecklist(checklistID uint, userID uint) error {
	db := database.GetDB()
	
	var checklist models.Checklist
	if err := db.Preload("Task").Where("id = ?", checklistID).First(&checklist).Error; err != nil {
		return errors.New("清单项不存在")
	}
	
	// 检查权限
	if !s.canModifyTask(userID, &checklist.Task) {
		return errors.New("无权修改此任务")
	}
	
	return db.Delete(&checklist).Error
}

// GetOverdueTasks 获取超期任务
func (s *TaskService) GetOverdueTasks() ([]models.Task, error) {
	db := database.GetDB()
	
	var tasks []models.Task
	now := time.Now()
	
	err := db.Where("due_date < ? AND status NOT IN (?)", now, []models.TaskStatus{
		models.TaskStatusCompleted,
		models.TaskStatusCancelled,
	}).Preload("Creator").Preload("Assignee").Find(&tasks).Error
	
	return tasks, err
}

// MarkTasksOverdue 标记任务为超期
func (s *TaskService) MarkTasksOverdue() error {
	db := database.GetDB()
	now := time.Now()
	
	return db.Model(&models.Task{}).
		Where("due_date < ? AND status NOT IN (?)", now, []models.TaskStatus{
			models.TaskStatusCompleted,
			models.TaskStatusCancelled,
			models.TaskStatusOverdue,
		}).
		Update("status", models.TaskStatusOverdue).Error
}

// 辅助方法

func (s *TaskService) loadTaskRelations(db *gorm.DB, task *models.Task) error {
	return db.Preload("Creator").
		Preload("Assignee").
		Preload("Members.User").
		Preload("Checklists").
		Preload("Comments.User").
		Preload("Attachments").
		First(task, task.ID).Error
}

func (s *TaskService) canAccessTask(userID uint, task *models.Task) bool {
	user, err := s.getUserByID(userID)
	if err != nil {
		return false
	}
	
	// 管理员可以访问所有任务
	if user.Role == models.RoleAdmin {
		return true
	}
	
	// 任务创建者、负责人或成员可以访问
	if task.CreatorID == userID || (task.AssigneeID != nil && *task.AssigneeID == userID) {
		return true
	}
	
	// 检查是否是任务成员
	db := database.GetDB()
	var count int64
	db.Model(&models.TaskMember{}).Where("task_id = ? AND user_id = ?", task.ID, userID).Count(&count)
	return count > 0
}

func (s *TaskService) canModifyTask(userID uint, task *models.Task) bool {
	user, err := s.getUserByID(userID)
	if err != nil {
		return false
	}
	
	// 管理员可以修改所有任务
	if user.Role == models.RoleAdmin {
		return true
	}
	
	// 任务创建者可以修改
	if task.CreatorID == userID {
		return true
	}
	
	// 负责人可以修改部分内容
	if task.AssigneeID != nil && *task.AssigneeID == userID {
		return true
	}
	
	return false
}

func (s *TaskService) getUserByID(userID uint) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *TaskService) sendTaskNotification(task *models.Task, notificationType string) {
	// 这里可以实现发送通知的逻辑
	// 例如邮件、站内消息等
	fmt.Printf("Sending notification: %s for task %d\n", notificationType, task.ID)
}