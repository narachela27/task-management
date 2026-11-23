package handler

import (
	"errors"
	"net/http"
	"strconv"

	"task-management-backend/internal/dto"
	"task-management-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	CategoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: categoryService}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input kategori tidak valid",
			"error":   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.CategoryService.Create(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal membuat kategori",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kategori berhasil dibuat",
		"data":    res,
	})
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	userID := c.GetInt64("user_id")
	res, err := h.CategoryService.GetAll(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil daftar kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar kategori",
		"data":    res,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || categoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID kategori tidak valid",
		})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input kategori tidak valid",
			"error":   err.Error(),
		})
		return
	}

	userID := c.GetInt64("user_id")
	res, err := h.CategoryService.Update(c.Request.Context(), userID, categoryID, req)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if errors.Is(err, service.ErrCategoryAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kategori berhasil diperbarui",
		"data":    res,
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || categoryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID kategori tidak valid",
		})
		return
	}

	userID := c.GetInt64("user_id")
	if err := h.CategoryService.Delete(c.Request.Context(), userID, categoryID); err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kategori berhasil dihapus",
	})
}