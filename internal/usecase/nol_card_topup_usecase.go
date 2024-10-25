package usecase

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
    "time"
    "fmt"
)

type NolCardTopupUsecase interface {
    AddTopupAndUpdateBalance(topup models.NolCardTopup) error 
    GetTopupsByCardID(nolCardID int) ([]models.NolCardTopup, error)
    GetTopupByID(topupID int) (models.NolCardTopup, error)
    GetCardTypeByCardID(nolCardID int) (string, error)
    GetNolCardByID(nolCardID int) (models.NolCard, error)
    GetNolCardByUserID(userID int) (models.NolCard, error)
    // UpdateNolCardBalance(topup models.NolCardTopup, nolCard models.NolCard) error
}

type NolCardTopupUsecaseImpl struct {
    NolCardTopupRepo repository.NolCardTopupRepository
}

func NewNolCardTopupUsecase(repo repository.NolCardTopupRepository) NolCardTopupUsecase {
    return &NolCardTopupUsecaseImpl{NolCardTopupRepo: repo}
}

func (u *NolCardTopupUsecaseImpl) AddTopupAndUpdateBalance(topup models.NolCardTopup) error {
    // Set the top-up date
    topup.TopupDate = time.Now()

    // Get the card type and perform validations (as you already have)
    cardType, err := u.GetCardTypeByCardID(topup.NolCardID)
    if err != nil {
        return err
    }

    // Minimum top-up amount check
    minTopup := map[string]float64{"gold": 100, "silver": 50, "ordinary": 20}[cardType]
    if topup.Amount < minTopup {
        return fmt.Errorf("minimum top-up for %s card is %.2f", cardType, minTopup)
    }

    // Retrieve the Nol card
    nolCard, err := u.GetNolCardByID(topup.NolCardID)
    if err != nil {
        return err
    }

    // Add top-up to the database and update balance
    err = u.NolCardTopupRepo.AddTopupAndUpdateBalance(topup, nolCard)
    if err != nil {
        return err
    }

    // Update the balance in the NolCard table after successful top-up record
    nolCard.Balance += topup.Amount // This logic is correctly placed here

    return nil
}

func (u *NolCardTopupUsecaseImpl) GetTopupsByCardID(nolCardID int) ([]models.NolCardTopup, error) {
    return u.NolCardTopupRepo.GetTopupsByCardID(nolCardID)
}

func (u *NolCardTopupUsecaseImpl) GetTopupByID(topupID int) (models.NolCardTopup, error) {
    return u.NolCardTopupRepo.GetTopupByID(topupID)
}

func (u *NolCardTopupUsecaseImpl) GetCardTypeByCardID(nolCardID int) (string, error) {
    return u.NolCardTopupRepo.GetCardTypeByCardID(nolCardID)
}

func (u *NolCardTopupUsecaseImpl) GetNolCardByID(nolCardID int) (models.NolCard, error) {
    return u.NolCardTopupRepo.GetNolCardByID(nolCardID)
}

func (u *NolCardTopupUsecaseImpl) GetNolCardByUserID(userID int) (models.NolCard, error) {
    return u.NolCardTopupRepo.GetNolCardByUserID(userID) 
}

// func (u *NolCardTopupUsecaseImpl) UpdateNolCardBalance(topup models.NolCardTopup, nolCard models.NolCard) error {
//     return u.NolCardTopupRepo.AddTopupAndUpdateBalance(topup, nolCard)
// }