package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/domain"
	"github.com/Prototype-1/xtrace/internal/handler"
	"github.com/Prototype-1/xtrace/internal/middleware"
	"github.com/Prototype-1/xtrace/internal/repository"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/razorpay/razorpay-go"
)

func createRazorpayOrder() error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	payload := map[string]interface{}{
		"amount":   10000,
		"currency": "INR",
		"receipt":  "receipt#1",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET"))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	fmt.Printf("Razorpay Order Created: %v\n", result)
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := createRazorpayOrder(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	razorpayClient := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_KEY_SECRET"))
	

	config.Connect()
	router := gin.Default()
	router.Use(gin.Logger()) 

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Static("/static", "./static")

	razorpayRepo := repository.NewRazorpayPaymentRepository(config.DB)
	couponRepo := repository.NewCouponRepository(config.DB)
	razorpayUsecase := usecase.NewRazorpayPaymentUsecase(razorpayRepo, razorpayClient, couponRepo)

	userRepo := repository.NewUserRepository(config.DB)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	ticker := time.NewTicker(24 * time.Hour)
    go func() {
        for {
            <-ticker.C
            err := userUsecase.UnsuspendInactiveUsers()
            if err != nil {
                log.Printf("Error running unsuspend task: %v\n", err)
            }
        }
    }()

	categoryRepo := &repository.CategoryRepositoryImpl{
		DB: config.DB,
	}
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := &handler.CategoryHandler{CategoryUsecase: categoryUsecase}

	routeRepo := repository.NewRouteRepository(config.DB)
	routeUsecase := usecase.NewRouteUsecase(routeRepo)
	routeHandler := handler.NewRouteHandler(routeUsecase)

	userFavoritesRepo := repository.NewUserFavoritesRepository(config.DB)
	userFavoritesUsecase := usecase.NewUserFavoritesUsecase(userFavoritesRepo)
	userFavoritesHandler := handler.NewUserFavoritesHandler(userFavoritesUsecase)

	stopRepo := repository.NewStopRepository(config.DB)
	stopUsecase := usecase.NewStopUsecase(stopRepo)
	stopHandler := handler.NewStopHandler(stopUsecase)

	routeStopRepo := repository.NewRouteStopRepository(config.DB)
	routeStopUsecase := usecase.NewRouteStopUsecase(routeStopRepo)
	routeStopHandler := handler.NewRouteStopHandler(routeStopUsecase)

	fareRuleRepo := repository.NewFareRuleRepository(config.DB)
	osrmService := domain.NewOSRMService()
	fareRuleUsecase := usecase.NewFareRuleUsecase(fareRuleRepo, osrmService)
	fareRuleHandler := handler.NewFareRuleHandler(fareRuleUsecase)

	couponUsecase := usecase.NewCouponUsecase(couponRepo)
	couponHandler := handler.NewCouponHandler(couponUsecase)

	nolCardTopupRepo := repository.NewNolCardTopupRepository(config.DB)
	nolCardTopupUsecase := usecase.NewNolCardTopupUsecase(nolCardTopupRepo)
	nolCardTopupHandler := handler.NewNolCardTopupHandler(nolCardTopupUsecase)

	nolCardRepo := repository.NewNolCardRepository(config.DB)
	nolCardUsecase := usecase.NewNolCardUsecase(nolCardRepo)
	nolCardHandler := handler.NewNolCardHandler(nolCardUsecase)

	walletRepo := repository.NewWalletRepository(config.DB)
	walletTransactionRepo := repository.NewWalletTransactionRepository(config.DB)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, walletTransactionRepo)
	walletHandler := handler.NewWalletHandler(walletUsecase, razorpayUsecase)

	subscriptionRepo := repository.NewSubscriptionRepository(config.DB)
	subscriptionPlanRepo := repository.NewSubscriptionPlanRepository(config.DB)
	subscriptionUsecase := usecase.NewSubscriptionUsecase(subscriptionRepo, subscriptionPlanRepo, razorpayClient)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionUsecase, nolCardRepo, subscriptionRepo,  razorpayUsecase, walletUsecase)

	bookingRepo := repository.NewBookingRepository(config.DB)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo)
	bookingHandler := handler.NewBookingHandler(bookingUsecase)

	invoiceRepo := repository.NewInvoiceRepository(config.DB) 
invoiceUsecase := usecase.NewInvoiceUsecase(invoiceRepo)
invoiceHandler := handler.NewInvoiceHandler(userRepo, invoiceRepo)

	razorpayHandler := handler.NewRazorpayHandler(walletUsecase, razorpayUsecase, bookingUsecase, subscriptionUsecase, razorpayClient, nolCardTopupUsecase, invoiceUsecase)

	revenueHandler := handler.NewRevenueHandler()

	router.POST("/admin/signup", handler.AdminSignUp)
	router.POST("/admin/login", handler.AdminLogin)
	router.POST("/admin/logout", middleware.TokenAuthMiddleware(), middleware.AdminAuthMiddleware(), handler.AdminLogout)

	router.POST("/user/signup", handler.UserSignUp)
	router.POST("/user/verify-otp", handler.VerifyOTP)
	router.POST("/user/login", handler.UserLogin)
	router.POST("/user/logout", middleware.TokenAuthMiddleware(), middleware.UserAuthMiddleware(), handler.UserLogout)
	router.POST("/user/resend-otp", handler.ResendOTP)

	router.GET("/user/google/login", handler.GoogleLogin)
	router.GET("/user/google/callback", handler.GoogleCallback)

	router.GET("/:id/status", userHandler.GetUserStatus)
	router.POST("/user/:userID/payment/create", razorpayHandler.CreatePayment)
	router.GET("/user/payment/:type/:id/amount", razorpayHandler.GetAmountByPaymentType)

	router.GET("/user/:userID/wallet/show", walletHandler.GetWallet)
	router.GET("/user/:userID/nol-card/balance/show", nolCardTopupHandler.GetNolCardBalance)
	router.POST("/user/:userID/nol-card/topup", nolCardTopupHandler.AddTopup)

	router.POST("/user/payment/verify", razorpayHandler.VerifyPayment)

	router.GET("/coupons/:paymentType", razorpayHandler.FetchApplicableCoupons)
	router.POST("/coupons/apply", razorpayHandler.ApplyCoupon)

	router.GET("/invoices/getUserEmail", invoiceHandler.GetUserEmail)
	router.POST("/invoices/sendInvoice", invoiceHandler.SendInvoice)
	
	router.GET("/view/users", userHandler.GetAllUsers)
	router.GET("/categories", categoryHandler.GetAllCategoriesAdmin)
	router.GET("/services/:category", handler.ListServices)

	router.GET("/revenue/total", revenueHandler.GetTotalRevenue)
	router.GET("/revenue/monthly", revenueHandler.GetMonthlyRevenue)

	adminRoutes := router.Group("/admin").Use(middleware.TokenAuthMiddleware()).Use(middleware.AdminAuthMiddleware())
	{
		adminRoutes.GET("/users", userHandler.GetAllUsers)
		adminRoutes.PUT("/users/:id/block", userHandler.BlockUser)
		adminRoutes.PUT("/users/:id/unblock", userHandler.UnblockUser)
		adminRoutes.PUT("/users/:id/suspend", userHandler.SuspendUser)

		adminRoutes.POST("/add/categories", categoryHandler.AddCategory)
		adminRoutes.PUT("/categories/update/:id", categoryHandler.UpdateCategory)
		adminRoutes.DELETE("/categories/delete/:id", categoryHandler.DeleteCategory)
		adminRoutes.GET("/categories", categoryHandler.GetAllCategoriesAdmin)

		adminRoutes.POST("/add/routes", routeHandler.AddRoute)
		adminRoutes.PUT("/update/routes/:id", routeHandler.UpdateRoute)
		adminRoutes.DELETE("/delete/routes/:id", routeHandler.DeleteRoute)
		adminRoutes.GET("/routes", routeHandler.GetAllRoutes)

		adminRoutes.POST("/add/stops", stopHandler.AddStop)
		adminRoutes.PUT("/update/stops/:id", stopHandler.UpdateStop)
		adminRoutes.DELETE("/delete/stops/:id", stopHandler.DeleteStop)
		adminRoutes.GET("/stops", stopHandler.GetAllStops)

		adminRoutes.POST("/add/route-stops", routeStopHandler.AddRouteStop)
		adminRoutes.PUT("/update/route-stops/:id", routeStopHandler.UpdateRouteStop)
		adminRoutes.DELETE("/delete/route-stops/:id", routeStopHandler.DeleteRouteStop)
		adminRoutes.GET("/route-stops", routeStopHandler.GetAllRouteStops)

		adminRoutes.POST("/add/fare-rule", fareRuleHandler.CreateFareRule)
		adminRoutes.PUT("/update/fare-rule/:id", fareRuleHandler.UpdateFareRule)
		adminRoutes.DELETE("/delete/fare-rule/:id", fareRuleHandler.DeleteFareRule)
		adminRoutes.GET("/fare-rules", fareRuleHandler.GetAllFareRules)
		adminRoutes.GET("/fare-rule/:id", fareRuleHandler.GetFareRuleByID)

		adminRoutes.POST("/add/coupons", couponHandler.CreateCoupon)
		adminRoutes.PUT("/update/coupons/:id", couponHandler.UpdateCoupon)
		adminRoutes.DELETE("/delete/coupons/:id", couponHandler.DeleteCoupon)
		adminRoutes.GET("/coupons/:id", couponHandler.GetCouponByID)
		adminRoutes.GET("/coupons", couponHandler.GetAllCoupons)

		adminRoutes.POST("/add/nolcard", nolCardHandler.AddNolCard)
		adminRoutes.POST("/add/topup", nolCardTopupHandler.AddTopup)
		adminRoutes.GET("/topups/:nol_card_id", nolCardTopupHandler.GetTopupsByCardID)
		adminRoutes.GET("/topup/:topup_id", nolCardTopupHandler.GetTopupByID)

		adminRoutes.POST("/add/subscription/plans", subscriptionHandler.CreateSubscriptionPlan)
		adminRoutes.PUT("/update/subscription/plans/:id", subscriptionHandler.UpdateSubscriptionPlan)
		adminRoutes.DELETE("/delete/subscription/plans/:id", subscriptionHandler.DeleteSubscriptionPlan)
		adminRoutes.GET("/view/subscription/plans", subscriptionHandler.GetAllSubscriptionPlans)
		adminRoutes.GET("/subscriptions/all", subscriptionHandler.GetAllSubscriptions)

		adminRoutes.POST("/wallet/topup", walletHandler.TopUpWalletByAdmin)
		adminRoutes.GET("/wallet/transactions/:wallet_id", walletHandler.GetWalletTransactionsAdmin)
	}

	userRoutes := router.Group("/user").Use(middleware.TokenAuthMiddleware()).Use(middleware.UserAuthMiddleware())
	{
		userRoutes.GET("/categories", categoryHandler.GetAllCategoriesUser)
		userRoutes.GET("/services/:category", handler.ListServices)

		userRoutes.GET("/routes", routeHandler.GetAllRoutesUser)

		userRoutes.POST("/:userID/favorites/:routeID", userFavoritesHandler.AddFavoriteRoute)
		userRoutes.GET("/:userID/favorites", userFavoritesHandler.GetUserFavoriteRoutes)
		userRoutes.DELETE("/:userID/favorites/:routeID", userFavoritesHandler.RemoveFavoriteRoute)

		userRoutes.GET("/route/stops/:route_id", routeStopHandler.GetOrderedStopsByRoute)
		userRoutes.GET("/nearest-stop", routeStopHandler.FindNearestStop)
		userRoutes.GET("/fare/calculate/:route_id/:start_stop_sequence/:end_stop_sequence", fareRuleHandler.CalculateFare)
		userRoutes.POST("/travel-time", fareRuleHandler.CalculateTravelTimes)

		userRoutes.POST("/add/topup", nolCardTopupHandler.AddTopup)
		userRoutes.GET("/nol-card/:nol_card_id", nolCardHandler.GetNolCardDetails)

		userRoutes.POST("/add/subscriptions", subscriptionHandler.CreateSubscription)
		userRoutes.GET("/subscriptions/:id", subscriptionHandler.GetUserSubscriptions)
		userRoutes.PUT("/extend/subscriptions/:id", subscriptionHandler.ExtendSubscription)

		userRoutes.POST("/:userID/bookings", bookingHandler.CreateBooking)
		
		userRoutes.POST("/:userID/wallet", walletHandler.CreateWallet)
		userRoutes.GET("/:userID/wallet", walletHandler.GetWallet)
		userRoutes.POST("/:userID/wallet/topup", walletHandler.TopUpWalletByUser)
		userRoutes.POST("/:userID/wallet/payment", walletHandler.MakePayment)
		userRoutes.GET("/:userID/wallet/transactions", walletHandler.GetWalletTransactions)
	}

	router.Run(":8000")
}

