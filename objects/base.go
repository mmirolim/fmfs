package object

import "time"

//Base technical data required to manage lifecycle
//@TODO add index
type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	DeletedBy string // UUID of user
	CreatedBy string // UUID of user
	UpdatedBy string // UUID of user
}

//Geo location struct
//@TODO add index
type Geo struct {
	Type        string    // type of geo location for mongo
	Coordinates []float32 // store lat and long properties
}
