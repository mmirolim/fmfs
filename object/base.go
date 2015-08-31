package object

import "time"

//Base technical data required to manage lifecycle
type Base struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	DeletedBy string // UUID of user
	CreatedBy string // UUID of user
	UpdatedBy string // UUID of user
}

//Geo location struct
type Geo struct {
	Type        string    // type of geo location for mongo
	Coordinates []float32 // store lat and long properties
}

// set who created and when
func (b *Base) Created(by string, at ...time.Time) {
	b.CreatedBy = by
	if len(at) == 1 {
		b.CreatedAt = at[0]
	} else {
		b.CreatedAt = time.Now()
	}
}

// set who updated and when
func (b *Base) Updated(by string, at ...time.Time) {
	b.UpdatedBy = by
	if len(at) == 1 {
		b.UpdatedAt = at[0]
	} else {
		b.UpdatedAt = time.Now()
	}
}

// set who deleted and when
func (b *Base) Deleted(by string, at ...time.Time) {
	b.DeletedBy = by
	if len(at) == 1 {
		b.DeletedAt = at[0]
	} else {
		b.DeletedAt = time.Now()
	}
}
