package usecase

import (
    "github.com/Prototype-1/xtrace/internal/repository"
    "github.com/Prototype-1/xtrace/internal/models"
)

type BookingUsecase interface {
    CreateBooking(userID uint, routeID uint, serviceType string, bookingAmount float64, cardType string) error
    GetBookingByID(bookingID uint) (*models.Booking, error) 
    GetBookingByPaymentID(paymentID string) (*models.Booking, error)
    IsPaymentMadeForBooking(bookingID uint) (bool, error) 
}

type bookingUsecase struct {
    bookingRepo repository.BookingRepository
}

func NewBookingUsecase(bookingRepo repository.BookingRepository) BookingUsecase {
    return &bookingUsecase{
        bookingRepo: bookingRepo,
    }
}

func (u *bookingUsecase) CreateBooking(userID uint, routeID uint, serviceType string, bookingAmount float64, cardType string) error {
    return u.bookingRepo.CreateBooking(userID, routeID, serviceType, bookingAmount, cardType)
}

func (u *bookingUsecase) GetBookingByID(bookingID uint) (*models.Booking, error) {
	return u.bookingRepo.GetBookingByID(bookingID) 
}

func (u *bookingUsecase) GetBookingByPaymentID(paymentID string) (*models.Booking, error) {
    return u.bookingRepo.GetBookingByPaymentID(paymentID)
}

func (u *bookingUsecase) IsPaymentMadeForBooking(bookingID uint) (bool, error) {
    booking, err := u.bookingRepo.GetBookingByID(bookingID)
    if err != nil {
        return false, err
    }
    return booking.PaymentID != nil, nil 
}