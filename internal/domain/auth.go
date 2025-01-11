package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OSRMService struct {
	BaseURL string
}

func NewOSRMService() *OSRMService {
	return &OSRMService{
		BaseURL: "http://localhost:5000",
	}
}

// Func to get travel time
func (s *OSRMService) GetTravelTime(fromLat, fromLon, toLat, toLon float64) (float64, float64, error) {
	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?overview=false", s.BaseURL, fromLon, fromLat, toLon, toLat)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, 0, fmt.Errorf("failed to decode JSON response: %w", err)
	}
	routes, ok := result["routes"].([]interface{})
	if !ok || len(routes) == 0 {
		return 0, 0, fmt.Errorf("no routes found in response")
	}
	route, ok := routes[0].(map[string]interface{})
	if !ok {
		return 0, 0, fmt.Errorf("unexpected route structure in response")
	}
	legs, ok := route["legs"].([]interface{})
	if !ok || len(legs) == 0 {
		return 0, 0, fmt.Errorf("no legs found in the route")
	}
	leg, ok := legs[0].(map[string]interface{})
	if !ok {
		return 0, 0, fmt.Errorf("unexpected leg structure in response")
	}

	duration, ok := leg["duration"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("unable to extract duration")
	}

	distance, ok := leg["distance"].(float64)
	if !ok {
		return 0, 0, fmt.Errorf("unable to extract distance")
	}

	return duration, distance, nil
}
