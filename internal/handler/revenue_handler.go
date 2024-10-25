package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/config"
)

type RevenueHandler struct{}

func NewRevenueHandler() *RevenueHandler {
    return &RevenueHandler{}
}

func (h *RevenueHandler) GetTotalRevenue(c *gin.Context) {
    var totalRevenue float64

    result := config.DB.Table("payments").Where("status = ?", "captured").Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total revenue"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total": totalRevenue})
}

func (h *RevenueHandler) GetMonthlyRevenue(c *gin.Context) {
    type MonthlyRevenue struct {
        Month string  `json:"month"`
        Total float64 `json:"total"`
    }

    var monthlyRevenues []MonthlyRevenue

    query := `
        SELECT TO_CHAR(payment_date, 'YYYY-MM') AS month, COALESCE(SUM(amount), 0) AS total
        FROM payments
        WHERE status = 'captured'
        GROUP BY month
        ORDER BY month;
    `
    result := config.DB.Raw(query).Scan(&monthlyRevenues)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch monthly revenue"})
        return
    }

    c.JSON(http.StatusOK, monthlyRevenues)
}
