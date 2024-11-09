// repositories/reservations/reservation_mongo.go
package reservations

import (
	"context"
	"fmt"
	"hotels-api/domain/reservations"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewMongo(client *mongo.Client, database, collection string) Mongo {
	return Mongo{client: client, database: database, collection: collection}
}

func (m Mongo) Create(ctx context.Context, reservation reservations.Reservation) (string, error) {
	result, err := m.client.Database(m.database).Collection(m.collection).InsertOne(ctx, reservation)
	if err != nil {
		return "", fmt.Errorf("error creating reservation: %w", err)
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (m Mongo) GetByUserID(ctx context.Context, userID string) ([]reservations.Reservation, error) {
	cursor, err := m.client.Database(m.database).Collection(m.collection).Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, fmt.Errorf("error getting reservations by user ID: %w", err)
	}

	var reservations []reservations.Reservation
	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, fmt.Errorf("error decoding reservations: %w", err)
	}
	return reservations, nil
}
