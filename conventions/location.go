package conventions

// Location is a domain type of location
type Location struct {
	Latitude       float64 `json:"latitude" validate:"required,gte=-90.0,lte=90.0"`
	Longitude      float64 `json:"longitude" validate:"required,gte=-180.0,lte=180.0"`
	LocationTypeID int8    `json:"location_type_id" validate:"omitempty,oneof=0 1"`
}

// NewLocation is a constructor of Location
func NewLocation(latitude, longitude float64, locationTypeID int8) Location {
	return Location{
		Latitude:       latitude,
		Longitude:      longitude,
		LocationTypeID: locationTypeID,
	}
}
