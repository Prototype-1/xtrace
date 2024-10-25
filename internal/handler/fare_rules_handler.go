package handler

import (
	"math"
	"net/http"
	"strconv"
    "fmt"
	"github.com/Prototype-1/xtrace/config"
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/usecase"
	"github.com/gin-gonic/gin"
)

type FareRuleHandler struct {
    FareRuleUsecase usecase.FareRuleUsecase
}

func NewFareRuleHandler(fareRuleUsecase usecase.FareRuleUsecase) *FareRuleHandler {
    return &FareRuleHandler{FareRuleUsecase: fareRuleUsecase}
}

func (h *FareRuleHandler) CreateFareRule(c *gin.Context) {
    var fareRule models.FareRule
    if err := c.ShouldBindJSON(&fareRule); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if err := h.FareRuleUsecase.CreateFareRule(fareRule); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create fare rule"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Fare rule created successfully"})
}

func (h *FareRuleHandler) UpdateFareRule(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    var fareRule models.FareRule
    if err := c.ShouldBindJSON(&fareRule); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    fareRule.FareRuleID = id
    if err := h.FareRuleUsecase.UpdateFareRule(fareRule); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fare rule"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Fare rule updated successfully"})
}

func (h *FareRuleHandler) DeleteFareRule(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.FareRuleUsecase.DeleteFareRule(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete fare rule"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Fare rule deleted successfully"})
}

func (h *FareRuleHandler) GetFareRuleByID(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    fareRule, err := h.FareRuleUsecase.GetFareRuleByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Fare rule not found"})
        return
    }
    c.JSON(http.StatusOK, fareRule)
}

func (h *FareRuleHandler) GetAllFareRules(c *gin.Context) {
    fareRules, err := h.FareRuleUsecase.GetAllFareRules()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get fare rules"})
        return
    }
    c.JSON(http.StatusOK, fareRules)
}

//Haversine equation
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 
    lat1Rad := lat1 * math.Pi / 180
    lon1Rad := lon1 * math.Pi / 180
    lat2Rad := lat2 * math.Pi / 180
    lon2Rad := lon2 * math.Pi / 180

    dlat := lat2Rad - lat1Rad
    dlon := lon2Rad - lon1Rad

    a := math.Sin(dlat/2)*math.Sin(dlat/2) +
        math.Cos(lat1Rad)*math.Cos(lat2Rad)*
            math.Sin(dlon/2)*math.Sin(dlon/2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

    distance := R * c
    return distance
}

func (h *FareRuleHandler) CalculateFare(c *gin.Context) {
    routeID, _ := strconv.Atoi(c.Param("route_id"))
    startStopSeq, _ := strconv.Atoi(c.Param("start_stop_sequence"))
    endStopSeq, _ := strconv.Atoi(c.Param("end_stop_sequence"))

    cardType := c.Query("cardType")
    if cardType == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Card type is required"})
        return
    }

    fareRule, err := h.FareRuleUsecase.GetFareRuleByRouteID(routeID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Fare rule not found"})
        return
    }

    // Get the latitude and longitude of the start and end stops
    var startStop, endStop models.Stop
    if err := config.DB.Where("stop_id = ?", startStopSeq).First(&startStop).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Start stop not found"})
        return
    }
    if err := config.DB.Where("stop_id = ?", endStopSeq).First(&endStop).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "End stop not found"})
        return
    }
    traveledKm := Haversine(float64(startStop.Latitude), float64(startStop.Longitude),
        float64(endStop.Latitude), float64(endStop.Longitude))

    numberOfStops := int(math.Abs(float64(endStopSeq - startStopSeq)))
    additionalStops := numberOfStops - fareRule.BaseStops

    var baseFare float64
    switch cardType {
    case "Ordinary":
        baseFare = fareRule.OrdinaryFare
    case "Silver":
        baseFare = fareRule.SilverFare
    case "Gold":
        baseFare = fareRule.GoldFare
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card type"})
        return
    }

    totalFare := baseFare

    additionalKm := traveledKm - fareRule.BaseKm
    if additionalKm > 0 {
        totalFare += additionalKm * fareRule.FarePerKm
    }
    if additionalStops > 0 {
        totalFare += float64(additionalStops) * fareRule.FarePerStop
    }
    if totalFare < fareRule.OrdinaryFare {
        totalFare = fareRule.OrdinaryFare
    }
    c.JSON(http.StatusOK, gin.H{"total_fare": totalFare})
}

func (h *FareRuleHandler) CalculateTravelTimes(c *gin.Context) {
    var input struct {
        StopPairs []struct {
            FromLat float64 `json:"from_lat"`
            FromLon float64 `json:"from_lon"`
            ToLat   float64 `json:"to_lat"`
            ToLon   float64 `json:"to_lon"`
        } `json:"stop_pairs"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
        return
    }

    fmt.Printf("Received input: %+v\n", input)

    var travelTimes []map[string]interface{}
    var errors []string

    for _, pair := range input.StopPairs {
        fmt.Printf("Processing travel time for %f,%f -> %f,%f\n", pair.FromLat, pair.FromLon, pair.ToLat, pair.ToLon)

        travelTime, err := h.FareRuleUsecase.GetTravelTimeByCoordinates(pair.FromLat, pair.FromLon, pair.ToLat, pair.ToLon)
        if err != nil {
            // Log the error and set travel time to -1 to indicate failure
            errorMsg := fmt.Sprintf("Error calculating time between %f,%f and %f,%f: %v", pair.FromLat, pair.FromLon, pair.ToLat, pair.ToLon, err)
            fmt.Println(errorMsg)
            errors = append(errors, errorMsg)

            travelTime = -1 // Indicate an error with a placeholder value
        }

        travelTimes = append(travelTimes, map[string]interface{}{
            "from_lat":    pair.FromLat,
            "from_lon":    pair.FromLon,
            "to_lat":      pair.ToLat,
            "to_lon":      pair.ToLon,
            "travel_time": travelTime,
            "error":       err != nil, // Set true if there was an error
        })
    }

    // Return results
    c.JSON(http.StatusOK, gin.H{
        "travel_times": travelTimes,
        "errors":       errors,
    })
}

