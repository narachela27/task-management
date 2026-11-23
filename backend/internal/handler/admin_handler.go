package handler

import (
	"errors"
	"net/http"
	"strconv"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	AdminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{AdminService: adminService}
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := dto.AdminUserFilterQuery{
		Role:     c.Query("role"),
		IsActive: c.Query("is_active"),
		Search:   c.Query("search"),
		Page:     page,
		Limit:    limit,
	}

	items, meta, err := h.AdminService.GetAllUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil daftar user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar user",
		"data":    items,
		"meta":    meta,
	})
}

func (h *AdminHandler) GetUserByID(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID user tidak valid",
		})
		return
	}

	user, err := h.AdminService.GetUserByID(c.Request.Context(), targetID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil detail user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail user",
		"data":    user,
	})
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID user tidak valid",
		})
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input role tidak valid",
			"error":   err.Error(),
		})
		return
	}

	currentAdminID := c.GetInt64("user_id")
	user, err := h.AdminService.UpdateUserRole(c.Request.Context(), currentAdminID, targetID, req)
	if err != nil {
		if errors.Is(err, service.ErrCannotModifySelf) || errors.Is(err, service.ErrInvalidRole) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengubah role user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role user berhasil diperbarui",
		"data":    user,
	})
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	targetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID user tidak valid",
		})
		return
	}

	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input status tidak valid",
			"error":   err.Error(),
		})
		return
	}

	currentAdminID := c.GetInt64("user_id")
	user, err := h.AdminService.UpdateUserStatus(c.Request.Context(), currentAdminID, targetID, req)
	if err != nil {
		if errors.Is(err, service.ErrCannotModifySelf) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengubah status user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status user berhasil diperbarui",
		"data":    user,
	})
}

func (h *AdminHandler) GetAllTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	filter := dto.AdminTaskFilterQuery{
		UserID:   userID,
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
		Search:   c.Query("search"),
		Page:     page,
		Limit:    limit,
	}

	items, meta, err := h.AdminService.GetAllTasks(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil seluruh task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil seluruh task",
		"data":    items,
		"meta":    meta,
	})
}