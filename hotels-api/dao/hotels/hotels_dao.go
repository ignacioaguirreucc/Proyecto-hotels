package hotels

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Address   string             `bson:"address"`
	City      string             `bson:"city"`
	State     string             `bson:"state"`
	Rating    float64            `bson:"rating"`
	Amenities []string           `bson:"amenities"`
}
