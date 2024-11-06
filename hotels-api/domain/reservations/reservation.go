package reservations

type Reservation struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	HotelID   string `json:"hotel_id" bson:"hotel_id"`
	UserID    string `json:"user_id" bson:"user_id"`
	StartDate string `json:"start_date" bson:"start_date"`
	EndDate   string `json:"end_date" bson:"end_date"`
	Status    string `json:"status" bson:"status"`
}
