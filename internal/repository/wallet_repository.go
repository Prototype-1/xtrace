package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
	"fmt"
	"log"
)

type WalletRepository interface {
    CreateWallet(userID uint) (*models.Wallet, error)                
    GetWalletByUserID(userID uint) (*models.Wallet, error)  
    GetWalletByID(walletID uint) (*models.Wallet, error)          
    UpdateWalletBalance(walletID uint, newBalance float64) error       
    DeductWalletBalance(walletID uint, amount float64) error          
}

type WalletTransactionRepository interface {
    CreateTransaction(transaction *models.WalletTransaction) error     
    GetTransactionsByWalletID(walletID uint) ([]models.WalletTransaction, error) 
}

type walletRepositoryImpl struct {
    DB *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
    return &walletRepositoryImpl{DB: db}
}

//walletTransactionRepositoryImpl
type walletTransactionRepositoryImpl struct {
    DB *gorm.DB
}

func NewWalletTransactionRepository(db *gorm.DB) WalletTransactionRepository {
    return &walletTransactionRepositoryImpl{DB: db}
}

func (r *walletRepositoryImpl) CreateWallet(userID uint) (*models.Wallet, error) {
    wallet := &models.Wallet{
        UserID:  userID,
        Balance: 0.0,
    }
    err := r.DB.Create(wallet).Error
    if err != nil {
        return nil, err
    }
    return wallet, nil
}

func (r *walletRepositoryImpl) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	log.Printf("Retrieving wallet for user ID: %d", userID)
    var wallet models.Wallet
    err := r.DB.Where("user_id = ?", userID).First(&wallet).Error
    if err != nil {
		log.Printf("Error retrieving wallet: %v", err)
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepositoryImpl) GetWalletByID(walletID uint) (*models.Wallet, error) {
    var wallet models.Wallet
    err := r.DB.Where("wallet_id = ?", walletID).First(&wallet).Error
    if err != nil {
        log.Printf("Error retrieving wallet: %v", err)
        return nil, err
    }
    return &wallet, nil
}

func (r *walletRepositoryImpl) UpdateWalletBalance(walletID uint, newBalance float64) error {
    return r.DB.Model(&models.Wallet{}).Where("wallet_id = ?", walletID).Update("balance", newBalance).Error
}

func (r *walletRepositoryImpl) DeductWalletBalance(walletID uint, amount float64) error {
    var wallet models.Wallet
    if err := r.DB.Where("wallet_id = ?", walletID).First(&wallet).Error; err != nil {
        return err
    }

    if wallet.Balance < amount {
        return fmt.Errorf("insufficient balance")
    }

    wallet.Balance -= amount
    return r.DB.Save(&wallet).Error
}

func (r *walletTransactionRepositoryImpl) CreateTransaction(transaction *models.WalletTransaction) error {
    return r.DB.Create(transaction).Error
}

func (r *walletTransactionRepositoryImpl) GetTransactionsByWalletID(walletID uint) ([]models.WalletTransaction, error) {
    var transactions []models.WalletTransaction
    err := r.DB.Where("wallet_id = ?", walletID).Order("created_at DESC").Find(&transactions).Error
    if err != nil {
        return nil, err
    }
    return transactions, nil
}