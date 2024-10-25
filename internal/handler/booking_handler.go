package handler

import (
    "net/http"
    "strconv"
    "math/rand"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/internal/usecase"
)

type BookingHandler struct {
    bookingUsecase usecase.BookingUsecase
}

func NewBookingHandler(bookingUsecase usecase.BookingUsecase) *BookingHandler {
    return &BookingHandler{bookingUsecase: bookingUsecase}
}


func (h *BookingHandler) CreateBooking(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var bookingInput struct {
        RouteID     uint   `json:"route_id"`
        ServiceType string `json:"service_type"`
        CardType    string `json:"card_type"`
    }

    if err := c.ShouldBindJSON(&bookingInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    var bookingAmount float64
    if bookingInput.CardType == "Silver" {
        bookingAmount = 50.00
    } else if bookingInput.CardType == "Gold" {
        bookingAmount = 30.00
    } else {
        c.JSON(http.StatusForbidden, gin.H{"message": "Please upgrade to a Silver or Gold NolCard to make a booking."})
        return
    }

    if bookingInput.ServiceType != "Metro" && bookingInput.ServiceType != "Cycle Rental" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service type. Choose 'Metro' or 'Cycle Rental'."})
        return
    }

    if bookingInput.ServiceType == "Cycle Rental" {
        c.JSON(http.StatusOK, gin.H{
            "message": "The service will be available from 01/01/2025.",
        })
        return
    }

    err = h.bookingUsecase.CreateBooking(uint(userID), bookingInput.RouteID, bookingInput.ServiceType, bookingAmount, bookingInput.CardType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
        return
    }

    if bookingInput.ServiceType == "Metro" {
        source := rand.NewSource(time.Now().UnixNano())
        random := rand.New(source)
        seatNumber := random.Intn(60) + 1

        var cabin string
        if seatNumber <= 30 {
            cabin = "the first cabin"
        } else {
            cabin = "the last cabin"
        }

        c.JSON(http.StatusOK, gin.H{
            "message":     "Please find your seat " + strconv.Itoa(seatNumber) + " in " + cabin,
            "Alert": "If payment unsuccessful your booking will be automatically cancelled.",
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Please proceed with the payment.",
        "booking_amount": bookingAmount,
    })
}

func (h *BookingHandler) VerifyPayment(c *gin.Context) {
    paymentID := c.Param("paymentID")
    
    booking, err := h.bookingUsecase.GetBookingByPaymentID(paymentID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Booking confirmed successfully.",
        "booking_amount": booking.BookingAmount,
    })
}
