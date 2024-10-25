package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"os"
    "fmt"
	"strconv"
	"time"
    "gorm.io/gorm"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
	"github.com/razorpay/razorpay-go"
)

type RazorpayPaymentUsecase interface {
    CreatePayment(userID uint, amount float64, currency string, couponCode string, paymentType string, walletID *uint, nolCardID *uint, subscriptionID *uint, bookingID *uint, orderID string) (*models.RazorpayPayment, error)

    CreateRazorpayOrder(amount float64, currency string, userID uint64) (string, error)
    VerifyPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature string) error
    GetPaymentStatus(paymentID uint) (*models.RazorpayPayment, error)
    GetPaymentByOrderID(orderID string) (*models.RazorpayPayment, error)


    UpdatePaymentStatus(orderID, razorpayID, status string) error
    GetPaymentByRazorpayID(razorpayID string) (*models.RazorpayPayment, error)
    GetExistingPayment(userID uint, paymentType string, walletID, nolCardID, subscriptionID, bookingID *uint) (*models.RazorpayPayment, error)

    ApplyCoupon(couponCode string, amount float64) (float64, error)
    fetchCoupon(couponCode string) (*models.Coupon, error)
    GetApplicableCoupons(paymentType string) ([]*models.Coupon, error)
    
    ProcessRefund(paymentID string) error 
}

type razorpayPaymentUsecaseImpl struct {
    razorpayRepo repository.RazorpayPaymentRepository
    client       *razorpay.Client
    keySecret    string 
    couponRepo   repository.CouponRepository
}

func NewRazorpayPaymentUsecase(razorpayRepo repository.RazorpayPaymentRepository, client *razorpay.Client, couponRepo repository.CouponRepository) RazorpayPaymentUsecase {
    return &razorpayPaymentUsecaseImpl{
        razorpayRepo: razorpayRepo, 
        client:       client,
        keySecret:    os.Getenv("RAZORPAY_KEY_SECRET"), 
        couponRepo:   couponRepo,
    }
}

func (u *razorpayPaymentUsecaseImpl) CreatePayment(userID uint, amount float64, currency string, couponCode string, paymentType string, walletID *uint, nolCardID *uint, subscriptionID *uint, bookingID *uint, orderID string) (*models.RazorpayPayment, error) {
    
    if userID == 0 {
        return nil, errors.New("invalid user ID")
    }
    if amount <= 0 {
        return nil, errors.New("amount must be greater than zero")
    }
    if currency == "" {
        return nil, errors.New("currency cannot be empty")
    }
    if orderID == "" {
        return nil, errors.New("order ID cannot be empty")
    }

    if couponCode != "" {
        isValid, err := u.couponRepo.IsCouponValid(couponCode)
        if err != nil {
            return nil, fmt.Errorf("error validating coupon: %v", err)
        }
        if !isValid {
            return nil, errors.New("invalid or expired coupon")
        }

        coupon, err := u.couponRepo.GetCouponByCode(couponCode)
        if err != nil {
            return nil, fmt.Errorf("error fetching coupon details: %v", err)
        }

        discount := (coupon.DiscountAmount / 100) * amount
        amount -= discount 
        if amount < 0 {
            amount = 0
        }

        log.Printf("Discount applied: %.2f, Final amount: %.2f", discount, amount)
    }

    var walletIDPtr, nolCardIDPtr, subscriptionIDPtr, bookingIDPtr *uint

    if walletID != nil {  
        walletIDPtr = walletID
    }
    if nolCardID != nil {
        nolCardIDPtr = nolCardID
    }
    if subscriptionID != nil {
        subscriptionIDPtr = subscriptionID
    }
    if bookingID != nil {
        bookingIDPtr = bookingID
    }

    payment := &models.RazorpayPayment{
        UserID:         userID,
        RazorpayID:     orderID, 
        OrderID:        orderID,
        Amount:         amount,
        Currency:       currency,
        Status:         "created",
        Method:         "razorpay",
        PaymentType:    paymentType,
        CouponCode:     couponCode,
        WalletID:       walletIDPtr,
        NolCardID:      nolCardIDPtr,
        SubscriptionID: subscriptionIDPtr,
        BookingID:      bookingIDPtr,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
        PaymentDate: time.Now(),
    }

    log.Printf("Saving payment to DB: %v", payment)

    err := u.razorpayRepo.CreatePayment(payment)
    if err != nil {
        log.Printf("Error saving payment to database: %v", err)
        return nil, err
    }

    return payment, nil
}


func (u *razorpayPaymentUsecaseImpl) VerifyPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature string) error {
    log.Printf("Verifying payment - Order ID: %s, Payment ID: %s", razorpayOrderID, razorpayPaymentID)

    signatureData := razorpayOrderID + "|" + razorpayPaymentID
    log.Printf("Signature data: %s", signatureData)

    h := hmac.New(sha256.New, []byte(u.keySecret))
    h.Write([]byte(signatureData))
    generatedSignature := hex.EncodeToString(h.Sum(nil))
    log.Printf("Generated signature: %s", generatedSignature)
    log.Printf("Received signature: %s", razorpaySignature)

    if generatedSignature != razorpaySignature {
        return errors.New("invalid Razorpay signature")
    }
    paymentDetails, err := u.client.Payment.Fetch(razorpayPaymentID, nil, nil) 
    if err != nil {
        return errors.New("failed to fetch payment details: " + err.Error())
    }
    log.Printf("Verifying payment - Order ID: %s, Payment ID: %s", razorpayOrderID, razorpayPaymentID)
log.Printf("Payment details from Razorpay: %v", paymentDetails)

    if paymentDetails["status"] != "captured" {
        return errors.New("payment not captured")
    }
    log.Printf("Updating payment status in database for Order ID: %s and Payment ID: %s", razorpayOrderID, razorpayPaymentID)
 
    err = u.razorpayRepo.UpdatePaymentStatus(paymentDetails["order_id"].(string), razorpayPaymentID, "verified") 
if err != nil {
    log.Printf("Error updating payment status for order ID %s: %v", paymentDetails["order_id"].(string), err)
    return err
    }

    return nil
}

// GetPaymentStatus retrieves the status of a payment by its internal ID
func (u *razorpayPaymentUsecaseImpl) GetPaymentStatus(paymentID uint) (*models.RazorpayPayment, error) {
    payment, err := u.razorpayRepo.GetPaymentByID(paymentID)
    if err != nil {
        log.Printf("Error retrieving payment status for payment ID %d: %v", paymentID, err)
        return nil, err
    }
    return payment, nil
}

func (u *razorpayPaymentUsecaseImpl) UpdatePaymentStatus(orderID, razorpayID, status string) error {
    return u.razorpayRepo.UpdatePaymentStatus(orderID, razorpayID, status)
}

func (u *razorpayPaymentUsecaseImpl) GetPaymentByOrderID(orderID string) (*models.RazorpayPayment, error) {
    payment, err := u.razorpayRepo.GetPaymentByOrderID(orderID)
    if err != nil {
        log.Printf("Error retrieving payment by Order ID %s: %v", orderID, err)
        return nil, err
    }
    return payment, nil
}

func (u *razorpayPaymentUsecaseImpl) GetPaymentByRazorpayID(razorpayID string) (*models.RazorpayPayment, error) {
    payment, err := u.razorpayRepo.GetPaymentByRazorpayID(razorpayID)
    if err != nil {
        log.Printf("Error retrieving payment by Razorpay ID %s: %v", razorpayID, err)
        return nil, err
    }
    return payment, nil
}

func (u *razorpayPaymentUsecaseImpl) GetExistingPayment(userID uint, paymentType string, walletID, nolCardID, subscriptionID, bookingID *uint) (*models.RazorpayPayment, error) {
	var payment *models.RazorpayPayment
	query := config.DB.Where("user_id = ? AND payment_type = ?", userID, paymentType)

	if walletID != nil {
		query = query.Where("wallet_id = ?", *walletID)
	}
	if nolCardID != nil {
		query = query.Where("nol_card_id = ?", *nolCardID)
	}
	if subscriptionID != nil {
		query = query.Where("subscription_id = ?", *subscriptionID)
	}
	if bookingID != nil {
		query = query.Where("booking_id = ?", *bookingID)
	}

    err := query.Table("payments").First(&payment).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil 
        }
        return nil, err 
    }

    return payment, nil
}

func (u *razorpayPaymentUsecaseImpl) CreateRazorpayOrder(amount float64, currency string, userID uint64) (string, error) {
    if amount <= 0 {
        return "", errors.New("amount must be greater than zero")
    }
    if currency == "" {
        return "", errors.New("currency cannot be empty")
    }

    orderData := map[string]interface{}{
        "amount":   amount * 100, 
        "currency": currency,
        "receipt":  "rcpt_" + strconv.FormatUint(userID, 10), 
    }

    // Create Razorpay order
    razorpayOrder, err := u.client.Order.Create(orderData, nil)
    if err != nil {
        log.Printf("Error creating Razorpay order: %v", err)
        return "", err
    }

    razorpayOrderID := razorpayOrder["id"].(string)
    log.Printf("Razorpay order created successfully. Order ID: %s", razorpayOrderID)

    return razorpayOrderID, nil
}

func (u *razorpayPaymentUsecaseImpl) ApplyCoupon(couponCode string, amount float64) (float64, error) {
    coupon, err := u.fetchCoupon(couponCode)
    if err != nil {
        return 0, err 
    }

    if !isCouponValid(coupon) {
        return 0, errors.New("invalid or expired coupon")
    }
    var discount float64
    if coupon.DiscountType == "fixed" {
        discount = coupon.DiscountAmount
    } else if coupon.DiscountType == "percentage" {
        discount = (amount * coupon.DiscountAmount) / 100
    }

    if discount > amount {
        discount = amount
    }

    return discount, nil
}

func (u *razorpayPaymentUsecaseImpl) fetchCoupon(couponCode string) (*models.Coupon, error) {
    return u.couponRepo.GetCouponByCode(couponCode)
}

func isCouponValid(coupon *models.Coupon) bool {
    currentDate := time.Now()
    return coupon.StartDate.Before(currentDate) && coupon.EndDate.After(currentDate)
}

func (u *razorpayPaymentUsecaseImpl) GetApplicableCoupons(paymentType string) ([]*models.Coupon, error) {
	return u.couponRepo.GetCouponsByPaymentType(paymentType)
}

func (u *razorpayPaymentUsecaseImpl) ProcessRefund(paymentID string) error {
    paymentDetails, err := u.client.Payment.Fetch(paymentID, nil, nil)
    if err != nil {
        return fmt.Errorf("failed to fetch payment details: %v", err)
    }

    if paymentDetails["status"] != "captured" {
        return errors.New("payment is not eligible for refund")
    }

    refundRequest := map[string]interface{}{
        "payment_id": paymentID,
    }

    _, err = u.client.Refund.Create(refundRequest, nil)
    if err != nil {
        return fmt.Errorf("failed to create refund: %v", err)
    }

    return nil
}
