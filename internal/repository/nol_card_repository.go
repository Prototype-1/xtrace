package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
)

type NolCardRepository interface {
    GetNolCardByID(nolCardID int) (models.NolCard, error)
    AddNolCard(nolCard models.NolCard) error
    GetNolCardByNumber(cardNumber string) (*models.NolCard, error)
}

type NolCardRepositoryImpl struct {
    DB *gorm.DB
}

func NewNolCardRepository(db *gorm.DB) NolCardRepository {
    return &NolCardRepositoryImpl{DB: db}
}

func (r *NolCardRepositoryImpl) GetNolCardByID(nolCardID int) (models.NolCard, error) {
    var nolCard models.NolCard
    err := r.DB.First(&nolCard, nolCardID).Error
    return nolCard, err
}

func (r *NolCardRepositoryImpl) AddNolCard(nolCard models.NolCard) error {
    return r.DB.Create(&nolCard).Error
}

func (r *NolCardRepositoryImpl) GetNolCardByNumber(cardNumber string) (*models.NolCard, error) {
    var nolCard models.NolCard
    err := r.DB.Where("card_number = ?", cardNumber).First(&nolCard).Error
    if err != nil {
        return nil, err
    }
    return &nolCard, nil
}

