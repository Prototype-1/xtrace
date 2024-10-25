package repository

import (
    "gorm.io/gorm"
    "github.com/Prototype-1/xtrace/internal/models"
    "time"
)

type CouponRepository interface {
    CreateCoupon(coupon models.Coupon) error
    UpdateCoupon(coupon models.Coupon) error
    DeleteCoupon(id int) error
    GetCouponByID(id int) (models.Coupon, error)
    GetAllCoupons() ([]models.Coupon, error)

    GetCouponByCode(code string) (*models.Coupon, error)
    IsCouponValid(code string) (bool, error)   
    GetCouponsByPaymentType(paymentType string) ([]*models.Coupon, error)
}

type CouponRepositoryImpl struct {
    DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) CouponRepository {
    return &CouponRepositoryImpl{DB: db}
}

func (r *CouponRepositoryImpl) CreateCoupon(coupon models.Coupon) error {
    return r.DB.Create(&coupon).Error
}

func (r *CouponRepositoryImpl) UpdateCoupon(coupon models.Coupon) error {
    return r.DB.Save(&coupon).Error
}

func (r *CouponRepositoryImpl) DeleteCoupon(id int) error {
    return r.DB.Delete(&models.Coupon{}, id).Error
}

func (r *CouponRepositoryImpl) GetCouponByID(id int) (models.Coupon, error) {
    var coupon models.Coupon
    err := r.DB.First(&coupon, id).Error
    return coupon, err
}

func (r *CouponRepositoryImpl) GetAllCoupons() ([]models.Coupon, error) {
    var coupons []models.Coupon
    err := r.DB.Find(&coupons).Error
    return coupons, err
}

func (c *CouponRepositoryImpl) GetCouponByCode(code string) (*models.Coupon, error) {
    var coupon models.Coupon
    err := c.DB.Table("coupons").Where("code = ?", code).First(&coupon).Error
    if err != nil {
        return nil, err
    }
    return &coupon, nil
}

func (c *CouponRepositoryImpl) IsCouponValid(code string) (bool, error) {
    coupon, err := c.GetCouponByCode(code)
    if err != nil {
        return false, err
    }
    if coupon == nil {
        return false, nil 
    }
    now := time.Now()
    if coupon.DiscountType != "" && coupon.DiscountAmount > 0 && 
       coupon.StartDate.Before(now) && coupon.EndDate.After(now) {
        return true, nil 
    }
    
    return false, nil 
}

func (r *CouponRepositoryImpl) GetCouponsByPaymentType(paymentType string) ([]*models.Coupon, error) {
	var coupons []*models.Coupon
	if err := r.DB.Table("coupons").Where("payment_type = ?", paymentType).Find(&coupons).Error; err != nil {
		return nil, err
	}

	return coupons, nil
}