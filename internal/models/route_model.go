package models

import "time"

type Route struct {
    RouteID    int       `gorm:"primaryKey;autoIncrement" json:"route_id"`
    RouteName  string    `gorm:"type:varchar(255);not null" json:"route_name"`
    StartStopID int      `gorm:"not null" json:"start_stop_id"`
    EndStopID   int      `gorm:"not null" json:"end_stop_id"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    CategoryID int      `gorm:"not null" json:"category_id"`  
}

type UserFavorite struct {
    UserID    uint `gorm:"primaryKey"`
    RouteID   uint `gorm:"primaryKey"`
    CreatedAt time.Time
}

type Stop struct {
	StopID     int       `gorm:"primaryKey" json:"stop_id"`
	StopName   string    `json:"stop_name"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}


type RouteStop struct {
	RouteStopID  int       `gorm:"primaryKey"`
	RouteID      int      `json:"route_id" binding:"required"`
	StopID       int       `json:"stop_id" binding:"required"`
	StopSequence int       `json:"stop_sequence" binding:"required"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

type OrderedStop struct {
    StopName     string `json:"stop_name"`
    StopSequence int    `json:"stop_sequence"`
    Category string `json:"category"`
}


type FareRule struct {
    FareRuleID  int       `gorm:"primary_key;auto_increment" json:"fare_rule_id"`
    RouteID     int       `json:"route_id"`
    OrdinaryFare float64  `json:"ordinary_fare"`   
    SilverFare  float64   `json:"silver_fare"`  
    GoldFare    float64   `json:"gold_fare"`
    FarePerKm   float64   `json:"fare_per_km"`
    FarePerStop float64   `json:"fare_per_stop"`
    BaseKm      float64   `json:"base_km"`
    BaseStops   int       `json:"base_stops"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
    CategoryID  int       `json:"category_id"`
}

type StopDuration struct {
    StopDurationID    uint      `gorm:"primaryKey;autoIncrement" json:"stop_duration_id"`
    RouteID           uint      `json:"route_id"`              
    FromStopID        uint      `json:"from_stop_id"`           
    ToStopID          uint      `json:"to_stop_id"`             
    TravelTimeMinutes int       `json:"travel_time_minutes"`   
    CreatedAt         time.Time `json:"created_at"`             
    UpdatedAt         time.Time `json:"updated_at"`             
}