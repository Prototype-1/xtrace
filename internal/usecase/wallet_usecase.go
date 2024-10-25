package usecase

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
	"fmt"
    "log"
)

type WalletUsecase interface {
    CreateWallet(userID uint) (*models.Wallet, error)                    
    GetWalletByUserID(userID uint) (*models.Wallet, error)     
	 // Top up wallet by admin          
    TopUpWallet(walletID, adminID *uint, amount float64, description string, transactionType string) error
    RecordWalletTransaction(transaction *models.WalletTransaction) error
    MakePayment(walletID uint, amount float64, transactionType string) error                 
    GetWalletTransactions(walletID uint) ([]models.WalletTransaction, error) 
    GetWalletByID(walletID uint) (*models.Wallet, error) 
}


type walletUsecaseImpl struct {
    walletRepo           repository.WalletRepository
    walletTransactionRepo repository.WalletTransactionRepository
}

func NewWalletUsecase(walletRepo repository.WalletRepository, walletTransactionRepo repository.WalletTransactionRepository) WalletUsecase {
    return &walletUsecaseImpl{
        walletRepo:           walletRepo,
        walletTransactionRepo: walletTransactionRepo,
    }
}

func (u *walletUsecaseImpl) CreateWallet(userID uint) (*models.Wallet, error) {
    return u.walletRepo.CreateWallet(userID)
}

func (u *walletUsecaseImpl) GetWalletByUserID(userID uint) (*models.Wallet, error) {
	wallet, err := u.walletRepo.GetWalletByUserID(userID)
    if err != nil {
        return nil, err
    }
    if wallet == nil {
        return nil, fmt.Errorf("wallet not found for user")
    }
    return wallet, nil
}

func (u *walletUsecaseImpl) TopUpWallet(walletID *uint, adminID *uint, amount float64, description string, transactionType string) error {
    wallet, err := u.walletRepo.GetWalletByID(*walletID)
    if err != nil {
        return fmt.Errorf("wallet not found: %v", err)
    }
    newBalance := wallet.Balance + amount

    if err := u.walletRepo.UpdateWalletBalance(wallet.WalletID, newBalance); err != nil {
        return fmt.Errorf("failed to update wallet balance: %v", err)
    }

    var adminIDPtr *uint
    if adminID != nil {
        adminIDPtr = adminID
    } else {
        adminIDPtr = nil
    }

    transaction := models.WalletTransaction{
        WalletID:    wallet.WalletID,
        AdminID:     adminIDPtr,
        Amount:      amount,
        TransactionType:  transactionType,
        Description: description,
    }
    return u.walletTransactionRepo.CreateTransaction(&transaction)
}

func (u *walletUsecaseImpl) MakePayment(walletID uint, amount float64, transactionType string) error {
    wallet, err := u.walletRepo.GetWalletByID(walletID)
    if err != nil {
        log.Printf("Error retrieving wallet (ID: %d): %v", walletID, err)
        return fmt.Errorf("wallet not found for ID %d: %v", walletID, err)
    }

    log.Printf("Retrieved wallet: %+v", wallet)

    log.Printf("Wallet Balance: %.2f, Payment Amount: %.2f", wallet.Balance, amount)

    if wallet.Balance < amount {
        return fmt.Errorf("insufficient balance in wallet ID %d", walletID)
    }

    newBalance := wallet.Balance - amount
    if err := u.walletRepo.UpdateWalletBalance(wallet.WalletID, newBalance); err != nil {
        log.Printf("Failed to deduct balance from wallet (ID: %d): %v", wallet.WalletID, err)
        return fmt.Errorf("failed to deduct wallet balance for ID %d: %v", wallet.WalletID, err)
    }

    var adminID *uint = nil 
    transaction := models.WalletTransaction{
        WalletID:       wallet.WalletID,
        AdminID:        adminID, 
        Amount:         amount,
        TransactionType: transactionType, // Use the transactionType parameter
        Description:     "User made a payment",
    }
    if err := u.walletTransactionRepo.CreateTransaction(&transaction); err != nil {
        log.Printf("Error creating wallet transaction for wallet (ID: %d): %v", wallet.WalletID, err)
        return fmt.Errorf("failed to create wallet transaction for ID %d: %v", wallet.WalletID, err)
    }
    log.Printf("Transaction created successfully: %+v", transaction)
    
    return nil
}

func (u *walletUsecaseImpl) GetWalletTransactions(walletID uint) ([]models.WalletTransaction, error) {
    return u.walletTransactionRepo.GetTransactionsByWalletID(walletID)
}

func (u *walletUsecaseImpl) GetWalletByID(walletID uint) (*models.Wallet, error) {
    return u.walletRepo.GetWalletByID(walletID)
}

func (u *walletUsecaseImpl) RecordWalletTransaction(transaction *models.WalletTransaction) error {
    return u.walletTransactionRepo.CreateTransaction(transaction)
}
