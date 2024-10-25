package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
    "errors"
    "log"
    "fmt"
    "time"
)

type RazorpayPaymentRepository interface {
    CreatePayment(payment *models.RazorpayPayment) error
    GetPaymentByRazorpayID(razorpayID string) (*models.RazorpayPayment, error)

    GetPaymentByOrderID(orderID string) (*models.RazorpayPayment, error)
    UpdatePaymentStatus(orderID, razorpayID string, status string) error
    GetPaymentByID(paymentID uint) (*models.RazorpayPayment, error)
    GetPaymentsByBookingID(bookingID uint) ([]models.RazorpayPayment, error)
    AddTransaction(transaction *models.WalletTransaction) error 
}

type razorpayPaymentRepositoryImpl struct {
    DB *gorm.DB
}

func NewRazorpayPaymentRepository(db *gorm.DB) RazorpayPaymentRepository {
    return &razorpayPaymentRepositoryImpl{DB: db}
}

func (r *razorpayPaymentRepositoryImpl) CreatePayment(payment *models.RazorpayPayment) error {
    if payment.WalletID != nil && *payment.WalletID == 0 {
        payment.WalletID = nil
    }
    if payment.NolCardID != nil && *payment.NolCardID == 0 {
        payment.NolCardID = nil
    }
    if payment.SubscriptionID != nil && *payment.SubscriptionID == 0 {
        payment.SubscriptionID = nil
    }
    if payment.BookingID != nil && *payment.BookingID == 0 {
        payment.BookingID = nil
    }
    return r.DB.Table("payments").Create(payment).Error
}

func (r *razorpayPaymentRepositoryImpl) GetPaymentByRazorpayID(razorpayID string) (*models.RazorpayPayment, error) {
    var payment models.RazorpayPayment
    err := r.DB.Table("payments").Where("razorpay_id = ?", razorpayID).First(&payment).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil 
        }
        return nil, err 
    }
    return &payment, nil
}

func (r *razorpayPaymentRepositoryImpl) GetPaymentByOrderID(orderID string) (*models.RazorpayPayment, error) {
    var payment models.RazorpayPayment
    log.Printf("Fetching payment by OrderID: %s", orderID)
    err := r.DB.Table("payments").Where("order_id = ?", orderID).First(&payment).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            log.Printf("No payment found with OrderID: %s", orderID)
            return nil, nil
        }
        log.Printf("Error fetching payment by OrderID: %v", err)
        return nil, err
    }
    log.Printf("Payment found with OrderID: %s", orderID)
    return &payment, nil
}

func (r *razorpayPaymentRepositoryImpl) UpdatePaymentStatus(orderID, razorpayID, status string) error {
    var payment models.RazorpayPayment
    log.Printf("Fetching payment by order ID: %s", orderID)
    if err := r.DB.Table("payments").Where("order_id = ?", orderID).First(&payment).Error; err != nil {
        log.Printf("Error fetching payment: %v", err)
        return fmt.Errorf("no payment found with order ID: %s", orderID)
    }
    log.Printf("Fetched payment: %+v", payment)
    
    payment.RazorpayID = razorpayID
    payment.Status = status
    payment.UpdatedAt = time.Now()
    
    log.Printf("Updating payment with Razorpay ID: %s and Status: %s", razorpayID, status)
    if err := r.DB.Table("payments").Save(&payment).Error; err != nil {
        log.Printf("Error updating payment: %v", err)
        return fmt.Errorf("failed to update payment: %v", err)
    }
    
    log.Printf("Payment updated successfully")
    return nil
}

func (r *razorpayPaymentRepositoryImpl) GetPaymentByID(paymentID uint) (*models.RazorpayPayment, error) {
    var payment models.RazorpayPayment
    err := r.DB.Table("payments").Where("payment_id = ?", paymentID).First(&payment).Error
    if err != nil {
        return nil, err
    }
    return &payment, nil
}

func (r *razorpayPaymentRepositoryImpl) GetPaymentsByBookingID(bookingID uint) ([]models.RazorpayPayment, error) {
    var payments []models.RazorpayPayment
    err := r.DB.Table("payments").Where("booking_id = ?", bookingID).Find(&payments).Error
    if err != nil {
        return nil, err
    }
    return payments, nil
}

func (r *razorpayPaymentRepositoryImpl) AddTransaction(transaction *models.WalletTransaction) error {
    return r.DB.Create(transaction).Error
}