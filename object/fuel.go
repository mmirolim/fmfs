package object

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Fuel object for vehicle
// Collection where Document fuel stored is defined by fleet UUID
// like collection name fuel_fleet-uuid
type Fuel struct {
	ID        bson.ObjectId `bson:"_id"` // _id of document in db
	Vehicle   string        // vehicle uuid
	Fleet     string        // fleet uuid
	FuelUnit                // unit measurment unit Litres, Gallon
	Amount    int           // Number of units
	TankSize  int           // fuel tank size in vehicle
	Info      string        // extra information
	TotalCost float32       // total cost of fuel
	Currency  string        // currency of cost
	FillDate  time.Time     // when fuel filled
	Partial   bool          // is it partial or fuel tank filled
	Vendor    string        // which gas station Lukoil, sibneft
	Geo                     // geo location of fuel filled
	Base                    // system fields who, when created, updated and deleted
	FuelType                // fuel type Gas, Diesel, Gasoline
}

// define custom type for Fuel Types
type FuelType string

// define units
// @todo temp solution actually it should be
// taken from some service with all possible
// liquid measuring units
type FuelUnit string

const (
	DIESEL      = FuelType("DIESEL")
	GASOLINE    = FuelType("GASOLINE")
	GAS         = FuelType("GAS")
	ELECTRICITY = FuelType("ELECTRICITY")

	LITRES  = FuelUnit("LITRES")
	GALLONS = FuelUnit("GALLONS")
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
	return "fuel"
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
			Key: []string{"filldate"}, // index date of fuel filled for reports and daily graphs
		},
		mgo.Index{
			Key: []string{"$2dsphere:loc"}, // add index for location search
		},
		mgo.Index{
			Key: []string{"deletedby"}, // index deleted by, all default queries will use this key because of soft delete
		},
	}
}
