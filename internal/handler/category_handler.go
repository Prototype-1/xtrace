package handler

import (
    "net/http"
    "strconv"
    "github.com/Prototype-1/xtrace/internal/usecase"
    "github.com/gin-gonic/gin"
)

type CategoryHandler struct {
    CategoryUsecase *usecase.CategoryUsecase
}

func (h *CategoryHandler) AddCategory(c *gin.Context) {
    var req struct {
        CategoryName string `json:"category_name" binding:"required"`
        ImageURL     string `json:"image_url" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := h.CategoryUsecase.AddCategory(req.CategoryName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Category created successfully"})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
        return
    }

    var req struct {
        CategoryName string `json:"category_name" binding:"required"`
        ImageURL     string `json:"image_url" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.CategoryUsecase.UpdateCategory(id, req.CategoryName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
        return
    }

    err = h.CategoryUsecase.DeleteCategory(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (h *CategoryHandler) GetAllCategoriesAdmin(c *gin.Context) {
    categories, err := h.CategoryUsecase.GetAllCategories()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetAllCategoriesUser(c *gin.Context) {
    categories, err := h.CategoryUsecase.GetAllCategories()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    type UserCategoryResponse struct {
        CategoryID   int    `json:"category_id"`
        CategoryName string `json:"category_name"`
        ImageURL     string `json:"image_url"`
    }

    userCategories := make([]UserCategoryResponse, len(categories))
    for i, category := range categories {
        userCategories[i] = UserCategoryResponse{
            CategoryID:   category.CategoryID,
            CategoryName: category.CategoryName,
            ImageURL:     category.ImageURL,
        }
    }

    c.JSON(http.StatusOK, userCategories)
}
