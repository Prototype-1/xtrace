package usecase

import (
	"github.com/Prototype-1/xtrace/internal/models"
	"github.com/Prototype-1/xtrace/internal/repository"
    "math"
)

type RouteStopUsecase struct {
	RouteStopRepo repository.RouteStopRepository
}

func NewRouteStopUsecase(repo repository.RouteStopRepository) *RouteStopUsecase {
	return &RouteStopUsecase{RouteStopRepo: repo}
}

func (u *RouteStopUsecase) AddRouteStop(routeStop models.RouteStop) error {
	return u.RouteStopRepo.AddRouteStop(routeStop)
}

func (u *RouteStopUsecase) UpdateRouteStop(id int, routeStop models.RouteStop) error {
	return u.RouteStopRepo.UpdateRouteStop(routeStop)
}

func (u *RouteStopUsecase) DeleteRouteStop(id int) error {
	return u.RouteStopRepo.DeleteRouteStop(id)
}

func (u *RouteStopUsecase) GetAllRouteStops() ([]models.RouteStop, error) {
	return u.RouteStopRepo.GetAllRouteStops()
}

func (u *RouteStopUsecase) GetOrderedStops(routeID uint) ([]models.OrderedStop, error) {
    routeStops, err := u.RouteStopRepo.GetOrderedStopsByRouteID(routeID)
    if err != nil {
        return nil, err
    }

    var orderedStops []models.OrderedStop
    for _, routeStop := range routeStops {
        stop, categoryName, err := u.RouteStopRepo.GetStopByID(routeStop.StopID)
        if err != nil {
            return nil, err
        }

        orderedStops = append(orderedStops, models.OrderedStop{
            StopName:     stop.StopName,
            StopSequence: routeStop.StopSequence,
            Category:     categoryName, 
        })
    }

    return orderedStops, nil
}

func (u *RouteStopUsecase) GetOrderedStopsByCategory(routeID uint, category string) ([]models.OrderedStop, error) {
	routeStops, err := u.RouteStopRepo.GetOrderedStopsByRouteID(routeID)
	if err != nil {
		return nil, err
	}

	var orderedStops []models.OrderedStop
	for _, routeStop := range routeStops {
		stop, categoryName, err := u.RouteStopRepo.GetStopByID(routeStop.StopID)
		if err != nil {
			return nil, err
		}
		if categoryName == category {
			orderedStops = append(orderedStops, models.OrderedStop{
				StopName:     stop.StopName,
				StopSequence: routeStop.StopSequence,
				Category:     categoryName,
			})
		}
	}

	return orderedStops, nil
}

func (u *RouteStopUsecase) FindNearestStop(userLat, userLon float64, routeID int) (models.Stop, error) {
	stops, err := u.RouteStopRepo.GetStopsByRouteID(uint(routeID))
	if err != nil {
		return models.Stop{}, err
	}

	var nearestStop models.Stop
	minDistance := math.MaxFloat64

	for _, stop := range stops {
		distance := haversine(userLat, userLon, stop.Latitude, stop.Longitude)
		if distance < minDistance {
			minDistance = distance
			nearestStop = stop
		}
	}

	return nearestStop, nil
}

// Haversine formula to calculate distance between two latitude/longitude points
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth radius in kilometers

	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}
