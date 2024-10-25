package usecase

import (
    "time"
    "fmt"
    "github.com/Prototype-1/xtrace/internal/models"
    "github.com/Prototype-1/xtrace/internal/repository"
    "github.com/razorpay/razorpay-go" 
)

type SubscriptionUsecase interface {
    // Subscription methods
    CreateSubscription(userID uint, planID uint, serviceType string, durationDays int, cardType string, nolCardID uint) error
    GetUserSubscriptions(userID uint) ([]models.Subscription, error)
    ExtendSubscription(subscriptionID uint, planID uint) error
    UpdateSubscription(subscriptionID uint, price float64, durationDays int) error
    GetAllSubscriptions() ([]models.Subscription, error)
    
    // Subscription Plan methods
    CreateSubscriptionPlan(plan *models.SubscriptionPlan) error
    GetSubscriptionPlan(planID uint) (*models.SubscriptionPlan, error)
    UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error
    DeleteSubscriptionPlan(planID uint) error
    GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error)
    HasActiveSubscription(userID uint) (bool, error)
    GetSubscriptionByID(subscriptionID uint) (*models.Subscription, error)
    GetPaymentIDByOrderID(orderID string) (string, error)
}



type subscriptionUsecase struct {
    subscriptionRepo repository.SubscriptionRepository
    planRepo          repository.SubscriptionPlanRepository
    razorpayClient   *razorpay.Client 
}

func NewSubscriptionUsecase(subscriptionRepo repository.SubscriptionRepository, planRepo repository.SubscriptionPlanRepository, razorpayClient *razorpay.Client) SubscriptionUsecase {
    return &subscriptionUsecase{
        subscriptionRepo: subscriptionRepo,
        planRepo:         planRepo,
        razorpayClient:   razorpayClient, 
    }
}

func (u *subscriptionUsecase) CreateSubscription(userID, planID uint, serviceType string, durationDays int, cardType string, nolCardID uint) error {
    activeSubscription, err := u.subscriptionRepo.GetActiveSubscriptionByNolCardID(nolCardID)
    if err != nil {
        return err
    }

    if activeSubscription != nil {
        return fmt.Errorf("this NolCard already has an active subscription")
    }

    plan, err := u.subscriptionRepo.GetSubscriptionPlan(planID)
    if err != nil {
        return fmt.Errorf("invalid plan ID: %v", err)
    }

    if plan.CardType != cardType {
        return fmt.Errorf("card type '%s' does not match the plan's card type '%s'", cardType, plan.CardType)
    }

    newSubscription := &models.Subscription{
        UserID:       userID,
        PlanID:       planID,
        NolCardID:    nolCardID,   
        ServiceType:  plan.PlanName,
        StartDate:    time.Now(),
        EndDate:      time.Now().AddDate(0, 0, durationDays),
        Price:        plan.Price,
        DurationDays: plan.DurationDays,
        CardType:     cardType,
    }

    return u.subscriptionRepo.CreateSubscription(newSubscription)
}

func (u *subscriptionUsecase) ExtendSubscription(subscriptionID uint, planID uint) error {
    // Retrieve the subscription by ID
    subscription, err := u.subscriptionRepo.GetSubscriptionByID(subscriptionID)
    if err != nil {
        return fmt.Errorf("subscription not found: %w", err)
    }

    // Retrieve the plan details by ID
    plan, err := u.planRepo.GetSubscriptionPlanByID(planID) // Updated method name
    if err != nil {
        return fmt.Errorf("plan not found: %w", err)
    }

    // Update subscription details with new plan information
    subscription.PlanID = plan.PlanID
    subscription.Price = plan.Price
    subscription.DurationDays = plan.DurationDays
    subscription.EndDate = subscription.EndDate.AddDate(0, 0, plan.DurationDays)

    // Update the subscription in the repository
    err = u.subscriptionRepo.UpdateSubscription(subscription)
    if err != nil {
        return fmt.Errorf("failed to update subscription: %w", err)
    }

    return nil
}

func (u *subscriptionUsecase) UpdateSubscription(subscriptionID uint, price float64, durationDays int) error {
    subscription, err := u.subscriptionRepo.GetSubscriptionByID(subscriptionID)
    if err != nil {
        return err
    }

    subscription.Price = price
    subscription.EndDate = time.Now().AddDate(0, 0, durationDays)

    return u.subscriptionRepo.UpdateSubscription(subscription)
}

func (u *subscriptionUsecase) GetAllSubscriptions() ([]models.Subscription, error) {
    return u.subscriptionRepo.GetAllSubscriptions()
}

func (u *subscriptionUsecase) GetUserSubscriptions(userID uint) ([]models.Subscription, error) {
    return u.subscriptionRepo.GetUserSubscriptions(userID)
}

// Subscription Plan methods
func (u *subscriptionUsecase) CreateSubscriptionPlan(plan *models.SubscriptionPlan) error {
    return u.planRepo.CreateSubscriptionPlan(plan)
}

func (u *subscriptionUsecase) GetSubscriptionPlan(planID uint) (*models.SubscriptionPlan, error) {
    return u.planRepo.GetSubscriptionPlanByID(planID)
}

func (u *subscriptionUsecase) UpdateSubscriptionPlan(plan *models.SubscriptionPlan) error {
    return u.planRepo.UpdateSubscriptionPlan(plan)
}

func (u *subscriptionUsecase) DeleteSubscriptionPlan(planID uint) error {
    return u.planRepo.DeleteSubscriptionPlan(planID)
}

func (u *subscriptionUsecase) GetAllSubscriptionPlans() ([]models.SubscriptionPlan, error) {
    return u.planRepo.GetAllSubscriptionPlans()
}

func (u *subscriptionUsecase) GetSubscriptionByID(subscriptionID uint) (*models.Subscription, error) {
    return u.subscriptionRepo.GetSubscriptionByID(subscriptionID)
}

func (u *subscriptionUsecase) HasActiveSubscription(userID uint) (bool, error) {
    subscriptions, err := u.subscriptionRepo.GetUserSubscriptions(userID)
    if err != nil {
        return false, err
    }

    for _, subscription := range subscriptions {
        if subscription.EndDate.After(time.Now()) {
            return true, nil 
        }
    }

    return false, nil 
}

func (u *subscriptionUsecase) GetPaymentIDByOrderID(orderID string) (string, error) {
    paymentID, err := u.subscriptionRepo.GetPaymentIDByOrderID(orderID)
    if err != nil {
        return "", err
    }
    return paymentID, nil
}




