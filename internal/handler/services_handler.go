package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func ListServices(c *gin.Context) {
	category := c.Param("category") 

	var services []string

	switch category {
	case "Metro":
		services = []string{
			"Routes",
			"Stops",
			"Nearest Stop",
			"Fare Calculation In Regards WithCard Type",
			"Get Time It Take",
			"NolCard Topups",
			"NolCard Balance",
			"Bookings",
		}
	case "Bus":
		services = []string{
			"Routes",
			"Stops",
			"Fare Calculation In Regards WithCard Type",
			"Get Time It Take",
			"NolCard Topups",
			"Balance",
		}
	case "Rental":
		c.JSON(http.StatusOK, gin.H{
			"message": "Service will be available soon",
		})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please choose a valid category (Metro/Bus/Rental)",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"services": services,
	})
}
