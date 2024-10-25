package handler

import (
    "net/http"
    "strconv"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/internal/models" 
    "github.com/Prototype-1/xtrace/internal/usecase"
)

type NolCardHandler struct {
    NolCardUsecase usecase.NolCardUsecase
}

func NewNolCardHandler(u usecase.NolCardUsecase) *NolCardHandler {
    return &NolCardHandler{NolCardUsecase: u}
}

func (h *NolCardHandler) GetNolCardDetails(c *gin.Context) {
    nolCardID, err := strconv.Atoi(c.Param("nol_card_id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Nol Card ID"})
        return
    }

    nolCard, err := h.NolCardUsecase.GetNolCardByID(nolCardID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Nol Card details"})
        return
    }
    var message string
    switch nolCard.CardType {
    case "gold":
        if nolCard.Balance < 50 {
            message = "Insufficient balance in your Gold Pass card. Please top up."
        } else {
            message = "There is sufficient balance in your card, Happy Journey."
        }
    case "silver":
        if nolCard.Balance < 30 {
            message = "Insufficient balance in your Silver Pass card. Please top up."
        } else {
            message = "There is sufficient balance in your card, Happy Journey."
        }
    default:
        if nolCard.Balance < 20 {
            message = fmt.Sprintf("Your Nol Card balance is %.2f, which is below the minimum balance. Please top up.", nolCard.Balance)
        } else {
            message = fmt.Sprintf("Your Nol Card balance is %.2f, Happy Journey", nolCard.Balance)
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "nol_card_id": nolCard.NolCardID,
        "card_number": nolCard.CardNumber,
        "balance":     nolCard.Balance,
        "card_type":   nolCard.CardType,
        "message":     message,
    })
}

func (h *NolCardHandler) AddNolCard(c *gin.Context) {
    var nolCard models.NolCard
    if err := c.ShouldBindJSON(&nolCard); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    if nolCard.CardType == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Card type is required - Ordinary/Silver/Gold"})
        return
    }

    err := h.NolCardUsecase.AddNolCard(nolCard)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Nol card added successfully"})
}


