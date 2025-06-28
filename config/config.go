package config

import (
    "fmt"
    "log"
    "os"
	"github.com/joho/godotenv"
    "github.com/Prototype-1/xtrace/internal/models"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

var GoogleOAuthConfig = &oauth2.Config{
     ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
    RedirectURL:  "http://localhost:8000/user/google/callback",
    Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
    Endpoint:     google.Endpoint,
}

var DB *gorm.DB

func Connect() {
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
        log.Fatal("One or more required environment variables are not set")
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    log.Println("Successfully connected to the database")

    err = DB.AutoMigrate(
        &models.User{}, 
        &models.Booking{}, 
        &models.Category{}, 
        &models.Coupon{}, 
        &models.Invoice{}, 
        &models.NolCardTopup{}, 
        &models.RazorpayPayment{}, 
        &models.Route{}, 
        &models.Subscription{}, 
        &models.SubscriptionPlan{}, 
        &models.UserSession{}, 
        &models.OTP{}, 
        &models.Wallet{}, 
        &models.WalletTransaction{}, 
        &models.NolCard{}, 
        &models.UserFavorite{}, 
        &models.Stop{}, 
        &models.RouteStop{}, 
        &models.OrderedStop{}, 
        &models.FareRule{}, 
        &models.StopDuration{}, 
    )
    if err != nil {
        log.Fatalf("Error running migrations: %v", err)
    }
    log.Println("Database migration completed")
}



