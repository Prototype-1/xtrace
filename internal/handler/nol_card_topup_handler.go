package handler

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/usecase"
    "fmt"
)

type NolCardTopupHandler struct {
    NolCardTopupUsecase usecase.NolCardTopupUsecase
}

func NewNolCardTopupHandler(u usecase.NolCardTopupUsecase) *NolCardTopupHandler {
    return &NolCardTopupHandler{NolCardTopupUsecase: u}
}

func (h *NolCardTopupHandler) AddTopup(c *gin.Context) {
    
    var topup models.NolCardTopup

    if err := c.ShouldBindJSON(&topup); err != nil {
        fmt.Println("Error while binding JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    fmt.Printf("Received topup request with NolCardID: %d, Amount: %.2f\n", topup.NolCardID, topup.Amount)

    err := h.NolCardTopupUsecase.AddTopupAndUpdateBalance(topup)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Payment Successful",
        "amount": topup.Amount,
    })
}

func (h *NolCardTopupHandler) GetTopupsByCardID(c *gin.Context) {
    nolCardID, err := strconv.Atoi(c.Param("nol_card_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nol Card ID"})
        return
    }

    topups, err := h.NolCardTopupUsecase.GetTopupsByCardID(nolCardID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top-ups"})
        return
    }

    c.JSON(http.StatusOK, topups)
}

func (h *NolCardTopupHandler) GetTopupByID(c *gin.Context) {
    topupID, err := strconv.Atoi(c.Param("topup_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Top-up ID"})
        return
    }

    topup, err := h.NolCardTopupUsecase.GetTopupByID(topupID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get top-up"})
        return
    }

    c.JSON(http.StatusOK, topup)
}

func (h *NolCardTopupHandler) GetNolCardBalance(c *gin.Context) {
    userIDParam := c.Param("userID")
    userID, err := strconv.Atoi(userIDParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }
    nolCard, err := h.NolCardTopupUsecase.GetNolCardByUserID(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Nol card details"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"balance": nolCard.Balance})
}