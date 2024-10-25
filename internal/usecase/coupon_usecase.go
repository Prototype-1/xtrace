package usecase

import (
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
)

type CouponUsecase interface {
    CreateCoupon(coupon models.Coupon) error
    UpdateCoupon(coupon models.Coupon) error
    DeleteCoupon(id int) error
    GetCouponByID(id int) (models.Coupon, error)
    GetAllCoupons() ([]models.Coupon, error)
}

type CouponUsecaseImpl struct {
    repo repository.CouponRepository
}

func NewCouponUsecase(repo repository.CouponRepository) CouponUsecase {
    return &CouponUsecaseImpl{repo: repo}
}

func (uc *CouponUsecaseImpl) CreateCoupon(coupon models.Coupon) error {
    return uc.repo.CreateCoupon(coupon)
}

func (uc *CouponUsecaseImpl) UpdateCoupon(coupon models.Coupon) error {
    return uc.repo.UpdateCoupon(coupon)
}

func (uc *CouponUsecaseImpl) DeleteCoupon(id int) error {
    return uc.repo.DeleteCoupon(id)
}

func (uc *CouponUsecaseImpl) GetCouponByID(id int) (models.Coupon, error) {
    return uc.repo.GetCouponByID(id)
}

func (uc *CouponUsecaseImpl) GetAllCoupons() ([]models.Coupon, error) {
    return uc.repo.GetAllCoupons()
}
