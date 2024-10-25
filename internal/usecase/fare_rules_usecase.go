package usecase

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/domain"
    "github.com/Prototype-1/xtrace/internal/repository"
    "fmt"
)

type FareRuleUsecase interface {
    CreateFareRule(fareRule models.FareRule) error
    UpdateFareRule(fareRule models.FareRule) error
    DeleteFareRule(id int) error
    GetFareRuleByID(id int) (models.FareRule, error)
    GetAllFareRules() ([]models.FareRule, error)
    GetFareRuleByRouteID(routeID int) (models.FareRule, error)
    GetTravelTimeByCoordinates(fromLat, fromLon, toLat, toLon float64) (int, error) 
    CreateStopDuration(stopDuration models.StopDuration) error
    UpdateStopDuration(stopDuration models.StopDuration) error
    DeleteStopDuration(id uint) error
    GetAllStopDurations() ([]models.StopDuration, error)
}

type FareRuleUsecaseImpl struct {
    repo repository.FareRuleRepository
    osrmService *domain.OSRMService
}

func NewFareRuleUsecase(repo repository.FareRuleRepository, osrmService *domain.OSRMService) FareRuleUsecase {
    return &FareRuleUsecaseImpl{
        repo: repo,
        osrmService: osrmService,
    }
}

func (u *FareRuleUsecaseImpl) CreateFareRule(fareRule models.FareRule) error {
    return u.repo.CreateFareRule(fareRule)
}

func (u *FareRuleUsecaseImpl) UpdateFareRule(fareRule models.FareRule) error {
    return u.repo.UpdateFareRule(fareRule)
}

func (u *FareRuleUsecaseImpl) DeleteFareRule(id int) error {
    return u.repo.DeleteFareRule(id)
}

func (u *FareRuleUsecaseImpl) GetFareRuleByID(id int) (models.FareRule, error) {
    return u.repo.GetFareRuleByID(id)
}

func (u *FareRuleUsecaseImpl) GetAllFareRules() ([]models.FareRule, error) {
    return u.repo.GetAllFareRules()
}

func (u *FareRuleUsecaseImpl) GetFareRuleByRouteID(routeID int) (models.FareRule, error) {
    return u.repo.GetFareRuleByRouteID(routeID)
}

func (u *FareRuleUsecaseImpl) GetTravelTimeByCoordinates(fromLat, fromLon, toLat, toLon float64) (int, error) {
    duration, _, err := u.osrmService.GetTravelTime(fromLat, fromLon, toLat, toLon)
    if err != nil {
        // Return -1 for duration to indicate an error condition
        return -1, fmt.Errorf("failed to get travel time: %w", err)
    }

    // Convert the duration from seconds to minutes
    return int(duration / 60), nil
}


func (u *FareRuleUsecaseImpl) CreateStopDuration(stopDuration models.StopDuration) error {
    return u.repo.CreateStopDuration(stopDuration)
}

func (u *FareRuleUsecaseImpl) UpdateStopDuration(stopDuration models.StopDuration) error {
    return u.repo.UpdateStopDuration(stopDuration)
}

func (u *FareRuleUsecaseImpl) DeleteStopDuration(id uint) error {
    return u.repo.DeleteStopDuration(id)
}

func (u *FareRuleUsecaseImpl) GetAllStopDurations() ([]models.StopDuration, error) {
    return u.repo.GetAllStopDurations()
}



