package object

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Fuel object for vehicle
// Collection where Document fuel stored is defined by fleet UUID
// like collection name fuel_fleet-uuid
type Fuel struct {
	ID        bson.ObjectId // id of document in db
	Vehicle   string        // vehicle uuid
	Fleet     string        // fleet uuid
	Unit      string        // unit measurment unit L, Gallon
	Amount    int           // Number of units
	TankSize  int           // fuel tank size in vehicle
	Info      string        // extra information
	TotalCost float32       // total cost of fuel
	Currency  string        // currency of cost
	Date      time.Time     // when fuel filled
	Partial   bool          // is it partial or fuel tank filled
	Vendor    string        // which gas station Lukoil, sibneft
	Geo                     // geo location of fuel filled
	Base                    // system fields who, when created, updated and deleted
	FuelType                // fuel type Gas, Diesel, Gasoline
}

// define custom type for Fuel Types
type FuelType string

const (
	DIESEL      = FuelType("DIESEL")
	GASOLINE    = FuelType("GASOLINE")
	GAS         = FuelType("GAS")
	ELECTRICITY = FuelType("ELECTRICITY")
)

var (
	ErrNoFleet = errors.New("no fleet set")
)

// get Geo property of fuel
func (f Fuel) Location() Geo {
	return f.Geo
}

// get entry collection name
// in form fuel_fleet-uuid if fleet missing return error
func (f Fuel) Collection() string {
	// check fleet should not be empty
	if f.Fleet == "" {
		return f.Fleet
	}
	return fmt.Sprintf("%s_%s", "fuel", f.Fleet)
}

// set ObjectID
func (f *Fuel) SetID(oid bson.ObjectId) {
	f.ID = oid
}

// @todo do not use mgo Index type ot not be depended
// set field which should be indexed
func (f Fuel) Index() []mgo.Index {
	return []mgo.Index{
		mgo.Index{
			Key: []string{"vehicle"}, // many queries maybe for particular vehicle so it is better to index this field
		},
		mgo.Index{
			Key: []string{"fleet"}, // actually fuel entries fill be in collection with fleet namespace so it is not required
		},
		mgo.Index{
			Key: []string{"date"}, // index date of fuel filled for reports and daily graphs
		},
		mgo.Index{
			Key: []string{"$2dsphere:loc"}, // add index for location search
		},
		mgo.Index{
			Key: []string{"deletedby"}, // index deleted by, all default queries will use this key because of soft delete
		},
	}
}
