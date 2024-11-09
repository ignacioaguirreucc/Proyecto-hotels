package reservations

import (
	"context"
	"hotels-api/domain/reservations"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateReservation(ctx context.Context, reservation reservations.Reservation) (string, error)
	GetReservationsByUserID(ctx context.Context, userID string) ([]reservations.Reservation, error)
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

	// Llamar al servicio para crear la reserva
	id, err := c.service.CreateReservation(ctx.Request.Context(), reservation)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error creating reservation"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// Obtener reservas del usuario autenticado
func (c Controller) GetReservationsByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	reservations, err := c.service.GetReservationsByUserID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching reservations"})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}
