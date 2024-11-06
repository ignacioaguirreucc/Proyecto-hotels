package reservations

import (
	"context"
	"hotels-api/domain/reservations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateReservation(ctx context.Context, reservation reservations.Reservation) (string, error)
	GetReservationsByHotelID(ctx context.Context, hotelID string) ([]reservations.Reservation, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{service: service}
}

// Endpoint para crear una reserva
func (c Controller) CreateReservation(ctx *gin.Context) {
	var reservation reservations.Reservation
	if err := ctx.ShouldBindJSON(&reservation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := c.service.CreateReservation(ctx.Request.Context(), reservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating reservation"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// Endpoint para obtener reservas de un hotel específico
func (c Controller) GetReservationsByHotelID(ctx *gin.Context) {
	hotelID := ctx.Param("hotel_id")
	reservations, err := c.service.GetReservationsByHotelID(ctx.Request.Context(), hotelID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching reservations"})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}
