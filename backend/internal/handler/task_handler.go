package handler

import (
	"errors"
	"net/http"
	"strconv"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{TaskService: taskService}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input task tidak valid",
			"error":   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.TaskService.Create(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPriority) ||
			errors.Is(err, service.ErrInvalidStatus) ||
			errors.Is(err, service.ErrInvalidDueDate) ||
			errors.Is(err, service.ErrCategoryNotOwned) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat task",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Task berhasil dibuat",
		"data":    res,
	})
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	categoryID, _ := strconv.ParseInt(c.Query("category_id"), 10, 64)

	filter := dto.TaskFilterQuery{
		Status:     c.Query("status"),
		Priority:   c.Query("priority"),
		CategoryID: categoryID,
		Search:     c.Query("search"),
		Page:       page,
		Limit:      limit,
	}

	userID := c.GetInt64("user_id")
	items, meta, err := h.TaskService.GetAll(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil daftar task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar task",
		"data":    items,
		"meta":    meta,
	})
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID task tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.TaskService.GetByID(c.Request.Context(), userID, taskID)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil detail task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail task",
		"data":    res,
	})
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID task tidak valid",
		})
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input task tidak valid",
			"error":   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.TaskService.Update(c.Request.Context(), userID, taskID, req)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrInvalidPriority) ||
			errors.Is(err, service.ErrInvalidStatus) ||
			errors.Is(err, service.ErrInvalidDueDate) ||
			errors.Is(err, service.ErrCategoryNotOwned) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task berhasil diperbarui",
		"data":    res,
	})
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID task tidak valid",
		})
		return
	}

	var req dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input status tidak valid",
			"error":   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.TaskService.UpdateStatus(c.Request.Context(), userID, taskID, req)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrInvalidStatus) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui status task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status task berhasil diperbarui",
		"data":    res,
	})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || taskID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID task tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	if err := h.TaskService.Delete(c.Request.Context(), userID, taskID); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Task berhasil dihapus",
	})
}

func (h *TaskHandler) Summary(c *gin.Context) {
	userID := c.GetInt64("user_id")
	summary, err := h.TaskService.GetSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil ringkasan task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil ringkasan task",
		"data":    summary,
	})
}