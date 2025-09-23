package service

import (
	"errors"
	"time"

	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/config"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.JWTConfig
}

func NewAuthService(db *gorm.DB, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		db:     db,
		config: cfg,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // 支持用户名、手机号、邮箱登录
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken           string      `json:"access_token"`
	RefreshToken          string      `json:"refresh_token"`
	User                  *model.User `json:"user"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username   string `json:"username" binding:"required"`
	Phone      string `json:"phone" binding:"required,len=11"` // 手机号必填，11位
	Email      string `json:"email"`                           // 邮箱非必填
	Password   string `json:"password" binding:"required,min=6"`
	Nickname   string `json:"nickname"`
	Department string `json:"department"`
	Company    string `json:"company"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	var user model.User

	// 查找用户（支持用户名、手机号、邮箱登录）
	if err := s.db.Preload("Roles").Where("username = ? OR phone = ? OR email = ?", req.Login, req.Login, req.Login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT token
	accessToken, refreshToken, accessTokenExpiresAt, refreshTokenExpiresAt, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	return &LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		User:                  &user,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*model.User, error) {
	// 检查用户名或手机号是否已存在
	var existingUser model.User
	if err := s.db.Where("username = ? OR phone = ?", req.Username, req.Phone).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名或手机号已存在")
	}

	// 如果提供了邮箱，检查邮箱是否已存在
	if req.Email != "" {
		if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username:   req.Username,
		Phone:      req.Phone,
		Email:      req.Email,
		Password:   string(hashedPassword),
		Nickname:   req.Nickname,
		Department: req.Department,
		Company:    req.Company,
		Status:     1,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// generateToken 生成JWT token
func (s *AuthService) generateToken(user *model.User) (string, string, time.Time, time.Time, error) {
	accessToken, accessTokenExpiresAt, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}
	refreshToken, refreshTokenExpiresAt, err := s.generateRefreshToken(user)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return accessToken, refreshToken, accessTokenExpiresAt, refreshTokenExpiresAt, nil
}

func (s *AuthService) generateAccessToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.config.AccessTokenExpireTime) * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"type":     "access",
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) generateRefreshToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.config.RefreshTokenExpireTime) * time.Hour)

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     expiresAt.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken 验证JWT token
func (s *AuthService) ValidateToken(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 验证token类型
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			return nil, errors.New("invalid token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		var user model.User
		if err := s.db.Preload("Roles").First(&user, uint(userID)).Error; err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken 验证刷新token
func (s *AuthService) ValidateRefreshToken(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 验证token类型
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "refresh" {
			return nil, errors.New("invalid refresh token type")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		var user model.User
		if err := s.db.Preload("Roles").First(&user, uint(userID)).Error; err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, errors.New("invalid refresh token")
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
	// 验证刷新token
	user, err := s.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 生成新的token对
	accessToken, newRefreshToken, accessTokenExpiresAt, refreshTokenExpiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(user)

	return &LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          newRefreshToken,
		User:                  user,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}
