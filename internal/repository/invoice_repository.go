package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
)

type InvoiceRepository interface {
    Create(invoice *models.Invoice) (*models.Invoice, error)
    GetInvoiceByUserID(userID uint) (*models.Invoice, error)
}

type invoiceRepositoryImpl struct {
    DB *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
    return &invoiceRepositoryImpl{
        DB: db,
    }
}

func (r *invoiceRepositoryImpl) Create(invoice *models.Invoice) (*models.Invoice, error) {
    if err := r.DB.Create(invoice).Error; err != nil {
        return nil, err
    }
    return invoice, nil
}

func (r *invoiceRepositoryImpl) GetInvoiceByUserID(userID uint) (*models.Invoice, error) {
    var invoice models.Invoice
    result := r.DB.Where("user_id = ?", userID).First(&invoice)
    if result.Error != nil {
        return nil, result.Error
    }
    return &invoice, nil
}
