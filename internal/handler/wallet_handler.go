package handler

import (
	"net/http"
	"strconv"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"
    "fmt"
)

type WalletHandler struct {
	WalletUsecase usecase.WalletUsecase
    RazorpayPaymentUsecase usecase.RazorpayPaymentUsecase
}

func NewWalletHandler(walletUsecase usecase.WalletUsecase, razorpayPaymentUsecase usecase.RazorpayPaymentUsecase) *WalletHandler {
	return &WalletHandler{
        WalletUsecase: walletUsecase,
        RazorpayPaymentUsecase: razorpayPaymentUsecase,
    }
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	wallet, err := h.WalletUsecase.CreateWallet(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet created successfully", "wallet": wallet})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	wallet, err := h.WalletUsecase.GetWalletByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Wallet not found for user ID: %d", userID)})
		return
	}

    response := gin.H{"balance": wallet.Balance} 
	c.JSON(http.StatusOK, response)
}

func (h *WalletHandler) TopUpWalletByAdmin(c *gin.Context) {
    var input struct {
        WalletID    uint    `json:"wallet_id" binding:"required"`
        AdminID     uint    `json:"admin_id" binding:"required"`
        Amount      float64 `json:"amount" binding:"required,numeric"`
        Description string   `json:"description" binding:"required"`
        TransactionType string   `json:"transaction_type" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.WalletUsecase.TopUpWallet(&input.WalletID, &input.AdminID, input.Amount, input.Description, input.TransactionType); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up wallet"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Wallet topped up successfully by admin"})
}

func (h *WalletHandler) TopUpWalletByUser(c *gin.Context) {
    userIDParam := c.Param("userID") 
    userID, err := strconv.Atoi(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    log.Printf("User %d is topping up their wallet", userID)

    var input struct {
        Amount      float64 `json:"amount" binding:"required,numeric"`
        Description string   `json:"description" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    wallet, err := h.WalletUsecase.GetWalletByUserID(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wallet"})
        return
    }
    var adminID *uint = nil 
    transactionType := "top-up" 

    if err := h.WalletUsecase.TopUpWallet(&wallet.WalletID, adminID, input.Amount, input.Description, transactionType); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up wallet"})
        return
    }
    walletIDPtr := &wallet.WalletID

    var nolCardIDPtr, subscriptionIDPtr, bookingIDPtr *uint
    nolCardIDPtr = nil
    subscriptionIDPtr = nil
    bookingIDPtr = nil

    payment, err := h.RazorpayPaymentUsecase.CreatePayment(
        uint(userID),
        input.Amount,
        "INR",
        "", 
        "wallet_topup", 
        walletIDPtr,   
        nolCardIDPtr, 
        subscriptionIDPtr, 
        bookingIDPtr, 
        "",           
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
        return
    }

	updatedWallet, err := h.WalletUsecase.GetWalletByUserID(uint(userID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated wallet"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":      "Please proceed with the payment",
        "new_balance":  updatedWallet.Balance,
        "payment":      payment, 
    })
}

func (h *WalletHandler) MakePayment(c *gin.Context) {
    var input struct {
        WalletID      uint    `json:"wallet_id" binding:"required"`
        Amount        float64 `json:"amount" binding:"required,numeric"`
        TransactionType string `json:"transaction_type" binding:"required"` // New field
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        log.Printf("Error binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.WalletUsecase.MakePayment(input.WalletID, input.Amount, input.TransactionType); err != nil {
        log.Printf("Failed to make payment: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Payment made successfully"})
}


func (h *WalletHandler) GetWalletTransactions(c *gin.Context) {
    userIDParam := c.Param("userID")
    userID, err := strconv.ParseUint(userIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    wallet, err := h.WalletUsecase.GetWalletByUserID(uint(userID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
        return
    }
    transactions, err := h.WalletUsecase.GetWalletTransactions(wallet.WalletID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *WalletHandler) GetWalletTransactionsAdmin(c *gin.Context) {
    walletIDParam := c.Param("wallet_id")
    walletID, err := strconv.ParseUint(walletIDParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
        return
    }

    transactions, err := h.WalletUsecase.GetWalletTransactions(uint(walletID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "No transactions found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}


