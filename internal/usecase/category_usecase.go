package usecase

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
)

type CategoryUsecase struct {
    CategoryRepo repository.CategoryRepository
}


func NewCategoryUsecase(repo repository.CategoryRepository) *CategoryUsecase {
    return &CategoryUsecase{CategoryRepo: repo}
}

func (u *CategoryUsecase) AddCategory(name string) error {
    category := models.Category{
        CategoryName: name,
    }
    return u.CategoryRepo.AddCategory(category)
}

func (u *CategoryUsecase) UpdateCategory(id int, name string) error {
    category, err := u.CategoryRepo.GetCategoryByID(id)
    if err != nil {
        return err
    }
    category.CategoryName = name
    return u.CategoryRepo.UpdateCategory(category)
}

func (u *CategoryUsecase) DeleteCategory(id int) error {
    return u.CategoryRepo.DeleteCategory(id)
}

func (u *CategoryUsecase) GetAllCategories() ([]models.Category, error) {
    return u.CategoryRepo.GetAllCategories()
}
