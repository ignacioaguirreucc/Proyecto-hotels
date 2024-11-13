package hotels

type Hotel struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Address     string   `json:"address"`
	City        string   `json:"city"`
	State       string   `json:"state"`
	Rating      float64  `json:"rating"`
	Amenities   []string `json:"amenities"`
	Descripcion []string `json:"descripcion"`
}

type HotelNew struct {
	Operation string `json:"operation"`
	HotelID   string `json:"hotel_id"`
}
