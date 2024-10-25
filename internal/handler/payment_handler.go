package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

type RazorpayHandler struct {
	RazorpayPaymentUsecase usecase.RazorpayPaymentUsecase
	BookingUsecase         usecase.BookingUsecase 
	SubscriptionUsecase    usecase.SubscriptionUsecase
    WalletUsecase          usecase.WalletUsecase
    razorpayClient         *razorpay.Client 
	NolCardTopupUsecase         usecase.NolCardTopupUsecase
	InvoiceUsecase         usecase.InvoiceUsecase 
}

func NewRazorpayHandler(walletUsecase usecase.WalletUsecase, razorpayPaymentUsecase usecase.RazorpayPaymentUsecase, bookingUsecase usecase.BookingUsecase, subscriptionUsecase usecase.SubscriptionUsecase, razorpayClient *razorpay.Client, NolCardTopupUsecase         usecase.NolCardTopupUsecase, invoiceUsecase usecase.InvoiceUsecase) *RazorpayHandler {
	return &RazorpayHandler{
		RazorpayPaymentUsecase: razorpayPaymentUsecase,
		BookingUsecase:         bookingUsecase, 
        WalletUsecase:          walletUsecase,
		SubscriptionUsecase:    subscriptionUsecase,
        razorpayClient:         razorpayClient,
		NolCardTopupUsecase: NolCardTopupUsecase,
		InvoiceUsecase:         invoiceUsecase,
	}
}

func (h *RazorpayHandler) CreatePayment(c *gin.Context) {
	log.Println("CreatePayment endpoint hit")

var originalAmount float64 
var finalAmount float64

	userIDParam := c.Param("userID")
	userID, err := strconv.ParseUint(userIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var input struct {		
		Amount         float64 `json:"amount" binding:"required,numeric"`
		Currency       string  `json:"currency" binding:"required"`
		PaymentType    string  `json:"payment_type" binding:"required"`
		CouponCode     string  `json:"coupon_code"`
		WalletID       *uint   `json:"wallet_id,omitempty"`
		NolCardID      *uint   `json:"nol_card_id,omitempty"`
		SubscriptionID *uint   `json:"subscription_id,omitempty"`
		BookingID      *uint   `json:"booking_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Binding JSON failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Input received: %+v", input)

	existingPayment, err := h.RazorpayPaymentUsecase.GetExistingPayment(uint(userID), input.PaymentType, input.WalletID, input.NolCardID, input.SubscriptionID, input.BookingID)
	if err != nil {
		log.Printf("Error checking for existing payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing payment: " + err.Error()})
		return
	}
	if existingPayment != nil {
		if input.PaymentType != "wallet_topup" {
			c.JSON(http.StatusConflict, gin.H{"error": "A payment already exists for this Payment Type and related ID"})
			return
		}
	}

	var walletIDPtr, nolCardIDPtr, subscriptionIDPtr, bookingIDPtr *uint
	if input.WalletID != nil {
		walletIDPtr = new(uint)
		*walletIDPtr = *input.WalletID
	}
	if input.NolCardID != nil {
		nolCardIDPtr = new(uint)
		*nolCardIDPtr = *input.NolCardID
	}
	if input.SubscriptionID != nil {
		subscriptionIDPtr = new(uint)
		*subscriptionIDPtr = *input.SubscriptionID
	}
	if input.BookingID != nil {
		bookingIDPtr = new(uint)
		*bookingIDPtr = *input.BookingID
	}

	switch input.PaymentType {
	case "wallet_topup":
		if walletIDPtr == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet ID is required for Wallet Topup"})
			return
		}
		originalAmount = input.Amount

	case "nol_card_topup":
		if nolCardIDPtr == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nol Card ID is required for Nol Card Topup"})
			return
		}
		originalAmount = input.Amount

	case "subscription":
		if subscriptionIDPtr == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Subscription ID is required for Subscription payment"})
            return
        }
        subscription, err := h.SubscriptionUsecase.GetSubscriptionByID(*subscriptionIDPtr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Subscription ID"})
            return
        }
        originalAmount = subscription.Price

	case "booking":
		if bookingIDPtr == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID is required for Booking payment"})
			return
		}
		originalAmount = input.Amount
		
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Payment Type"})
		return
	}

	finalAmount = originalAmount
	var discountAmount float64

if input.CouponCode != "" {
    discountAmount, err = h.RazorpayPaymentUsecase.ApplyCoupon(input.CouponCode, finalAmount)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    finalAmount -= discountAmount
}

if finalAmount <= 0 {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Amount must be greater than zero after applying the coupon"})
    return
}

orderID, err := h.RazorpayPaymentUsecase.CreateRazorpayOrder(finalAmount, input.Currency, userID)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating Razorpay order: " + err.Error()})
    return
}


	log.Printf("Razorpay order created successfully. Order ID: %s", orderID)

	payment, err := h.RazorpayPaymentUsecase.CreatePayment(
		uint(userID),
		finalAmount,
		input.Currency,
		input.CouponCode,
		input.PaymentType,
		walletIDPtr,       
		nolCardIDPtr,     
		subscriptionIDPtr,  
		bookingIDPtr,     
		orderID,          
	)

log.Printf("Subscription ID: %v", subscriptionIDPtr)
log.Printf("Creating payment for user %d with amount %.2f", userID, input.Amount)
log.Printf("Original amount: %.2f, Discount: %.2f, Final amount: %.2f", input.Amount, discountAmount, finalAmount)

	if err != nil {
		log.Printf("Failed to create payment record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment: " + err.Error()})
		return
	}

	log.Printf("Payment record created successfully. Payment ID: %d, Razorpay Order ID: %s", payment.PaymentID, orderID)

	invoice, err := h.InvoiceUsecase.CreateInvoice(
		uint(userID),        
		payment.PaymentID,  
		originalAmount,      
		input.PaymentType,   
		discountAmount,      
	)
	if err != nil {
		log.Printf("Failed to create invoice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice: " + err.Error()})
		return
	}
	
	// Log or process the created invoice as needed
	log.Printf("Invoice created successfully. Invoice ID: %d", invoice.InvoiceID)
	

	c.JSON(http.StatusOK, gin.H{
		"payment":     payment,
		"order_id":    orderID,
		"razorpay_id": payment.RazorpayID,
		"original_amount":  input.Amount,
        "discounted_amount": finalAmount,
	})
}

func (h *RazorpayHandler) GetAmountByPaymentType(c *gin.Context) {
    paymentType := c.Param("type")
    id := c.Param("id")

    var amount float64

    switch paymentType {
    case "booking":
        bookingID, _ := strconv.ParseUint(id, 10, 64)
		log.Printf("Fetching booking with ID: %d", bookingID) 
        booking, err := h.BookingUsecase.GetBookingByID(uint(bookingID))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
            return
        }
        amount = booking.BookingAmount
    case "wallet_topup":
     //Not needed now
    case "nol_card_topup":
    // Not needed now
    case "subscription":
		subscriptionID, _ := strconv.ParseUint(id, 10, 64)
        log.Printf("Fetching subscription with ID: %d", subscriptionID)
        subscription, err := h.SubscriptionUsecase.GetSubscriptionByID(uint(subscriptionID)) 
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
            return
        }
        amount = subscription.Price 
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment type"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"amount": amount})
}

func (h *RazorpayHandler) VerifyPayment(c *gin.Context) {
    log.Println("Payment verification endpoint hit")
    var input struct {
        OrderID           string `json:"order_id" binding:"required"`
        PaymentID         string `json:"payment_id" binding:"required"`
        RazorpaySignature string `json:"razorpay_signature" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "verified": false,
            "error": "Invalid input data",
        })
        return
    }

    log.Printf("Verifying payment - Order ID: %s, Payment ID: %s", input.OrderID, input.PaymentID)

    err := h.RazorpayPaymentUsecase.VerifyPayment(input.OrderID, input.PaymentID, input.RazorpaySignature)
    if err != nil {
        log.Printf("Payment verification failed: %v", err)
        refundErr := h.RazorpayPaymentUsecase.ProcessRefund(input.PaymentID)
        if refundErr != nil {
            log.Printf("Refund process failed: %v", refundErr)
            c.JSON(http.StatusInternalServerError, gin.H{
                "verified": false,
                "error": "Payment verification and refund failed",
            })
            return
        }
        log.Println("Refund processed due to payment verification failure")
        c.JSON(http.StatusInternalServerError, gin.H{
            "verified": false,
            "error": "Payment verification failed and refund processed",
        })
        return
    }
    log.Println("Payment verified successfully")

    payment, err := h.RazorpayPaymentUsecase.GetPaymentByOrderID(input.OrderID)
    if err != nil {
        log.Printf("Error fetching payment details: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "verified": false,
            "error": "Failed to fetch payment details",
        })
        return
    }

    var processingErr error
    switch payment.PaymentType {
    case "nol_card_topup":
        processingErr = h.handleNOLCardTopup(payment)
    case "wallet_topup":
        processingErr = h.handleWalletTopup(payment)
    case "subscription", "booking":
        log.Printf("Payment verified successfully for %s", payment.PaymentType)
        c.JSON(http.StatusOK, gin.H{
            "verified": true,
            "message": "Payment verified successfully",
            "payment_type": payment.PaymentType,
        })
        return
    default:
        log.Printf("Unknown payment type: %s", payment.PaymentType)
        c.JSON(http.StatusBadRequest, gin.H{
            "verified": false,
            "error": "Invalid payment type",
        })
        return
    }

    if processingErr != nil {
        log.Printf("Error handling %s: %v", payment.PaymentType, processingErr)
        c.JSON(http.StatusInternalServerError, gin.H{
            "verified": false,
            "error": fmt.Sprintf("Failed to process %s", payment.PaymentType),
        })
        return
    }

    log.Println("Payment processed successfully")
    c.JSON(http.StatusOK, gin.H{
        "verified": true,
        "message": "Payment verified and processed successfully",
        "payment_type": payment.PaymentType,
    })
}

func (h *RazorpayHandler) handleNOLCardTopup( payment *models.RazorpayPayment) error {
    if payment.NolCardID == nil {
        return fmt.Errorf("nol_card_id is nil")
    }

    nolCardID := int(*payment.NolCardID)

    nolCard, err := h.NolCardTopupUsecase.GetNolCardByID(nolCardID)
    if err != nil {
        return fmt.Errorf("failed to retrieve NOL card: %w", err)
    }

    nolCardTopup := models.NolCardTopup{
        NolCardID: nolCard.NolCardID,
        Amount:    payment.Amount,
        TopupDate: time.Now(),
    }

    err = h.NolCardTopupUsecase.AddTopupAndUpdateBalance(nolCardTopup)
    if err != nil {
		log.Printf("Error processing NOL card topup: %v", err)
        return fmt.Errorf("failed to process NOL card top-up: %w", err)
    }

    return nil
}


func (h *RazorpayHandler) handleWalletTopup(payment *models.RazorpayPayment) error {
    wallet, err := h.WalletUsecase.GetWalletByUserID(payment.UserID)
    if err != nil {
        return fmt.Errorf("failed to retrieve wallet: %w", err)
    }

    err = h.WalletUsecase.TopUpWallet(&wallet.WalletID, nil, payment.Amount, "Top-up via successful payment", "top-up")
    if err != nil {
        return fmt.Errorf("failed to update wallet balance: %w", err)
    }

    transaction := &models.WalletTransaction{
        WalletID:        wallet.WalletID,
        AdminID:         nil,
        Amount:          payment.Amount,
        TransactionType: "top-up",
        Description:     "Wallet topped up via Razorpay payment",
    }
   
    err = h.WalletUsecase.RecordWalletTransaction(transaction)
    if err != nil {
        return fmt.Errorf("failed to record wallet transaction: %w", err)
    }

    return nil
}

func (h *RazorpayHandler) GetPaymentStatus(c *gin.Context) {
	paymentIDParam := c.Param("payment_id")
	paymentID, err := strconv.Atoi(paymentIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.RazorpayPaymentUsecase.GetPaymentStatus(uint(paymentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payment status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment": payment})
}

func (h *RazorpayHandler) FetchApplicableCoupons(c *gin.Context) {
	paymentType := c.Param("paymentType")

	coupons, err := h.RazorpayPaymentUsecase.GetApplicableCoupons(paymentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupons": coupons})
}

func (h *RazorpayHandler) ApplyCoupon(c *gin.Context) {
	var req struct {
		CouponCode  string  `json:"coupon_code"`
		Amount      float64 `json:"amount"`
		PaymentType string  `json:"payment_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	discount, err := h.RazorpayPaymentUsecase.ApplyCoupon(req.CouponCode, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	finalAmount := req.Amount - discount
	if finalAmount < 0 {
		finalAmount = 0
	}

	log.Printf("Coupon applied: %s, Discount: %.2f, Final Amount: %.2f", req.CouponCode, discount, finalAmount)


	c.JSON(http.StatusOK, gin.H{
		"original_amount": req.Amount,
		"discount_amount": discount,
		"final_amount":    finalAmount,
		"coupon_code":     req.CouponCode,
	})
	
}
