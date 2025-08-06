package services

import (
	"errors"

	"task-management-platform/internal/database"
	"task-management-platform/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type OrganizationService struct{}

type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateOrganizationRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CreateDepartmentRequest struct {
	OrganizationID uint   `json:"organization_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
}

type UpdateDepartmentRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type CreateUserRequest struct {
	DepartmentID *uint       `json:"department_id"`
	Username     string      `json:"username" binding:"required"`
	Email        string      `json:"email" binding:"required,email"`
	Password     string      `json:"password" binding:"required,min=6"`
	FullName     string      `json:"full_name" binding:"required"`
	Role         models.Role `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	DepartmentID *uint       `json:"department_id"`
	FullName     *string     `json:"full_name"`
	Role         *models.Role `json:"role"`
	IsActive     *bool       `json:"is_active"`
}

func NewOrganizationService() *OrganizationService {
	return &OrganizationService{}
}

// Organization management

// CreateOrganization 创建组织
func (s *OrganizationService) CreateOrganization(req *CreateOrganizationRequest) (*models.Organization, error) {
	db := database.GetDB()
	
	organization := &models.Organization{
		Name:        req.Name,
		Description: req.Description,
	}
	
	if err := db.Create(organization).Error; err != nil {
		return nil, err
	}
	
	return organization, nil
}

// GetOrganizations 获取组织列表
func (s *OrganizationService) GetOrganizations() ([]models.Organization, error) {
	db := database.GetDB()
	
	var organizations []models.Organization
	if err := db.Preload("Departments").Find(&organizations).Error; err != nil {
		return nil, err
	}
	
	return organizations, nil
}

// GetOrganizationByID 根据ID获取组织
func (s *OrganizationService) GetOrganizationByID(orgID uint) (*models.Organization, error) {
	db := database.GetDB()
	
	var organization models.Organization
	if err := db.Preload("Departments.Users").First(&organization, orgID).Error; err != nil {
		return nil, errors.New("组织不存在")
	}
	
	return &organization, nil
}

// UpdateOrganization 更新组织
func (s *OrganizationService) UpdateOrganization(orgID uint, req *UpdateOrganizationRequest) (*models.Organization, error) {
	db := database.GetDB()
	
	var organization models.Organization
	if err := db.First(&organization, orgID).Error; err != nil {
		return nil, errors.New("组织不存在")
	}
	
	if req.Name != nil {
		organization.Name = *req.Name
	}
	if req.Description != nil {
		organization.Description = *req.Description
	}
	
	if err := db.Save(&organization).Error; err != nil {
		return nil, err
	}
	
	return &organization, nil
}

// DeleteOrganization 删除组织
func (s *OrganizationService) DeleteOrganization(orgID uint) error {
	db := database.GetDB()
	
	var organization models.Organization
	if err := db.First(&organization, orgID).Error; err != nil {
		return errors.New("组织不存在")
	}
	
	// 检查是否有部门
	var deptCount int64
	db.Model(&models.Department{}).Where("organization_id = ?", orgID).Count(&deptCount)
	if deptCount > 0 {
		return errors.New("组织下还有部门，无法删除")
	}
	
	return db.Delete(&organization).Error
}

// Department management

// CreateDepartment 创建部门
func (s *OrganizationService) CreateDepartment(req *CreateDepartmentRequest) (*models.Department, error) {
	db := database.GetDB()
	
	// 检查组织是否存在
	var organization models.Organization
	if err := db.First(&organization, req.OrganizationID).Error; err != nil {
		return nil, errors.New("组织不存在")
	}
	
	department := &models.Department{
		OrganizationID: req.OrganizationID,
		Name:           req.Name,
		Description:    req.Description,
	}
	
	if err := db.Create(department).Error; err != nil {
		return nil, err
	}
	
	// 预加载组织信息
	db.Preload("Organization").First(department, department.ID)
	
	return department, nil
}

// GetDepartments 获取部门列表
func (s *OrganizationService) GetDepartments(orgID *uint) ([]models.Department, error) {
	db := database.GetDB()
	
	query := db.Preload("Organization").Preload("Users")
	if orgID != nil {
		query = query.Where("organization_id = ?", *orgID)
	}
	
	var departments []models.Department
	if err := query.Find(&departments).Error; err != nil {
		return nil, err
	}
	
	return departments, nil
}

// GetDepartmentByID 根据ID获取部门
func (s *OrganizationService) GetDepartmentByID(deptID uint) (*models.Department, error) {
	db := database.GetDB()
	
	var department models.Department
	if err := db.Preload("Organization").Preload("Users").First(&department, deptID).Error; err != nil {
		return nil, errors.New("部门不存在")
	}
	
	return &department, nil
}

// UpdateDepartment 更新部门
func (s *OrganizationService) UpdateDepartment(deptID uint, req *UpdateDepartmentRequest) (*models.Department, error) {
	db := database.GetDB()
	
	var department models.Department
	if err := db.First(&department, deptID).Error; err != nil {
		return nil, errors.New("部门不存在")
	}
	
	if req.Name != nil {
		department.Name = *req.Name
	}
	if req.Description != nil {
		department.Description = *req.Description
	}
	
	if err := db.Save(&department).Error; err != nil {
		return nil, err
	}
	
	// 预加载关联数据
	db.Preload("Organization").First(&department, department.ID)
	
	return &department, nil
}

// DeleteDepartment 删除部门
func (s *OrganizationService) DeleteDepartment(deptID uint) error {
	db := database.GetDB()
	
	var department models.Department
	if err := db.First(&department, deptID).Error; err != nil {
		return errors.New("部门不存在")
	}
	
	// 检查是否有用户
	var userCount int64
	db.Model(&models.User{}).Where("department_id = ?", deptID).Count(&userCount)
	if userCount > 0 {
		return errors.New("部门下还有用户，无法删除")
	}
	
	return db.Delete(&department).Error
}

// User management

// CreateUser 创建用户
func (s *OrganizationService) CreateUser(req *CreateUserRequest) (*models.User, error) {
	db := database.GetDB()
	
	// 检查用户名是否已存在
	var existingUser models.User
	if err := db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}
	
	// 如果指定了部门，检查部门是否存在
	if req.DepartmentID != nil {
		var department models.Department
		if err := db.First(&department, *req.DepartmentID).Error; err != nil {
			return nil, errors.New("部门不存在")
		}
	}
	
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	user := &models.User{
		DepartmentID: req.DepartmentID,
		Username:     req.Username,
		Email:        req.Email,
		Password:     string(hashedPassword),
		FullName:     req.FullName,
		Role:         req.Role,
		IsActive:     true,
	}
	
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	
	// 预加载部门信息
	db.Preload("Department").First(user, user.ID)
	
	return user, nil
}

// GetUsers 获取用户列表
func (s *OrganizationService) GetUsers(deptID *uint, role *models.Role, isActive *bool) ([]models.User, error) {
	db := database.GetDB()
	
	query := db.Preload("Department.Organization")
	
	if deptID != nil {
		query = query.Where("department_id = ?", *deptID)
	}
	if role != nil {
		query = query.Where("role = ?", *role)
	}
	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}
	
	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	
	return users, nil
}

// GetUserByID 根据ID获取用户
func (s *OrganizationService) GetUserByID(userID uint) (*models.User, error) {
	db := database.GetDB()
	
	var user models.User
	if err := db.Preload("Department.Organization").First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	
	return &user, nil
}

// UpdateUser 更新用户
func (s *OrganizationService) UpdateUser(userID uint, req *UpdateUserRequest) (*models.User, error) {
	db := database.GetDB()
	
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	
	// 如果要更改部门，检查部门是否存在
	if req.DepartmentID != nil && *req.DepartmentID != 0 {
		var department models.Department
		if err := db.First(&department, *req.DepartmentID).Error; err != nil {
			return nil, errors.New("部门不存在")
		}
		user.DepartmentID = req.DepartmentID
	}
	
	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}
	
	// 预加载关联数据
	db.Preload("Department.Organization").First(&user, user.ID)
	
	return &user, nil
}

// DeleteUser 删除用户
func (s *OrganizationService) DeleteUser(userID uint) error {
	db := database.GetDB()
	
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	
	// 检查用户是否有未完成的任务
	var taskCount int64
	db.Model(&models.Task{}).Where("(creator_id = ? OR assignee_id = ?) AND status NOT IN (?)", 
		userID, userID, []models.TaskStatus{
			models.TaskStatusCompleted,
			models.TaskStatusCancelled,
		}).Count(&taskCount)
	
	if taskCount > 0 {
		return errors.New("用户还有未完成的任务，无法删除")
	}
	
	return db.Delete(&user).Error
}

// GetOrganizationStructure 获取组织架构
func (s *OrganizationService) GetOrganizationStructure() ([]models.Organization, error) {
	db := database.GetDB()
	
	var organizations []models.Organization
	if err := db.Preload("Departments.Users").Find(&organizations).Error; err != nil {
		return nil, err
	}
	
	return organizations, nil
}

// GetUsersByRole 根据角色获取用户
func (s *OrganizationService) GetUsersByRole(role models.Role, deptID *uint) ([]models.User, error) {
	db := database.GetDB()
	
	query := db.Preload("Department").Where("role = ? AND is_active = ?", role, true)
	
	if deptID != nil {
		query = query.Where("department_id = ?", *deptID)
	}
	
	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	
	return users, nil
}

// TransferUser 转移用户到其他部门
func (s *OrganizationService) TransferUser(userID uint, newDeptID uint) error {
	db := database.GetDB()
	
	// 检查用户是否存在
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	
	// 检查目标部门是否存在
	var department models.Department
	if err := db.First(&department, newDeptID).Error; err != nil {
		return errors.New("目标部门不存在")
	}
	
	// 更新用户部门
	user.DepartmentID = &newDeptID
	return db.Save(&user).Error
}

// GetDepartmentStats 获取部门统计信息
func (s *OrganizationService) GetDepartmentStats(deptID uint) (map[string]interface{}, error) {
	db := database.GetDB()
	
	// 检查部门是否存在
	var department models.Department
	if err := db.First(&department, deptID).Error; err != nil {
		return nil, errors.New("部门不存在")
	}
	
	stats := make(map[string]interface{})
	
	// 用户统计
	var totalUsers int64
	db.Model(&models.User{}).Where("department_id = ?", deptID).Count(&totalUsers)
	stats["total_users"] = totalUsers
	
	var activeUsers int64
	db.Model(&models.User{}).Where("department_id = ? AND is_active = ?", deptID, true).Count(&activeUsers)
	stats["active_users"] = activeUsers
	
	// 角色统计
	var roleStats []struct {
		Role  models.Role `json:"role"`
		Count int64       `json:"count"`
	}
	db.Model(&models.User{}).
		Select("role, count(*) as count").
		Where("department_id = ? AND is_active = ?", deptID, true).
		Group("role").
		Scan(&roleStats)
	stats["role_stats"] = roleStats
	
	// 任务统计（该部门用户创建或负责的任务）
	var taskStats []struct {
		Status models.TaskStatus `json:"status"`
		Count  int64             `json:"count"`
	}
	db.Table("tasks").
		Select("status, count(*) as count").
		Where("creator_id IN (SELECT id FROM users WHERE department_id = ?) OR assignee_id IN (SELECT id FROM users WHERE department_id = ?)", 
			deptID, deptID).
		Group("status").
		Scan(&taskStats)
	stats["task_stats"] = taskStats
	
	return stats, nil
}