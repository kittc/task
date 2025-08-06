package database

import (
	"fmt"
	"log"

	"task-management-platform/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(dbPath string) error {
	var err error
	
	// 配置 GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	
	// 连接 SQLite 数据库
	DB, err = gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	
	// 自动迁移数据库表
	err = AutoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}
	
	// 初始化默认数据
	err = SeedData()
	if err != nil {
		return fmt.Errorf("failed to seed data: %v", err)
	}
	
	log.Println("Database initialized successfully")
	return nil
}

// AutoMigrate 自动迁移数据库表
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Organization{},
		&models.Department{},
		&models.User{},
		&models.Task{},
		&models.TaskMember{},
		&models.Checklist{},
		&models.Comment{},
		&models.Attachment{},
		&models.Notification{},
		&models.SystemConfig{},
	)
}

// SeedData 初始化默认数据
func SeedData() error {
	// 检查是否已有数据
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil // 已有数据，跳过初始化
	}
	
	// 创建默认组织
	org := &models.Organization{
		Name:        "示例公司",
		Description: "这是一个示例组织",
	}
	if err := DB.Create(org).Error; err != nil {
		return err
	}
	
	// 创建默认部门
	departments := []*models.Department{
		{
			OrganizationID: org.ID,
			Name:          "技术部",
			Description:   "负责技术开发和维护",
		},
		{
			OrganizationID: org.ID,
			Name:          "产品部",
			Description:   "负责产品设计和规划",
		},
		{
			OrganizationID: org.ID,
			Name:          "运营部",
			Description:   "负责日常运营管理",
		},
	}
	
	for _, dept := range departments {
		if err := DB.Create(dept).Error; err != nil {
			return err
		}
	}
	
	// 创建默认管理员用户
	admin := &models.User{
		DepartmentID: &departments[0].ID,
		Username:     "admin",
		Email:        "admin@example.com",
		Password:     "$2a$10$8X2vPf1M2QJ3bHv9g8qQQu1.C3d4R5a6b7c8d9e0f1g2h3i4j5k6l", // password: admin123
		FullName:     "系统管理员",
		Role:         models.RoleAdmin,
		IsActive:     true,
	}
	if err := DB.Create(admin).Error; err != nil {
		return err
	}
	
	// 创建示例用户
	users := []*models.User{
		{
			DepartmentID: &departments[0].ID,
			Username:     "manager1",
			Email:        "manager1@example.com",
			Password:     "$2a$10$8X2vPf1M2QJ3bHv9g8qQQu1.C3d4R5a6b7c8d9e0f1g2h3i4j5k6l",
			FullName:     "张经理",
			Role:         models.RoleManager,
			IsActive:     true,
		},
		{
			DepartmentID: &departments[0].ID,
			Username:     "supervisor1",
			Email:        "supervisor1@example.com",
			Password:     "$2a$10$8X2vPf1M2QJ3bHv9g8qQQu1.C3d4R5a6b7c8d9e0f1g2h3i4j5k6l",
			FullName:     "李主管",
			Role:         models.RoleSupervisor,
			IsActive:     true,
		},
		{
			DepartmentID: &departments[0].ID,
			Username:     "employee1",
			Email:        "employee1@example.com",
			Password:     "$2a$10$8X2vPf1M2QJ3bHv9g8qQQu1.C3d4R5a6b7c8d9e0f1g2h3i4j5k6l",
			FullName:     "王职员",
			Role:         models.RoleEmployee,
			IsActive:     true,
		},
	}
	
	for _, user := range users {
		if err := DB.Create(user).Error; err != nil {
			return err
		}
	}
	
	// 创建系统配置
	configs := []*models.SystemConfig{
		{
			Key:         "task_reminder_hours",
			Value:       "24",
			Description: "任务到期前提醒时间（小时）",
		},
		{
			Key:         "overdue_check_interval",
			Value:       "60",
			Description: "超时检查间隔（分钟）",
		},
		{
			Key:         "max_file_size",
			Value:       "10485760",
			Description: "最大文件上传大小（字节）",
		},
	}
	
	for _, config := range configs {
		if err := DB.Create(config).Error; err != nil {
			return err
		}
	}
	
	log.Println("Default data seeded successfully")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}