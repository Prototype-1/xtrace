package repository

import (
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/gorm"
    "time"
)

type BookingRepository interface {
    CreateBooking(userID uint, routeID uint, serviceType string, bookingAmount float64, cardType string) error
    GetBookingByID(bookingID uint) (*models.Booking, error)
    GetBookingByPaymentID(paymentID string) (*models.Booking, error)
}

type bookingRepository struct {
    DB *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
    return &bookingRepository{DB: db}
}

func (r *bookingRepository) CreateBooking(userID uint, routeID uint, serviceType string, bookingAmount float64, cardType string) error {
    booking := models.Booking{
        UserID:        userID,
        RouteID:       routeID,
        ServiceType:   serviceType,
        BookingAmount: bookingAmount,
        CardType:      cardType,
        Status:        "Pending Payment", 
        PaymentID:     nil, 
    }

    return r.DB.Create(&booking).Error
}

func (r *bookingRepository) UpdatePaymentStatus(bookingID uint, paymentID uint, status string) error {
    return r.DB.Model(&models.Booking{}).Where("booking_id = ?", bookingID).
        Updates(map[string]interface{}{
            "payment_id": &paymentID, 
            "status":     status,
            "booking_date": time.Now(),
        }).Error
}

func (r *bookingRepository) GetBookingByID(bookingID uint) (*models.Booking, error) {
	var booking models.Booking
    err := r.DB.First(&booking, bookingID).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetBookingByPaymentID(paymentID string) (*models.Booking, error) {
    var booking models.Booking
    err := r.DB.Where("payment_id = ?", paymentID).First(&booking).Error
    if err != nil {
        return nil, err
    }
    return &booking, nil
}