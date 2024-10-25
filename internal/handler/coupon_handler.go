package handler

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/usecase"
)

type CouponHandler struct {
    CouponUsecase usecase.CouponUsecase
}

func NewCouponHandler(cu usecase.CouponUsecase) *CouponHandler {
    return &CouponHandler{CouponUsecase: cu}
}

func (h *CouponHandler) CreateCoupon(c *gin.Context) {
    var coupon models.Coupon
    if err := c.ShouldBindJSON(&coupon); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.CouponUsecase.CreateCoupon(coupon); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Coupon created successfully"})
}

func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var coupon models.Coupon
    if err := c.ShouldBindJSON(&coupon); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    coupon.CouponID = id
    if err := h.CouponUsecase.UpdateCoupon(coupon); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Coupon updated successfully"})
}

func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.CouponUsecase.DeleteCoupon(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Coupon deleted successfully"})
}

func (h *CouponHandler) GetCouponByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    coupon, err := h.CouponUsecase.GetCouponByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
        return
    }
    c.JSON(http.StatusOK, coupon)
}

func (h *CouponHandler) GetAllCoupons(c *gin.Context) {
    coupons, err := h.CouponUsecase.GetAllCoupons()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupons"})
        return
    }
    c.JSON(http.StatusOK, coupons)
}
