package repository

import (
    "time"
    "fmt"
    "log"
    "gorm.io/gorm"
    "github.com/Prototype-1/xtrace/internal/models"
)

type SubscriptionRepository interface {
    CreateSubscription(subscription *models.Subscription) error
    GetActiveSubscriptionByNolCardID(nolCardID uint) (*models.Subscription, error)
    GetSubscriptionByID(subscriptionID uint) (*models.Subscription, error)
    GetUserSubscriptions(userID uint) ([]models.Subscription, error)
    UpdateSubscription(subscription *models.Subscription) error
    UpdateSubscriptionEndDate(subscriptionID uint, newEndDate time.Time) error
    GetAllSubscriptions() ([]models.Subscription, error)
    GetSubscriptionPlan(planID uint) (*models.SubscriptionPlan, error)
    GetSubscriptionByUserAndCard(userID uint, nolCardID uint) (*models.Subscription, error)
    GetPaymentIDByOrderID(orderID string) (string, error)
}

type subscriptionRepository struct {
    db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
    return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
    return r.db.Create(subscription).Error
}

func (r *subscriptionRepository) GetSubscriptionByID(subscriptionID uint) (*models.Subscription, error) {
    var subscription models.Subscription
    err := r.db.Where("subscription_id = ?", subscriptionID).First(&subscription).Error
    if err != nil {
        return nil, err
    }
    return &subscription, nil
}

func (r *subscriptionRepository) GetActiveSubscriptionByNolCardID(nolCardID uint) (*models.Subscription, error) {
    var subscription models.Subscription
    result := r.db.Where("nol_card_id = ? AND end_date > ?", nolCardID, time.Now()).First(&subscription)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil 
        }
        return nil, result.Error
    }
    return &subscription, nil
}

func (r *subscriptionRepository) GetUserSubscriptions(userID uint) ([]models.Subscription, error) {
    var subscriptions []models.Subscription
    result := r.db.Where("user_id = ?", userID).Find(&subscriptions)
    if result.Error != nil {
        return nil, result.Error
    }
    return subscriptions, nil
}

func (r *subscriptionRepository) UpdateSubscriptionEndDate(subscriptionID uint, newEndDate time.Time) error {
    result := r.db.Model(&models.Subscription{}).
        Where("subscription_id = ?", subscriptionID).
        Update("end_date", newEndDate)
    if result.Error != nil {
        log.Printf("Error updating end_date: %v", result.Error)
        return result.Error
    }
    return nil
}

func (r *subscriptionRepository) UpdateSubscription(subscription *models.Subscription) error {
    r.db = r.db.Debug()

    result := r.db.Table("subscriptions").
        Where("subscription_id = ?", subscription.SubscriptionID).
        Updates(map[string]interface{}{
            "plan_id":        subscription.PlanID,
            "price":          subscription.Price,
            "end_date":      subscription.EndDate,
            "duration_days": subscription.DurationDays,
            "updated_at":    time.Now(),
        })

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("no rows updated")
    }

    return nil
}

func (r *subscriptionRepository) GetAllSubscriptions() ([]models.Subscription, error) {
    var subscriptions []models.Subscription
    result := r.db.Find(&subscriptions)
    if result.Error != nil {
        return nil, result.Error
    }
    return subscriptions, nil
}

func (r *subscriptionRepository) GetSubscriptionPlan(planID uint) (*models.SubscriptionPlan, error) {
    var plan models.SubscriptionPlan
    err := r.db.Where("plan_id = ?", planID).First(&plan).Error
    if err != nil {
        return nil, err
    }
    return &plan, nil
}

func (r *subscriptionRepository) GetSubscriptionByUserAndCard(userID uint, nolCardID uint) (*models.Subscription, error) {
    var subscription models.Subscription
    err := r.db.Where("user_id = ? AND nol_card_id = ? AND end_date > ?", userID, nolCardID, time.Now()).First(&subscription).Error
    if err != nil {
        return nil, err
    }
    return &subscription, nil
}

func (r *subscriptionRepository) GetPaymentIDByOrderID(orderID string) (string, error) {
    var payment models.RazorpayPayment
    err := r.db.Table("payments").Where("order_id = ?", orderID).First(&payment).Error
    if err != nil {
        return "", err
    }
    return payment.RazorpayID, nil
}




