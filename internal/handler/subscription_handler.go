package handler

import (
	"net/http"
	"strconv"
    "log"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"github.com/gin-gonic/gin"
    "time"
)

type SubscriptionHandler struct {
    SubscriptionUsecase usecase.SubscriptionUsecase
    nolCardRepo         repository.NolCardRepository 
    subscriptionRepo    repository.SubscriptionRepository
    RazorpayPaymentUsecase   usecase.RazorpayPaymentUsecase 
    WalletUsecase     usecase.WalletUsecase
}

func NewSubscriptionHandler(subscriptionUsecase usecase.SubscriptionUsecase, nolCardRepo repository.NolCardRepository, subscriptionRepo repository.SubscriptionRepository, razorpayPaymentUsecase usecase.RazorpayPaymentUsecase, walletUsecase usecase.WalletUsecase) *SubscriptionHandler {
    return &SubscriptionHandler{
        SubscriptionUsecase: subscriptionUsecase,
        nolCardRepo:         nolCardRepo,
        subscriptionRepo:    subscriptionRepo,
        RazorpayPaymentUsecase: razorpayPaymentUsecase,
        WalletUsecase:     walletUsecase, 
    }
}

func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
    var input struct {
        UserID        uint   `json:"user_id" binding:"required"`
        WalletID      uint   `json:"wallet_id"`
        PlanID        uint   `json:"plan_id" binding:"required"`
        ServiceType   string `json:"service_type" binding:"required"`
        DurationDays  int    `json:"duration_days" binding:"required"`
        CardType      string `json:"card_type" binding:"required"`
        NolCardNumber string `json:"nol_card_number" binding:"required"`
        PaymentMethod string `json:"payment_method" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    validCardTypes := map[string]bool{
        "Ordinary": true,
        "Silver":   true,
        "Gold":     true,
    }
    if _, valid := validCardTypes[input.CardType]; !valid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card type"})
        return
    }
    nolCard, err := h.nolCardRepo.GetNolCardByNumber(input.NolCardNumber)
    if err != nil || nolCard == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Nol Card not found"})
        return
    }

    existingSubscription, err := h.subscriptionRepo.GetSubscriptionByUserAndCard(input.UserID, uint(nolCard.NolCardID))
    if err == nil && existingSubscription != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "You already have an active subscription with this NolCard"})
        return
    }
    plan, err := h.SubscriptionUsecase.GetSubscriptionPlan(input.PlanID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
        return
    }

    newSubscription := &models.Subscription{
        UserID:       input.UserID,
        PlanID:       input.PlanID,
        ServiceType:  plan.PlanName,
        DurationDays: plan.DurationDays,
        CardType:     plan.CardType,
        NolCardID:    uint(nolCard.NolCardID),
        Price:        plan.Price,
        StartDate:    time.Now(), 
        EndDate:      time.Now().AddDate(0, 0, plan.DurationDays), 
    }
    
    if input.PaymentMethod == "wallet" {
       // err := h.WalletUsecase.MakePayment(input.WalletID, plan.Price)
        // if err != nil {
        //     c.JSON(http.StatusPaymentRequired, gin.H{"error": "Insufficient wallet balance or payment failed"})
        //     return
        // }
        err := h.subscriptionRepo.CreateSubscription(newSubscription)
        log.Println(err)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "message":             "Subscription created successfully using wallet.",
            "subscription_amount": newSubscription.Price,
            "subscription_id":     newSubscription.SubscriptionID,
        })

    } else if input.PaymentMethod == "razorpay" {

        err := h.subscriptionRepo.CreateSubscription(newSubscription)
        log.Println(err)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message":             "Please proceed with the payment.",
            "subscription_amount": newSubscription.Price,
            "subscription_id":     newSubscription.SubscriptionID,
        })
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
    }
}


func (h *SubscriptionHandler) ConfirmSubscriptionPayment(c *gin.Context) {
    var input struct {
        OrderID       string `json:"order_id" binding:"required"`
         PaymentID       string `json:"payment_id" binding:"required"`
        NolCardNumber string `json:"nol_card_number" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Fetch payments for the given order ID
    paymentID, err := h.SubscriptionUsecase.GetPaymentIDByOrderID(input.OrderID)
    if err != nil || paymentID == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment verification failed"})
        return
    }

    // Verify payment using RazorpayPaymentUsecase
    err = h.RazorpayPaymentUsecase.VerifyPayment(input.OrderID, paymentID, "")
    if err != nil {
        log.Printf("Payment verification failed for payment ID: %v, error: %v", paymentID, err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment verification failed"})
        return
    }

    nolCard, err := h.nolCardRepo.GetNolCardByNumber(input.NolCardNumber)
    if err != nil || nolCard == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Nol Card not found"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Subscription created successfully"})
}

func (h *SubscriptionHandler) ExtendSubscription(c *gin.Context) {
    subscriptionIDStr := c.Param("id")
    subscriptionID, err := strconv.ParseUint(subscriptionIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Please choose a valid subscription ID"})
        return
    }
    var input struct {
        PlanID uint `json:"plan_id"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.SubscriptionUsecase.ExtendSubscription(uint(subscriptionID), input.PlanID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Subscription extended successfully"})
}

// UpdateSubscription updates an existing subscription
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
    subscriptionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Please choose a valid subscription ID"})
        return
    }

    var input struct {
        Price         float64 `json:"price"`
        DurationDays  int     `json:"duration_days"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = h.SubscriptionUsecase.UpdateSubscription(uint(subscriptionID), input.Price, input.DurationDays)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// GetUserSubscriptions retrieves subscriptions for a user
func (h *SubscriptionHandler) GetUserSubscriptions(c *gin.Context) {
    userIDStr := c.Param("id")
    userID, err := strconv.ParseUint(userIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    subscriptions, err := h.SubscriptionUsecase.GetUserSubscriptions(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, subscriptions)
}


// GetAllSubscriptions retrieves all subscriptions
func (h *SubscriptionHandler) GetAllSubscriptions(c *gin.Context) {
    subscriptions, err := h.SubscriptionUsecase.GetAllSubscriptions()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, subscriptions)
}

// GetSubscriptionPlan retrieves a subscription plan by its ID
func (h *SubscriptionHandler) GetSubscriptionPlan(c *gin.Context) {
    planID, err := strconv.ParseUint(c.Param("plan_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
        return
    }

    plan, err := h.SubscriptionUsecase.GetSubscriptionPlan(uint(planID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if plan == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
        return
    }

    c.JSON(http.StatusOK, plan)
}

// CreateSubscriptionPlan creates a new subscription plan
func (h *SubscriptionHandler) CreateSubscriptionPlan(c *gin.Context) {
    var input struct {
        PlanName        string  `json:"plan_name"`
        Price       float64 `json:"price"`
        DurationDays int    `json:"duration_days"`
        CardType      string  `json:"card_type"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    validCardTypes := map[string]bool{
        "Ordinary": true,
        "Silver":   true,
        "Gold":     true,
    }
    if _, valid := validCardTypes[input.CardType]; !valid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card type"})
        return
    }

    plan := &models.SubscriptionPlan{
        PlanName:         input.PlanName,
        Price:        input.Price,
        DurationDays: input.DurationDays,
        CardType:      input.CardType,
    }

    err := h.SubscriptionUsecase.CreateSubscriptionPlan(plan)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Subscription plan created successfully"})
}

// UpdateSubscriptionPlan updates an existing subscription plan
func (h *SubscriptionHandler) UpdateSubscriptionPlan(c *gin.Context) {
    planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
        return
    }

    var input struct {
        PlanName         string  `json:"plan_name"`
        Price            float64 `json:"price"`
        DurationDays     int     `json:"duration_days"`
        CardType        string  `json:"card_type"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    validCardTypes := map[string]bool{
        "Ordinary": true,
        "Silver":   true,
        "Gold":     true,
    }
    if _, valid := validCardTypes[input.CardType]; !valid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card type"})
        return
    }

    plan := &models.SubscriptionPlan{
        PlanID:           uint(planID),
        PlanName:         input.PlanName,
        Price:            input.Price,
        DurationDays:     input.DurationDays,
        CardType:      input.CardType,
    }

    log.Printf("Updating plan with ID %d: %+v", plan.PlanID, plan) // Debug log

    err = h.SubscriptionUsecase.UpdateSubscriptionPlan(plan)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Subscription plan updated successfully"})
}


// DeleteSubscriptionPlan deletes a subscription plan by its ID
func (h *SubscriptionHandler) DeleteSubscriptionPlan(c *gin.Context) {
    planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
        return
    }

    err = h.SubscriptionUsecase.DeleteSubscriptionPlan(uint(planID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Subscription plan deleted successfully"})
}

// GetAllSubscriptionPlans retrieves all subscription plans
func (h *SubscriptionHandler) GetAllSubscriptionPlans(c *gin.Context) {
    plans, err := h.SubscriptionUsecase.GetAllSubscriptionPlans()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, plans)
}







