package services

import (
	"errors"
	"time"

	"task-management-platform/internal/database"
	"task-management-platform/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret []byte
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	FullName     string `json:"full_name" binding:"required"`
	DepartmentID *uint  `json:"department_id"`
	Role         models.Role `json:"role"`
}

type AuthResponse struct {
	Token     string      `json:"token"`
	User      *models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
}

type Claims struct {
	UserID uint        `json:"user_id"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	db := database.GetDB()
	
	var user models.User
	if err := db.Where("username = ? AND is_active = ?", req.Username, true).
		Preload("Department").
		First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}
	
	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	db.Save(&user)
	
	// 生成 JWT token
	token, expiresAt, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	
	return &AuthResponse{
		Token:     token,
		User:      &user,
		ExpiresAt: expiresAt,
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
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
	
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	// 创建新用户
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		Password:     string(hashedPassword),
		FullName:     req.FullName,
		DepartmentID: req.DepartmentID,
		Role:         req.Role,
		IsActive:     true,
	}
	
	// 如果没有指定角色，默认为员工
	if user.Role == "" {
		user.Role = models.RoleEmployee
	}
	
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	
	// 预加载部门信息
	db.Preload("Department").First(user, user.ID)
	
	// 生成 JWT token
	token, expiresAt, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	
	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}, nil
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(tokenString string) (*AuthResponse, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	db := database.GetDB()
	var user models.User
	if err := db.Preload("Department").First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	
	if !user.IsActive {
		return nil, errors.New("用户已被禁用")
	}
	
	// 生成新的token
	token, expiresAt, err := s.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	
	return &AuthResponse{
		Token:     token,
		User:      &user,
		ExpiresAt: expiresAt,
	}, nil
}

// GetUserByToken 通过token获取用户信息
func (s *AuthService) GetUserByToken(tokenString string) (*models.User, error) {
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}
	
	db := database.GetDB()
	var user models.User
	if err := db.Preload("Department").First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	
	if !user.IsActive {
		return nil, errors.New("用户已被禁用")
	}
	
	return &user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	db := database.GetDB()
	
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}
	
	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}
	
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	// 更新密码
	user.Password = string(hashedPassword)
	return db.Save(&user).Error
}

// generateToken 生成JWT token
func (s *AuthService) generateToken(userID uint, role models.Role) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour) // 24小时过期
	
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "task-management-platform",
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}
	
	return tokenString, expiresAt, nil
}

// parseToken 解析JWT token
func (s *AuthService) parseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}