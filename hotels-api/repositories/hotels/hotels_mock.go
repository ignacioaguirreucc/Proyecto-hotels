package hotels

import (
	"context"
	"fmt"
	hotelsDAO "hotels-api/dao/hotels"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mock struct {
	docs map[string]hotelsDAO.Hotel
}

func NewMock() Mock {
	return Mock{
		docs: make(map[string]hotelsDAO.Hotel),
	}
}

func (repository Mock) GetHotelByID(ctx context.Context, id string) (hotelsDAO.Hotel, error) {
	return repository.docs[id], nil
}

func (repository Mock) Create(ctx context.Context, hotel hotelsDAO.Hotel) (string, error) {
	hotelID := primitive.NewObjectID() // Genera un nuevo ObjectID
	hotel.ID = hotelID                 // Asigna el ObjectID generado al hotel

	repository.docs[hotelID.Hex()] = hotel // Guarda el hotel en el mapa usando el ObjectID en formato hexadecimal como clave
	return hotelID.Hex(), nil              // Retorna el ID como cadena en formato hexadecimal
}

func (repository Mock) Update(ctx context.Context, hotel hotelsDAO.Hotel) error {
	// Check if the hotel exists in the mock storage
	currentHotel, exists := repository.docs[hotel.ID.Hex()]
	if !exists {
		return fmt.Errorf("hotel with ID %s not found", hotel.ID)
	}

	// Update only the fields that are non-zero or non-empty
	if hotel.Name != "" {
		currentHotel.Name = hotel.Name
	}
	if hotel.Address != "" {
		currentHotel.Address = hotel.Address
	}
	if hotel.City != "" {
		currentHotel.City = hotel.City
	}
	if hotel.State != "" {
		currentHotel.State = hotel.State
	}

	if hotel.Rating != 0 {
		currentHotel.Rating = hotel.Rating
	}
	if len(hotel.Amenities) > 0 {
		currentHotel.Amenities = hotel.Amenities
	}
	if len(hotel.Descripcion) > 0 {
		currentHotel.Descripcion = hotel.Descripcion
	}
	// Save the updated hotel back to the mock storage
	repository.docs[hotel.ID.Hex()] = currentHotel
	return nil
}

func (repository Mock) Delete(ctx context.Context, id string) error {
	if _, exists := repository.docs[id]; !exists {
		return fmt.Errorf("hotel with ID %s not found", id)
	}
	// Remove the hotel from the mock storage
	delete(repository.docs, id)
	return nil
}
