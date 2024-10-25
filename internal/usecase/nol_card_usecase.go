package usecase

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
)

type NolCardUsecase interface {
    GetNolCardByID(nolCardID int) (models.NolCard, error)
    AddNolCard(nolCard models.NolCard) error
}

type NolCardUsecaseImpl struct {
    NolCardRepo repository.NolCardRepository
}

func NewNolCardUsecase(repo repository.NolCardRepository) NolCardUsecase {
    return &NolCardUsecaseImpl{NolCardRepo: repo}
}

func (u *NolCardUsecaseImpl) GetNolCardByID(nolCardID int) (models.NolCard, error) {
    return u.NolCardRepo.GetNolCardByID(nolCardID)
}

func (u *NolCardUsecaseImpl) AddNolCard(nolCard models.NolCard) error {
    return u.NolCardRepo.AddNolCard(nolCard)
}