package models

import (
	"time"

	"gorm.io/gorm"
)

// 用户角色枚举
type Role string

const (
	RoleEmployee Role = "employee" // 职员
	RoleSupervisor Role = "supervisor" // 主管
	RoleManager Role = "manager" // 经理
	RoleAdmin Role = "admin" // 系统管理员
)

// 任务状态枚举
type TaskStatus string

const (
	TaskStatusTodo TaskStatus = "todo" // 代办
	TaskStatusInProgress TaskStatus = "in_progress" // 处理中
	TaskStatusOverdue TaskStatus = "overdue" // 延期
	TaskStatusPaused TaskStatus = "paused" // 暂停
	TaskStatusCancelled TaskStatus = "cancelled" // 取消
	TaskStatusCompleted TaskStatus = "completed" // 完成
)

// 优先级枚举
type Priority string

const (
	PriorityLow Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// 组织机构
type Organization struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	Description string         `json:"description" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Departments []Department `json:"departments" gorm:"foreignKey:OrganizationID"`
}

// 部门
type Department struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	OrganizationID uint           `json:"organization_id" gorm:"not null"`
	Name           string         `json:"name" gorm:"not null;size:100"`
	Description    string         `json:"description" gorm:"size:500"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Organization Organization `json:"organization" gorm:"foreignKey:OrganizationID"`
	Users        []User      `json:"users" gorm:"foreignKey:DepartmentID"`
}

// 用户
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	DepartmentID *uint          `json:"department_id"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null;size:100"`
	Password     string         `json:"-" gorm:"not null"`
	FullName     string         `json:"full_name" gorm:"not null;size:100"`
	Avatar       string         `json:"avatar" gorm:"size:255"`
	Role         Role           `json:"role" gorm:"not null;default:'employee'"`
	IsActive     bool           `json:"is_active" gorm:"default:true"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Department        *Department      `json:"department" gorm:"foreignKey:DepartmentID"`
	CreatedTasks      []Task          `json:"created_tasks" gorm:"foreignKey:CreatorID"`
	AssignedTasks     []Task          `json:"assigned_tasks" gorm:"foreignKey:AssigneeID"`
	TaskMembers       []TaskMember    `json:"task_members"`
	Comments          []Comment       `json:"comments"`
	Notifications     []Notification  `json:"notifications"`
}

// 任务
type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;size:200"`
	Description string         `json:"description" gorm:"type:text"`
	Status      TaskStatus     `json:"status" gorm:"not null;default:'todo'"`
	Priority    Priority       `json:"priority" gorm:"not null;default:'medium'"`
	CreatorID   uint           `json:"creator_id" gorm:"not null"`
	AssigneeID  *uint          `json:"assignee_id"`
	DueDate     *time.Time     `json:"due_date"`
	CompletedAt *time.Time     `json:"completed_at"`
	EstimatedHours *float64    `json:"estimated_hours"`
	ActualHours    *float64    `json:"actual_hours"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Creator     User          `json:"creator" gorm:"foreignKey:CreatorID"`
	Assignee    *User         `json:"assignee" gorm:"foreignKey:AssigneeID"`
	Members     []TaskMember  `json:"members"`
	Checklists  []Checklist   `json:"checklists"`
	Comments    []Comment     `json:"comments"`
	Attachments []Attachment  `json:"attachments"`
}

// 任务成员
type TaskMember struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	TaskID   uint      `json:"task_id" gorm:"not null"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	JoinedAt time.Time `json:"joined_at"`
	
	// 关联关系
	Task Task `json:"task" gorm:"foreignKey:TaskID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// 任务清单
type Checklist struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TaskID      uint           `json:"task_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null;size:200"`
	IsCompleted bool           `json:"is_completed" gorm:"default:false"`
	Position    int            `json:"position" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Task Task `json:"task" gorm:"foreignKey:TaskID"`
}

// 评论
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TaskID    uint           `json:"task_id" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Content   string         `json:"content" gorm:"not null;type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Task Task `json:"task" gorm:"foreignKey:TaskID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// 附件
type Attachment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	TaskID    uint           `json:"task_id" gorm:"not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	FileName  string         `json:"file_name" gorm:"not null;size:255"`
	FilePath  string         `json:"file_path" gorm:"not null;size:500"`
	FileSize  int64          `json:"file_size"`
	MimeType  string         `json:"mime_type" gorm:"size:100"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// 关联关系
	Task Task `json:"task" gorm:"foreignKey:TaskID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// 通知
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	TaskID    *uint          `json:"task_id"`
	Type      string         `json:"type" gorm:"not null;size:50"` // task_assigned, task_due, task_overdue, task_completed
	Title     string         `json:"title" gorm:"not null;size:200"`
	Content   string         `json:"content" gorm:"type:text"`
	IsRead    bool           `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	
	// 关联关系
	User User  `json:"user" gorm:"foreignKey:UserID"`
	Task *Task `json:"task" gorm:"foreignKey:TaskID"`
}

// 系统配置
type SystemConfig struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"uniqueIndex;not null;size:100"`
	Value string `json:"value" gorm:"type:text"`
	Description string `json:"description" gorm:"size:500"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}