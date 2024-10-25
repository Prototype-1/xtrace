package repository

import (
	"errors"
	"time"
	"github.com/Prototype-1/xtrace/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
    AddCategory(category models.Category) error
    UpdateCategory(category models.Category) error
    DeleteCategory(id int) error
    GetCategoryByID(id int) (models.Category, error)
    GetAllCategories() ([]models.Category, error)
}

type CategoryRepositoryImpl struct {
    DB *gorm.DB
}

func (r *CategoryRepositoryImpl) AddCategory(category models.Category) error {
    return r.DB.Create(&category).Error
}

func (r *CategoryRepositoryImpl) UpdateCategory(category models.Category) error {
    return r.DB.Model(&models.Category{}).
        Where("category_id = ?", category.CategoryID).
        Updates(map[string]interface{}{
            "category_name": category.CategoryName,
            "updated_at":    time.Now(),
        }).Error
}

func (r *CategoryRepositoryImpl) DeleteCategory(id int) error {
    var category models.Category
    if err := r.DB.Where("category_id = ?", id).First(&category).Error; err != nil {
        return err
    }
    category.IsDeleted = true
    return r.DB.Save(&category).Error
}

func (r *CategoryRepositoryImpl) GetCategoryByID(id int) (models.Category, error) {
    var category models.Category
	if err := r.DB.Where("category_id = ?", id).First(&category).Error; err != nil {
        return category, err
    }
    if category.IsDeleted {
        return category, errors.New("category is deleted")
    }
    return category, nil
}

func (r *CategoryRepositoryImpl) GetAllCategories() ([]models.Category, error) {
    var categories []models.Category
    err := r.DB.Where("is_deleted = ?", false).Find(&categories).Error
    return categories, err
}
