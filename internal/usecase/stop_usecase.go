package usecase

import (
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
)

type StopUsecase struct {
	StopRepo repository.StopRepository
}

func NewStopUsecase(repo repository.StopRepository) *StopUsecase {
	return &StopUsecase{StopRepo: repo}
}

func (u *StopUsecase) AddStop(stop models.Stop) error {
	return u.StopRepo.AddStop(stop)
}

func (u *StopUsecase) UpdateStop(id int, stop models.Stop) error {
	return u.StopRepo.UpdateStop(stop)
}

func (u *StopUsecase) DeleteStop(id int) error {
	return u.StopRepo.DeleteStop(id)
}

func (u *StopUsecase) GetAllStops() ([]models.Stop, error) {
	return u.StopRepo.GetAllStops()
}
