package usecase

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
    "time"
	"fmt"
	"log"
)

type InvoiceUsecase interface {
    CreateInvoice(userID uint, paymentID uint, amount float64, paymentType string, discountedAmount float64) (*models.Invoice, error)
}

type invoiceUsecaseImpl struct {
    invoiceRepo repository.InvoiceRepository
}

func NewInvoiceUsecase(invoiceRepo repository.InvoiceRepository) InvoiceUsecase {
    return &invoiceUsecaseImpl{
        invoiceRepo: invoiceRepo,
    }
}

func (u *invoiceUsecaseImpl) CreateInvoice(userID uint, paymentID uint, amount float64, paymentType string, discountedAmount float64) (*models.Invoice, error) {
    if amount <= 0 {
        return nil, fmt.Errorf("invalid amount: must be greater than 0")
    }
    if paymentType == "" {
        return nil, fmt.Errorf("payment type is required")
    }

    invoice := &models.Invoice{
        UserID:         userID,
        PaymentID:      paymentID, 
        OriginalAmount: amount,
        DiscountAmount: discountedAmount,
        Amount:         amount - discountedAmount,
        PaymentType:    paymentType,
        Status:         "Paid", 
        InvoiceDate:    time.Now(),
    }

    createdInvoice, err := u.invoiceRepo.Create(invoice)
    if err != nil {
        log.Printf("Failed to create invoice: %v, Invoice: %+v", err, invoice)
        return nil, err
    }

    return createdInvoice, nil
}



