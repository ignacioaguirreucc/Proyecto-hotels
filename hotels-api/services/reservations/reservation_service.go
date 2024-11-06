// services/reservations/reservation_service.go
package reservations

import (
	"context"
	"fmt"
	"hotels-api/domain/reservations"
)

type Repository interface {
	Create(ctx context.Context, reservation reservations.Reservation) (string, error)
	GetByHotelID(ctx context.Context, hotelID string) ([]reservations.Reservation, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository: repository}
}

func (s Service) CreateReservation(ctx context.Context, reservation reservations.Reservation) (string, error) {
	// Se puede agregar lógica para validar disponibilidad aquí.
	id, err := s.repository.Create(ctx, reservation)
	if err != nil {
		return "", fmt.Errorf("error creating reservation: %w", err)
	}
	return id, nil
}

func (s Service) GetReservationsByHotelID(ctx context.Context, hotelID string) ([]reservations.Reservation, error) {
	return s.repository.GetByHotelID(ctx, hotelID)
}
