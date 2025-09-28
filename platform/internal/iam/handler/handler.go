package handler

import (
	"errors"
	"net/http"
	"strconv"

	model "gitee.com/sichuan-shutong-zhihui-data/k-base/internal/common/models"
	"gitee.com/sichuan-shutong-zhihui-data/k-base/internal/iam/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler 处理器结构体
type Handler struct {
	db          *gorm.DB
	authService *service.AuthService
}

// NewHandler 创建新的处理器
func NewHandler(db *gorm.DB, authService *service.AuthService) *Handler {
	return &Handler{
		db:          db,
		authService: authService,
	}
}

// 认证相关处理器
// @Summary 登录
// @Description 登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param login_request body service.LoginRequest true "登录请求"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data":    response,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	// JWT是无状态的，登出只需要客户端删除token
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// @Summary 刷新token
// @Description 使用refresh token获取新的access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token body map[string]string true "刷新token请求"
// @Success 200 {object} service.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "token刷新成功",
		"data":    response,
	})
}

// @Summary 修改密码
// @Description 修改密码
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param change_password_request body service.ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/change-password [patch]
func (h *Handler) ChangePassword(c *gin.Context) {
	// 获取当前用户
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangePassword(userModel.ID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// 用户管理处理器

func (h *Handler) GetUsers(c *gin.Context) {
	var users []model.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := h.db.Preload("Roles").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var total int64
	h.db.Model(&model.User{}).Count(&total)

	// 脱敏
	for i := range users {
		users[i].Password = ""
		users[i].Roles = []model.Role{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"data": gin.H{
			"users": users,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user model.User
	if err := h.db.Preload("Roles").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 脱敏
	user.Password = ""
	user.Roles = []model.Role{}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户信息成功",
		"data":    user,
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建用户成功",
		"data":    user,
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user model.User
	if err := h.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateData struct {
		Nickname   string `json:"nickname"`
		Department string `json:"department"`
		Company    string `json:"company"`
		Status     int    `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户信息
	user.Nickname = updateData.Nickname
	user.Department = updateData.Department
	user.Company = updateData.Company
	user.Status = updateData.Status

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新用户成功",
		"data":    user,
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := h.db.Delete(&model.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除用户成功"})
}

// 角色管理处理器

func (h *Handler) GetRoles(c *gin.Context) {
	var roles []model.Role
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := h.db.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var total int64
	h.db.Model(&model.Role{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取角色列表成功",
		"data": gin.H{
			"roles": roles,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

func (h *Handler) GetRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var role model.Role
	if err := h.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取角色信息成功",
		"data":    role,
	})
}

func (h *Handler) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建角色成功",
		"data":    role,
	})
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新角色成功",
		"data":    role,
	})
}

func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	if err := h.db.Delete(&model.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除角色成功"})
}

// 权限管理处理器

func (h *Handler) GetPermissions(c *gin.Context) {
	var permissions []model.Permission
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := h.db.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var total int64
	h.db.Model(&model.Permission{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取权限列表成功",
		"data": gin.H{
			"permissions": permissions,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

func (h *Handler) GetPermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	var permission model.Permission
	if err := h.db.First(&permission, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "权限不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取权限信息成功",
		"data":    permission,
	})
}

// @Summary 检查用户权限
// @Description 检查当前用户是否有指定权限
// @Tags Permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param permission_check body map[string]interface{} true "权限检查请求"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/permissions/check [post]
func (h *Handler) CheckPermission(c *gin.Context) {
	// 从上下文获取当前用户
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return
	}

	var req struct {
		SpaceID  uint   `json:"space_id"`
		Resource string `json:"resource" binding:"required"`
		Action   string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户权限
	hasPermission, err := h.authService.CheckPermission(userModel.ID, req.SpaceID, req.Resource, req.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "权限检查完成",
		"data": gin.H{
			"has_permission": hasPermission,
			"user_id":        userModel.ID,
			"space_id":       req.SpaceID,
			"resource":       req.Resource,
			"action":         req.Action,
		},
	})
}

// 空间管理处理器

// @Summary 获取空间列表
// @Description 分页获取空间列表
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/v1/spaces [get]
func (h *Handler) GetSpaces(c *gin.Context) {
	var spaces []model.Space
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	if err := h.db.Preload("Creator").Offset(offset).Limit(pageSize).Find(&spaces).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var total int64
	h.db.Model(&model.Space{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取空间列表成功",
		"data": gin.H{
			"spaces": spaces,
			"pagination": gin.H{
				"page":      page,
				"page_size": pageSize,
				"total":     total,
			},
		},
	})
}

// @Summary 获取空间详情
// @Description 获取指定空间的详细信息
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id} [get]
func (h *Handler) GetSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	var space model.Space
	if err := h.db.Preload("Creator").Preload("Members").First(&space, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "空间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取空间信息成功",
		"data":    space,
	})
}

// @Summary 创建空间
// @Description 创建新的知识空间
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param space body model.Space true "空间信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/spaces [post]
func (h *Handler) CreateSpace(c *gin.Context) {
	// 获取当前用户
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	userModel, ok := user.(*model.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户信息格式错误"})
		return
	}

	var space model.Space
	if err := c.ShouldBindJSON(&space); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置创建者
	space.CreatedBy = userModel.ID

	if err := h.db.Create(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 将创建者添加为空间管理员
	spaceMember := model.SpaceMember{
		SpaceID: space.ID,
		UserID:  userModel.ID,
		Role:    "admin",
	}
	h.db.Create(&spaceMember)

	c.JSON(http.StatusOK, gin.H{
		"message": "创建空间成功",
		"data":    space,
	})
}

// @Summary 更新空间
// @Description 更新空间信息
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Param space body model.Space true "空间信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id} [put]
func (h *Handler) UpdateSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	var space model.Space
	if err := h.db.First(&space, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "空间不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&space); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&space).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新空间成功",
		"data":    space,
	})
}

// @Summary 删除空间
// @Description 删除指定空间
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id} [delete]
func (h *Handler) DeleteSpace(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	if err := h.db.Delete(&model.Space{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除空间成功"})
}

// @Summary 获取空间成员列表
// @Description 获取指定空间的所有成员
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id}/members [get]
func (h *Handler) GetSpaceMembers(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	var members []model.SpaceMember
	if err := h.db.Preload("User").Where("space_id = ?", id).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取空间成员成功",
		"data":    members,
	})
}

// @Summary 添加空间成员
// @Description 将用户添加到指定空间
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Param member body map[string]interface{} true "成员信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/spaces/{id}/members [post]
func (h *Handler) AddSpaceMember(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已经在空间中
	var existingMember model.SpaceMember
	if err := h.db.Where("space_id = ? AND user_id = ?", id, req.UserID).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已在该空间中"})
		return
	}

	// 添加成员
	member := model.SpaceMember{
		SpaceID: uint(id),
		UserID:  req.UserID,
		Role:    req.Role,
	}

	if err := h.db.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "添加空间成员成功",
		"data":    member,
	})
}

// @Summary 移除空间成员
// @Description 从指定空间中移除用户
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Param user_id path int true "用户ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id}/members/{user_id} [delete]
func (h *Handler) RemoveSpaceMember(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := h.db.Where("space_id = ? AND user_id = ?", id, userID).Delete(&model.SpaceMember{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除空间成员成功"})
}

// @Summary 更新空间成员角色
// @Description 更新用户在空间中的角色
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "空间ID"
// @Param user_id path int true "用户ID"
// @Param role body map[string]string true "新角色"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id}/members/{user_id} [put]
func (h *Handler) UpdateSpaceMemberRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var member model.SpaceMember
	if err := h.db.Where("space_id = ? AND user_id = ?", id, userID).First(&member).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不在该空间中"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	member.Role = req.Role
	if err := h.db.Save(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新空间成员角色成功",
		"data":    member,
	})
}

// 角色权限管理处理器

// @Summary 分配角色权限
// @Description 为角色分配权限
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "角色ID"
// @Param permission body map[string]interface{} true "权限信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id}/permissions [post]
func (h *Handler) AssignRolePermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var req struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查角色是否存在
	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查权限是否存在
	var permission model.Permission
	if err := h.db.First(&permission, req.PermissionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "权限不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 分配权限
	rolePermission := model.RolePermission{
		RoleID:       uint(id),
		PermissionID: req.PermissionID,
	}

	if err := h.db.Create(&rolePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分配角色权限成功"})
}

// @Summary 移除角色权限
// @Description 从角色中移除权限
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "角色ID"
// @Param permission_id path int true "权限ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id}/permissions/{permission_id} [delete]
func (h *Handler) RemoveRolePermission(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	permissionID, err := strconv.Atoi(c.Param("permission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	if err := h.db.Where("role_id = ? AND permission_id = ?", id, permissionID).Delete(&model.RolePermission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除角色权限成功"})
}

// @Summary 获取角色权限列表
// @Description 获取角色的所有权限
// @Tags Roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/roles/{id}/permissions [get]
func (h *Handler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var role model.Role
	if err := h.db.Preload("Permissions").First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "获取角色权限成功",
		"role_id":     role.ID,
		"role_name":   role.Name,
		"permissions": role.Permissions,
	})
}

// @Summary 分配用户角色
// @Description 为用户分配角色
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "用户ID"
// @Param role body map[string]interface{} true "角色信息"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id}/roles [post]
func (h *Handler) AssignUserRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否存在
	var user model.User
	if err := h.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查角色是否存在
	var role model.Role
	if err := h.db.First(&role, req.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "角色不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查用户是否已经有该角色
	var existingUserRole model.UserRole
	if err := h.db.Where("user_id = ? AND role_id = ?", id, req.RoleID).First(&existingUserRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已有该角色"})
		return
	}

	// 分配角色
	userRole := model.UserRole{
		UserID: uint(id),
		RoleID: req.RoleID,
	}

	if err := h.db.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分配用户角色成功"})
}

// @Summary 移除用户角色
// @Description 从用户中移除角色
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "用户ID"
// @Param role_id path int true "角色ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id}/roles/{role_id} [delete]
func (h *Handler) RemoveUserRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	if err := h.db.Where("user_id = ? AND role_id = ?", id, roleID).Delete(&model.UserRole{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除用户角色成功"})
}

// @Summary 获取空间成员列表
// @Description 获取空间中的所有成员
// @Tags Spaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path int true "空间ID"
// @Param role_id path int true "角色ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/spaces/{id}/members/{role_id} [get]
func (h *Handler) GetSpaceMembersByRole(c *gin.Context) {
	spaceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的空间ID"})
		return
	}

	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色ID"})
		return
	}

	var members []model.SpaceMember
	if err := h.db.Where("space_id = ? AND role_id = ?", spaceID, roleID).Find(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取空间成员成功", "data": members})
}

func (h *Handler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
		return
	}

	user, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token有效", "data": user})
}
