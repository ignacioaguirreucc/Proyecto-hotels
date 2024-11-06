package reservations

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Reservation struct {
	ID        string `bson:"_id,omitempty"`
	HotelID   string `bson:"hotel_id"`
	UserID    string `bson:"user_id"`
	StartDate string `bson:"start_date"`
	EndDate   string `bson:"end_date"`
	Status    string `bson:"status"`
}

type MongoConfig struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

type Mongo struct {
	collection *mongo.Collection
}

func NewMongo(config MongoConfig) Mongo {
	return Mongo{
		collection: config.Client.Database(config.Database).Collection(config.Collection),
	}
}

// CRUD básico
func (m Mongo) CreateReservation(ctx context.Context, reservation Reservation) (string, error) {
	result, err := m.collection.InsertOne(ctx, reservation)
	if err != nil {
		return "", fmt.Errorf("error creating reservation: %w", err)
	}
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (m Mongo) GetReservationsByHotelID(ctx context.Context, hotelID string) ([]Reservation, error) {
	cursor, err := m.collection.Find(ctx, bson.M{"hotel_id": hotelID})
	if err != nil {
		return nil, fmt.Errorf("error finding reservations: %w", err)
	}
	defer cursor.Close(ctx)

	var reservations []Reservation
	for cursor.Next(ctx) {
		var res Reservation
		if err := cursor.Decode(&res); err != nil {
			return nil, fmt.Errorf("error decoding reservation: %w", err)
		}
		reservations = append(reservations, res)
	}
	return reservations, nil
}
