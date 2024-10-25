package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
    "fmt"
    "log"
)

type NolCardTopupRepository interface {
    AddTopupAndUpdateBalance(topup models.NolCardTopup, nolCard models.NolCard) error
    GetTopupsByCardID(nolCardID int) ([]models.NolCardTopup, error)
    GetTopupByID(topupID int) (models.NolCardTopup, error)
    GetCardTypeByCardID(nolCardID int) (string, error)
    GetNolCardByID(nolCardID int) (models.NolCard, error)
    GetNolCardByUserID(userID int) (models.NolCard, error)
}

type NolCardTopupRepositoryImpl struct {
    db *gorm.DB
}

func NewNolCardTopupRepository(db *gorm.DB) NolCardTopupRepository {
    return &NolCardTopupRepositoryImpl{db: db}
}

func (r *NolCardTopupRepositoryImpl) AddTopupAndUpdateBalance(topup models.NolCardTopup, nolCard models.NolCard) error {
    tx := r.db.Begin() // Start the transaction

    // Attempt to create the top-up record
    if err := tx.Create(&topup).Error; err != nil {
        tx.Rollback() // Rollback on error
        return err
    }

    currentBalance := nolCard.Balance 

    newBalance := currentBalance + topup.Amount 

    // Update the balance in the NolCard table
    if err := tx.Model(&models.NolCard{}).Where("nol_card_id = ?", nolCard.NolCardID).Update("balance", newBalance).Error; err != nil {
        tx.Rollback() // Rollback on error
        return err
    }

    // Commit the transaction, but only if no errors occurred
    if err := tx.Commit().Error; err != nil {
        return err
    }
    log.Printf("Current Balance: %f, Topup Amount: %f, New Balance: %f", currentBalance, topup.Amount, newBalance)

    return nil // Successful operation
}

func (r *NolCardTopupRepositoryImpl) GetTopupsByCardID(nolCardID int) ([]models.NolCardTopup, error) {
    var topups []models.NolCardTopup
    err := r.db.Where("nol_card_id = ?", nolCardID).Find(&topups).Error
    return topups, err
}

func (r *NolCardTopupRepositoryImpl) GetTopupByID(topupID int) (models.NolCardTopup, error) {
    var topup models.NolCardTopup
    err := r.db.Where("topup_id = ?", topupID).First(&topup).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return topup, fmt.Errorf("topup with ID %d not found", topupID)
        }
        return topup, err
    }
    return topup, nil
}

func (r *NolCardTopupRepositoryImpl) GetCardTypeByCardID(nolCardID int) (string, error) {
    var card models.NolCard
    err := r.db.Select("card_type").Where("nol_card_id = ?", nolCardID).First(&card).Error
    if err != nil {
        return "", err
    }
    return card.CardType, nil
}

func (r *NolCardTopupRepositoryImpl) GetNolCardByID(nolCardID int) (models.NolCard, error) {
    var nolCard models.NolCard
    err := r.db.Where("nol_card_id = ?", nolCardID).First(&nolCard).Error
    return nolCard, err
}


func (r *NolCardTopupRepositoryImpl) GetNolCardByUserID(userID int) (models.NolCard, error) {
    var nolCard models.NolCard
    err := r.db.Where("user_id = ?", userID).First(&nolCard).Error
    return nolCard, err
}